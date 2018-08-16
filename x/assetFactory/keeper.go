package assetFactory

//Keeper : asset keeper
type Keeper struct {
	am AssetPegMapper
}

//NewKeeper : return a new keeper
func NewKeeper(am AssetPegMapper) Keeper {
	return Keeper{am: am}
}
