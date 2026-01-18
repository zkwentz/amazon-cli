package amazon

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// TestClientGet tests basic GET requests
func TestClientGet(t *testing.T) {
	// Create mock server
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
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	// Create client with minimal rate limiting for tests
	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10 // Reduce delay for tests
	client := NewClient(cfg)

	// Make request
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Verify response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expected := `{"status": "ok"}`
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestClientRetryOn429 tests that the client retries on 429 responses
func TestClientRetryOn429(t *testing.T) {
	attemptCount := atomic.Int32{}

	// Create mock server that returns 429 twice, then 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := attemptCount.Add(1)
		if count <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limited"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	// Create client with fast retries for testing
	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	cfg.RateLimiting.MaxRetries = 3
	client := NewClient(cfg)

	// Make request
	start := time.Now()
	resp, err := client.Get(server.URL)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Verify response is eventually successful
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 after retries, got %d", resp.StatusCode)
	}

	// Verify we made 3 attempts
	if attemptCount.Load() != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount.Load())
	}

	// Verify backoff was applied (should take at least 2^0 + 2^1 = 3 seconds)
	// With our fast test config, it will be much shorter but still have some delay
	if duration < 10*time.Millisecond {
		t.Errorf("Expected some backoff delay, but request completed in %v", duration)
	}

	t.Logf("Request completed in %v with %d attempts", duration, attemptCount.Load())
}

// TestClientRetryOn503 tests that the client retries on 503 responses
func TestClientRetryOn503(t *testing.T) {
	attemptCount := atomic.Int32{}

	// Create mock server that returns 503 once, then 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := attemptCount.Add(1)
		if count == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "service unavailable"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	cfg.RateLimiting.MaxRetries = 3
	client := NewClient(cfg)

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 after retry, got %d", resp.StatusCode)
	}

	if attemptCount.Load() != 2 {
		t.Errorf("Expected 2 attempts, got %d", attemptCount.Load())
	}
}

// TestClientMaxRetries tests that the client stops retrying after max retries
func TestClientMaxRetries(t *testing.T) {
	attemptCount := atomic.Int32{}

	// Create mock server that always returns 429
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount.Add(1)
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error": "rate limited"}`))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	cfg.RateLimiting.MaxRetries = 2
	client := NewClient(cfg)

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Should return 429 after max retries
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429 after max retries, got %d", resp.StatusCode)
	}

	// Should have made MaxRetries + 1 attempts (initial + retries)
	expectedAttempts := int32(cfg.RateLimiting.MaxRetries + 1)
	if attemptCount.Load() != expectedAttempts {
		t.Errorf("Expected %d attempts, got %d", expectedAttempts, attemptCount.Load())
	}
}

// TestClientNoRetryOn404 tests that the client doesn't retry on 404
func TestClientNoRetryOn404(t *testing.T) {
	attemptCount := atomic.Int32{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount.Add(1)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	cfg.RateLimiting.MaxRetries = 3
	client := NewClient(cfg)

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	// Should only make one attempt (no retries on 404)
	if attemptCount.Load() != 1 {
		t.Errorf("Expected 1 attempt, got %d", attemptCount.Load())
	}
}

// TestClientRateLimiting tests that rate limiting delays are applied
func TestClientRateLimiting(t *testing.T) {
	requestTimes := make([]time.Time, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTimes = append(requestTimes, time.Now())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 100 // 100ms minimum delay
	client := NewClient(cfg)

	// Make 3 requests
	for i := 0; i < 3; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}

	// Verify delays between requests
	if len(requestTimes) != 3 {
		t.Fatalf("Expected 3 request times, got %d", len(requestTimes))
	}

	// Check delay between first and second request
	delay1 := requestTimes[1].Sub(requestTimes[0])
	if delay1 < 100*time.Millisecond {
		t.Errorf("Expected at least 100ms delay between requests 1-2, got %v", delay1)
	}

	// Check delay between second and third request
	delay2 := requestTimes[2].Sub(requestTimes[1])
	if delay2 < 100*time.Millisecond {
		t.Errorf("Expected at least 100ms delay between requests 2-3, got %v", delay2)
	}

	t.Logf("Delays: %v, %v", delay1, delay2)
}

// TestClientUserAgentRotation tests that user agents are rotated
func TestClientUserAgentRotation(t *testing.T) {
	userAgents := make([]string, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgents = append(userAgents, r.Header.Get("User-Agent"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	client := NewClient(cfg)

	// Make multiple requests
	numRequests := 15
	for i := 0; i < numRequests; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}

	// Verify we got different user agents
	if len(userAgents) != numRequests {
		t.Fatalf("Expected %d user agents, got %d", numRequests, len(userAgents))
	}

	// Verify rotation (after 10 requests, we should see the first UA again)
	if userAgents[0] != userAgents[10] {
		t.Error("User agent rotation not working correctly")
	}

	// Count unique user agents
	uniqueUAs := make(map[string]bool)
	for _, ua := range userAgents {
		uniqueUAs[ua] = true
	}

	if len(uniqueUAs) != 10 {
		t.Errorf("Expected 10 unique user agents, got %d", len(uniqueUAs))
	}

	t.Logf("Got %d unique user agents across %d requests", len(uniqueUAs), numRequests)
}

// TestClientPostForm tests POST requests with form data
func TestClientPostForm(t *testing.T) {
	var receivedMethod string
	var receivedContentType string
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedContentType = r.Header.Get("Content-Type")
		receivedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	client := NewClient(cfg)

	// Make POST request
	formData := map[string][]string{
		"key1": {"value1"},
		"key2": {"value2"},
	}
	resp, err := client.PostForm(server.URL, formData)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if receivedMethod != http.MethodPost {
		t.Errorf("Expected POST method, got %s", receivedMethod)
	}

	if receivedContentType != "application/x-www-form-urlencoded" {
		t.Errorf("Expected application/x-www-form-urlencoded content type, got %s", receivedContentType)
	}

	if !strings.Contains(receivedQuery, "key1=value1") || !strings.Contains(receivedQuery, "key2=value2") {
		t.Errorf("Form data not properly encoded in query: %s", receivedQuery)
	}
}

// TestClientCookieJar tests that cookies are persisted across requests
func TestClientCookieJar(t *testing.T) {
	requestCount := atomic.Int32{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count == 1 {
			// First request: set a cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: "test-session-id",
			})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "cookie_set"}`))
		} else {
			// Subsequent requests: verify cookie is sent
			cookie, err := r.Cookie("session")
			if err != nil {
				t.Error("Cookie not sent in subsequent request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if cookie.Value != "test-session-id" {
				t.Errorf("Expected cookie value 'test-session-id', got '%s'", cookie.Value)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "cookie_verified"}`))
		}
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	client := NewClient(cfg)

	// First request - cookie should be set
	resp1, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}
	resp1.Body.Close()

	// Second request - cookie should be sent
	resp2, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}
	resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Error("Cookie was not properly persisted and sent")
	}
}

// TestClientTimeout tests that requests timeout after the configured duration
func TestClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep longer than the client timeout
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	client := NewClient(cfg)

	// Override timeout to be very short for testing
	client.httpClient.Timeout = 100 * time.Millisecond

	_, err := client.Get(server.URL)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "deadline") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

// TestClientConcurrentRequests tests that the client can handle concurrent requests
func TestClientConcurrentRequests(t *testing.T) {
	requestCount := atomic.Int32{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.RateLimiting.MinDelayMs = 10
	client := NewClient(cfg)

	numGoroutines := 5
	done := make(chan bool, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Launch concurrent requests
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			resp, err := client.Get(server.URL)
			if err != nil {
				errors <- fmt.Errorf("goroutine %d error: %v", id, err)
				done <- false
				return
			}
			resp.Body.Close()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errors:
			t.Error(err)
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent requests")
		}
	}

	if requestCount.Load() != int32(numGoroutines) {
		t.Errorf("Expected %d requests, got %d", numGoroutines, requestCount.Load())
	}
}
