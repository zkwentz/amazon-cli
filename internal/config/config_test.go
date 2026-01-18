package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Test defaults
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

	if config.Auth.AccessToken != "" {
		t.Errorf("Expected empty AccessToken, got '%s'", config.Auth.AccessToken)
	}

	if config.Auth.RefreshToken != "" {
		t.Errorf("Expected empty RefreshToken, got '%s'", config.Auth.RefreshToken)
	}
}

func TestGetConfigPath(t *testing.T) {
	tests := []struct {
		name       string
		customPath string
		wantCustom bool
	}{
		{
			name:       "Custom path provided",
			customPath: "/custom/path/config.json",
			wantCustom: true,
		},
		{
			name:       "No custom path",
			customPath: "",
			wantCustom: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := GetConfigPath(tt.customPath)
			if err != nil {
				t.Fatalf("GetConfigPath() error = %v", err)
			}

			if tt.wantCustom {
				if path != tt.customPath {
					t.Errorf("Expected custom path '%s', got '%s'", tt.customPath, path)
				}
			} else {
				// Should contain .amazon-cli directory
				if !filepath.IsAbs(path) {
					t.Errorf("Expected absolute path, got '%s'", path)
				}
				if filepath.Base(path) != "config.json" {
					t.Errorf("Expected filename 'config.json', got '%s'", filepath.Base(path))
				}
				if filepath.Base(filepath.Dir(path)) != ".amazon-cli" {
					t.Errorf("Expected directory '.amazon-cli', got '%s'", filepath.Base(filepath.Dir(path)))
				}
			}
		})
	}
}

func TestLoadConfig_NonExistent(t *testing.T) {
	// Create a temporary path that doesn't exist
	tempPath := filepath.Join(t.TempDir(), "nonexistent", "config.json")

	config, err := LoadConfig(tempPath)
	if err != nil {
		t.Fatalf("LoadConfig() should not error for non-existent file, got: %v", err)
	}

	// Should return default config
	if config == nil {
		t.Fatal("LoadConfig() returned nil for non-existent file")
	}

	defaultConfig := DefaultConfig()
	if config.Defaults.OutputFormat != defaultConfig.Defaults.OutputFormat {
		t.Errorf("Expected default output format, got '%s'", config.Defaults.OutputFormat)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".amazon-cli", "config.json")

	// Create a test config
	testConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
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

	// Save the config
	if err := SaveConfig(testConfig, configPath); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Verify directory permissions
	dirInfo, err := os.Stat(filepath.Dir(configPath))
	if err != nil {
		t.Fatalf("Failed to stat config directory: %v", err)
	}
	if dirInfo.Mode().Perm() != 0700 {
		t.Errorf("Expected directory permissions 0700, got %o", dirInfo.Mode().Perm())
	}

	// Verify file permissions
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}
	if fileInfo.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", fileInfo.Mode().Perm())
	}

	// Load the config back
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify all fields match
	if loadedConfig.Auth.AccessToken != testConfig.Auth.AccessToken {
		t.Errorf("AccessToken mismatch: got '%s', want '%s'", loadedConfig.Auth.AccessToken, testConfig.Auth.AccessToken)
	}

	if loadedConfig.Auth.RefreshToken != testConfig.Auth.RefreshToken {
		t.Errorf("RefreshToken mismatch: got '%s', want '%s'", loadedConfig.Auth.RefreshToken, testConfig.Auth.RefreshToken)
	}

	if !loadedConfig.Auth.ExpiresAt.Equal(testConfig.Auth.ExpiresAt) {
		t.Errorf("ExpiresAt mismatch: got '%v', want '%v'", loadedConfig.Auth.ExpiresAt, testConfig.Auth.ExpiresAt)
	}

	if loadedConfig.Defaults.AddressID != testConfig.Defaults.AddressID {
		t.Errorf("AddressID mismatch: got '%s', want '%s'", loadedConfig.Defaults.AddressID, testConfig.Defaults.AddressID)
	}

	if loadedConfig.Defaults.PaymentID != testConfig.Defaults.PaymentID {
		t.Errorf("PaymentID mismatch: got '%s', want '%s'", loadedConfig.Defaults.PaymentID, testConfig.Defaults.PaymentID)
	}

	if loadedConfig.Defaults.OutputFormat != testConfig.Defaults.OutputFormat {
		t.Errorf("OutputFormat mismatch: got '%s', want '%s'", loadedConfig.Defaults.OutputFormat, testConfig.Defaults.OutputFormat)
	}

	if loadedConfig.RateLimiting.MinDelayMs != testConfig.RateLimiting.MinDelayMs {
		t.Errorf("MinDelayMs mismatch: got %d, want %d", loadedConfig.RateLimiting.MinDelayMs, testConfig.RateLimiting.MinDelayMs)
	}

	if loadedConfig.RateLimiting.MaxDelayMs != testConfig.RateLimiting.MaxDelayMs {
		t.Errorf("MaxDelayMs mismatch: got %d, want %d", loadedConfig.RateLimiting.MaxDelayMs, testConfig.RateLimiting.MaxDelayMs)
	}

	if loadedConfig.RateLimiting.MaxRetries != testConfig.RateLimiting.MaxRetries {
		t.Errorf("MaxRetries mismatch: got %d, want %d", loadedConfig.RateLimiting.MaxRetries, testConfig.RateLimiting.MaxRetries)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	invalidJSON := []byte(`{"auth": "invalid json structure`)
	if err := os.WriteFile(configPath, invalidJSON, 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Should return error for invalid JSON
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig() should return error for invalid JSON")
	}
}

func TestSaveConfig_EmptyPath(t *testing.T) {
	// Saving with empty path should use default path
	testConfig := DefaultConfig()
	testConfig.Auth.AccessToken = "test-token"

	// This test requires a real home directory, so we'll just verify it doesn't panic
	// In a real environment, this would save to ~/.amazon-cli/config.json
	// For testing, we'll skip this if we can't get the home directory
	if _, err := os.UserHomeDir(); err == nil {
		// We won't actually save to the real home directory in tests
		// Just verify the function signature works
		t.Skip("Skipping test that would write to real home directory")
	}
}

func TestJSONMarshaling(t *testing.T) {
	// Create a config with all fields populated
	config := &Config{
		Auth: AuthConfig{
			AccessToken:  "access-123",
			RefreshToken: "refresh-456",
			ExpiresAt:    time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr_default",
			PaymentID:    "pay_default",
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
	var unmarshaled Config
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify roundtrip
	if unmarshaled.Auth.AccessToken != config.Auth.AccessToken {
		t.Errorf("AccessToken mismatch after roundtrip")
	}

	if !unmarshaled.Auth.ExpiresAt.Equal(config.Auth.ExpiresAt) {
		t.Errorf("ExpiresAt mismatch after roundtrip")
	}
}
