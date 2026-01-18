package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test configs
	tmpDir := t.TempDir()

	t.Run("loads existing config", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "test-config.json")

		// Create a test config
		testConfig := &Config{
			Auth: AuthConfig{
				AccessToken:  "test_token",
				RefreshToken: "refresh_token",
				ExpiresAt:    time.Now().Add(1 * time.Hour),
			},
			Defaults: DefaultsConfig{
				AddressID:    "addr_123",
				PaymentID:    "pay_123",
				OutputFormat: "json",
			},
			RateLimiting: RateLimitConfig{
				MinDelayMs: 1000,
				MaxDelayMs: 5000,
				MaxRetries: 3,
			},
		}

		// Save it
		err := SaveConfig(testConfig, configPath)
		if err != nil {
			t.Fatalf("SaveConfig failed: %v", err)
		}

		// Load it back
		loadedConfig, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		// Verify loaded config matches
		if loadedConfig.Auth.AccessToken != testConfig.Auth.AccessToken {
			t.Errorf("AccessToken = %v, want %v", loadedConfig.Auth.AccessToken, testConfig.Auth.AccessToken)
		}
		if loadedConfig.Defaults.AddressID != testConfig.Defaults.AddressID {
			t.Errorf("AddressID = %v, want %v", loadedConfig.Defaults.AddressID, testConfig.Defaults.AddressID)
		}
		if loadedConfig.RateLimiting.MinDelayMs != testConfig.RateLimiting.MinDelayMs {
			t.Errorf("MinDelayMs = %v, want %v", loadedConfig.RateLimiting.MinDelayMs, testConfig.RateLimiting.MinDelayMs)
		}
	})

	t.Run("returns default config for non-existent file", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "nonexistent", "config.json")

		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		// Should return default config
		if config.Defaults.OutputFormat != "json" {
			t.Errorf("OutputFormat = %v, want json", config.Defaults.OutputFormat)
		}
		if config.RateLimiting.MinDelayMs != 1000 {
			t.Errorf("MinDelayMs = %v, want 1000", config.RateLimiting.MinDelayMs)
		}
	})

	t.Run("creates directory if not exists", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "subdir", "nested", "config.json")

		config := DefaultConfig()
		err := SaveConfig(config, configPath)
		if err != nil {
			t.Fatalf("SaveConfig failed: %v", err)
		}

		// Verify directory was created
		dir := filepath.Dir(configPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Error("Directory was not created")
		}
	})
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("saves config with correct permissions", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "config.json")

		config := &Config{
			Auth: AuthConfig{
				AccessToken:  "secret_token",
				RefreshToken: "secret_refresh",
				ExpiresAt:    time.Now(),
			},
			Defaults: DefaultsConfig{
				OutputFormat: "json",
			},
			RateLimiting: RateLimitConfig{
				MinDelayMs: 1000,
				MaxDelayMs: 5000,
				MaxRetries: 3,
			},
		}

		err := SaveConfig(config, configPath)
		if err != nil {
			t.Fatalf("SaveConfig failed: %v", err)
		}

		// Verify file exists
		info, err := os.Stat(configPath)
		if err != nil {
			t.Fatalf("Config file not created: %v", err)
		}

		// Verify permissions (should be 0600)
		if info.Mode().Perm() != 0600 {
			t.Errorf("Config file permissions = %o, want 0600", info.Mode().Perm())
		}
	})

	t.Run("creates valid JSON", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "json-test.json")

		config := DefaultConfig()
		err := SaveConfig(config, configPath)
		if err != nil {
			t.Fatalf("SaveConfig failed: %v", err)
		}

		// Load it back to verify it's valid JSON
		loadedConfig, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig failed, JSON might be invalid: %v", err)
		}

		if loadedConfig == nil {
			t.Error("Loaded config is nil")
		}
	})
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	// Verify default values
	if config.Defaults.OutputFormat != "json" {
		t.Errorf("OutputFormat = %v, want json", config.Defaults.OutputFormat)
	}

	if config.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("MinDelayMs = %v, want 1000", config.RateLimiting.MinDelayMs)
	}

	if config.RateLimiting.MaxDelayMs != 5000 {
		t.Errorf("MaxDelayMs = %v, want 5000", config.RateLimiting.MaxDelayMs)
	}

	if config.RateLimiting.MaxRetries != 3 {
		t.Errorf("MaxRetries = %v, want 3", config.RateLimiting.MaxRetries)
	}

	// Auth tokens should be empty
	if config.Auth.AccessToken != "" {
		t.Error("AccessToken should be empty in default config")
	}

	if config.Auth.RefreshToken != "" {
		t.Error("RefreshToken should be empty in default config")
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()

	if path == "" {
		t.Error("GetConfigPath returned empty string")
	}

	// Should contain .amazon-cli
	if !contains(path, ".amazon-cli") {
		t.Errorf("Path %v does not contain .amazon-cli", path)
	}

	// Should end with config.json
	if !contains(path, "config.json") {
		t.Errorf("Path %v does not contain config.json", path)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
