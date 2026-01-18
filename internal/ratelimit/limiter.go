package ratelimit

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// RateLimiter manages request rate limiting
type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          config,
		lastRequestTime: time.Time{},
	}
}

// Wait waits for the minimum delay before allowing the next request
func (r *RateLimiter) Wait() error {
	if r.lastRequestTime.IsZero() {
		r.lastRequestTime = time.Now()
		return nil
	}

	// Calculate time since last request
	elapsed := time.Since(r.lastRequestTime)
	minDelay := time.Duration(r.config.MinDelayMs) * time.Millisecond

	// If not enough time has passed, sleep for the difference
	if elapsed < minDelay {
		time.Sleep(minDelay - elapsed)
	}

	// Add random jitter (0-500ms)
	jitter, err := r.getJitter(500)
	if err != nil {
		return fmt.Errorf("failed to calculate jitter: %w", err)
	}
	time.Sleep(jitter)

	// Update last request time
	r.lastRequestTime = time.Now()

	return nil
}

// WaitWithBackoff waits with exponential backoff
func (r *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := math.Min(math.Pow(2, float64(attempt))*1000, 60000)
	backoffDuration := time.Duration(backoffMs) * time.Millisecond

	time.Sleep(backoffDuration)

	// Update last request time
	r.lastRequestTime = time.Now()

	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt
func (r *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Don't retry if we've exceeded max retries
	if attempt >= r.config.MaxRetries {
		return false
	}

	// Retry on rate limited (429) or service unavailable (503)
	return statusCode == 429 || statusCode == 503
}

// getJitter returns a random duration between 0 and maxMs milliseconds
func (r *RateLimiter) getJitter(maxMs int) (time.Duration, error) {
	if maxMs <= 0 {
		return 0, nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(maxMs)))
	if err != nil {
		return 0, err
	}

	return time.Duration(n.Int64()) * time.Millisecond, nil
}
