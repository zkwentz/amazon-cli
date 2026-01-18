package amazon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
)

// OAuth constants
const (
	// Amazon Login with Amazon (LWA) endpoints
	TokenURL = "https://api.amazon.com/auth/o2/token"

	// Time buffer before token expiry to trigger refresh
	RefreshBuffer = 5 * time.Minute
)

// TokenResponse represents the response from Amazon's token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RefreshTokenIfNeeded checks if the access token needs to be refreshed
// and refreshes it if it expires within 5 minutes
func RefreshTokenIfNeeded(cfg *config.Config) error {
	// Check if we have a refresh token
	if cfg.Auth.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	// Check if access token is expired or will expire soon
	if time.Until(cfg.Auth.ExpiresAt) > RefreshBuffer {
		// Token is still valid for more than 5 minutes
		return nil
	}

	// Token needs refresh - make request to Amazon
	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {cfg.Auth.RefreshToken},
		// Note: client_id and client_secret would be needed here for actual OAuth
		// These would be obtained from Amazon Developer Console
		// "client_id":     {os.Getenv("AMAZON_CLIENT_ID")},
		// "client_secret": {os.Getenv("AMAZON_CLIENT_SECRET")},
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", TokenURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Update config with new tokens
	cfg.Auth.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		// Amazon may issue a new refresh token
		cfg.Auth.RefreshToken = tokenResp.RefreshToken
	}
	cfg.Auth.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Save config to disk
	configPath := config.GetConfigPath()
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save refreshed config: %w", err)
	}

	return nil
}

// IsAuthenticated checks if the config has valid authentication
func IsAuthenticated(cfg *config.Config) bool {
	return cfg.Auth.AccessToken != "" && cfg.Auth.RefreshToken != ""
}

// NeedsRefresh checks if the token needs to be refreshed
func NeedsRefresh(cfg *config.Config) bool {
	return time.Until(cfg.Auth.ExpiresAt) <= RefreshBuffer
}
