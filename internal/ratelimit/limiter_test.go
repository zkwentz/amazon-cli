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

	if limiter.config.MinDelayMs != cfg.MinDelayMs {
		t.Errorf("MinDelayMs = %v, want %v", limiter.config.MinDelayMs, cfg.MinDelayMs)
	}

	if limiter.config.MaxRetries != cfg.MaxRetries {
		t.Errorf("MaxRetries = %v, want %v", limiter.config.MaxRetries, cfg.MaxRetries)
	}

	if !limiter.lastRequestTime.IsZero() {
		t.Error("lastRequestTime should be zero initially")
	}

	if limiter.retryCount != 0 {
		t.Errorf("retryCount = %v, want 0", limiter.retryCount)
	}
}

func TestWait(t *testing.T) {
	t.Run("first wait completes quickly", func(t *testing.T) {
		cfg := config.RateLimitConfig{
			MinDelayMs: 100,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		}

		limiter := NewRateLimiter(cfg)

		start := time.Now()
		err := limiter.Wait()
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("Wait failed: %v", err)
		}

		// First wait should only have jitter (0-500ms)
		if elapsed > 600*time.Millisecond {
			t.Errorf("First wait took %v, expected less than 600ms", elapsed)
		}

		// lastRequestTime should be set
		if limiter.lastRequestTime.IsZero() {
			t.Error("lastRequestTime not set after Wait")
		}
	})

	t.Run("enforces minimum delay between requests", func(t *testing.T) {
		cfg := config.RateLimitConfig{
			MinDelayMs: 200,
			MaxDelayMs: 5000,
			MaxRetries: 3,
		}

		limiter := NewRateLimiter(cfg)

		// First wait
		limiter.Wait()

		// Second wait should enforce delay
		start := time.Now()
		limiter.Wait()
		elapsed := time.Since(start)

		// Should wait at least MinDelayMs (200ms) plus some jitter
		// Being lenient with lower bound due to timing precision
		if elapsed < 150*time.Millisecond {
			t.Errorf("Second wait took %v, expected at least 150ms", elapsed)
		}

		// Should not wait too long (200ms + 500ms jitter max = 700ms)
		if elapsed > 800*time.Millisecond {
			t.Errorf("Second wait took %v, expected less than 800ms", elapsed)
		}
	})
}

func TestWaitWithBackoff(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	tests := []struct {
		attempt         int
		minExpected     time.Duration
		maxExpected     time.Duration
		description     string
	}{
		{0, 900 * time.Millisecond, 1100 * time.Millisecond, "attempt 0 should wait ~1s (2^0 * 1000ms)"},
		{1, 1900 * time.Millisecond, 2100 * time.Millisecond, "attempt 1 should wait ~2s (2^1 * 1000ms)"},
		{2, 3900 * time.Millisecond, 4100 * time.Millisecond, "attempt 2 should wait ~4s (2^2 * 1000ms)"},
		{3, 7900 * time.Millisecond, 8100 * time.Millisecond, "attempt 3 should wait ~8s (2^3 * 1000ms)"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			start := time.Now()
			err := limiter.WaitWithBackoff(tt.attempt)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("WaitWithBackoff failed: %v", err)
			}

			if elapsed < tt.minExpected {
				t.Errorf("WaitWithBackoff took %v, expected at least %v", elapsed, tt.minExpected)
			}

			if elapsed > tt.maxExpected {
				t.Errorf("WaitWithBackoff took %v, expected at most %v", elapsed, tt.maxExpected)
			}
		})
	}
}

func TestWaitWithBackoffMaximum(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 10,
	}

	limiter := NewRateLimiter(cfg)

	// Test with very high attempt number (should cap at 60 seconds)
	start := time.Now()
	err := limiter.WaitWithBackoff(10) // 2^10 * 1000 = 1024000ms, should cap at 60000ms
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("WaitWithBackoff failed: %v", err)
	}

	// Should wait approximately 60 seconds (with some tolerance)
	if elapsed < 59*time.Second {
		t.Errorf("WaitWithBackoff took %v, expected at least 59s", elapsed)
	}

	if elapsed > 61*time.Second {
		t.Errorf("WaitWithBackoff took %v, expected at most 61s", elapsed)
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
		want       bool
		description string
	}{
		{429, 0, true, "429 on first attempt should retry"},
		{429, 1, true, "429 on second attempt should retry"},
		{429, 2, true, "429 on third attempt should retry"},
		{429, 3, false, "429 on fourth attempt (>= MaxRetries) should not retry"},
		{503, 0, true, "503 on first attempt should retry"},
		{503, 2, true, "503 within retry limit should retry"},
		{503, 3, false, "503 exceeding retry limit should not retry"},
		{500, 0, false, "500 should not retry"},
		{404, 0, false, "404 should not retry"},
		{200, 0, false, "200 should not retry"},
		{403, 2, false, "403 should not retry even within limit"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := limiter.ShouldRetry(tt.statusCode, tt.attempt)
			if got != tt.want {
				t.Errorf("ShouldRetry(%d, %d) = %v, want %v", tt.statusCode, tt.attempt, got, tt.want)
			}
		})
	}
}

func TestShouldRetryWithDifferentMaxRetries(t *testing.T) {
	tests := []struct {
		maxRetries int
		attempt    int
		statusCode int
		want       bool
	}{
		{1, 0, 429, true},
		{1, 1, 429, false},
		{5, 4, 429, true},
		{5, 5, 429, false},
		{0, 0, 429, false}, // MaxRetries of 0 means no retries
	}

	for _, tt := range tests {
		cfg := config.RateLimitConfig{
			MinDelayMs: 1000,
			MaxDelayMs: 5000,
			MaxRetries: tt.maxRetries,
		}

		limiter := NewRateLimiter(cfg)
		got := limiter.ShouldRetry(tt.statusCode, tt.attempt)

		if got != tt.want {
			t.Errorf("With MaxRetries=%d, ShouldRetry(%d, %d) = %v, want %v",
				tt.maxRetries, tt.statusCode, tt.attempt, got, tt.want)
		}
	}
}
