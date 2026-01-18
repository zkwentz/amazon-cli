package cmd

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name      string
		errorCode string
		wantExit  int
	}{
		{
			name:      "auth required",
			errorCode: models.ErrCodeAuthRequired,
			wantExit:  3,
		},
		{
			name:      "auth expired",
			errorCode: models.ErrCodeAuthExpired,
			wantExit:  3,
		},
		{
			name:      "network error",
			errorCode: models.ErrCodeNetworkError,
			wantExit:  4,
		},
		{
			name:      "rate limited",
			errorCode: models.ErrCodeRateLimited,
			wantExit:  5,
		},
		{
			name:      "not found",
			errorCode: models.ErrCodeNotFound,
			wantExit:  6,
		},
		{
			name:      "invalid input",
			errorCode: models.ErrCodeInvalidInput,
			wantExit:  2,
		},
		{
			name:      "unknown error",
			errorCode: "SOME_OTHER_ERROR",
			wantExit:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := getExitCode(tt.errorCode)
			if exitCode != tt.wantExit {
				t.Errorf("getExitCode() = %v, want %v", exitCode, tt.wantExit)
			}
		})
	}
}
