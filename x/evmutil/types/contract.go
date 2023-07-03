// Derived from https://github.com/tharsis/evmos/blob/0bfaf0db7be47153bc651e663176ba8deca960b5/contracts/erc20.go
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	// Embed ERC20 JSON files
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/nautilus/x/evm/types"
)

var (
	//go:embed nautilus_json/ERC20MintableBurnable.json
	ERC20MintableBurnableJSON []byte

	// ERC20MintableBurnableContract is the compiled erc20 contract
	ERC20MintableBurnableContract evmtypes.CompiledContract

	// ERC20MintableBurnableAddress is the erc20 module address
	ERC20MintableBurnableAddress common.Address

	//go:embed nautilus_json/ERC20NemoWrappedCosmosCoin.json
	ERC20NemoWrappedCosmosCoinJSON []byte

	// ERC20NemoWrappedCosmosCoinContract is the compiled erc20 contract
	ERC20NemoWrappedCosmosCoinContract evmtypes.CompiledContract
)

func init() {
	ERC20MintableBurnableAddress = ModuleEVMAddress

	err := json.Unmarshal(ERC20MintableBurnableJSON, &ERC20MintableBurnableContract)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal ERC20MintableBurnableJSON: %s. %s", err, string(ERC20MintableBurnableJSON)))
	}

	if len(ERC20MintableBurnableContract.Bin) == 0 {
		panic("loading ERC20MintableBurnable contract failed")
	}

	err = json.Unmarshal(ERC20NemoWrappedCosmosCoinJSON, &ERC20NemoWrappedCosmosCoinContract)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal ERC20NemoWrappedCosmosCoinJSON: %s. %s", err, string(ERC20NemoWrappedCosmosCoinJSON)))
	}

	if len(ERC20NemoWrappedCosmosCoinContract.Bin) == 0 {
		panic("loading ERC20NemoWrappedCosmosCoin contract failed")
	}
}
