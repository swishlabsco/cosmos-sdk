package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/cmd/thorchain/app"
	"github.com/cosmos/cosmos-sdk/server"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "thorchaind",
		Short:             "Thorchain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, app.ThorchainAppInit(),
		server.ConstructAppCreator(newApp, "thorchain"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "thorchain"))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewThorchainApp(logger, db, traceStore, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	gApp := app.NewThorchainApp(logger, db, traceStore)
	return gApp.ExportAppStateAndValidators()
}
