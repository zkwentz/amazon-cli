package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfigNewFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Load config from non-existent file
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Should return default config
	if cfg.Defaults.OutputFormat != "json" {
		t.Errorf("Expected default output format 'json', got '%s'", cfg.Defaults.OutputFormat)
	}

	if cfg.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("Expected default MinDelayMs 1000, got %d", cfg.RateLimiting.MinDelayMs)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create config
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
		RateLimiting: RateLimitingConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	// Save config
	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load config
	loadedCfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify loaded config matches saved config
	if loadedCfg.Auth.AccessToken != cfg.Auth.AccessToken {
		t.Errorf("AccessToken mismatch: expected %s, got %s", cfg.Auth.AccessToken, loadedCfg.Auth.AccessToken)
	}

	if loadedCfg.Auth.RefreshToken != cfg.Auth.RefreshToken {
		t.Errorf("RefreshToken mismatch: expected %s, got %s", cfg.Auth.RefreshToken, loadedCfg.Auth.RefreshToken)
	}

	if loadedCfg.Defaults.AddressID != cfg.Defaults.AddressID {
		t.Errorf("AddressID mismatch: expected %s, got %s", cfg.Defaults.AddressID, loadedCfg.Defaults.AddressID)
	}
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name: "authenticated with access token",
			config: Config{
				Auth: AuthConfig{
					AccessToken: "test-token",
				},
			},
			expected: true,
		},
		{
			name: "authenticated with refresh token",
			config: Config{
				Auth: AuthConfig{
					RefreshToken: "test-refresh",
				},
			},
			expected: true,
		},
		{
			name: "not authenticated",
			config: Config{
				Auth: AuthConfig{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsAuthenticated()
			if result != tt.expected {
				t.Errorf("IsAuthenticated() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsTokenExpired(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name: "token not expired",
			config: Config{
				Auth: AuthConfig{
					ExpiresAt: time.Now().Add(1 * time.Hour),
				},
			},
			expected: false,
		},
		{
			name: "token expired",
			config: Config{
				Auth: AuthConfig{
					ExpiresAt: time.Now().Add(-1 * time.Hour),
				},
			},
			expected: true,
		},
		{
			name: "zero time means expired",
			config: Config{
				Auth: AuthConfig{
					ExpiresAt: time.Time{},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsTokenExpired()
			if result != tt.expected {
				t.Errorf("IsTokenExpired() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestConfigFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &Config{
		Auth: AuthConfig{
			AccessToken: "secret-token",
		},
	}

	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// Should be 0600 (read/write for owner only)
	expectedPerm := os.FileMode(0600)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("Config file has incorrect permissions: got %v, expected %v", info.Mode().Perm(), expectedPerm)
	}
}
