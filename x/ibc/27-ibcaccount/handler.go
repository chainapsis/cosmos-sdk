package ibcaccount

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
)

func HandleRegisterIBCAccount(ctx sdk.Context, k Keeper, sourcePort, sourceChannel string, packet RegisterIBCAccountPacketData) sdk.Result {
	err := k.RegisterIBCAccount(ctx, sourcePort, sourceChannel, packet.Salt)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleRunTx(ctx sdk.Context, k Keeper, packet RunTxPacketData) sdk.Result {
	interchainAccountTx, err := k.DeserializeTx(ctx, packet.TxBytes)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	return k.RunTx(ctx, interchainAccountTx)
}
