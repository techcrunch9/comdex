package app

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	crypto "github.com/tendermint/go-crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

//GenesisState : State to Unmarshal
type GenesisState struct {
	Accounts  []GenesisAccount   `json:"accounts"`
	Assets    []GenesisAssetPeg  `json:"assets"`
	StakeData stake.GenesisState `json:"stake"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.Address `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
}

//NewGenesisAccount : returns a new genesis state account
func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
	}
}

//NewGenesisAccountI : new genesis account from already existing account
func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address: acc.GetAddress(),
		Coins:   acc.GetCoins(),
	}
}

//ToAccount : convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
}

//GenesisAssetPeg : genesis state of assets
type GenesisAssetPeg struct {
	PegHash       sdk.PegHash `json:"pegHash"`
	OblHash       string      `json:"oblHash"`
	AssetType     string      `json:"assetType"`
	AssetQuantity int64       `json:"assetQuantity"`
	QuantityType  string      `json:"quantityType"`
	OwnerAddress  string      `json:"ownerAddress"`
}

//NewGenesisAsset : returns a new genesis state account
func NewGenesisAsset(assetPeg *sdk.BaseAssetPeg) GenesisAssetPeg {
	return GenesisAssetPeg{
		PegHash:       assetPeg.PegHash,
		OblHash:       assetPeg.OblHash,
		AssetType:     assetPeg.AssetType,
		AssetQuantity: assetPeg.AssetQuantity,
		QuantityType:  assetPeg.QuantityType,
		OwnerAddress:  assetPeg.OwnerAddress,
	}
}

var (
	flagName       = "name"
	flagClientHome = "home-client"
	flagOWK        = "owk"

	// bonded tokens given to genesis validators/accounts
	freeFermionVal  = int64(100)
	freeFermionsAcc = int64(50)
)

//ComdexAssetAppInit : get app init parameters for server init command
func ComdexAssetAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagName, "", "validator moniker, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         ComdexAssetAppGenTx,
		AppGenState:      ComdexAssetAppGenStateJSON,
	}
}

//ComdexAssetGenTx : simple genesis tx
type ComdexAssetGenTx struct {
	Name    string        `json:"name"`
	Address sdk.Address   `json:"address"`
	PubKey  crypto.PubKey `json:"pub_key"`
}

//ComdexAssetAppGenTx : Generate a comdex asset genesis transaction with flags
func ComdexAssetAppGenTx(cdc *wire.Codec, pk crypto.PubKey) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	clientRoot := viper.GetString(flagClientHome)
	overwrite := viper.GetBool(flagOWK)
	name := viper.GetString(flagName)
	if name == "" {
		return nil, nil, tmtypes.GenesisValidator{}, errors.New("Must specify --name (validator moniker)")
	}

	var addr sdk.Address
	var secret string
	addr, secret, err = server.GenerateSaveCoinKey(clientRoot, name, "1234567890", overwrite)
	if err != nil {
		return
	}
	mm := map[string]string{"secret": secret}
	var bz []byte
	bz, err = cdc.MarshalJSON(mm)
	if err != nil {
		return
	}
	cliPrint = json.RawMessage(bz)
	appGenTx, _, validator, err = ComdexAssetAppGenTxNF(cdc, pk, addr, name, overwrite)
	return
}

// ComdexAssetAppGenTxNF : Generate a comdex asset genesis transaction without flags
func ComdexAssetAppGenTxNF(cdc *wire.Codec, pk crypto.PubKey, addr sdk.Address, name string, overwrite bool) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	var bz []byte
	comdexAssetGenTx := ComdexAssetGenTx{
		Name:    name,
		Address: addr,
		PubKey:  pk,
	}
	bz, err = wire.MarshalJSONIndent(cdc, comdexAssetGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  freeFermionVal,
	}
	return
}

//ComdexAssetAppGenState : Create the core parameters for genesis initialization for comdex asset
// note that the pubkey input is this machines pubkey
func ComdexAssetAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {

	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// start with the default staking genesis state
	stakeData := stake.DefaultGenesisState()

	// get genesis flag account information
	genaccs := make([]GenesisAccount, len(appGenTxs))
	for i, appGenTx := range appGenTxs {

		var genTx ComdexAssetGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}

		// create the genesis account, give'm few steaks and a buncha token with there name
		accAuth := auth.NewBaseAccountWithAddress(genTx.Address)

		accAuth.Coins = sdk.Coins{
			{
				Denom:  "comdex" + genTx.Name,
				Amount: 1000,
			},
			{
				Denom:  "steak",
				Amount: freeFermionsAcc,
			},
		}
		acc := NewGenesisAccount(&accAuth)
		genaccs[i] = acc
		stakeData.Pool.LooseUnbondedTokens += freeFermionsAcc // increase the supply

		// add the validator
		if len(genTx.Name) > 0 {
			desc := stake.NewDescription(genTx.Name, "", "", "")
			validator := stake.NewValidator(genTx.Address, genTx.PubKey, desc)
			validator.PoolShares = stake.NewBondedShares(sdk.NewRat(freeFermionVal))
			stakeData.Validators = append(stakeData.Validators, validator)

			// pool logic
			stakeData.Pool.BondedTokens += freeFermionVal
			stakeData.Pool.BondedShares = sdk.NewRat(stakeData.Pool.BondedTokens)
		}
	}

	//Generate empty asset tokens
	genesisAssetPegs := make([]GenesisAssetPeg, 1024)
	for i := 0; i < 1024; i++ {
		pegHash, err := sdk.GetAssetPegHashHex(strconv.Itoa(i))
		if err == nil {
			genesisAssetPegs[i] = NewGenesisAsset(&sdk.BaseAssetPeg{
				PegHash:       pegHash,
				OblHash:       "",
				AssetType:     "",
				AssetQuantity: 0,
				QuantityType:  "",
			})
		}
	}

	// create the final app state
	genesisState = GenesisState{
		Accounts:  genaccs,
		Assets:    genesisAssetPegs,
		StakeData: stakeData,
	}
	return
}

//ComdexAssetAppGenStateJSON : ComdexAssetAppGenState but with JSON
func ComdexAssetAppGenStateJSON(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	// create the final app state
	genesisState, err := ComdexAssetAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}
