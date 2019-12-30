package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const CosmosSdkChainType = "cosmos-sdk"

type InterchainAccountTx struct {
	Msgs []sdk.Msg
}

type IbcPacketData interface {
	ValidateBasic() sdk.Error
}
