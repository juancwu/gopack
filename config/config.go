package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

// Config represents the structure of the gopack.json configuration file
type Config struct {
	Scripts map[string]string `json:"scripts,omitempty"`
}

type ScriptError struct {
	err error
}

func (e ScriptError) Error() string {
	return e.err.Error()
}

const (
	// DefaultConfigName is the name of the configuration file
	DefaultConfigName = "gopack.json"
)

// LoadConfig loads the configuration from the specified file path
// If no path is provided, it looks for gopack.json in the current directory
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		// Try to find config in current directory
		path = DefaultConfigName
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// CreateDefaultConfig creates a default configuration file in the current directory
func CreateDefaultConfig() error {
	// Check if file already exists
	if _, err := os.Stat(DefaultConfigName); err == nil {
		return fmt.Errorf("config file already exists")
	}

	config := Config{
		Scripts: map[string]string{
			"build": "go build",
			"test":  "go test ./...",
			"run":   "go run main.go",
		},
	}

	// Convert to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(DefaultConfigName, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// RunScript executes a script from the configuration
func RunScript(config *Config, scriptName string, args []string) error {
	// Check if script exists
	script, ok := config.Scripts[scriptName]
	if !ok {
		return fmt.Errorf("script not found: %s", scriptName)
	}

	// Append any additional arguments
	if len(args) > 0 {
		script = fmt.Sprintf("%s %s", script, strings.Join(args, " "))
	}

	log.Debug("running script", "name", scriptName, "command", script)

	// Create command
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run command
	err := cmd.Run()
	if err != nil {
		return ScriptError{err: err}
	}
	return nil
}

// ListScripts returns a list of available scripts in the configuration
func ListScripts(config *Config) []string {
	scripts := make([]string, 0, len(config.Scripts))
	for name := range config.Scripts {
		scripts = append(scripts, name)
	}
	return scripts
}
