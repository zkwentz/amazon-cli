package amazon

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// AuthManager handles both OAuth and cookie-based authentication
type AuthManager struct {
	config     *config.Config
	cookieAuth *CookieAuth
	httpClient *http.Client
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(cfg *config.Config) (*AuthManager, error) {
	// Create cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Create HTTP client
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	am := &AuthManager{
		config:     cfg,
		cookieAuth: NewCookieAuth(cfg),
		httpClient: httpClient,
	}

	// If using cookie auth, apply cookies to client
	if cfg.Auth.AuthMethod == "cookie" && len(cfg.Auth.Cookies) > 0 {
		if err := am.cookieAuth.ApplyCookies(httpClient); err != nil {
			return nil, fmt.Errorf("failed to apply cookies: %w", err)
		}
	}

	return am, nil
}

// Login initiates the login flow based on the authentication method
func (am *AuthManager) Login(useBrowser bool) error {
	if useBrowser {
		// Use cookie-based auth
		if err := am.cookieAuth.LoginWithBrowser(); err != nil {
			return fmt.Errorf("browser login failed: %w", err)
		}

		// Apply cookies to HTTP client
		if err := am.cookieAuth.ApplyCookies(am.httpClient); err != nil {
			return fmt.Errorf("failed to apply cookies: %w", err)
		}

		return nil
	}

	// OAuth flow (to be implemented)
	return fmt.Errorf("OAuth authentication not yet implemented, use --browser flag for cookie-based auth")
}

// IsAuthenticated checks if the user is authenticated
func (am *AuthManager) IsAuthenticated() bool {
	switch am.config.Auth.AuthMethod {
	case "cookie":
		return len(am.config.Auth.Cookies) > 0
	case "oauth":
		return am.config.Auth.AccessToken != ""
	default:
		return false
	}
}

// NeedsRefresh checks if authentication needs to be refreshed
func (am *AuthManager) NeedsRefresh() bool {
	switch am.config.Auth.AuthMethod {
	case "cookie":
		return am.cookieAuth.NeedsRefresh()
	case "oauth":
		// Check if OAuth token expires within 5 minutes
		return time.Until(am.config.Auth.ExpiresAt) < 5*time.Minute
	default:
		return true
	}
}

// ValidateAuth validates that the current authentication is still valid
func (am *AuthManager) ValidateAuth() error {
	if !am.IsAuthenticated() {
		return fmt.Errorf("not authenticated, please run 'amazon-cli auth login --browser'")
	}

	switch am.config.Auth.AuthMethod {
	case "cookie":
		return am.cookieAuth.ValidateCookies(am.httpClient)
	case "oauth":
		// OAuth validation would go here
		return fmt.Errorf("OAuth validation not yet implemented")
	default:
		return fmt.Errorf("unknown authentication method: %s", am.config.Auth.AuthMethod)
	}
}

// RefreshAuthIfNeeded refreshes authentication if needed
func (am *AuthManager) RefreshAuthIfNeeded() error {
	if !am.NeedsRefresh() {
		return nil
	}

	fmt.Println("Authentication needs refresh...")

	switch am.config.Auth.AuthMethod {
	case "cookie":
		fmt.Println("Please re-authenticate using: amazon-cli auth login --browser")
		return fmt.Errorf("cookies expired, re-authentication required")
	case "oauth":
		// OAuth refresh would go here
		return fmt.Errorf("OAuth refresh not yet implemented")
	default:
		return fmt.Errorf("unknown authentication method: %s", am.config.Auth.AuthMethod)
	}
}

// Logout clears stored authentication
func (am *AuthManager) Logout() error {
	am.config.Auth.AccessToken = ""
	am.config.Auth.RefreshToken = ""
	am.config.Auth.ExpiresAt = time.Time{}
	am.config.Auth.Cookies = nil
	am.config.Auth.CookiesSetAt = time.Time{}
	am.config.Auth.AuthMethod = ""

	return nil
}

// GetHTTPClient returns the authenticated HTTP client
func (am *AuthManager) GetHTTPClient() *http.Client {
	return am.httpClient
}

// GetAuthStatus returns the current authentication status
func (am *AuthManager) GetAuthStatus() map[string]interface{} {
	status := make(map[string]interface{})
	status["authenticated"] = am.IsAuthenticated()
	status["auth_method"] = am.config.Auth.AuthMethod

	switch am.config.Auth.AuthMethod {
	case "cookie":
		status["cookies_count"] = len(am.config.Auth.Cookies)
		if !am.config.Auth.CookiesSetAt.IsZero() {
			status["cookies_set_at"] = am.config.Auth.CookiesSetAt
			age := time.Since(am.config.Auth.CookiesSetAt)
			status["cookies_age_hours"] = int(age.Hours())
		}
		status["needs_refresh"] = am.cookieAuth.NeedsRefresh()

	case "oauth":
		if !am.config.Auth.ExpiresAt.IsZero() {
			status["expires_at"] = am.config.Auth.ExpiresAt
			status["expires_in_seconds"] = int(time.Until(am.config.Auth.ExpiresAt).Seconds())
		}
	}

	return status
}
