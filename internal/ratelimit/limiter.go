package ratelimit

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
}

func NewRateLimiter(config config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          config,
		lastRequestTime: time.Time{},
	}
}

func (rl *RateLimiter) Wait() error {
	if rl.lastRequestTime.IsZero() {
		rl.lastRequestTime = time.Now()
		return nil
	}

	elapsed := time.Since(rl.lastRequestTime)
	minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

	if elapsed < minDelay {
		time.Sleep(minDelay - elapsed)
	}

	jitterMax := big.NewInt(500)
	jitter, _ := rand.Int(rand.Reader, jitterMax)
	time.Sleep(time.Duration(jitter.Int64()) * time.Millisecond)

	rl.lastRequestTime = time.Now()
	return nil
}

func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	backoff := time.Duration(1<<uint(attempt)) * time.Second
	maxBackoff := 60 * time.Second

	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	time.Sleep(backoff)
	rl.lastRequestTime = time.Now()
	return nil
}

func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	if attempt >= rl.config.MaxRetries {
		return false
	}

	return statusCode == 429 || statusCode == 503
}
