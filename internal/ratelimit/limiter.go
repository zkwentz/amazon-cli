package ratelimit

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// RateLimiter manages request rate limiting with configurable delays and retries
type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
	retryCount      int
}

// NewRateLimiter creates a new RateLimiter with the specified configuration
func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          cfg,
		lastRequestTime: time.Time{}, // Zero value means no previous request
		retryCount:      0,
	}
}

// Wait enforces minimum delay between requests with random jitter
func (rl *RateLimiter) Wait() error {
	// Calculate time since last request
	if !rl.lastRequestTime.IsZero() {
		elapsed := time.Since(rl.lastRequestTime)
		minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

		// If less than minimum delay, sleep for the difference
		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	// Add random jitter (0-500ms)
	jitter, err := rl.randomJitter(500)
	if err != nil {
		// If random generation fails, use a fixed jitter
		jitter = 250 * time.Millisecond
	}
	time.Sleep(jitter)

	// Update last request timestamp
	rl.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff applies exponential backoff for retry attempts
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := 1000 * (1 << uint(attempt)) // 2^attempt * 1000
	if backoffMs > 60000 {
		backoffMs = 60000
	}

	backoff := time.Duration(backoffMs) * time.Millisecond
	time.Sleep(backoff)

	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt count
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Check if we've exceeded max retries
	if attempt >= rl.config.MaxRetries {
		return false
	}

	// Retry on rate limited (429) or service unavailable (503)
	return statusCode == 429 || statusCode == 503
}

// randomJitter generates a random duration between 0 and maxMs milliseconds
func (rl *RateLimiter) randomJitter(maxMs int) (time.Duration, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(maxMs)))
	if err != nil {
		return 0, err
	}
	return time.Duration(n.Int64()) * time.Millisecond, nil
}

// ResetRetryCount resets the retry counter (useful after a successful request)
func (rl *RateLimiter) ResetRetryCount() {
	rl.retryCount = 0
}

// IncrementRetryCount increments the retry counter
func (rl *RateLimiter) IncrementRetryCount() {
	rl.retryCount++
}

// GetRetryCount returns the current retry count
func (rl *RateLimiter) GetRetryCount() int {
	return rl.retryCount
}
