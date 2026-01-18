package main

import (
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "invalid input error",
			err:      models.NewCLIError(models.ErrInvalidInput, "test"),
			expected: 2,
		},
		{
			name:     "auth required error",
			err:      models.NewCLIError(models.ErrAuthRequired, "test"),
			expected: 3,
		},
		{
			name:     "auth expired error",
			err:      models.NewCLIError(models.ErrAuthExpired, "test"),
			expected: 3,
		},
		{
			name:     "network error",
			err:      models.NewCLIError(models.ErrNetworkError, "test"),
			expected: 4,
		},
		{
			name:     "rate limited error",
			err:      models.NewCLIError(models.ErrRateLimited, "test"),
			expected: 5,
		},
		{
			name:     "not found error",
			err:      models.NewCLIError(models.ErrNotFound, "test"),
			expected: 6,
		},
		{
			name:     "general cli error",
			err:      models.NewCLIError(models.ErrAmazonError, "test"),
			expected: 1,
		},
		{
			name:     "non-cli error",
			err:      &testError{msg: "test error"},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getExitCode(tt.err)
			if result != tt.expected {
				t.Errorf("Expected exit code %d, got %d", tt.expected, result)
			}
		})
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
