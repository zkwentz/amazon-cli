package amazon

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "creates client with default config",
			config: &config.Config{
				Auth: config.AuthConfig{
					AccessToken:  "test_token",
					RefreshToken: "refresh_token",
					ExpiresAt:    time.Now().Add(1 * time.Hour),
				},
				Defaults: config.DefaultsConfig{
					AddressID:    "addr_123",
					PaymentID:    "pay_123",
					OutputFormat: "json",
				},
				RateLimiting: config.RateLimitConfig{
					MinDelayMs: 1000,
					MaxDelayMs: 5000,
					MaxRetries: 3,
				},
			},
		},
		{
			name: "creates client with custom rate limiting config",
			config: &config.Config{
				Auth: config.AuthConfig{
					AccessToken:  "",
					RefreshToken: "",
					ExpiresAt:    time.Time{},
				},
				Defaults: config.DefaultsConfig{
					AddressID:    "",
					PaymentID:    "",
					OutputFormat: "json",
				},
				RateLimiting: config.RateLimitConfig{
					MinDelayMs: 2000,
					MaxDelayMs: 10000,
					MaxRetries: 5,
				},
			},
		},
		{
			name:   "creates client with nil-like empty config",
			config: config.DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)

			// Verify client is not nil
			if client == nil {
				t.Fatal("NewClient returned nil")
			}

			// Verify httpClient is initialized
			if client.httpClient == nil {
				t.Error("httpClient is nil")
			}

			// Verify httpClient has correct timeout
			expectedTimeout := 30 * time.Second
			if client.httpClient.Timeout != expectedTimeout {
				t.Errorf("httpClient timeout = %v, want %v", client.httpClient.Timeout, expectedTimeout)
			}

			// Verify cookie jar is initialized
			if client.httpClient.Jar == nil {
				t.Error("httpClient cookie jar is nil")
			}

			// Verify rateLimiter is initialized
			if client.rateLimiter == nil {
				t.Error("rateLimiter is nil")
			}

			// Verify config is stored
			if client.config == nil {
				t.Error("config is nil")
			}
			if client.config != tt.config {
				t.Error("config does not match input config")
			}

			// Verify user agents list is populated
			if len(client.userAgents) == 0 {
				t.Error("userAgents list is empty")
			}

			// Verify we have at least 10 user agents (as per PRD)
			if len(client.userAgents) < 10 {
				t.Errorf("userAgents has %d entries, want at least 10", len(client.userAgents))
			}

			// Verify currentUA is initialized to 0
			if client.currentUA != 0 {
				t.Errorf("currentUA = %d, want 0", client.currentUA)
			}

			// Verify all user agents are non-empty strings
			for i, ua := range client.userAgents {
				if ua == "" {
					t.Errorf("userAgent[%d] is empty", i)
				}
			}
		})
	}
}

func TestNewClientIsolation(t *testing.T) {
	// Create two clients with different configs
	config1 := &config.Config{
		RateLimiting: config.RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	config2 := &config.Config{
		RateLimiting: config.RateLimitConfig{
			MinDelayMs: 2000,
			MaxDelayMs: 10000,
			MaxRetries: 5,
		},
	}

	client1 := NewClient(config1)
	client2 := NewClient(config2)

	// Verify clients are separate instances
	if client1 == client2 {
		t.Error("NewClient returned the same instance for different configs")
	}

	// Verify http clients are separate
	if client1.httpClient == client2.httpClient {
		t.Error("httpClient instances are the same")
	}

	// Verify rate limiters are separate
	if client1.rateLimiter == client2.rateLimiter {
		t.Error("rateLimiter instances are the same")
	}

	// Verify configs are separate
	if client1.config == client2.config {
		t.Error("config instances are the same")
	}
}

func TestNewClientUserAgents(t *testing.T) {
	client := NewClient(config.DefaultConfig())

	// Verify all user agents contain expected browser identifiers
	expectedBrowsers := []string{"Chrome", "Firefox", "Safari", "Edg", "OPR"}
	foundBrowsers := make(map[string]bool)

	for _, ua := range client.userAgents {
		for _, browser := range expectedBrowsers {
			if contains(ua, browser) {
				foundBrowsers[browser] = true
			}
		}
	}

	// Should have variety of browsers represented
	if len(foundBrowsers) < 3 {
		t.Errorf("user agents only contain %d different browsers, want at least 3", len(foundBrowsers))
	}

	// Verify all user agents look like valid user agent strings
	for i, ua := range client.userAgents {
		if len(ua) < 50 {
			t.Errorf("userAgent[%d] seems too short: %s", i, ua)
		}
		if !contains(ua, "Mozilla") {
			t.Errorf("userAgent[%d] missing Mozilla prefix: %s", i, ua)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
