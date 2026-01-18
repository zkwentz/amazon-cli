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

func TestAuthStatus(t *testing.T) {
	// Create a temporary directory for test config
	tmpDir, err := os.MkdirTemp("", "amazon-cli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set HOME to temp directory so config goes there
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	tests := []struct {
		name           string
		setupConfig    func() error
		expectAuth     bool
		expectExpired  *bool
	}{
		{
			name: "no config file",
			setupConfig: func() error {
				return nil
			},
			expectAuth: false,
		},
		{
			name: "no tokens",
			setupConfig: func() error {
				cfg := config.GetDefaultConfig()
				return config.SaveConfig(cfg)
			},
			expectAuth: false,
		},
		{
			name: "valid tokens not expired",
			setupConfig: func() error {
				cfg := config.GetDefaultConfig()
				cfg.Auth.AccessToken = "test-access-token"
				cfg.Auth.RefreshToken = "test-refresh-token"
				cfg.Auth.ExpiresAt = time.Now().Add(1 * time.Hour)
				return config.SaveConfig(cfg)
			},
			expectAuth:    true,
			expectExpired: boolPtr(false),
		},
		{
			name: "valid tokens but expired",
			setupConfig: func() error {
				cfg := config.GetDefaultConfig()
				cfg.Auth.AccessToken = "test-access-token"
				cfg.Auth.RefreshToken = "test-refresh-token"
				cfg.Auth.ExpiresAt = time.Now().Add(-1 * time.Hour)
				return config.SaveConfig(cfg)
			},
			expectAuth:    true,
			expectExpired: boolPtr(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up config directory
			configDir := filepath.Join(tmpDir, ".amazon-cli")
			os.RemoveAll(configDir)

			// Setup config
			if err := tt.setupConfig(); err != nil {
				t.Fatalf("Failed to setup config: %v", err)
			}

			// Capture output
			var buf bytes.Buffer
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run command directly
			err := runAuthStatus(authStatusCmd, []string{})

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout
			buf.ReadFrom(r)

			if err != nil {
				t.Logf("Command output: %s", buf.String())
			}

			// Parse output
			var response map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
			}

			// Check authenticated status
			authenticated, ok := response["authenticated"].(bool)
			if !ok {
				t.Fatalf("authenticated field missing or not a bool")
			}

			if authenticated != tt.expectAuth {
				t.Errorf("Expected authenticated=%v, got %v", tt.expectAuth, authenticated)
			}

			// If we expect authentication, check expired status
			if tt.expectAuth && tt.expectExpired != nil {
				expired, ok := response["expired"].(bool)
				if !ok {
					t.Fatalf("expired field missing or not a bool")
				}

				if expired != *tt.expectExpired {
					t.Errorf("Expected expired=%v, got %v", *tt.expectExpired, expired)
				}

				// Check that expires_at and expires_in_seconds are present
				if _, ok := response["expires_at"].(string); !ok {
					t.Errorf("expires_at field missing or not a string")
				}

				if _, ok := response["expires_in_seconds"].(float64); !ok {
					t.Errorf("expires_in_seconds field missing or not a number")
				}
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
