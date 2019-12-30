package types

import (
	"fmt"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

// IBC ibcaccount events
const (
	AttributeKeyRegisterIBCAccount = "register_ibc_account"
	AttributeRunTx                 = "run_tx"
)

// IBC ibcaccount events vars
var (
	EventTypeRegisterIBCAccount = MsgRegisterIBCAccount{}.Type()
	EventTypeRunTx              = MsgRunTx{}.Type()

	AttributeValueCategory = fmt.Sprintf("%s_%s", ibctypes.ModuleName, SubModuleName)
)
