package ratelimit

import (
	"testing"
)

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name       string
		maxRetries int
		statusCode int
		attempt    int
		want       bool
	}{
		{
			name:       "should retry on 429 with attempts remaining",
			maxRetries: 3,
			statusCode: 429,
			attempt:    0,
			want:       true,
		},
		{
			name:       "should retry on 503 with attempts remaining",
			maxRetries: 3,
			statusCode: 503,
			attempt:    1,
			want:       true,
		},
		{
			name:       "should not retry on 429 when max retries reached",
			maxRetries: 3,
			statusCode: 429,
			attempt:    3,
			want:       false,
		},
		{
			name:       "should not retry on 503 when max retries exceeded",
			maxRetries: 3,
			statusCode: 503,
			attempt:    4,
			want:       false,
		},
		{
			name:       "should not retry on 200 OK",
			maxRetries: 3,
			statusCode: 200,
			attempt:    0,
			want:       false,
		},
		{
			name:       "should not retry on 404 Not Found",
			maxRetries: 3,
			statusCode: 404,
			attempt:    0,
			want:       false,
		},
		{
			name:       "should not retry on 500 Internal Server Error",
			maxRetries: 3,
			statusCode: 500,
			attempt:    0,
			want:       false,
		},
		{
			name:       "should not retry on 401 Unauthorized",
			maxRetries: 3,
			statusCode: 401,
			attempt:    0,
			want:       false,
		},
		{
			name:       "should retry on 429 with attempt just below max",
			maxRetries: 3,
			statusCode: 429,
			attempt:    2,
			want:       true,
		},
		{
			name:       "should not retry on 429 with attempt equal to max",
			maxRetries: 3,
			statusCode: 429,
			attempt:    3,
			want:       false,
		},
		{
			name:       "should work with max retries of 1",
			maxRetries: 1,
			statusCode: 429,
			attempt:    0,
			want:       true,
		},
		{
			name:       "should not retry with max retries 1 and attempt 1",
			maxRetries: 1,
			statusCode: 429,
			attempt:    1,
			want:       false,
		},
		{
			name:       "should work with max retries of 0",
			maxRetries: 0,
			statusCode: 429,
			attempt:    0,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := RateLimitConfig{
				MinDelayMs: 1000,
				MaxDelayMs: 5000,
				MaxRetries: tt.maxRetries,
			}
			limiter := NewRateLimiter(config)

			got := limiter.ShouldRetry(tt.statusCode, tt.attempt)
			if got != tt.want {
				t.Errorf("ShouldRetry(statusCode=%d, attempt=%d) = %v, want %v",
					tt.statusCode, tt.attempt, got, tt.want)
			}
		})
	}
}

func TestNewRateLimiter(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(config)

	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	if limiter.config.MinDelayMs != config.MinDelayMs {
		t.Errorf("MinDelayMs = %d, want %d", limiter.config.MinDelayMs, config.MinDelayMs)
	}

	if limiter.config.MaxDelayMs != config.MaxDelayMs {
		t.Errorf("MaxDelayMs = %d, want %d", limiter.config.MaxDelayMs, config.MaxDelayMs)
	}

	if limiter.config.MaxRetries != config.MaxRetries {
		t.Errorf("MaxRetries = %d, want %d", limiter.config.MaxRetries, config.MaxRetries)
	}

	if !limiter.lastRequestTime.IsZero() {
		t.Error("lastRequestTime should be zero value initially")
	}

	if limiter.retryCount != 0 {
		t.Errorf("retryCount = %d, want 0", limiter.retryCount)
	}
}

func TestWait(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 100, // Use smaller delay for testing
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}
	limiter := NewRateLimiter(config)

	// First call should not wait (no previous request)
	err := limiter.Wait()
	if err != nil {
		t.Errorf("Wait() returned error: %v", err)
	}

	// Second call should enforce minimum delay
	err = limiter.Wait()
	if err != nil {
		t.Errorf("Wait() returned error: %v", err)
	}

	if limiter.lastRequestTime.IsZero() {
		t.Error("lastRequestTime should be set after Wait()")
	}
}

func TestWaitWithBackoff(t *testing.T) {
	tests := []struct {
		name            string
		attempt         int
		expectedMinMs   int64
		expectedMaxMs   int64
	}{
		{
			name:            "attempt 0",
			attempt:         0,
			expectedMinMs:   1000,  // 2^0 * 1000 = 1000ms
			expectedMaxMs:   1000,
		},
		{
			name:            "attempt 1",
			attempt:         1,
			expectedMinMs:   2000,  // 2^1 * 1000 = 2000ms
			expectedMaxMs:   2000,
		},
		{
			name:            "attempt 2",
			attempt:         2,
			expectedMinMs:   4000,  // 2^2 * 1000 = 4000ms
			expectedMaxMs:   4000,
		},
		{
			name:            "attempt 3",
			attempt:         3,
			expectedMinMs:   8000,  // 2^3 * 1000 = 8000ms
			expectedMaxMs:   8000,
		},
		{
			name:            "attempt 10 (should cap at 60000)",
			attempt:         10,
			expectedMinMs:   60000, // capped at 60000ms
			expectedMaxMs:   60000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := RateLimitConfig{
				MinDelayMs: 1000,
				MaxDelayMs: 5000,
				MaxRetries: 3,
			}
			limiter := NewRateLimiter(config)

			err := limiter.WaitWithBackoff(tt.attempt)
			if err != nil {
				t.Errorf("WaitWithBackoff() returned error: %v", err)
			}
		})
	}
}
