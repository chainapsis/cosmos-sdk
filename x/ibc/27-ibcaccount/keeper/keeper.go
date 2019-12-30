package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/tendermint/tendermint/libs/log"
)

// DefaultPacketTimeout is the default packet timeout relative to the current block height
const (
	DefaultPacketTimeout = 1000 // NOTE: in blocks
)

// Keeper defines the IBC transfer keeper
type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             *codec.Codec
	counterpartyCdc *codec.Codec
	codespace       sdk.CodespaceType

	// Capability key and port to which ICS20 is binded. Used for packet relaying.
	boundedCapability sdk.CapabilityKey

	clientKeeper     types.ClientKeeper
	connectionKeeper types.ConnectionKeeper
	channelKeeper    types.ChannelKeeper
	accountKeeper    types.AccountKeeper

	router sdk.Router
}

// NewKeeper creates a new IBC interchain-account Keeper instance
func NewKeeper(
	cdc *codec.Codec, counterpartyCdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType,
	capKey sdk.CapabilityKey, clientKeeper types.ClientKeeper,
	connnectionKeeper types.ConnectionKeeper, channelKeeper types.ChannelKeeper,
	accountKeeper types.AccountKeeper, router sdk.Router,
) Keeper {
	return Keeper{
		storeKey:          key,
		cdc:               cdc,
		counterpartyCdc:   counterpartyCdc,
		codespace:         sdk.CodespaceType(fmt.Sprintf("%s/%s", codespace, types.DefaultCodespace)), // "ibc/interchain-account",
		boundedCapability: capKey,
		clientKeeper:      clientKeeper,
		connectionKeeper:  connnectionKeeper,
		channelKeeper:     channelKeeper,

		accountKeeper: accountKeeper,
		router:        router,
	}
}

// todo: This will be removed when the routing module(ics26) is implemented.
func (k Keeper) UnmarshalPacketData(packet channelexported.PacketI) (types.IbcPacketData, error) {
	var data types.IbcPacketData
	err := k.counterpartyCdc.UnmarshalBinaryBare(packet.GetData(), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s/%s", ibctypes.ModuleName, types.SubModuleName))
}
