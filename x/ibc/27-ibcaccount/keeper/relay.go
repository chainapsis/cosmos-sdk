package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strings"
)

func (k Keeper) RegisterIBCAccount(ctx sdk.Context, sourcePort, sourceChannel, salt string) error {
	address, err := k.GenerateAddress(types.GetIdentifier(sourcePort, sourceChannel), salt)
	if err != nil {
		return err
	}
	sdkErr := k.CreateAccount(ctx, address)
	if sdkErr != nil {
		return sdkErr
	}

	return nil
}

func (k Keeper) CreateAccount(ctx sdk.Context, address sdk.AccAddress) error {
	// Currently, it seems that there is no way to get the information of counterparty chain.
	// So, just don't use path for hackathon version.
	account := k.accountKeeper.GetAccount(ctx, address)
	if account != nil {
		if account.GetSequence() != 1 || account.GetPubKey() != nil {
			return types.ErrAccountAlreadyExist(types.DefaultCodespace, address.String())
		}
	} else {
		account = k.accountKeeper.NewAccountWithAddress(ctx, address)
		err := account.SetSequence(1)
		if err != nil {
			return err
		}
		account = k.accountKeeper.NewAccount(ctx, account)
	}

	k.accountKeeper.SetAccount(ctx, account)

	store := ctx.KVStore(k.storeKey)
	// Ignore that which chain makes the interchain account.
	// Assume that only one to one communication exists for prototype version.
	store.Set(address, []byte{1})

	return nil
}

// Determine account's address that will be created.
func (k Keeper) GenerateAddress(identifier string, salt string) ([]byte, error) {
	hash := tmhash.NewTruncated()
	hashsum := hash.Sum([]byte(identifier + salt))
	return hashsum, nil
}

func (k Keeper) CreateOutgoingPacket(
	ctx sdk.Context,
	seq uint64,
	chainType,
	sourcePort,
	sourceChannel,
	destinationPort,
	destinationChannel string,
	msgs []sdk.Msg,
) error {
	if chainType == types.CosmosSdkChainType {
		interchainAccountTx := types.InterchainAccountTx{Msgs: msgs}

		txBytes, err := k.counterpartyCdc.MarshalBinaryBare(interchainAccountTx)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid packet data or codec")
		}

		// todo: consider
		packetData := types.RunTxPacketData{
			TxBytes: txBytes,
		}
		pdBytes, err := k.counterpartyCdc.MarshalBinaryBare(packetData)

		packet := channel.NewPacket(
			seq,
			uint64(ctx.BlockHeight())+DefaultPacketTimeout,
			sourcePort,
			sourceChannel,
			destinationPort,
			destinationChannel,
			pdBytes,
		)

		return k.channelKeeper.SendPacket(ctx, packet, k.boundedCapability)
	} else {
		return types.ErrUnsupportedChainType(types.DefaultCodespace, chainType)
	}
}

func (k Keeper) DeserializeTx(ctx sdk.Context, txBytes []byte) (types.InterchainAccountTx, error) {
	tx := types.InterchainAccountTx{}

	err := k.counterpartyCdc.UnmarshalBinaryBare(txBytes, &tx)
	return tx, err
}

func (k Keeper) RunTx(ctx sdk.Context, tx types.InterchainAccountTx) sdk.Result {
	err := k.AuthenticateTx(ctx, tx)
	if err != nil {
		return err.Result()
	}

	msgs := tx.Msgs

	logs := make([]string, 0, len(msgs))
	data := make([]byte, 0, len(msgs))
	var (
		code      sdk.CodeType
		codespace sdk.CodespaceType
	)
	events := ctx.EventManager().Events()

	for _, msg := range msgs {
		result := k.RunMsg(ctx, msg)
		if result.IsOK() == false {
			return result
		}

		data = append(data, result.Data...)

		events = events.AppendEvents(result.Events)

		if len(result.Log) > 0 {
			logs = append(logs, result.Log)
		}

		if !result.IsOK() {
			code = result.Code
			codespace = result.Codespace
			break
		}
	}

	return sdk.Result{
		Code:      code,
		Codespace: codespace,
		Data:      data,
		Log:       strings.TrimSpace(strings.Join(logs, ",")),
		GasUsed:   ctx.GasMeter().GasConsumed(),
		Events:    events,
	}
}

func (k Keeper) AuthenticateTx(ctx sdk.Context, tx types.InterchainAccountTx) sdk.Error {
	msgs := tx.Msgs

	seen := map[string]bool{}
	var signers []sdk.AccAddress
	for _, msg := range msgs {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	store := ctx.KVStore(k.storeKey)

	for _, signer := range signers {
		path := store.Get(signer)
		// Ignore that which chain makes the interchain account.
		// Assume that only one to one communication exists for prototype version.
		if len(path) == 0 {
			return sdk.ErrUnauthorized("unauthorized")
		}
	}

	return nil
}

func (k Keeper) RunMsg(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	//hander := k.router.Route(msg.Route())
	//if hander == nil {
	//	return sdk.ErrInternal("invalid route").Result()
	//}
	//
	//return hander(ctx, msg)
	return sdk.Result{}
}
