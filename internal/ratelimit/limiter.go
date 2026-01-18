package ratelimit

import (
	"crypto/rand"
	"math/big"
	"time"
)

// RateLimiter manages rate limiting with exponential backoff
type RateLimiter struct {
	minDelayMs   int
	maxDelayMs   int
	maxRetries   int
	lastRequest  time.Time
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(minDelayMs, maxDelayMs, maxRetries int) *RateLimiter {
	return &RateLimiter{
		minDelayMs:  minDelayMs,
		maxDelayMs:  maxDelayMs,
		maxRetries:  maxRetries,
		lastRequest: time.Time{},
	}
}

// Wait ensures minimum delay between requests with random jitter
func (rl *RateLimiter) Wait() error {
	if !rl.lastRequest.IsZero() {
		elapsed := time.Since(rl.lastRequest)
		minDelay := time.Duration(rl.minDelayMs) * time.Millisecond

		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	// Add random jitter (0-500ms)
	jitter, err := rand.Int(rand.Reader, big.NewInt(501))
	if err != nil {
		// If random fails, use a fixed small jitter
		jitter = big.NewInt(250)
	}
	time.Sleep(time.Duration(jitter.Int64()) * time.Millisecond)

	rl.lastRequest = time.Now()
	return nil
}

// WaitWithBackoff implements exponential backoff for retries
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := 1 << uint(attempt) * 1000 // 2^attempt * 1000
	if backoffMs > 60000 {
		backoffMs = 60000
	}

	time.Sleep(time.Duration(backoffMs) * time.Millisecond)
	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	if attempt >= rl.maxRetries {
		return false
	}

	// Retry on rate limited or service unavailable
	return statusCode == 429 || statusCode == 503
}
