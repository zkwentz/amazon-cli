package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	// Check defaults are set correctly
	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("expected output format 'json', got '%s'", cfg.Defaults.OutputFormat)
	}

	if cfg.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("expected MinDelayMs 1000, got %d", cfg.RateLimiting.MinDelayMs)
	}

	if cfg.RateLimiting.MaxDelayMs != 5000 {
		t.Errorf("expected MaxDelayMs 5000, got %d", cfg.RateLimiting.MaxDelayMs)
	}

	if cfg.RateLimiting.MaxRetries != 3 {
		t.Errorf("expected MaxRetries 3, got %d", cfg.RateLimiting.MaxRetries)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create a config with custom values
	cfg := &Config{
		Auth: AuthConfig{
			AccessToken:  "test-token",
			RefreshToken: "test-refresh",
			ExpiresAt:    "2024-12-31T23:59:59Z",
		},
		Defaults: DefaultsConfig{
			OutputFormat: "table",
			AddressID:    "addr123",
			PaymentID:    "pay456",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 2000,
			MaxDelayMs: 10000,
			MaxRetries: 5,
		},
	}

	// Save config
	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	// Load config
	loaded, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify values match
	if loaded.Auth.AccessToken != cfg.Auth.AccessToken {
		t.Errorf("AccessToken mismatch: got %s, want %s", loaded.Auth.AccessToken, cfg.Auth.AccessToken)
	}

	if loaded.Defaults.OutputFormat != cfg.Defaults.OutputFormat {
		t.Errorf("OutputFormat mismatch: got %s, want %s", loaded.Defaults.OutputFormat, cfg.Defaults.OutputFormat)
	}

	if loaded.RateLimiting.MaxRetries != cfg.RateLimiting.MaxRetries {
		t.Errorf("MaxRetries mismatch: got %d, want %d", loaded.RateLimiting.MaxRetries, cfg.RateLimiting.MaxRetries)
	}
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent", "config.json")

	// Loading non-existent config should return default config
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected default config, got nil")
	}

	// Should have default values
	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("expected default output format, got %s", cfg.Defaults.OutputFormat)
	}
}
