package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strings"
)

func (k Keeper) RegisterIBCAccount(
	ctx sdk.Context,
	sourcePort,
	sourceChannel,
	salt string,
) error {
	_, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channel.ErrChannelNotFound, sourceChannel)
	}

	address, err := k.GenerateAddress(types.GetIdentifier(sourcePort, sourceChannel), salt)
	if err != nil {
		return err
	}
	err = k.CreateAccount(ctx, address)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) CreateAccount(ctx sdk.Context, address sdk.AccAddress) error {
	// Currently, it seems that there is no way to get the information of counterparty chain.
	// So, just don't use path for hackathon version.
	account := k.accountKeeper.GetAccount(ctx, address)
	if account != nil {
		if account.GetSequence() != 1 || account.GetPubKey() != nil {
			return sdkerrors.Wrap(types.ErrAccountAlreadyExist, account.String())
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
		return sdkerrors.Wrap(types.ErrUnsupportedChainType, chainType)
	}
}

func (k Keeper) DeserializeTx(ctx sdk.Context, txBytes []byte) (types.InterchainAccountTx, error) {
	tx := types.InterchainAccountTx{}

	err := k.counterpartyCdc.UnmarshalBinaryBare(txBytes, &tx)
	return tx, err
}

func (k Keeper) RunTx(ctx sdk.Context, tx types.InterchainAccountTx) (*sdk.Result, error) {
	err := k.AuthenticateTx(ctx, tx)
	if err != nil {
		return nil, err
	}

	msgs := tx.Msgs

	logs := make([]string, 0, len(msgs))
	data := make([]byte, 0, len(msgs))

	events := ctx.EventManager().Events()

	for _, msg := range msgs {
		result, err := k.RunMsg(ctx, msg)
		if err != nil {
			return result, err
		}

		data = append(data, result.Data...)

		events = events.AppendEvents(result.Events)

		if len(result.Log) > 0 {
			logs = append(logs, result.Log)
		}
	}

	return &sdk.Result{
		Data:   data,
		Log:    strings.TrimSpace(strings.Join(logs, ",")),
		Events: events,
	}, nil
}

func (k Keeper) AuthenticateTx(ctx sdk.Context, tx types.InterchainAccountTx) error {
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
			return sdkerrors.ErrUnauthorized
		}
	}

	return nil
}

func (k Keeper) RunMsg(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	hander := k.router.Route(ctx, msg.Route())
	if hander == nil {
		return nil, types.ErrInvalidRoute
	}

	return hander(ctx, msg)
}
