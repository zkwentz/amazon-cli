package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_FileNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent", "config.json")

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected config to be returned")
	}

	// Should return default config
	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("expected default output format 'json', got %s", cfg.Defaults.OutputFormat)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create a config
	cfg := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
		Defaults: DefaultsConfig{
			AddressID:    "addr-123",
			PaymentID:    "pay-456",
			OutputFormat: "json",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	// Save config
	err := SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	// Load config
	loadedCfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Verify values
	if loadedCfg.Auth.AccessToken != cfg.Auth.AccessToken {
		t.Errorf("expected access token %s, got %s", cfg.Auth.AccessToken, loadedCfg.Auth.AccessToken)
	}

	if loadedCfg.Auth.RefreshToken != cfg.Auth.RefreshToken {
		t.Errorf("expected refresh token %s, got %s", cfg.Auth.RefreshToken, loadedCfg.Auth.RefreshToken)
	}

	if loadedCfg.Defaults.AddressID != cfg.Defaults.AddressID {
		t.Errorf("expected address ID %s, got %s", cfg.Defaults.AddressID, loadedCfg.Defaults.AddressID)
	}

	if loadedCfg.RateLimiting.MinDelayMs != cfg.RateLimiting.MinDelayMs {
		t.Errorf("expected min delay %d, got %d", cfg.RateLimiting.MinDelayMs, loadedCfg.RateLimiting.MinDelayMs)
	}
}

func TestSaveConfig_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nested", "dir", "config.json")

	cfg := GetDefaultConfig()
	err := SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify directory was created
	dir := filepath.Dir(configPath)
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Fatal("expected directory to be created")
	}

	// Verify permissions (0700)
	mode := info.Mode().Perm()
	if mode != 0700 {
		t.Errorf("expected directory permissions 0700, got %o", mode)
	}
}

func TestSaveConfig_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := GetDefaultConfig()
	err := SaveConfig(cfg, configPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file permissions (0600)
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("expected file permissions 0600, got %o", mode)
	}
}

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg == nil {
		t.Fatal("expected config to be returned")
	}

	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("expected default output format 'json', got %s", cfg.Defaults.OutputFormat)
	}

	if cfg.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("expected default min delay 1000, got %d", cfg.RateLimiting.MinDelayMs)
	}

	if cfg.RateLimiting.MaxDelayMs != 5000 {
		t.Errorf("expected default max delay 5000, got %d", cfg.RateLimiting.MaxDelayMs)
	}

	if cfg.RateLimiting.MaxRetries != 3 {
		t.Errorf("expected default max retries 3, got %d", cfg.RateLimiting.MaxRetries)
	}
}

func TestLoadConfig_EmptyPath(t *testing.T) {
	// This should use default path
	// Since we can't control the home directory in tests easily,
	// we'll just verify it doesn't panic
	_, err := LoadConfig("")
	// It's okay if this errors (e.g., file not found), we just want to ensure it doesn't panic
	_ = err
}
