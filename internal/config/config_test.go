package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGetConfigPath(t *testing.T) {
	// Save original viper state and reset after test
	defer viper.Reset()

	t.Run("returns custom path when config flag is set", func(t *testing.T) {
		viper.Reset()
		customPath := "/custom/path/to/config.json"
		viper.Set("config", customPath)

		result := GetConfigPath()
		if result != customPath {
			t.Errorf("expected %s, got %s", customPath, result)
		}
	})

	t.Run("returns default path when config flag is not set", func(t *testing.T) {
		viper.Reset()

		result := GetConfigPath()

		// Should contain .amazon-cli/config.json
		if !strings.Contains(result, ".amazon-cli") || !strings.HasSuffix(result, "config.json") {
			t.Errorf("expected path to contain .amazon-cli/config.json, got %s", result)
		}

		// Should be an absolute path (unless home dir is not available)
		if !filepath.IsAbs(result) && !strings.HasPrefix(result, ".amazon-cli") {
			t.Errorf("expected absolute path or relative .amazon-cli path, got %s", result)
		}
	})

	t.Run("handles relative custom path", func(t *testing.T) {
		viper.Reset()
		customPath := "./my-config.json"
		viper.Set("config", customPath)

		result := GetConfigPath()
		if result != customPath {
			t.Errorf("expected %s, got %s", customPath, result)
		}
	})

	t.Run("handles absolute custom path", func(t *testing.T) {
		viper.Reset()
		customPath := "/etc/amazon-cli/config.json"
		viper.Set("config", customPath)

		result := GetConfigPath()
		if result != customPath {
			t.Errorf("expected %s, got %s", customPath, result)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("returns default config when file does not exist", func(t *testing.T) {
		nonExistentPath := "/tmp/this-should-not-exist-12345.json"

		config, err := LoadConfig(nonExistentPath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if config == nil {
			t.Fatal("expected config to be non-nil")
		}

		// Check default values
		if config.Defaults.OutputFormat != "json" {
			t.Errorf("expected default output format 'json', got %s", config.Defaults.OutputFormat)
		}
		if config.RateLimiting.MinDelayMs != 1000 {
			t.Errorf("expected default MinDelayMs 1000, got %d", config.RateLimiting.MinDelayMs)
		}
		if config.RateLimiting.MaxDelayMs != 5000 {
			t.Errorf("expected default MaxDelayMs 5000, got %d", config.RateLimiting.MaxDelayMs)
		}
		if config.RateLimiting.MaxRetries != 3 {
			t.Errorf("expected default MaxRetries 3, got %d", config.RateLimiting.MaxRetries)
		}
	})

	t.Run("loads valid config file", func(t *testing.T) {
		// Create temporary config file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.json")

		testConfig := &Config{
			Auth: AuthConfig{
				AccessToken:  "test_access_token",
				RefreshToken: "test_refresh_token",
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

		// Write test config
		data, _ := json.MarshalIndent(testConfig, "", "  ")
		if err := os.WriteFile(configPath, data, 0600); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		// Load config
		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify loaded values
		if config.Auth.AccessToken != "test_access_token" {
			t.Errorf("expected access token 'test_access_token', got %s", config.Auth.AccessToken)
		}
		if config.Defaults.AddressID != "addr_123" {
			t.Errorf("expected address ID 'addr_123', got %s", config.Defaults.AddressID)
		}
		if config.RateLimiting.MinDelayMs != 2000 {
			t.Errorf("expected MinDelayMs 2000, got %d", config.RateLimiting.MinDelayMs)
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "invalid.json")

		// Write invalid JSON
		if err := os.WriteFile(configPath, []byte("invalid json {{{"), 0600); err != nil {
			t.Fatalf("failed to write invalid config: %v", err)
		}

		_, err := LoadConfig(configPath)
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})
}

func TestSaveConfig(t *testing.T) {
	t.Run("saves config successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "test-config.json")

		testConfig := &Config{
			Auth: AuthConfig{
				AccessToken:  "save_test_token",
				RefreshToken: "save_test_refresh",
				ExpiresAt:    time.Now().Add(1 * time.Hour),
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
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("config file was not created")
		}

		// Verify file permissions
		info, _ := os.Stat(configPath)
		mode := info.Mode()
		if mode != 0600 {
			t.Errorf("expected file mode 0600, got %o", mode)
		}

		// Load and verify content
		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("failed to load saved config: %v", err)
		}

		if config.Auth.AccessToken != "save_test_token" {
			t.Errorf("expected access token 'save_test_token', got %s", config.Auth.AccessToken)
		}
	})

	t.Run("creates directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		nestedPath := filepath.Join(tmpDir, "nested", "dir", "config.json")

		testConfig := &Config{
			Defaults: DefaultsConfig{
				OutputFormat: "json",
			},
			RateLimiting: RateLimitConfig{
				MinDelayMs: 1000,
				MaxDelayMs: 5000,
				MaxRetries: 3,
			},
		}

		// Save config to nested path
		err := SaveConfig(testConfig, nestedPath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify directory was created
		dir := filepath.Dir(nestedPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Error("directory was not created")
		}

		// Verify directory permissions
		info, _ := os.Stat(dir)
		mode := info.Mode()
		if mode&0700 != 0700 {
			t.Errorf("expected directory mode to have at least 0700, got %o", mode)
		}

		// Verify file exists
		if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
			t.Error("config file was not created")
		}
	})
}
