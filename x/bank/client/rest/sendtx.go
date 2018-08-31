package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	lcdhelpers "github.com/cosmos/cosmos-sdk/client/lcd/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank/client"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/accounts/{address}/send", lcdhelpers.RequestHandlerFn(cdc, kb, ctx, buildSendMsg)).Methods("POST")

}

type sendBody struct {
	Amount sdk.Coins `json:"amount"`
}

func buildSendMsg(w http.ResponseWriter, cdc *wire.Codec, from sdk.AccAddress, body []byte, routeVars map[string]string) (sdk.Msg, error) {
	var m sendBody

	bech32addr := routeVars["address"]

	to, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return nil, err
	}

	err = cdc.UnmarshalJSON(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return nil, err
	}
	msg := client.BuildMsg(from, to, m.Amount)
	return msg, nil
}
