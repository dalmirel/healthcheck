package types

import (
	"fmt"
	commonTypes "healthcheck/x/common"

	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId:    commonTypes.HealthcheckPortID,
		ChainList: []Chain{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	// Check for duplicated index in chain
	chainIndexMap := make(map[string]struct{})

	for _, elem := range gs.ChainList {
		index := string(ChainKey(elem.ChainId))
		if _, ok := chainIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for chain")
		}
		chainIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
