package cli

import (
	flag "github.com/spf13/pflag"
)

//noLint
const (
	FlagTo            = "to"
	FlagAmount        = "amount"
	FlagOblHash       = "oblHash"
	FlagAssetType     = "assetType"
	FlagAssetQuantity = "assetQuantity"
	FlagQuantityType  = "quantityType"
)

var (
	fsTo            = flag.NewFlagSet("", flag.ContinueOnError)
	fsAmount        = flag.NewFlagSet("", flag.ContinueOnError)
	fsOblHash       = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetType     = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetQuantity = flag.NewFlagSet("", flag.ContinueOnError)
	fsQuantityType  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsTo.String(flagTo, "", "Address to send coins")
	fsAmount.String(flagAmount, "", "Amount of coins to send")
	fsOblHash.String(FlagOblHash, "", "Hash of the OBl doccument of the asset")
	fsAssetType.String(FlagAssetType, "", "Type of the asset")
	fsAssetQuantity.String(FlagAssetQuantity, "", "Quantity of the assent in integer")
	fsQuantityType.String(FlagQuantityType, "", "The unit of the qunatity")
}
