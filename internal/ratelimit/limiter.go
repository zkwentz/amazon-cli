package ratelimit

import (
	"crypto/rand"
	"math/big"
	"time"
)

// RateLimitConfig holds the rate limiting configuration
type RateLimitConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

// RateLimiter manages rate limiting and retry logic
type RateLimiter struct {
	config          RateLimitConfig
	lastRequestTime time.Time
	retryCount      int
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          config,
		lastRequestTime: time.Time{},
		retryCount:      0,
	}
}

// Wait ensures the minimum delay between requests is respected
func (r *RateLimiter) Wait() error {
	if !r.lastRequestTime.IsZero() {
		elapsed := time.Since(r.lastRequestTime)
		minDelay := time.Duration(r.config.MinDelayMs) * time.Millisecond

		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	// Add random jitter (0-500ms)
	jitter, err := rand.Int(rand.Reader, big.NewInt(500))
	if err != nil {
		// If crypto/rand fails, use a fixed small jitter
		jitter = big.NewInt(100)
	}
	time.Sleep(time.Duration(jitter.Int64()) * time.Millisecond)

	r.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff waits with exponential backoff based on the attempt number
func (r *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := 1000 << attempt // 2^attempt * 1000
	if backoffMs > 60000 {
		backoffMs = 60000
	}

	time.Sleep(time.Duration(backoffMs) * time.Millisecond)
	return nil
}

// ShouldRetry determines whether a request should be retried based on the status code and attempt number
func (r *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Return false if we've exceeded max retries
	if attempt >= r.config.MaxRetries {
		return false
	}

	// Return true for 429 (rate limited) or 503 (service unavailable)
	if statusCode == 429 || statusCode == 503 {
		return true
	}

	// Return false for all other status codes
	return false
}
