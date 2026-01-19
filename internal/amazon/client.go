package amazon

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// userAgents contains a list of common browser User-Agent strings
// to help mimic real browser behavior when making HTTP requests
var userAgents = []string{
	// Chrome on Windows
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	// Chrome on macOS
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	// Chrome on Android mobile
	"Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	// Firefox on Windows
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	// Firefox on macOS
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
	// Firefox on Android mobile
	"Mozilla/5.0 (Android 13; Mobile; rv:121.0) Gecko/121.0 Firefox/121.0",
	// Safari on macOS
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	// Safari on iPhone
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	// Safari on iPad
	"Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	// Edge on Windows
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
}

var rng *rand.Rand

func init() {
	// Initialize random number generator with current time as seed
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// getRandomUserAgent returns a random User-Agent string from the userAgents slice
func getRandomUserAgent() string {
	return userAgents[rng.Intn(len(userAgents))]
}

// Do executes an HTTP request with rate limiting, retries, and proper headers
// It enforces rate limiting, sets browser-like headers, and automatically retries
// requests that fail with 429 (Too Many Requests) or 503 (Service Unavailable)
// status codes using exponential backoff.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Enforce rate limiting before making the request
	c.rateLimiter.Wait()

	// Set headers to mimic a real browser request
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Execute the initial request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}

	// Check if we should retry based on status code
	attempt := 0
	for c.rateLimiter.ShouldRetry(resp.StatusCode, attempt) {
		// Close the previous response body to avoid resource leaks
		resp.Body.Close()

		// Increment attempt counter and wait with exponential backoff
		attempt++
		c.rateLimiter.WaitWithBackoff(attempt)

		// Set a new random User-Agent for the retry to avoid detection
		req.Header.Set("User-Agent", getRandomUserAgent())

		// Retry the request
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("network request failed: %w", err)
		}
	}

	return resp, nil
}

// detectCAPTCHA checks if the response body contains CAPTCHA indicators
// It looks for common CAPTCHA-related strings and Amazon-specific CAPTCHA patterns
func (c *Client) detectCAPTCHA(body []byte) bool {
	// Convert body to lowercase for case-insensitive matching
	lowerBody := bytes.ToLower(body)

	// Common CAPTCHA indicators
	captchaIndicators := [][]byte{
		[]byte("captcha"),
		[]byte("robot check"),
		[]byte("automated access"),
		[]byte("enter the characters you see"),
		[]byte("type the characters"),
		[]byte("sorry, we just need to make sure you're not a robot"),
		[]byte("to continue shopping, please type the characters"),
		[]byte("api.captcha.com"),
		[]byte("api-secure.recaptcha.net"),
		[]byte("g-recaptcha"),
		[]byte("data-sitekey"),
		[]byte("amazoncaptcha"),
	}

	// Check for any of the CAPTCHA indicators in the response body
	for _, indicator := range captchaIndicators {
		if bytes.Contains(lowerBody, indicator) {
			return true
		}
	}

	return false
}
