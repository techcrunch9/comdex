package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/bank/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOblHash       = "oblHash"
	flagAssetType     = "assetType"
	flagAssetQuantity = "assetQuantity"
	flagQuantityType  = "quantityType"
)

//IssueAssetCmd : create a init asset tx and sign it with the give key
func IssueAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Initializes asset with the given details and issues to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			from, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.GetAccAddressBech32(toStr)
			if err != nil {
				return nil
			}

			oblHashStr := viper.GetString(flagOblHash)
			assetTypeStr := viper.GetString(flagAssetType)
			assetQuantityStr := viper.GetInt64(flagAssetQuantity)
			quantityTypeStr := viper.GetString(flagQuantityType)

			assetPeg := sdk.BaseAssetPeg{
				AssetQuantity: assetQuantityStr,
				AssetType:     assetTypeStr,
				OblHash:       oblHashStr,
				QuantityType:  quantityTypeStr,
			}
			msg := client.BuildIssueAssetMsg(from, to, assetPeg)
			res, err := ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, msg, cdc)
			if err != nil {
				return err
			}
			fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
			return nil
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsAmount)
	cmd.Flags().AddFlagSet(fsOblHash)
	cmd.Flags().AddFlagSet(fsAssetType)
	cmd.Flags().AddFlagSet(fsAssetQuantity)
	cmd.Flags().AddFlagSet(fsQuantityType)
	return cmd
}
