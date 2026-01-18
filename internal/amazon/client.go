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

// Client is the HTTP client for making requests to Amazon
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
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
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
	}

	// Create rate limiter from config
	rateLimiter := ratelimit.NewRateLimiter(
		cfg.RateLimiting.MinDelayMs,
		cfg.RateLimiting.MaxDelayMs,
		cfg.RateLimiting.MaxRetries,
	)

	return &Client{
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		config:      cfg,
		userAgents:  defaultUserAgents,
	}, nil
}

// getRandomUserAgent returns a random user agent from the rotation list
func (c *Client) getRandomUserAgent() string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(c.userAgents))))
	if err != nil {
		// Fallback to first user agent if random fails
		return c.userAgents[0]
	}
	return c.userAgents[n.Int64()]
}

// Do executes an HTTP request with rate limiting and retry logic
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Wait for rate limiter before making request
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Set random User-Agent from rotation list
	req.Header.Set("User-Agent", c.getRandomUserAgent())

	// Set common headers
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	}
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	}
	if req.Header.Get("Accept-Encoding") == "" {
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Execute request with retry logic
	var resp *http.Response
	var err error
	attempt := 0

	for {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("http request failed: %w", err)
		}

		// Check if we should retry based on status code
		if c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
			resp.Body.Close() // Close response body before retry

			// Wait with exponential backoff
			if err := c.rateLimiter.WaitWithBackoff(attempt); err != nil {
				return nil, fmt.Errorf("backoff wait error: %w", err)
			}

			attempt++
			continue
		}

		// Success or non-retryable error - return response
		return resp, nil
	}
}

// Get is a convenience method for GET requests
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	return c.Do(req)
}
