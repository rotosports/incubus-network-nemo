package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/incubus-network/nemo/app"
)

func init() {
	nemoConfig := sdk.GetConfig()
	app.SetBech32AddressPrefixes(nemoConfig)
	app.SetBip44CoinType(nemoConfig)
	nemoConfig.Seal()
}
