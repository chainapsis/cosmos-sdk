package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

var _ sdk.Msg = MsgRegisterIBCAccount{}
var _ sdk.Msg = MsgRunTx{}

//todo: temp
type MsgRegisterIBCAccount struct {
	Salt          string         `json:"salt"`
	Signer        sdk.AccAddress `json:"signer" yaml:"signer"`
	SourcePort    string         `json:"source_port" yaml:"source_port"`
	SourceChannel string         `json:"source_channel" yaml:"source_channel"`
}

// Route implements sdk.Msg
func (MsgRegisterIBCAccount) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (MsgRegisterIBCAccount) Type() string {
	return "register_ibc_account"
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterIBCAccount) ValidateBasic() sdk.Error {
	if len(msg.Salt) == 0 {
		return sdk.ConvertError(ErrContentIsEmpty(DefaultCodespace, "salt"))
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgRegisterIBCAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterIBCAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

//todo: temp
type MsgRunTx struct {
	TxBytes []byte         `json:"tx_bytes"`
	Signer  sdk.AccAddress `json:"signer" yaml:"signer"`
}

// Route implements sdk.Msg
func (MsgRunTx) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (MsgRunTx) Type() string {
	return "run_tx"
}

// ValidateBasic implements sdk.Msg
func (msg MsgRunTx) ValidateBasic() sdk.Error {
	if len(msg.TxBytes) == 0 {
		return sdk.ConvertError(ErrContentIsEmpty(DefaultCodespace, "tx bytes"))
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgRunTx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRunTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
