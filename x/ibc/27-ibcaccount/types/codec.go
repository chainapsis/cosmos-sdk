package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*IbcPacketData)(nil), nil)

	cdc.RegisterConcrete(RegisterIBCAccountPacketData{}, "ibc/ibcaccount/RegisterIBCAccountPacketData", nil)
	cdc.RegisterConcrete(RunTxPacketData{}, "ibc/ibcaccount/RunTxPacketData", nil)
	cdc.RegisterConcrete(InterchainAccountTx{}, "ibc/ibcaccount/InterchainAccountTx", nil)
}

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}
