package keeper_test

import (
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
)

func (suite *KeeperTestSuite) TestOnChanOpenInit() {
	invalidOrder := channel.UNORDERED

	counterparty := channel.NewCounterparty(testPort2, testChannel2)
	err := suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenInit(suite.ctx, invalidOrder, []string{testConnection}, testPort1, testChannel1, counterparty, "")
	suite.EqualError(err, "invalid channel: channel must be ORDERED")

	err = suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenInit(suite.ctx, testChannelOrdered, []string{testConnection}, testPort1, testChannel1, counterparty, "")
	suite.EqualError(err, "invalid port: counterparty port ID doesn't match the capability key (test-port2 ≠ interchainaccount)")

	counterparty = channel.NewCounterparty(testPort1, testChannel2)
	err = suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenInit(suite.ctx, testChannelOrdered, []string{testConnection}, testPort1, testChannel1, counterparty, testChannelVersion)
	suite.EqualError(err, "invalid version: version must be blank")

	err = suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenInit(suite.ctx, testChannelOrdered, []string{testConnection}, testPort1, testChannel1, counterparty, "")
	suite.NoError(err) // successfully executed
}

func (suite *KeeperTestSuite) TestOnChanOpenTry() {
	invalidOrder := channel.UNORDERED

	counterparty := channel.NewCounterparty(testPort2, testChannel2)
	err := suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenTry(suite.ctx, invalidOrder, []string{testConnection}, testPort1, testChannel1, counterparty, "", "")
	suite.EqualError(err, "invalid channel: channel must be ORDERED")

	err = suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenTry(suite.ctx, testChannelOrdered, []string{testConnection}, testPort1, testChannel1, counterparty, "", "")
	suite.EqualError(err, "invalid port: counterparty port ID doesn't match the capability key (test-port2 ≠ interchainaccount)")

	counterparty = channel.NewCounterparty(testPort1, testChannel2)
	err = suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenTry(suite.ctx, testChannelOrdered, []string{testConnection}, testPort1, testChannel1, counterparty, testChannelVersion, "")
	suite.EqualError(err, "invalid version: version must be blank")

	err = suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenTry(suite.ctx, testChannelOrdered, []string{testConnection}, testPort1, testChannel1, counterparty, "", testChannelVersion)
	suite.EqualError(err, "invalid version: counterparty version must be blank")

	err = suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenTry(suite.ctx, testChannelOrdered, []string{testConnection}, testPort1, testChannel1, counterparty, "", "")
	suite.NoError(err) // successfully executed
}

func (suite *KeeperTestSuite) TestOnChanOpenAck() {
	err := suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenAck(suite.ctx, testPort1, testChannel1, testChannelVersion)
	suite.EqualError(err, "invalid version: version must be blank")

	err = suite.app.IBCKeeper.IbcaccountKeeper.OnChanOpenAck(suite.ctx, testPort1, testChannel1, "")
	suite.NoError(err)
}
