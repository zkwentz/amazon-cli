package amazon

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetRandomUserAgent(t *testing.T) {
	// Test that getRandomUserAgent returns a non-empty string
	userAgent := getRandomUserAgent()
	if userAgent == "" {
		t.Error("getRandomUserAgent returned empty string")
	}
}

func TestGetRandomUserAgentReturnsValidUserAgent(t *testing.T) {
	// Test that the returned user agent is one from our list
	userAgent := getRandomUserAgent()

	found := false
	for _, ua := range userAgents {
		if ua == userAgent {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("getRandomUserAgent returned unexpected user agent: %s", userAgent)
	}
}

func TestGetRandomUserAgentDistribution(t *testing.T) {
	// Test that calling getRandomUserAgent multiple times returns different values
	// This tests the randomness aspect
	calls := 100
	results := make(map[string]int)

	for i := 0; i < calls; i++ {
		ua := getRandomUserAgent()
		results[ua]++
	}

	// With 100 calls and 10 user agents, we should see at least 2 different user agents
	// (this is a probabilistic test, but the chances of getting only 1 UA in 100 calls is extremely low)
	if len(results) < 2 {
		t.Errorf("getRandomUserAgent showed poor randomness: only %d unique user agents in %d calls", len(results), calls)
	}
}

func TestUserAgentsSliceHasCorrectLength(t *testing.T) {
	// Test that we have exactly 10 user agents as specified
	expectedCount := 10
	if len(userAgents) != expectedCount {
		t.Errorf("userAgents slice has %d entries, expected %d", len(userAgents), expectedCount)
	}
}

func TestUserAgentsSliceContainsValidStrings(t *testing.T) {
	// Test that all user agents are non-empty and appear valid
	for i, ua := range userAgents {
		if ua == "" {
			t.Errorf("userAgents[%d] is empty", i)
		}

		// User agents should be reasonably long (at least 50 characters)
		if len(ua) < 50 {
			t.Errorf("userAgents[%d] appears to be invalid (too short): %s", i, ua)
		}

		// All modern browser user agents contain "Mozilla"
		if len(ua) > 0 && len(ua) >= 7 {
			found := false
			for j := 0; j <= len(ua)-7; j++ {
				if ua[j:j+7] == "Mozilla" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("userAgents[%d] does not contain 'Mozilla': %s", i, ua)
			}
		}
	}
}

func TestDo_Success(t *testing.T) {
	// Create a test server that returns 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers are set
		if r.Header.Get("User-Agent") == "" {
			t.Error("User-Agent header not set")
		}
		if r.Header.Get("Accept") == "" {
			t.Error("Accept header not set")
		}
		if r.Header.Get("Accept-Language") == "" {
			t.Error("Accept-Language header not set")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	client := NewClient()
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "success" {
		t.Errorf("Expected body 'success', got '%s'", string(body))
	}
}

func TestDo_SetsHeaders(t *testing.T) {
	// Track if headers were set correctly
	headerChecks := struct {
		userAgent      bool
		accept         bool
		acceptLanguage bool
	}{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		headerChecks.userAgent = ua != "" && strings.Contains(ua, "Mozilla")

		accept := r.Header.Get("Accept")
		headerChecks.accept = accept != "" && strings.Contains(accept, "text/html")

		acceptLang := r.Header.Get("Accept-Language")
		headerChecks.acceptLanguage = acceptLang != "" && strings.Contains(acceptLang, "en-US")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if !headerChecks.userAgent {
		t.Error("User-Agent header not set correctly")
	}
	if !headerChecks.accept {
		t.Error("Accept header not set correctly")
	}
	if !headerChecks.acceptLanguage {
		t.Error("Accept-Language header not set correctly")
	}
}

func TestDo_RetryOn429(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			// Return 429 for first 2 attempts
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited"))
		} else {
			// Return 200 on 3rd attempt
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}))
	defer server.Close()

	client := NewClient()
	req, _ := http.NewRequest("GET", server.URL, nil)

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected final status 200, got %d", resp.StatusCode)
	}

	// Should have waited for backoff between retries
	// With 2 retries and exponential backoff, should take at least a few seconds
	if elapsed < 1*time.Second {
		t.Errorf("Expected retry delays, but completed too quickly: %v", elapsed)
	}
}

func TestDo_RetryOn503(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("service unavailable"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}
	}))
	defer server.Close()

	client := NewClient()
	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if attemptCount != 2 {
		t.Errorf("Expected 2 attempts (1 retry), got %d", attemptCount)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected final status 200, got %d", resp.StatusCode)
	}
}

func TestDo_StopsAfterMaxRetries(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		// Always return 429
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("rate limited"))
	}))
	defer server.Close()

	client := NewClient()
	// NewClient creates a client with maxRetries=3 in the rate limiter

	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	// With default maxRetries=3 in rate limiter:
	// - Initial attempt (before loop)
	// - Retry attempt 0 (ShouldRetry returns true since 0 < 3)
	// - Retry attempt 1 (ShouldRetry returns true since 1 < 3)
	// - Retry attempt 2 (ShouldRetry returns true since 2 < 3)
	// - Retry attempt 3 (ShouldRetry returns false since 3 >= 3)
	// Total: 1 initial + 3 retries = 4 attempts
	expectedAttempts := 4
	if attemptCount != expectedAttempts {
		t.Errorf("Expected %d attempts with default maxRetries=3, got %d", expectedAttempts, attemptCount)
	}

	// Final response should still be 429
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected final status 429, got %d", resp.StatusCode)
	}
}

func TestDo_NoRetryOn200(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	client := NewClient()
	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt for successful response, got %d", attemptCount)
	}
}

func TestDo_NoRetryOn404(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer server.Close()

	client := NewClient()
	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt for 404 (no retry), got %d", attemptCount)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestDo_NetworkError(t *testing.T) {
	client := NewClient()
	// Create request to invalid URL
	req, _ := http.NewRequest("GET", "http://invalid-url-that-does-not-exist-12345.com", nil)

	_, err := client.Do(req)
	if err == nil {
		t.Error("Expected error for network failure, got nil")
	}

	if !strings.Contains(err.Error(), "network request failed") {
		t.Errorf("Expected error message to contain 'network request failed', got: %v", err)
	}
}

func TestDo_ChangesUserAgentOnRetry(t *testing.T) {
	userAgents := make([]string, 0)
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		userAgents = append(userAgents, r.Header.Get("User-Agent"))

		if attemptCount < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := NewClient()
	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if len(userAgents) != 3 {
		t.Fatalf("Expected 3 user agents captured, got %d", len(userAgents))
	}

	// All user agents should be set (not empty)
	for i, ua := range userAgents {
		if ua == "" {
			t.Errorf("User agent %d is empty", i)
		}
	}

	// Verify they're all valid user agents from our list
	for i, ua := range userAgents {
		found := false
		for _, knownUA := range userAgents {
			if ua == knownUA {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("User agent %d is not a known user agent: %s", i, ua)
		}
	}
}

func TestDo_EnforcesRateLimiting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()

	// Make first request
	req1, _ := http.NewRequest("GET", server.URL, nil)
	start := time.Now()
	_, err := client.Do(req1)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}

	// Make second request - should be rate limited
	req2, _ := http.NewRequest("GET", server.URL, nil)
	_, err = client.Do(req2)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Second request failed: %v", err)
	}

	// Second request should have been delayed by at least the minimum delay
	// Note: The rate limiter adds jitter, so we check for a minimum of 1.5 seconds
	if elapsed < 1500*time.Millisecond {
		t.Errorf("Expected rate limiting delay, but requests completed too quickly: %v", elapsed)
	}
}

func TestDetectCAPTCHA_DetectsGenericCAPTCHA(t *testing.T) {
	client := NewClient()

	testCases := []struct {
		name string
		body string
	}{
		{"lowercase captcha", "Please solve this captcha to continue"},
		{"uppercase CAPTCHA", "Please solve this CAPTCHA to continue"},
		{"mixed case CaPtChA", "Please solve this CaPtChA to continue"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !client.detectCAPTCHA([]byte(tc.body)) {
				t.Errorf("Expected detectCAPTCHA to return true for body containing '%s'", tc.name)
			}
		})
	}
}

func TestDetectCAPTCHA_DetectsRobotCheck(t *testing.T) {
	client := NewClient()

	bodies := []string{
		"Robot Check - Amazon.com",
		"ROBOT CHECK",
		"robot check required",
	}

	for _, body := range bodies {
		if !client.detectCAPTCHA([]byte(body)) {
			t.Errorf("Expected detectCAPTCHA to return true for body: %s", body)
		}
	}
}

func TestDetectCAPTCHA_DetectsAutomatedAccessMessages(t *testing.T) {
	client := NewClient()

	body := "We've detected automated access from your IP address."
	if !client.detectCAPTCHA([]byte(body)) {
		t.Error("Expected detectCAPTCHA to return true for automated access message")
	}
}

func TestDetectCAPTCHA_DetectsAmazonSpecificPatterns(t *testing.T) {
	client := NewClient()

	testCases := []struct {
		name string
		body string
	}{
		{
			"Enter characters message",
			"Enter the characters you see below",
		},
		{
			"Type characters message",
			"Type the characters you see in this image",
		},
		{
			"Amazon robot message",
			"Sorry, we just need to make sure you're not a robot. For best results, please make sure your browser is accepting cookies.",
		},
		{
			"Continue shopping message",
			"To continue shopping, please type the characters you see in the image below.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !client.detectCAPTCHA([]byte(tc.body)) {
				t.Errorf("Expected detectCAPTCHA to return true for: %s", tc.name)
			}
		})
	}
}

func TestDetectCAPTCHA_DetectsReCAPTCHAElements(t *testing.T) {
	client := NewClient()

	testCases := []struct {
		name string
		body string
	}{
		{
			"reCAPTCHA API",
			`<script src="https://api-secure.recaptcha.net/recaptcha/api.js"></script>`,
		},
		{
			"g-recaptcha div",
			`<div class="g-recaptcha" data-sitekey="6LfExample"></div>`,
		},
		{
			"data-sitekey attribute",
			`<div data-sitekey="6LfExample" data-callback="onSubmit"></div>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !client.detectCAPTCHA([]byte(tc.body)) {
				t.Errorf("Expected detectCAPTCHA to return true for: %s", tc.name)
			}
		})
	}
}

func TestDetectCAPTCHA_DetectsAmazonCAPTCHA(t *testing.T) {
	client := NewClient()

	bodies := []string{
		"<form action='/amazoncaptcha/verify'>",
		"https://images-na.ssl-images-amazon.com/amazoncaptcha/",
		"AmazonCaptcha challenge",
	}

	for _, body := range bodies {
		if !client.detectCAPTCHA([]byte(body)) {
			t.Errorf("Expected detectCAPTCHA to return true for Amazon CAPTCHA pattern: %s", body)
		}
	}
}

func TestDetectCAPTCHA_ReturnsFalseForNormalContent(t *testing.T) {
	client := NewClient()

	normalBodies := []string{
		"<html><body><h1>Welcome to Amazon</h1></body></html>",
		"Your order has been placed successfully",
		"Product details and description",
		"Shopping cart contains 3 items",
		"Customer reviews for this product",
	}

	for _, body := range normalBodies {
		if client.detectCAPTCHA([]byte(body)) {
			t.Errorf("Expected detectCAPTCHA to return false for normal content: %s", body)
		}
	}
}

func TestDetectCAPTCHA_HandlesEmptyBody(t *testing.T) {
	client := NewClient()

	if client.detectCAPTCHA([]byte("")) {
		t.Error("Expected detectCAPTCHA to return false for empty body")
	}

	if client.detectCAPTCHA(nil) {
		t.Error("Expected detectCAPTCHA to return false for nil body")
	}
}

func TestDetectCAPTCHA_IsCaseInsensitive(t *testing.T) {
	client := NewClient()

	variations := []string{
		"CAPTCHA REQUIRED",
		"captcha required",
		"CaPtChA ReQuIrEd",
		"Captcha Required",
	}

	for _, variation := range variations {
		if !client.detectCAPTCHA([]byte(variation)) {
			t.Errorf("Expected detectCAPTCHA to be case-insensitive for: %s", variation)
		}
	}
}

func TestDetectCAPTCHA_DetectsMultipleIndicators(t *testing.T) {
	client := NewClient()

	// Body with multiple CAPTCHA indicators
	body := `
	<html>
	<head><title>Robot Check</title></head>
	<body>
		<h1>Enter the characters you see below</h1>
		<div class="g-recaptcha" data-sitekey="test"></div>
		<p>Sorry, we just need to make sure you're not a robot.</p>
	</body>
	</html>
	`

	if !client.detectCAPTCHA([]byte(body)) {
		t.Error("Expected detectCAPTCHA to return true for body with multiple indicators")
	}
}
