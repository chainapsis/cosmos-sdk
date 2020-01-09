package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"strings"
)

// nolint: unused
func (k Keeper) OnChanOpenInit(
	ctx sdk.Context,
	order channel.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channel.Counterparty,
	version string,
) error {
	// only ordered channels allowed
	if order != channel.ORDERED {
		return sdkerrors.Wrap(channel.ErrInvalidChannel, "channel must be ORDERED")
	}

	// NOTE: here the capability key name defines the port ID of the counterparty
	// only allow channels to "interchain-account" port on counterparty chain
	if counterparty.PortID != k.boundedCapability.Name() {
		return sdkerrors.Wrapf(
			port.ErrInvalidPort,
			"counterparty port ID doesn't match the capability key (%s ≠ %s)", counterparty.PortID, k.boundedCapability.Name(),
		)
	}

	if strings.TrimSpace(version) != "" {
		return sdkerrors.Wrap(ibctypes.ErrInvalidVersion, "version must be blank")
	}

	// NOTE: as the escrow address is generated from both the port and channel IDs
	// there's no need to store it on a map.
	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenTry(
	ctx sdk.Context,
	order channel.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channel.Counterparty,
	version string,
	counterpartyVersion string,
) error {
	// only ordered channels allowed
	if order != channel.ORDERED {
		return sdkerrors.Wrap(channel.ErrInvalidChannel, "channel must be ORDERED")
	}

	// NOTE: here the capability key name defines the port ID of the counterparty
	// only allow channels to "interchain-account" port on counterparty chain
	if counterparty.PortID != k.boundedCapability.Name() {
		return sdkerrors.Wrapf(
			port.ErrInvalidPort,
			"counterparty port ID doesn't match the capability key (%s ≠ %s)", counterparty.PortID, k.boundedCapability.Name(),
		)
	}

	if strings.TrimSpace(version) != "" {
		return sdkerrors.Wrap(ibctypes.ErrInvalidVersion, "version must be blank")
	}

	if strings.TrimSpace(counterpartyVersion) != "" {
		return sdkerrors.Wrap(ibctypes.ErrInvalidVersion, "counterparty version must be blank")
	}

	// NOTE: as the escrow address is generated from both the port and channel IDs
	// there's no need to store it on a map.
	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenAck(ctx sdk.Context, portID, channelID string, version string) error {
	if strings.TrimSpace(version) != "" {
		return sdkerrors.Wrap(ibctypes.ErrInvalidVersion, "version must be blank")
	}
	return nil
}

// nolint: unused
func (k Keeper) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	// no-op
	return nil
}

// onRecvPacket is called when an FTTransfer packet is received
// nolint: unused
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channelexported.PacketI) error {
	// TODO
	return nil
}

// nolint: unused
func (k Keeper) OnAcknowledgePacket(ctx sdk.Context, packet channelexported.PacketI, acknowledgement []byte) error {
	// no-op
	return nil
}

// nolint: unused
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channelexported.PacketI) error {
	// TODO
	return nil
}

// nolint: unused
func (k Keeper) OnTimeoutPacketClose(_ sdk.Context, _ channelexported.PacketI) {
	// TODO
}
