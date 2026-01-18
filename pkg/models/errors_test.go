package models

import (
	"errors"
	"testing"
)

func TestCLIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		cliError *CLIError
		want     string
	}{
		{
			name:     "error without details",
			cliError: NewCLIError(ErrCodeAuthRequired, "Authentication required"),
			want:     "AUTH_REQUIRED: Authentication required",
		},
		{
			name: "error with details",
			cliError: &CLIError{
				Code:    ErrCodeNetworkError,
				Message: "Connection failed",
				Details: map[string]interface{}{"host": "amazon.com"},
			},
			want: "NETWORK_ERROR: Connection failed (details: map[host:amazon.com])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cliError.Error(); got != tt.want {
				t.Errorf("CLIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCLIError(t *testing.T) {
	code := ErrCodeInvalidInput
	message := "Invalid ASIN format"

	err := NewCLIError(code, message)

	if err.Code != code {
		t.Errorf("NewCLIError() Code = %v, want %v", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("NewCLIError() Message = %v, want %v", err.Message, message)
	}
	if err.Details == nil {
		t.Error("NewCLIError() Details should be initialized")
	}
}

func TestCLIError_WithDetails(t *testing.T) {
	err := NewCLIError(ErrCodeNotFound, "Order not found")
	err.WithDetails("order_id", "123-456")
	err.WithDetails("timestamp", "2024-01-15")

	if len(err.Details) != 2 {
		t.Errorf("WithDetails() expected 2 details, got %v", len(err.Details))
	}
	if err.Details["order_id"] != "123-456" {
		t.Errorf("WithDetails() order_id = %v, want 123-456", err.Details["order_id"])
	}
	if err.Details["timestamp"] != "2024-01-15" {
		t.Errorf("WithDetails() timestamp = %v, want 2024-01-15", err.Details["timestamp"])
	}
}

func TestCLIError_GetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		wantExit int
	}{
		{
			name:     "AUTH_REQUIRED maps to ExitAuthError",
			code:     ErrCodeAuthRequired,
			wantExit: ExitAuthError,
		},
		{
			name:     "AUTH_EXPIRED maps to ExitAuthError",
			code:     ErrCodeAuthExpired,
			wantExit: ExitAuthError,
		},
		{
			name:     "NETWORK_ERROR maps to ExitNetworkError",
			code:     ErrCodeNetworkError,
			wantExit: ExitNetworkError,
		},
		{
			name:     "RATE_LIMITED maps to ExitRateLimited",
			code:     ErrCodeRateLimited,
			wantExit: ExitRateLimited,
		},
		{
			name:     "NOT_FOUND maps to ExitNotFound",
			code:     ErrCodeNotFound,
			wantExit: ExitNotFound,
		},
		{
			name:     "INVALID_INPUT maps to ExitInvalidArguments",
			code:     ErrCodeInvalidInput,
			wantExit: ExitInvalidArguments,
		},
		{
			name:     "PURCHASE_FAILED maps to ExitGeneralError",
			code:     ErrCodePurchaseFailed,
			wantExit: ExitGeneralError,
		},
		{
			name:     "AMAZON_ERROR maps to ExitGeneralError",
			code:     ErrCodeAmazonError,
			wantExit: ExitGeneralError,
		},
		{
			name:     "GENERAL_ERROR maps to ExitGeneralError",
			code:     ErrCodeGeneralError,
			wantExit: ExitGeneralError,
		},
		{
			name:     "unknown error code maps to ExitGeneralError",
			code:     "UNKNOWN_CODE",
			wantExit: ExitGeneralError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCLIError(tt.code, "test message")
			if got := err.GetExitCode(); got != tt.wantExit {
				t.Errorf("GetExitCode() = %v, want %v", got, tt.wantExit)
			}
		})
	}
}

func TestGetExitCodeFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantExit int
	}{
		{
			name:     "nil error returns ExitSuccess",
			err:      nil,
			wantExit: ExitSuccess,
		},
		{
			name:     "CLIError with AUTH_REQUIRED returns ExitAuthError",
			err:      NewCLIError(ErrCodeAuthRequired, "auth required"),
			wantExit: ExitAuthError,
		},
		{
			name:     "CLIError with NETWORK_ERROR returns ExitNetworkError",
			err:      NewCLIError(ErrCodeNetworkError, "network error"),
			wantExit: ExitNetworkError,
		},
		{
			name:     "CLIError with RATE_LIMITED returns ExitRateLimited",
			err:      NewCLIError(ErrCodeRateLimited, "rate limited"),
			wantExit: ExitRateLimited,
		},
		{
			name:     "CLIError with NOT_FOUND returns ExitNotFound",
			err:      NewCLIError(ErrCodeNotFound, "not found"),
			wantExit: ExitNotFound,
		},
		{
			name:     "CLIError with INVALID_INPUT returns ExitInvalidArguments",
			err:      NewCLIError(ErrCodeInvalidInput, "invalid input"),
			wantExit: ExitInvalidArguments,
		},
		{
			name:     "standard error returns ExitGeneralError",
			err:      errors.New("standard error"),
			wantExit: ExitGeneralError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExitCodeFromError(tt.err); got != tt.wantExit {
				t.Errorf("GetExitCodeFromError() = %v, want %v", got, tt.wantExit)
			}
		})
	}
}

func TestExitCodeConstants(t *testing.T) {
	// Verify exit codes match PRD specification
	tests := []struct {
		name  string
		code  int
		value int
	}{
		{"ExitSuccess", ExitSuccess, 0},
		{"ExitGeneralError", ExitGeneralError, 1},
		{"ExitInvalidArguments", ExitInvalidArguments, 2},
		{"ExitAuthError", ExitAuthError, 3},
		{"ExitNetworkError", ExitNetworkError, 4},
		{"ExitRateLimited", ExitRateLimited, 5},
		{"ExitNotFound", ExitNotFound, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.value {
				t.Errorf("%s = %v, want %v", tt.name, tt.code, tt.value)
			}
		})
	}
}

func TestErrorCodeConstants(t *testing.T) {
	// Verify error codes match PRD specification
	tests := []struct {
		name  string
		code  string
		value string
	}{
		{"ErrCodeAuthRequired", ErrCodeAuthRequired, "AUTH_REQUIRED"},
		{"ErrCodeAuthExpired", ErrCodeAuthExpired, "AUTH_EXPIRED"},
		{"ErrCodeNotFound", ErrCodeNotFound, "NOT_FOUND"},
		{"ErrCodeRateLimited", ErrCodeRateLimited, "RATE_LIMITED"},
		{"ErrCodeInvalidInput", ErrCodeInvalidInput, "INVALID_INPUT"},
		{"ErrCodePurchaseFailed", ErrCodePurchaseFailed, "PURCHASE_FAILED"},
		{"ErrCodeNetworkError", ErrCodeNetworkError, "NETWORK_ERROR"},
		{"ErrCodeAmazonError", ErrCodeAmazonError, "AMAZON_ERROR"},
		{"ErrCodeGeneralError", ErrCodeGeneralError, "GENERAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.value {
				t.Errorf("%s = %v, want %v", tt.name, tt.code, tt.value)
			}
		})
	}
}

// TestCLIError_ChainedWithDetails tests that WithDetails returns the same error
// so it can be chained
func TestCLIError_ChainedWithDetails(t *testing.T) {
	err := NewCLIError(ErrCodeNetworkError, "Connection timeout").
		WithDetails("host", "amazon.com").
		WithDetails("timeout", "30s").
		WithDetails("retries", 3)

	if len(err.Details) != 3 {
		t.Errorf("Expected 3 details, got %d", len(err.Details))
	}

	if err.Details["host"] != "amazon.com" {
		t.Errorf("Expected host to be amazon.com, got %v", err.Details["host"])
	}
	if err.Details["timeout"] != "30s" {
		t.Errorf("Expected timeout to be 30s, got %v", err.Details["timeout"])
	}
	if err.Details["retries"] != 3 {
		t.Errorf("Expected retries to be 3, got %v", err.Details["retries"])
	}
}
