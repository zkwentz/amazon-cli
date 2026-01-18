package amazon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// OAuth constants (placeholder - would need real Amazon OAuth endpoints)
const (
	AuthURL = "https://www.amazon.com/ap/oa"
	// RedirectURI would be localhost for OAuth callback
	RedirectURI = "http://localhost:8085/callback"
)

// TokenURL is a variable so it can be overridden in tests
var TokenURL = "https://api.amazon.com/auth/o2/token"

// AuthTokens represents OAuth tokens
type AuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshTokenIfNeeded checks if the access token is expired or about to expire
// and refreshes it if necessary. This is the auth check middleware function.
func RefreshTokenIfNeeded(cfg *config.Config, configPath string) error {
	// Check if tokens exist
	if cfg.Auth.AccessToken == "" {
		return models.NewCLIError(
			models.ErrCodeAuthRequired,
			"Authentication required. Run 'amazon-cli auth login' to authenticate.",
			nil,
		)
	}

	// Check if token expires within 5 minutes
	expiresIn := time.Until(cfg.Auth.ExpiresAt)
	if expiresIn > 5*time.Minute {
		// Token is still valid
		return nil
	}

	// Check if we have a refresh token
	if cfg.Auth.RefreshToken == "" {
		return models.NewCLIError(
			models.ErrCodeAuthExpired,
			"Authentication token has expired and no refresh token available. Run 'amazon-cli auth login' to re-authenticate.",
			nil,
		)
	}

	// Refresh the token
	newTokens, err := refreshAccessToken(cfg.Auth.RefreshToken)
	if err != nil {
		return models.NewCLIError(
			models.ErrCodeAuthExpired,
			fmt.Sprintf("Failed to refresh authentication token: %v. Run 'amazon-cli auth login' to re-authenticate.", err),
			map[string]interface{}{
				"error": err.Error(),
			},
		)
	}

	// Update config with new tokens
	cfg.Auth.AccessToken = newTokens.AccessToken
	cfg.Auth.RefreshToken = newTokens.RefreshToken
	cfg.Auth.ExpiresAt = newTokens.ExpiresAt

	// Save updated config to disk
	if err := config.SaveConfig(cfg, configPath); err != nil {
		// Log the error but don't fail - we have the tokens in memory
		fmt.Fprintf(os.Stderr, "Warning: Failed to save refreshed tokens to config: %v\n", err)
	}

	return nil
}

// refreshAccessToken makes a request to Amazon's OAuth token endpoint to refresh the access token
func refreshAccessToken(refreshToken string) (*AuthTokens, error) {
	// Prepare the refresh token request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	// Note: In a real implementation, you would need client_id and client_secret
	// data.Set("client_id", clientID)
	// data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send refresh token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh token request failed with status %d", resp.StatusCode)
	}

	var tokens AuthTokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to decode refresh token response: %w", err)
	}

	// Calculate expiry time
	tokens.ExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)

	return &tokens, nil
}

// IsAuthenticated checks if the user has valid authentication credentials
func IsAuthenticated(cfg *config.Config) bool {
	return cfg.Auth.AccessToken != "" && time.Now().Before(cfg.Auth.ExpiresAt)
}

// CheckAuth is a convenience function that checks authentication and refreshes if needed
// This is the main middleware function to be called before any authenticated request
func CheckAuth(cfg *config.Config, configPath string) error {
	return RefreshTokenIfNeeded(cfg, configPath)
}
