package types

//FiatPeg : pegs issued against fiat
type FiatPeg struct {
	PegHash         string `json:"pegHash"`
	OriginatingTxID string `json:"originatingTxID"`
	RedeemingTxID   string `json:"redeemingTxID"`
}

//FiatPegWallet : A wallet of fiat peg tokens
type FiatPegWallet []FiatPeg

//TODO: comdex sort
