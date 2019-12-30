package ibcaccount

import (
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
)

const (
	DefaultPacketTimeout = keeper.DefaultPacketTimeout
	DefaultCodespace     = types.DefaultCodespace
	SubModuleName        = types.SubModuleName
	StoreKey             = types.StoreKey
	RouterKey            = types.RouterKey
	QuerierRoute         = types.QuerierRoute
)

var (
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
	GetIdentifier = types.GetIdentifier

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper           = keeper.Keeper
	ClientKeeper     = types.ClientKeeper
	ConnectionKeeper = types.ConnectionKeeper
	ChannelKeeper    = types.ChannelKeeper
	AccountKeeper    = types.AccountKeeper

	RegisterIBCAccountPacketData = types.RegisterIBCAccountPacketData
	RunTxPacketData              = types.RunTxPacketData
)
