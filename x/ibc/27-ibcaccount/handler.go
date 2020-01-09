package ibcaccount

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
)

// HandleMsgRecvPacket defines the sdk.Handler for MsgRecvPacket
func HandleMsgRecvPacket(ctx sdk.Context, k Keeper, msg MsgRecvPacket) (*sdk.Result, error) {
	err := k.ReceivePacket(ctx, msg.Packet, msg.Proofs[0], msg.Height)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
