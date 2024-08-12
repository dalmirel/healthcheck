package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"healthcheck/x/healthcheck/types"
)

func (k Keeper) ChainAll(goCtx context.Context, req *types.QueryAllChainRequest) (*types.QueryAllChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var chains []types.Chain
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	chainStore := prefix.NewStore(store, types.KeyPrefix(types.ChainKey))

	pageRes, err := query.Paginate(chainStore, req.Pagination, func(key []byte, value []byte) error {
		var chain types.Chain
		if err := k.cdc.Unmarshal(value, &chain); err != nil {
			return err
		}

		chains = append(chains, chain)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllChainResponse{Chain: chains, Pagination: pageRes}, nil
}

func (k Keeper) Chain(goCtx context.Context, req *types.QueryGetChainRequest) (*types.QueryGetChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	chain, found := k.GetChain(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetChainResponse{Chain: chain}, nil
}
