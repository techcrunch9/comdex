package bank

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg)
		case MsgIssue:
			return handleMsgIssue(ctx, k, msg)
		case MsgIssueAssets:
			return handleMsgIssueAsset(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

//Hande MsgIssueAsset
func handleMsgIssueAsset(ctx sdk.Context, k Keeper, msg MsgIssueAssets) sdk.Result {
	tags, err := k.IssueAssetsToWallets(ctx, msg.IssueAssets)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k Keeper, msg MsgSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked

	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgIssue.
func handleMsgIssue(ctx sdk.Context, k Keeper, msg MsgIssue) sdk.Result {
	panic("not implemented yet")
}
