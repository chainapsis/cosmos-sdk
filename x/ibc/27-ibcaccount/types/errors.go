package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// interchain-account errors
var (
	ErrAccountAlreadyExist  = sdkerrors.Register(SubModuleName, 1, "account already exist")
	ErrUnsupportedChainType = sdkerrors.Register(SubModuleName, 2, "unsupported chain type")
	ErrContentIsEmpty       = sdkerrors.Register(SubModuleName, 3, "content is empty")
	ErrInvalidRoute         = sdkerrors.Register(SubModuleName, 4, "invalid route")
)
