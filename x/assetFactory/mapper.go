package assetFactory

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

//AssetPegMapper : encoder decoder for asset type
type AssetPegMapper struct {
	key   sdk.StoreKey
	proto sdk.AssetPeg
	cdc   *wire.Codec
}

//NewAssetPegMapper : returns asset mapper
func NewAssetPegMapper(cdc *wire.Codec, key sdk.StoreKey, proto sdk.AssetPeg) AssetPegMapper {
	return AssetPegMapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}
