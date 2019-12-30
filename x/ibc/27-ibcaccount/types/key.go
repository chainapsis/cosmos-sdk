package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// SubModuleName defines the IBC interchain-account name
	SubModuleName = "interchainaccount"

	// StoreKey is the store key string for IBC interchain-account
	StoreKey = SubModuleName

	// RouterKey is the message route for IBC interchain-account
	RouterKey = SubModuleName

	// QuerierRoute is the querier route for IBC interchain-account
	QuerierRoute = SubModuleName

	// DefaultCodespace is the default error codespace for IBC interchain-account
	DefaultCodespace sdk.CodespaceType = SubModuleName
)

func GetIdentifier(sourcePort, sourceChannel string) string {
	return fmt.Sprintf("%s/%s", sourcePort, sourceChannel)
}
