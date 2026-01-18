package ratelimit

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

type RateLimiter struct {
	config          config.RateLimitConfig
	lastRequestTime time.Time
}

func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          cfg,
		lastRequestTime: time.Time{},
	}
}

func (rl *RateLimiter) Wait() error {
	if !rl.lastRequestTime.IsZero() {
		elapsed := time.Since(rl.lastRequestTime)
		minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	jitter, err := rand.Int(rand.Reader, big.NewInt(500))
	if err != nil {
		jitter = big.NewInt(0)
	}
	time.Sleep(time.Duration(jitter.Int64()) * time.Millisecond)

	rl.lastRequestTime = time.Now()
	return nil
}

func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	backoffMs := math.Min(math.Pow(2, float64(attempt))*1000, 60000)
	time.Sleep(time.Duration(backoffMs) * time.Millisecond)
	return nil
}

func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	if attempt >= rl.config.MaxRetries {
		return false
	}
	return statusCode == 429 || statusCode == 503
}
