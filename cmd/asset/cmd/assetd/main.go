package main

import (
	"encoding/json"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/cmd/asset/app"
	"github.com/cosmos/cosmos-sdk/server"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "assetd",
		Short:             "Comdex Asset Chain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, app.ComdexAssetAppInit(),
		server.ConstructAppCreator(newApp, "comdexAsset"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "comdexAsset"))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "CA", app.DefaultNodeHome)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB) abci.Application {
	return app.NewComdexAssetApp(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	gapp := app.NewComdexAssetApp(logger, db)
	return gapp.ExportAppStateAndValidators()
}
