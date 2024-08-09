package monitored_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "healthcheck/testutil/keeper"
	"healthcheck/testutil/nullify"
	"healthcheck/x/monitored"
	"healthcheck/x/monitored/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MonitoredKeeper(t)
	monitored.InitGenesis(ctx, *k, genesisState)
	got := monitored.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
