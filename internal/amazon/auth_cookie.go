package amazon

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/zkwentz/amazon-cli/internal/config"
)

const (
	amazonLoginURL = "https://www.amazon.com/ap/signin"
	amazonHomeURL  = "https://www.amazon.com"
	cookieMaxAge   = 365 * 24 * time.Hour // Cookies older than this need refresh
)

// CookieAuth handles cookie-based authentication
type CookieAuth struct {
	config *config.Config
}

// NewCookieAuth creates a new cookie authentication handler
func NewCookieAuth(cfg *config.Config) *CookieAuth {
	return &CookieAuth{
		config: cfg,
	}
}

// LoginWithBrowser opens a browser for the user to login and captures session cookies
func (ca *CookieAuth) LoginWithBrowser() error {
	// Launch browser
	url := launcher.New().
		Headless(false). // Show browser so user can login
		MustLaunch()

	browser := rod.New().
		ControlURL(url).
		MustConnect()
	defer browser.MustClose()

	// Navigate to Amazon login page
	page := browser.MustPage(amazonLoginURL)
	defer page.MustClose()

	fmt.Println("Browser opened. Please login to Amazon...")
	fmt.Println("Waiting for successful login...")

	// Wait for user to login by checking if we've reached the home page
	// or any page that's not the login page
	maxWaitTime := 5 * time.Minute
	startTime := time.Now()

	for {
		if time.Since(startTime) > maxWaitTime {
			return fmt.Errorf("login timeout: no successful login detected within %v", maxWaitTime)
		}

		// Check current URL
		currentURL := page.MustInfo().URL

		// If we're no longer on the login page, assume successful login
		if currentURL != amazonLoginURL &&
		   !contains(currentURL, "/ap/signin") &&
		   !contains(currentURL, "/ap/mfa") &&
		   !contains(currentURL, "/ap/cvf") {
			fmt.Println("Login detected! Capturing session cookies...")
			break
		}

		time.Sleep(1 * time.Second)
	}

	// Give it a moment to ensure cookies are set
	time.Sleep(2 * time.Second)

	// Get all cookies
	cookies, err := page.Cookies([]string{amazonHomeURL})
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	// Convert Rod cookies to our config cookies
	configCookies := make([]config.Cookie, 0, len(cookies))
	for _, cookie := range cookies {
		configCookies = append(configCookies, config.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  time.Unix(int64(cookie.Expires), 0),
			Secure:   cookie.Secure,
			HttpOnly: cookie.HTTPOnly,
		})
	}

	// Update config with cookies
	ca.config.Auth.Cookies = configCookies
	ca.config.Auth.CookiesSetAt = time.Now()
	ca.config.Auth.AuthMethod = "cookie"

	fmt.Printf("Successfully captured %d session cookies\n", len(configCookies))
	return nil
}

// ApplyCookies applies stored cookies to an HTTP client's cookie jar
func (ca *CookieAuth) ApplyCookies(client *http.Client) error {
	if len(ca.config.Auth.Cookies) == 0 {
		return fmt.Errorf("no cookies stored, please login first")
	}

	// Get the cookie jar from the client
	jar := client.Jar
	if jar == nil {
		return fmt.Errorf("http client has no cookie jar")
	}

	// Convert config cookies to http.Cookie and add to jar
	amazonURL := mustParseURL(amazonHomeURL)
	httpCookies := make([]*http.Cookie, 0, len(ca.config.Auth.Cookies))

	for _, cookie := range ca.config.Auth.Cookies {
		httpCookies = append(httpCookies, &http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  cookie.Expires,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HttpOnly,
		})
	}

	jar.SetCookies(amazonURL, httpCookies)
	return nil
}

// NeedsRefresh checks if cookies need to be refreshed
func (ca *CookieAuth) NeedsRefresh() bool {
	if len(ca.config.Auth.Cookies) == 0 {
		return true
	}

	// Check if cookies are getting old (refresh when 80% of maxAge is reached)
	age := time.Since(ca.config.Auth.CookiesSetAt)
	return age > cookieMaxAge*8/10
}

// ValidateCookies makes a test request to verify cookies are still valid
func (ca *CookieAuth) ValidateCookies(client *http.Client) error {
	req, err := http.NewRequest("GET", amazonHomeURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check if we're being redirected to login page
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther {
		location := resp.Header.Get("Location")
		if contains(location, "/ap/signin") {
			return fmt.Errorf("cookies invalid: redirected to login page")
		}
	}

	// Check if response indicates we're not authenticated
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("cookies invalid: unauthorized response")
	}

	return nil
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		   (s == substr ||
		    (len(s) > len(substr) &&
		     (s[:len(substr)] == substr ||
		      findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
