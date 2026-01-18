package amazon

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/ratelimit"
)

// Client represents the Amazon API client
type Client struct {
	httpClient  *http.Client
	rateLimiter *ratelimit.RateLimiter
	config      *config.Config
	userAgents  []string
}

// Common browser User-Agent strings for rotation
var defaultUserAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
}

// NewClient creates a new Amazon API client
func NewClient(cfg *config.Config) (*Client, error) {
	// Create cookie jar for session management
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Create rate limiter
	rateLimiter := ratelimit.NewRateLimiter(cfg.RateLimit)

	return &Client{
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		config:      cfg,
		userAgents:  defaultUserAgents,
	}, nil
}

// Do executes an HTTP request with rate limiting and retries
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Set a random User-Agent from the rotation list
	userAgent, err := c.getRandomUserAgent()
	if err != nil {
		return nil, fmt.Errorf("failed to get user agent: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	// Set common headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Execute request with retry logic
	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.config.RateLimit.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait with exponential backoff for retries
			if err := c.rateLimiter.WaitWithBackoff(attempt); err != nil {
				return nil, fmt.Errorf("backoff wait failed: %w", err)
			}
		}

		resp, lastErr = c.httpClient.Do(req)
		if lastErr != nil {
			// Network error, retry if within max retries
			if attempt < c.config.RateLimit.MaxRetries {
				continue
			}
			return nil, fmt.Errorf("request failed after %d attempts: %w", attempt+1, lastErr)
		}

		// Check if we should retry based on status code
		if c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
			resp.Body.Close()
			continue
		}

		// Success or non-retryable error
		return resp, nil
	}

	return resp, lastErr
}

// Get performs a GET request
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	return c.Do(req)
}

// PostForm performs a POST request with form data
func (c *Client) PostForm(url string, data map[string]string) (*http.Response, error) {
	form := make(map[string][]string)
	for k, v := range data {
		form[k] = []string{v}
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.Do(req)
}

// getRandomUserAgent returns a random user agent from the list
func (c *Client) getRandomUserAgent() (string, error) {
	if len(c.userAgents) == 0 {
		return "", fmt.Errorf("no user agents available")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(c.userAgents))))
	if err != nil {
		return c.userAgents[0], nil // Fallback to first user agent
	}

	return c.userAgents[n.Int64()], nil
}
