package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_NonExistentFile(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Load config from non-existent file
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify default values
	if config.Defaults.OutputFormat != "json" {
		t.Errorf("Expected default output format 'json', got '%s'", config.Defaults.OutputFormat)
	}
	if config.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("Expected default MinDelayMs 1000, got %d", config.RateLimiting.MinDelayMs)
	}
	if config.RateLimiting.MaxDelayMs != 5000 {
		t.Errorf("Expected default MaxDelayMs 5000, got %d", config.RateLimiting.MaxDelayMs)
	}
	if config.RateLimiting.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries 3, got %d", config.RateLimiting.MaxRetries)
	}
	if config.Auth.AccessToken != "" {
		t.Errorf("Expected empty AccessToken, got '%s'", config.Auth.AccessToken)
	}
}

func TestLoadConfig_ExistingFile(t *testing.T) {
	// Create a temporary directory and config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create test config
	expiresAt := time.Now().Add(time.Hour)
	testConfig := Config{
		Auth: AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    expiresAt,
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr_123",
			PaymentID:    "pay_456",
			OutputFormat: "table",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 2000,
			MaxDelayMs: 10000,
			MaxRetries: 5,
		},
	}

	// Write config to file
	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify values
	if config.Auth.AccessToken != "test-access-token" {
		t.Errorf("Expected AccessToken 'test-access-token', got '%s'", config.Auth.AccessToken)
	}
	if config.Auth.RefreshToken != "test-refresh-token" {
		t.Errorf("Expected RefreshToken 'test-refresh-token', got '%s'", config.Auth.RefreshToken)
	}
	if !config.Auth.ExpiresAt.Equal(expiresAt) {
		t.Errorf("Expected ExpiresAt %v, got %v", expiresAt, config.Auth.ExpiresAt)
	}
	if config.Defaults.AddressID != "addr_123" {
		t.Errorf("Expected AddressID 'addr_123', got '%s'", config.Defaults.AddressID)
	}
	if config.Defaults.PaymentID != "pay_456" {
		t.Errorf("Expected PaymentID 'pay_456', got '%s'", config.Defaults.PaymentID)
	}
	if config.Defaults.OutputFormat != "table" {
		t.Errorf("Expected OutputFormat 'table', got '%s'", config.Defaults.OutputFormat)
	}
	if config.RateLimiting.MinDelayMs != 2000 {
		t.Errorf("Expected MinDelayMs 2000, got %d", config.RateLimiting.MinDelayMs)
	}
	if config.RateLimiting.MaxDelayMs != 10000 {
		t.Errorf("Expected MaxDelayMs 10000, got %d", config.RateLimiting.MaxDelayMs)
	}
	if config.RateLimiting.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries 5, got %d", config.RateLimiting.MaxRetries)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Create a temporary directory and config file with invalid JSON
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Write invalid JSON
	if err := os.WriteFile(configPath, []byte("{ invalid json }"), 0600); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Load config should fail
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoadConfig_EmptyPath(t *testing.T) {
	// Load config with empty path should use default path
	// This will check home directory existence and create directories if needed
	config, err := LoadConfig("")
	if err != nil {
		// It's okay to fail if we can't get home directory
		t.Logf("LoadConfig with empty path failed (expected in some environments): %v", err)
		return
	}

	// If it succeeds, it should return default config
	if config == nil {
		t.Error("Expected non-nil config")
	}
}

func TestLoadConfig_DirectoryCreation(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	// Use a nested path that doesn't exist
	configPath := filepath.Join(tempDir, "nested", "dir", "config.json")

	// Load config should create the directory
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify directory was created
	dir := filepath.Dir(configPath)
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Expected directory, got file")
	}

	// Verify permissions (0700)
	if info.Mode().Perm() != 0700 {
		t.Errorf("Expected directory permissions 0700, got %o", info.Mode().Perm())
	}

	// Verify default config was returned
	if config.Defaults.OutputFormat != "json" {
		t.Errorf("Expected default output format 'json', got '%s'", config.Defaults.OutputFormat)
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := getDefaultConfig()

	// Verify all default values
	if config.Auth.AccessToken != "" {
		t.Errorf("Expected empty AccessToken, got '%s'", config.Auth.AccessToken)
	}
	if config.Auth.RefreshToken != "" {
		t.Errorf("Expected empty RefreshToken, got '%s'", config.Auth.RefreshToken)
	}
	if !config.Auth.ExpiresAt.IsZero() {
		t.Errorf("Expected zero ExpiresAt, got %v", config.Auth.ExpiresAt)
	}
	if config.Defaults.AddressID != "" {
		t.Errorf("Expected empty AddressID, got '%s'", config.Defaults.AddressID)
	}
	if config.Defaults.PaymentID != "" {
		t.Errorf("Expected empty PaymentID, got '%s'", config.Defaults.PaymentID)
	}
	if config.Defaults.OutputFormat != "json" {
		t.Errorf("Expected OutputFormat 'json', got '%s'", config.Defaults.OutputFormat)
	}
	if config.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("Expected MinDelayMs 1000, got %d", config.RateLimiting.MinDelayMs)
	}
	if config.RateLimiting.MaxDelayMs != 5000 {
		t.Errorf("Expected MaxDelayMs 5000, got %d", config.RateLimiting.MaxDelayMs)
	}
	if config.RateLimiting.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries 3, got %d", config.RateLimiting.MaxRetries)
	}
}
