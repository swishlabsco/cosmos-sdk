package counter

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// name to identify transaction types
const MsgType = "counter"

var cdc = wire.NewCodec()

//-----------------------------------------------------------
// MsgAdd
type MsgAdd struct {
	SenderAddress sdk.AccAddress
	Amount        int //  Amount to add
}

func NewMsgAdd(sender sdk.AccAddress, amount int) MsgAdd {
	return MsgAdd{
		SenderAddress: sender,
		Amount:        amount,
	}
}

// Implements Msg.
func (msg MsgAdd) Type() string { return MsgType }

// Implements Msg.
func (msg MsgAdd) ValidateBasic() sdk.Error {
	if msg.Amount <= 0 {
		return ErrInvalidAmount(DefaultCodespace, msg.Amount)
	}
	return nil
}

func (msg MsgAdd) String() string {
	return fmt.Sprintf("MsgAdd{%s, %s, %s, %v}", msg.Amount)
}

// Implements Msg.
func (msg MsgAdd) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgAdd) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgAdd) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SenderAddress}
}

//-----------------------------------------------------------
// MsgSubtract
type MsgSubtract struct {
	SenderAddress sdk.AccAddress
	Amount        int //  Amount to add
}

func NewMsgSubtract(sender sdk.AccAddress, amount int) MsgSubtract {
	return MsgSubtract{
		SenderAddress: sender,
		Amount:        amount,
	}
}

// Implements Msg.
func (msg MsgSubtract) Type() string { return MsgType }

// Implements Msg.
func (msg MsgSubtract) ValidateBasic() sdk.Error {
	if msg.Amount <= 0 {
		return ErrInvalidAmount(DefaultCodespace, msg.Amount)
	}
	return nil
}

func (msg MsgSubtract) String() string {
	return fmt.Sprintf("MsgSubtract{%s, %s, %s, %v}", msg.Amount)
}

// Implements Msg.
func (msg MsgSubtract) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgSubtract) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSubtract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SenderAddress}
}
