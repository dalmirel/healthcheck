package healthcheck

import (
	"fmt"

	commonTypes "healthcheck/x/common"
	"healthcheck/x/healthcheck/keeper"
	"healthcheck/x/healthcheck/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

type IBCModule struct {
	keeper keeper.Keeper
}

func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
// Health check / registry chain should not initialize handshake! Return error.
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {

	return version, sdkerrors.Wrap(types.ErrInvalidHandshakeFlow, "handshake should be initiated by monitored chain!")
}

// OnChanOpenTry implements the IBCModule interface
// Since monitored chain is initializing the handshake
// health check / registry chain should try to answer the init
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {

	// Require portID is the portID module is bound to
	boundPort := im.keeper.GetPort(ctx)
	if boundPort != portID {
		return "", sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	if counterpartyVersion != commonTypes.Version {
		return "", sdkerrors.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: got: %s, expected %s", counterpartyVersion, commonTypes.Version)
	}

	if counterparty.PortId != commonTypes.MonitoredPortID {
		return "", sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, commonTypes.MonitoredPortID)
	}

	// TODO Mirel?  - this is not needed, since the Init does not initiate handshake:
	// Module may have already claimed capability in OnChanOpenInit in the case of crossing hellos
	// (ie chainA and chainB both call ChanOpenInit before one of them calls ChanOpenTry)
	// If module can already authenticate the capability then module already owns it so we don't need to claim
	// Otherwise, module does not have channel capability and we must claim it from IBC
	//if !im.keeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {

	// Only claim channel capability passed back by IBC module if we do not already own it
	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", err
	}
	//}

	// initiate monitoring for chain - if not registered yet:
	// 1. check connection hops?
	if len(connectionHops) != 1 {
		return "", sdkerrors.Wrapf(types.ErrInvalidNumOfConnectionHops, "only 1 hop expected")
	}
	// 2. check if chain is already monitored - skip registration then
	chainId, err := im.keeper.GetCounterpartyChainIDFromConnection(ctx, connectionHops[0])
	if err != nil {
		return "", err
	}

	monitoredChain, found := im.keeper.GetChain(ctx, chainId)
	if !found {
		return "", sdkerrors.Wrapf(types.ErrChainNotRegistered, "chain with the chain ID %s isn't registered yet", monitoredChain.GetChainId())
	}

	// initialize CHAIN data for monitoring, without channelId data ->
	// This is set on OnChanOpenConfirm (finalized channel opening)
	// another option? Set everything on OnChanOpeneConfirm?
	im.keeper.SetChain(ctx, monitoredChain)

	// here we can read updateInterval and TimeoutInterval
	// set by the chain initiating the handshake - this should be received
	// or leave those on DEFAULT VALUES!
	monitoredChain.UpdateInterval = types.DefaultUpdateInterval
	monitoredChain.TimeoutInterval = types.DefaultTimeoutInterval

	// return types.Version, nil // where version is healthcheck-1
	return commonTypes.Version, nil
}

// OnChanOpenAck implements the IBCModule interface
// MONITORED CHAIN
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	_,
	counterpartyVersion string,
) error {
	return sdkerrors.Wrap(types.ErrInvalidHandshakeFlow, "handshake should be initiated by monitored chain!")

}

// OnChanOpenConfirm implements the IBCModule interface
// HEALTH-CHECK/REGISTRY chain step: add chain to list of monitored chains!
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	monitoredChainID, err := im.keeper.GetCounterpartyChainIDFromChannel(ctx, portID, channelID)
	if err != nil {
		return err
	}

	monitoredChain, found := im.keeper.GetChain(ctx, monitoredChainID)
	if !found {
		return sdkerrors.Wrapf(types.ErrChainNotRegistered, "chain with the chain ID %s isn't registered yet", monitoredChainID)
	}

	// I think I should probably store channel Id here, since this is the last step of the handshake?
	monitoredChain.ChannelId = channelID
	// TODO Mirel: do we need registration height for some reason?
	im.keeper.SetChain(ctx, monitoredChain)

	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for channels
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface
// health check module / chain should process received packets -> when monitored chain sends the packed, healthcheck (regristry) chain should update
// registry data for the monitored chain
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	//TODO Mirel?: what should we set as result of the acknowledgment
	var ack channeltypes.Acknowledgement

	// this line is used by starport scaffolding # oracle/packet/module/recv

	var modulePacketData commonTypes.HealthcheckPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error()))
	}

	// packet is received, we need to check from which chain - it contains heart beat update
	// TODO Mirel: should this be source or destination?
	chainId, err := im.keeper.GetCounterpartyChainIDFromChannel(ctx, modulePacket.DestinationPort, modulePacket.DestinationChannel)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	monitoredChain, found := im.keeper.GetChain(ctx, chainId)
	if !found {
		return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrChainNotRegistered, "chain with the chain ID %s isn't registered yet", chainId))
	}

	// Dispatch packet
	switch packet := modulePacketData.Packet.(type) {
	case *commonTypes.HealthcheckPacketData_HealthCheckUpdate:
		// check timestamp and block data from healtCheck update received from the monitored chain!
		// maybe rename this to heart beat?
		if monitoredChain.Status.Timestamp >= packet.HealthCheckUpdate.Timestamp ||
			monitoredChain.Status.Block >= packet.HealthCheckUpdate.Block {

			err := fmt.Errorf("old heart beat update has already been received for chain with chain ID %s", monitoredChain.ChainId)
			return channeltypes.NewErrorAcknowledgement(err)
		}

		// if this is newer heart beat IBC update:
		monitoredChain.Status.Block = packet.HealthCheckUpdate.Block
		monitoredChain.Status.Timestamp = packet.HealthCheckUpdate.Timestamp

		// set activity status to "Active"
		monitoredChain.Status.Activity = uint64(types.Active)

		monitoredChain.Status.HealthCheckBlockHeight = uint64(ctx.BlockHeight())
		im.keeper.SetChain(ctx, monitoredChain)

	// this line is used by starport scaffolding # ibc/packet/module/recv
	default:
		err := fmt.Errorf("unrecognized %s packet type: %T", types.ModuleName, packet)
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return ack
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet acknowledgement: %v", err)
	}

	// this line is used by starport scaffolding # oracle/packet/module/ack

	var modulePacketData commonTypes.HealthcheckPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error())
	}

	var eventType string

	// Dispatch packet
	switch packet := modulePacketData.Packet.(type) {
	// this line is used by starport scaffolding # ibc/packet/module/ack
	default:
		errMsg := fmt.Sprintf("unrecognized %s packet type: %T", types.ModuleName, packet)
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAck, fmt.Sprintf("%v", ack)),
		),
	)

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
				sdk.NewAttribute(types.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
				sdk.NewAttribute(types.AttributeKeyAckError, resp.Error),
			),
		)
	}

	return nil
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var modulePacketData commonTypes.HealthcheckPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error())
	}

	// Dispatch packet
	switch packet := modulePacketData.Packet.(type) {
	// this line is used by starport scaffolding # ibc/packet/module/timeout
	default:
		errMsg := fmt.Sprintf("unrecognized %s packet type: %T", types.ModuleName, packet)
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	return nil
}
