package amazon

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// HTTPDoer interface for making HTTP requests (useful for testing)
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents an Amazon API client
type Client struct {
	httpClient HTTPDoer
	userAgents []string
	config     *Config
}

// Config holds client configuration
type Config struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// NewClient creates a new Amazon client
func NewClient(config *Config) *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		userAgents: getUserAgents(),
		config:     config,
	}
}

// Do executes an HTTP request with rate limiting and retries
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Set a random user agent
	req.Header.Set("User-Agent", c.getRandomUserAgent())

	// Set common headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Add authentication if available
	if c.config != nil && c.config.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	}

	// Execute the request with retry logic
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		// Check if we should retry based on status code
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			if backoff > 60*time.Second {
				backoff = 60 * time.Second
			}
			time.Sleep(backoff)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries, lastErr)
}

// getRandomUserAgent returns a random user agent from the list
func (c *Client) getRandomUserAgent() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(c.userAgents))))
	return c.userAgents[n.Int64()]
}

// getUserAgents returns a list of common browser user agents
func getUserAgents() []string {
	return []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	}
}
