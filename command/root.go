package command

import (
	"context"
	"fmt"

	"github.com/juancwu/gopack/config"
	"github.com/spf13/cobra"
)

const (
	cfgFile    = ".gonttrc"
	cfgType    = "yaml"
	timeFormat = "20060102150405"
	timezone   = "UTC"
)

func Execute() error {
	var showVersion bool

	rootCmd := &cobra.Command{
		Use:           "gop",
		Short:         "A simple go package installer",
		Args:          cobra.ArbitraryArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if version flag was provided
			if showVersion {
				fmt.Println(version)
				return nil
			}

			// If no arguments, show help
			if len(args) == 0 {
				return cmd.Help()
			}

			// Check if the first argument is a known subcommand
			scriptName := args[0]
			for _, subcmd := range cmd.Commands() {
				if subcmd.Name() == scriptName || subcmd.HasAlias(scriptName) {
					// Let Cobra handle it if this is a known command
					return nil
				}
			}

			// Not a command, try to run it as a script
			cfg, err := config.LoadConfig("")
			if err != nil {
				// Config not found, display help
				return cmd.Help()
			}

			// Check if the script exists
			if _, ok := cfg.Scripts[scriptName]; !ok {
				// Script not found, display help
				return cmd.Help()
			}

			// Run the script
			return config.RunScript(cfg, scriptName, args[1:])
		},
	}

	// Add global version flag
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version information")

	rootCmd.AddCommand(get())
	rootCmd.AddCommand(run())
	rootCmd.AddCommand(list())
	rootCmd.AddCommand(update())
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(create())

	return rootCmd.ExecuteContext(context.Background())
}
