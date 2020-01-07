package keeper_test

import (
	"fmt"
	clienttypestm "github.com/cosmos/cosmos-sdk/x/ibc/02-client/types/tendermint"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	abci "github.com/tendermint/tendermint/abci/types"
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

func (suite *KeeperTestSuite) createConnection(state connection.State) {
	conn := connection.ConnectionEnd{
		State:    state,
		ClientID: testClient,
		Counterparty: connection.Counterparty{
			ClientID:     testClient,
			ConnectionID: testConnection,
			Prefix:       suite.app.IBCKeeper.ConnectionKeeper.GetCommitmentPrefix(),
		},
		Versions: connection.GetCompatibleVersions(),
	}

	suite.app.IBCKeeper.ConnectionKeeper.SetConnection(suite.ctx, testConnection, conn)
}

func (suite *KeeperTestSuite) createChannel(portID, chanID, connID, counterpartyPort, counterpartyChan string, state channel.State) {
	ch := channel.Channel{
		State:    state,
		Ordering: testChannelOrdered,
		Counterparty: channel.Counterparty{
			PortID:    counterpartyPort,
			ChannelID: counterpartyChan,
		},
		ConnectionHops: []string{connID},
		Version:        testChannelVersion,
	}

	suite.app.IBCKeeper.ChannelKeeper.SetChannel(suite.ctx, portID, chanID, ch)
}

func (suite *KeeperTestSuite) TestRegisterIBCAccount() {
	err := suite.app.IBCKeeper.IbcaccountKeeper.RegisterIBCAccount(suite.ctx, testPort1, testChannel1, testSalt)
	suite.EqualError(err, "channel not found: test-channel1")

	suite.createChannel(testPort1, testChannel1, testConnection, testPort2, testChannel2, channel.OPEN)
	err = suite.app.IBCKeeper.IbcaccountKeeper.RegisterIBCAccount(suite.ctx, testPort1, testChannel1, testSalt)
	fmt.Println(err)
}
