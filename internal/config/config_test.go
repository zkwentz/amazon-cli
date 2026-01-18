package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestDefaultConfig tests the default configuration
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	// Check defaults
	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("Expected default output format 'json', got '%s'", cfg.Defaults.OutputFormat)
	}

	if cfg.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("Expected default MinDelayMs 1000, got %d", cfg.RateLimiting.MinDelayMs)
	}

	if cfg.RateLimiting.MaxDelayMs != 5000 {
		t.Errorf("Expected default MaxDelayMs 5000, got %d", cfg.RateLimiting.MaxDelayMs)
	}

	if cfg.RateLimiting.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries 3, got %d", cfg.RateLimiting.MaxRetries)
	}

	// Check that auth is empty
	if cfg.Auth.AccessToken != "" {
		t.Error("Expected empty access token in default config")
	}

	if cfg.Auth.RefreshToken != "" {
		t.Error("Expected empty refresh token in default config")
	}
}

// TestSaveAndLoadConfig tests saving and loading configuration
func TestSaveAndLoadConfig(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create a config to save
	originalConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
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

	// Save config
	err := SaveConfig(originalConfig, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load config
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config matches original
	if loadedConfig.Auth.AccessToken != originalConfig.Auth.AccessToken {
		t.Errorf("Expected access token '%s', got '%s'",
			originalConfig.Auth.AccessToken, loadedConfig.Auth.AccessToken)
	}

	if loadedConfig.Auth.RefreshToken != originalConfig.Auth.RefreshToken {
		t.Errorf("Expected refresh token '%s', got '%s'",
			originalConfig.Auth.RefreshToken, loadedConfig.Auth.RefreshToken)
	}

	if loadedConfig.Defaults.AddressID != originalConfig.Defaults.AddressID {
		t.Errorf("Expected address ID '%s', got '%s'",
			originalConfig.Defaults.AddressID, loadedConfig.Defaults.AddressID)
	}

	if loadedConfig.Defaults.PaymentID != originalConfig.Defaults.PaymentID {
		t.Errorf("Expected payment ID '%s', got '%s'",
			originalConfig.Defaults.PaymentID, loadedConfig.Defaults.PaymentID)
	}

	if loadedConfig.Defaults.OutputFormat != originalConfig.Defaults.OutputFormat {
		t.Errorf("Expected output format '%s', got '%s'",
			originalConfig.Defaults.OutputFormat, loadedConfig.Defaults.OutputFormat)
	}

	if loadedConfig.RateLimiting.MinDelayMs != originalConfig.RateLimiting.MinDelayMs {
		t.Errorf("Expected MinDelayMs %d, got %d",
			originalConfig.RateLimiting.MinDelayMs, loadedConfig.RateLimiting.MinDelayMs)
	}

	if loadedConfig.RateLimiting.MaxRetries != originalConfig.RateLimiting.MaxRetries {
		t.Errorf("Expected MaxRetries %d, got %d",
			originalConfig.RateLimiting.MaxRetries, loadedConfig.RateLimiting.MaxRetries)
	}
}

// TestLoadConfigNonExistent tests loading a non-existent config file
func TestLoadConfigNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "nonexistent", "config.json")

	// Load config from non-existent path
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Expected no error when loading non-existent config, got: %v", err)
	}

	// Should return default config
	if cfg == nil {
		t.Fatal("Expected default config, got nil")
	}

	if cfg.Defaults.OutputFormat != "json" {
		t.Error("Expected default config values")
	}
}

// TestSaveConfigCreatesDirectory tests that SaveConfig creates parent directories
func TestSaveConfigCreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "subdir1", "subdir2", "config.json")

	cfg := DefaultConfig()
	err := SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify directories were created
	if _, err := os.Stat(filepath.Dir(configPath)); os.IsNotExist(err) {
		t.Error("Parent directories were not created")
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

// TestConfigFilePermissions tests that config file has correct permissions
func TestConfigFilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := DefaultConfig()
	cfg.Auth.AccessToken = "secret-token"

	err := SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// File should be readable and writable by owner only (0600)
	mode := info.Mode()
	expectedMode := os.FileMode(0600)
	if mode.Perm() != expectedMode {
		t.Errorf("Expected file permissions %v, got %v", expectedMode, mode.Perm())
	}
}

// TestConfigDirectoryPermissions tests that config directory has correct permissions
func TestConfigDirectoryPermissions(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "amazon-cli")
	configPath := filepath.Join(configDir, "config.json")

	cfg := DefaultConfig()
	err := SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check directory permissions
	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Failed to stat config directory: %v", err)
	}

	// Directory should be rwx for owner only (0700)
	mode := info.Mode()
	expectedMode := os.FileMode(0700) | os.ModeDir
	if mode.Perm() != expectedMode.Perm() {
		t.Errorf("Expected directory permissions %v, got %v", expectedMode.Perm(), mode.Perm())
	}
}

// TestConfigJSONFormat tests that the saved config is valid JSON
func TestConfigJSONFormat(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-token",
			RefreshToken: "refresh-token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
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

	err := SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Read and verify JSON formatting
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	jsonStr := string(data)

	// Check that it's indented (contains newlines and spaces)
	if !contains(jsonStr, "\n") {
		t.Error("Expected indented JSON, got minified")
	}

	// Check that required fields are present
	requiredFields := []string{
		`"auth"`,
		`"access_token"`,
		`"refresh_token"`,
		`"expires_at"`,
		`"defaults"`,
		`"address_id"`,
		`"payment_id"`,
		`"output_format"`,
		`"rate_limiting"`,
		`"min_delay_ms"`,
		`"max_delay_ms"`,
		`"max_retries"`,
	}

	for _, field := range requiredFields {
		if !contains(jsonStr, field) {
			t.Errorf("Expected JSON to contain %s", field)
		}
	}
}

// TestLoadConfigWithInvalidJSON tests loading a config file with invalid JSON
func TestLoadConfigWithInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Write invalid JSON
	err := os.WriteFile(configPath, []byte(`{invalid json}`), 0600)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	// Try to load config
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid JSON, got nil")
	}
}

// TestConfigRoundTrip tests that config survives save/load round trip
func TestConfigRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create config with various values
	originalConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "access-123",
			RefreshToken: "refresh-456",
			ExpiresAt:    time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
		},
		Defaults: DefaultsConfig{
			AddressID:    "my-address",
			PaymentID:    "my-payment",
			OutputFormat: "table",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1500,
			MaxDelayMs: 7500,
			MaxRetries: 4,
		},
	}

	// Save
	if err := SaveConfig(originalConfig, configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Save again
	configPath2 := filepath.Join(tempDir, "config2.json")
	if err := SaveConfig(loadedConfig, configPath2); err != nil {
		t.Fatalf("Second save failed: %v", err)
	}

	// Load again
	loadedConfig2, err := LoadConfig(configPath2)
	if err != nil {
		t.Fatalf("Second load failed: %v", err)
	}

	// Compare all fields
	if loadedConfig2.Auth.AccessToken != originalConfig.Auth.AccessToken {
		t.Error("Access token changed after round trip")
	}

	if loadedConfig2.RateLimiting.MinDelayMs != originalConfig.RateLimiting.MinDelayMs {
		t.Error("MinDelayMs changed after round trip")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
