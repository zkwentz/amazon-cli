package models

import (
	"errors"
	"testing"
)

func TestNewCLIError(t *testing.T) {
	err := NewCLIError(ErrInvalidInput, "test message")

	if err.Code != ErrInvalidInput {
		t.Errorf("Expected code %s, got %s", ErrInvalidInput, err.Code)
	}

	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message)
	}

	if err.Details == nil {
		t.Error("Expected Details map to be initialized")
	}
}

func TestNewCLIErrorWithCause(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewCLIErrorWithCause(ErrNetworkError, "network failed", cause)

	if err.Code != ErrNetworkError {
		t.Errorf("Expected code %s, got %s", ErrNetworkError, err.Code)
	}

	if err.Message != "network failed" {
		t.Errorf("Expected message 'network failed', got '%s'", err.Message)
	}

	if err.Cause() != cause {
		t.Errorf("Expected cause to be %v, got %v", cause, err.Cause())
	}
}

func TestCLIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *CLIError
		expected string
	}{
		{
			name:     "error without cause",
			err:      NewCLIError(ErrNotFound, "resource not found"),
			expected: "NOT_FOUND: resource not found",
		},
		{
			name:     "error with cause",
			err:      NewCLIErrorWithCause(ErrAmazonError, "API failed", errors.New("connection timeout")),
			expected: "AMAZON_ERROR: API failed (caused by: connection timeout)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Expected error string '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestCLIError_WithDetails(t *testing.T) {
	err := NewCLIError(ErrInvalidInput, "invalid ASIN")
	err.WithDetails("asin", "INVALID123")
	err.WithDetails("reason", "invalid format")

	if len(err.Details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(err.Details))
	}

	if err.Details["asin"] != "INVALID123" {
		t.Errorf("Expected asin detail 'INVALID123', got '%v'", err.Details["asin"])
	}

	if err.Details["reason"] != "invalid format" {
		t.Errorf("Expected reason detail 'invalid format', got '%v'", err.Details["reason"])
	}
}

func TestCLIError_WithStackTrace(t *testing.T) {
	err := NewCLIError(ErrAmazonError, "test error")
	trace := "goroutine 1 [running]:\nmain.test()"
	err.WithStackTrace(trace)

	if err.StackTrace() != trace {
		t.Errorf("Expected stack trace '%s', got '%s'", trace, err.StackTrace())
	}
}

func TestCLIError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewCLIErrorWithCause(ErrNetworkError, "network failed", cause)

	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Errorf("Expected unwrapped error to be %v, got %v", cause, unwrapped)
	}
}

func TestErrorCodes(t *testing.T) {
	// Test that all error codes are defined correctly
	expectedCodes := map[string]string{
		"ErrAuthRequired":   "AUTH_REQUIRED",
		"ErrAuthExpired":    "AUTH_EXPIRED",
		"ErrNotFound":       "NOT_FOUND",
		"ErrRateLimited":    "RATE_LIMITED",
		"ErrInvalidInput":   "INVALID_INPUT",
		"ErrPurchaseFailed": "PURCHASE_FAILED",
		"ErrNetworkError":   "NETWORK_ERROR",
		"ErrAmazonError":    "AMAZON_ERROR",
	}

	actualCodes := map[string]string{
		"ErrAuthRequired":   ErrAuthRequired,
		"ErrAuthExpired":    ErrAuthExpired,
		"ErrNotFound":       ErrNotFound,
		"ErrRateLimited":    ErrRateLimited,
		"ErrInvalidInput":   ErrInvalidInput,
		"ErrPurchaseFailed": ErrPurchaseFailed,
		"ErrNetworkError":   ErrNetworkError,
		"ErrAmazonError":    ErrAmazonError,
	}

	for name, expected := range expectedCodes {
		if actual := actualCodes[name]; actual != expected {
			t.Errorf("Expected %s to be '%s', got '%s'", name, expected, actual)
		}
	}
}
