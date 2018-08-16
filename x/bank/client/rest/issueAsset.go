package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-crypto/keys"
)

var msgWireCdc = wire.NewCodec()

func init() {
	bank.RegisterWire(msgWireCdc)
}

type IssueAssetBody struct {
	Name          string `json:"name"`
	Gas           int64  `json:"gas"`
	OblHash       string `json:"oblhash"`
	AssetType     string `json:"asset_type"`
	QuantityType  string `json:"quantity_type"`
	AssetQuantity int64  `json:"asset-quantity"`
	ChainID       string `json:"chain-id"`
	AccountNumber int64  `json:"accountNumber"`
	Password      string `json:"password"`
	Sequence      int64  `json:"sequence"`
}

func IssueAssetHandlerFunction(ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var msg IssueAssetBody
		vars := mux.Vars(r)
		to := vars["address"]
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		toAddress, err := sdk.GetAccAddressBech32(to)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = msgWireCdc.UnmarshalJSON(body, &msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		info, err := kb.Get(msg.Name)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			return
		}
		// w.Write([]byte(msg.Name))
		ctx = ctx.WithFromAddressName(msg.Name)
		// from, err := ctx.GetFromAddress()
		// w.Write([]byte(from))
		ctx = ctx.WithGas(msg.Gas)
		ctx = ctx.WithAccountNumber(msg.AccountNumber)
		ctx = ctx.WithDecoder(authcmd.GetAccountDecoder(cdc))
		// , _ = ctx.NextSequence(toAddress)
		// a, _ := ctx.NextSequence(from)
		// w.Write([]byte(strconv.FormatInt(a, 10)))
		ctx = ctx.WithSequence(msg.Sequence)
		ctx = ctx.WithChainID(msg.ChainID)
		assetPeg := sdk.BaseAssetPeg{
			AssetQuantity: msg.AssetQuantity,
			AssetType:     msg.AssetType,
			OblHash:       msg.OblHash,
			QuantityType:  msg.QuantityType,
		}
		buildMsg := client.BuildIssueAssetMsg(info.PubKey.Address(), toAddress, assetPeg)
		if err != nil { // XXX rechecking same error ?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		ctx, err = context.EnsureAccountNumber(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// default to next sequence number if none provided
		ctx, err = context.EnsureSequence(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		txBytes, err := ctx.SignAndBuild(msg.Name, msg.Password, buildMsg, cdc)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// send
		res, err := ctx.BroadcastTx(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)

	}

}
