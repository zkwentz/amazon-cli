package amazon

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Client is the main Amazon API client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Amazon client
func NewClient() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    "https://www.amazon.com",
	}, nil
}
