package amazon

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Client represents an Amazon API client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Amazon client
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    "https://www.amazon.com",
	}
}
