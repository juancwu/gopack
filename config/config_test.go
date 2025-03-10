package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a test config file
	testConfig := Config{
		Scripts: map[string]string{
			"test": "go test",
		},
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	// Write to temp file
	configPath := filepath.Join(tmpDir, DefaultConfigName)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Test loading with path
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify loaded config
	if len(cfg.Scripts) != len(testConfig.Scripts) {
		t.Errorf("Expected %d scripts, got %d", len(testConfig.Scripts), len(cfg.Scripts))
	}

	// Verify test script value
	if cfg.Scripts["test"] != testConfig.Scripts["test"] {
		t.Errorf("Expected test script '%s', got '%s'", testConfig.Scripts["test"], cfg.Scripts["test"])
	}

	// Test loading non-existent file
	_, err = LoadConfig(filepath.Join(tmpDir, "nonexistent.json"))
	if err == nil {
		t.Error("Expected error when loading non-existent file, got nil")
	}
}

func TestCreateDefaultConfig(t *testing.T) {
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

	// Create default config
	if err := CreateDefaultConfig(); err != nil {
		t.Fatalf("CreateDefaultConfig failed: %v", err)
	}

	// Check if file exists
	configPath := filepath.Join(tmpDir, DefaultConfigName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}

	// Load and verify the created config
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load created config: %v", err)
	}

	// Basic verification
	if len(cfg.Scripts) != 3 {
		t.Errorf("Expected 3 default scripts, got %d", len(cfg.Scripts))
	}

	// Verify default scripts
	expectedScripts := map[string]string{
		"build": "go build",
		"test":  "go test ./...",
		"run":   "go run main.go",
	}

	for name, command := range expectedScripts {
		if cfg.Scripts[name] != command {
			t.Errorf("Expected %s script to be '%s', got '%s'", name, command, cfg.Scripts[name])
		}
	}

	// Test creating config when it already exists
	err = CreateDefaultConfig()
	if err == nil {
		t.Error("Expected error when creating config that already exists, got nil")
	}
}

func TestListScripts(t *testing.T) {
	// Create test config
	cfg := &Config{
		Scripts: map[string]string{
			"build": "go build",
			"test":  "go test",
			"run":   "go run",
		},
	}

	// List scripts
	scripts := ListScripts(cfg)

	// Check length
	if len(scripts) != 3 {
		t.Errorf("Expected 3 scripts, got %d", len(scripts))
	}

	// Check script names
	scriptMap := make(map[string]bool)
	for _, script := range scripts {
		scriptMap[script] = true
	}

	expectedScripts := []string{"build", "test", "run"}
	for _, expected := range expectedScripts {
		if !scriptMap[expected] {
			t.Errorf("Expected script %s not found", expected)
		}
	}
}

func TestRunScript(t *testing.T) {
	// Create test config with echo command for easy testing
	cfg := &Config{
		Scripts: map[string]string{
			"echo": "echo test",
		},
	}

	// Test running a valid script
	err := RunScript(cfg, "echo", nil)
	if err != nil {
		t.Errorf("RunScript failed for valid script: %v", err)
	}

	// Test running a non-existent script
	err = RunScript(cfg, "nonexistent", nil)
	if err == nil {
		t.Error("Expected error when running non-existent script, got nil")
	}
}
