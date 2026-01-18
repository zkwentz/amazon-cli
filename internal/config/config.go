package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/zkwentz/amazon-cli/internal/ratelimit"
)

// Config represents the application configuration
type Config struct {
	Auth         AuthConfig                `json:"auth"`
	Defaults     DefaultsConfig            `json:"defaults"`
	RateLimiting ratelimit.RateLimitConfig `json:"rate_limiting"`
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

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
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
	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Write to file with restricted permissions
	return os.WriteFile(path, data, 0600)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Auth: AuthConfig{},
		Defaults: DefaultsConfig{
			OutputFormat: "json",
		},
		RateLimiting: ratelimit.RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}
}
