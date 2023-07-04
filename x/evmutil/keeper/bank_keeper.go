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
	EvmDenom = "atfury"

	// CosmosDenom is the gas denom used by the nemo app
	CosmosDenom = "unemo"
)

// ConversionMultiplier is the conversion multiplier between atfury and unemo
var ConversionMultiplier = sdkmath.NewInt(1_000_000_000_000)

var _ evmtypes.BankKeeper = EvmBankKeeper{}

// EvmBankKeeper is a BankKeeper wrapper for the x/evm module to allow the use
// of the 18 decimal atfury coin on the evm.
// x/evm consumes gas and send coins by minting and burning atfury coins in its module
// account and then sending the funds to the target account.
// This keeper uses both the unemo coin and a separate atfury balance to manage the
// extra percision needed by the evm.
type EvmBankKeeper struct {
	atfuryKeeper Keeper
	bk          types.BankKeeper
	ak          types.AccountKeeper
}

func NewEvmBankKeeper(atfuryKeeper Keeper, bk types.BankKeeper, ak types.AccountKeeper) EvmBankKeeper {
	return EvmBankKeeper{
		atfuryKeeper: atfuryKeeper,
		bk:          bk,
		ak:          ak,
	}
}

// GetBalance returns the total **spendable** balance of atfury for a given account by address.
func (k EvmBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	if denom != EvmDenom {
		panic(fmt.Errorf("only evm denom %s is supported by EvmBankKeeper", EvmDenom))
	}

	spendableCoins := k.bk.SpendableCoins(ctx, addr)
	unemo := spendableCoins.AmountOf(CosmosDenom)
	atfury := k.atfuryKeeper.GetBalance(ctx, addr)
	total := unemo.Mul(ConversionMultiplier).Add(atfury)
	return sdk.NewCoin(EvmDenom, total)
}

// SendCoins transfers atfury coins from a AccAddress to an AccAddress.
func (k EvmBankKeeper) SendCoins(ctx sdk.Context, senderAddr sdk.AccAddress, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	// SendCoins method is not used by the evm module, but is required by the
	// evmtypes.BankKeeper interface. This must be updated if the evm module
	// is updated to use SendCoins.
	panic("not implemented")
}

// SendCoinsFromModuleToAccount transfers atfury coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist. An error is returned if the recipient
// address is black-listed or if sending the tokens fails.
func (k EvmBankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	unemo, atfury, err := SplitAtfuryCoins(amt)
	if err != nil {
		return err
	}

	if unemo.Amount.IsPositive() {
		if err := k.bk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, sdk.NewCoins(unemo)); err != nil {
			return err
		}
	}

	senderAddr := k.GetModuleAddress(senderModule)
	if err := k.ConvertOneUnemoToAtfuryIfNeeded(ctx, senderAddr, atfury); err != nil {
		return err
	}

	if err := k.atfuryKeeper.SendBalance(ctx, senderAddr, recipientAddr, atfury); err != nil {
		return err
	}

	return k.ConvertAtfuryToUnemo(ctx, recipientAddr)
}

// SendCoinsFromAccountToModule transfers atfury coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (k EvmBankKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	unemo, atfuryNeeded, err := SplitAtfuryCoins(amt)
	if err != nil {
		return err
	}

	if unemo.IsPositive() {
		if err := k.bk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, sdk.NewCoins(unemo)); err != nil {
			return err
		}
	}

	if err := k.ConvertOneUnemoToAtfuryIfNeeded(ctx, senderAddr, atfuryNeeded); err != nil {
		return err
	}

	recipientAddr := k.GetModuleAddress(recipientModule)
	if err := k.atfuryKeeper.SendBalance(ctx, senderAddr, recipientAddr, atfuryNeeded); err != nil {
		return err
	}

	return k.ConvertAtfuryToUnemo(ctx, recipientAddr)
}

// MintCoins mints atfury coins by minting the equivalent unemo coins and any remaining atfury coins.
// It will panic if the module account does not exist or is unauthorized.
func (k EvmBankKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	unemo, atfury, err := SplitAtfuryCoins(amt)
	if err != nil {
		return err
	}

	if unemo.IsPositive() {
		if err := k.bk.MintCoins(ctx, moduleName, sdk.NewCoins(unemo)); err != nil {
			return err
		}
	}

	recipientAddr := k.GetModuleAddress(moduleName)
	if err := k.atfuryKeeper.AddBalance(ctx, recipientAddr, atfury); err != nil {
		return err
	}

	return k.ConvertAtfuryToUnemo(ctx, recipientAddr)
}

// BurnCoins burns atfury coins by burning the equivalent unemo coins and any remaining atfury coins.
// It will panic if the module account does not exist or is unauthorized.
func (k EvmBankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	unemo, atfury, err := SplitAtfuryCoins(amt)
	if err != nil {
		return err
	}

	if unemo.IsPositive() {
		if err := k.bk.BurnCoins(ctx, moduleName, sdk.NewCoins(unemo)); err != nil {
			return err
		}
	}

	moduleAddr := k.GetModuleAddress(moduleName)
	if err := k.ConvertOneUnemoToAtfuryIfNeeded(ctx, moduleAddr, atfury); err != nil {
		return err
	}

	return k.atfuryKeeper.RemoveBalance(ctx, moduleAddr, atfury)
}

// ConvertOneUnemoToAtfuryIfNeeded converts 1 unemo to atfury for an address if
// its atfury balance is smaller than the atfuryNeeded amount.
func (k EvmBankKeeper) ConvertOneUnemoToAtfuryIfNeeded(ctx sdk.Context, addr sdk.AccAddress, atfuryNeeded sdkmath.Int) error {
	atfuryBal := k.atfuryKeeper.GetBalance(ctx, addr)
	if atfuryBal.GTE(atfuryNeeded) {
		return nil
	}

	unemoToStore := sdk.NewCoins(sdk.NewCoin(CosmosDenom, sdk.OneInt()))
	if err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, unemoToStore); err != nil {
		return err
	}

	// add 1unemo equivalent of atfury to addr
	atfuryToReceive := ConversionMultiplier
	if err := k.atfuryKeeper.AddBalance(ctx, addr, atfuryToReceive); err != nil {
		return err
	}

	return nil
}

// ConvertAtfuryToUnemo converts all available atfury to unemo for a given AccAddress.
func (k EvmBankKeeper) ConvertAtfuryToUnemo(ctx sdk.Context, addr sdk.AccAddress) error {
	totalAtfury := k.atfuryKeeper.GetBalance(ctx, addr)
	unemo, _, err := SplitAtfuryCoins(sdk.NewCoins(sdk.NewCoin(EvmDenom, totalAtfury)))
	if err != nil {
		return err
	}

	// do nothing if account does not have enough atfury for a single unemo
	unemoToReceive := unemo.Amount
	if !unemoToReceive.IsPositive() {
		return nil
	}

	// remove atfury used for converting to unemo
	atfuryToBurn := unemoToReceive.Mul(ConversionMultiplier)
	finalBal := totalAtfury.Sub(atfuryToBurn)
	if err := k.atfuryKeeper.SetBalance(ctx, addr, finalBal); err != nil {
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

// SplitAtfuryCoins splits atfury coins to the equivalent unemo coins and any remaining atfury balance.
// An error will be returned if the coins are not valid or if the coins are not the atfury denom.
func SplitAtfuryCoins(coins sdk.Coins) (sdk.Coin, sdkmath.Int, error) {
	atfury := sdk.ZeroInt()
	unemo := sdk.NewCoin(CosmosDenom, sdk.ZeroInt())

	if len(coins) == 0 {
		return unemo, atfury, nil
	}

	if err := ValidateEvmCoins(coins); err != nil {
		return unemo, atfury, err
	}

	// note: we should always have len(coins) == 1 here since coins cannot have dup denoms after we validate.
	coin := coins[0]
	remainingBalance := coin.Amount.Mod(ConversionMultiplier)
	if remainingBalance.IsPositive() {
		atfury = remainingBalance
	}
	unemoAmount := coin.Amount.Quo(ConversionMultiplier)
	if unemoAmount.IsPositive() {
		unemo = sdk.NewCoin(CosmosDenom, unemoAmount)
	}

	return unemo, atfury, nil
}

// ValidateEvmCoins validates the coins from evm is valid and is the EvmDenom (atfury).
func ValidateEvmCoins(coins sdk.Coins) error {
	if len(coins) == 0 {
		return nil
	}

	// validate that coins are non-negative, sorted, and no dup denoms
	if err := coins.Validate(); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, coins.String())
	}

	// validate that coin denom is atfury
	if len(coins) != 1 || coins[0].Denom != EvmDenom {
		errMsg := fmt.Sprintf("invalid evm coin denom, only %s is supported", EvmDenom)
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, errMsg)
	}

	return nil
}
