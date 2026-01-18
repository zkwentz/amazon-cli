package main

import (
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// TestExitCodeMapping verifies that different error types result in correct exit codes
func TestExitCodeMapping(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantExit int
	}{
		{
			name:     "nil error should exit with 0",
			err:      nil,
			wantExit: models.ExitSuccess,
		},
		{
			name:     "AUTH_REQUIRED should exit with 3",
			err:      models.NewCLIError(models.ErrCodeAuthRequired, "Not authenticated"),
			wantExit: models.ExitAuthError,
		},
		{
			name:     "AUTH_EXPIRED should exit with 3",
			err:      models.NewCLIError(models.ErrCodeAuthExpired, "Token expired"),
			wantExit: models.ExitAuthError,
		},
		{
			name:     "NOT_FOUND should exit with 6",
			err:      models.NewCLIError(models.ErrCodeNotFound, "Resource not found"),
			wantExit: models.ExitNotFound,
		},
		{
			name:     "RATE_LIMITED should exit with 5",
			err:      models.NewCLIError(models.ErrCodeRateLimited, "Too many requests"),
			wantExit: models.ExitRateLimited,
		},
		{
			name:     "INVALID_INPUT should exit with 2",
			err:      models.NewCLIError(models.ErrCodeInvalidInput, "Invalid input"),
			wantExit: models.ExitInvalidArgs,
		},
		{
			name:     "NETWORK_ERROR should exit with 4",
			err:      models.NewCLIError(models.ErrCodeNetworkError, "Network error"),
			wantExit: models.ExitNetworkError,
		},
		{
			name:     "PURCHASE_FAILED should exit with 1",
			err:      models.NewCLIError(models.ErrCodePurchaseFailed, "Purchase failed"),
			wantExit: models.ExitGeneralError,
		},
		{
			name:     "AMAZON_ERROR should exit with 1",
			err:      models.NewCLIError(models.ErrCodeAmazonError, "Amazon error"),
			wantExit: models.ExitGeneralError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExit := models.GetExitCodeFromError(tt.err)
			if gotExit != tt.wantExit {
				t.Errorf("GetExitCodeFromError(%v) = %d, want %d", tt.err, gotExit, tt.wantExit)
			}
		})
	}
}
