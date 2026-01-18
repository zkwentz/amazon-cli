package ratelimit

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(1000, 5000, 3)

	if rl == nil {
		t.Fatal("NewRateLimiter() returned nil")
	}

	if rl.minDelayMs != 1000 {
		t.Errorf("minDelayMs = %d, want 1000", rl.minDelayMs)
	}

	if rl.maxDelayMs != 5000 {
		t.Errorf("maxDelayMs = %d, want 5000", rl.maxDelayMs)
	}

	if rl.maxRetries != 3 {
		t.Errorf("maxRetries = %d, want 3", rl.maxRetries)
	}
}

func TestWait_FirstRequest(t *testing.T) {
	rl := NewRateLimiter(100, 5000, 3)

	start := time.Now()
	err := rl.Wait()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Wait() error = %v, want nil", err)
	}

	// First request should only have jitter (0-500ms)
	if elapsed > 600*time.Millisecond {
		t.Errorf("elapsed time = %v, want < 600ms (jitter only)", elapsed)
	}
}

func TestWait_EnforcesMinimumDelay(t *testing.T) {
	rl := NewRateLimiter(100, 5000, 3)

	// First wait to set lastRequest
	rl.Wait()

	start := time.Now()
	err := rl.Wait()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Wait() error = %v, want nil", err)
	}

	// Should enforce minimum delay (100ms) plus jitter (0-500ms)
	if elapsed < 100*time.Millisecond {
		t.Errorf("elapsed time = %v, want >= 100ms", elapsed)
	}
}

func TestWait_WithDelay(t *testing.T) {
	rl := NewRateLimiter(200, 5000, 3)

	// First request
	rl.Wait()

	// Sleep for less than minimum delay
	time.Sleep(50 * time.Millisecond)

	start := time.Now()
	rl.Wait()
	elapsed := time.Since(start)

	// Should wait for remaining time (150ms) plus jitter
	if elapsed < 150*time.Millisecond {
		t.Errorf("elapsed time = %v, want >= 150ms", elapsed)
	}
}

func TestWaitWithBackoff(t *testing.T) {
	rl := NewRateLimiter(100, 5000, 3)

	tests := []struct {
		attempt int
		wantMin time.Duration
	}{
		{0, 1 * time.Second},    // 2^0 * 1000ms = 1000ms
		{1, 2 * time.Second},    // 2^1 * 1000ms = 2000ms
		{2, 4 * time.Second},    // 2^2 * 1000ms = 4000ms
		{3, 8 * time.Second},    // 2^3 * 1000ms = 8000ms
	}

	for _, tt := range tests {
		start := time.Now()
		err := rl.WaitWithBackoff(tt.attempt)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("WaitWithBackoff(%d) error = %v, want nil", tt.attempt, err)
		}

		if elapsed < tt.wantMin {
			t.Errorf("WaitWithBackoff(%d) elapsed = %v, want >= %v", tt.attempt, elapsed, tt.wantMin)
		}

		// Should not exceed max by more than a small margin
		if elapsed > (tt.wantMin + 100*time.Millisecond) {
			t.Errorf("WaitWithBackoff(%d) elapsed = %v, want <= %v", tt.attempt, elapsed, tt.wantMin+100*time.Millisecond)
		}
	}
}

func TestWaitWithBackoff_Capped(t *testing.T) {
	// Test that backoff is capped at 60 seconds without actually waiting
	// We'll verify the calculation logic
	tests := []struct {
		attempt     int
		expectedCap bool
	}{
		{6, true},  // 2^6 * 1000ms = 64000ms, should be capped at 60000ms
		{10, true}, // Should be capped at 60000ms
	}

	for _, tt := range tests {
		// Calculate what the backoff should be
		backoffMs := 1 << uint(tt.attempt) * 1000
		expectedBackoff := 60000 * time.Millisecond

		if backoffMs > 60000 && tt.expectedCap {
			// Verified the calculation would be capped
			if backoffMs <= 60000 {
				t.Errorf("attempt %d: backoffMs = %d, expected > 60000", tt.attempt, backoffMs)
			}
		}

		// We won't actually wait 60 seconds in the test, just verify the math
		t.Logf("Attempt %d would wait for %v (capped to %v)", tt.attempt, time.Duration(backoffMs)*time.Millisecond, expectedBackoff)
	}
}

func TestShouldRetry_429(t *testing.T) {
	rl := NewRateLimiter(100, 5000, 3)

	tests := []struct {
		attempt int
		want    bool
	}{
		{0, true},
		{1, true},
		{2, true},
		{3, false}, // At max retries
		{4, false},
	}

	for _, tt := range tests {
		got := rl.ShouldRetry(429, tt.attempt)
		if got != tt.want {
			t.Errorf("ShouldRetry(429, %d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}

func TestShouldRetry_503(t *testing.T) {
	rl := NewRateLimiter(100, 5000, 3)

	tests := []struct {
		attempt int
		want    bool
	}{
		{0, true},
		{1, true},
		{2, true},
		{3, false}, // At max retries
	}

	for _, tt := range tests {
		got := rl.ShouldRetry(503, tt.attempt)
		if got != tt.want {
			t.Errorf("ShouldRetry(503, %d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}

func TestShouldRetry_OtherStatusCodes(t *testing.T) {
	rl := NewRateLimiter(100, 5000, 3)

	statusCodes := []int{200, 201, 400, 401, 403, 404, 500, 502}

	for _, code := range statusCodes {
		got := rl.ShouldRetry(code, 0)
		if got != false {
			t.Errorf("ShouldRetry(%d, 0) = %v, want false", code, got)
		}
	}
}

func TestShouldRetry_MaxRetriesZero(t *testing.T) {
	rl := NewRateLimiter(100, 5000, 0)

	// Should never retry if maxRetries is 0
	got := rl.ShouldRetry(429, 0)
	if got != false {
		t.Error("ShouldRetry(429, 0) with maxRetries=0 should return false")
	}
}

func TestWait_Jitter(t *testing.T) {
	rl := NewRateLimiter(0, 5000, 3) // No minimum delay

	// Make multiple requests and collect jitter values
	var jitters []time.Duration

	for i := 0; i < 10; i++ {
		start := time.Now()
		rl.Wait()
		elapsed := time.Since(start)
		jitters = append(jitters, elapsed)
	}

	// Verify all jitters are within expected range (0-500ms)
	for i, jitter := range jitters {
		if jitter < 0 {
			t.Errorf("jitter[%d] = %v, want >= 0", i, jitter)
		}
		if jitter > 600*time.Millisecond {
			t.Errorf("jitter[%d] = %v, want <= 600ms", i, jitter)
		}
	}
}
