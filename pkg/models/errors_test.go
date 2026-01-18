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
			name: "simple error",
			cliError: &CLIError{
				Code:    ErrCodeAuthRequired,
				Message: "Not logged in",
			},
			want: "AUTH_REQUIRED: Not logged in",
		},
		{
			name: "error with details",
			cliError: &CLIError{
				Code:    ErrCodeNotFound,
				Message: "Order not found",
				Details: map[string]interface{}{"order_id": "123"},
			},
			want: "NOT_FOUND: Order not found",
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
	code := ErrCodeNetworkError
	message := "Connection timeout"

	err := NewCLIError(code, message)

	if err.Code != code {
		t.Errorf("NewCLIError() Code = %v, want %v", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("NewCLIError() Message = %v, want %v", err.Message, message)
	}
	if err.Details == nil {
		t.Error("NewCLIError() Details should be initialized, got nil")
	}
}

func TestNewCLIErrorWithDetails(t *testing.T) {
	code := ErrCodeInvalidInput
	message := "Invalid ASIN format"
	details := map[string]interface{}{
		"asin":  "INVALID",
		"field": "asin",
	}

	err := NewCLIErrorWithDetails(code, message, details)

	if err.Code != code {
		t.Errorf("NewCLIErrorWithDetails() Code = %v, want %v", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("NewCLIErrorWithDetails() Message = %v, want %v", err.Message, message)
	}
	if err.Details == nil {
		t.Error("NewCLIErrorWithDetails() Details should not be nil")
	}
	if err.Details["asin"] != "INVALID" {
		t.Errorf("NewCLIErrorWithDetails() Details[asin] = %v, want INVALID", err.Details["asin"])
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
			name:     "NOT_FOUND maps to ExitNotFound",
			code:     ErrCodeNotFound,
			wantExit: ExitNotFound,
		},
		{
			name:     "RATE_LIMITED maps to ExitRateLimited",
			code:     ErrCodeRateLimited,
			wantExit: ExitRateLimited,
		},
		{
			name:     "INVALID_INPUT maps to ExitInvalidArgs",
			code:     ErrCodeInvalidInput,
			wantExit: ExitInvalidArgs,
		},
		{
			name:     "NETWORK_ERROR maps to ExitNetworkError",
			code:     ErrCodeNetworkError,
			wantExit: ExitNetworkError,
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
			name:     "unknown error code maps to ExitGeneralError",
			code:     "UNKNOWN_ERROR",
			wantExit: ExitGeneralError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &CLIError{
				Code:    tt.code,
				Message: "test message",
			}
			if got := err.GetExitCode(); got != tt.wantExit {
				t.Errorf("CLIError.GetExitCode() = %v, want %v", got, tt.wantExit)
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
			name:     "CLIError with AUTH_REQUIRED",
			err:      NewCLIError(ErrCodeAuthRequired, "Not logged in"),
			wantExit: ExitAuthError,
		},
		{
			name:     "CLIError with NOT_FOUND",
			err:      NewCLIError(ErrCodeNotFound, "Resource not found"),
			wantExit: ExitNotFound,
		},
		{
			name:     "CLIError with RATE_LIMITED",
			err:      NewCLIError(ErrCodeRateLimited, "Too many requests"),
			wantExit: ExitRateLimited,
		},
		{
			name:     "CLIError with NETWORK_ERROR",
			err:      NewCLIError(ErrCodeNetworkError, "Connection failed"),
			wantExit: ExitNetworkError,
		},
		{
			name:     "standard error returns ExitGeneralError",
			err:      errors.New("some standard error"),
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
		name     string
		got      int
		expected int
	}{
		{"ExitSuccess", ExitSuccess, 0},
		{"ExitGeneralError", ExitGeneralError, 1},
		{"ExitInvalidArgs", ExitInvalidArgs, 2},
		{"ExitAuthError", ExitAuthError, 3},
		{"ExitNetworkError", ExitNetworkError, 4},
		{"ExitRateLimited", ExitRateLimited, 5},
		{"ExitNotFound", ExitNotFound, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestErrorCodeConstants(t *testing.T) {
	// Verify error codes match PRD specification
	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"ErrCodeAuthRequired", ErrCodeAuthRequired, "AUTH_REQUIRED"},
		{"ErrCodeAuthExpired", ErrCodeAuthExpired, "AUTH_EXPIRED"},
		{"ErrCodeNotFound", ErrCodeNotFound, "NOT_FOUND"},
		{"ErrCodeRateLimited", ErrCodeRateLimited, "RATE_LIMITED"},
		{"ErrCodeInvalidInput", ErrCodeInvalidInput, "INVALID_INPUT"},
		{"ErrCodePurchaseFailed", ErrCodePurchaseFailed, "PURCHASE_FAILED"},
		{"ErrCodeNetworkError", ErrCodeNetworkError, "NETWORK_ERROR"},
		{"ErrCodeAmazonError", ErrCodeAmazonError, "AMAZON_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %s, want %s", tt.name, tt.got, tt.expected)
			}
		})
	}
}
