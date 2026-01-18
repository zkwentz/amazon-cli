package amazon

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// HTTPClient is an interface for HTTP operations
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

// Client is the Amazon API client
type Client struct {
	httpClient HTTPClient
	config     *config.Config
	userAgents []string
}

// defaultHTTPClient implements HTTPClient using standard http.Client
type defaultHTTPClient struct {
	client     *http.Client
	userAgents []string
}

func (d *defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Set headers
	req.Header.Set("User-Agent", d.getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	return d.client.Do(req)
}

func (d *defaultHTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return d.Do(req)
}

func (d *defaultHTTPClient) getRandomUserAgent() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(d.userAgents))))
	return d.userAgents[n.Int64()]
}

// NewClient creates a new Amazon API client
func NewClient(cfg *config.Config) *Client {
	jar, _ := cookiejar.New(nil)

	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
	}

	return &Client{
		httpClient: &defaultHTTPClient{
			client: &http.Client{
				Timeout: 30 * time.Second,
				Jar:     jar,
			},
			userAgents: userAgents,
		},
		config:     cfg,
		userAgents: userAgents,
	}
}

// NewClientWithHTTPClient creates a client with a custom HTTP client (useful for testing)
func NewClientWithHTTPClient(cfg *config.Config, httpClient HTTPClient) *Client {
	return &Client{
		httpClient: httpClient,
		config:     cfg,
	}
}
