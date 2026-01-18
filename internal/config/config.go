package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Auth         AuthConfig        `json:"auth"`
	Defaults     DefaultsConfig    `json:"defaults"`
	RateLimiting RateLimitConfig   `json:"rate_limiting"`
}

// AuthConfig holds authentication credentials
type AuthConfig struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// DefaultsConfig holds default settings
type DefaultsConfig struct {
	AddressID    string `json:"address_id"`
	PaymentID    string `json:"payment_id"`
	OutputFormat string `json:"output_format"`
}

// RateLimitConfig holds rate limiting settings
type RateLimitConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	// Expand home directory if necessary
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, path[1:])
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	// Return default config if file doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read and parse config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves configuration to the specified path
func SaveConfig(config *Config, path string) error {
	// Expand home directory if necessary
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, path[1:])
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file with restrictive permissions
	return os.WriteFile(path, data, 0600)
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
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
	}
}

// GetConfigPath returns the default config path
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".amazon-cli/config.json"
	}
	return filepath.Join(home, ".amazon-cli", "config.json")
}
