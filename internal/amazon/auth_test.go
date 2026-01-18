package amazon

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestRefreshTokenIfNeeded_NoRefreshToken(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "", // No refresh token
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
	}

	err := RefreshTokenIfNeeded(cfg)
	if err == nil {
		t.Error("expected error when refresh token is missing")
	}
	if err.Error() != "no refresh token available" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRefreshTokenIfNeeded_TokenStillValid(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Now().Add(10 * time.Minute), // Valid for 10 more minutes
		},
		RateLimiting: config.RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	// Store original values
	originalAccessToken := cfg.Auth.AccessToken
	originalExpiresAt := cfg.Auth.ExpiresAt

	err := RefreshTokenIfNeeded(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify token wasn't changed (no refresh needed)
	if cfg.Auth.AccessToken != originalAccessToken {
		t.Error("access token should not have changed")
	}
	if !cfg.Auth.ExpiresAt.Equal(originalExpiresAt) {
		t.Error("expiry time should not have changed")
	}
}

func TestRefreshTokenIfNeeded_TokenNeedsRefresh(t *testing.T) {
	// Note: This is a placeholder test for token refresh functionality
	// In a production implementation, we would refactor RefreshTokenIfNeeded to accept
	// a custom HTTP client or token URL to enable proper testing with a mock server

	t.Skip("Skipping integration test - requires refactoring to support dependency injection")
}

func TestRefreshTokenIfNeeded_ServerError(t *testing.T) {
	// Note: This is a placeholder test for error handling
	// In a production implementation, we would refactor RefreshTokenIfNeeded to accept
	// a custom HTTP client to enable proper testing with a mock server

	t.Skip("Skipping integration test - requires refactoring to support dependency injection")
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected bool
	}{
		{
			name: "valid authentication",
			config: &config.Config{
				Auth: config.AuthConfig{
					AccessToken:  "test-access-token",
					RefreshToken: "test-refresh-token",
				},
			},
			expected: true,
		},
		{
			name: "missing access token",
			config: &config.Config{
				Auth: config.AuthConfig{
					AccessToken:  "",
					RefreshToken: "test-refresh-token",
				},
			},
			expected: false,
		},
		{
			name: "missing refresh token",
			config: &config.Config{
				Auth: config.AuthConfig{
					AccessToken:  "test-access-token",
					RefreshToken: "",
				},
			},
			expected: false,
		},
		{
			name: "no tokens",
			config: &config.Config{
				Auth: config.AuthConfig{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAuthenticated(tt.config)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNeedsRefresh(t *testing.T) {
	tests := []struct {
		name     string
		expiresAt time.Time
		expected bool
	}{
		{
			name:     "needs refresh - expired",
			expiresAt: time.Now().Add(-1 * time.Minute),
			expected: true,
		},
		{
			name:     "needs refresh - expires in 3 minutes",
			expiresAt: time.Now().Add(3 * time.Minute),
			expected: true,
		},
		{
			name:     "does not need refresh - expires in 10 minutes",
			expiresAt: time.Now().Add(10 * time.Minute),
			expected: false,
		},
		{
			name:     "does not need refresh - expires in 1 hour",
			expiresAt: time.Now().Add(1 * time.Hour),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Auth: config.AuthConfig{
					ExpiresAt: tt.expiresAt,
				},
			}
			result := NeedsRefresh(cfg)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConfigPersistence(t *testing.T) {
	// Create temp directory for config
	tempDir, err := os.MkdirTemp("", "amazon-cli-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")

	// Create and save a config
	originalConfig := &config.Config{
		Auth: config.AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
		Defaults: config.DefaultsConfig{
			OutputFormat: "json",
		},
		RateLimiting: config.RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	// Save config
	err = config.SaveConfig(originalConfig, configPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Load config back
	loadedConfig, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Verify values
	if loadedConfig.Auth.AccessToken != originalConfig.Auth.AccessToken {
		t.Error("access token mismatch")
	}
	if loadedConfig.Auth.RefreshToken != originalConfig.Auth.RefreshToken {
		t.Error("refresh token mismatch")
	}
	if loadedConfig.Defaults.OutputFormat != originalConfig.Defaults.OutputFormat {
		t.Error("output format mismatch")
	}
	if loadedConfig.RateLimiting.MinDelayMs != originalConfig.RateLimiting.MinDelayMs {
		t.Error("min delay mismatch")
	}
}
