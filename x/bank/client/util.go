package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
)

// build the sendTx msg
func BuildMsg(from sdk.Address, to sdk.Address, coins sdk.Coins) sdk.Msg {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}

//BuildIssueAssetMsg : butild the issueAssetTx
func BuildIssueAssetMsg(from sdk.Address, to sdk.Address, assetPeg sdk.BaseAssetPeg) sdk.Msg {
	issueAsset := bank.NewIssueAsset(from, to, sdk.AssetPegWallet{assetPeg})
	msg := bank.NewMsgIssueAssets([]bank.IssueAsset{issueAsset})
	return msg
}
