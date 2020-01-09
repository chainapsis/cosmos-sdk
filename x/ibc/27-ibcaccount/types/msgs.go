package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

type MsgRecvPacket struct {
	Packet channelexported.PacketI `json:"packet" yaml:"packet"`
	Proofs []commitment.Proof      `json:"proofs" yaml:"proofs"`
	Height uint64                  `json:"height" yaml:"height"`
	Signer sdk.AccAddress          `json:"signer" yaml:"signer"`
}

// NewMsgRecvPacket creates a new MsgRecvPacket instance
func NewMsgRecvPacket(packet channelexported.PacketI, proofs []commitment.Proof, height uint64, signer sdk.AccAddress) MsgRecvPacket {
	return MsgRecvPacket{
		Packet: packet,
		Proofs: proofs,
		Height: height,
		Signer: signer,
	}
}

// Route implements sdk.Msg
func (MsgRecvPacket) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (MsgRecvPacket) Type() string {
	return "recv_packet"
}

// ValidateBasic implements sdk.Msg
func (msg MsgRecvPacket) ValidateBasic() error {
	if msg.Height == 0 {
		return sdkerrors.Wrap(ibctypes.ErrInvalidHeight, "height must be > 0")
	}
	if msg.Proofs == nil || len(msg.Proofs) == 0 {
		return sdkerrors.Wrap(commitment.ErrInvalidProof, "missing proof")
	}
	for _, proof := range msg.Proofs {
		if err := proof.ValidateBasic(); err != nil {
			return err
		}
	}
	if msg.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	return msg.Packet.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (msg MsgRecvPacket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRecvPacket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
