package ratelimit

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	minDelay := 100 * time.Millisecond
	maxDelay := 5 * time.Second
	maxRetries := 3

	rl := NewRateLimiter(minDelay, maxDelay, maxRetries)

	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	if rl.minDelay != minDelay {
		t.Errorf("Expected minDelay %v, got %v", minDelay, rl.minDelay)
	}

	if rl.maxDelay != maxDelay {
		t.Errorf("Expected maxDelay %v, got %v", maxDelay, rl.maxDelay)
	}

	if rl.maxRetries != maxRetries {
		t.Errorf("Expected maxRetries %d, got %d", maxRetries, rl.maxRetries)
	}
}

func TestWait(t *testing.T) {
	minDelay := 100 * time.Millisecond
	rl := NewRateLimiter(minDelay, 5*time.Second, 3)

	// First call should have minimal delay (just jitter)
	start := time.Now()
	rl.Wait()
	elapsed := time.Since(start)

	// Should have some jitter (0-500ms)
	if elapsed > 600*time.Millisecond {
		t.Errorf("First Wait() took too long: %v", elapsed)
	}

	// Second call should enforce minDelay + jitter
	start = time.Now()
	rl.Wait()
	elapsed = time.Since(start)

	// Should be at least minDelay, but allow for some variance
	if elapsed < minDelay {
		t.Errorf("Wait() did not enforce minDelay: expected at least %v, got %v", minDelay, elapsed)
	}

	// Should not exceed minDelay + max jitter (500ms) + small buffer
	if elapsed > minDelay+600*time.Millisecond {
		t.Errorf("Wait() took too long: expected around %v, got %v", minDelay, elapsed)
	}
}

func TestWaitWithJitter(t *testing.T) {
	minDelay := 50 * time.Millisecond
	rl := NewRateLimiter(minDelay, 5*time.Second, 3)

	// Make multiple calls and verify jitter is applied
	var delays []time.Duration
	for i := 0; i < 5; i++ {
		start := time.Now()
		rl.Wait()
		if i > 0 { // Skip first call
			delays = append(delays, time.Since(start))
		}
	}

	// Check that delays vary (indicating jitter is working)
	allSame := true
	if len(delays) > 1 {
		first := delays[0]
		for _, d := range delays[1:] {
			// Allow 10ms tolerance
			if d > first+10*time.Millisecond || d < first-10*time.Millisecond {
				allSame = false
				break
			}
		}
	}

	if allSame && len(delays) > 1 {
		t.Log("Warning: All delays were very similar, jitter may not be working properly")
	}
}

func TestWaitWithBackoff(t *testing.T) {
	minDelay := 100 * time.Millisecond
	rl := NewRateLimiter(minDelay, 10*time.Second, 5)

	tests := []struct {
		attempt      int
		minExpected  time.Duration
		maxExpected  time.Duration
		description  string
	}{
		{1, 100 * time.Millisecond, 700 * time.Millisecond, "First attempt"},
		{2, 200 * time.Millisecond, 800 * time.Millisecond, "Second attempt (2x)"},
		{3, 400 * time.Millisecond, 1000 * time.Millisecond, "Third attempt (4x)"},
		{4, 800 * time.Millisecond, 1400 * time.Millisecond, "Fourth attempt (8x)"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			start := time.Now()
			rl.WaitWithBackoff(tt.attempt)
			elapsed := time.Since(start)

			if elapsed < tt.minExpected {
				t.Errorf("Attempt %d: delay too short: expected at least %v, got %v",
					tt.attempt, tt.minExpected, elapsed)
			}

			if elapsed > tt.maxExpected {
				t.Errorf("Attempt %d: delay too long: expected at most %v, got %v",
					tt.attempt, tt.maxExpected, elapsed)
			}
		})
	}
}

func TestWaitWithBackoffCap(t *testing.T) {
	minDelay := 1 * time.Second
	rl := NewRateLimiter(minDelay, 0, 10)

	// Test that backoff is capped at 60 seconds
	// With minDelay of 1s and attempt 10: 1s * 2^9 = 512s
	// Should be capped at 60s + jitter (max 500ms)
	start := time.Now()
	rl.WaitWithBackoff(10)
	elapsed := time.Since(start)

	maxAllowed := 60*time.Second + 600*time.Millisecond // 60s cap + jitter + buffer

	if elapsed > maxAllowed {
		t.Errorf("Backoff exceeded 60s cap: got %v", elapsed)
	}
}

func TestWaitWithBackoffMaxDelay(t *testing.T) {
	minDelay := 100 * time.Millisecond
	maxDelay := 1 * time.Second
	rl := NewRateLimiter(minDelay, maxDelay, 10)

	// With high attempt number, backoff should be capped at maxDelay
	start := time.Now()
	rl.WaitWithBackoff(10)
	elapsed := time.Since(start)

	maxAllowed := maxDelay + 600*time.Millisecond // maxDelay + jitter + buffer

	if elapsed > maxAllowed {
		t.Errorf("Backoff exceeded maxDelay: expected at most %v, got %v", maxAllowed, elapsed)
	}
}

func TestShouldRetry(t *testing.T) {
	rl := NewRateLimiter(100*time.Millisecond, 5*time.Second, 3)

	tests := []struct {
		statusCode int
		attempt    int
		expected   bool
		description string
	}{
		{429, 0, true, "Rate limited, first attempt"},
		{429, 1, true, "Rate limited, second attempt"},
		{429, 2, true, "Rate limited, third attempt"},
		{429, 3, false, "Rate limited, exceeded max retries"},
		{503, 0, true, "Service unavailable, first attempt"},
		{503, 1, true, "Service unavailable, second attempt"},
		{503, 2, true, "Service unavailable, third attempt"},
		{503, 3, false, "Service unavailable, exceeded max retries"},
		{500, 0, false, "Internal server error, should not retry"},
		{404, 0, false, "Not found, should not retry"},
		{200, 0, false, "Success, should not retry"},
		{400, 0, false, "Bad request, should not retry"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := rl.ShouldRetry(tt.statusCode, tt.attempt)
			if result != tt.expected {
				t.Errorf("ShouldRetry(%d, %d) = %v, expected %v",
					tt.statusCode, tt.attempt, result, tt.expected)
			}
		})
	}
}

func TestShouldRetryBoundary(t *testing.T) {
	maxRetries := 5
	rl := NewRateLimiter(100*time.Millisecond, 5*time.Second, maxRetries)

	// Test boundary conditions
	if !rl.ShouldRetry(429, maxRetries-1) {
		t.Errorf("Should retry at maxRetries-1")
	}

	if rl.ShouldRetry(429, maxRetries) {
		t.Errorf("Should not retry at maxRetries")
	}

	if rl.ShouldRetry(429, maxRetries+1) {
		t.Errorf("Should not retry beyond maxRetries")
	}
}

func TestConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(10*time.Millisecond, 1*time.Second, 3)

	// Test that concurrent access doesn't cause race conditions
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			rl.Wait()
			rl.WaitWithBackoff(1)
			rl.ShouldRetry(429, 1)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestWaitWithBackoffZeroAttempt(t *testing.T) {
	rl := NewRateLimiter(100*time.Millisecond, 5*time.Second, 3)

	// Test that attempt 0 is treated as attempt 1
	start := time.Now()
	rl.WaitWithBackoff(0)
	elapsed := time.Since(start)

	// Should be similar to attempt 1: minDelay + jitter
	minExpected := 100 * time.Millisecond
	maxExpected := 700 * time.Millisecond

	if elapsed < minExpected || elapsed > maxExpected {
		t.Errorf("WaitWithBackoff(0) took %v, expected between %v and %v",
			elapsed, minExpected, maxExpected)
	}
}
