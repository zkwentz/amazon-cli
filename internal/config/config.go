package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuthConfig holds authentication tokens and expiration information
type AuthConfig struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
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

// Config represents the complete configuration structure
type Config struct {
	Auth         AuthConfig      `json:"auth"`
	Defaults     DefaultsConfig  `json:"defaults"`
	RateLimiting RateLimitConfig `json:"rate_limiting"`
}

// DefaultConfig returns a new Config with sensible defaults
func DefaultConfig() *Config {
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

// GetConfigPath returns the path to the config file
// If a custom path is provided, it uses that; otherwise, it uses the default path
func GetConfigPath(customPath string) (string, error) {
	if customPath != "" {
		return customPath, nil
	}

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".amazon-cli", "config.json"), nil
}

// LoadConfig loads the configuration from the specified path
// If the file doesn't exist, it returns a default configuration
func LoadConfig(path string) (*Config, error) {
	// If path is empty, use default
	if path == "" {
		var err error
		path, err = GetConfigPath("")
		if err != nil {
			return nil, err
		}
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default config if file doesn't exist (first run)
		return DefaultConfig(), nil
	}

	// Read the file
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
	// If path is empty, use default
	if path == "" {
		var err error
		path, err = GetConfigPath("")
		if err != nil {
			return err
		}
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

	// Write to file with restricted permissions (0600 = rw-------)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
