package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clientexported "github.com/cosmos/cosmos-sdk/x/ibc/02-client/exported"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	ibcaccount "github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"reflect"
	"testing"
)

// define constants used for testing
const (
	testChainID    = "test-chain-id"
	testClient     = "test-client"
	testClientType = clientexported.Tendermint

	testSeq                = 123
	testConnection         = "test-connection"
	testSourcePort         = "test-source-port"
	testDestinationPort    = "test-destination-port"
	testSourceChannel      = "test-source-channel"
	testDestinationChannel = "test-destination-channel"

	testChannelOrdered = channel.ORDERED
	testChannelVersion = "1.0"
)

var (
	testAddr1 = sdk.AccAddress([]byte("testaddr1"))
)

type KeeperTestSuite struct {
	suite.Suite

	cdc             *codec.Codec
	counterpartyCdc *codec.Codec

	ctx    sdk.Context
	app    *simapp.SimApp
	valSet *tmtypes.ValidatorSet
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	//todo: more consider
	suite.counterpartyCdc = app.Codec()

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.app = app

	privVal := tmtypes.NewMockPV()

	validator := tmtypes.NewValidator(privVal.GetPubKey(), 1)
	suite.valSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	suite.createClient()
	suite.createConnection(connection.OPEN)
}

func (suite *KeeperTestSuite) TestUnmarshalPacketData_RunTx() {
	var msgs []sdk.Msg
	interchainAccountTx := types.InterchainAccountTx{Msgs: msgs}

	txBytes, err := suite.counterpartyCdc.MarshalBinaryBare(interchainAccountTx)
	suite.NoError(err)

	RunTxPD := types.RunTxPacketData{
		TxBytes: txBytes,
	}
	runTxPdBytes, err := suite.counterpartyCdc.MarshalBinaryBare(RunTxPD)
	suite.NoError(err)

	RunTxPacket := channel.NewPacket(
		testSeq,
		uint64(suite.ctx.BlockHeight())+1000,
		testSourcePort,
		testSourceChannel,
		testDestinationPort,
		testDestinationChannel,
		runTxPdBytes,
	)

	runTxData, err := suite.app.IBCKeeper.IbcaccountKeeper.UnmarshalPacketData(RunTxPacket)
	suite.NoError(err)
	suite.Equal(reflect.TypeOf(ibcaccount.RunTxPacketData{}), reflect.TypeOf(runTxData))
	suite.Equal(RunTxPD, runTxData)
}

func (suite *KeeperTestSuite) TestUnmarshalPacketData_RegisterIbcAccount() {
	RegisterIbcAccountPD := types.RegisterIBCAccountPacketData{
		Salt: "test salt",
	}
	registerIbcAccountPdBytes, err := suite.counterpartyCdc.MarshalBinaryBare(RegisterIbcAccountPD)
	suite.NoError(err)

	RegisterIbcAccountPacket := channel.NewPacket(
		testSeq,
		uint64(suite.ctx.BlockHeight())+1000,
		testSourcePort,
		testSourceChannel,
		testDestinationPort,
		testDestinationChannel,
		registerIbcAccountPdBytes,
	)

	registerData, err := suite.app.IBCKeeper.IbcaccountKeeper.UnmarshalPacketData(RegisterIbcAccountPacket)
	suite.NoError(err)
	suite.Equal(reflect.TypeOf(ibcaccount.RegisterIBCAccountPacketData{}), reflect.TypeOf(registerData))
	suite.Equal(RegisterIbcAccountPD, registerData)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
