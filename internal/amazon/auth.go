package amazon

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// OAuth constants for Amazon Login with Amazon (LWA)
const (
	// AuthURL is the Amazon OAuth authorization endpoint
	AuthURL = "https://www.amazon.com/ap/oa"

	// TokenURL is the Amazon OAuth token endpoint
	TokenURL = "https://api.amazon.com/auth/o2/token"

	// DefaultRedirectURI is the local callback URI for OAuth flow
	DefaultRedirectURI = "http://localhost:8085/callback"

	// DefaultPort is the default port for the local OAuth callback server
	DefaultPort = 8085

	// OAuthTimeout is the maximum time to wait for user to complete OAuth flow
	OAuthTimeout = 2 * time.Minute
)

// OAuth scopes required for Amazon CLI functionality
var (
	// RequiredScopes lists the OAuth scopes needed for order/profile access
	RequiredScopes = []string{
		"profile",                    // Basic profile information
		"postal_code",                // Shipping address access
		"profile:user_id",            // User ID access
	}
)

// AuthTokens represents the OAuth tokens returned by Amazon
type AuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// TokenResponse represents the response from Amazon's token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// AuthConfig holds the authentication configuration
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

// Authenticator manages the OAuth authentication flow
type Authenticator struct {
	config AuthConfig
	client *http.Client
}

// NewAuthenticator creates a new authenticator with the given configuration
func NewAuthenticator(config AuthConfig) *Authenticator {
	if config.RedirectURI == "" {
		config.RedirectURI = DefaultRedirectURI
	}
	if len(config.Scopes) == 0 {
		config.Scopes = RequiredScopes
	}

	return &Authenticator{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateAuthURL generates the OAuth authorization URL with state parameter
func (a *Authenticator) GenerateAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", a.config.ClientID)
	params.Set("scope", joinScopes(a.config.Scopes))
	params.Set("response_type", "code")
	params.Set("redirect_uri", a.config.RedirectURI)
	params.Set("state", state)

	return fmt.Sprintf("%s?%s", AuthURL, params.Encode())
}

// GenerateState generates a random state parameter for CSRF protection
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ExchangeCodeForTokens exchanges an authorization code for access and refresh tokens
func (a *Authenticator) ExchangeCodeForTokens(ctx context.Context, code string) (*AuthTokens, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", a.config.RedirectURI)
	data.Set("client_id", a.config.ClientID)
	data.Set("client_secret", a.config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", TokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.URL.RawQuery = data.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	tokens := &AuthTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	return tokens, nil
}

// RefreshAccessToken uses the refresh token to obtain a new access token
func (a *Authenticator) RefreshAccessToken(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", a.config.ClientID)
	data.Set("client_secret", a.config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", TokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.URL.RawQuery = data.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	tokens := &AuthTokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: refreshToken, // Keep existing refresh token if not returned
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	// If a new refresh token is provided, use it
	if tokenResp.RefreshToken != "" {
		tokens.RefreshToken = tokenResp.RefreshToken
	}

	return tokens, nil
}

// IsTokenExpired checks if the access token is expired or will expire within the given duration
func IsTokenExpired(tokens *AuthTokens, buffer time.Duration) bool {
	if tokens == nil || tokens.AccessToken == "" {
		return true
	}
	return time.Now().Add(buffer).After(tokens.ExpiresAt)
}

// ShouldRefreshToken checks if the token should be refreshed (expires within 5 minutes)
func ShouldRefreshToken(tokens *AuthTokens) bool {
	return IsTokenExpired(tokens, 5*time.Minute)
}

// ValidateTokens validates that the tokens are present and not expired
func ValidateTokens(tokens *AuthTokens) error {
	if tokens == nil {
		return fmt.Errorf("no authentication tokens provided")
	}
	if tokens.AccessToken == "" {
		return fmt.Errorf("access token is empty")
	}
	if tokens.RefreshToken == "" {
		return fmt.Errorf("refresh token is empty")
	}
	if IsTokenExpired(tokens, 0) {
		return fmt.Errorf("access token is expired")
	}
	return nil
}

// OAuthCallbackResult represents the result of the OAuth callback
type OAuthCallbackResult struct {
	Code  string
	State string
	Error string
}

// StartCallbackServer starts a local HTTP server to handle the OAuth callback
func StartCallbackServer(ctx context.Context, expectedState string) (*AuthTokens, error) {
	resultChan := make(chan OAuthCallbackResult, 1)
	errChan := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		errorParam := r.URL.Query().Get("error")

		result := OAuthCallbackResult{
			Code:  code,
			State: state,
			Error: errorParam,
		}

		if errorParam != "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "<html><body><h1>Authentication Failed</h1><p>Error: %s</p><p>You can close this window.</p></body></html>", errorParam)
		} else if state != expectedState {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "<html><body><h1>Authentication Failed</h1><p>Invalid state parameter (CSRF protection)</p><p>You can close this window.</p></body></html>")
			result.Error = "invalid_state"
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "<html><body><h1>Authentication Successful!</h1><p>You can close this window and return to the terminal.</p></body></html>")
		}

		resultChan <- result
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", DefaultPort),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("callback server error: %w", err)
		}
	}()

	// Wait for callback or timeout
	select {
	case result := <-resultChan:
		// Shutdown server
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)

		if result.Error != "" {
			return nil, fmt.Errorf("authentication error: %s", result.Error)
		}

		return nil, nil // Return nil to indicate callback was received, tokens will be exchanged separately

	case err := <-errChan:
		return nil, err

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
		return nil, fmt.Errorf("authentication timeout: user did not complete login within %v", OAuthTimeout)
	}
}

// joinScopes joins OAuth scopes with space delimiter
func joinScopes(scopes []string) string {
	if len(scopes) == 0 {
		return ""
	}
	result := scopes[0]
	for i := 1; i < len(scopes); i++ {
		result += " " + scopes[i]
	}
	return result
}
