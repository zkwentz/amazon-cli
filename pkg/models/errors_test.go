package models

import (
	"errors"
	"strings"
	"testing"
)

func TestCLIError_Error(t *testing.T) {
	testCases := []struct {
		name     string
		err      *CLIError
		expected string
	}{
		{
			name: "error without details",
			err: &CLIError{
				Code:    ErrCodeAuthRequired,
				Message: "Authentication required",
			},
			expected: "AUTH_REQUIRED: Authentication required",
		},
		{
			name: "error with details",
			err: &CLIError{
				Code:    ErrCodeCaptchaRequired,
				Message: "CAPTCHA required",
				Details: map[string]interface{}{
					"url": "https://www.amazon.com",
				},
			},
			expected: "CAPTCHA_REQUIRED: CAPTCHA required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.err.Error()
			if !strings.Contains(result, tc.expected) {
				t.Errorf("Expected error string to contain %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestCLIError_ImplementsError(t *testing.T) {
	var err error = &CLIError{
		Code:    ErrCodeNetworkError,
		Message: "Network error occurred",
	}

	if err == nil {
		t.Error("CLIError should implement error interface")
	}
}

func TestNewCLIError(t *testing.T) {
	code := "TEST_CODE"
	message := "Test message"
	details := map[string]interface{}{
		"key": "value",
	}

	err := NewCLIError(code, message, details)

	if err.Code != code {
		t.Errorf("Expected code %s, got %s", code, err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message %s, got %s", message, err.Message)
	}
	if err.Details["key"] != "value" {
		t.Error("Expected details to be set correctly")
	}
}

func TestNewAuthRequiredError(t *testing.T) {
	err := NewAuthRequiredError()

	if err.Code != ErrCodeAuthRequired {
		t.Errorf("Expected code %s, got %s", ErrCodeAuthRequired, err.Code)
	}
	if !strings.Contains(err.Message, "auth login") {
		t.Error("Expected message to mention auth login command")
	}
}

func TestNewAuthExpiredError(t *testing.T) {
	err := NewAuthExpiredError()

	if err.Code != ErrCodeAuthExpired {
		t.Errorf("Expected code %s, got %s", ErrCodeAuthExpired, err.Code)
	}
	if !strings.Contains(err.Message, "expired") {
		t.Error("Expected message to mention expiration")
	}
}

func TestNewCaptchaRequiredError(t *testing.T) {
	details := map[string]interface{}{
		"url":     "https://www.amazon.com/cart",
		"snippet": "<html>CAPTCHA</html>",
	}

	err := NewCaptchaRequiredError(details)

	if err.Code != ErrCodeCaptchaRequired {
		t.Errorf("Expected code %s, got %s", ErrCodeCaptchaRequired, err.Code)
	}
	if !strings.Contains(err.Message, "CAPTCHA") {
		t.Error("Expected message to mention CAPTCHA")
	}
	if err.Details == nil {
		t.Error("Expected details to be set")
	}
	if err.Details["url"] != details["url"] {
		t.Error("Expected details to be passed through")
	}
}

func TestNewLoginRequiredError(t *testing.T) {
	details := map[string]interface{}{
		"url": "https://www.amazon.com/signin",
	}

	err := NewLoginRequiredError(details)

	if err.Code != ErrCodeLoginRequired {
		t.Errorf("Expected code %s, got %s", ErrCodeLoginRequired, err.Code)
	}
	if !strings.Contains(err.Message, "auth") || !strings.Contains(err.Message, "login") {
		t.Error("Expected message to mention authentication/login")
	}
	if err.Details == nil {
		t.Error("Expected details to be set")
	}
}

func TestNewHTMLResponseError(t *testing.T) {
	details := map[string]interface{}{
		"content_type": "text/html",
	}

	err := NewHTMLResponseError(details)

	if err.Code != ErrCodeHTMLResponse {
		t.Errorf("Expected code %s, got %s", ErrCodeHTMLResponse, err.Code)
	}
	if !strings.Contains(err.Message, "HTML") {
		t.Error("Expected message to mention HTML")
	}
	if err.Details == nil {
		t.Error("Expected details to be set")
	}
}

func TestNewRateLimitedError(t *testing.T) {
	err := NewRateLimitedError()

	if err.Code != ErrCodeRateLimited {
		t.Errorf("Expected code %s, got %s", ErrCodeRateLimited, err.Code)
	}
	if !strings.Contains(err.Message, "Rate limited") {
		t.Error("Expected message to mention rate limiting")
	}
}

func TestNewNetworkError(t *testing.T) {
	originalErr := errors.New("connection timeout")
	err := NewNetworkError(originalErr)

	if err.Code != ErrCodeNetworkError {
		t.Errorf("Expected code %s, got %s", ErrCodeNetworkError, err.Code)
	}
	if !strings.Contains(err.Message, originalErr.Error()) {
		t.Errorf("Expected message to contain original error: %s", originalErr.Error())
	}
}

func TestNewAmazonError(t *testing.T) {
	message := "Invalid request"
	details := map[string]interface{}{
		"status_code": 400,
	}

	err := NewAmazonError(message, details)

	if err.Code != ErrCodeAmazonError {
		t.Errorf("Expected code %s, got %s", ErrCodeAmazonError, err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message %s, got %s", message, err.Message)
	}
	if err.Details == nil {
		t.Error("Expected details to be set")
	}
}

func TestNewNotFoundError(t *testing.T) {
	resource := "Order"
	err := NewNotFoundError(resource)

	if err.Code != ErrCodeNotFound {
		t.Errorf("Expected code %s, got %s", ErrCodeNotFound, err.Code)
	}
	if !strings.Contains(err.Message, resource) {
		t.Errorf("Expected message to contain resource name: %s", resource)
	}
}

func TestNewInvalidInputError(t *testing.T) {
	message := "ASIN must be 10 characters"
	err := NewInvalidInputError(message)

	if err.Code != ErrCodeInvalidInput {
		t.Errorf("Expected code %s, got %s", ErrCodeInvalidInput, err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message %s, got %s", message, err.Message)
	}
}

func TestNewPurchaseFailedError(t *testing.T) {
	message := "Payment declined"
	details := map[string]interface{}{
		"reason": "insufficient_funds",
	}

	err := NewPurchaseFailedError(message, details)

	if err.Code != ErrCodePurchaseFailed {
		t.Errorf("Expected code %s, got %s", ErrCodePurchaseFailed, err.Code)
	}
	if err.Message != message {
		t.Errorf("Expected message %s, got %s", message, err.Message)
	}
	if err.Details == nil {
		t.Error("Expected details to be set")
	}
}

func TestErrorCodes_AreUnique(t *testing.T) {
	codes := []string{
		ErrCodeAuthRequired,
		ErrCodeAuthExpired,
		ErrCodeNotFound,
		ErrCodeRateLimited,
		ErrCodeNetworkError,
		ErrCodeAmazonError,
		ErrCodeInvalidInput,
		ErrCodePurchaseFailed,
		ErrCodeCaptchaRequired,
		ErrCodeLoginRequired,
		ErrCodeHTMLResponse,
	}

	seen := make(map[string]bool)
	for _, code := range codes {
		if seen[code] {
			t.Errorf("Duplicate error code found: %s", code)
		}
		seen[code] = true
	}
}

func TestErrorCodes_AllUppercase(t *testing.T) {
	codes := []string{
		ErrCodeAuthRequired,
		ErrCodeAuthExpired,
		ErrCodeNotFound,
		ErrCodeRateLimited,
		ErrCodeNetworkError,
		ErrCodeAmazonError,
		ErrCodeInvalidInput,
		ErrCodePurchaseFailed,
		ErrCodeCaptchaRequired,
		ErrCodeLoginRequired,
		ErrCodeHTMLResponse,
	}

	for _, code := range codes {
		if code != strings.ToUpper(code) {
			t.Errorf("Error code should be uppercase: %s", code)
		}
	}
}
