package counter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgAdd:
			return handleMsgAdd(ctx, msg, k)
		case MsgSubtract:
			return handleMsgSubtract(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in counter module").Result()
		}
	}
}

func handleMsgAdd(ctx sdk.Context, msg MsgAdd, k Keeper) sdk.Result {

	// Nonbasic validation, eg, max counter

	err := k.Add(ctx, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}

}

func handleMsgSubtract(ctx sdk.Context, msg MsgSubtract, k Keeper) sdk.Result {

	// Nonbasic validation, eg, max counter

	err := k.Subtract(ctx, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}

}
