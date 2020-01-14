package keeper_test

import (
	"fmt"
	clienttypestm "github.com/cosmos/cosmos-sdk/x/ibc/02-client/types/tendermint"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
	err := suite.app.IBCKeeper.IbcaccountKeeper.RegisterIBCAccount(suite.ctx, testPort1, testChannel1, testSalt)
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
