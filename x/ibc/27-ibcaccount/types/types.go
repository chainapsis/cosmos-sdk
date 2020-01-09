package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const CosmosSdkChainType = "cosmos-sdk"

type InterchainAccountTx struct {
	Msgs []sdk.Msg
}

// This is just any(interface{}).
// This is used as interface for interchain account packet data.
// Currently, packet data is fundamentally just a byte slice.
// But, we need type swtiching for identifying type of packet data.
// TODO: This will be removed when ics-026 (routing module) is implemented.
type InterchainAccountPacketData interface{}
