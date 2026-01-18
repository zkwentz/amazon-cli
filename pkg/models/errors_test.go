package models

import (
	"errors"
	"testing"
)

func TestCLIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		message  string
		expected string
	}{
		{
			name:     "auth required error",
			code:     ErrCodeAuthRequired,
			message:  "Please authenticate",
			expected: "AUTH_REQUIRED: Please authenticate",
		},
		{
			name:     "rate limited error",
			code:     ErrCodeRateLimited,
			message:  "Too many requests",
			expected: "RATE_LIMITED: Too many requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCLIError(tt.code, tt.message, ExitGeneralError)
			if err.Error() != tt.expected {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expected)
			}
		})
	}
}

func TestCLIError_WithDetails(t *testing.T) {
	err := NewCLIError(ErrCodeNotFound, "Resource not found", ExitNotFound)
	err = err.WithDetails("resource_id", "12345")
	err = err.WithDetails("resource_type", "order")

	if len(err.Details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(err.Details))
	}

	if err.Details["resource_id"] != "12345" {
		t.Errorf("Expected resource_id=12345, got %v", err.Details["resource_id"])
	}

	if err.Details["resource_type"] != "order" {
		t.Errorf("Expected resource_type=order, got %v", err.Details["resource_type"])
	}
}

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "nil error returns success",
			err:      nil,
			expected: ExitSuccess,
		},
		{
			name:     "auth error returns auth exit code",
			err:      NewAuthRequiredError(""),
			expected: ExitAuthError,
		},
		{
			name:     "network error returns network exit code",
			err:      NewNetworkError(""),
			expected: ExitNetworkError,
		},
		{
			name:     "not found error returns not found exit code",
			err:      NewNotFoundError(""),
			expected: ExitNotFound,
		},
		{
			name:     "rate limited error returns rate limited exit code",
			err:      NewRateLimitedError(""),
			expected: ExitRateLimited,
		},
		{
			name:     "invalid input error returns invalid args exit code",
			err:      NewInvalidInputError(""),
			expected: ExitInvalidArgs,
		},
		{
			name:     "generic error returns general error exit code",
			err:      errors.New("generic error"),
			expected: ExitGeneralError,
		},
		{
			name:     "amazon error returns general error exit code",
			err:      NewAmazonError(""),
			expected: ExitGeneralError,
		},
		{
			name:     "purchase failed error returns general error exit code",
			err:      NewPurchaseFailedError(""),
			expected: ExitGeneralError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := GetExitCode(tt.err)
			if exitCode != tt.expected {
				t.Errorf("GetExitCode() = %v, want %v", exitCode, tt.expected)
			}
		})
	}
}

func TestNewAuthRequiredError(t *testing.T) {
	tests := []struct {
		name            string
		message         string
		expectedCode    string
		expectedExitCode int
	}{
		{
			name:             "custom message",
			message:          "Custom auth message",
			expectedCode:     ErrCodeAuthRequired,
			expectedExitCode: ExitAuthError,
		},
		{
			name:             "default message",
			message:          "",
			expectedCode:     ErrCodeAuthRequired,
			expectedExitCode: ExitAuthError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAuthRequiredError(tt.message)
			if err.Code != tt.expectedCode {
				t.Errorf("Code = %v, want %v", err.Code, tt.expectedCode)
			}
			if err.ExitCode != tt.expectedExitCode {
				t.Errorf("ExitCode = %v, want %v", err.ExitCode, tt.expectedExitCode)
			}
			if tt.message != "" && err.Message != tt.message {
				t.Errorf("Message = %v, want %v", err.Message, tt.message)
			}
			if tt.message == "" && err.Message == "" {
				t.Error("Expected default message, got empty string")
			}
		})
	}
}

func TestNewAuthExpiredError(t *testing.T) {
	err := NewAuthExpiredError("Token expired")
	if err.Code != ErrCodeAuthExpired {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeAuthExpired)
	}
	if err.ExitCode != ExitAuthError {
		t.Errorf("ExitCode = %v, want %v", err.ExitCode, ExitAuthError)
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("Order not found")
	if err.Code != ErrCodeNotFound {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeNotFound)
	}
	if err.ExitCode != ExitNotFound {
		t.Errorf("ExitCode = %v, want %v", err.ExitCode, ExitNotFound)
	}
}

func TestNewRateLimitedError(t *testing.T) {
	err := NewRateLimitedError("Too many requests")
	if err.Code != ErrCodeRateLimited {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeRateLimited)
	}
	if err.ExitCode != ExitRateLimited {
		t.Errorf("ExitCode = %v, want %v", err.ExitCode, ExitRateLimited)
	}
}

func TestNewInvalidInputError(t *testing.T) {
	err := NewInvalidInputError("Invalid ASIN")
	if err.Code != ErrCodeInvalidInput {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeInvalidInput)
	}
	if err.ExitCode != ExitInvalidArgs {
		t.Errorf("ExitCode = %v, want %v", err.ExitCode, ExitInvalidArgs)
	}
}

func TestNewNetworkError(t *testing.T) {
	err := NewNetworkError("Connection failed")
	if err.Code != ErrCodeNetworkError {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeNetworkError)
	}
	if err.ExitCode != ExitNetworkError {
		t.Errorf("ExitCode = %v, want %v", err.ExitCode, ExitNetworkError)
	}
}

func TestNewAmazonError(t *testing.T) {
	err := NewAmazonError("Amazon API error")
	if err.Code != ErrCodeAmazonError {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeAmazonError)
	}
	if err.ExitCode != ExitGeneralError {
		t.Errorf("ExitCode = %v, want %v", err.ExitCode, ExitGeneralError)
	}
}

func TestNewPurchaseFailedError(t *testing.T) {
	err := NewPurchaseFailedError("Payment declined")
	if err.Code != ErrCodePurchaseFailed {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodePurchaseFailed)
	}
	if err.ExitCode != ExitGeneralError {
		t.Errorf("ExitCode = %v, want %v", err.ExitCode, ExitGeneralError)
	}
}

func TestExitCodeConstants(t *testing.T) {
	// Verify exit code constants match PRD specification
	tests := []struct {
		name     string
		code     int
		expected int
	}{
		{"Success", ExitSuccess, 0},
		{"General error", ExitGeneralError, 1},
		{"Invalid args", ExitInvalidArgs, 2},
		{"Auth error", ExitAuthError, 3},
		{"Network error", ExitNetworkError, 4},
		{"Rate limited", ExitRateLimited, 5},
		{"Not found", ExitNotFound, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.code, tt.expected)
			}
		})
	}
}
