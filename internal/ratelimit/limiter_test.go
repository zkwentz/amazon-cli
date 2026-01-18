package ratelimit

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

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
		expectedDelayMs float64
		maxCapMs        float64
	}{
		{
			name:            "attempt 0 - 1 second delay",
			attempt:         0,
			expectedDelayMs: 1000,  // 2^0 * 1000 = 1000ms
			maxCapMs:        60000,
		},
		{
			name:            "attempt 1 - 2 second delay",
			attempt:         1,
			expectedDelayMs: 2000,  // 2^1 * 1000 = 2000ms
			maxCapMs:        60000,
		},
		{
			name:            "attempt 2 - 4 second delay",
			attempt:         2,
			expectedDelayMs: 4000,  // 2^2 * 1000 = 4000ms
			maxCapMs:        60000,
		},
		{
			name:            "attempt 3 - 8 second delay",
			attempt:         3,
			expectedDelayMs: 8000,  // 2^3 * 1000 = 8000ms
			maxCapMs:        60000,
		},
		{
			name:            "attempt 4 - 16 second delay",
			attempt:         4,
			expectedDelayMs: 16000, // 2^4 * 1000 = 16000ms
			maxCapMs:        60000,
		},
		{
			name:            "attempt 5 - 32 second delay",
			attempt:         5,
			expectedDelayMs: 32000, // 2^5 * 1000 = 32000ms
			maxCapMs:        60000,
		},
		{
			name:            "attempt 6 - capped at 60 seconds",
			attempt:         6,
			expectedDelayMs: 60000, // 2^6 * 1000 = 64000ms, but capped at 60000ms
			maxCapMs:        60000,
		},
		{
			name:            "attempt 10 - capped at 60 seconds",
			attempt:         10,
			expectedDelayMs: 60000, // 2^10 * 1000 = 1024000ms, but capped at 60000ms
			maxCapMs:        60000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			err := rl.WaitWithBackoff(tt.attempt)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("WaitWithBackoff returned error: %v", err)
			}

			// Allow for 100ms tolerance due to system scheduling
			expectedDuration := time.Duration(tt.expectedDelayMs) * time.Millisecond
			tolerance := 100 * time.Millisecond

			if elapsed < expectedDuration-tolerance {
				t.Errorf("WaitWithBackoff slept too short: got %v, expected at least %v", elapsed, expectedDuration-tolerance)
			}

			// Check that it doesn't sleep way too long (accounting for system overhead)
			maxDuration := time.Duration(tt.maxCapMs)*time.Millisecond + tolerance
			if elapsed > maxDuration {
				t.Errorf("WaitWithBackoff slept too long: got %v, expected at most %v", elapsed, maxDuration)
			}
		})
	}
}

func TestWaitWithBackoff_NegativeAttempt(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}
	rl := NewRateLimiter(cfg)

	err := rl.WaitWithBackoff(-1)
	if err == nil {
		t.Error("Expected error for negative attempt, got nil")
	}

	expectedError := "attempt must be non-negative, got: -1"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, err.Error())
	}
}

func TestWaitWithBackoff_UpdatesLastRequestTime(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}
	rl := NewRateLimiter(cfg)

	// Initially, lastRequestTime should be zero
	if !rl.lastRequestTime.IsZero() {
		t.Error("Expected lastRequestTime to be zero initially")
	}

	// After WaitWithBackoff, lastRequestTime should be updated
	before := time.Now()
	err := rl.WaitWithBackoff(0)
	after := time.Now()

	if err != nil {
		t.Fatalf("WaitWithBackoff returned error: %v", err)
	}

	if rl.lastRequestTime.IsZero() {
		t.Error("Expected lastRequestTime to be updated after WaitWithBackoff")
	}

	// lastRequestTime should be between before and after
	if rl.lastRequestTime.Before(before) || rl.lastRequestTime.After(after) {
		t.Errorf("lastRequestTime %v not in expected range [%v, %v]", rl.lastRequestTime, before, after)
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
		expected   bool
	}{
		{
			name:       "429 rate limited - attempt 0",
			statusCode: 429,
			attempt:    0,
			expected:   true,
		},
		{
			name:       "503 service unavailable - attempt 0",
			statusCode: 503,
			attempt:    0,
			expected:   true,
		},
		{
			name:       "429 rate limited - attempt at max retries",
			statusCode: 429,
			attempt:    3,
			expected:   false,
		},
		{
			name:       "429 rate limited - attempt beyond max retries",
			statusCode: 429,
			attempt:    4,
			expected:   false,
		},
		{
			name:       "404 not found - should not retry",
			statusCode: 404,
			attempt:    0,
			expected:   false,
		},
		{
			name:       "500 internal server error - should not retry",
			statusCode: 500,
			attempt:    0,
			expected:   false,
		},
		{
			name:       "200 OK - should not retry",
			statusCode: 200,
			attempt:    0,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rl.ShouldRetry(tt.statusCode, tt.attempt)
			if result != tt.expected {
				t.Errorf("ShouldRetry(%d, %d) = %v, expected %v", tt.statusCode, tt.attempt, result, tt.expected)
			}
		})
	}
}

func TestWait(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100, // Use shorter delay for tests
		MaxDelayMs: 500,
		MaxRetries: 3,
	}
	rl := NewRateLimiter(cfg)

	// First call should not wait (no previous request)
	start := time.Now()
	err := rl.Wait()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}

	// Should be nearly instant (within 50ms)
	if elapsed > 50*time.Millisecond {
		t.Errorf("First Wait() took too long: %v", elapsed)
	}

	// Second call should wait at least MinDelayMs + some jitter
	start = time.Now()
	err = rl.Wait()
	elapsed = time.Since(start)

	if err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}

	// Should wait at least MinDelayMs (100ms) but account for jitter (0-500ms)
	minExpected := 100 * time.Millisecond
	maxExpected := 700 * time.Millisecond // 100ms + 500ms jitter + 100ms tolerance

	if elapsed < minExpected {
		t.Errorf("Second Wait() too short: got %v, expected at least %v", elapsed, minExpected)
	}

	if elapsed > maxExpected {
		t.Errorf("Second Wait() too long: got %v, expected at most %v", elapsed, maxExpected)
	}
}

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

	if rl.config.MinDelayMs != 1000 {
		t.Errorf("Expected MinDelayMs to be 1000, got %d", rl.config.MinDelayMs)
	}

	if rl.config.MaxDelayMs != 5000 {
		t.Errorf("Expected MaxDelayMs to be 5000, got %d", rl.config.MaxDelayMs)
	}

	if rl.config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", rl.config.MaxRetries)
	}

	if !rl.lastRequestTime.IsZero() {
		t.Error("Expected lastRequestTime to be zero for new RateLimiter")
	}

	if rl.verbose {
		t.Error("Expected verbose to be false by default")
	}
}

func TestSetVerbose(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}
	rl := NewRateLimiter(cfg)

	if rl.verbose {
		t.Error("Expected verbose to be false initially")
	}

	rl.SetVerbose(true)
	if !rl.verbose {
		t.Error("Expected verbose to be true after SetVerbose(true)")
	}

	rl.SetVerbose(false)
	if rl.verbose {
		t.Error("Expected verbose to be false after SetVerbose(false)")
	}
}
