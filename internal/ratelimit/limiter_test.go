package ratelimit

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	if rl.config.MinDelayMs != config.MinDelayMs {
		t.Errorf("Expected MinDelayMs %d, got %d", config.MinDelayMs, rl.config.MinDelayMs)
	}

	if rl.config.MaxDelayMs != config.MaxDelayMs {
		t.Errorf("Expected MaxDelayMs %d, got %d", config.MaxDelayMs, rl.config.MaxDelayMs)
	}

	if rl.config.MaxRetries != config.MaxRetries {
		t.Errorf("Expected MaxRetries %d, got %d", config.MaxRetries, rl.config.MaxRetries)
	}

	if !rl.lastRequestTime.IsZero() {
		t.Error("Expected lastRequestTime to be zero initially")
	}
}

func TestWait_FirstRequest(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	start := time.Now()
	err := rl.Wait()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Wait() returned error: %v", err)
	}

	// First request should not wait (or minimal time for jitter)
	if elapsed > 100*time.Millisecond {
		t.Errorf("First Wait() took too long: %v", elapsed)
	}

	if rl.lastRequestTime.IsZero() {
		t.Error("lastRequestTime should be set after Wait()")
	}
}

func TestWait_EnforcesMinimumDelay(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 100, // Short delay for testing
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	// First request
	err := rl.Wait()
	if err != nil {
		t.Fatalf("First Wait() returned error: %v", err)
	}

	// Second request should enforce delay
	start := time.Now()
	err = rl.Wait()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Second Wait() returned error: %v", err)
	}

	// Should wait at least MinDelayMs (plus jitter up to 500ms)
	minExpected := time.Duration(config.MinDelayMs) * time.Millisecond
	maxExpected := minExpected + 600*time.Millisecond // MinDelay + 500ms jitter + tolerance

	if elapsed < minExpected {
		t.Errorf("Wait() did not enforce minimum delay. Expected >= %v, got %v", minExpected, elapsed)
	}

	if elapsed > maxExpected {
		t.Errorf("Wait() took too long. Expected <= %v, got %v", maxExpected, elapsed)
	}
}

func TestWait_WithPreExistingDelay(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 50, // Very short for testing
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	// First request
	err := rl.Wait()
	if err != nil {
		t.Fatalf("First Wait() returned error: %v", err)
	}

	// Wait longer than MinDelayMs before second request
	time.Sleep(100 * time.Millisecond)

	// Second request should only add jitter, not wait for MinDelayMs
	start := time.Now()
	err = rl.Wait()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Second Wait() returned error: %v", err)
	}

	// Should only wait for jitter (0-500ms)
	maxExpected := 600 * time.Millisecond // 500ms jitter + 100ms tolerance

	if elapsed > maxExpected {
		t.Errorf("Wait() took too long when delay already satisfied. Expected <= %v, got %v", maxExpected, elapsed)
	}
}

func TestWaitWithBackoff(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	tests := []struct {
		attempt      int
		expectedMin  time.Duration
		expectedMax  time.Duration
		description  string
	}{
		{0, 1000 * time.Millisecond, 1100 * time.Millisecond, "First retry: 2^0 * 1000ms = 1000ms"},
		{1, 2000 * time.Millisecond, 2100 * time.Millisecond, "Second retry: 2^1 * 1000ms = 2000ms"},
		{2, 4000 * time.Millisecond, 4100 * time.Millisecond, "Third retry: 2^2 * 1000ms = 4000ms"},
		{3, 8000 * time.Millisecond, 8100 * time.Millisecond, "Fourth retry: 2^3 * 1000ms = 8000ms"},
		{6, 60000 * time.Millisecond, 60100 * time.Millisecond, "Max backoff: capped at 60000ms"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			start := time.Now()
			err := rl.WaitWithBackoff(tt.attempt)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("WaitWithBackoff(%d) returned error: %v", tt.attempt, err)
			}

			if elapsed < tt.expectedMin {
				t.Errorf("WaitWithBackoff(%d) too short. Expected >= %v, got %v",
					tt.attempt, tt.expectedMin, elapsed)
			}

			if elapsed > tt.expectedMax {
				t.Errorf("WaitWithBackoff(%d) too long. Expected <= %v, got %v",
					tt.attempt, tt.expectedMax, elapsed)
			}
		})
	}
}

func TestShouldRetry_StatusCodes(t *testing.T) {
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
		description string
	}{
		{429, 0, true, "429 (rate limited) on first attempt should retry"},
		{503, 0, true, "503 (service unavailable) on first attempt should retry"},
		{500, 0, false, "500 (internal server error) should not retry"},
		{404, 0, false, "404 (not found) should not retry"},
		{200, 0, false, "200 (success) should not retry"},
		{429, 2, true, "429 on third attempt (< MaxRetries) should retry"},
		{429, 3, false, "429 on fourth attempt (>= MaxRetries) should not retry"},
		{503, 3, false, "503 on fourth attempt (>= MaxRetries) should not retry"},
		{429, 5, false, "429 after MaxRetries should not retry"},
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

func TestGenerateJitter(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	// Test jitter generation multiple times
	for i := 0; i < 100; i++ {
		jitter, err := rl.generateJitter(500)
		if err != nil {
			t.Fatalf("generateJitter returned error: %v", err)
		}

		if jitter < 0 {
			t.Errorf("Jitter is negative: %v", jitter)
		}

		if jitter > 500*time.Millisecond {
			t.Errorf("Jitter exceeds maximum: %v > 500ms", jitter)
		}
	}
}

func TestGenerateJitter_Zero(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	jitter, err := rl.generateJitter(0)
	if err != nil {
		t.Fatalf("generateJitter(0) returned error: %v", err)
	}

	if jitter != 0 {
		t.Errorf("Expected zero jitter for maxMs=0, got %v", jitter)
	}
}

func TestGenerateJitter_Negative(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	jitter, err := rl.generateJitter(-100)
	if err != nil {
		t.Fatalf("generateJitter(-100) returned error: %v", err)
	}

	if jitter != 0 {
		t.Errorf("Expected zero jitter for negative maxMs, got %v", jitter)
	}
}

func TestReset(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	// Make a request to set lastRequestTime
	err := rl.Wait()
	if err != nil {
		t.Fatalf("Wait() returned error: %v", err)
	}

	if rl.lastRequestTime.IsZero() {
		t.Error("lastRequestTime should be set after Wait()")
	}

	// Reset the limiter
	rl.Reset()

	if !rl.lastRequestTime.IsZero() {
		t.Error("lastRequestTime should be zero after Reset()")
	}
}

func TestGetConfig(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	retrievedConfig := rl.GetConfig()

	if retrievedConfig.MinDelayMs != config.MinDelayMs {
		t.Errorf("GetConfig() MinDelayMs mismatch: expected %d, got %d",
			config.MinDelayMs, retrievedConfig.MinDelayMs)
	}

	if retrievedConfig.MaxDelayMs != config.MaxDelayMs {
		t.Errorf("GetConfig() MaxDelayMs mismatch: expected %d, got %d",
			config.MaxDelayMs, retrievedConfig.MaxDelayMs)
	}

	if retrievedConfig.MaxRetries != config.MaxRetries {
		t.Errorf("GetConfig() MaxRetries mismatch: expected %d, got %d",
			config.MaxRetries, retrievedConfig.MaxRetries)
	}
}

func TestSetVerbose(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

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

func TestConcurrentWait(t *testing.T) {
	config := RateLimitConfig{
		MinDelayMs: 50,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	// Test concurrent access to ensure thread safety
	done := make(chan bool)
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		go func() {
			err := rl.Wait()
			if err != nil {
				t.Errorf("Wait() returned error: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func BenchmarkWait(b *testing.B) {
	config := RateLimitConfig{
		MinDelayMs: 1,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rl.Wait()
	}
}

func BenchmarkShouldRetry(b *testing.B) {
	config := RateLimitConfig{
		MinDelayMs: 1000,
		MaxDelayMs: 5000,
		MaxRetries: 3,
	}

	rl := NewRateLimiter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rl.ShouldRetry(429, 0)
	}
}
