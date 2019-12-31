package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/types"
	"testing"
)

func TestDecoding(t *testing.T) {
	var msgs []sdk.Msg
	interchainAccountTx := types.InterchainAccountTx{Msgs: msgs}

	txBytes, err := types.ModuleCdc.MarshalBinaryBare(interchainAccountTx)
	if err != nil {
		fmt.Println("1. err = ", err)
	}

	packetData := types.RunTxPacketData{
		TxBytes: txBytes,
	}
	pdBytes, err := types.ModuleCdc.MarshalBinaryBare(packetData)
	if err != nil {
		fmt.Println("1-1. err = ", err)
	}

	packet := channel.NewPacket(
		11,
		123123,
		"sourcePort",
		"sourceChannel",
		"destinationPort",
		"destinationChannel",
		pdBytes,
	)
	recvPacket(packet)
}

func recvPacket(packet channelexported.PacketI) {
	fmt.Printf("packet = %v\n", packet)

	//var data types.InterchainAccountTx
	//var data types.RunTxPacketData
	var data types.IbcPacketData
	//var data interface{}

	err := types.ModuleCdc.UnmarshalBinaryBare(packet.GetData(), &data)
	if err != nil {
		fmt.Println("2. err = ", err)
	}

	switch packetData := data.(type) {
	case types.RegisterIBCAccountPacketData:
		fmt.Println("RegisterIBCAccountPacketData")
		fmt.Println("packetData=", packetData)
	case types.RunTxPacketData:
		fmt.Println("RunTxPacketData")
		fmt.Println("packetData=", packetData)
	}

	//// interface test
	//packetBytes, err := types.ModuleCdc.MarshalBinaryBare(packet)
	//if err != nil {
	//	fmt.Println("3. err = ", err)
	//}
	//fmt.Println("packetBytes=", packetBytes)
	//
	////var rePacket interface{}
	//var rePacket channelexported.PacketI
	//err = types.ModuleCdc.UnmarshalBinaryBare(packetBytes, &rePacket)
	//if err != nil {
	//	fmt.Println("4. err = ", err)
	//}
	//fmt.Println("rePacket=", rePacket)
}
