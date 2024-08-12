package healthcheck_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "healthcheck/testutil/keeper"
	"healthcheck/testutil/nullify"
	"healthcheck/x/healthcheck"
	"healthcheck/x/healthcheck/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		ChainList: []types.Chain{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		ChainCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.HealthcheckKeeper(t)
	healthcheck.InitGenesis(ctx, *k, genesisState)
	got := healthcheck.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	require.ElementsMatch(t, genesisState.ChainList, got.ChainList)
	require.Equal(t, genesisState.ChainCount, got.ChainCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
