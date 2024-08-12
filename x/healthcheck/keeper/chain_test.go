package keeper_test

import (
	"strconv"
	"testing"

	keepertest "healthcheck/testutil/keeper"
	"healthcheck/testutil/nullify"
	"healthcheck/x/healthcheck/keeper"
	"healthcheck/x/healthcheck/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNChain(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Chain {
	items := make([]types.Chain, n)
	for i := range items {
		items[i].ChainId = strconv.Itoa(i)

		keeper.SetChain(ctx, items[i])
	}
	return items
}

func TestChainGet(t *testing.T) {
	keeper, ctx := keepertest.HealthcheckKeeper(t)
	items := createNChain(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetChain(ctx, item.ChainId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestChainRemove(t *testing.T) {
	keeper, ctx := keepertest.HealthcheckKeeper(t)
	items := createNChain(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveChain(ctx, item.ChainId)
		_, found := keeper.GetChain(ctx, item.ChainId)
		require.False(t, found)
	}
}

func TestChainGetAll(t *testing.T) {
	keeper, ctx := keepertest.HealthcheckKeeper(t)
	items := createNChain(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllChain(ctx)),
	)
}
