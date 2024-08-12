package monitored

import (
	"fmt"

	commonTypes "healthcheck/x/common"

	"healthcheck/x/monitored/keeper"
	"healthcheck/x/monitored/types"

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
// MONITORED CHAIN step
// update
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

	// TODO Mirel: maybe this is not important?
	if order != channeltypes.ORDERED {
		return "", sdkerrors.Wrapf(channeltypes.ErrInvalidChannelOrdering, "expected %s channel, got %s ", channeltypes.ORDERED, order)
	}

	// Require portID is the portID module is bound to
	boundPort := im.keeper.GetPort(ctx)
	if boundPort != portID {
		return "", sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	if version != commonTypes.Version {
		return "", sdkerrors.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", version, commonTypes.Version)
	}

	// Counterparty portID is the portID of the healthcheck module
	if counterparty.PortId != commonTypes.HealthcheckPortID {
		return "", sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, commonTypes.HealthcheckPortID)
	}

	// Claim channel capability passed back by IBC module
	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", err
	}

	return version, nil
}

// OnChanOpenTry implements the IBCModule interface
// Monitored chain initiates connection to the registry - healthcheck chain
// HEALTH CHECK /REGISTRY chain step:
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

	return "", sdkerrors.Wrap(types.ErrInvalidHandshakeFlow, "channel handshake must be initiated by monitored chain")
}

// OnChanOpenAck implements the IBCModule interface
// MONITORED chain step
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	_,
	counterpartyVersion string,
) error {
	if counterpartyVersion != commonTypes.Version {
		return sdkerrors.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: %s, expected %s", counterpartyVersion, commonTypes.Version)
	}

	im.keeper.SetRegistryChainChannelID(ctx, channelID)
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
// HEALTH CHECK /REGISTRY chain step!
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return sdkerrors.Wrap(types.ErrInvalidHandshakeFlow, "channel handshake must be initiated by monitored chain")
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
	registryChainChannelID := im.keeper.GetRegistryChainChannelID(ctx)
	if registryChainChannelID != channelID {
		// should not happen since only one monitored-healthcheck channel is allowed
		return sdkerrors.Wrap(types.ErrUnexpectedChannelID, fmt.Sprintf("expected: %s, got: %s", registryChainChannelID, channelID))
	}

	im.keeper.SetRegistryChainChannelID(ctx, "")

	return nil
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(
		sdkerrors.ErrUnknownRequest,
		"cannot send packet data from monitored chain: port %v and source channel %t",
		modulePacket.SourcePort,
		modulePacket.SourceChannel,
	))
}

// OnAcknowledgementPacket implements the IBCModule interface
// TODO MIREL? - MONITORED CHAIN SHOULD process info about ack?
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

	return nil
}

// OnTimeoutPacket implements the IBCModule interface
// TODO MIREL? - MONITORED CHAIN SHOULD process info about timeout?
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var modulePacketData types.MonitoredPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error())
	}

	return nil
}
