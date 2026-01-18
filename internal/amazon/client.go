package amazon

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/ratelimit"
)

// Client represents an Amazon API client
type Client struct {
	httpClient  *http.Client
	rateLimiter *ratelimit.RateLimiter
	config      *config.Config
	userAgents  []string
	uaIndex     int
}

// NewClient creates a new Amazon client with the given configuration
func NewClient(cfg *config.Config) *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		rateLimiter: ratelimit.NewRateLimiter(cfg.RateLimiting),
		config:      cfg,
		userAgents: []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
		},
		uaIndex: 0,
	}
}

// Do executes an HTTP request with rate limiting and retry logic
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.config.RateLimiting.MaxRetries; attempt++ {
		// Wait for rate limiter before request
		if err := c.rateLimiter.Wait(); err != nil {
			return nil, err
		}

		// Set User-Agent rotation
		req.Header.Set("User-Agent", c.getUserAgent())

		// Set common headers
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("DNT", "1")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		// Execute request
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		// Check if we should retry
		if c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
			resp.Body.Close()
			if err := c.rateLimiter.WaitWithBackoff(attempt); err != nil {
				return nil, err
			}
			continue
		}

		// Success or non-retryable error
		return resp, nil
	}

	return resp, nil
}

// Get is a convenience method for making GET requests
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// getUserAgent returns a user agent string from the rotation list
func (c *Client) getUserAgent() string {
	ua := c.userAgents[c.uaIndex]
	c.uaIndex = (c.uaIndex + 1) % len(c.userAgents)
	return ua
}
