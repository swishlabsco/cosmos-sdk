package lcd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Ensure Account Number is added to context for transaction
func EnsureAccountNumber(ctx context.CoreContext, accountNumber int64, from sdk.AccAddress) (context.CoreContext, error) {
	if accountNumber != 0 {
		return ctx, nil
	}
	accountNumber, err := ctx.GetAccountNumber(from)
	if err != nil {
		return ctx, err
	}
	fmt.Printf("Defaulting to account number: %d\n", accountNumber)
	ctx = ctx.WithAccountNumber(accountNumber)
	return ctx, nil
}

//Ensure Sequence is added to context for transaction
func EnsureSequence(ctx context.CoreContext, sequence int64, from sdk.AccAddress) (context.CoreContext, error) {
	if sequence != 0 {
		return ctx, nil
	}
	sequence, err := ctx.NextSequence(from)
	if err != nil {
		return ctx, err
	}
	fmt.Printf("Defaulting to next sequence number: %d\n", sequence)
	ctx = ctx.WithSequence(sequence)
	return ctx, nil
}

// MakeRequestHandlerFn - create http request handler to send coins to a address with correct json type
// func MakeRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext, sendBody interface{}) {

// 	return RequestHandlerFn()
// }

// RequestHandlerFn - http request handler to send coins to a address
// func RequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext, sendBody interface{}, msgBuilder func(sdk.AccAddress, interface{}) sdk.Msg) http.HandlerFunc {
// 	ctx = ctx.WithDecoder(authcmd.GetAccountDecoder(cdc))
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var m sendBody
// 		body, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}
// 		err = msgCdc.UnmarshalJSON(body, &m)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}

// 		info, err := kb.Get(m.LocalAccountName)
// 		if err != nil {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}

// 		from := sdk.AccAddress(info.GetPubKey().Address())

// 		// build message
// 		msg := msgBuilder(from, m)

// 		if err != nil { // XXX rechecking same error ?
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}

// 		// add gas to context
// 		ctx = ctx.WithGas(m.Gas)

// 		// add chain-id to context
// 		ctx = ctx.WithChainID(m.ChainID)

// 		//add account number and sequence
// 		ctx, err = EnsureAccountNumber(ctx, m, from)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}
// 		ctx, err = EnsureSequence(ctx, m, from)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}

// 		//sign
// 		txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, []sdk.Msg{msg}, cdc)
// 		if err != nil {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}

// 		// send
// 		res, err := ctx.BroadcastTx(txBytes)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}

// 		output, err := wire.MarshalJSONIndent(cdc, res)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(err.Error()))
// 			return
// 		}

// 		w.Write(output)
// 	}
// }
