package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/healthcheck module sentinel errors
var (
	ErrSample                     = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout       = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion             = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrInvalidHandshakeFlow       = sdkerrors.Register(ModuleName, 1502, "invalid handshake flow initialization")
	ErrInvalidNumOfConnectionHops = sdkerrors.Register(ModuleName, 1503, "invalid num of hops")
	ErrChainNotRegistered         = sdkerrors.Register(ModuleName, 1504, "chain not registered for monitoring")
)
