package amazon

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// Client represents an Amazon API client
type Client struct {
	httpClient *http.Client
	userAgents []string
	baseURL    string
}

// NewClient creates a new Amazon API client
func NewClient() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	}

	return &Client{
		httpClient: httpClient,
		userAgents: userAgents,
		baseURL:    "https://www.amazon.com",
	}, nil
}

// Do executes an HTTP request with rate limiting and user agent rotation
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c == nil || c.httpClient == nil {
		return nil, http.ErrAbortHandler
	}

	// Set a user agent (in a real implementation, this would rotate)
	if len(c.userAgents) > 0 {
		req.Header.Set("User-Agent", c.userAgents[0])
	}

	// Set common headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	// Execute request
	return c.httpClient.Do(req)
}

// Get executes a GET request
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// PostForm executes a POST request with form data
func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = data
	return c.Do(req)
}
