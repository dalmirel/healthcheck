package monitored_test

import (
	"testing"

	keepertest "healthcheck/testutil/keeper"
	"healthcheck/testutil/nullify"
	commonTypes "healthcheck/x/common"
	"healthcheck/x/monitored"
	"healthcheck/x/monitored/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: commonTypes.MonitoredPortID,
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
