package amazon

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestNewClient(t *testing.T) {
	cfg := config.NewDefaultConfig()
	client, err := NewClient(cfg)

	if err != nil {
		t.Fatalf("NewClient() error = %v, want nil", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}

	if client.httpClient == nil {
		t.Error("client.httpClient is nil")
	}

	if client.rateLimiter == nil {
		t.Error("client.rateLimiter is nil")
	}

	if client.config == nil {
		t.Error("client.config is nil")
	}

	if len(client.userAgents) == 0 {
		t.Error("client.userAgents is empty")
	}
}

func TestDo_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0 // No delay for testing
	client, _ := NewClient(cfg)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	resp.Body.Close()
}

func TestDo_CustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check custom Accept header is preserved
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Accept header = %s, want application/json", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	client, _ := NewClient(cfg)

	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}

	resp.Body.Close()
}

func TestDo_RetryOn429(t *testing.T) {
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	cfg.RateLimiting.MaxRetries = 3
	client, _ := NewClient(cfg)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}

	if attemptCount != 3 {
		t.Errorf("attemptCount = %d, want 3", attemptCount)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	resp.Body.Close()
}

func TestDo_RetryOn503(t *testing.T) {
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	cfg.RateLimiting.MaxRetries = 3
	client, _ := NewClient(cfg)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}

	if attemptCount != 2 {
		t.Errorf("attemptCount = %d, want 2", attemptCount)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	resp.Body.Close()
}

func TestDo_MaxRetriesExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	cfg.RateLimiting.MaxRetries = 2
	client, _ := NewClient(cfg)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}

	// Should return 429 after max retries
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode, http.StatusTooManyRequests)
	}

	resp.Body.Close()
}

func TestDo_NoRetryOn404(t *testing.T) {
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	client, _ := NewClient(cfg)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}

	// Should only attempt once (no retry on 404)
	if attemptCount != 1 {
		t.Errorf("attemptCount = %d, want 1", attemptCount)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	resp.Body.Close()
}

func TestDo_UserAgentRotation(t *testing.T) {
	userAgents := make(map[string]bool)
	requestCount := 10

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		userAgents[ua] = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	client, _ := NewClient(cfg)

	for i := 0; i < requestCount; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Do() error = %v, want nil", err)
		}
		resp.Body.Close()
	}

	// Should have used at least one user agent
	if len(userAgents) == 0 {
		t.Error("No user agents were used")
	}

	// All user agents should be from our list
	for ua := range userAgents {
		found := false
		for _, validUA := range defaultUserAgents {
			if ua == validUA {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("User agent %s not in default list", ua)
		}
	}
}

func TestGet_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Method = %s, want GET", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	client, _ := NewClient(cfg)

	resp, err := client.Get(server.URL)

	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Get() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	resp.Body.Close()
}

func TestDo_RateLimiting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 100 // 100ms minimum delay
	client, _ := NewClient(cfg)

	start := time.Now()

	// Make two requests
	req1, _ := http.NewRequest("GET", server.URL, nil)
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}
	resp1.Body.Close()

	req2, _ := http.NewRequest("GET", server.URL, nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}
	resp2.Body.Close()

	elapsed := time.Since(start)

	// Should take at least 100ms due to rate limiting (plus jitter)
	if elapsed < 100*time.Millisecond {
		t.Errorf("elapsed time = %v, want >= 100ms", elapsed)
	}
}

func TestGetRandomUserAgent(t *testing.T) {
	cfg := config.NewDefaultConfig()
	client, _ := NewClient(cfg)

	// Call multiple times to ensure it doesn't panic and returns valid UAs
	for i := 0; i < 20; i++ {
		ua := client.getRandomUserAgent()
		if ua == "" {
			t.Error("getRandomUserAgent() returned empty string")
		}

		// Verify it's from the list
		found := false
		for _, validUA := range defaultUserAgents {
			if ua == validUA {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("User agent %s not in default list", ua)
		}
	}
}

func TestDo_InvalidURL(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cfg.RateLimiting.MinDelayMs = 0
	client, _ := NewClient(cfg)

	req, _ := http.NewRequest("GET", "http://invalid-url-that-does-not-exist-12345.com", nil)
	_, err := client.Do(req)

	if err == nil {
		t.Error("Do() with invalid URL should return error")
	}

	if !strings.Contains(err.Error(), "http request failed") {
		t.Errorf("error message = %v, want to contain 'http request failed'", err)
	}
}
