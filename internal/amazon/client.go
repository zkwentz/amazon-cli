package amazon

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
	uaIndex     int
}

// NewClient creates a new Amazon API client
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
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/120.0.0.0",
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
		},
		uaIndex: 0,
	}
}

// Do executes an HTTP request with rate limiting and retry logic
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, err
	}

	// Set user agent (rotate through list)
	req.Header.Set("User-Agent", c.userAgents[c.uaIndex])
	c.uaIndex = (c.uaIndex + 1) % len(c.userAgents)

	// Set common headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	// Execute request with retry logic
	var resp *http.Response
	var err error
	attempt := 0

	for {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		// Check if we should retry
		if c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
			resp.Body.Close()
			c.rateLimiter.WaitWithBackoff(attempt)
			attempt++
			continue
		}

		break
	}

	return resp, nil
}

// Get performs a GET request
func (c *Client) Get(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// PostForm performs a POST request with form data
func (c *Client) PostForm(urlStr string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	return c.Do(req)
}
