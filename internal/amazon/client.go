package amazon

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/ratelimit"
)

// Client represents an HTTP client for Amazon API interactions
type Client struct {
	httpClient  *http.Client
	rateLimiter *ratelimit.RateLimiter
	config      *config.Config
	configPath  string
	userAgents  []string
}

// NewClient creates a new Amazon API client
func NewClient(cfg *config.Config, configPath string) (*Client, error) {
	// Create cookie jar for session management
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	rateLimiter := ratelimit.NewRateLimiter(cfg.RateLimiting)

	// Common browser User-Agent strings for rotation
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	}

	return &Client{
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		config:      cfg,
		configPath:  configPath,
		userAgents:  userAgents,
	}, nil
}

// Do executes an HTTP request with rate limiting, auth checks, and retry logic
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Check authentication before making request
	if err := CheckAuth(c.config, c.configPath); err != nil {
		return nil, err
	}

	// Wait for rate limiter
	if err := c.rateLimiter.Wait(); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Set random User-Agent
	req.Header.Set("User-Agent", c.getRandomUserAgent())

	// Set common headers
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	}
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	}

	// Add authorization header if we have an access token
	if c.config.Auth.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.Auth.AccessToken)
	}

	// Execute request with retry logic
	var resp *http.Response
	var err error
	attempt := 0

	for {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		// Check if we should retry
		if c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
			resp.Body.Close()
			attempt++
			if err := c.rateLimiter.WaitWithBackoff(attempt); err != nil {
				return nil, fmt.Errorf("backoff error: %w", err)
			}
			continue
		}

		// Success or non-retryable error
		break
	}

	return resp, nil
}

// Get performs a GET request with auth middleware
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	return c.Do(req)
}

// PostForm performs a POST request with form data and auth middleware
func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	return c.Do(req)
}

// getRandomUserAgent returns a random User-Agent from the list
func (c *Client) getRandomUserAgent() string {
	return c.userAgents[rand.Intn(len(c.userAgents))]
}
