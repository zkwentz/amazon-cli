package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the complete configuration structure for amazon-cli
type Config struct {
	Auth         AuthConfig        `json:"auth"`
	Defaults     DefaultsConfig    `json:"defaults"`
	RateLimiting RateLimitConfig   `json:"rate_limiting"`
}

// AuthConfig holds authentication-related configuration
type AuthConfig struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// DefaultsConfig holds default values for common operations
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

// LoadConfig loads configuration from the specified path.
// If the file doesn't exist, it returns a default configuration.
// If path is empty, it uses the default path (~/.amazon-cli/config.json).
func LoadConfig(path string) (*Config, error) {
	// Use default path if not provided
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = filepath.Join(homeDir, ".amazon-cli", "config.json")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Return default config if file doesn't exist (first run)
		return getDefaultConfig(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to stat config file: %w", err)
	}

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &config, nil
}

// getDefaultConfig returns a Config with default values
func getDefaultConfig() *Config {
	return &Config{
		Auth: AuthConfig{
			AccessToken:  "",
			RefreshToken: "",
			ExpiresAt:    time.Time{},
		},
		Defaults: DefaultsConfig{
			AddressID:    "",
			PaymentID:    "",
			OutputFormat: "json",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}
}
