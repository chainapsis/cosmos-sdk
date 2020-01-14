package keeper

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strings"
)

// ReceivePacket handles receiving packet
func (k Keeper) ReceivePacket(ctx sdk.Context, packet channelexported.PacketI, proof commitment.ProofI, height uint64) error {
	_, err := k.channelKeeper.RecvPacket(ctx, packet, proof, height, []byte{}, k.boundedCapability)
	if err != nil {
		return err
	}

	var data types.InterchainAccountPacketData
	err = k.cdc.UnmarshalBinaryBare(packet.GetData(), &data)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid packet data")
	}

	switch packetData := data.(type) {
	case types.RegisterIBCAccountPacketData:
		err := packetData.ValidateBasic()
		if err != nil {
			return err
		}
		// TODO: Send acknowledgement packet as generated address
		return k.RegisterIBCAccount(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), packetData.Salt)
	case types.RunTxPacketData:
		err := packetData.ValidateBasic()
		if err != nil {
			return err
		}
		tx, err := k.DeserializeTx(ctx, packetData.TxBytes)
		if err != nil {
			return err
		}
		// TODO: Send acknowledgement packet as 0x0 if run tx succeeds or as non 0x0 if run tx fails.
		_, err = k.RunTx(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), tx)
		return err
	}

	panic("unexpected packet data")
}

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

	identifier := types.GetIdentifier(sourcePort, sourceChannel)
	address, err := k.GenerateAddress(identifier, salt)
	if err != nil {
		return err
	}
	err = k.CreateAccount(ctx, address, identifier)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) CreateAccount(ctx sdk.Context, address sdk.AccAddress, identifier string) error {
	account := k.accountKeeper.GetAccount(ctx, address)
	// Don't block even if there is normal account,
	// because attackers can distrupt to create an interchain account
	// by sending some assets to estimated address in advance.
	if account != nil {
		if account.GetSequence() != 0 || account.GetPubKey() != nil {
			// If account is interchain account or is usable by someone.
			return sdkerrors.Wrap(types.ErrAccountAlreadyExist, account.String())
		}
		err := account.SetSequence(1)
		if err != nil {
			return err
		}
	} else {
		account = k.accountKeeper.NewAccountWithAddress(ctx, address)
		err := account.SetSequence(1)
		if err != nil {
			return err
		}
	}

	// Interchain accounts have the sequence "1" and nil public key.
	// Sequence never be increased without signing tx and sending this tx.
	// But, it is impossible to send tx without publishing the public key.
	// So, accounts that have the sequence "1" and nill public key are explicitly interchain accounts.
	k.accountKeeper.SetAccount(ctx, account)

	store := ctx.KVStore(k.storeKey)
	// Save the identifier for each address to check where the interchain account is made from.
	store.Set(address, []byte(identifier))

	return nil
}

// Determine account's address that will be created.
func (k Keeper) GenerateAddress(identifier string, salt string) ([]byte, error) {
	hash := tmhash.NewTruncated()
	hashsum := hash.Sum([]byte(identifier + salt))
	return hashsum, nil
}

func (k Keeper) CreateInterchainAccount(ctx sdk.Context, sourcePort, sourceChannel, salt string) error {
	// get the port and channel of the counterparty
	sourceChan, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channel.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChan.Counterparty.PortID
	destinationChannel := sourceChan.Counterparty.ChannelID

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channel.ErrSequenceSendNotFound
	}

	packetData := types.RegisterIBCAccountPacketData{
		Salt: salt,
	}
	pdBytes, err := k.counterpartyCdc.MarshalBinaryBare(packetData)
	if err != nil {
		return err
	}

	packet := channel.NewPacket(
		sequence,
		uint64(ctx.BlockHeight())+DefaultPacketTimeout,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		pdBytes,
	)

	return k.channelKeeper.SendPacket(ctx, packet, k.boundedCapability)
}

func (k Keeper) CreateOutgoingPacket(
	ctx sdk.Context,
	seq uint64,
	chainType,
	sourcePort,
	sourceChannel,
	destinationPort,
	destinationChannel string,
	data interface{},
) error {
	if chainType == types.CosmosSdkChainType {
		if data == nil {
			return types.ErrInvalidOutgoingData
		}

		var msgs []sdk.Msg

		switch data := data.(type) {
		case []sdk.Msg:
			msgs = data
		case sdk.Msg:
			msgs = []sdk.Msg{data}
		default:
			return types.ErrInvalidOutgoingData
		}

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
		if err != nil {
			return err
		}

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

	err := k.cdc.UnmarshalBinaryBare(txBytes, &tx)
	return tx, err
}

func (k Keeper) RunTx(ctx sdk.Context, sourcePort, sourceChannel string, tx types.InterchainAccountTx) (*sdk.Result, error) {
	identifier := types.GetIdentifier(sourcePort, sourceChannel)
	err := k.AuthenticateTx(ctx, tx, identifier)
	if err != nil {
		return nil, err
	}

	msgs := tx.Msgs

	logs := make([]string, 0, len(msgs))
	data := make([]byte, 0, len(msgs))

	events := ctx.EventManager().Events()

	// Use cache context.
	// Receive packet msg should succeed regardless of the result of logic.
	// But, if we just return the success even though handler is failed,
	// the leftovers of state transition in handler will remain.
	// However, this can make the unexpected error.
	// To solve this problem, use cache context instead of context,
	// and write the state transition if handler succeeds.
	cacheContext, writeFn := ctx.CacheContext()
	err = nil
	for _, msg := range msgs {
		result, msgErr := k.RunMsg(cacheContext, msg)
		if msgErr != nil {
			err = msgErr
			break
		}

		data = append(data, result.Data...)

		events = events.AppendEvents(result.Events)

		if len(result.Log) > 0 {
			logs = append(logs, result.Log)
		}
	}

	if err != nil {
		return nil, err
	}

	// Write the state transitions if all handlers succeed.
	writeFn()

	return &sdk.Result{
		Data:   data,
		Log:    strings.TrimSpace(strings.Join(logs, ",")),
		Events: events,
	}, nil
}

func (k Keeper) AuthenticateTx(ctx sdk.Context, tx types.InterchainAccountTx, identifier string) error {
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
		// Check where the interchain account is made from.
		path := store.Get(signer)
		if !bytes.Equal(path, []byte(identifier)) {
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
