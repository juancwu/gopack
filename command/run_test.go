package command

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/juancwu/gopack/config"
)

func TestRunCommand(t *testing.T) {
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

	t.Run("Init flag", func(t *testing.T) {
		// Get the run command
		cmd := run()

		// Set init flag
		cmd.SetArgs([]string{"--init"})

		// Run command
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		// Check if config file was created
		if _, err := os.Stat(filepath.Join(tmpDir, config.DefaultConfigName)); os.IsNotExist(err) {
			t.Fatal("Config file was not created")
		}
	})

	// Create test config for remaining tests
	testConfig := config.Config{
		Scripts: map[string]string{
			"hello": "echo Hello World",
		},
	}

	// Write custom config file
	configPath := filepath.Join(tmpDir, config.DefaultConfigName)
	writeTestConfig(t, configPath, testConfig)

	t.Run("List flag", func(t *testing.T) {
		// Get a fresh run command
		cmd := run()

		// Set list flag
		cmd.SetArgs([]string{"--list"})

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run command
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Check output
		if !strings.Contains(output, "hello") {
			t.Errorf("Expected output to contain script 'hello', got: %s", output)
		}
	})

	t.Run("No script specified", func(t *testing.T) {
		// Get a fresh run command
		cmd := run()

		// Set no args
		cmd.SetArgs([]string{})

		// Run command (should fail)
		if err := cmd.Execute(); err == nil {
			t.Fatal("Expected error when no script specified, got nil")
		}
	})

	t.Run("Non-existent script", func(t *testing.T) {
		// Get a fresh run command
		cmd := run()

		// Set non-existent script
		cmd.SetArgs([]string{"nonexistent"})

		// Run command (should fail)
		if err := cmd.Execute(); err == nil {
			t.Fatal("Expected error when script doesn't exist, got nil")
		}
	})
}

func writeTestConfig(t *testing.T, configPath string, cfg config.Config) {
	t.Helper()

	// Convert config to JSON
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
}
