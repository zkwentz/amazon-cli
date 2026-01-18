package amazon

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestRefreshTokenIfNeeded_NoTokens(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Auth.AccessToken = ""
	cfg.Auth.RefreshToken = ""

	err := RefreshTokenIfNeeded(cfg, "")
	if err == nil {
		t.Fatal("expected error when no tokens present")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrCodeAuthRequired {
		t.Errorf("expected error code %s, got %s", models.ErrCodeAuthRequired, cliErr.Code)
	}
}

func TestRefreshTokenIfNeeded_ValidToken(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Auth.AccessToken = "valid-token"
	cfg.Auth.RefreshToken = "refresh-token"
	cfg.Auth.ExpiresAt = time.Now().Add(10 * time.Minute) // Token valid for 10 minutes

	err := RefreshTokenIfNeeded(cfg, "")
	if err != nil {
		t.Fatalf("unexpected error with valid token: %v", err)
	}
}

func TestRefreshTokenIfNeeded_ExpiredNoRefresh(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Auth.AccessToken = "expired-token"
	cfg.Auth.RefreshToken = "" // No refresh token
	cfg.Auth.ExpiresAt = time.Now().Add(-1 * time.Minute) // Token expired

	err := RefreshTokenIfNeeded(cfg, "")
	if err == nil {
		t.Fatal("expected error when token expired and no refresh token")
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}

	if cliErr.Code != models.ErrCodeAuthExpired {
		t.Errorf("expected error code %s, got %s", models.ErrCodeAuthExpired, cliErr.Code)
	}
}

func TestRefreshTokenIfNeeded_ExpiringToken(t *testing.T) {
	// Create a test server to mock token refresh endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/o2/token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "new-access-token",
			"refresh_token": "new-refresh-token",
			"expires_in": 3600,
			"token_type": "Bearer"
		}`))
	}))
	defer server.Close()

	// Temporarily override TokenURL for testing
	originalTokenURL := TokenURL
	TokenURL = server.URL + "/auth/o2/token"
	defer func() { TokenURL = originalTokenURL }()

	cfg := config.GetDefaultConfig()
	cfg.Auth.AccessToken = "expiring-token"
	cfg.Auth.RefreshToken = "valid-refresh-token"
	cfg.Auth.ExpiresAt = time.Now().Add(2 * time.Minute) // Token expires in 2 minutes

	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	err := RefreshTokenIfNeeded(cfg, configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify tokens were updated
	if cfg.Auth.AccessToken != "new-access-token" {
		t.Errorf("expected new access token, got %s", cfg.Auth.AccessToken)
	}

	if cfg.Auth.RefreshToken != "new-refresh-token" {
		t.Errorf("expected new refresh token, got %s", cfg.Auth.RefreshToken)
	}
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Config
		expected bool
	}{
		{
			name: "no token",
			cfg: &config.Config{
				Auth: config.AuthConfig{
					AccessToken: "",
				},
			},
			expected: false,
		},
		{
			name: "valid token",
			cfg: &config.Config{
				Auth: config.AuthConfig{
					AccessToken: "valid-token",
					ExpiresAt:   time.Now().Add(10 * time.Minute),
				},
			},
			expected: true,
		},
		{
			name: "expired token",
			cfg: &config.Config{
				Auth: config.AuthConfig{
					AccessToken: "expired-token",
					ExpiresAt:   time.Now().Add(-1 * time.Minute),
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAuthenticated(tt.cfg)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCheckAuth(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Auth.AccessToken = "valid-token"
	cfg.Auth.RefreshToken = "refresh-token"
	cfg.Auth.ExpiresAt = time.Now().Add(10 * time.Minute)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	err := CheckAuth(cfg, configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRefreshAccessToken_ServerError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Temporarily override TokenURL for testing
	originalTokenURL := TokenURL
	TokenURL = server.URL
	defer func() { TokenURL = originalTokenURL }()

	_, err := refreshAccessToken("test-refresh-token")
	if err == nil {
		t.Fatal("expected error when server returns 500")
	}
}

func TestRefreshAccessToken_Success(t *testing.T) {
	// Create a test server to mock successful token refresh
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "new-token",
			"refresh_token": "new-refresh",
			"expires_in": 3600,
			"token_type": "Bearer"
		}`))
	}))
	defer server.Close()

	// Temporarily override TokenURL for testing
	originalTokenURL := TokenURL
	TokenURL = server.URL
	defer func() { TokenURL = originalTokenURL }()

	tokens, err := refreshAccessToken("test-refresh-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tokens.AccessToken != "new-token" {
		t.Errorf("expected access token 'new-token', got %s", tokens.AccessToken)
	}

	if tokens.RefreshToken != "new-refresh" {
		t.Errorf("expected refresh token 'new-refresh', got %s", tokens.RefreshToken)
	}

	if tokens.ExpiresIn != 3600 {
		t.Errorf("expected expires_in 3600, got %d", tokens.ExpiresIn)
	}
}
