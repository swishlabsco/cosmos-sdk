package counter

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgAdd{}, "cosmos-sdk/MsgAdd", nil)
	cdc.RegisterConcrete(MsgSubtract{}, "cosmos-sdk/MsgSubtract", nil)
}

var cdcEmpty = wire.NewCodec()
