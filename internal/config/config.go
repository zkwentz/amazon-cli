package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Auth         AuthConfig        `json:"auth"`
	Defaults     DefaultsConfig    `json:"defaults"`
	RateLimiting RateLimitConfig   `json:"rate_limiting"`
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

type RateLimitConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".amazon-cli", "config.json")
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = GetConfigPath()
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{
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

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
