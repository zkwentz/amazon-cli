package amazon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		responseStatus int
		wantErr        bool
		checkBody      bool
	}{
		{
			name:           "successful GET request",
			responseBody:   "test response",
			responseStatus: http.StatusOK,
			wantErr:        false,
			checkBody:      true,
		},
		{
			name:           "404 not found",
			responseBody:   "not found",
			responseStatus: http.StatusNotFound,
			wantErr:        false,
			checkBody:      true,
		},
		{
			name:           "500 internal server error",
			responseBody:   "error",
			responseStatus: http.StatusInternalServerError,
			wantErr:        false,
			checkBody:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify it's a GET request
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				// Verify headers are set
				if r.Header.Get("User-Agent") == "" {
					t.Error("User-Agent header not set")
				}
				if r.Header.Get("Accept") == "" {
					t.Error("Accept header not set")
				}
				if r.Header.Get("Accept-Language") == "" {
					t.Error("Accept-Language header not set")
				}

				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create client with minimal rate limiting for faster tests
			cfg := config.DefaultConfig()
			cfg.RateLimiting.MinDelayMs = 0
			cfg.RateLimiting.MaxRetries = 0
			client := NewClient(cfg)

			// Make GET request
			resp, err := client.Get(server.URL)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tt.responseStatus {
				t.Errorf("Expected status %d, got %d", tt.responseStatus, resp.StatusCode)
			}

			// Check body if requested
			if tt.checkBody {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("Failed to read response body: %v", err)
				}
				if string(body) != tt.responseBody {
					t.Errorf("Expected body %q, got %q", tt.responseBody, string(body))
				}
			}
		})
	}
}

func TestClient_Get_InvalidURL(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	client := NewClient(cfg)

	_, err := client.Get("://invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestClient_Get_UserAgentRotation(t *testing.T) {
	requestCount := 0
	userAgents := make([]string, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgents = append(userAgents, r.Header.Get("User-Agent"))
		requestCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	client := NewClient(cfg)

	// Make multiple requests
	for i := 0; i < 5; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}

	// Verify we got all requests
	if requestCount != 5 {
		t.Errorf("Expected 5 requests, got %d", requestCount)
	}

	// Verify user agents are rotating (at least 2 different ones in 5 requests)
	uniqueUA := make(map[string]bool)
	for _, ua := range userAgents {
		uniqueUA[ua] = true
	}

	if len(uniqueUA) < 2 {
		t.Error("Expected user agent rotation, but all requests used the same UA")
	}
}

func TestClient_Get_Retry429(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount <= 2 {
			// Return 429 for first 2 attempts
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited"))
		} else {
			// Success on third attempt
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	cfg.RateLimiting.MaxRetries = 3
	client := NewClient(cfg)

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected successful retry, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 after retry, got %d", resp.StatusCode)
	}

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "success" {
		t.Errorf("Expected 'success', got %q", string(body))
	}
}

func TestClient_Get_Retry503(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("service unavailable"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	cfg.RateLimiting.MaxRetries = 2
	client := NewClient(cfg)

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected successful retry, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 after retry, got %d", resp.StatusCode)
	}

	if attemptCount != 2 {
		t.Errorf("Expected 2 attempts, got %d", attemptCount)
	}
}

func TestClient_Get_MaxRetriesExceeded(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("rate limited"))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	cfg.RateLimiting.MaxRetries = 2
	client := NewClient(cfg)

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected response even after max retries, got error: %v", err)
	}
	defer resp.Body.Close()

	// Should still return the last 429 response
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected final status 429, got %d", resp.StatusCode)
	}

	// Should have attempted maxRetries + 1 times (initial + retries)
	expectedAttempts := cfg.RateLimiting.MaxRetries + 1
	if attemptCount != expectedAttempts {
		t.Errorf("Expected %d attempts, got %d", expectedAttempts, attemptCount)
	}
}
