package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Config represents the complete application configuration
type Config struct {
	Auth AuthConfig `json:"auth"`
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".amazon-cli", "config.json")
}

// LoadConfig reads configuration from the specified path
// If the file doesn't exist, it returns a default empty config
func LoadConfig(path string) (*Config, error) {
	// Expand ~ to home directory if present
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// If path is empty, use default
	if path == "" {
		path = DefaultConfigPath()
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		return &Config{}, nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON using a temporary struct to handle empty time strings
	var raw struct {
		Auth struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresAt    string `json:"expires_at"`
		} `json:"auth"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Build config with proper time parsing
	config := &Config{
		Auth: AuthConfig{
			AccessToken:  raw.Auth.AccessToken,
			RefreshToken: raw.Auth.RefreshToken,
		},
	}

	// Parse time if not empty
	if raw.Auth.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, raw.Auth.ExpiresAt)
		if err != nil {
			// If parsing fails, leave as zero time (treat as invalid)
			config.Auth.ExpiresAt = time.Time{}
		} else {
			config.Auth.ExpiresAt = expiresAt
		}
	}

	return config, nil
}

// SaveConfig writes configuration to the specified path with 0600 permissions
func SaveConfig(config *Config, path string) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Expand ~ to home directory if present
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// If path is empty, use default
	if path == "" {
		path = DefaultConfigPath()
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write file with 0600 permissions (read/write for owner only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// IsAuthenticated checks if the user has valid authentication
func (c *Config) IsAuthenticated() bool {
	if c == nil || c.Auth.AccessToken == "" {
		return false
	}

	// Check if token is expired
	if time.Now().After(c.Auth.ExpiresAt) {
		return false
	}

	return true
}

// ClearAuth clears all authentication data
func (c *Config) ClearAuth() {
	c.Auth.AccessToken = ""
	c.Auth.RefreshToken = ""
	c.Auth.ExpiresAt = time.Time{}
}
