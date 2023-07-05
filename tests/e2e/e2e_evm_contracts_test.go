package e2e_test

import (
	"context"
	"math/big"
	"time"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/incubus-network/fury/app"
	earntypes "github.com/incubus-network/fury/x/earn/types"
	evmutiltypes "github.com/incubus-network/fury/x/evmutil/types"

	"github.com/incubus-network/fury/tests/e2e/contracts/greeter"
	"github.com/incubus-network/fury/tests/util"
)

func (suite *IntegrationTestSuite) TestEthCallToGreeterContract() {
	// this test manipulates state of the Greeter contract which means other tests shouldn't use it.

	// setup funded account to interact with contract
	user := suite.Fury.NewFundedAccount("greeter-contract-user", sdk.NewCoins(ufury(1e6)))

	greeterAddr := suite.Fury.ContractAddrs["greeter"]
	contract, err := greeter.NewGreeter(greeterAddr, suite.Fury.EvmClient)
	suite.NoError(err)

	beforeGreeting, err := contract.Greet(nil)
	suite.NoError(err)

	updatedGreeting := "look at me, using the evm"
	tx, err := contract.SetGreeting(user.EvmAuth, updatedGreeting)
	suite.NoError(err)

	_, err = util.WaitForEvmTxReceipt(suite.Fury.EvmClient, tx.Hash(), 10*time.Second)
	suite.NoError(err)

	afterGreeting, err := contract.Greet(nil)
	suite.NoError(err)

	suite.Equal("what's up!", beforeGreeting)
	suite.Equal(updatedGreeting, afterGreeting)
}

func (suite *IntegrationTestSuite) TestEthCallToErc20() {
	randoReceiver := util.SdkToEvmAddress(app.RandomAddress())
	amount := big.NewInt(1)

	// make unauthenticated eth_call query to check balance
	beforeBalance := suite.Fury.GetErc20Balance(suite.DeployedErc20.Address, randoReceiver)

	// make authenticate eth_call to transfer tokens
	res := suite.FundFuryErc20Balance(randoReceiver, amount)
	suite.NoError(res.Err)

	// make another unauthenticated eth_call query to check new balance
	afterBalance := suite.Fury.GetErc20Balance(suite.DeployedErc20.Address, randoReceiver)

	suite.BigIntsEqual(big.NewInt(0), beforeBalance, "expected before balance to be zero")
	suite.BigIntsEqual(amount, afterBalance, "unexpected post-transfer balance")
}

func (suite *IntegrationTestSuite) TestEip712BasicMessageAuthorization() {
	// create new funded account
	sender := suite.Fury.NewFundedAccount("eip712-msgSend", sdk.NewCoins(ufury(2e4)))
	receiver := app.RandomAddress()

	// setup message for sending some fury to random receiver
	msgs := []sdk.Msg{
		banktypes.NewMsgSend(sender.SdkAddress, receiver, sdk.NewCoins(ufury(1e3))),
	}

	// create tx
	tx := suite.NewEip712TxBuilder(
		sender,
		suite.Fury,
		1e6,
		sdk.NewCoins(ufury(1e4)),
		msgs,
		"this is a memo",
	).GetTx()

	txBytes, err := suite.Fury.EncodingConfig.TxConfig.TxEncoder()(tx)
	suite.NoError(err)

	// broadcast tx
	res, err := suite.Fury.Tx.BroadcastTx(context.Background(), &txtypes.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
	})
	suite.NoError(err)
	suite.Equal(sdkerrors.SuccessABCICode, res.TxResponse.Code)

	_, err = util.WaitForSdkTxCommit(suite.Fury.Tx, res.TxResponse.TxHash, 6*time.Second)
	suite.NoError(err)

	// check that the message was processed & the fury is transferred.
	balRes, err := suite.Fury.Bank.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: receiver.String(),
		Denom:   "ufury",
	})
	suite.NoError(err)
	suite.Equal(sdk.NewInt(1e3), balRes.Balance.Amount)
}

// Note that this test works because the deployed erc20 is configured in evmutil & earn params.
func (suite *IntegrationTestSuite) TestEip712ConvertToCoinAndDepositToEarn() {
	amount := sdk.NewInt(1e2) // 0.0002 USDC
	sdkDenom := suite.DeployedErc20.CosmosDenom

	// create new funded account
	depositor := suite.Fury.NewFundedAccount("eip712-earn-depositor", sdk.NewCoins(ufury(1e5)))
	// give them erc20 balance to deposit
	fundRes := suite.FundFuryErc20Balance(depositor.EvmAddress, amount.BigInt())
	suite.NoError(fundRes.Err)

	// setup messages for convert to coin & deposit into earn
	convertMsg := evmutiltypes.NewMsgConvertERC20ToCoin(
		evmutiltypes.NewInternalEVMAddress(depositor.EvmAddress),
		depositor.SdkAddress,
		evmutiltypes.NewInternalEVMAddress(suite.DeployedErc20.Address),
		amount,
	)
	depositMsg := earntypes.NewMsgDeposit(
		depositor.SdkAddress.String(),
		sdk.NewCoin(sdkDenom, amount),
		earntypes.STRATEGY_TYPE_SAVINGS,
	)
	msgs := []sdk.Msg{
		// convert to coin
		&convertMsg,
		// deposit into earn
		depositMsg,
	}

	// create tx
	tx := suite.NewEip712TxBuilder(
		depositor,
		suite.Fury,
		1e6,
		sdk.NewCoins(ufury(1e4)),
		msgs,
		"depositing my USDC into Earn!",
	).GetTx()

	txBytes, err := suite.Fury.EncodingConfig.TxConfig.TxEncoder()(tx)
	suite.NoError(err)

	// broadcast tx
	res, err := suite.Fury.Tx.BroadcastTx(context.Background(), &txtypes.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
	})
	suite.NoError(err)
	suite.Equal(sdkerrors.SuccessABCICode, res.TxResponse.Code)

	_, err = util.WaitForSdkTxCommit(suite.Fury.Tx, res.TxResponse.TxHash, 6*time.Second)
	suite.NoError(err)

	// check that depositor no longer has erc20 balance
	balance := suite.Fury.GetErc20Balance(suite.DeployedErc20.Address, depositor.EvmAddress)
	suite.BigIntsEqual(big.NewInt(0), balance, "expected no erc20 balance")

	// check that account has an earn deposit position
	earnRes, err := suite.Fury.Earn.Deposits(context.Background(), &earntypes.QueryDepositsRequest{
		Depositor: depositor.SdkAddress.String(),
		Denom:     sdkDenom,
	})
	suite.NoError(err)
	suite.Len(earnRes.Deposits, 1)
	suite.Equal(sdk.NewDecFromInt(amount), earnRes.Deposits[0].Shares.AmountOf(sdkDenom))

	// withdraw deposit & convert back to erc20 (this allows refund to recover erc20s used in test)
	withdraw := earntypes.NewMsgWithdraw(
		depositor.SdkAddress.String(),
		sdk.NewCoin(sdkDenom, amount),
		earntypes.STRATEGY_TYPE_SAVINGS,
	)
	convertBack := evmutiltypes.NewMsgConvertCoinToERC20(
		depositor.SdkAddress.String(),
		depositor.EvmAddress.Hex(),
		sdk.NewCoin(sdkDenom, amount),
	)
	withdrawAndConvertBack := util.FuryMsgRequest{
		Msgs:      []sdk.Msg{withdraw, &convertBack},
		GasLimit:  3e5,
		FeeAmount: sdk.NewCoins(ufury(300)),
		Data:      "withdrawing from earn & converting back to erc20",
	}
	lastRes := depositor.SignAndBroadcastFuryTx(withdrawAndConvertBack)
	suite.NoError(lastRes.Err)
}
