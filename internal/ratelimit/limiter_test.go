package ratelimit

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestNewRateLimiter(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	if limiter.config.MinDelayMs != 1000 {
		t.Errorf("Expected MinDelayMs to be 1000, got %d", limiter.config.MinDelayMs)
	}

	if limiter.config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", limiter.config.MaxRetries)
	}
}

func TestWait_FirstCall(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100, // Use smaller delay for testing
		MaxDelayMs: 500,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)
	start := time.Now()

	err := limiter.Wait()
	if err != nil {
		t.Fatalf("Wait() returned error: %v", err)
	}

	elapsed := time.Since(start)

	// First call should complete quickly (only jitter, no min delay wait)
	// Jitter is 0-500ms, so should be less than 600ms with overhead
	if elapsed > 600*time.Millisecond {
		t.Errorf("First Wait() took too long: %v", elapsed)
	}
}

func TestWait_MinimumDelay(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100, // 100ms minimum delay
		MaxDelayMs: 500,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// First call to set lastRequestTime
	err := limiter.Wait()
	if err != nil {
		t.Fatalf("First Wait() returned error: %v", err)
	}

	// Second call should wait for minimum delay
	start := time.Now()
	err = limiter.Wait()
	if err != nil {
		t.Fatalf("Second Wait() returned error: %v", err)
	}
	elapsed := time.Since(start)

	// Should take at least MinDelayMs (100ms)
	// Adding jitter (0-500ms), total should be between 100ms and 650ms
	if elapsed < 100*time.Millisecond {
		t.Errorf("Wait() was too fast: %v, expected at least 100ms", elapsed)
	}

	if elapsed > 650*time.Millisecond {
		t.Errorf("Wait() was too slow: %v, expected less than 650ms", elapsed)
	}
}

func TestWait_ConsecutiveCalls(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 50,
		MaxDelayMs: 500,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Make multiple consecutive calls
	for i := 0; i < 3; i++ {
		err := limiter.Wait()
		if err != nil {
			t.Fatalf("Wait() call %d returned error: %v", i, err)
		}
	}

	// If we got here without errors, the rate limiter is working
}

func TestWaitWithBackoff(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	tests := []struct {
		attempt      int
		expectedMin  time.Duration
		expectedMax  time.Duration
		description  string
	}{
		{0, 1000 * time.Millisecond, 1100 * time.Millisecond, "attempt 0: 2^0 * 1000 = 1000ms"},
		{1, 2000 * time.Millisecond, 2100 * time.Millisecond, "attempt 1: 2^1 * 1000 = 2000ms"},
		{2, 4000 * time.Millisecond, 4100 * time.Millisecond, "attempt 2: 2^2 * 1000 = 4000ms"},
		{3, 8000 * time.Millisecond, 8100 * time.Millisecond, "attempt 3: 2^3 * 1000 = 8000ms"},
		{10, 60000 * time.Millisecond, 60100 * time.Millisecond, "attempt 10: capped at 60000ms"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			start := time.Now()
			err := limiter.WaitWithBackoff(tt.attempt)
			if err != nil {
				t.Fatalf("WaitWithBackoff() returned error: %v", err)
			}
			elapsed := time.Since(start)

			if elapsed < tt.expectedMin {
				t.Errorf("Backoff was too short: %v, expected at least %v", elapsed, tt.expectedMin)
			}

			if elapsed > tt.expectedMax {
				t.Errorf("Backoff was too long: %v, expected at most %v", elapsed, tt.expectedMax)
			}
		})
	}
}

func TestShouldRetry(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	tests := []struct {
		statusCode int
		attempt    int
		expected   bool
		description string
	}{
		{429, 0, true, "429 on first attempt should retry"},
		{429, 2, true, "429 on second attempt should retry"},
		{429, 3, false, "429 on third attempt should not retry (max reached)"},
		{429, 4, false, "429 on fourth attempt should not retry (exceeded max)"},
		{503, 0, true, "503 on first attempt should retry"},
		{503, 2, true, "503 on second attempt should retry"},
		{503, 3, false, "503 on third attempt should not retry (max reached)"},
		{404, 0, false, "404 should not retry"},
		{500, 0, false, "500 should not retry"},
		{200, 0, false, "200 should not retry"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := limiter.ShouldRetry(tt.statusCode, tt.attempt)
			if result != tt.expected {
				t.Errorf("ShouldRetry(%d, %d) = %v, expected %v",
					tt.statusCode, tt.attempt, result, tt.expected)
			}
		})
	}
}

func TestResetRetryCount(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Set retry count to non-zero
	limiter.retryCount = 5

	// Reset it
	limiter.ResetRetryCount()

	if limiter.retryCount != 0 {
		t.Errorf("ResetRetryCount() did not reset count, got %d", limiter.retryCount)
	}
}

func TestGenerateJitter(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Test that jitter is within expected range
	for i := 0; i < 100; i++ {
		jitter, err := limiter.generateJitter(500)
		if err != nil {
			t.Fatalf("generateJitter() returned error: %v", err)
		}

		if jitter < 0 || jitter > 500 {
			t.Errorf("Jitter %d out of range [0, 500]", jitter)
		}
	}
}

func TestGenerateJitter_ZeroMax(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	jitter, err := limiter.generateJitter(0)
	if err != nil {
		t.Fatalf("generateJitter(0) returned error: %v", err)
	}

	if jitter != 0 {
		t.Errorf("generateJitter(0) = %d, expected 0", jitter)
	}
}

func TestGenerateJitter_NegativeMax(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	jitter, err := limiter.generateJitter(-100)
	if err != nil {
		t.Fatalf("generateJitter(-100) returned error: %v", err)
	}

	if jitter != 0 {
		t.Errorf("generateJitter(-100) = %d, expected 0", jitter)
	}
}

func TestWait_ThreadSafety(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 10,
		MaxDelayMs: 100,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Run Wait() concurrently from multiple goroutines
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 5; j++ {
				_ = limiter.Wait()
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we got here without deadlock or panic, thread safety is working
}
