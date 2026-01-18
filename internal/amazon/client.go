package amazon

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MinDelayMs int
	MaxDelayMs int
	MaxRetries int
}

// Config holds the client configuration
type Config struct {
	Auth struct {
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		ExpiresAt    time.Time `json:"expires_at"`
	} `json:"auth"`
	RateLimiting RateLimitConfig `json:"rate_limiting"`
}

// RateLimiter handles rate limiting with jitter and exponential backoff
type RateLimiter struct {
	config          RateLimitConfig
	lastRequestTime time.Time
	mu              sync.Mutex
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:          config,
		lastRequestTime: time.Time{},
	}
}

// Wait enforces minimum delay with jitter between requests
func (rl *RateLimiter) Wait() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if !rl.lastRequestTime.IsZero() {
		elapsed := time.Since(rl.lastRequestTime)
		minDelay := time.Duration(rl.config.MinDelayMs) * time.Millisecond

		if elapsed < minDelay {
			time.Sleep(minDelay - elapsed)
		}
	}

	// Add random jitter (0-500ms)
	jitter, err := rand.Int(rand.Reader, big.NewInt(501))
	if err != nil {
		jitter = big.NewInt(250) // fallback to 250ms on error
	}
	time.Sleep(time.Duration(jitter.Int64()) * time.Millisecond)

	rl.lastRequestTime = time.Now()
	return nil
}

// WaitWithBackoff sleeps for exponential backoff duration
func (rl *RateLimiter) WaitWithBackoff(attempt int) error {
	// Calculate exponential backoff: min(2^attempt * 1000ms, 60000ms)
	backoffMs := math.Min(math.Pow(2, float64(attempt))*1000, 60000)
	duration := time.Duration(backoffMs) * time.Millisecond
	time.Sleep(duration)
	return nil
}

// ShouldRetry determines if a request should be retried based on status code and attempt count
func (rl *RateLimiter) ShouldRetry(statusCode int, attempt int) bool {
	if attempt >= rl.config.MaxRetries {
		return false
	}
	// Retry on rate limited (429) or service unavailable (503)
	return statusCode == 429 || statusCode == 503
}

// Client is the main Amazon API client
type Client struct {
	httpClient  *http.Client
	rateLimiter *RateLimiter
	config      *Config
	userAgents  []string
	uaIndex     int
	mu          sync.Mutex
}

// Common browser User-Agent strings for rotation
var defaultUserAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
}

// NewClient creates a new Amazon API client
func NewClient(config *Config) (*Client, error) {
	// Create cookie jar for session management
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Create HTTP client with 30 second timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	// Initialize rate limiter with config
	rateLimiter := NewRateLimiter(config.RateLimiting)

	return &Client{
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		config:      config,
		userAgents:  defaultUserAgents,
		uaIndex:     0,
	}, nil
}

// getNextUserAgent returns the next User-Agent string in rotation
func (c *Client) getNextUserAgent() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	ua := c.userAgents[c.uaIndex]
	c.uaIndex = (c.uaIndex + 1) % len(c.userAgents)
	return ua
}

// Do executes an HTTP request with rate limiting and retry logic
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Try request with retries
	for attempt := 0; attempt <= c.rateLimiter.config.MaxRetries; attempt++ {
		// Wait for rate limiter before each attempt
		if err := c.rateLimiter.Wait(); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}

		// Set rotating User-Agent
		req.Header.Set("User-Agent", c.getNextUserAgent())

		// Set common headers
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("DNT", "1")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		// Add authorization if access token exists
		if c.config.Auth.AccessToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.Auth.AccessToken))
		}

		// Execute request
		resp, err = c.httpClient.Do(req)
		if err != nil {
			// Network error - retry if we have attempts left
			if attempt < c.rateLimiter.config.MaxRetries {
				if err := c.rateLimiter.WaitWithBackoff(attempt); err != nil {
					return nil, err
				}
				continue
			}
			return nil, fmt.Errorf("request failed after %d attempts: %w", attempt+1, err)
		}

		// Check if we should retry based on status code
		if c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
			resp.Body.Close()
			if err := c.rateLimiter.WaitWithBackoff(attempt); err != nil {
				return nil, err
			}
			continue
		}

		// Request successful (or failed with non-retryable error)
		return resp, nil
	}

	return resp, err
}

// Get is a convenience method for GET requests
func (c *Client) Get(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	return c.Do(req)
}

// PostForm is a convenience method for POST requests with form data
func (c *Client) PostForm(urlStr string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()

	return c.Do(req)
}
