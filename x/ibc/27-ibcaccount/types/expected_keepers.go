package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clientexported "github.com/cosmos/cosmos-sdk/x/ibc/02-client/exported"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	ibcaccountexported "github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/exported"
)

// ClientKeeper defines the expected IBC client keeper
type ClientKeeper interface {
	GetConsensusState(ctx sdk.Context, clientID string) (connection clientexported.ConsensusState, found bool)
}

// ConnectionKeeper defines the expected IBC connection keeper
type ConnectionKeeper interface {
	GetConnection(ctx sdk.Context, connectionID string) (connection connection.ConnectionEnd, found bool)
}

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channel.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, packet channelexported.PacketI, portCapability sdk.CapabilityKey) error
	RecvPacket(ctx sdk.Context, packet channelexported.PacketI, proof commitment.ProofI, proofHeight uint64, acknowledgement []byte, portCapability sdk.CapabilityKey) (channelexported.PacketI, error)
}

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) ibcaccountexported.Account
	NewAccount(ctx sdk.Context, acc ibcaccountexported.Account) ibcaccountexported.Account
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) ibcaccountexported.Account
	SetAccount(ctx sdk.Context, acc ibcaccountexported.Account)
}
