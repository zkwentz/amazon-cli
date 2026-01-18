package amazon

import (
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
	userAgentIdx int
}

// NewClient creates a new Amazon API client
func NewClient(cfg *config.Config) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	limiter := ratelimit.NewRateLimiter(cfg.RateLimiting)

	// Common browser user agents for rotation
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	return &Client{
		httpClient:  httpClient,
		rateLimiter: limiter,
		config:      cfg,
		userAgents:  userAgents,
		userAgentIdx: 0,
	}, nil
}

// Do executes an HTTP request with rate limiting and retries
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Set user agent
	req.Header.Set("User-Agent", c.userAgents[c.userAgentIdx])
	c.userAgentIdx = (c.userAgentIdx + 1) % len(c.userAgents)

	// Set common headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Execute request with retries
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.config.RateLimiting.MaxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt < c.config.RateLimiting.MaxRetries {
				c.rateLimiter.WaitWithBackoff(attempt)
				continue
			}
			return nil, err
		}

		// Check if we should retry
		if c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
			resp.Body.Close()
			c.rateLimiter.WaitWithBackoff(attempt)
			continue
		}

		return resp, nil
	}

	return resp, err
}

// Get performs a GET request
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
