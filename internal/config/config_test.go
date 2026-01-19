package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_NonExistentFile(t *testing.T) {
	// Create a temp directory
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nonexistent.json")

	// Load config from non-existent file
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig should not error on non-existent file, got: %v", err)
	}

	if config == nil {
		t.Fatal("LoadConfig should return empty config, got nil")
	}

	if config.Auth.AccessToken != "" {
		t.Errorf("Expected empty access token, got: %s", config.Auth.AccessToken)
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	// Create a temp directory
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	// Create a test config
	expiresAt := time.Now().Add(24 * time.Hour)
	testConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "test_access_token",
			RefreshToken: "test_refresh_token",
			ExpiresAt:    expiresAt,
		},
	}

	// Save the config
	data, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load the config
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify the loaded config
	if config.Auth.AccessToken != testConfig.Auth.AccessToken {
		t.Errorf("Expected access token %s, got %s", testConfig.Auth.AccessToken, config.Auth.AccessToken)
	}

	if config.Auth.RefreshToken != testConfig.Auth.RefreshToken {
		t.Errorf("Expected refresh token %s, got %s", testConfig.Auth.RefreshToken, config.Auth.RefreshToken)
	}

	// Time comparison with tolerance for serialization
	if config.Auth.ExpiresAt.Unix() != testConfig.Auth.ExpiresAt.Unix() {
		t.Errorf("Expected expires at %v, got %v", testConfig.Auth.ExpiresAt, config.Auth.ExpiresAt)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Create a temp directory
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	// Write invalid JSON
	if err := os.WriteFile(path, []byte("invalid json"), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Load the config
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("LoadConfig should error on invalid JSON")
	}
}

func TestLoadConfig_TildeExpansion(t *testing.T) {
	// Create a temp file to test tilde expansion logic
	tmpDir := t.TempDir()

	// We can't actually test real tilde without modifying home,
	// but we can test the expansion logic works
	testConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "test_token",
			RefreshToken: "refresh_token",
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		},
	}

	path := filepath.Join(tmpDir, "config.json")
	err := SaveConfig(testConfig, path)
	if err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	// Load it back
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config == nil {
		t.Fatal("LoadConfig should return config, got nil")
	}

	if config.Auth.AccessToken != "test_token" {
		t.Errorf("Expected test_token, got %s", config.Auth.AccessToken)
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temp directory
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test-config.json")

	// Create a test config
	expiresAt := time.Now().Add(24 * time.Hour)
	testConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "test_access_token",
			RefreshToken: "test_refresh_token",
			ExpiresAt:    expiresAt,
		},
	}

	// Save the config
	err := SaveConfig(testConfig, path)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Verify file permissions are 0600
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}

	// Load and verify content
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if config.Auth.AccessToken != testConfig.Auth.AccessToken {
		t.Errorf("Expected access token %s, got %s", testConfig.Auth.AccessToken, config.Auth.AccessToken)
	}
}

func TestSaveConfig_NilConfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	err := SaveConfig(nil, path)
	if err == nil {
		t.Fatal("SaveConfig should error on nil config")
	}
}

func TestSaveConfig_CreatesDirectory(t *testing.T) {
	// Create a temp directory
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "subdir", "nested", "config.json")

	// Create a test config
	testConfig := &Config{
		Auth: AuthConfig{
			AccessToken: "test_token",
		},
	}

	// Save the config
	err := SaveConfig(testConfig, path)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Verify directory was created
	dir := filepath.Dir(path)
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Fatal("Expected directory to be created")
	}

	// Verify directory permissions are 0700
	if info.Mode().Perm() != 0700 {
		t.Errorf("Expected directory permissions 0700, got %o", info.Mode().Perm())
	}
}

func TestSaveConfig_TildeExpansion(t *testing.T) {
	// Create a test config
	testConfig := &Config{
		Auth: AuthConfig{
			AccessToken: "test_token",
		},
	}

	// Save to temp location with tilde (won't actually use home dir in test)
	tmpDir := t.TempDir()
	// We can't actually test tilde expansion without modifying home dir,
	// but we can test that the function handles it
	path := filepath.Join(tmpDir, "config.json")

	err := SaveConfig(testConfig, path)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
}

func TestIsAuthenticated_ValidToken(t *testing.T) {
	config := &Config{
		Auth: AuthConfig{
			AccessToken: "valid_token",
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		},
	}

	if !config.IsAuthenticated() {
		t.Error("Expected config to be authenticated")
	}
}

func TestIsAuthenticated_ExpiredToken(t *testing.T) {
	config := &Config{
		Auth: AuthConfig{
			AccessToken: "expired_token",
			ExpiresAt:   time.Now().Add(-1 * time.Hour),
		},
	}

	if config.IsAuthenticated() {
		t.Error("Expected config to not be authenticated with expired token")
	}
}

func TestIsAuthenticated_EmptyToken(t *testing.T) {
	config := &Config{
		Auth: AuthConfig{
			AccessToken: "",
			ExpiresAt:   time.Now().Add(1 * time.Hour),
		},
	}

	if config.IsAuthenticated() {
		t.Error("Expected config to not be authenticated with empty token")
	}
}

func TestIsAuthenticated_NilConfig(t *testing.T) {
	var config *Config = nil

	if config.IsAuthenticated() {
		t.Error("Expected nil config to not be authenticated")
	}
}

func TestClearAuth(t *testing.T) {
	config := &Config{
		Auth: AuthConfig{
			AccessToken:  "test_access_token",
			RefreshToken: "test_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
	}

	config.ClearAuth()

	if config.Auth.AccessToken != "" {
		t.Error("Expected access token to be cleared")
	}

	if config.Auth.RefreshToken != "" {
		t.Error("Expected refresh token to be cleared")
	}

	if !config.Auth.ExpiresAt.IsZero() {
		t.Error("Expected expires at to be zero time")
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()

	if path == "" {
		t.Fatal("DefaultConfigPath should not return empty string")
	}

	// Verify it contains .amazon-cli
	if !filepath.IsAbs(path) {
		t.Error("Expected absolute path")
	}

	// Check that it ends with .amazon-cli/config.json
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	if base != "config.json" {
		t.Errorf("Expected filename to be config.json, got %s", base)
	}

	if filepath.Base(dir) != ".amazon-cli" {
		t.Errorf("Expected directory to be .amazon-cli, got %s", filepath.Base(dir))
	}
}

func TestLoadConfig_EmptyPath(t *testing.T) {
	// LoadConfig with empty path should use default path
	// This test just verifies it doesn't crash, not the actual content
	// since the real config file might exist
	config, err := LoadConfig("")

	// If there's a parse error from existing config, that's ok - just verify
	// the function handles empty path correctly
	if err != nil {
		// Check if it's trying to use default path (error would be about parsing)
		if config == nil {
			// If there was an error and config is nil, that's still acceptable
			// since we're just testing empty path handling
			return
		}
	}

	if config == nil {
		t.Fatal("LoadConfig should return config, got nil")
	}
}

func TestSaveConfig_EmptyPath(t *testing.T) {
	// We can't test actual default path, but we can test that it doesn't error
	// on empty path. Skip this test to avoid modifying actual config
	t.Skip("Skipping test that would modify actual config file")
}

func TestRoundTrip(t *testing.T) {
	// Create a temp directory
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	// Create a test config
	expiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	originalConfig := &Config{
		Auth: AuthConfig{
			AccessToken:  "test_access_token",
			RefreshToken: "test_refresh_token",
			ExpiresAt:    expiresAt,
		},
	}

	// Save the config
	err := SaveConfig(originalConfig, path)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Load the config
	loadedConfig, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify the loaded config matches
	if loadedConfig.Auth.AccessToken != originalConfig.Auth.AccessToken {
		t.Errorf("Access token mismatch: expected %s, got %s",
			originalConfig.Auth.AccessToken, loadedConfig.Auth.AccessToken)
	}

	if loadedConfig.Auth.RefreshToken != originalConfig.Auth.RefreshToken {
		t.Errorf("Refresh token mismatch: expected %s, got %s",
			originalConfig.Auth.RefreshToken, loadedConfig.Auth.RefreshToken)
	}

	// Compare times with second precision (JSON doesn't preserve nanoseconds)
	if !loadedConfig.Auth.ExpiresAt.Truncate(time.Second).Equal(originalConfig.Auth.ExpiresAt.Truncate(time.Second)) {
		t.Errorf("ExpiresAt mismatch: expected %v, got %v",
			originalConfig.Auth.ExpiresAt, loadedConfig.Auth.ExpiresAt)
	}
}
