package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Should return default config
	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("Default output format incorrect: got %v, want json", cfg.Defaults.OutputFormat)
	}
	if cfg.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("Default min delay incorrect: got %v, want 1000", cfg.RateLimiting.MinDelayMs)
	}
	if cfg.RateLimiting.MaxDelayMs != 5000 {
		t.Errorf("Default max delay incorrect: got %v, want 5000", cfg.RateLimiting.MaxDelayMs)
	}
	if cfg.RateLimiting.MaxRetries != 3 {
		t.Errorf("Default max retries incorrect: got %v, want 3", cfg.RateLimiting.MaxRetries)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create test config
	expiresAt := time.Now().Add(time.Hour)
	testCfg := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-access",
			RefreshToken: "test-refresh",
			ExpiresAt:    expiresAt,
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr123",
			PaymentID:    "pay456",
			OutputFormat: "table",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 2000,
			MaxDelayMs: 10000,
			MaxRetries: 5,
		},
	}

	// Save config
	if err := SaveConfig(testCfg, configPath); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Config file not created: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Config file has wrong permissions: got %v, want 0600", info.Mode().Perm())
	}

	// Load config
	loadedCfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify auth config
	if loadedCfg.Auth.AccessToken != "test-access" {
		t.Errorf("AccessToken incorrect: got %v, want test-access", loadedCfg.Auth.AccessToken)
	}
	if loadedCfg.Auth.RefreshToken != "test-refresh" {
		t.Errorf("RefreshToken incorrect: got %v, want test-refresh", loadedCfg.Auth.RefreshToken)
	}
	// Time comparison with small tolerance
	if loadedCfg.Auth.ExpiresAt.Unix() != expiresAt.Unix() {
		t.Errorf("ExpiresAt incorrect: got %v, want %v", loadedCfg.Auth.ExpiresAt, expiresAt)
	}

	// Verify defaults
	if loadedCfg.Defaults.AddressID != "addr123" {
		t.Errorf("AddressID incorrect: got %v, want addr123", loadedCfg.Defaults.AddressID)
	}
	if loadedCfg.Defaults.PaymentID != "pay456" {
		t.Errorf("PaymentID incorrect: got %v, want pay456", loadedCfg.Defaults.PaymentID)
	}
	if loadedCfg.Defaults.OutputFormat != "table" {
		t.Errorf("OutputFormat incorrect: got %v, want table", loadedCfg.Defaults.OutputFormat)
	}

	// Verify rate limiting
	if loadedCfg.RateLimiting.MinDelayMs != 2000 {
		t.Errorf("MinDelayMs incorrect: got %v, want 2000", loadedCfg.RateLimiting.MinDelayMs)
	}
	if loadedCfg.RateLimiting.MaxDelayMs != 10000 {
		t.Errorf("MaxDelayMs incorrect: got %v, want 10000", loadedCfg.RateLimiting.MaxDelayMs)
	}
	if loadedCfg.RateLimiting.MaxRetries != 5 {
		t.Errorf("MaxRetries incorrect: got %v, want 5", loadedCfg.RateLimiting.MaxRetries)
	}
}

func TestSaveConfig_CreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "subdir", "config.json")

	cfg := &Config{
		Defaults: DefaultsConfig{
			OutputFormat: "json",
		},
	}

	// Save should create the directory
	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify directory was created
	dirInfo, err := os.Stat(filepath.Dir(configPath))
	if err != nil {
		t.Fatalf("Directory not created: %v", err)
	}
	if !dirInfo.IsDir() {
		t.Errorf("Path is not a directory")
	}
	if dirInfo.Mode().Perm() != 0700 {
		t.Errorf("Directory has wrong permissions: got %v, want 0700", dirInfo.Mode().Perm())
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file not created")
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	// Should contain .amazon-cli
	if !contains(path, ".amazon-cli") {
		t.Errorf("Config path doesn't contain .amazon-cli: %v", path)
	}

	// Should end with config.json
	if filepath.Base(path) != "config.json" {
		t.Errorf("Config path doesn't end with config.json: %v", path)
	}
}

func contains(s, substr string) bool {
	return filepath.Base(filepath.Dir(s)) == substr
}
