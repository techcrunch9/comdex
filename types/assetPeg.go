package types

import (
	"encoding/hex"
	"errors"

	cmn "github.com/tendermint/tmlibs/common"
)

//PegHash : reference address of asset peg
type PegHash cmn.HexBytes

//AssetPeg : peg issued against assets
// type AssetPeg struct {
// 	PegHash       PegHash `json:"pegHash"`
// 	OblHash       string  `json:"oblHash"`
// 	AssetType     string  `json:"assetType"`
// 	AssetQuantity int64   `json:"assetQuantity"`
// 	QuantityType  string  `json:"quantityType"`
// }

//AssetPegWallet : A wallet of AssetPegTokens
type AssetPegWallet []BaseAssetPeg

//BaseAssetPeg : base asset type
type BaseAssetPeg struct {
	PegHash       PegHash `json:"pegHash"`
	OblHash       string  `json:"oblHash"`
	AssetType     string  `json:"assetType"`
	AssetQuantity int64   `json:"assetQuantity"`
	QuantityType  string  `json:"quantityType"`
	OwnerAddress  string  `json:"ownerAddress"`
}

//NewBaseAssetPegWithPegHash : return a base asset peg with peg hash
func NewBaseAssetPegWithPegHash(pegHash PegHash) BaseAssetPeg {
	return BaseAssetPeg{
		PegHash: pegHash,
	}
}

//AssetPeg : comdex asset interface
type AssetPeg interface{}

//GetAssetPegHashHex : convert string to hex peg hash
func GetAssetPegHashHex(pegHashStr string) (pegHash PegHash, err error) {
	if len(pegHashStr) == 0 {
		return pegHash, errors.New("must use provide pegHash")
	}
	bz, err := hex.DecodeString(pegHashStr)
	if err != nil {
		return nil, err
	}
	return PegHash(bz), nil
}
