package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ IbcPacketData = RegisterIBCAccountPacketData{}
var _ IbcPacketData = RunTxPacketData{}

type RegisterIBCAccountPacketData struct {
	Salt string `json:"salt"`
}

func (pd RegisterIBCAccountPacketData) ValidateBasic() error {
	if len(pd.Salt) == 0 {
		return sdkerrors.Wrap(ErrContentIsEmpty, "salt")
	}
	return nil
}

type RunTxPacketData struct {
	TxBytes []byte `json:"tx_bytes"`
}

func (pd RunTxPacketData) ValidateBasic() error {
	if len(pd.TxBytes) == 0 {
		return sdkerrors.Wrap(ErrContentIsEmpty, "txBytes")
	}
	return nil
}
