package ratelimit

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"
)

// RateLimitConfig defines the rate limiting parameters
type RateLimitConfig struct {
	MinDelayMs int `json:"min_delay_ms"`
	MaxDelayMs int `json:"max_delay_ms"`
	MaxRetries int `json:"max_retries"`
}

// RateLimiter implements rate limiting with jitter and exponential backoff
type RateLimiter struct {
	config          RateLimitConfig
	lastRequestTime time.Time
	mu              sync.Mutex
	verbose         bool
}

// NewRateLimiter creates a new RateLimiter with the given configuration
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          config,
		lastRequestTime: time.Time{},
		verbose:         false,
	}
}

// SetVerbose enables or disables verbose logging
func (rl *RateLimiter) SetVerbose(verbose bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.verbose = verbose
}

// Wait enforces rate limiting by sleeping if necessary
// It ensures minimum delay between requests and adds random jitter
func (rl *RateLimiter) Wait() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// If this is the first request, no need to wait
	if rl.lastRequestTime.IsZero() {
		rl.lastRequestTime = time.Now()
		return nil
	}

	// Calculate time since last request
	elapsed := time.Since(rl.lastRequestTime)
	minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

	// If less than minimum delay has passed, sleep for the difference
	if elapsed < minDelay {
		sleepDuration := minDelay - elapsed

		// Add random jitter (0-500ms)
		jitter, err := rl.generateJitter(500)
		if err != nil {
			return fmt.Errorf("failed to generate jitter: %w", err)
		}
		sleepDuration += jitter

		if rl.verbose {
			fmt.Printf("[RateLimiter] Waiting %v (base: %v, jitter: %v)\n",
				sleepDuration, minDelay-elapsed, jitter)
		}

		time.Sleep(sleepDuration)
	} else {
		// Even if enough time has passed, add jitter
		jitter, err := rl.generateJitter(500)
		if err != nil {
			return fmt.Errorf("failed to generate jitter: %w", err)
		}

		if jitter > 0 {
			if rl.verbose {
				fmt.Printf("[RateLimiter] Adding jitter: %v\n", jitter)
			}
			time.Sleep(jitter)
		}
	}

	// Update last request timestamp
	rl.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff implements exponential backoff for retries
// attempt is 0-indexed (0 for first retry, 1 for second retry, etc.)
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Calculate exponential backoff: 2^attempt * 1000ms, capped at 60000ms
	backoffMs := int(math.Pow(2, float64(attempt))) * 1000
	if backoffMs > 60000 {
		backoffMs = 60000
	}

	backoffDuration := time.Duration(backoffMs) * time.Millisecond

	if rl.verbose {
		fmt.Printf("[RateLimiter] Backoff attempt %d: waiting %v\n", attempt, backoffDuration)
	}

	time.Sleep(backoffDuration)

	// Update last request timestamp
	rl.lastRequestTime = time.Now()
	return nil
}

// ShouldRetry determines whether a request should be retried based on status code and attempt number
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Don't retry if we've exceeded max retries
	if attempt >= rl.config.MaxRetries {
		if rl.verbose {
			fmt.Printf("[RateLimiter] Max retries (%d) exceeded\n", rl.config.MaxRetries)
		}
		return false
	}

	// Retry on rate limited (429) or service unavailable (503)
	shouldRetry := statusCode == 429 || statusCode == 503

	if rl.verbose && shouldRetry {
		fmt.Printf("[RateLimiter] Status %d: will retry (attempt %d/%d)\n",
			statusCode, attempt+1, rl.config.MaxRetries)
	}

	return shouldRetry
}

// generateJitter generates a random duration between 0 and maxMs milliseconds
// Uses crypto/rand for cryptographically secure randomness
func (rl *RateLimiter) generateJitter(maxMs int) (time.Duration, error) {
	if maxMs <= 0 {
		return 0, nil
	}

	// Generate random number between 0 and maxMs
	maxBig := big.NewInt(int64(maxMs))
	n, err := rand.Int(rand.Reader, maxBig)
	if err != nil {
		return 0, err
	}

	return time.Duration(n.Int64()) * time.Millisecond, nil
}

// Reset resets the rate limiter state (useful for testing)
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.lastRequestTime = time.Time{}
}

// GetConfig returns the current rate limit configuration
func (rl *RateLimiter) GetConfig() RateLimitConfig {
	return rl.config
}
