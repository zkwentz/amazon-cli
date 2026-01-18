package amazon

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Client represents an Amazon API client with rate limiting and session management
type Client struct {
	httpClient *http.Client
	userAgents []string
	lastReqTime time.Time
	minDelayMs int
	maxRetries int
}

// Common browser User-Agent strings for rotation
var defaultUserAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
}

// NewClient creates a new Amazon API client
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		userAgents: defaultUserAgents,
		minDelayMs: 1000,
		maxRetries: 3,
	}
}

// Do executes an HTTP request with rate limiting, retries, and user agent rotation
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Apply rate limiting
	c.wait()

	// Set random user agent
	req.Header.Set("User-Agent", c.getRandomUserAgent())

	// Set common headers
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Execute request with retry logic
	var resp *http.Response
	var err error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt < c.maxRetries-1 {
				time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
				continue
			}
			return nil, err
		}

		// If rate limited or service unavailable, retry with backoff
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			if attempt < c.maxRetries-1 {
				backoff := time.Duration(1<<uint(attempt)) * time.Second
				if backoff > 60*time.Second {
					backoff = 60 * time.Second
				}
				time.Sleep(backoff)
				continue
			}
		}

		break
	}

	return resp, err
}

// wait implements rate limiting by ensuring minimum delay between requests
func (c *Client) wait() {
	if c.lastReqTime.IsZero() {
		c.lastReqTime = time.Now()
		return
	}

	elapsed := time.Since(c.lastReqTime)
	minDelay := time.Duration(c.minDelayMs) * time.Millisecond

	if elapsed < minDelay {
		time.Sleep(minDelay - elapsed)
	}

	// Add random jitter (0-500ms)
	jitter := c.getRandomJitter(500)
	time.Sleep(time.Duration(jitter) * time.Millisecond)

	c.lastReqTime = time.Now()
}

// getRandomUserAgent returns a random user agent from the list
func (c *Client) getRandomUserAgent() string {
	if len(c.userAgents) == 0 {
		return defaultUserAgents[0]
	}

	idx := c.getRandomInt(len(c.userAgents))
	return c.userAgents[idx]
}

// getRandomInt returns a random integer between 0 and max (exclusive)
func (c *Client) getRandomInt(max int) int {
	if max <= 0 {
		return 0
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}

	return int(n.Int64())
}

// getRandomJitter returns a random jitter value between 0 and maxMs
func (c *Client) getRandomJitter(maxMs int) int {
	return c.getRandomInt(maxMs)
}
