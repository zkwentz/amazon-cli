package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Config represents the entire configuration structure
type Config struct {
	Auth         AuthConfig        `json:"auth"`
	Defaults     DefaultsConfig    `json:"defaults"`
	RateLimiting RateLimitConfig   `json:"rate_limiting"`
}

// AuthConfig holds authentication tokens and expiry
type AuthConfig struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// DefaultsConfig holds default settings
type DefaultsConfig struct {
	AddressID    string `json:"address_id"`
	PaymentID    string `json:"payment_id"`
	OutputFormat string `json:"output_format"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

// GetDefaultConfig returns a config with sensible defaults
func GetDefaultConfig() *Config {
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

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".amazon-cli")
	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	// Create directory if it doesn't exist
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, err
	}

	// If file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return GetDefaultConfig(), nil
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
	// Create directory if it doesn't exist
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file with restricted permissions
	return os.WriteFile(path, data, 0600)
}
