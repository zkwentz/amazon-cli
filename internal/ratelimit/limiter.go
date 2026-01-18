package ratelimit

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// RateLimiter manages rate limiting for HTTP requests
type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
	retryCount      int
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          cfg,
		lastRequestTime: time.Time{},
		retryCount:      0,
	}
}

// Wait blocks until the minimum delay has elapsed since the last request
func (rl *RateLimiter) Wait() error {
	// Calculate time since last request
	if !rl.lastRequestTime.IsZero() {
		elapsed := time.Since(rl.lastRequestTime)
		minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

		// If less than MinDelayMs, sleep for the difference
		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	// Add random jitter (0-500ms)
	jitter, err := rand.Int(rand.Reader, big.NewInt(500))
	if err != nil {
		jitter = big.NewInt(0)
	}
	time.Sleep(time.Duration(jitter.Int64()) * time.Millisecond)

	// Update last request timestamp
	rl.lastRequestTime = time.Now()

	return nil
}

// WaitWithBackoff waits with exponential backoff based on attempt number
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := 1 << uint(attempt) * 1000
	maxBackoffMs := 60000
	if backoffMs > maxBackoffMs {
		backoffMs = maxBackoffMs
	}

	duration := time.Duration(backoffMs) * time.Millisecond
	time.Sleep(duration)

	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt count
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Return false if attempt >= MaxRetries
	if attempt >= rl.config.MaxRetries {
		return false
	}

	// Return true for 429 (rate limited) or 503 (service unavailable)
	return statusCode == 429 || statusCode == 503
}
