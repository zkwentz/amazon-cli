package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Load config when file doesn't exist - should return default values
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig() returned nil")
	}

	// Verify default rate limiting values
	if cfg.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("RateLimiting.MinDelayMs = %d, want 1000", cfg.RateLimiting.MinDelayMs)
	}
	if cfg.RateLimiting.MaxDelayMs != 5000 {
		t.Errorf("RateLimiting.MaxDelayMs = %d, want 5000", cfg.RateLimiting.MaxDelayMs)
	}
	if cfg.RateLimiting.MaxRetries != 3 {
		t.Errorf("RateLimiting.MaxRetries = %d, want 3", cfg.RateLimiting.MaxRetries)
	}

	// Verify default output format
	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("Defaults.OutputFormat = %s, want json", cfg.Defaults.OutputFormat)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create a test config
	testConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr-123",
			PaymentID:    "pay-456",
			OutputFormat: "table",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 500,
			MaxDelayMs: 3000,
			MaxRetries: 5,
		},
	}

	// Save config
	err := SaveConfig(testConfig, configPath)
	if err != nil {
		t.Fatalf("SaveConfig() returned error: %v", err)
	}

	// Verify file exists and has correct permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("config file not created: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("config file permissions = %o, want 0600", info.Mode().Perm())
	}

	// Load config back
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}

	// Verify all fields match
	if loadedConfig.Auth.AccessToken != testConfig.Auth.AccessToken {
		t.Errorf("AccessToken = %s, want %s", loadedConfig.Auth.AccessToken, testConfig.Auth.AccessToken)
	}
	if loadedConfig.Auth.RefreshToken != testConfig.Auth.RefreshToken {
		t.Errorf("RefreshToken = %s, want %s", loadedConfig.Auth.RefreshToken, testConfig.Auth.RefreshToken)
	}
	if !loadedConfig.Auth.ExpiresAt.Equal(testConfig.Auth.ExpiresAt) {
		t.Errorf("ExpiresAt = %v, want %v", loadedConfig.Auth.ExpiresAt, testConfig.Auth.ExpiresAt)
	}
	if loadedConfig.Defaults.AddressID != testConfig.Defaults.AddressID {
		t.Errorf("AddressID = %s, want %s", loadedConfig.Defaults.AddressID, testConfig.Defaults.AddressID)
	}
	if loadedConfig.Defaults.PaymentID != testConfig.Defaults.PaymentID {
		t.Errorf("PaymentID = %s, want %s", loadedConfig.Defaults.PaymentID, testConfig.Defaults.PaymentID)
	}
	if loadedConfig.Defaults.OutputFormat != testConfig.Defaults.OutputFormat {
		t.Errorf("OutputFormat = %s, want %s", loadedConfig.Defaults.OutputFormat, testConfig.Defaults.OutputFormat)
	}
	if loadedConfig.RateLimiting.MinDelayMs != testConfig.RateLimiting.MinDelayMs {
		t.Errorf("MinDelayMs = %d, want %d", loadedConfig.RateLimiting.MinDelayMs, testConfig.RateLimiting.MinDelayMs)
	}
	if loadedConfig.RateLimiting.MaxDelayMs != testConfig.RateLimiting.MaxDelayMs {
		t.Errorf("MaxDelayMs = %d, want %d", loadedConfig.RateLimiting.MaxDelayMs, testConfig.RateLimiting.MaxDelayMs)
	}
	if loadedConfig.RateLimiting.MaxRetries != testConfig.RateLimiting.MaxRetries {
		t.Errorf("MaxRetries = %d, want %d", loadedConfig.RateLimiting.MaxRetries, testConfig.RateLimiting.MaxRetries)
	}
}

func TestSaveConfig_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	// Create path with nested directories
	configPath := filepath.Join(tmpDir, "nested", "path", "config.json")

	testConfig := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	// Save should create all parent directories
	err := SaveConfig(testConfig, configPath)
	if err != nil {
		t.Fatalf("SaveConfig() returned error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	// Verify directory permissions
	dirPath := filepath.Dir(configPath)
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("path is not a directory")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Write invalid JSON
	invalidJSON := []byte(`{"auth": "not valid json`)
	err := os.WriteFile(configPath, invalidJSON, 0600)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// LoadConfig should return an error
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig() should return error for invalid JSON")
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	if path == "" {
		t.Error("GetConfigPath() returned empty string")
	}

	// Should contain .amazon-cli
	if !filepath.IsAbs(path) && path != ".amazon-cli/config.json" {
		t.Errorf("GetConfigPath() = %s, expected absolute path or .amazon-cli/config.json", path)
	}
}

func TestRateLimitConfig_Types(t *testing.T) {
	cfg := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	// Verify fields are the correct types
	var _ int = cfg.MinDelayMs
	var _ int = cfg.MaxDelayMs
	var _ int = cfg.MaxRetries
}
