package types_test

import (
	"testing"

	commonTypes "healthcheck/x/common"
	"healthcheck/x/healthcheck/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				PortId: commonTypes.HealthcheckPortID,
				ChainList: []types.Chain{
					{
						ChainId: "0",
					},
					{
						ChainId: "1",
					},
				},
				ChainCount: 2,
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated chain",
			genState: &types.GenesisState{
				ChainList: []types.Chain{
					{
						ChainId: "0",
					},
					{
						ChainId: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
