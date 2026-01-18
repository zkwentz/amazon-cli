package models

import (
	"testing"
)

func TestCLIError_Error(t *testing.T) {
	err := &CLIError{
		Code:    ErrCodeAuthRequired,
		Message: "Authentication required",
		Details: map[string]interface{}{
			"foo": "bar",
		},
	}

	expected := "[AUTH_REQUIRED] Authentication required"
	if err.Error() != expected {
		t.Errorf("expected %s, got %s", expected, err.Error())
	}
}

func TestNewCLIError(t *testing.T) {
	details := map[string]interface{}{
		"key": "value",
	}

	err := NewCLIError(ErrCodeAuthExpired, "Token expired", details)

	if err.Code != ErrCodeAuthExpired {
		t.Errorf("expected code %s, got %s", ErrCodeAuthExpired, err.Code)
	}

	if err.Message != "Token expired" {
		t.Errorf("expected message 'Token expired', got %s", err.Message)
	}

	if err.Details["key"] != "value" {
		t.Errorf("expected details to contain key=value")
	}
}

func TestNewCLIError_NilDetails(t *testing.T) {
	err := NewCLIError(ErrCodeNotFound, "Not found", nil)

	if err.Details == nil {
		t.Error("expected details to be initialized as empty map")
	}

	if len(err.Details) != 0 {
		t.Errorf("expected empty details map, got %d entries", len(err.Details))
	}
}

func TestErrorCodes(t *testing.T) {
	// Verify error code constants match PRD specification
	expectedCodes := map[string]string{
		"AUTH_REQUIRED":   ErrCodeAuthRequired,
		"AUTH_EXPIRED":    ErrCodeAuthExpired,
		"NOT_FOUND":       ErrCodeNotFound,
		"RATE_LIMITED":    ErrCodeRateLimited,
		"INVALID_INPUT":   ErrCodeInvalidInput,
		"PURCHASE_FAILED": ErrCodePurchaseFailed,
		"NETWORK_ERROR":   ErrCodeNetworkError,
		"AMAZON_ERROR":    ErrCodeAmazonError,
	}

	for expected, actual := range expectedCodes {
		if expected != actual {
			t.Errorf("expected error code %s, got %s", expected, actual)
		}
	}
}
