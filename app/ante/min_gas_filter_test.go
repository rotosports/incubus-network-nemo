package ante_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/evmos/nautilus/x/evm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/incubus-network/nemo/app"
	"github.com/incubus-network/nemo/app/ante"
)

func mustParseDecCoins(value string) sdk.DecCoins {
	coins, err := sdk.ParseDecCoins(strings.ReplaceAll(value, ";", ","))
	if err != nil {
		panic(err)
	}

	return coins
}

func TestEvmMinGasFilter(t *testing.T) {
	tApp := app.NewTestApp()
	handler := ante.NewEvmMinGasFilter(tApp.GetEvmKeeper())

	ctx := tApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})
	tApp.GetEvmKeeper().SetParams(ctx, evmtypes.Params{
		EvmDenom: "anemo",
	})

	testCases := []struct {
		name                 string
		minGasPrices         sdk.DecCoins
		expectedMinGasPrices sdk.DecCoins
	}{
		{
			"no min gas prices",
			mustParseDecCoins(""),
			mustParseDecCoins(""),
		},
		{
			"zero unemo gas price",
			mustParseDecCoins("0unemo"),
			mustParseDecCoins("0unemo"),
		},
		{
			"non-zero unemo gas price",
			mustParseDecCoins("0.001unemo"),
			mustParseDecCoins("0.001unemo"),
		},
		{
			"zero unemo gas price, min anemo price",
			mustParseDecCoins("0unemo;100000anemo"),
			mustParseDecCoins("0unemo"), // anemo is removed
		},
		{
			"zero unemo gas price, min anemo price, other token",
			mustParseDecCoins("0unemo;100000anemo;0.001other"),
			mustParseDecCoins("0unemo;0.001other"), // anemo is removed
		},
		{
			"non-zero unemo gas price, min anemo price",
			mustParseDecCoins("0.25unemo;100000anemo;0.001other"),
			mustParseDecCoins("0.25unemo;0.001other"), // anemo is removed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})

			ctx = ctx.WithMinGasPrices(tc.minGasPrices)
			mmd := MockAnteHandler{}

			_, err := handler.AnteHandle(ctx, nil, false, mmd.AnteHandle)
			require.NoError(t, err)
			require.True(t, mmd.WasCalled)

			assert.NoError(t, mmd.CalledCtx.MinGasPrices().Validate())
			assert.Equal(t, tc.expectedMinGasPrices, mmd.CalledCtx.MinGasPrices())
		})
	}
}
