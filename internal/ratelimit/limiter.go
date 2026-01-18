package ratelimit

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// RateLimiter manages rate limiting for API requests
type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
	mu              sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          cfg,
		lastRequestTime: time.Time{},
	}
}

// Wait ensures minimum delay between requests with jitter
func (rl *RateLimiter) Wait() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Calculate time since last request
	if !rl.lastRequestTime.IsZero() {
		elapsed := time.Since(rl.lastRequestTime)
		minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

		// If less than minimum delay has passed, sleep for the difference
		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	// Add random jitter (0-500ms)
	jitter, err := rl.generateJitter(500)
	if err != nil {
		return fmt.Errorf("failed to generate jitter: %w", err)
	}
	time.Sleep(jitter)

	// Update last request timestamp
	rl.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff implements exponential backoff for retry attempts
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := 1 << uint(attempt) * 1000 // 2^attempt * 1000
	if backoffMs > 60000 {
		backoffMs = 60000
	}

	backoffDuration := time.Duration(backoffMs) * time.Millisecond
	time.Sleep(backoffDuration)

	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt count
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Don't retry if max retries exceeded
	if attempt >= rl.config.MaxRetries {
		return false
	}

	// Retry on rate limit (429) or service unavailable (503)
	return statusCode == 429 || statusCode == 503
}

// generateJitter creates a random duration between 0 and maxMs milliseconds
func (rl *RateLimiter) generateJitter(maxMs int) (time.Duration, error) {
	if maxMs <= 0 {
		return 0, nil
	}

	// Generate cryptographically secure random number
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(maxMs)))
	if err != nil {
		return 0, err
	}

	return time.Duration(nBig.Int64()) * time.Millisecond, nil
}

// Reset resets the rate limiter state (useful for testing)
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.lastRequestTime = time.Time{}
}

// GetLastRequestTime returns the last request time (useful for testing)
func (rl *RateLimiter) GetLastRequestTime() time.Time {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.lastRequestTime
}
