//nolint
package counter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Local code type
type CodeType = sdk.CodeType

const (
	// Default slashing codespace
	DefaultCodespace sdk.CodespaceType = 14

	CodeInvalidAmount CodeType = 141
)

func ErrInvalidAmount(codespace sdk.CodespaceType, amount int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAmount, "that amount is not positive")
}
