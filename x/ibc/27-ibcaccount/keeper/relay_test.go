package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypestm "github.com/cosmos/cosmos-sdk/x/ibc/02-client/types/tendermint"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/cosmos/cosmos-sdk/x/bank"
)

func (suite *KeeperTestSuite) createClient() {
	suite.app.Commit()
	commitID := suite.app.LastCommitID()

	suite.app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: suite.app.LastBlockHeight() + 1}})
	suite.ctx = suite.app.BaseApp.NewContext(false, abci.Header{})

	consensusState := clienttypestm.ConsensusState{
		ChainID:          testChainID,
		Height:           uint64(commitID.Version),
		Root:             commitment.NewRoot(commitID.Hash),
		ValidatorSet:     suite.valSet,
		NextValidatorSet: suite.valSet,
	}

	_, err := suite.app.IBCKeeper.ClientKeeper.CreateClient(suite.ctx, testClient, testClientType, consensusState)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) updateClient() {
	// always commit and begin a new block on updateClient
	suite.app.Commit()
	commitID := suite.app.LastCommitID()

	suite.app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: suite.app.LastBlockHeight() + 1}})
	suite.ctx = suite.app.BaseApp.NewContext(false, abci.Header{})

	state := clienttypestm.ConsensusState{
		ChainID: testChainID,
		Height:  uint64(commitID.Version),
		Root:    commitment.NewRoot(commitID.Hash),
	}

	suite.app.IBCKeeper.ClientKeeper.SetConsensusState(suite.ctx, testClient, state)
	suite.app.IBCKeeper.ClientKeeper.SetVerifiedRoot(suite.ctx, testClient, state.GetHeight(), state.GetRoot())
}

func (suite *KeeperTestSuite) createConnection(state connection.State) {
	connection := connection.ConnectionEnd{
		State:    state,
		ClientID: testClient,
		Counterparty: connection.Counterparty{
			ClientID:     testClient,
			ConnectionID: testConnection,
			Prefix:       suite.app.IBCKeeper.ConnectionKeeper.GetCommitmentPrefix(),
		},
		Versions: connection.GetCompatibleVersions(),
	}

	suite.app.IBCKeeper.ConnectionKeeper.SetConnection(suite.ctx, testConnection, connection)
}

func (suite *KeeperTestSuite) createChannel(portID string, chanID string, connID string, counterpartyPort string, counterpartyChan string, state channel.State) {
	ch := channel.Channel{
		State:    state,
		Ordering: testChannelOrder,
		Counterparty: channel.Counterparty{
			PortID:    counterpartyPort,
			ChannelID: counterpartyChan,
		},
		ConnectionHops: []string{connID},
		Version:        testChannelVersion,
	}

	suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, portID, chanID, ch)
}

func (suite *KeeperTestSuite) queryProof(key []byte) (proof commitment.Proof, height int64) {
	res := suite.app.Query(abci.RequestQuery{
		Path:  fmt.Sprintf("store/%s/key", ibctypes.StoreKey),
		Data:  key,
		Prove: true,
	})

	height = res.Height
	proof = commitment.Proof{
		Proof: res.Proof,
	}

	return
}

func (suite *KeeperTestSuite) TestSendRegisterIBCAccount() {
	err := suite.app.IBCKeeper.IbcaccountKeeper.CreateInterchainAccount(suite.ctx, testPort1, testChannel1, testSalt)
	suite.Error(err) // channel does not exist

	suite.createChannel(testPort1, testChannel1, testConnection, testPort2, testChannel2, channel.OPEN)
	err = suite.app.IBCKeeper.IbcaccountKeeper.CreateInterchainAccount(suite.ctx, testPort1, testChannel1, testSalt)
	suite.Error(err) // next send sequence not found

	nextSeqSend := uint64(1)
	suite.app.IBCKeeper.ChannelKeeper.SetNextSequenceSend(suite.ctx, testPort1, testChannel1, nextSeqSend)
	err = suite.app.IBCKeeper.IbcaccountKeeper.CreateInterchainAccount(suite.ctx, testPort1, testChannel1, testSalt)
	suite.NoError(err)

	packetCommitment := suite.app.IBCKeeper.ChannelKeeper.GetPacketCommitment(suite.ctx, testPort1, testChannel1, nextSeqSend)
	suite.NotNil(packetCommitment)
}

func (suite *KeeperTestSuite) TestReceiveRegisterIBCAccount() {
	suite.createChannel(testPort1, testChannel1, testConnection, testPort2, testChannel2, channel.OPEN)
	err := suite.app.IBCKeeper.IbcaccountKeeper.RegisterIBCAccount(suite.ctx, testPort1, testChannel1, testSalt)
	suite.NoError(err)

	identifier := fmt.Sprintf("%s/%s", testPort1, testChannel1)
	hash := tmhash.NewTruncated()
	expectedAddress := hash.Sum([]byte(identifier + testSalt))

	// Interchain account should be created with expected address
	createdAccount := suite.app.AccountKeeper.GetAccount(suite.ctx, expectedAddress)
	suite.NotNil(createdAccount)
	// Sequence should be 1
	suite.Equal(uint64(1), createdAccount.GetSequence())
	// Public key should be nil
	suite.Equal(crypto.PubKey(nil), createdAccount.GetPubKey())
}

func (suite *KeeperTestSuite) TestReceivePacket() {
	packetSeq := uint64(1)
	packetTimeout := uint64(100)

	packetDataBz := []byte("invaliddata")
	packet := channel.NewPacket(packetSeq, packetTimeout, testPort2, testChannel2, testPort2, testChannel1, packetDataBz)
	packetCommitmentPath := channel.KeyPacketCommitment(testPort2, testChannel2, packetSeq)

	suite.app.IBCKeeper.ChannelKeeper.SetPacketCommitment(suite.ctx, testPort2, testChannel2, packetSeq, []byte("invalidcommitment"))
	suite.updateClient()
	proofPacket, proofHeight := suite.queryProof(packetCommitmentPath)

	suite.createChannel(testPort2, testChannel1, testConnection, testPort2, testChannel2, channel.OPEN)
	err := suite.app.IBCKeeper.TransferKeeper.ReceivePacket(suite.ctx, packet, proofPacket, uint64(proofHeight))
	suite.Error(err) // invalid port id

	packet.DestinationPort = testPort1
	suite.createChannel(testPort1, testChannel1, testConnection, testPort2, testChannel2, channel.OPEN)
	err = suite.app.IBCKeeper.TransferKeeper.ReceivePacket(suite.ctx, packet, proofPacket, uint64(proofHeight))
	suite.Error(err) // packet membership verification failed due to invalid counterparty packet commitment

	suite.app.IBCKeeper.ChannelKeeper.SetPacketCommitment(suite.ctx, testPort2, testChannel2, packetSeq, packetDataBz)
	suite.updateClient()
	proofPacket, proofHeight = suite.queryProof(packetCommitmentPath)
	err = suite.app.IBCKeeper.TransferKeeper.ReceivePacket(suite.ctx, packet, proofPacket, uint64(proofHeight))
	suite.Error(err) // invalid packet data

	registerIAAccountPacketData := types.RegisterIBCAccountPacketData{
		Salt: testSalt,
	}
	packetDataBz, _ = suite.cdc.MarshalBinaryBare(registerIAAccountPacketData)
	packet = channel.NewPacket(packetSeq, packetTimeout, testPort2, testChannel2, testPort1, testChannel1, packetDataBz)

	suite.app.IBCKeeper.ChannelKeeper.SetPacketCommitment(suite.ctx, testPort2, testChannel2, packetSeq, packetDataBz)
	suite.updateClient()
	proofPacket, proofHeight = suite.queryProof(packetCommitmentPath)
	err = suite.app.IBCKeeper.IbcaccountKeeper.ReceivePacket(suite.ctx, packet, proofPacket, uint64(proofHeight))
	suite.NoError(err) // successfully executed

	identifier := fmt.Sprintf("%s/%s", testPort2, testChannel2)
	hash := tmhash.NewTruncated()
	expectedAddress := hash.Sum([]byte(identifier + testSalt))

	// Interchain account should be created with expected address
	createdAccount := suite.app.AccountKeeper.GetAccount(suite.ctx, expectedAddress)
	suite.NotNil(createdAccount)
	// Sequence should be 1
	suite.Equal(uint64(1), createdAccount.GetSequence())
	// Public key should be nil
	suite.Equal(crypto.PubKey(nil), createdAccount.GetPubKey())

	// Add some test coin to created interchain account
	err = createdAccount.SetCoins(sdk.Coins{
		sdk.Coin{Amount: sdk.NewInt(1000), Denom: "test"},
	})
	suite.NoError(err)
	suite.app.AccountKeeper.SetAccount(suite.ctx, createdAccount)

	sendFromIBCAccountMsg := bank.NewMsgSend(expectedAddress, []byte("normalAccount"), sdk.Coins{
		sdk.Coin{Amount: sdk.NewInt(500), Denom: "test"},
	})
	tx := types.InterchainAccountTx{Msgs: []sdk.Msg{sendFromIBCAccountMsg}}
	txBytes, err := suite.cdc.MarshalBinaryBare(tx)
	runTxPacketData := types.RunTxPacketData{TxBytes: txBytes}
	packetDataBz, _ = suite.cdc.MarshalBinaryBare(runTxPacketData)
	packet = channel.NewPacket(packetSeq, packetTimeout, testPort2, testChannel2, testPort1, testChannel1, packetDataBz)

	suite.app.IBCKeeper.ChannelKeeper.SetPacketCommitment(suite.ctx, testPort2, testChannel2, packetSeq, packetDataBz)
	suite.updateClient()
	proofPacket, proofHeight = suite.queryProof(packetCommitmentPath)
	err = suite.app.IBCKeeper.IbcaccountKeeper.ReceivePacket(suite.ctx, packet, proofPacket, uint64(proofHeight))
	suite.NoError(err) // successfully executed

	// Interchain Account should have sent some asset
	createdAccount = suite.app.AccountKeeper.GetAccount(suite.ctx, expectedAddress)
	suite.NotNil(createdAccount)
	suite.Equal(sdk.NewCoin("test", sdk.NewInt(500)).String(), createdAccount.GetCoins()[0].String())

	// Normal account should have received some asset
	normalAccount := suite.app.AccountKeeper.GetAccount(suite.ctx, []byte("normalAccount"))
	suite.NotNil(normalAccount)
	suite.Equal(sdk.NewCoin("test", sdk.NewInt(500)).String(), normalAccount.GetCoins()[0].String())

	// Normal account can't send msgs via IBC run tx
	sendFromNormalAccountMsg := bank.NewMsgSend([]byte("normalAccount"), expectedAddress, sdk.Coins{
		sdk.Coin{Amount: sdk.NewInt(500), Denom: "test"},
	})
	tx = types.InterchainAccountTx{Msgs: []sdk.Msg{sendFromNormalAccountMsg}}
	txBytes, err = suite.cdc.MarshalBinaryBare(tx)
	runTxPacketData = types.RunTxPacketData{TxBytes: txBytes}
	packetDataBz, _ = suite.cdc.MarshalBinaryBare(runTxPacketData)
	packet = channel.NewPacket(packetSeq, packetTimeout, testPort2, testChannel2, testPort1, testChannel1, packetDataBz)

	suite.app.IBCKeeper.ChannelKeeper.SetPacketCommitment(suite.ctx, testPort2, testChannel2, packetSeq, packetDataBz)
	suite.updateClient()
	proofPacket, proofHeight = suite.queryProof(packetCommitmentPath)
	err = suite.app.IBCKeeper.IbcaccountKeeper.ReceivePacket(suite.ctx, packet, proofPacket, uint64(proofHeight))
	suite.Error(err) // unauthorized
}
