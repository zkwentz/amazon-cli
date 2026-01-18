package ratelimit

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// RateLimiter handles rate limiting for API requests
type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          cfg,
		lastRequestTime: time.Time{},
	}
}

// Wait waits according to the rate limiting rules
func (r *RateLimiter) Wait() error {
	// Calculate time since last request
	if !r.lastRequestTime.IsZero() {
		elapsed := time.Since(r.lastRequestTime)
		minDelay := time.Duration(r.config.MinDelayMs) * time.Millisecond

		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	// Add jitter (0-500ms)
	jitter, err := rand.Int(rand.Reader, big.NewInt(500))
	if err == nil {
		time.Sleep(time.Duration(jitter.Int64()) * time.Millisecond)
	}

	r.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff waits with exponential backoff
func (r *RateLimiter) WaitWithBackoff(attempt int) {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := math.Min(math.Pow(2, float64(attempt))*1000, 60000)
	time.Sleep(time.Duration(backoffMs) * time.Millisecond)
}

// ShouldRetry determines if a request should be retried
func (r *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	if attempt >= r.config.MaxRetries {
		return false
	}

	// Retry on rate limited or service unavailable
	return statusCode == 429 || statusCode == 503
}
