package models

import (
	"testing"
)

// TestErrorCodesExist verifies that all required error codes are defined
func TestErrorCodesExist(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"AUTH_REQUIRED", ErrAuthRequired, "AUTH_REQUIRED"},
		{"AUTH_EXPIRED", ErrAuthExpired, "AUTH_EXPIRED"},
		{"NOT_FOUND", ErrNotFound, "NOT_FOUND"},
		{"RATE_LIMITED", ErrRateLimited, "RATE_LIMITED"},
		{"INVALID_INPUT", ErrInvalidInput, "INVALID_INPUT"},
		{"PURCHASE_FAILED", ErrPurchaseFailed, "PURCHASE_FAILED"},
		{"NETWORK_ERROR", ErrNetworkError, "NETWORK_ERROR"},
		{"AMAZON_ERROR", ErrAmazonError, "AMAZON_ERROR"},
		{"CAPTCHA_REQUIRED", ErrCaptchaRequired, "CAPTCHA_REQUIRED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.expected {
				t.Errorf("Expected error code %s to equal %s, got %s", tt.name, tt.expected, tt.code)
			}
		})
	}
}

// TestCLIErrorCreation tests creating and using CLI errors
func TestCLIErrorCreation(t *testing.T) {
	err := NewCLIError(ErrCaptchaRequired, "CAPTCHA verification required")

	if err.Code != ErrCaptchaRequired {
		t.Errorf("Expected error code %s, got %s", ErrCaptchaRequired, err.Code)
	}

	if err.Message != "CAPTCHA verification required" {
		t.Errorf("Expected message 'CAPTCHA verification required', got %s", err.Message)
	}

	// Error() now returns JSON formatted string
	errorStr := err.Error()
	if errorStr == "" {
		t.Error("Expected non-empty error string")
	}
}

// TestCLIErrorWithDetails tests adding details to CLI errors
func TestCLIErrorWithDetails(t *testing.T) {
	err := NewCLIError(ErrCaptchaRequired, "CAPTCHA verification required").
		WithDetails(map[string]interface{}{
			"url": "https://amazon.com/captcha",
		})

	if err.Details["url"] != "https://amazon.com/captcha" {
		t.Errorf("Expected detail 'url' to be 'https://amazon.com/captcha', got %v", err.Details["url"])
	}
}

// TestExitCodeForError tests the exit code mapping for all error codes
func TestExitCodeForError(t *testing.T) {
	tests := []struct {
		name     string
		errCode  string
		expected int
	}{
		{"AUTH_REQUIRED returns ExitAuthError", ErrAuthRequired, ExitAuthError},
		{"AUTH_EXPIRED returns ExitAuthError", ErrAuthExpired, ExitAuthError},
		{"NETWORK_ERROR returns ExitNetworkError", ErrNetworkError, ExitNetworkError},
		{"RATE_LIMITED returns ExitRateLimited", ErrRateLimited, ExitRateLimited},
		{"NOT_FOUND returns ExitNotFound", ErrNotFound, ExitNotFound},
		{"INVALID_INPUT returns ExitInvalidArgs", ErrInvalidInput, ExitInvalidArgs},
		{"PURCHASE_FAILED returns ExitGeneralError", ErrPurchaseFailed, ExitGeneralError},
		{"AMAZON_ERROR returns ExitGeneralError", ErrAmazonError, ExitGeneralError},
		{"CAPTCHA_REQUIRED returns ExitGeneralError", ErrCaptchaRequired, ExitGeneralError},
		{"Unknown error returns ExitGeneralError", "UNKNOWN_ERROR", ExitGeneralError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := ExitCodeForError(tt.errCode)
			if exitCode != tt.expected {
				t.Errorf("Expected exit code %d for %s, got %d", tt.expected, tt.errCode, exitCode)
			}
		})
	}
}
