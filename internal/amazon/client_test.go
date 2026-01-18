package amazon

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestNewClient verifies client initialization
func TestNewClient(t *testing.T) {
	config := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		},
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.httpClient == nil {
		t.Fatal("Expected non-nil HTTP client")
	}

	if client.rateLimiter == nil {
		t.Fatal("Expected non-nil rate limiter")
	}

	if len(client.userAgents) == 0 {
		t.Fatal("Expected user agents to be populated")
	}
}

// TestRateLimiterWait verifies minimum delay enforcement
func TestRateLimiterWait(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 100, // 100ms for faster test
		MaxDelayMs: 1000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	// First call should complete quickly (no previous request)
	start := time.Now()
	if err := rl.Wait(); err != nil {
		t.Fatalf("Wait failed: %v", err)
	}
	elapsed := time.Since(start)

	// Should complete in under 1 second (jitter is 0-500ms)
	if elapsed > time.Second {
		t.Errorf("First Wait took too long: %v", elapsed)
	}

	// Second call should enforce minimum delay
	start = time.Now()
	if err := rl.Wait(); err != nil {
		t.Fatalf("Wait failed: %v", err)
	}
	elapsed = time.Since(start)

	// Should take at least minDelay (100ms) + jitter (up to 500ms)
	if elapsed < 100*time.Millisecond {
		t.Errorf("Wait did not enforce minimum delay: %v", elapsed)
	}
}

// TestRateLimiterWaitWithBackoff verifies exponential backoff calculation
func TestRateLimiterWaitWithBackoff(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	tests := []struct {
		attempt      int
		minExpected  time.Duration
		maxExpected  time.Duration
		description  string
	}{
		{0, 1 * time.Second, 1 * time.Second, "first attempt: 2^0 * 1000ms = 1s"},
		{1, 2 * time.Second, 2 * time.Second, "second attempt: 2^1 * 1000ms = 2s"},
		{2, 4 * time.Second, 4 * time.Second, "third attempt: 2^2 * 1000ms = 4s"},
		{3, 8 * time.Second, 8 * time.Second, "fourth attempt: 2^3 * 1000ms = 8s"},
		{6, 60 * time.Second, 60 * time.Second, "capped at 60s"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			start := time.Now()
			if err := rl.WaitWithBackoff(tt.attempt); err != nil {
				t.Fatalf("WaitWithBackoff failed: %v", err)
			}
			elapsed := time.Since(start)

			// Allow 50ms margin for execution time
			if elapsed < tt.minExpected-50*time.Millisecond {
				t.Errorf("Backoff too short: expected ~%v, got %v", tt.minExpected, elapsed)
			}
			if elapsed > tt.maxExpected+100*time.Millisecond {
				t.Errorf("Backoff too long: expected ~%v, got %v", tt.maxExpected, elapsed)
			}
		})
	}
}

// TestRateLimiterShouldRetry verifies retry logic
func TestRateLimiterShouldRetry(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	tests := []struct {
		statusCode int
		attempt    int
		expected   bool
		name       string
	}{
		{429, 0, true, "rate limited, first attempt"},
		{429, 2, true, "rate limited, second attempt"},
		{429, 3, false, "rate limited, max retries reached"},
		{503, 0, true, "service unavailable, first attempt"},
		{503, 3, false, "service unavailable, max retries reached"},
		{200, 0, false, "success status"},
		{404, 0, false, "not found"},
		{500, 0, false, "internal server error"},
		{401, 0, false, "unauthorized"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rl.ShouldRetry(tt.statusCode, tt.attempt)
			if result != tt.expected {
				t.Errorf("ShouldRetry(%d, %d) = %v, expected %v",
					tt.statusCode, tt.attempt, result, tt.expected)
			}
		})
	}
}

// TestClientUserAgentRotation verifies user agent rotation
func TestClientUserAgentRotation(t *testing.T) {
	config := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 10, // Short delay for testing
			MaxDelayMs: 100,
			MaxRetries: 3,
		},
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Collect user agents from multiple calls
	seenUserAgents := make(map[string]bool)
	totalAgents := len(client.userAgents)

	for i := 0; i < totalAgents*2; i++ {
		ua := client.getNextUserAgent()
		seenUserAgents[ua] = true
	}

	// Should have seen all user agents
	if len(seenUserAgents) != totalAgents {
		t.Errorf("Expected to see %d unique user agents, got %d", totalAgents, len(seenUserAgents))
	}
}

// TestClientDoSuccess verifies successful request execution
func TestClientDoSuccess(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers are set
		if r.Header.Get("User-Agent") == "" {
			t.Error("User-Agent not set")
		}
		if r.Header.Get("Accept") == "" {
			t.Error("Accept header not set")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	config := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 10,
			MaxDelayMs: 100,
			MaxRetries: 3,
		},
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestClientDoRetry verifies retry logic on 429/503
func TestClientDoRetry(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		if count <= 2 {
			// First two requests return 429
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			// Third request succeeds
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	config := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 10,
			MaxDelayMs: 100,
			MaxRetries: 3,
		},
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if requestCount.Load() != 3 {
		t.Errorf("Expected 3 requests, got %d", requestCount.Load())
	}
}

// TestClientDoMaxRetries verifies max retry limit
func TestClientDoMaxRetries(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		// Always return 429
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	config := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 10,
			MaxDelayMs: 100,
			MaxRetries: 2, // Only allow 2 retries
		},
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	// Should get 429 response after exhausting retries
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", resp.StatusCode)
	}

	// Should make initial request + MaxRetries attempts
	expectedRequests := int32(config.RateLimiting.MaxRetries + 1)
	if requestCount.Load() != expectedRequests {
		t.Errorf("Expected %d requests, got %d", expectedRequests, requestCount.Load())
	}
}

// TestClientDoNoRetryOn404 verifies non-retryable status codes
func TestClientDoNoRetryOn404(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 10,
			MaxDelayMs: 100,
			MaxRetries: 3,
		},
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	// Should only make one request (no retries on 404)
	if requestCount.Load() != 1 {
		t.Errorf("Expected 1 request, got %d", requestCount.Load())
	}
}

// TestClientAuthHeader verifies authorization header is set
func TestClientAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token-123" {
			t.Errorf("Expected Authorization header 'Bearer test-token-123', got '%s'", authHeader)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 10,
			MaxDelayMs: 100,
			MaxRetries: 3,
		},
	}
	config.Auth.AccessToken = "test-token-123"

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	defer resp.Body.Close()
}

// BenchmarkClientDo benchmarks request execution
func BenchmarkClientDo(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &Config{
		RateLimiting: RateLimitConfig{
			MinDelayMs: 1,
			MaxDelayMs: 10,
			MaxRetries: 3,
		},
	}

	client, err := NewClient(config)
	if err != nil {
		b.Fatalf("NewClient failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			b.Fatalf("Get failed: %v", err)
		}
		resp.Body.Close()
	}
}
