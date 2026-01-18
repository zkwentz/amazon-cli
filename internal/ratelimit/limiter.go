package ratelimit

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// RateLimiter manages request rate limiting with delays and backoff
type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
	retryCount      int
	mu              sync.Mutex
}

// NewRateLimiter creates a new RateLimiter with the given configuration
func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          cfg,
		lastRequestTime: time.Time{},
		retryCount:      0,
	}
}

// Wait implements rate limiting by ensuring minimum delay between requests
// with random jitter to avoid predictable patterns
func (rl *RateLimiter) Wait() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Calculate time since last request
	now := time.Now()
	if !rl.lastRequestTime.IsZero() {
		elapsed := now.Sub(rl.lastRequestTime)
		minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

		// If less than minimum delay has passed, sleep for the difference
		if elapsed < minDelay {
			sleepDuration := minDelay - elapsed
			time.Sleep(sleepDuration)
		}
	}

	// Add random jitter (0-500ms) using crypto/rand for better randomness
	jitter, err := rl.generateJitter(500)
	if err != nil {
		// Fall back to no jitter if random generation fails
		jitter = 0
	}

	if jitter > 0 {
		time.Sleep(time.Duration(jitter) * time.Millisecond)
	}

	// Update last request timestamp
	rl.lastRequestTime = time.Now()

	return nil
}

// generateJitter generates a random jitter value between 0 and maxMs milliseconds
func (rl *RateLimiter) generateJitter(maxMs int64) (int64, error) {
	if maxMs <= 0 {
		return 0, nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(maxMs+1))
	if err != nil {
		return 0, err
	}

	return n.Int64(), nil
}

// WaitWithBackoff implements exponential backoff for retry scenarios
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := int64(1000) << attempt // 2^attempt * 1000
	maxBackoffMs := int64(60000)

	if backoffMs > maxBackoffMs {
		backoffMs = maxBackoffMs
	}

	// Sleep for calculated duration
	time.Sleep(time.Duration(backoffMs) * time.Millisecond)

	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt count
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Don't retry if we've exceeded max retries
	if attempt >= rl.config.MaxRetries {
		return false
	}

	// Retry for rate limited (429) or service unavailable (503)
	return statusCode == 429 || statusCode == 503
}

// ResetRetryCount resets the retry counter (call after successful request)
func (rl *RateLimiter) ResetRetryCount() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.retryCount = 0
}
