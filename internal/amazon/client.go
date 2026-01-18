package amazon

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Client handles all Amazon API/scraping operations
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Amazon API client
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		baseURL: "https://www.amazon.com",
	}
}
