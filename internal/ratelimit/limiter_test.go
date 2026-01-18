package ratelimit

import (
	"sync"
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
		t.Errorf("expected MinDelayMs=1000, got %d", limiter.config.MinDelayMs)
	}

	if limiter.config.MaxDelayMs != 5000 {
		t.Errorf("expected MaxDelayMs=5000, got %d", limiter.config.MaxDelayMs)
	}

	if limiter.config.MaxRetries != 3 {
		t.Errorf("expected MaxRetries=3, got %d", limiter.config.MaxRetries)
	}

	if !limiter.lastRequestTime.IsZero() {
		t.Error("expected lastRequestTime to be zero initially")
	}
}

func TestWait_FirstCall(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
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

	// First call should only add jitter (0-500ms), no minimum delay wait
	// Allow some tolerance for execution time
	if elapsed > 600*time.Millisecond {
		t.Errorf("first Wait() took too long: %v (expected <600ms)", elapsed)
	}

	if limiter.lastRequestTime.IsZero() {
		t.Error("lastRequestTime should be set after Wait()")
	}
}

func TestWait_EnforcesMinimumDelay(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 200,
		MaxDelayMs: 500,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// First call
	err := limiter.Wait()
	if err != nil {
		t.Fatalf("first Wait() returned error: %v", err)
	}

	// Second call immediately after
	start := time.Now()
	err = limiter.Wait()
	if err != nil {
		t.Fatalf("second Wait() returned error: %v", err)
	}
	elapsed := time.Since(start)

	// Should wait at least MinDelayMs (200ms) plus jitter (0-500ms)
	minExpected := 200 * time.Millisecond
	maxExpected := 800 * time.Millisecond // 200ms + 500ms jitter + tolerance

	if elapsed < minExpected {
		t.Errorf("Wait() took %v, expected at least %v", elapsed, minExpected)
	}

	if elapsed > maxExpected {
		t.Errorf("Wait() took %v, expected at most %v", elapsed, maxExpected)
	}
}

func TestWait_SuccessiveCallsTiming(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Make three successive calls
	for i := 0; i < 3; i++ {
		start := time.Now()
		err := limiter.Wait()
		if err != nil {
			t.Fatalf("Wait() call %d returned error: %v", i, err)
		}

		if i > 0 {
			elapsed := time.Since(start)
			// Should enforce minimum delay
			if elapsed < 100*time.Millisecond {
				t.Errorf("Wait() call %d took %v, expected at least 100ms", i, elapsed)
			}
		}
	}
}

func TestWait_WithDelayBetweenCalls(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// First call
	err := limiter.Wait()
	if err != nil {
		t.Fatalf("first Wait() returned error: %v", err)
	}

	// Wait longer than minimum delay
	time.Sleep(150 * time.Millisecond)

	// Second call should only add jitter since we already waited
	start := time.Now()
	err = limiter.Wait()
	if err != nil {
		t.Fatalf("second Wait() returned error: %v", err)
	}
	elapsed := time.Since(start)

	// Should only add jitter, not full minimum delay
	// Jitter is 0-500ms, give some tolerance
	if elapsed > 600*time.Millisecond {
		t.Errorf("Wait() took %v, expected mostly just jitter (<600ms)", elapsed)
	}
}

func TestWait_ThreadSafety(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 50,
		MaxDelayMs: 200,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Run 10 concurrent Wait() calls
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := limiter.Wait(); err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("concurrent Wait() returned error: %v", err)
	}
}

func TestWaitWithBackoff_ExponentialGrowth(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	testCases := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},    // 2^0 * 1000ms = 1s
		{1, 2 * time.Second},    // 2^1 * 1000ms = 2s
		{2, 4 * time.Second},    // 2^2 * 1000ms = 4s
		{3, 8 * time.Second},    // 2^3 * 1000ms = 8s
		{4, 16 * time.Second},   // 2^4 * 1000ms = 16s
		{5, 32 * time.Second},   // 2^5 * 1000ms = 32s
		{6, 60 * time.Second},   // 2^6 * 1000ms = 64s, capped at 60s
		{10, 60 * time.Second},  // way beyond, should cap at 60s
	}

	for _, tc := range testCases {
		start := time.Now()
		err := limiter.WaitWithBackoff(tc.attempt)
		if err != nil {
			t.Fatalf("WaitWithBackoff(%d) returned error: %v", tc.attempt, err)
		}
		elapsed := time.Since(start)

		// Allow 50ms tolerance for execution time
		tolerance := 50 * time.Millisecond
		if elapsed < tc.expected-tolerance || elapsed > tc.expected+tolerance {
			t.Errorf("WaitWithBackoff(%d): expected ~%v, got %v", tc.attempt, tc.expected, elapsed)
		}
	}
}

func TestWaitWithBackoff_MaximumCap(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Test that backoff is capped at 60 seconds
	start := time.Now()
	err := limiter.WaitWithBackoff(10) // Large attempt number
	if err != nil {
		t.Fatalf("WaitWithBackoff(10) returned error: %v", err)
	}
	elapsed := time.Since(start)

	expected := 60 * time.Second
	tolerance := 50 * time.Millisecond

	if elapsed < expected-tolerance || elapsed > expected+tolerance {
		t.Errorf("WaitWithBackoff(10): expected ~%v (capped), got %v", expected, elapsed)
	}
}

func TestShouldRetry_StatusCode429(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	testCases := []struct {
		attempt  int
		expected bool
	}{
		{0, true},  // First retry
		{1, true},  // Second retry
		{2, true},  // Third retry
		{3, false}, // Exceeded max retries
		{4, false}, // Way beyond max retries
	}

	for _, tc := range testCases {
		result := limiter.ShouldRetry(429, tc.attempt)
		if result != tc.expected {
			t.Errorf("ShouldRetry(429, %d): expected %v, got %v", tc.attempt, tc.expected, result)
		}
	}
}

func TestShouldRetry_StatusCode503(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	testCases := []struct {
		attempt  int
		expected bool
	}{
		{0, true},  // First retry
		{1, true},  // Second retry
		{2, true},  // Third retry
		{3, false}, // Exceeded max retries
	}

	for _, tc := range testCases {
		result := limiter.ShouldRetry(503, tc.attempt)
		if result != tc.expected {
			t.Errorf("ShouldRetry(503, %d): expected %v, got %v", tc.attempt, tc.expected, result)
		}
	}
}

func TestShouldRetry_OtherStatusCodes(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// These status codes should never trigger retry
	statusCodes := []int{200, 201, 400, 401, 403, 404, 500, 502}

	for _, code := range statusCodes {
		result := limiter.ShouldRetry(code, 0)
		if result {
			t.Errorf("ShouldRetry(%d, 0): expected false, got true", code)
		}
	}
}

func TestShouldRetry_CustomMaxRetries(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 5, // Custom max retries
	}

	limiter := NewRateLimiter(cfg)

	testCases := []struct {
		attempt  int
		expected bool
	}{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, false}, // Exceeded custom max retries
		{6, false},
	}

	for _, tc := range testCases {
		result := limiter.ShouldRetry(429, tc.attempt)
		if result != tc.expected {
			t.Errorf("ShouldRetry(429, %d) with MaxRetries=5: expected %v, got %v", tc.attempt, tc.expected, result)
		}
	}
}

func TestGenerateJitter(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Test multiple times to check randomness
	for i := 0; i < 100; i++ {
		jitter, err := limiter.generateJitter(500)
		if err != nil {
			t.Fatalf("generateJitter(500) returned error: %v", err)
		}

		if jitter < 0 {
			t.Errorf("generateJitter returned negative value: %v", jitter)
		}

		if jitter > 500*time.Millisecond {
			t.Errorf("generateJitter returned value > 500ms: %v", jitter)
		}
	}
}

func TestGenerateJitter_ZeroMax(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	jitter, err := limiter.generateJitter(0)
	if err != nil {
		t.Fatalf("generateJitter(0) returned error: %v", err)
	}

	if jitter != 0 {
		t.Errorf("generateJitter(0) expected 0, got %v", jitter)
	}
}

func TestGenerateJitter_NegativeMax(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	jitter, err := limiter.generateJitter(-10)
	if err != nil {
		t.Fatalf("generateJitter(-10) returned error: %v", err)
	}

	if jitter != 0 {
		t.Errorf("generateJitter(-10) expected 0, got %v", jitter)
	}
}

func TestReset(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Make a call to set lastRequestTime
	err := limiter.Wait()
	if err != nil {
		t.Fatalf("Wait() returned error: %v", err)
	}

	if limiter.GetLastRequestTime().IsZero() {
		t.Error("lastRequestTime should be set after Wait()")
	}

	// Reset the limiter
	limiter.Reset()

	if !limiter.GetLastRequestTime().IsZero() {
		t.Error("lastRequestTime should be zero after Reset()")
	}
}

func TestGetLastRequestTime(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// Initially should be zero
	if !limiter.GetLastRequestTime().IsZero() {
		t.Error("GetLastRequestTime() should return zero time initially")
	}

	// After Wait(), should be set
	before := time.Now()
	err := limiter.Wait()
	if err != nil {
		t.Fatalf("Wait() returned error: %v", err)
	}
	after := time.Now()

	lastReq := limiter.GetLastRequestTime()
	if lastReq.Before(before) || lastReq.After(after) {
		t.Errorf("GetLastRequestTime() returned %v, expected between %v and %v", lastReq, before, after)
	}
}

func TestWait_ZeroMinDelay(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 0,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// First call
	err := limiter.Wait()
	if err != nil {
		t.Fatalf("first Wait() returned error: %v", err)
	}

	// Second call immediately - should only add jitter
	start := time.Now()
	err = limiter.Wait()
	if err != nil {
		t.Fatalf("second Wait() returned error: %v", err)
	}
	elapsed := time.Since(start)

	// With zero min delay, should only have jitter (0-500ms)
	if elapsed > 600*time.Millisecond {
		t.Errorf("Wait() with MinDelayMs=0 took %v, expected <600ms (jitter only)", elapsed)
	}
}

func TestRateLimiter_IntegrationScenario(t *testing.T) {
	// Simulate a realistic scenario with retries
	cfg := config.RateLimitConfig{
		MinDelayMs: 100,
		MaxDelayMs: 300,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)

	// First request
	err := limiter.Wait()
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	// Simulate getting a 429 response and retrying
	statusCode := 429
	attempt := 0

	for limiter.ShouldRetry(statusCode, attempt) {
		err := limiter.WaitWithBackoff(attempt)
		if err != nil {
			t.Fatalf("backoff wait failed: %v", err)
		}
		attempt++

		// Simulate success on third retry
		if attempt >= 2 {
			statusCode = 200
			break
		}
	}

	if statusCode != 200 {
		t.Error("integration scenario: expected eventual success")
	}

	if attempt != 2 {
		t.Errorf("integration scenario: expected 2 retries, got %d", attempt)
	}
}

func TestRateLimiter_ConcurrentRequests(t *testing.T) {
	cfg := config.RateLimitConfig{
		MinDelayMs: 50,
		MaxDelayMs: 200,
		MaxRetries: 3,
	}

	limiter := NewRateLimiter(cfg)
	var wg sync.WaitGroup
	var mu sync.Mutex
	requestTimes := make([]time.Time, 0, 5)

	// Make 5 concurrent requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			limiter.Wait()
			now := time.Now()
			mu.Lock()
			requestTimes = append(requestTimes, now)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Verify all requests completed
	if len(requestTimes) != 5 {
		t.Fatalf("expected 5 request times, got %d", len(requestTimes))
	}

	// Verify that requests were properly rate limited
	// The total duration should account for rate limiting
	minTime := requestTimes[0]
	maxTime := requestTimes[0]
	for _, rt := range requestTimes {
		if rt.Before(minTime) {
			minTime = rt
		}
		if rt.After(maxTime) {
			maxTime = rt
		}
	}

	totalDuration := maxTime.Sub(minTime)
	// With 5 requests and min delay of 50ms, should take at least 200ms total
	// (first request + 4 more with delays)
	// Add tolerance for jitter and execution time
	minExpected := 200 * time.Millisecond
	if totalDuration < minExpected {
		t.Errorf("concurrent requests completed too quickly: %v (expected at least %v)", totalDuration, minExpected)
	}
}
