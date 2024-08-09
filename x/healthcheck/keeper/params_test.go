package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "healthcheck/testutil/keeper"
	"healthcheck/x/healthcheck/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.HealthcheckKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
