package main

import (
	"fmt"
	"log"

	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/wiring"
	"github.com/spf13/cobra"
)

var (
	version    string = "0"
	commitHash string = "0"
)

const (
	flagConfigFilePath = "config-file-path"
)

func server() *cobra.Command {
	command := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFilePath, err := cmd.Flags().GetString(flagConfigFilePath)
			if err != nil {
				return err
			}

			app, cleanup, err := wiring.InitializeAppServer(configs.ConfigFilePath(configFilePath))
			if err != nil {
				return err
			}
			defer cleanup()

			app.Start()
			return nil
		},
	}

	command.Flags().String(flagConfigFilePath, "", "If provided, will use the provided config file")

	return command
}

func main() {
	rootCommand := &cobra.Command{
		Version: fmt.Sprintf("%s-%s", version, commitHash),
	}
	rootCommand.AddCommand(
		server(),
	)

	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
