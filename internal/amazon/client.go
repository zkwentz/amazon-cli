package amazon

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/ratelimit"
)

type Client struct {
	httpClient  *http.Client
	rateLimiter *ratelimit.RateLimiter
	config      *config.Config
	userAgents  []string
	currentUA   int
}

func NewClient(cfg *config.Config) *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		rateLimiter: ratelimit.NewRateLimiter(cfg.RateLimiting),
		config:      cfg,
		userAgents: []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		},
		currentUA: 0,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.rateLimiter.Wait()

	req.Header.Set("User-Agent", c.userAgents[c.currentUA%len(c.userAgents)])
	c.currentUA++

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	var resp *http.Response
	var err error
	attempt := 0
	maxRetries := c.config.RateLimiting.MaxRetries

	for attempt < maxRetries {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
			resp.Body.Close()
			c.rateLimiter.WaitWithBackoff(attempt)
			attempt++
			continue
		}

		break
	}

	return resp, nil
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
