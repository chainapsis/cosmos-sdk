package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ IbcPacketData = RegisterIBCAccountPacketData{}
var _ IbcPacketData = RunTxPacketData{}

type RegisterIBCAccountPacketData struct {
	Salt string `json:"salt"`
}

func (pd RegisterIBCAccountPacketData) ValidateBasic() sdk.Error {
	if len(pd.Salt) == 0 {
		return sdk.ConvertError(ErrContentIsEmpty(DefaultCodespace, "salt"))
	}
	return nil
}

type RunTxPacketData struct {
	TxBytes []byte `json:"tx_bytes"`
}

func (pd RunTxPacketData) ValidateBasic() sdk.Error {
	if len(pd.TxBytes) == 0 {
		return sdk.ConvertError(ErrContentIsEmpty(DefaultCodespace, "tx bytes"))
	}
	return nil
}
