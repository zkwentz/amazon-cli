package ratelimit

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net/http"
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
		lastRequestTime: time.Now().Add(-time.Duration(cfg.MinDelayMs) * time.Millisecond),
	}
}

// Wait waits for the appropriate time before allowing the next request
func (rl *RateLimiter) Wait() error {
	timeSinceLastRequest := time.Since(rl.lastRequestTime)
	minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

	if timeSinceLastRequest < minDelay {
		sleepDuration := minDelay - timeSinceLastRequest
		time.Sleep(sleepDuration)
	}

	// Add random jitter (0-500ms)
	jitter, err := rl.getJitter(500)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(jitter) * time.Millisecond)

	rl.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff waits with exponential backoff based on attempt number
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := int(math.Min(math.Pow(2, float64(attempt))*1000, 60000))

	// Log backoff if verbose (would need logging infrastructure)
	// For now, we'll just sleep
	time.Sleep(time.Duration(backoffMs) * time.Millisecond)

	rl.lastRequestTime = time.Now()
	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt count
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Don't retry if we've exceeded max retries
	if attempt >= rl.config.MaxRetries {
		return false
	}

	// Retry on rate limit (429) or service unavailable (503)
	if statusCode == http.StatusTooManyRequests || statusCode == http.StatusServiceUnavailable {
		return true
	}

	return false
}

// getJitter generates a random jitter value between 0 and maxMs milliseconds
func (rl *RateLimiter) getJitter(maxMs int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(maxMs)))
	if err != nil {
		return 0, fmt.Errorf("failed to generate jitter: %w", err)
	}
	return int(n.Int64()), nil
}
