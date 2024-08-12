package types

import (
	"fmt"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId:    PortID,
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
	// Check for duplicated ID in chain
	chainIdMap := make(map[uint64]bool)
	chainCount := gs.GetChainCount()
	for _, elem := range gs.ChainList {
		if _, ok := chainIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for chain")
		}
		if elem.Id >= chainCount {
			return fmt.Errorf("chain id should be lower or equal than the last id")
		}
		chainIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
