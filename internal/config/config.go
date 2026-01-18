package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the application configuration
type Config struct {
	Auth      AuthConfig      `json:"auth"`
	Defaults  DefaultsConfig  `json:"defaults"`
	RateLimit RateLimitConfig `json:"rate_limiting"`
}

// AuthConfig holds authentication tokens
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

// RateLimitConfig holds rate limiting settings
type RateLimitConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

// LoadConfig loads the configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	// Expand home directory if needed
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Return default config if file doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return defaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the specified path
func SaveConfig(config *Config, path string) error {
	// Expand home directory if needed
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the default config path
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".amazon-cli/config.json"
	}
	return filepath.Join(home, ".amazon-cli", "config.json")
}

// defaultConfig returns a config with default values
func defaultConfig() *Config {
	return &Config{
		Auth: AuthConfig{},
		Defaults: DefaultsConfig{
			OutputFormat: "json",
		},
		RateLimit: RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}
}
