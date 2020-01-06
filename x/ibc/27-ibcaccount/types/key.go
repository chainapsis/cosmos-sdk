package types

import (
	"fmt"
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
)

func GetIdentifier(sourcePort, sourceChannel string) string {
	return fmt.Sprintf("%s/%s", sourcePort, sourceChannel)
}
