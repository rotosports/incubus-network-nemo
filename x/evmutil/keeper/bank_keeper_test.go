package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmtime "github.com/tendermint/tendermint/types/time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/incubus-network/nemo/x/evmutil/keeper"
	"github.com/incubus-network/nemo/x/evmutil/testutil"
	"github.com/incubus-network/nemo/x/evmutil/types"
)

type evmBankKeeperTestSuite struct {
	testutil.Suite
}

func (suite *evmBankKeeperTestSuite) SetupTest() {
	suite.Suite.SetupTest()
}

func (suite *evmBankKeeperTestSuite) TestGetBalance_ReturnsSpendable() {
	startingCoins := sdk.NewCoins(sdk.NewInt64Coin("unemo", 10))
	startingAnemo := sdkmath.NewInt(100)

	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	bacc := authtypes.NewBaseAccountWithAddress(suite.Addrs[0])
	vacc := vesting.NewContinuousVestingAccount(bacc, startingCoins, now.Unix(), endTime.Unix())
	suite.AccountKeeper.SetAccount(suite.Ctx, vacc)

	err := suite.App.FundAccount(suite.Ctx, suite.Addrs[0], startingCoins)
	suite.Require().NoError(err)
	err = suite.Keeper.SetBalance(suite.Ctx, suite.Addrs[0], startingAnemo)
	suite.Require().NoError(err)

	coin := suite.EvmBankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "anemo")
	suite.Require().Equal(startingAnemo, coin.Amount)

	ctx := suite.Ctx.WithBlockTime(now.Add(12 * time.Hour))
	coin = suite.EvmBankKeeper.GetBalance(ctx, suite.Addrs[0], "anemo")
	suite.Require().Equal(sdkmath.NewIntFromUint64(5_000_000_000_100), coin.Amount)
}

func (suite *evmBankKeeperTestSuite) TestGetBalance_NotEvmDenom() {
	suite.Require().Panics(func() {
		suite.EvmBankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "unemo")
	})
	suite.Require().Panics(func() {
		suite.EvmBankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "busd")
	})
}

func (suite *evmBankKeeperTestSuite) TestGetBalance() {
	tests := []struct {
		name           string
		startingAmount sdk.Coins
		expAmount      sdkmath.Int
	}{
		{
			"unemo with anemo",
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 100),
				sdk.NewInt64Coin("unemo", 10),
			),
			sdkmath.NewInt(10_000_000_000_100),
		},
		{
			"just anemo",
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 100),
				sdk.NewInt64Coin("busd", 100),
			),
			sdkmath.NewInt(100),
		},
		{
			"just unemo",
			sdk.NewCoins(
				sdk.NewInt64Coin("unemo", 10),
				sdk.NewInt64Coin("busd", 100),
			),
			sdkmath.NewInt(10_000_000_000_000),
		},
		{
			"no unemo or anemo",
			sdk.NewCoins(),
			sdk.ZeroInt(),
		},
		{
			"with avaka that is more than 1 unemo",
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 20_000_000_000_220),
				sdk.NewInt64Coin("unemo", 11),
			),
			sdkmath.NewInt(31_000_000_000_220),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.FundAccountWithNemo(suite.Addrs[0], tt.startingAmount)
			coin := suite.EvmBankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "anemo")
			suite.Require().Equal(tt.expAmount, coin.Amount)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestSendCoinsFromModuleToAccount() {
	startingModuleCoins := sdk.NewCoins(
		sdk.NewInt64Coin("anemo", 200),
		sdk.NewInt64Coin("unemo", 100),
	)
	tests := []struct {
		name           string
		sendCoins      sdk.Coins
		startingAccBal sdk.Coins
		expAccBal      sdk.Coins
		hasErr         bool
	}{
		{
			"send more than 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 12_000_000_000_010)),
			sdk.Coins{},
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 10),
				sdk.NewInt64Coin("unemo", 12),
			),
			false,
		},
		{
			"send less than 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 122)),
			sdk.Coins{},
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 122),
				sdk.NewInt64Coin("unemo", 0),
			),
			false,
		},
		{
			"send an exact amount of unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 98_000_000_000_000)),
			sdk.Coins{},
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 0o0),
				sdk.NewInt64Coin("unemo", 98),
			),
			false,
		},
		{
			"send no anemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 0)),
			sdk.Coins{},
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 0),
				sdk.NewInt64Coin("unemo", 0),
			),
			false,
		},
		{
			"errors if sending other coins",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 500), sdk.NewInt64Coin("busd", 1000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if not enough total anemo to cover",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100_000_000_001_000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if not enough unemo to cover",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 200_000_000_000_000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"converts receiver's anemo to unemo if there's enough anemo after the transfer",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 99_000_000_000_200)),
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 999_999_999_900),
				sdk.NewInt64Coin("unemo", 1),
			),
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 100),
				sdk.NewInt64Coin("unemo", 101),
			),
			false,
		},
		{
			"converts all of receiver's anemo to unemo even if somehow receiver has more than 1unemo of anemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 12_000_000_000_100)),
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 5_999_999_999_990),
				sdk.NewInt64Coin("unemo", 1),
			),
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 90),
				sdk.NewInt64Coin("unemo", 19),
			),
			false,
		},
		{
			"swap 1 unemo for anemo if module account doesn't have enough anemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 99_000_000_001_000)),
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 200),
				sdk.NewInt64Coin("unemo", 1),
			),
			sdk.NewCoins(
				sdk.NewInt64Coin("anemo", 1200),
				sdk.NewInt64Coin("unemo", 100),
			),
			false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.FundAccountWithNemo(suite.Addrs[0], tt.startingAccBal)
			suite.FundModuleAccountWithNemo(evmtypes.ModuleName, startingModuleCoins)

			// fund our module with some unemo to account for converting extra anemo back to unemo
			suite.FundModuleAccountWithNemo(types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("unemo", 10)))

			err := suite.EvmBankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, evmtypes.ModuleName, suite.Addrs[0], tt.sendCoins)
			if tt.hasErr {
				suite.Require().Error(err)
				return
			} else {
				suite.Require().NoError(err)
			}

			// check unemo
			unemoSender := suite.BankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "unemo")
			suite.Require().Equal(tt.expAccBal.AmountOf("unemo").Int64(), unemoSender.Amount.Int64())

			// check anemo
			actualAnemo := suite.Keeper.GetBalance(suite.Ctx, suite.Addrs[0])
			suite.Require().Equal(tt.expAccBal.AmountOf("anemo").Int64(), actualAnemo.Int64())
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestSendCoinsFromAccountToModule() {
	startingAccCoins := sdk.NewCoins(
		sdk.NewInt64Coin("anemo", 200),
		sdk.NewInt64Coin("unemo", 100),
	)
	startingModuleCoins := sdk.NewCoins(
		sdk.NewInt64Coin("anemo", 100_000_000_000),
	)
	tests := []struct {
		name           string
		sendCoins      sdk.Coins
		expSenderCoins sdk.Coins
		expModuleCoins sdk.Coins
		hasErr         bool
	}{
		{
			"send more than 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 12_000_000_000_010)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 190), sdk.NewInt64Coin("unemo", 88)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100_000_000_010), sdk.NewInt64Coin("unemo", 12)),
			false,
		},
		{
			"send less than 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 122)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 78), sdk.NewInt64Coin("unemo", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100_000_000_122), sdk.NewInt64Coin("unemo", 0)),
			false,
		},
		{
			"send an exact amount of unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 98_000_000_000_000)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 200), sdk.NewInt64Coin("unemo", 2)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100_000_000_000), sdk.NewInt64Coin("unemo", 98)),
			false,
		},
		{
			"send no anemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 0)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 200), sdk.NewInt64Coin("unemo", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100_000_000_000), sdk.NewInt64Coin("unemo", 0)),
			false,
		},
		{
			"errors if sending other coins",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 500), sdk.NewInt64Coin("busd", 1000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if have dup coins",
			sdk.Coins{
				sdk.NewInt64Coin("anemo", 12_000_000_000_000),
				sdk.NewInt64Coin("anemo", 2_000_000_000_000),
			},
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if not enough total anemo to cover",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100_000_000_001_000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if not enough unemo to cover",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 200_000_000_000_000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"converts 1 unemo to anemo if not enough anemo to cover",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 99_001_000_000_000)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 999_000_000_200), sdk.NewInt64Coin("unemo", 0)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 101_000_000_000), sdk.NewInt64Coin("unemo", 99)),
			false,
		},
		{
			"converts receiver's anemo to unemo if there's enough anemo after the transfer",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 5_900_000_000_200)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100_000_000_000), sdk.NewInt64Coin("unemo", 94)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 200), sdk.NewInt64Coin("unemo", 6)),
			false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.FundAccountWithNemo(suite.Addrs[0], startingAccCoins)
			suite.FundModuleAccountWithNemo(evmtypes.ModuleName, startingModuleCoins)

			err := suite.EvmBankKeeper.SendCoinsFromAccountToModule(suite.Ctx, suite.Addrs[0], evmtypes.ModuleName, tt.sendCoins)
			if tt.hasErr {
				suite.Require().Error(err)
				return
			} else {
				suite.Require().NoError(err)
			}

			// check sender balance
			unemoSender := suite.BankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "unemo")
			suite.Require().Equal(tt.expSenderCoins.AmountOf("unemo").Int64(), unemoSender.Amount.Int64())
			actualAnemo := suite.Keeper.GetBalance(suite.Ctx, suite.Addrs[0])
			suite.Require().Equal(tt.expSenderCoins.AmountOf("anemo").Int64(), actualAnemo.Int64())

			// check module balance
			moduleAddr := suite.AccountKeeper.GetModuleAddress(evmtypes.ModuleName)
			unemoSender = suite.BankKeeper.GetBalance(suite.Ctx, moduleAddr, "unemo")
			suite.Require().Equal(tt.expModuleCoins.AmountOf("unemo").Int64(), unemoSender.Amount.Int64())
			actualAnemo = suite.Keeper.GetBalance(suite.Ctx, moduleAddr)
			suite.Require().Equal(tt.expModuleCoins.AmountOf("anemo").Int64(), actualAnemo.Int64())
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestBurnCoins() {
	startingUnemo := sdkmath.NewInt(100)
	tests := []struct {
		name       string
		burnCoins  sdk.Coins
		expUnemo   sdkmath.Int
		expAnemo   sdkmath.Int
		hasErr     bool
		anemoStart sdkmath.Int
	}{
		{
			"burn more than 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 12_021_000_000_002)),
			sdkmath.NewInt(88),
			sdkmath.NewInt(100_000_000_000),
			false,
			sdkmath.NewInt(121_000_000_002),
		},
		{
			"burn less than 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 122)),
			sdkmath.NewInt(100),
			sdkmath.NewInt(878),
			false,
			sdkmath.NewInt(1000),
		},
		{
			"burn an exact amount of unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 98_000_000_000_000)),
			sdkmath.NewInt(2),
			sdkmath.NewInt(10),
			false,
			sdkmath.NewInt(10),
		},
		{
			"burn no anemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 0)),
			startingUnemo,
			sdk.ZeroInt(),
			false,
			sdk.ZeroInt(),
		},
		{
			"errors if burning other coins",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 500), sdk.NewInt64Coin("busd", 1000)),
			startingUnemo,
			sdkmath.NewInt(100),
			true,
			sdkmath.NewInt(100),
		},
		{
			"errors if have dup coins",
			sdk.Coins{
				sdk.NewInt64Coin("anemo", 12_000_000_000_000),
				sdk.NewInt64Coin("anemo", 2_000_000_000_000),
			},
			startingUnemo,
			sdk.ZeroInt(),
			true,
			sdk.ZeroInt(),
		},
		{
			"errors if burn amount is negative",
			sdk.Coins{sdk.Coin{Denom: "anemo", Amount: sdkmath.NewInt(-100)}},
			startingUnemo,
			sdkmath.NewInt(50),
			true,
			sdkmath.NewInt(50),
		},
		{
			"errors if not enough anemo to cover burn",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100_999_000_000_000)),
			sdkmath.NewInt(0),
			sdkmath.NewInt(99_000_000_000),
			true,
			sdkmath.NewInt(99_000_000_000),
		},
		{
			"errors if not enough unemo to cover burn",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 200_000_000_000_000)),
			sdkmath.NewInt(100),
			sdk.ZeroInt(),
			true,
			sdk.ZeroInt(),
		},
		{
			"converts 1 unemo to anemo if not enough anemo to cover",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 12_021_000_000_002)),
			sdkmath.NewInt(87),
			sdkmath.NewInt(980_000_000_000),
			false,
			sdkmath.NewInt(1_000_000_002),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			startingCoins := sdk.NewCoins(
				sdk.NewCoin("unemo", startingUnemo),
				sdk.NewCoin("anemo", tt.anemoStart),
			)
			suite.FundModuleAccountWithNemo(evmtypes.ModuleName, startingCoins)

			err := suite.EvmBankKeeper.BurnCoins(suite.Ctx, evmtypes.ModuleName, tt.burnCoins)
			if tt.hasErr {
				suite.Require().Error(err)
				return
			} else {
				suite.Require().NoError(err)
			}

			// check unemo
			unemoActual := suite.BankKeeper.GetBalance(suite.Ctx, suite.EvmModuleAddr, "unemo")
			suite.Require().Equal(tt.expUnemo, unemoActual.Amount)

			// check anemo
			anemoActual := suite.Keeper.GetBalance(suite.Ctx, suite.EvmModuleAddr)
			suite.Require().Equal(tt.expAnemo, anemoActual)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestMintCoins() {
	tests := []struct {
		name       string
		mintCoins  sdk.Coins
		unemo      sdkmath.Int
		anemo      sdkmath.Int
		hasErr     bool
		anemoStart sdkmath.Int
	}{
		{
			"mint more than 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 12_021_000_000_002)),
			sdkmath.NewInt(12),
			sdkmath.NewInt(21_000_000_002),
			false,
			sdk.ZeroInt(),
		},
		{
			"mint less than 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 901_000_000_001)),
			sdk.ZeroInt(),
			sdkmath.NewInt(901_000_000_001),
			false,
			sdk.ZeroInt(),
		},
		{
			"mint an exact amount of unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 123_000_000_000_000_000)),
			sdkmath.NewInt(123_000),
			sdk.ZeroInt(),
			false,
			sdk.ZeroInt(),
		},
		{
			"mint no anemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 0)),
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			false,
			sdk.ZeroInt(),
		},
		{
			"errors if minting other coins",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 500), sdk.NewInt64Coin("busd", 1000)),
			sdk.ZeroInt(),
			sdkmath.NewInt(100),
			true,
			sdkmath.NewInt(100),
		},
		{
			"errors if have dup coins",
			sdk.Coins{
				sdk.NewInt64Coin("anemo", 12_000_000_000_000),
				sdk.NewInt64Coin("anemo", 2_000_000_000_000),
			},
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			true,
			sdk.ZeroInt(),
		},
		{
			"errors if mint amount is negative",
			sdk.Coins{sdk.Coin{Denom: "anemo", Amount: sdkmath.NewInt(-100)}},
			sdk.ZeroInt(),
			sdkmath.NewInt(50),
			true,
			sdkmath.NewInt(50),
		},
		{
			"adds to existing anemo balance",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 12_021_000_000_002)),
			sdkmath.NewInt(12),
			sdkmath.NewInt(21_000_000_102),
			false,
			sdkmath.NewInt(100),
		},
		{
			"convert anemo balance to unemo if it exceeds 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 10_999_000_000_000)),
			sdkmath.NewInt(12),
			sdkmath.NewInt(1_200_000_001),
			false,
			sdkmath.NewInt(1_002_200_000_001),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.FundModuleAccountWithNemo(types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("unemo", 10)))
			suite.FundModuleAccountWithNemo(evmtypes.ModuleName, sdk.NewCoins(sdk.NewCoin("anemo", tt.anemoStart)))

			err := suite.EvmBankKeeper.MintCoins(suite.Ctx, evmtypes.ModuleName, tt.mintCoins)
			if tt.hasErr {
				suite.Require().Error(err)
				return
			} else {
				suite.Require().NoError(err)
			}

			// check unemo
			unemoActual := suite.BankKeeper.GetBalance(suite.Ctx, suite.EvmModuleAddr, "unemo")
			suite.Require().Equal(tt.unemo, unemoActual.Amount)

			// check anemo
			anemoActual := suite.Keeper.GetBalance(suite.Ctx, suite.EvmModuleAddr)
			suite.Require().Equal(tt.anemo, anemoActual)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestValidateEvmCoins() {
	tests := []struct {
		name      string
		coins     sdk.Coins
		shouldErr bool
	}{
		{
			"valid coins",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 500)),
			false,
		},
		{
			"dup coins",
			sdk.Coins{sdk.NewInt64Coin("anemo", 500), sdk.NewInt64Coin("anemo", 500)},
			true,
		},
		{
			"not evm coins",
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 500)),
			true,
		},
		{
			"negative coins",
			sdk.Coins{sdk.Coin{Denom: "anemo", Amount: sdkmath.NewInt(-500)}},
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := keeper.ValidateEvmCoins(tt.coins)
			if tt.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestConvertOneUnemoToAnemoIfNeeded() {
	anemoNeeded := sdkmath.NewInt(200)
	tests := []struct {
		name          string
		startingCoins sdk.Coins
		expectedCoins sdk.Coins
		success       bool
	}{
		{
			"not enough unemo for conversion",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100)),
			false,
		},
		{
			"converts 1 unemo to anemo",
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 10), sdk.NewInt64Coin("anemo", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 9), sdk.NewInt64Coin("anemo", 1_000_000_000_100)),
			true,
		},
		{
			"conversion not needed",
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 10), sdk.NewInt64Coin("anemo", 200)),
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 10), sdk.NewInt64Coin("anemo", 200)),
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.FundAccountWithNemo(suite.Addrs[0], tt.startingCoins)
			err := suite.EvmBankKeeper.ConvertOneUnemoToAnemoIfNeeded(suite.Ctx, suite.Addrs[0], anemoNeeded)
			moduleNemo := suite.BankKeeper.GetBalance(suite.Ctx, suite.AccountKeeper.GetModuleAddress(types.ModuleName), "unemo")
			if tt.success {
				suite.Require().NoError(err)
				if tt.startingCoins.AmountOf("anemo").LT(anemoNeeded) {
					suite.Require().Equal(sdk.OneInt(), moduleNemo.Amount)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(sdk.ZeroInt(), moduleNemo.Amount)
			}

			anemo := suite.Keeper.GetBalance(suite.Ctx, suite.Addrs[0])
			suite.Require().Equal(tt.expectedCoins.AmountOf("anemo"), anemo)
			unemo := suite.BankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "unemo")
			suite.Require().Equal(tt.expectedCoins.AmountOf("unemo"), unemo.Amount)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestConvertAnemoToUnemo() {
	tests := []struct {
		name          string
		startingCoins sdk.Coins
		expectedCoins sdk.Coins
	}{
		{
			"not enough unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 100), sdk.NewInt64Coin("unemo", 0)),
		},
		{
			"converts anemo for 1 unemo",
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 10), sdk.NewInt64Coin("anemo", 1_000_000_000_003)),
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 11), sdk.NewInt64Coin("anemo", 3)),
		},
		{
			"converts more than 1 unemo of anemo",
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 10), sdk.NewInt64Coin("anemo", 8_000_000_000_123)),
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 18), sdk.NewInt64Coin("anemo", 123)),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			err := suite.App.FundModuleAccount(suite.Ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("unemo", 10)))
			suite.Require().NoError(err)
			suite.FundAccountWithNemo(suite.Addrs[0], tt.startingCoins)
			err = suite.EvmBankKeeper.ConvertAnemoToUnemo(suite.Ctx, suite.Addrs[0])
			suite.Require().NoError(err)
			anemo := suite.Keeper.GetBalance(suite.Ctx, suite.Addrs[0])
			suite.Require().Equal(tt.expectedCoins.AmountOf("anemo"), anemo)
			unemo := suite.BankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "unemo")
			suite.Require().Equal(tt.expectedCoins.AmountOf("unemo"), unemo.Amount)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestSplitAnemoCoins() {
	tests := []struct {
		name          string
		coins         sdk.Coins
		expectedCoins sdk.Coins
		shouldErr     bool
	}{
		{
			"invalid coins",
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 500)),
			nil,
			true,
		},
		{
			"empty coins",
			sdk.NewCoins(),
			sdk.NewCoins(),
			false,
		},
		{
			"unemo & anemo coins",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 8_000_000_000_123)),
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 8), sdk.NewInt64Coin("anemo", 123)),
			false,
		},
		{
			"only anemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 10_123)),
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 10_123)),
			false,
		},
		{
			"only unemo",
			sdk.NewCoins(sdk.NewInt64Coin("anemo", 5_000_000_000_000)),
			sdk.NewCoins(sdk.NewInt64Coin("unemo", 5)),
			false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			unemo, anemo, err := keeper.SplitAnemoCoins(tt.coins)
			if tt.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tt.expectedCoins.AmountOf("unemo"), unemo.Amount)
				suite.Require().Equal(tt.expectedCoins.AmountOf("anemo"), anemo)
			}
		})
	}
}

func TestEvmBankKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(evmBankKeeperTestSuite))
}
