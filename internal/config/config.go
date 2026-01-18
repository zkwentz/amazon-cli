package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Auth         AuthConfig         `json:"auth"`
	Defaults     DefaultsConfig     `json:"defaults"`
	RateLimiting RateLimitingConfig `json:"rate_limiting"`
}

type AuthConfig struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type DefaultsConfig struct {
	AddressID    string `json:"address_id"`
	PaymentID    string `json:"payment_id"`
	OutputFormat string `json:"output_format"`
}

type RateLimitingConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

func GetConfigPath() string {
	if configPath := os.Getenv("AMAZON_CLI_CONFIG"); configPath != "" {
		return configPath
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".amazon-cli", "config.json")
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = GetConfigPath()
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	// Return default config if file doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{
			Auth: AuthConfig{},
			Defaults: DefaultsConfig{
				OutputFormat: "json",
			},
			RateLimiting: RateLimitingConfig{
				MinDelayMs: 1000,
				MaxDelayMs: 5000,
				MaxRetries: 3,
			},
		}, nil
	}

	// Read and parse existing config
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

func SaveConfig(config *Config, path string) error {
	if path == "" {
		path = GetConfigPath()
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file with restrictive permissions
	return os.WriteFile(path, data, 0600)
}

func (c *Config) IsAuthenticated() bool {
	return c.Auth.AccessToken != "" || c.Auth.RefreshToken != ""
}

func (c *Config) IsTokenExpired() bool {
	if c.Auth.ExpiresAt.IsZero() {
		return true
	}
	return time.Now().After(c.Auth.ExpiresAt)
}
