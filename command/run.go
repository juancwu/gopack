package command

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/juancwu/gopack/config"
	"github.com/spf13/cobra"
)

func run() *cobra.Command {
	var configPath string
	var listScripts bool
	var initConfig bool

	runCmd := &cobra.Command{
		Use:     "run [script]",
		Short:   "Run a script from gopack.json",
		Long:    "Run a script defined in the gopack.json configuration file",
		Example: "gopack run build",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle init flag
			if initConfig {
				return handleInitConfig()
			}

			// Load config
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			// Handle list flag
			if listScripts {
				return handleListScripts(cfg)
			}

			// Handle running a script
			if len(args) == 0 {
				return fmt.Errorf("no script specified")
			}

			scriptName := args[0]
			scriptArgs := args[1:]

			return config.RunScript(cfg, scriptName, scriptArgs)
		},
	}

	runCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file (default: gopack.json in current directory)")
	runCmd.Flags().BoolVarP(&listScripts, "list", "l", false, "List available scripts")
	runCmd.Flags().BoolVarP(&initConfig, "init", "i", false, "Initialize a new gopack.json configuration file")

	return runCmd
}

func handleInitConfig() error {
	if err := config.CreateDefaultConfig(); err != nil {
		return fmt.Errorf("failed to initialize configuration: %v", err)
	}
	log.Info("created gopack.json configuration file")
	return nil
}

func handleListScripts(cfg *config.Config) error {
	scripts := config.ListScripts(cfg)
	if len(scripts) == 0 {
		fmt.Println("No scripts defined in configuration")
		return nil
	}

	fmt.Println("Available scripts:")
	for _, script := range scripts {
		fmt.Printf("  %s: %s\n", script, cfg.Scripts[script])
	}
	return nil
}
