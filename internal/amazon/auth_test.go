package amazon

import (
	"testing"
	"time"
)

func TestAuthTokens_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "token expired",
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      true,
		},
		{
			name:      "token not expired",
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      false,
		},
		{
			name:      "token expired just now",
			expiresAt: time.Now().Add(-1 * time.Second),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthTokens{
				AccessToken:  "test_token",
				RefreshToken: "test_refresh",
				ExpiresAt:    tt.expiresAt,
			}
			if got := a.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthTokens_ExpiresWithin(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		duration  time.Duration
		want      bool
	}{
		{
			name:      "expires within 30 minutes - true",
			expiresAt: time.Now().Add(15 * time.Minute),
			duration:  30 * time.Minute,
			want:      true,
		},
		{
			name:      "expires within 30 minutes - false",
			expiresAt: time.Now().Add(45 * time.Minute),
			duration:  30 * time.Minute,
			want:      false,
		},
		{
			name:      "already expired",
			expiresAt: time.Now().Add(-1 * time.Hour),
			duration:  30 * time.Minute,
			want:      true,
		},
		{
			name:      "expires beyond duration boundary",
			expiresAt: time.Now().Add(2 * time.Hour),
			duration:  1 * time.Hour,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthTokens{
				AccessToken:  "test_token",
				RefreshToken: "test_refresh",
				ExpiresAt:    tt.expiresAt,
			}
			if got := a.ExpiresWithin(tt.duration); got != tt.want {
				t.Errorf("ExpiresWithin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshTokens(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken string
	}{
		{
			name:         "refresh with valid token",
			refreshToken: "valid_refresh_token",
		},
		{
			name:         "refresh with empty token",
			refreshToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RefreshTokens(tt.refreshToken)
			if err != nil {
				t.Errorf("RefreshTokens() error = %v", err)
				return
			}
			if got == nil {
				t.Error("RefreshTokens() returned nil")
				return
			}
			if got.AccessToken == "" {
				t.Error("RefreshTokens() returned empty AccessToken")
			}
			if got.RefreshToken != tt.refreshToken {
				t.Errorf("RefreshTokens() RefreshToken = %v, want %v", got.RefreshToken, tt.refreshToken)
			}
			if got.ExpiresAt.Before(time.Now()) {
				t.Error("RefreshTokens() returned already expired token")
			}
			// Check that expiration is set to approximately 1 hour from now
			expectedExpiry := time.Now().Add(1 * time.Hour)
			if got.ExpiresAt.Before(expectedExpiry.Add(-1*time.Second)) || got.ExpiresAt.After(expectedExpiry.Add(1*time.Second)) {
				t.Errorf("RefreshTokens() ExpiresAt not set to approximately 1 hour from now")
			}
		})
	}
}
