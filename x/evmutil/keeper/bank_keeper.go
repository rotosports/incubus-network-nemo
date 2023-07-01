package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/incubus-network/nemo/x/evmutil/types"
)

const (
	// EvmDenom is the gas denom used by the evm
	EvmDenom = "anemo"

	// CosmosDenom is the gas denom used by the nemo app
	CosmosDenom = "unemo"
)

// ConversionMultiplier is the conversion multiplier between anemo and unemo
var ConversionMultiplier = sdkmath.NewInt(1_000_000_000_000)

var _ evmtypes.BankKeeper = EvmBankKeeper{}

// EvmBankKeeper is a BankKeeper wrapper for the x/evm module to allow the use
// of the 18 decimal anemo coin on the evm.
// x/evm consumes gas and send coins by minting and burning anemo coins in its module
// account and then sending the funds to the target account.
// This keeper uses both the unemo coin and a separate anemo balance to manage the
// extra percision needed by the evm.
type EvmBankKeeper struct {
	anemoKeeper Keeper
	bk          types.BankKeeper
	ak          types.AccountKeeper
}

func NewEvmBankKeeper(anemoKeeper Keeper, bk types.BankKeeper, ak types.AccountKeeper) EvmBankKeeper {
	return EvmBankKeeper{
		anemoKeeper: anemoKeeper,
		bk:          bk,
		ak:          ak,
	}
}

// GetBalance returns the total **spendable** balance of anemo for a given account by address.
func (k EvmBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	if denom != EvmDenom {
		panic(fmt.Errorf("only evm denom %s is supported by EvmBankKeeper", EvmDenom))
	}

	spendableCoins := k.bk.SpendableCoins(ctx, addr)
	unemo := spendableCoins.AmountOf(CosmosDenom)
	anemo := k.anemoKeeper.GetBalance(ctx, addr)
	total := unemo.Mul(ConversionMultiplier).Add(anemo)
	return sdk.NewCoin(EvmDenom, total)
}

// SendCoins transfers anemo coins from a AccAddress to an AccAddress.
func (k EvmBankKeeper) SendCoins(ctx sdk.Context, senderAddr sdk.AccAddress, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	// SendCoins method is not used by the evm module, but is required by the
	// evmtypes.BankKeeper interface. This must be updated if the evm module
	// is updated to use SendCoins.
	panic("not implemented")
}

// SendCoinsFromModuleToAccount transfers anemo coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist. An error is returned if the recipient
// address is black-listed or if sending the tokens fails.
func (k EvmBankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	unemo, anemo, err := SplitAnemoCoins(amt)
	if err != nil {
		return err
	}

	if unemo.Amount.IsPositive() {
		if err := k.bk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, sdk.NewCoins(unemo)); err != nil {
			return err
		}
	}

	senderAddr := k.GetModuleAddress(senderModule)
	if err := k.ConvertOneUnemoToAnemoIfNeeded(ctx, senderAddr, anemo); err != nil {
		return err
	}

	if err := k.anemoKeeper.SendBalance(ctx, senderAddr, recipientAddr, anemo); err != nil {
		return err
	}

	return k.ConvertAnemoToUnemo(ctx, recipientAddr)
}

// SendCoinsFromAccountToModule transfers anemo coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (k EvmBankKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	unemo, anemoNeeded, err := SplitAnemoCoins(amt)
	if err != nil {
		return err
	}

	if unemo.IsPositive() {
		if err := k.bk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, sdk.NewCoins(unemo)); err != nil {
			return err
		}
	}

	if err := k.ConvertOneUnemoToAnemoIfNeeded(ctx, senderAddr, anemoNeeded); err != nil {
		return err
	}

	recipientAddr := k.GetModuleAddress(recipientModule)
	if err := k.anemoKeeper.SendBalance(ctx, senderAddr, recipientAddr, anemoNeeded); err != nil {
		return err
	}

	return k.ConvertAnemoToUnemo(ctx, recipientAddr)
}

// MintCoins mints anemo coins by minting the equivalent unemo coins and any remaining anemo coins.
// It will panic if the module account does not exist or is unauthorized.
func (k EvmBankKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	unemo, anemo, err := SplitAnemoCoins(amt)
	if err != nil {
		return err
	}

	if unemo.IsPositive() {
		if err := k.bk.MintCoins(ctx, moduleName, sdk.NewCoins(unemo)); err != nil {
			return err
		}
	}

	recipientAddr := k.GetModuleAddress(moduleName)
	if err := k.anemoKeeper.AddBalance(ctx, recipientAddr, anemo); err != nil {
		return err
	}

	return k.ConvertAnemoToUnemo(ctx, recipientAddr)
}

// BurnCoins burns anemo coins by burning the equivalent unemo coins and any remaining anemo coins.
// It will panic if the module account does not exist or is unauthorized.
func (k EvmBankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	unemo, anemo, err := SplitAnemoCoins(amt)
	if err != nil {
		return err
	}

	if unemo.IsPositive() {
		if err := k.bk.BurnCoins(ctx, moduleName, sdk.NewCoins(unemo)); err != nil {
			return err
		}
	}

	moduleAddr := k.GetModuleAddress(moduleName)
	if err := k.ConvertOneUnemoToAnemoIfNeeded(ctx, moduleAddr, anemo); err != nil {
		return err
	}

	return k.anemoKeeper.RemoveBalance(ctx, moduleAddr, anemo)
}

// ConvertOneUnemoToAnemoIfNeeded converts 1 unemo to anemo for an address if
// its anemo balance is smaller than the anemoNeeded amount.
func (k EvmBankKeeper) ConvertOneUnemoToAnemoIfNeeded(ctx sdk.Context, addr sdk.AccAddress, anemoNeeded sdkmath.Int) error {
	anemoBal := k.anemoKeeper.GetBalance(ctx, addr)
	if anemoBal.GTE(anemoNeeded) {
		return nil
	}

	unemoToStore := sdk.NewCoins(sdk.NewCoin(CosmosDenom, sdk.OneInt()))
	if err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, unemoToStore); err != nil {
		return err
	}

	// add 1unemo equivalent of anemo to addr
	anemoToReceive := ConversionMultiplier
	if err := k.anemoKeeper.AddBalance(ctx, addr, anemoToReceive); err != nil {
		return err
	}

	return nil
}

// ConvertAnemoToUnemo converts all available anemo to unemo for a given AccAddress.
func (k EvmBankKeeper) ConvertAnemoToUnemo(ctx sdk.Context, addr sdk.AccAddress) error {
	totalAnemo := k.anemoKeeper.GetBalance(ctx, addr)
	unemo, _, err := SplitAnemoCoins(sdk.NewCoins(sdk.NewCoin(EvmDenom, totalAnemo)))
	if err != nil {
		return err
	}

	// do nothing if account does not have enough anemo for a single unemo
	unemoToReceive := unemo.Amount
	if !unemoToReceive.IsPositive() {
		return nil
	}

	// remove anemo used for converting to unemo
	anemoToBurn := unemoToReceive.Mul(ConversionMultiplier)
	finalBal := totalAnemo.Sub(anemoToBurn)
	if err := k.anemoKeeper.SetBalance(ctx, addr, finalBal); err != nil {
		return err
	}

	fromAddr := k.GetModuleAddress(types.ModuleName)
	if err := k.bk.SendCoins(ctx, fromAddr, addr, sdk.NewCoins(unemo)); err != nil {
		return err
	}

	return nil
}

func (k EvmBankKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	addr := k.ak.GetModuleAddress(moduleName)
	if addr == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}
	return addr
}

// SplitAnemoCoins splits anemo coins to the equivalent unemo coins and any remaining anemo balance.
// An error will be returned if the coins are not valid or if the coins are not the anemo denom.
func SplitAnemoCoins(coins sdk.Coins) (sdk.Coin, sdkmath.Int, error) {
	anemo := sdk.ZeroInt()
	unemo := sdk.NewCoin(CosmosDenom, sdk.ZeroInt())

	if len(coins) == 0 {
		return unemo, anemo, nil
	}

	if err := ValidateEvmCoins(coins); err != nil {
		return unemo, anemo, err
	}

	// note: we should always have len(coins) == 1 here since coins cannot have dup denoms after we validate.
	coin := coins[0]
	remainingBalance := coin.Amount.Mod(ConversionMultiplier)
	if remainingBalance.IsPositive() {
		anemo = remainingBalance
	}
	unemoAmount := coin.Amount.Quo(ConversionMultiplier)
	if unemoAmount.IsPositive() {
		unemo = sdk.NewCoin(CosmosDenom, unemoAmount)
	}

	return unemo, anemo, nil
}

// ValidateEvmCoins validates the coins from evm is valid and is the EvmDenom (anemo).
func ValidateEvmCoins(coins sdk.Coins) error {
	if len(coins) == 0 {
		return nil
	}

	// validate that coins are non-negative, sorted, and no dup denoms
	if err := coins.Validate(); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, coins.String())
	}

	// validate that coin denom is anemo
	if len(coins) != 1 || coins[0].Denom != EvmDenom {
		errMsg := fmt.Sprintf("invalid evm coin denom, only %s is supported", EvmDenom)
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, errMsg)
	}

	return nil
}
