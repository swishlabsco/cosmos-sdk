package coins

import (
	"encoding/json"
	"fmt"

	crypto "github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SendMsg - high level transaction of the coin module
type SendMsg struct {
	FromAddress  crypto.Address `json:"address"`
	ToAddress    crypto.Address `json:"address"`
	Coins        Coins          `json:"coins"`
}

// NewSendMsg - construct arbitrary multi-in, multi-out send msg.
func NewSendMsg(in []Input, out []Output) SendMsg {
	return SendMsg{Inputs: in, Outputs: out}
}

// Implements Msg.
func (msg SendMsg) Type() string { return "bank/send" } // TODO: "bank/send"

// Implements Msg.
func (msg SendMsg) ValidateBasic() sdk.Error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside

	if len(msg.FromAddress) == 0 {
		return ErrInvalidAddress(msg.FromAddress.String())
	}
	if len(msg.ToAddress) == 0 {
		return ErrInvalidAddress(msg.ToAddress.String())
	}
	if !msg.Coins.IsValid() {
		return ErrInvalidCoins(msg.Coins.String())
	}
	if !msg.Coins.IsPositive() {
		return ErrInvalidCoins(msg.Coins.String())
	}
	return nil
}

func (msg SendMsg) String() string {
	return fmt.Sprintf("SendMsg{%v->%v:%v}", msg.FromAddress, msg.ToAddress, msg.Coins)
}

// Implements Msg.
func (msg SendMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg SendMsg) GetSigners() []crypto.Address {
	addrs := make([]crypto.Address, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

//----------------------------------------
// IssueMsg

// IssueMsg - high level transaction of the coin module
type IssueMsg struct {
	Banker  crypto.Address `json:"banker"`
	Outputs []Output       `json:"outputs"`
}

// NewIssueMsg - construct arbitrary multi-in, multi-out send msg.
func NewIssueMsg(banker crypto.Address, out []Output) IssueMsg {
	return IssueMsg{Banker: banker, Outputs: out}
}

// Implements Msg.
func (msg IssueMsg) Type() string { return "bank" } // TODO: "bank/send"

// Implements Msg.
func (msg IssueMsg) ValidateBasic() sdk.Error {
	// XXX
	if len(msg.Outputs) == 0 {
		return ErrNoOutputs().Trace("")
	}
	for _, out := range msg.Outputs {
		if err := out.ValidateBasic(); err != nil {
			return err.Trace("")
		}
	}
	return nil
}

func (msg IssueMsg) String() string {
	return fmt.Sprintf("IssueMsg{%v#%v}", msg.Banker, msg.Outputs)
}

// Implements Msg.
func (msg IssueMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg IssueMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg IssueMsg) GetSigners() []crypto.Address {
	return []crypto.Address{msg.Banker}
}

//----------------------------------------
// Input

type Input struct {
	Address  crypto.Address `json:"address"`
	Coins    sdk.Coins      `json:"coins"`
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() sdk.Error {
	if len(in.Address) == 0 {
		return ErrInvalidAddress(in.Address.String())
	}
	if in.Sequence < 0 {
		return ErrInvalidSequence("negative sequence")
	}
	if !in.Coins.IsValid() {
		return ErrInvalidCoins(in.Coins.String())
	}
	if !in.Coins.IsPositive() {
		return ErrInvalidCoins(in.Coins.String())
	}
	return nil
}

func (in Input) String() string {
	return fmt.Sprintf("Input{%v,%v}", in.Address, in.Coins)
}

// NewInput - create a transaction input, used with SendMsg
func NewInput(addr crypto.Address, coins sdk.Coins) Input {
	input := Input{
		Address: addr,
		Coins:   coins,
	}
	return input
}