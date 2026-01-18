package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AuthConfig holds authentication credentials
type AuthConfig struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// DefaultsConfig holds user default preferences
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

// Config represents the application configuration
type Config struct {
	Auth         AuthConfig      `json:"auth"`
	Defaults     DefaultsConfig  `json:"defaults"`
	RateLimiting RateLimitConfig `json:"rate_limiting"`
}

// SaveConfig saves the configuration to the specified path
func SaveConfig(config *Config, path string) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with restrictive permissions (0600 = rw-------)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
