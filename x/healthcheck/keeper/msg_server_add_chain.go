package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"healthcheck/x/healthcheck/types"
)

func (k msgServer) AddChain(goCtx context.Context, msg *types.MsgAddChain) (*types.MsgAddChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgAddChainResponse{}, nil
}
