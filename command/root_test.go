package command

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/juancwu/gopack/config"
	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Save and restore working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a temporary directory and change to it
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create test command for evaluation
	var rootCmd *cobra.Command

	createTestRootCmd := func() *cobra.Command {
		rootCmd = &cobra.Command{
			Use:   "gop",
			Short: "A simple go package installer",
			Args:  cobra.ArbitraryArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
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

		rootCmd.AddCommand(get())
		rootCmd.AddCommand(run())

		return rootCmd
	}

	t.Run("No arguments", func(t *testing.T) {
		cmd := createTestRootCmd()
		cmd.SetArgs([]string{})

		// Capture output
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)

		// Run command should show help
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		output := buf.String()
		if len(output) == 0 {
			t.Error("Expected help output, got empty string")
		}
	})

	t.Run("Script from config", func(t *testing.T) {
		// Create config with test script
		testConfig := config.Config{
			Scripts: map[string]string{
				"echo": "echo test-output",
			},
		}

		// Write test config
		configPath := filepath.Join(tmpDir, config.DefaultConfigName)
		data, err := json.MarshalIndent(testConfig, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Create command
		cmd := createTestRootCmd()
		cmd.SetArgs([]string{"echo"})

		// Execute command (should run the script)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}
	})
}
