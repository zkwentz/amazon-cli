package ratelimit

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// TestNewRateLimiter tests creation of a new rate limiter
func TestNewRateLimiter(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(cfg)
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	if rl.config.MinDelayMs != cfg.MinDelayMs {
		t.Errorf("Expected MinDelayMs %d, got %d", cfg.MinDelayMs, rl.config.MinDelayMs)
	}

	if rl.lastRequestTime.IsZero() == false {
		t.Error("Expected lastRequestTime to be zero initially")
	}
}

// TestWaitFirstRequest tests that the first request doesn't wait
func TestWaitFirstRequest(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(cfg)
	start := time.Now()
	err := rl.Wait()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}

	// First request should complete quickly (only jitter, no min delay)
	// Jitter is 0-500ms, so expect less than 600ms
	if duration > 600*time.Millisecond {
		t.Errorf("First request took too long: %v", duration)
	}

	if rl.lastRequestTime.IsZero() {
		t.Error("lastRequestTime should be set after Wait()")
	}
}

// TestWaitSubsequentRequests tests that subsequent requests respect MinDelayMs
func TestWaitSubsequentRequests(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 200,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(cfg)

	// First request
	rl.Wait()

	// Second request should wait at least MinDelayMs
	start := time.Now()
	rl.Wait()
	duration := time.Since(start)

	// Should wait at least MinDelayMs (200ms) plus jitter (up to 500ms)
	minExpected := time.Duration(cfg.MinDelayMs) * time.Millisecond
	if duration < minExpected {
		t.Errorf("Expected at least %v delay, got %v", minExpected, duration)
	}

	// Should not wait excessively long (max MinDelayMs + jitter = 700ms)
	maxExpected := time.Duration(cfg.MinDelayMs+500) * time.Millisecond
	if duration > maxExpected {
		t.Errorf("Expected at most %v delay, got %v", maxExpected, duration)
	}

	t.Logf("Wait duration: %v", duration)
}

// TestWaitMultipleRequests tests waiting for multiple consecutive requests
func TestWaitMultipleRequests(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(cfg)
	numRequests := 5
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		err := rl.Wait()
		if err != nil {
			t.Fatalf("Wait %d returned error: %v", i, err)
		}
	}

	totalDuration := time.Since(startTime)

	// Total time should be at least (numRequests - 1) * MinDelayMs
	// (first request doesn't wait, rest do)
	minExpected := time.Duration((numRequests-1)*cfg.MinDelayMs) * time.Millisecond
	if totalDuration < minExpected {
		t.Errorf("Expected at least %v total time, got %v", minExpected, totalDuration)
	}

	t.Logf("Total duration for %d requests: %v", numRequests, totalDuration)
}

// TestWaitWithBackoff tests exponential backoff
func TestWaitWithBackoff(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(cfg)

	tests := []struct {
		attempt      int
		minExpected  time.Duration
		maxExpected  time.Duration
		description  string
	}{
		{0, 1000 * time.Millisecond, 1100 * time.Millisecond, "First backoff (2^0 * 1000ms)"},
		{1, 2000 * time.Millisecond, 2100 * time.Millisecond, "Second backoff (2^1 * 1000ms)"},
		{2, 4000 * time.Millisecond, 4100 * time.Millisecond, "Third backoff (2^2 * 1000ms)"},
		{3, 8000 * time.Millisecond, 8100 * time.Millisecond, "Fourth backoff (2^3 * 1000ms)"},
		{6, 60000 * time.Millisecond, 60100 * time.Millisecond, "Capped backoff (max 60s)"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			start := time.Now()
			err := rl.WaitWithBackoff(tt.attempt)
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("WaitWithBackoff returned error: %v", err)
			}

			if duration < tt.minExpected {
				t.Errorf("Expected at least %v, got %v", tt.minExpected, duration)
			}

			if duration > tt.maxExpected {
				t.Errorf("Expected at most %v, got %v", tt.maxExpected, duration)
			}

			t.Logf("Attempt %d: waited %v", tt.attempt, duration)
		})
	}
}

// TestShouldRetry tests retry logic for different status codes
func TestShouldRetry(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(cfg)

	tests := []struct {
		statusCode     int
		attempt        int
		shouldRetry    bool
		description    string
	}{
		{429, 0, true, "429 on first attempt should retry"},
		{429, 1, true, "429 on second attempt should retry"},
		{429, 2, true, "429 on third attempt should retry"},
		{429, 3, false, "429 on fourth attempt (max retries) should not retry"},
		{503, 0, true, "503 on first attempt should retry"},
		{503, 2, true, "503 on third attempt should retry"},
		{503, 3, false, "503 on fourth attempt (max retries) should not retry"},
		{404, 0, false, "404 should not retry"},
		{500, 0, false, "500 should not retry"},
		{200, 0, false, "200 should not retry"},
		{400, 0, false, "400 should not retry"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := rl.ShouldRetry(tt.statusCode, tt.attempt)
			if result != tt.shouldRetry {
				t.Errorf("Expected ShouldRetry(%d, %d) = %v, got %v",
					tt.statusCode, tt.attempt, tt.shouldRetry, result)
			}
		})
	}
}

// TestShouldRetryMaxRetries tests that max retries are respected
func TestShouldRetryMaxRetries(t *testing.T) {
	tests := []struct {
		maxRetries  int
		attempt     int
		shouldRetry bool
	}{
		{3, 0, true},
		{3, 1, true},
		{3, 2, true},
		{3, 3, false},
		{3, 4, false},
		{1, 0, true},
		{1, 1, false},
		{0, 0, false},
	}

	for _, tt := range tests {
		cfg := config.RateLimitConfig{
			MinDelayMs: 100,
			MaxDelayMs: 5000,
			MaxRetries: tt.maxRetries,
		}
		rl := NewRateLimiter(cfg)

		result := rl.ShouldRetry(429, tt.attempt)
		if result != tt.shouldRetry {
			t.Errorf("MaxRetries=%d, attempt=%d: expected %v, got %v",
				tt.maxRetries, tt.attempt, tt.shouldRetry, result)
		}
	}
}

// BenchmarkWait benchmarks the Wait function
func BenchmarkWait(b *testing.B) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 10, // Short delay for benchmarking
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Wait()
	}
}

// BenchmarkWaitWithBackoff benchmarks the WaitWithBackoff function
func BenchmarkWaitWithBackoff(b *testing.B) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 10,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.WaitWithBackoff(0) // Use attempt 0 for consistency
	}
}
