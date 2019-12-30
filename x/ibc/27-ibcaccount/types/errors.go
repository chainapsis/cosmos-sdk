package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// interchain-account error codes
const (
	CodeAccountAlreadyExist  sdk.CodeType = 301
	CodeUnsupportedChianType sdk.CodeType = 302
	CodeContentIsEmpty       sdk.CodeType = 303
)

// ErrAccountAlreadyExist implements sdk.Error
func ErrAccountAlreadyExist(codespace sdk.CodespaceType, account string) error {
	return sdkerrors.New(
		string(codespace),
		uint32(CodeAccountAlreadyExist),
		fmt.Sprintf("account (%s) already exists", account),
	)
}

// ErrUnsupportedChainType implements sdk.Error
func ErrUnsupportedChainType(codespace sdk.CodespaceType, chainType string) error {
	return sdkerrors.New(
		string(codespace),
		uint32(CodeUnsupportedChianType),
		fmt.Sprintf("type (%s) is unsupported chain type", chainType),
	)
}

// ErrContentIsEmpty implements sdk.Error
func ErrContentIsEmpty(codespace sdk.CodespaceType, content string) error {
	return sdkerrors.New(
		string(codespace),
		uint32(CodeContentIsEmpty),
		fmt.Sprintf("%s is empty", content),
	)
}
