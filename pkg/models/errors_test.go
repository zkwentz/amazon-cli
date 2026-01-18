package models

import (
	"testing"
)

func TestNewCLIError(t *testing.T) {
	err := NewCLIError(ErrCodeAuthRequired, "Authentication required")

	if err.Code != ErrCodeAuthRequired {
		t.Errorf("Expected code %s, got %s", ErrCodeAuthRequired, err.Code)
	}

	if err.Message != "Authentication required" {
		t.Errorf("Expected message 'Authentication required', got %s", err.Message)
	}

	if err.Details == nil {
		t.Error("Expected Details to be initialized")
	}

	if len(err.Details) != 0 {
		t.Errorf("Expected empty Details, got %d items", len(err.Details))
	}
}

func TestCLIError_Error(t *testing.T) {
	err := NewCLIError(ErrCodeNetworkError, "Connection failed")
	expected := "NETWORK_ERROR: Connection failed"

	if err.Error() != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
	}
}

func TestCLIError_WithDetails(t *testing.T) {
	details := map[string]interface{}{
		"url":        "https://amazon.com",
		"statusCode": 500,
	}

	err := NewCLIError(ErrCodeAmazonError, "Server error").WithDetails(details)

	if len(err.Details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(err.Details))
	}

	if err.Details["url"] != "https://amazon.com" {
		t.Errorf("Expected url detail, got %v", err.Details["url"])
	}

	if err.Details["statusCode"] != 500 {
		t.Errorf("Expected statusCode detail, got %v", err.Details["statusCode"])
	}
}

func TestCLIError_WithDetail(t *testing.T) {
	err := NewCLIError(ErrCodeRateLimited, "Too many requests").
		WithDetail("retryAfter", 60).
		WithDetail("requestId", "abc123")

	if len(err.Details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(err.Details))
	}

	if err.Details["retryAfter"] != 60 {
		t.Errorf("Expected retryAfter detail 60, got %v", err.Details["retryAfter"])
	}

	if err.Details["requestId"] != "abc123" {
		t.Errorf("Expected requestId detail 'abc123', got %v", err.Details["requestId"])
	}
}

func TestErrorCodes(t *testing.T) {
	testCases := []struct {
		code     string
		expected string
	}{
		{ErrCodeAuthRequired, "AUTH_REQUIRED"},
		{ErrCodeAuthExpired, "AUTH_EXPIRED"},
		{ErrCodeNotFound, "NOT_FOUND"},
		{ErrCodeRateLimited, "RATE_LIMITED"},
		{ErrCodeInvalidInput, "INVALID_INPUT"},
		{ErrCodePurchaseFailed, "PURCHASE_FAILED"},
		{ErrCodeNetworkError, "NETWORK_ERROR"},
		{ErrCodeAmazonError, "AMAZON_ERROR"},
	}

	for _, tc := range testCases {
		if tc.code != tc.expected {
			t.Errorf("Error code mismatch: expected %s, got %s", tc.expected, tc.code)
		}
	}
}
