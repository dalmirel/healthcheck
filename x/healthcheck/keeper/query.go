package keeper

import (
	"healthcheck/x/healthcheck/types"
)

var _ types.QueryServer = Keeper{}
