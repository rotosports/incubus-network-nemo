package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rotosports/fury/app"
)

func init() {
	furyConfig := sdk.GetConfig()
	app.SetBech32AddressPrefixes(furyConfig)
	app.SetBip44CoinType(furyConfig)
	furyConfig.Seal()
}
