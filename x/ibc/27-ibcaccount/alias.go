package ibcaccount

import (
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
)

const (
	DefaultPacketTimeout = keeper.DefaultPacketTimeout
	SubModuleName        = types.SubModuleName
	StoreKey             = types.StoreKey
	RouterKey            = types.RouterKey
	QuerierRoute         = types.QuerierRoute
)

var (
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
	GetIdentifier = types.GetIdentifier

	// Errors
	ErrAccountAlreadyExist  = types.ErrAccountAlreadyExist
	ErrUnsupportedChainType = types.ErrUnsupportedChainType
	ErrInvalidOutgoingData  = types.ErrInvalidOutgoingData
	ErrContentIsEmpty       = types.ErrContentIsEmpty

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper           = keeper.Keeper
	ClientKeeper     = types.ClientKeeper
	ConnectionKeeper = types.ConnectionKeeper
	ChannelKeeper    = types.ChannelKeeper
	AccountKeeper    = types.AccountKeeper

	MsgRecvPacket                = types.MsgRecvPacket
	RegisterIBCAccountPacketData = types.RegisterIBCAccountPacketData
	RunTxPacketData              = types.RunTxPacketData
)
