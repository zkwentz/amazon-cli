package ratelimit

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

func TestNewRateLimiter(t *testing.T) {
	tests := []struct {
		name   string
		config config.RateLimitConfig
	}{
		{
			name: "default configuration",
			config: config.RateLimitConfig{
				MinDelayMs: 1000,
				MaxDelayMs: 5000,
				MaxRetries: 3,
			},
		},
		{
			name: "custom configuration",
			config: config.RateLimitConfig{
				MinDelayMs: 500,
				MaxDelayMs: 3000,
				MaxRetries: 5,
			},
		},
		{
			name: "zero values",
			config: config.RateLimitConfig{
				MinDelayMs: 0,
				MaxDelayMs: 0,
				MaxRetries: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := NewRateLimiter(tt.config)

			if rl == nil {
				t.Fatal("NewRateLimiter returned nil")
			}

			if rl.config.MinDelayMs != tt.config.MinDelayMs {
				t.Errorf("MinDelayMs = %d, want %d", rl.config.MinDelayMs, tt.config.MinDelayMs)
			}

			if rl.config.MaxDelayMs != tt.config.MaxDelayMs {
				t.Errorf("MaxDelayMs = %d, want %d", rl.config.MaxDelayMs, tt.config.MaxDelayMs)
			}

			if rl.config.MaxRetries != tt.config.MaxRetries {
				t.Errorf("MaxRetries = %d, want %d", rl.config.MaxRetries, tt.config.MaxRetries)
			}

			if !rl.lastRequestTime.IsZero() {
				t.Errorf("lastRequestTime should be zero value, got %v", rl.lastRequestTime)
			}

			if rl.retryCount != 0 {
				t.Errorf("retryCount = %d, want 0", rl.retryCount)
			}
		})
	}
}

func TestWait(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100, // Use small value for testing
		MaxDelayMs: 500,
		MaxRetries: 3,
	}
	rl := NewRateLimiter(cfg)

	// First call should wait only for jitter (0-500ms)
	start := time.Now()
	err := rl.Wait()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Wait() returned error: %v", err)
	}

	// Should have at least some jitter
	if elapsed < 0 {
		t.Errorf("elapsed time %v is negative", elapsed)
	}

	// Second call should wait for minimum delay + jitter
	start = time.Now()
	err = rl.Wait()
	elapsed = time.Since(start)

	if err != nil {
		t.Fatalf("Wait() returned error on second call: %v", err)
	}

	// Should have waited at least minDelay (100ms) + some jitter
	// Using 80ms threshold to account for timing variations
	if elapsed < 80*time.Millisecond {
		t.Errorf("elapsed time %v is less than expected minimum", elapsed)
	}
}

func TestWaitWithBackoff(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}
	rl := NewRateLimiter(cfg)

	tests := []struct {
		name            string
		attempt         int
		minExpectedMs   int64
		maxExpectedMs   int64
	}{
		{
			name:            "first retry (2^0 * 1000)",
			attempt:         0,
			minExpectedMs:   900,  // Allow 10% margin
			maxExpectedMs:   1100,
		},
		{
			name:            "second retry (2^1 * 1000)",
			attempt:         1,
			minExpectedMs:   1800,
			maxExpectedMs:   2200,
		},
		{
			name:            "third retry (2^2 * 1000)",
			attempt:         2,
			minExpectedMs:   3600,
			maxExpectedMs:   4400,
		},
		{
			name:            "max backoff capped at 60s",
			attempt:         10,
			minExpectedMs:   59000,
			maxExpectedMs:   61000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			err := rl.WaitWithBackoff(tt.attempt)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("WaitWithBackoff(%d) returned error: %v", tt.attempt, err)
			}

			elapsedMs := elapsed.Milliseconds()
			if elapsedMs < tt.minExpectedMs || elapsedMs > tt.maxExpectedMs {
				t.Errorf("WaitWithBackoff(%d) elapsed %dms, want between %dms and %dms",
					tt.attempt, elapsedMs, tt.minExpectedMs, tt.maxExpectedMs)
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
	rl := NewRateLimiter(cfg)

	tests := []struct {
		name       string
		statusCode int
		attempt    int
		want       bool
	}{
		{
			name:       "429 rate limited, first attempt",
			statusCode: 429,
			attempt:    0,
			want:       true,
		},
		{
			name:       "503 service unavailable, first attempt",
			statusCode: 503,
			attempt:    0,
			want:       true,
		},
		{
			name:       "429 rate limited, max retries reached",
			statusCode: 429,
			attempt:    3,
			want:       false,
		},
		{
			name:       "500 internal server error, should not retry",
			statusCode: 500,
			attempt:    0,
			want:       false,
		},
		{
			name:       "404 not found, should not retry",
			statusCode: 404,
			attempt:    0,
			want:       false,
		},
		{
			name:       "200 success, should not retry",
			statusCode: 200,
			attempt:    0,
			want:       false,
		},
		{
			name:       "503 service unavailable, exceeded max retries",
			statusCode: 503,
			attempt:    5,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rl.ShouldRetry(tt.statusCode, tt.attempt)
			if got != tt.want {
				t.Errorf("ShouldRetry(%d, %d) = %v, want %v",
					tt.statusCode, tt.attempt, got, tt.want)
			}
		})
	}
}

func TestRetryCountMethods(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}
	rl := NewRateLimiter(cfg)

	// Initial count should be 0
	if count := rl.GetRetryCount(); count != 0 {
		t.Errorf("initial GetRetryCount() = %d, want 0", count)
	}

	// Increment retry count
	rl.IncrementRetryCount()
	if count := rl.GetRetryCount(); count != 1 {
		t.Errorf("GetRetryCount() after increment = %d, want 1", count)
	}

	// Increment again
	rl.IncrementRetryCount()
	if count := rl.GetRetryCount(); count != 2 {
		t.Errorf("GetRetryCount() after second increment = %d, want 2", count)
	}

	// Reset retry count
	rl.ResetRetryCount()
	if count := rl.GetRetryCount(); count != 0 {
		t.Errorf("GetRetryCount() after reset = %d, want 0", count)
	}
}

func TestRandomJitter(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}
	rl := NewRateLimiter(cfg)

	maxMs := 500
	for i := 0; i < 100; i++ {
		jitter, err := rl.randomJitter(maxMs)
		if err != nil {
			t.Fatalf("randomJitter(%d) returned error: %v", maxMs, err)
		}

		if jitter < 0 {
			t.Errorf("randomJitter(%d) = %v, should be non-negative", maxMs, jitter)
		}

		if jitter > time.Duration(maxMs)*time.Millisecond {
			t.Errorf("randomJitter(%d) = %v, exceeds maximum", maxMs, jitter)
		}
	}
}
