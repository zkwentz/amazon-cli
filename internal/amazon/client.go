package amazon

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Client represents an Amazon API client
type Client struct {
	httpClient *http.Client
	userAgents []string
	// config would be added here for authentication, rate limiting, etc.
}

// NewClient creates a new Amazon API client
func NewClient() *Client {
	// Create cookie jar for session management
	jar, _ := cookiejar.New(nil)

	// Create HTTP client with reasonable timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	// List of common browser User-Agent strings for rotation
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	}

	return &Client{
		httpClient: httpClient,
		userAgents: userAgents,
	}
}
