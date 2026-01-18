package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		config      *Config
		path        string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful save with valid config",
			config: &Config{
				Auth: AuthConfig{
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
					ExpiresAt:    "2024-01-20T12:00:00Z",
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
			},
			path:    filepath.Join(tempDir, "config.json"),
			wantErr: false,
		},
		{
			name: "successful save with nested directory creation",
			config: &Config{
				Auth: AuthConfig{
					AccessToken:  "token",
					RefreshToken: "refresh",
					ExpiresAt:    "2024-01-20T12:00:00Z",
				},
				Defaults: DefaultsConfig{
					OutputFormat: "json",
				},
				RateLimiting: RateLimitConfig{
					MinDelayMs: 1000,
					MaxDelayMs: 5000,
					MaxRetries: 3,
				},
			},
			path:    filepath.Join(tempDir, "nested", "dir", "config.json"),
			wantErr: false,
		},
		{
			name: "successful save with empty auth fields",
			config: &Config{
				Auth: AuthConfig{},
				Defaults: DefaultsConfig{
					OutputFormat: "json",
				},
				RateLimiting: RateLimitConfig{
					MinDelayMs: 1000,
					MaxDelayMs: 5000,
					MaxRetries: 3,
				},
			},
			path:    filepath.Join(tempDir, "empty-auth.json"),
			wantErr: false,
		},
		{
			name:        "nil config returns error",
			config:      nil,
			path:        filepath.Join(tempDir, "nil-config.json"),
			wantErr:     true,
			errContains: "config cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveConfig(tt.config, tt.path)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && err != nil {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("SaveConfig() error = %v, should contain %v", err, tt.errContains)
					}
				}
				return
			}

			// Verify file was created
			if _, err := os.Stat(tt.path); os.IsNotExist(err) {
				t.Errorf("Config file was not created at %s", tt.path)
				return
			}

			// Verify file permissions
			info, err := os.Stat(tt.path)
			if err != nil {
				t.Fatalf("Failed to stat config file: %v", err)
			}
			mode := info.Mode().Perm()
			expectedMode := os.FileMode(0600)
			if mode != expectedMode {
				t.Errorf("Config file permissions = %v, want %v", mode, expectedMode)
			}

			// Read back the file and verify content
			data, err := os.ReadFile(tt.path)
			if err != nil {
				t.Fatalf("Failed to read config file: %v", err)
			}

			// Unmarshal and compare
			var readConfig Config
			if err := json.Unmarshal(data, &readConfig); err != nil {
				t.Fatalf("Failed to unmarshal saved config: %v", err)
			}

			// Compare structs
			if !configsEqual(tt.config, &readConfig) {
				t.Errorf("Saved config does not match original.\nOriginal: %+v\nRead: %+v", tt.config, &readConfig)
			}

			// Verify JSON is indented (human-readable)
			// Check that the JSON contains newlines and indentation
			dataStr := string(data)
			if !contains(dataStr, "\n") {
				t.Error("Saved JSON does not contain newlines (not indented)")
			}
			if !contains(dataStr, "  ") {
				t.Error("Saved JSON does not contain indentation spaces")
			}
		})
	}
}

func TestSaveConfigDirectoryPermissions(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "subdir", "config.json")

	config := &Config{
		Auth: AuthConfig{
			AccessToken: "test",
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
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Check directory permissions
	dirInfo, err := os.Stat(filepath.Dir(configPath))
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}

	dirMode := dirInfo.Mode().Perm()
	expectedDirMode := os.FileMode(0700)
	if dirMode != expectedDirMode {
		t.Errorf("Directory permissions = %v, want %v", dirMode, expectedDirMode)
	}
}

func TestSaveConfigOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Save initial config
	config1 := &Config{
		Auth: AuthConfig{
			AccessToken: "token1",
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

	err := SaveConfig(config1, configPath)
	if err != nil {
		t.Fatalf("First SaveConfig() failed: %v", err)
	}

	// Overwrite with new config
	config2 := &Config{
		Auth: AuthConfig{
			AccessToken: "token2",
		},
		Defaults: DefaultsConfig{
			OutputFormat: "table",
		},
		RateLimiting: RateLimitConfig{
			MinDelayMs: 2000,
			MaxDelayMs: 10000,
			MaxRetries: 5,
		},
	}

	err = SaveConfig(config2, configPath)
	if err != nil {
		t.Fatalf("Second SaveConfig() failed: %v", err)
	}

	// Read back and verify it's the second config
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var readConfig Config
	if err := json.Unmarshal(data, &readConfig); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if readConfig.Auth.AccessToken != "token2" {
		t.Errorf("Expected overwritten config, got AccessToken = %v, want token2", readConfig.Auth.AccessToken)
	}
	if readConfig.Defaults.OutputFormat != "table" {
		t.Errorf("Expected overwritten config, got OutputFormat = %v, want table", readConfig.Defaults.OutputFormat)
	}
	if readConfig.RateLimiting.MinDelayMs != 2000 {
		t.Errorf("Expected overwritten config, got MinDelayMs = %v, want 2000", readConfig.RateLimiting.MinDelayMs)
	}
}

// Helper functions
func configsEqual(a, b *Config) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Auth.AccessToken == b.Auth.AccessToken &&
		a.Auth.RefreshToken == b.Auth.RefreshToken &&
		a.Auth.ExpiresAt == b.Auth.ExpiresAt &&
		a.Defaults.AddressID == b.Defaults.AddressID &&
		a.Defaults.PaymentID == b.Defaults.PaymentID &&
		a.Defaults.OutputFormat == b.Defaults.OutputFormat &&
		a.RateLimiting.MinDelayMs == b.RateLimiting.MinDelayMs &&
		a.RateLimiting.MaxDelayMs == b.RateLimiting.MaxDelayMs &&
		a.RateLimiting.MaxRetries == b.RateLimiting.MaxRetries
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
