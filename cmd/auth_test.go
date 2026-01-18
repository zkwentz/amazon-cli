package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/output"
)

func TestLogout(t *testing.T) {
	// Create a temporary directory for test config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create a config with auth tokens
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Now().Add(time.Hour),
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

	// Save the config
	if err := config.SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	// Verify config file exists and has correct permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Config file not created: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Config file has wrong permissions: got %v, want 0600", info.Mode().Perm())
	}

	// Load the config to verify tokens are set
	loadedCfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if loadedCfg.Auth.AccessToken != "test-access-token" {
		t.Errorf("AccessToken not set correctly: got %v, want test-access-token", loadedCfg.Auth.AccessToken)
	}
	if loadedCfg.Auth.RefreshToken != "test-refresh-token" {
		t.Errorf("RefreshToken not set correctly: got %v, want test-refresh-token", loadedCfg.Auth.RefreshToken)
	}

	// Set the config file path for the command
	cfgFile = configPath

	// Initialize printer for the test
	printer = output.NewPrinter("json", false)

	// Execute logout
	if err := runLogout(logoutCmd, []string{}); err != nil {
		t.Fatalf("runLogout failed: %v", err)
	}

	// Load config again and verify tokens are cleared
	loadedCfg, err = config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config after logout: %v", err)
	}

	if loadedCfg.Auth.AccessToken != "" {
		t.Errorf("AccessToken not cleared: got %v, want empty string", loadedCfg.Auth.AccessToken)
	}
	if loadedCfg.Auth.RefreshToken != "" {
		t.Errorf("RefreshToken not cleared: got %v, want empty string", loadedCfg.Auth.RefreshToken)
	}
	if !loadedCfg.Auth.ExpiresAt.IsZero() {
		t.Errorf("ExpiresAt not cleared: got %v, want zero time", loadedCfg.Auth.ExpiresAt)
	}

	// Verify other config sections remain intact
	if loadedCfg.Defaults.OutputFormat != "json" {
		t.Errorf("Defaults were modified: got %v, want json", loadedCfg.Defaults.OutputFormat)
	}
	if loadedCfg.RateLimiting.MinDelayMs != 1000 {
		t.Errorf("RateLimiting was modified: got %v, want 1000", loadedCfg.RateLimiting.MinDelayMs)
	}
}

func TestLogoutWithEmptyConfig(t *testing.T) {
	// Create a temporary directory for test config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Set the config file path for the command
	cfgFile = configPath

	// Initialize printer for the test
	printer = output.NewPrinter("json", false)

	// Run logout on non-existent config (should create empty config)
	if err := runLogout(logoutCmd, []string{}); err != nil {
		t.Fatalf("runLogout failed on empty config: %v", err)
	}

	// Verify config was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}

	// Load and verify
	loadedCfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedCfg.Auth.AccessToken != "" {
		t.Errorf("AccessToken should be empty: got %v", loadedCfg.Auth.AccessToken)
	}
}
