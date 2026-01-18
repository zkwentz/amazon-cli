package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config represents the CLI configuration structure
type Config struct {
	Auth         AuthConfig        `json:"auth"`
	Defaults     DefaultsConfig    `json:"defaults"`
	RateLimiting RateLimitConfig   `json:"rate_limiting"`
}

// AuthConfig stores authentication credentials
type AuthConfig struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// DefaultsConfig stores user preferences
type DefaultsConfig struct {
	AddressID    string `json:"address_id"`
	PaymentID    string `json:"payment_id"`
	OutputFormat string `json:"output_format"`
}

// RateLimitConfig controls request rate limiting
type RateLimitConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

// GetConfigPath returns the path to the config file, respecting the --config flag.
// If --config flag is set, it returns that path.
// Otherwise, it returns the default path: ~/.amazon-cli/config.json
func GetConfigPath() string {
	// Check if config flag was set
	configPath := viper.GetString("config")
	if configPath != "" {
		return configPath
	}

	// Return default path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir is not available
		return ".amazon-cli/config.json"
	}

	return filepath.Join(homeDir, ".amazon-cli", "config.json")
}

// LoadConfig loads the configuration from the specified path.
// If the file doesn't exist, it returns an empty config with default values.
func LoadConfig(path string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default config on first run
		return &Config{
			Auth: AuthConfig{},
			Defaults: DefaultsConfig{
				OutputFormat: "json",
			},
			RateLimiting: RateLimitConfig{
				MinDelayMs: 1000,
				MaxDelayMs: 5000,
				MaxRetries: 3,
			},
		}, nil
	}

	// Read the config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to the specified path.
// Creates the directory if it doesn't exist.
func SaveConfig(config *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}

	return nil
}
