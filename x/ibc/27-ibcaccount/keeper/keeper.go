package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	// This field us used to marshal transaction for counterparty chain.
	// Currently, we support only one counterparty chain per interchain account keeper.
	// TODO: support multiple counterparty codec.
	counterpartyCdc *codec.Codec

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
	cdc *codec.Codec, counterpartyCdc *codec.Codec, key sdk.StoreKey,
	capKey sdk.CapabilityKey, clientKeeper types.ClientKeeper,
	connnectionKeeper types.ConnectionKeeper, channelKeeper types.ChannelKeeper,
	accountKeeper types.AccountKeeper, router sdk.Router,
) Keeper {
	return Keeper{
		storeKey:          key,
		cdc:               cdc,
		counterpartyCdc:   counterpartyCdc,
		boundedCapability: capKey,
		clientKeeper:      clientKeeper,
		connectionKeeper:  connnectionKeeper,
		channelKeeper:     channelKeeper,

		accountKeeper: accountKeeper,
		router:        router,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s/%s", ibctypes.ModuleName, types.SubModuleName))
}
