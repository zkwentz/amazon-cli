package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg.Auth.AuthMethod != "cookie" {
		t.Errorf("Expected default auth method to be 'cookie', got '%s'", cfg.Auth.AuthMethod)
	}

	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("Expected default output format to be 'json', got '%s'", cfg.Defaults.OutputFormat)
	}

	if cfg.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("Expected default min delay to be 1000ms, got %d", cfg.RateLimiting.MinDelayMs)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create test config
	testConfig := &Config{
		Auth: AuthConfig{
			AuthMethod: "cookie",
			Cookies: []Cookie{
				{
					Name:     "session",
					Value:    "test123",
					Domain:   ".amazon.com",
					Path:     "/",
					Expires:  time.Now().Add(24 * time.Hour),
					Secure:   true,
					HttpOnly: true,
				},
			},
			CookiesSetAt: time.Now(),
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

	// Save config
	err := SaveConfig(testConfig, configPath)
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

	// Verify loaded config matches
	if loadedConfig.Auth.AuthMethod != testConfig.Auth.AuthMethod {
		t.Errorf("Auth method mismatch: expected %s, got %s",
			testConfig.Auth.AuthMethod, loadedConfig.Auth.AuthMethod)
	}

	if len(loadedConfig.Auth.Cookies) != len(testConfig.Auth.Cookies) {
		t.Errorf("Cookie count mismatch: expected %d, got %d",
			len(testConfig.Auth.Cookies), len(loadedConfig.Auth.Cookies))
	}

	if len(loadedConfig.Auth.Cookies) > 0 {
		cookie := loadedConfig.Auth.Cookies[0]
		if cookie.Name != "session" || cookie.Value != "test123" {
			t.Errorf("Cookie data mismatch: got name=%s, value=%s", cookie.Name, cookie.Value)
		}
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent.json")

	// Should return default config without error
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got: %v", err)
	}

	if cfg.Auth.AuthMethod != "cookie" {
		t.Errorf("Expected default auth method, got: %s", cfg.Auth.AuthMethod)
	}
}

func TestConfigFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	testConfig := GetDefaultConfig()

	err := SaveConfig(testConfig, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	mode := info.Mode()
	expected := os.FileMode(0600)

	if mode != expected {
		t.Errorf("Expected file permissions %v, got %v", expected, mode)
	}
}
