package keeper

import (
	"healthcheck/x/monitored/types"
)

var _ types.QueryServer = Keeper{}
