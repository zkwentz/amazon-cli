package ratelimit

import (
	"crypto/rand"
	"math/big"
	"time"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

// RateLimiter handles request rate limiting and backoff
type RateLimiter struct {
	config      RateLimitConfig
	lastRequest time.Time
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:      config,
		lastRequest: time.Time{},
	}
}

// Wait ensures the minimum delay between requests and adds jitter
func (rl *RateLimiter) Wait() error {
	if !rl.lastRequest.IsZero() {
		elapsed := time.Since(rl.lastRequest)
		minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	// Add random jitter (0-500ms)
	jitter, err := rand.Int(rand.Reader, big.NewInt(500))
	if err == nil {
		time.Sleep(time.Duration(jitter.Int64()) * time.Millisecond)
	}

	rl.lastRequest = time.Now()
	return nil
}

// WaitWithBackoff sleeps with exponential backoff based on attempt number
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := (1 << uint(attempt)) * 1000
	if backoffMs > 60000 {
		backoffMs = 60000
	}

	time.Sleep(time.Duration(backoffMs) * time.Millisecond)
	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt count
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	if attempt >= rl.config.MaxRetries {
		return false
	}

	// Retry on rate limited or service unavailable
	return statusCode == 429 || statusCode == 503
}
