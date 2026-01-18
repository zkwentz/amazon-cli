package amazon

import (
	"net/http"
	"time"
)

// Client represents an Amazon API client
type Client struct {
	httpClient *http.Client
	userAgents []string
}

// NewClient creates a new Amazon API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgents: []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
		},
	}
}

// Do executes an HTTP request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Set a user agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.userAgents[0])
	}

	// Set common headers
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	}
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	}

	return c.httpClient.Do(req)
}

// Get performs a GET request
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
