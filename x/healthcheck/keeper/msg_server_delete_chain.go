package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"healthcheck/x/healthcheck/types"
)

func (k msgServer) DeleteChain(goCtx context.Context, msg *types.MsgDeleteChain) (*types.MsgDeleteChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgDeleteChainResponse{}, nil
}
