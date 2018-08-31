package lcd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type baseSendBody struct {
	LocalAccountName string `json:"name"`
	Password         string `json:"password"`
	ChainID          string `json:"chain_id"`
	AccountNumber    int64  `json:"account_number"`
	Sequence         int64  `json:"sequence"`
	Gas              int64  `json:"gas"`
}

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

func extractRequest(w http.ResponseWriter, r *http.Request, cdc *wire.Codec) (baseSendBody, []byte, error) {
	var m baseSendBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return baseSendBody{}, nil, err
	}
	err = cdc.UnmarshalJSON(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return baseSendBody{}, nil, err
	}
	return m, body, nil
}

func setupContext(w http.ResponseWriter, ctx context.CoreContext, m baseSendBody, from sdk.AccAddress) (context.CoreContext, error) {
	// add gas to context
	ctx = ctx.WithGas(m.Gas)

	// add chain-id to context
	ctx = ctx.WithChainID(m.ChainID)

	//add account number and sequence
	ctx, err := EnsureAccountNumber(ctx, m.AccountNumber, from)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return ctx, err
	}
	ctx, err = EnsureSequence(ctx, m.Sequence, from)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return ctx, err
	}
	return ctx, nil
}

func getFromAddress(w http.ResponseWriter, kb keys.Keybase, localAccountName string) (sdk.AccAddress, error) {
	info, err := kb.Get(localAccountName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return sdk.AccAddress{}, err
	}

	from := sdk.AccAddress(info.GetPubKey().Address())
	return from, nil
}

func processMsg(w http.ResponseWriter, ctx context.CoreContext, localAccountName string, password string, cdc *wire.Codec, msg sdk.Msg) ([]byte, error) {
	//sign
	txBytes, err := ctx.SignAndBuild(localAccountName, password, []sdk.Msg{msg}, cdc)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return nil, err
	}

	// send
	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return nil, err
	}

	output, err := wire.MarshalJSONIndent(cdc, res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return nil, err
	}
	return output, err
}

// RequestHandlerFn - http request handler to handle generic transaction api call
func RequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext, msgBuilder func(http.ResponseWriter, *wire.Codec, sdk.AccAddress, []byte) (sdk.Msg, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, body, err := extractRequest(w, r, cdc)
		if err != nil {
			return
		}

		from, err := getFromAddress(w, kb, m.LocalAccountName)
		if err != nil {
			return
		}

		ctx, err = setupContext(w, ctx, m, from)
		if err != nil {
			return
		}

		// build message
		msg, err := msgBuilder(w, cdc, from, body)
		if err != nil {
			return
		}

		output, err := processMsg(w, ctx, m.LocalAccountName, m.Password, cdc, msg)
		if err != nil {
			return
		}

		w.Write(output)
	}
}
