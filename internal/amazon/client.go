package amazon

import (
	"net/http"
)

// Client represents the Amazon API client
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Amazon client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}
