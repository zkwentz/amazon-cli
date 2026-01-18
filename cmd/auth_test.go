package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestAuthStatusNotAuthenticated(t *testing.T) {
	// Create temp directory for config
	tmpDir := t.TempDir()
	cfgFile = filepath.Join(tmpDir, "config.json")

	// Create empty config
	cfg := &config.Config{
		Auth: config.AuthConfig{},
		Defaults: config.DefaultsConfig{
			OutputFormat: "json",
		},
		RateLimiting: config.RateLimitingConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}
	if err := config.SaveConfig(cfg, cfgFile); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run status command
	err := runStatus(nil, nil)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runStatus failed: %v", err)
	}

	// Read output
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify not authenticated
	authenticated, ok := result["authenticated"].(bool)
	if !ok {
		t.Fatal("authenticated field missing or not a bool")
	}
	if authenticated {
		t.Error("Expected authenticated to be false")
	}
}

func TestAuthStatusAuthenticated(t *testing.T) {
	// Create temp directory for config
	tmpDir := t.TempDir()
	cfgFile = filepath.Join(tmpDir, "config.json")

	// Create config with auth tokens
	expiresAt := time.Now().Add(1 * time.Hour)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    expiresAt,
		},
		Defaults: config.DefaultsConfig{
			OutputFormat: "json",
		},
		RateLimiting: config.RateLimitingConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}
	if err := config.SaveConfig(cfg, cfgFile); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run status command
	err := runStatus(nil, nil)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runStatus failed: %v", err)
	}

	// Read output
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify authenticated
	authenticated, ok := result["authenticated"].(bool)
	if !ok {
		t.Fatal("authenticated field missing or not a bool")
	}
	if !authenticated {
		t.Error("Expected authenticated to be true")
	}

	// Verify has expiry info
	if _, ok := result["expires_at"]; !ok {
		t.Error("Expected expires_at field")
	}

	if _, ok := result["expires_in_seconds"]; !ok {
		t.Error("Expected expires_in_seconds field")
	}
}

func TestAuthLogout(t *testing.T) {
	// Create temp directory for config
	tmpDir := t.TempDir()
	cfgFile = filepath.Join(tmpDir, "config.json")

	// Create config with auth tokens
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
		Defaults: config.DefaultsConfig{
			OutputFormat: "json",
		},
		RateLimiting: config.RateLimitingConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}
	if err := config.SaveConfig(cfg, cfgFile); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run logout command
	err := runLogout(nil, nil)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runLogout failed: %v", err)
	}

	// Read output
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify logout status
	status, ok := result["status"].(string)
	if !ok {
		t.Fatal("status field missing or not a string")
	}
	if status != "logged_out" {
		t.Errorf("Expected status 'logged_out', got '%s'", status)
	}

	// Verify config was updated
	loadedCfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedCfg.Auth.AccessToken != "" {
		t.Error("Expected access token to be cleared")
	}

	if loadedCfg.Auth.RefreshToken != "" {
		t.Error("Expected refresh token to be cleared")
	}
}
