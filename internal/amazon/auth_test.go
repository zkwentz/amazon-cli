package amazon

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestNewAuthManager(t *testing.T) {
	cfg := config.GetDefaultConfig()

	am, err := NewAuthManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth manager: %v", err)
	}

	if am == nil {
		t.Fatal("Auth manager is nil")
	}

	if am.httpClient == nil {
		t.Error("HTTP client is nil")
	}

	if am.cookieAuth == nil {
		t.Error("Cookie auth is nil")
	}
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected bool
	}{
		{
			name:     "No auth",
			config:   config.GetDefaultConfig(),
			expected: false,
		},
		{
			name: "Cookie auth",
			config: &config.Config{
				Auth: config.AuthConfig{
					AuthMethod: "cookie",
					Cookies: []config.Cookie{
						{Name: "session", Value: "test"},
					},
				},
			},
			expected: true,
		},
		{
			name: "OAuth auth",
			config: &config.Config{
				Auth: config.AuthConfig{
					AuthMethod:  "oauth",
					AccessToken: "test_token",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am, err := NewAuthManager(tt.config)
			if err != nil {
				t.Fatalf("Failed to create auth manager: %v", err)
			}

			result := am.IsAuthenticated()
			if result != tt.expected {
				t.Errorf("Expected IsAuthenticated() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNeedsRefresh(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected bool
	}{
		{
			name:     "No auth",
			config:   config.GetDefaultConfig(),
			expected: true,
		},
		{
			name: "Fresh cookies",
			config: &config.Config{
				Auth: config.AuthConfig{
					AuthMethod: "cookie",
					Cookies: []config.Cookie{
						{Name: "session", Value: "test"},
					},
					CookiesSetAt: time.Now(),
				},
			},
			expected: false,
		},
		{
			name: "Old cookies",
			config: &config.Config{
				Auth: config.AuthConfig{
					AuthMethod: "cookie",
					Cookies: []config.Cookie{
						{Name: "session", Value: "test"},
					},
					CookiesSetAt: time.Now().Add(-400 * 24 * time.Hour), // Over a year old
				},
			},
			expected: true,
		},
		{
			name: "OAuth expiring soon",
			config: &config.Config{
				Auth: config.AuthConfig{
					AuthMethod:  "oauth",
					AccessToken: "test",
					ExpiresAt:   time.Now().Add(2 * time.Minute), // Expires in 2 minutes
				},
			},
			expected: true,
		},
		{
			name: "OAuth not expiring soon",
			config: &config.Config{
				Auth: config.AuthConfig{
					AuthMethod:  "oauth",
					AccessToken: "test",
					ExpiresAt:   time.Now().Add(30 * time.Minute),
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am, err := NewAuthManager(tt.config)
			if err != nil {
				t.Fatalf("Failed to create auth manager: %v", err)
			}

			result := am.NeedsRefresh()
			if result != tt.expected {
				t.Errorf("Expected NeedsRefresh() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AuthMethod:  "cookie",
			AccessToken: "test",
			Cookies: []config.Cookie{
				{Name: "session", Value: "test"},
			},
			CookiesSetAt: time.Now(),
		},
	}

	am, err := NewAuthManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth manager: %v", err)
	}

	if !am.IsAuthenticated() {
		t.Error("Should be authenticated before logout")
	}

	err = am.Logout()
	if err != nil {
		t.Errorf("Logout failed: %v", err)
	}

	if am.IsAuthenticated() {
		t.Error("Should not be authenticated after logout")
	}

	if len(cfg.Auth.Cookies) != 0 {
		t.Error("Cookies should be cleared after logout")
	}
}

func TestGetAuthStatus(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			AuthMethod: "cookie",
			Cookies: []config.Cookie{
				{Name: "session", Value: "test"},
			},
			CookiesSetAt: time.Now(),
		},
	}

	am, err := NewAuthManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create auth manager: %v", err)
	}

	status := am.GetAuthStatus()

	if status["authenticated"] != true {
		t.Error("Status should show authenticated")
	}

	if status["auth_method"] != "cookie" {
		t.Errorf("Expected auth_method 'cookie', got %v", status["auth_method"])
	}

	if status["cookies_count"] != 1 {
		t.Errorf("Expected 1 cookie, got %v", status["cookies_count"])
	}

	if _, ok := status["cookies_set_at"]; !ok {
		t.Error("Status should include cookies_set_at")
	}
}
