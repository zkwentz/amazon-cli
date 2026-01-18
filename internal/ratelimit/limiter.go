package ratelimit

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// RateLimiter handles rate limiting for HTTP requests
type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          cfg,
		lastRequestTime: time.Time{},
	}
}

// Wait waits the appropriate amount of time before allowing the next request
func (rl *RateLimiter) Wait() error {
	if !rl.lastRequestTime.IsZero() {
		elapsed := time.Since(rl.lastRequestTime)
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

	rl.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff waits with exponential backoff based on attempt number
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := math.Min(math.Pow(2, float64(attempt))*1000, 60000)
	time.Sleep(time.Duration(backoffMs) * time.Millisecond)
	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt number
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Don't retry if we've exceeded max retries
	if attempt >= rl.config.MaxRetries {
		return false
	}

	// Retry on rate limited or service unavailable
	return statusCode == 429 || statusCode == 503
}
