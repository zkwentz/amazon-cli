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
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/zkwentz/amazon-cli/internal/config"
)

const (
	// Amazon OAuth endpoints
	AuthURL  = "https://www.amazon.com/ap/oa"
	TokenURL = "https://api.amazon.com/auth/o2/token"

	// OAuth scopes
	Scopes = "profile postal_code"

	// Redirect URI for local server
	RedirectURITemplate = "http://localhost:%d/callback"
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthClient struct {
	ClientID     string
	ClientSecret string
	Config       *config.Config
}

func NewAuthClient(cfg *config.Config) *AuthClient {
	// Note: In production, these should be set from environment variables or config
	// For now, using placeholder values - user needs to register app at https://developer.amazon.com/
	return &AuthClient{
		ClientID:     os.Getenv("AMAZON_CLIENT_ID"),
		ClientSecret: os.Getenv("AMAZON_CLIENT_SECRET"),
		Config:       cfg,
	}
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (a *AuthClient) Login() (*AuthTokens, error) {
	state, err := generateState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	// Start local HTTP server for OAuth callback
	port := 8085
	redirectURI := fmt.Sprintf(RedirectURITemplate, port)

	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Verify state parameter
		if r.URL.Query().Get("state") != state {
			errChan <- fmt.Errorf("invalid state parameter")
			fmt.Fprintf(w, "<html><body><h1>Authentication Failed</h1><p>Invalid state parameter. You can close this window.</p></body></html>")
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no authorization code received")
			fmt.Fprintf(w, "<html><body><h1>Authentication Failed</h1><p>No authorization code received. You can close this window.</p></body></html>")
			return
		}

		codeChan <- code
		fmt.Fprintf(w, "<html><body><h1>Authentication Successful!</h1><p>You can close this window and return to the CLI.</p></body></html>")
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Build authorization URL
	authURL := fmt.Sprintf("%s?client_id=%s&scope=%s&response_type=code&redirect_uri=%s&state=%s",
		AuthURL,
		url.QueryEscape(a.ClientID),
		url.QueryEscape(Scopes),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)

	// Open browser
	if err := browser.OpenURL(authURL); err != nil {
		return nil, fmt.Errorf("failed to open browser: %w", err)
	}

	fmt.Println("Opening browser for authentication...")
	fmt.Println("If the browser doesn't open automatically, please visit:")
	fmt.Println(authURL)

	// Wait for callback with timeout
	var code string
	select {
	case code = <-codeChan:
		// Success
	case err := <-errChan:
		return nil, err
	case <-time.After(2 * time.Minute):
		return nil, fmt.Errorf("authentication timeout - please try again")
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	// Exchange code for tokens
	tokens, err := a.exchangeCodeForTokens(code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}

	return tokens, nil
}

func (a *AuthClient) exchangeCodeForTokens(code, redirectURI string) (*AuthTokens, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"client_id":     {a.ClientID},
		"client_secret": {a.ClientSecret},
	}

	resp, err := http.PostForm(TokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokens AuthTokens
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (a *AuthClient) RefreshTokenIfNeeded() error {
	// Check if token expires within 5 minutes
	if time.Until(a.Config.Auth.ExpiresAt) > 5*time.Minute {
		return nil
	}

	if a.Config.Auth.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {a.Config.Auth.RefreshToken},
		"client_id":     {a.ClientID},
		"client_secret": {a.ClientSecret},
	}

	resp, err := http.PostForm(TokenURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokens AuthTokens
	if err := json.Unmarshal(body, &tokens); err != nil {
		return err
	}

	// Update config with new tokens
	a.Config.Auth.AccessToken = tokens.AccessToken
	if tokens.RefreshToken != "" {
		a.Config.Auth.RefreshToken = tokens.RefreshToken
	}
	a.Config.Auth.ExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)

	return config.SaveConfig(a.Config, config.GetConfigPath())
}
