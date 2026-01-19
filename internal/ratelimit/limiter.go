package ratelimit

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// RateLimiter manages rate limiting with exponential backoff and jitter
type RateLimiter struct {
	minDelay   time.Duration
	maxDelay   time.Duration
	maxRetries int
	lastCall   time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new RateLimiter with the specified parameters
func NewRateLimiter(minDelay, maxDelay time.Duration, maxRetries int) *RateLimiter {
	return &RateLimiter{
		minDelay:   minDelay,
		maxDelay:   maxDelay,
		maxRetries: maxRetries,
		lastCall:   time.Time{},
	}
}

// Wait enforces the minimum delay between calls with random jitter (0-500ms)
func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if !rl.lastCall.IsZero() {
		elapsed := time.Since(rl.lastCall)
		remaining := rl.minDelay - elapsed

		if remaining > 0 {
			// Add random jitter between 0-500ms
			jitter := time.Duration(rand.Int63n(501)) * time.Millisecond
			time.Sleep(remaining + jitter)
		} else {
			// Even if minimum delay has passed, add jitter
			jitter := time.Duration(rand.Int63n(501)) * time.Millisecond
			time.Sleep(jitter)
		}
	}

	rl.lastCall = time.Now()
}

// WaitWithBackoff implements exponential backoff with a cap at 60 seconds
func (rl *RateLimiter) WaitWithBackoff(attempt int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if attempt <= 0 {
		attempt = 1
	}

	// Calculate exponential backoff: minDelay * 2^(attempt-1)
	backoff := float64(rl.minDelay) * math.Pow(2, float64(attempt-1))
	delay := time.Duration(backoff)

	// Cap at 60 seconds
	maxBackoff := 60 * time.Second
	if delay > maxBackoff {
		delay = maxBackoff
	}

	// Ensure delay doesn't exceed maxDelay if set
	if rl.maxDelay > 0 && delay > rl.maxDelay {
		delay = rl.maxDelay
	}

	// Add random jitter between 0-500ms
	jitter := time.Duration(rand.Int63n(501)) * time.Millisecond

	time.Sleep(delay + jitter)
	rl.lastCall = time.Now()
}

// ShouldRetry determines if a request should be retried based on status code and attempt count
func (rl *RateLimiter) ShouldRetry(statusCode, attempt int) bool {
	// Check if we've exceeded max retries
	if attempt >= rl.maxRetries {
		return false
	}

	// Retry on rate limit (429) or service unavailable (503)
	return statusCode == 429 || statusCode == 503
}
