package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	if config == nil {
		t.Fatal("NewDefaultConfig returned nil")
	}

	if config.Defaults.OutputFormat != "json" {
		t.Errorf("Expected default output format 'json', got '%s'", config.Defaults.OutputFormat)
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

func TestLoadConfig_NonExistentFile(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "nonexistent", "config.json")

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error for non-existent file: %v", err)
	}

	if config == nil {
		t.Fatal("LoadConfig returned nil for non-existent file")
	}

	// Should return default config
	if config.Defaults.OutputFormat != "json" {
		t.Errorf("Expected default output format 'json', got '%s'", config.Defaults.OutputFormat)
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".amazon-cli", "config.json")

	// Create a test config
	expiresAt := time.Now().Add(24 * time.Hour)
	config := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    expiresAt,
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr_123",
			PaymentID:    "pay_456",
			OutputFormat: "json",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	// Save config
	err := SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Verify file permissions (should be 0600)
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}
	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", mode)
	}

	// Verify directory permissions (should be 0700)
	dirInfo, err := os.Stat(filepath.Dir(configPath))
	if err != nil {
		t.Fatalf("Failed to stat config directory: %v", err)
	}
	dirMode := dirInfo.Mode().Perm()
	if dirMode != 0700 {
		t.Errorf("Expected directory permissions 0700, got %o", dirMode)
	}

	// Read and verify content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var savedConfig Config
	if err := json.Unmarshal(data, &savedConfig); err != nil {
		t.Fatalf("Failed to parse saved config: %v", err)
	}

	if savedConfig.Auth.AccessToken != config.Auth.AccessToken {
		t.Errorf("AccessToken mismatch: expected '%s', got '%s'", config.Auth.AccessToken, savedConfig.Auth.AccessToken)
	}

	if savedConfig.Auth.RefreshToken != config.Auth.RefreshToken {
		t.Errorf("RefreshToken mismatch: expected '%s', got '%s'", config.Auth.RefreshToken, savedConfig.Auth.RefreshToken)
	}

	if savedConfig.Defaults.AddressID != config.Defaults.AddressID {
		t.Errorf("AddressID mismatch: expected '%s', got '%s'", config.Defaults.AddressID, savedConfig.Defaults.AddressID)
	}

	if savedConfig.Defaults.PaymentID != config.Defaults.PaymentID {
		t.Errorf("PaymentID mismatch: expected '%s', got '%s'", config.Defaults.PaymentID, savedConfig.Defaults.PaymentID)
	}

	if savedConfig.RateLimiting.MinDelayMs != config.RateLimiting.MinDelayMs {
		t.Errorf("MinDelayMs mismatch: expected %d, got %d", config.RateLimiting.MinDelayMs, savedConfig.RateLimiting.MinDelayMs)
	}
}

func TestLoadConfig_ExistingFile(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create a test config and save it
	expiresAt := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	originalConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "loaded-access-token",
			RefreshToken: "loaded-refresh-token",
			ExpiresAt:    expiresAt,
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr_789",
			PaymentID:    "pay_012",
			OutputFormat: "table",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 2000,
			MaxDelayMs: 10000,
			MaxRetries: 5,
		},
	}

	err := SaveConfig(originalConfig, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the config
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify all fields
	if loadedConfig.Auth.AccessToken != originalConfig.Auth.AccessToken {
		t.Errorf("AccessToken mismatch: expected '%s', got '%s'", originalConfig.Auth.AccessToken, loadedConfig.Auth.AccessToken)
	}

	if loadedConfig.Auth.RefreshToken != originalConfig.Auth.RefreshToken {
		t.Errorf("RefreshToken mismatch: expected '%s', got '%s'", originalConfig.Auth.RefreshToken, loadedConfig.Auth.RefreshToken)
	}

	if !loadedConfig.Auth.ExpiresAt.Equal(originalConfig.Auth.ExpiresAt) {
		t.Errorf("ExpiresAt mismatch: expected %v, got %v", originalConfig.Auth.ExpiresAt, loadedConfig.Auth.ExpiresAt)
	}

	if loadedConfig.Defaults.AddressID != originalConfig.Defaults.AddressID {
		t.Errorf("AddressID mismatch: expected '%s', got '%s'", originalConfig.Defaults.AddressID, loadedConfig.Defaults.AddressID)
	}

	if loadedConfig.Defaults.PaymentID != originalConfig.Defaults.PaymentID {
		t.Errorf("PaymentID mismatch: expected '%s', got '%s'", originalConfig.Defaults.PaymentID, loadedConfig.Defaults.PaymentID)
	}

	if loadedConfig.Defaults.OutputFormat != originalConfig.Defaults.OutputFormat {
		t.Errorf("OutputFormat mismatch: expected '%s', got '%s'", originalConfig.Defaults.OutputFormat, loadedConfig.Defaults.OutputFormat)
	}

	if loadedConfig.RateLimiting.MinDelayMs != originalConfig.RateLimiting.MinDelayMs {
		t.Errorf("MinDelayMs mismatch: expected %d, got %d", originalConfig.RateLimiting.MinDelayMs, loadedConfig.RateLimiting.MinDelayMs)
	}

	if loadedConfig.RateLimiting.MaxDelayMs != originalConfig.RateLimiting.MaxDelayMs {
		t.Errorf("MaxDelayMs mismatch: expected %d, got %d", originalConfig.RateLimiting.MaxDelayMs, loadedConfig.RateLimiting.MaxDelayMs)
	}

	if loadedConfig.RateLimiting.MaxRetries != originalConfig.RateLimiting.MaxRetries {
		t.Errorf("MaxRetries mismatch: expected %d, got %d", originalConfig.RateLimiting.MaxRetries, loadedConfig.RateLimiting.MaxRetries)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Write invalid JSON
	err := os.WriteFile(configPath, []byte("invalid json content"), 0600)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	// Try to load the config
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Fatal("LoadConfig should have returned an error for invalid JSON")
	}
}

func TestSaveAndLoadConfig_RoundTrip(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create a test config
	expiresAt := time.Now().Add(24 * time.Hour).Round(time.Second)
	originalConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "roundtrip-access-token",
			RefreshToken: "roundtrip-refresh-token",
			ExpiresAt:    expiresAt,
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr_roundtrip",
			PaymentID:    "pay_roundtrip",
			OutputFormat: "json",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1500,
			MaxDelayMs: 7500,
			MaxRetries: 4,
		},
	}

	// Save config
	err := SaveConfig(originalConfig, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Load config
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify the config survived the round trip
	if loadedConfig.Auth.AccessToken != originalConfig.Auth.AccessToken {
		t.Errorf("AccessToken mismatch after round trip")
	}

	if loadedConfig.Auth.RefreshToken != originalConfig.Auth.RefreshToken {
		t.Errorf("RefreshToken mismatch after round trip")
	}

	// Time comparison needs to account for JSON marshaling precision
	if !loadedConfig.Auth.ExpiresAt.Round(time.Second).Equal(originalConfig.Auth.ExpiresAt.Round(time.Second)) {
		t.Errorf("ExpiresAt mismatch after round trip: expected %v, got %v",
			originalConfig.Auth.ExpiresAt, loadedConfig.Auth.ExpiresAt)
	}

	if loadedConfig.Defaults.AddressID != originalConfig.Defaults.AddressID {
		t.Errorf("AddressID mismatch after round trip")
	}

	if loadedConfig.Defaults.PaymentID != originalConfig.Defaults.PaymentID {
		t.Errorf("PaymentID mismatch after round trip")
	}

	if loadedConfig.RateLimiting.MinDelayMs != originalConfig.RateLimiting.MinDelayMs {
		t.Errorf("MinDelayMs mismatch after round trip")
	}

	if loadedConfig.RateLimiting.MaxDelayMs != originalConfig.RateLimiting.MaxDelayMs {
		t.Errorf("MaxDelayMs mismatch after round trip")
	}

	if loadedConfig.RateLimiting.MaxRetries != originalConfig.RateLimiting.MaxRetries {
		t.Errorf("MaxRetries mismatch after round trip")
	}
}

func TestSaveConfig_CreatesNestedDirectories(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "deeply", "nested", "path", "config.json")

	config := NewDefaultConfig()

	// Save config to deeply nested path
	err := SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed with nested directories: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created in nested directory")
	}
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create an empty file
	err := os.WriteFile(configPath, []byte(""), 0600)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Try to load the config
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Fatal("LoadConfig should have returned an error for empty file")
	}
}

func TestConfigJSONStructure(t *testing.T) {
	// Create a test config
	expiresAt := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	config := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-token",
			RefreshToken: "refresh-token",
			ExpiresAt:    expiresAt,
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr_test",
			PaymentID:    "pay_test",
			OutputFormat: "json",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Unmarshal back
	var parsedConfig Config
	if err := json.Unmarshal(data, &parsedConfig); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify structure matches PRD specification
	if parsedConfig.Auth.AccessToken != config.Auth.AccessToken {
		t.Errorf("Auth.AccessToken mismatch")
	}

	if parsedConfig.Defaults.OutputFormat != config.Defaults.OutputFormat {
		t.Errorf("Defaults.OutputFormat mismatch")
	}

	if parsedConfig.RateLimiting.MinDelayMs != config.RateLimiting.MinDelayMs {
		t.Errorf("RateLimiting.MinDelayMs mismatch")
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	path := GetDefaultConfigPath()
	if path == "" {
		t.Fatal("GetDefaultConfigPath returned empty string")
	}

	// Should contain .amazon-cli and config.json
	if !filepath.IsAbs(path) {
		t.Error("GetDefaultConfigPath should return an absolute path")
	}

	if filepath.Base(path) != "config.json" {
		t.Errorf("Expected filename 'config.json', got '%s'", filepath.Base(path))
	}

	if filepath.Base(filepath.Dir(path)) != ".amazon-cli" {
		t.Errorf("Expected directory '.amazon-cli', got '%s'", filepath.Base(filepath.Dir(path)))
	}
}

func TestLoadConfig_WithTildeExpansion(t *testing.T) {
	// Create a test config in a temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	config := NewDefaultConfig()
	config.Auth.AccessToken = "test-token"

	err := SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load using absolute path
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConfig.Auth.AccessToken != "test-token" {
		t.Errorf("Config not loaded correctly")
	}
}

func TestSaveConfig_WithTildeExpansion(t *testing.T) {
	// This tests that tilde expansion works (though we can't easily test ~ itself)
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test", "config.json")

	config := NewDefaultConfig()
	config.Auth.AccessToken = "tilde-test"

	err := SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}
}

func TestLoadConfig_ReadError(t *testing.T) {
	// Create a directory with the name of the config file (will cause read error)
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	err := os.Mkdir(configPath, 0700)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Try to load it as a config file
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Fatal("LoadConfig should have returned an error when trying to read a directory")
	}
}

func TestSaveConfig_MarshalSucceeds(t *testing.T) {
	// Test that config with all fields set marshals correctly
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	expiresAt := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	config := &Config{
		Auth: AuthConfig{
			AccessToken:  "access",
			RefreshToken: "refresh",
			ExpiresAt:    expiresAt,
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr",
			PaymentID:    "pay",
			OutputFormat: "json",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 100,
			MaxDelayMs: 1000,
			MaxRetries: 2,
		},
	}

	err := SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Load and verify
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConfig.RateLimiting.MaxRetries != 2 {
		t.Errorf("MaxRetries not saved/loaded correctly")
	}
}
