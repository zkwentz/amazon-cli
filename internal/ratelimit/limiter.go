package ratelimit

import (
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
	verbose         bool
}

func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          cfg,
		lastRequestTime: time.Time{},
		verbose:         false,
	}
}

func (rl *RateLimiter) SetVerbose(verbose bool) {
	rl.verbose = verbose
}

func (rl *RateLimiter) Wait() error {
	if !rl.lastRequestTime.IsZero() {
		elapsed := time.Since(rl.lastRequestTime)
		minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

		if elapsed < minDelay {
			sleepDuration := minDelay - elapsed

			// Add random jitter (0-500ms)
			jitter, err := rl.randomJitter(500)
			if err != nil {
				return fmt.Errorf("failed to generate jitter: %w", err)
			}

			sleepDuration += jitter

			if rl.verbose {
				log.Printf("Rate limiter: sleeping for %v (min delay + jitter)", sleepDuration)
			}

			time.Sleep(sleepDuration)
		}
	}

	rl.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff implements exponential backoff for retry scenarios.
// It calculates backoff as min(2^attempt * 1000ms, 60000ms) and sleeps for that duration.
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	if attempt < 0 {
		return fmt.Errorf("attempt must be non-negative, got: %d", attempt)
	}

	// Calculate exponential backoff: 2^attempt * 1000ms, capped at 60000ms (60 seconds)
	backoffMs := math.Pow(2, float64(attempt)) * 1000
	if backoffMs > 60000 {
		backoffMs = 60000
	}

	backoffDuration := time.Duration(backoffMs) * time.Millisecond

	if rl.verbose {
		log.Printf("Rate limiter: exponential backoff for attempt %d, sleeping for %v", attempt, backoffDuration)
	}

	time.Sleep(backoffDuration)
	rl.lastRequestTime = time.Now()

	return nil
}

func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	// Don't retry if we've exceeded max retries
	if attempt >= rl.config.MaxRetries {
		return false
	}

	// Retry on rate limited (429) or service unavailable (503) responses
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
