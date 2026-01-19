package amazon

import "time"

// AuthTokens represents Amazon authentication tokens
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// IsExpired checks if the access token has expired
func (a *AuthTokens) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}

// ExpiresWithin checks if the access token will expire within the given duration
func (a *AuthTokens) ExpiresWithin(duration time.Duration) bool {
	return time.Now().Add(duration).After(a.ExpiresAt)
}

// RefreshTokens refreshes the authentication tokens using a refresh token
// This is a placeholder implementation that returns mock tokens for now
func RefreshTokens(refreshToken string) (*AuthTokens, error) {
	// Mock implementation - returns tokens valid for 1 hour
	return &AuthTokens{
		AccessToken:  "mock_access_token_" + refreshToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}, nil
}
