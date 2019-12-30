package test

import (
	//"github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/27-ibcaccount/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/28-test/types"
)

var (
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper
)
