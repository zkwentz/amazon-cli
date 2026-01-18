package ratelimit

import (
	"net/http"
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
		t.Fatal("expected limiter to be created")
	}

	if limiter.config.MinDelayMs != 1000 {
		t.Errorf("expected min delay 1000, got %d", limiter.config.MinDelayMs)
	}
}

func TestWait(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100, // Short delay for testing
		MaxDelayMs: 500,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	start := time.Now()
	err := limiter.Wait()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Second wait should enforce minimum delay
	err = limiter.Wait()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	elapsed := time.Since(start)

	// Should take at least MinDelayMs (100ms)
	minExpected := time.Duration(cfg.MinDelayMs) * time.Millisecond
	if elapsed < minExpected {
		t.Errorf("expected at least %v delay, got %v", minExpected, elapsed)
	}
}

func TestWaitWithBackoff(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	tests := []struct {
		attempt      int
		minExpectedMs int
	}{
		{attempt: 0, minExpectedMs: 1000},   // 2^0 * 1000 = 1000
		{attempt: 1, minExpectedMs: 2000},   // 2^1 * 1000 = 2000
		{attempt: 2, minExpectedMs: 4000},   // 2^2 * 1000 = 4000
		{attempt: 6, minExpectedMs: 60000},  // 2^6 * 1000 = 64000, capped at 60000
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.attempt)), func(t *testing.T) {
			start := time.Now()
			err := limiter.WaitWithBackoff(tt.attempt)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			elapsed := time.Since(start)

			minExpected := time.Duration(tt.minExpectedMs) * time.Millisecond
			// Allow some tolerance (10ms)
			tolerance := 10 * time.Millisecond
			if elapsed < minExpected-tolerance {
				t.Errorf("attempt %d: expected at least %v, got %v", tt.attempt, minExpected, elapsed)
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
		name       string
		statusCode int
		attempt    int
		expected   bool
	}{
		{
			name:       "429 on first attempt",
			statusCode: http.StatusTooManyRequests,
			attempt:    0,
			expected:   true,
		},
		{
			name:       "503 on first attempt",
			statusCode: http.StatusServiceUnavailable,
			attempt:    0,
			expected:   true,
		},
		{
			name:       "429 after max retries",
			statusCode: http.StatusTooManyRequests,
			attempt:    3,
			expected:   false,
		},
		{
			name:       "404 on first attempt",
			statusCode: http.StatusNotFound,
			attempt:    0,
			expected:   false,
		},
		{
			name:       "200 on first attempt",
			statusCode: http.StatusOK,
			attempt:    0,
			expected:   false,
		},
		{
			name:       "500 on first attempt",
			statusCode: http.StatusInternalServerError,
			attempt:    0,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := limiter.ShouldRetry(tt.statusCode, tt.attempt)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetJitter(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Test that jitter is within expected range
	for i := 0; i < 100; i++ {
		jitter, err := limiter.getJitter(500)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if jitter < 0 || jitter > 500 {
			t.Errorf("jitter %d out of range [0, 500]", jitter)
		}
	}
}
