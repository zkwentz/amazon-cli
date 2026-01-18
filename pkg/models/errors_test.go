package models

import (
	"encoding/json"
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
		t.Errorf("Expected Details to be initialized")
	}

	if len(err.Details) != 0 {
		t.Errorf("Expected Details to be empty initially")
	}
}

func TestCLIError_Error(t *testing.T) {
	err := NewCLIError(ErrCodeNetworkError, "Connection timeout")
	expected := "[NETWORK_ERROR] Connection timeout"

	if err.Error() != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
	}
}

func TestCLIError_WithDetails(t *testing.T) {
	err := NewCLIError(ErrCodeRateLimited, "Too many requests")
	details := map[string]interface{}{
		"retry_after": 60,
		"endpoint":    "/orders",
	}

	result := err.WithDetails(details)

	// Verify it returns the same error for chaining
	if result != err {
		t.Errorf("Expected WithDetails to return the same error instance")
	}

	// Verify details are set
	if retryAfter, ok := err.Details["retry_after"].(int); !ok || retryAfter != 60 {
		t.Errorf("Expected retry_after to be 60")
	}

	if endpoint, ok := err.Details["endpoint"].(string); !ok || endpoint != "/orders" {
		t.Errorf("Expected endpoint to be '/orders'")
	}
}

func TestCLIError_JSONMarshaling(t *testing.T) {
	err := NewCLIError(ErrCodeInvalidInput, "Invalid ASIN format")
	err.WithDetails(map[string]interface{}{
		"field": "asin",
		"value": "INVALID",
	})

	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal CLIError: %v", marshalErr)
	}

	// Unmarshal back to verify structure
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify fields
	if code := result["code"].(string); code != ErrCodeInvalidInput {
		t.Errorf("Expected code %s, got %s", ErrCodeInvalidInput, code)
	}

	if msg := result["message"].(string); msg != "Invalid ASIN format" {
		t.Errorf("Expected message 'Invalid ASIN format', got %s", msg)
	}

	details := result["details"].(map[string]interface{})
	if field := details["field"].(string); field != "asin" {
		t.Errorf("Expected field 'asin', got %s", field)
	}
}

func TestErrorCodes_Constants(t *testing.T) {
	// Verify all error codes are defined correctly
	expectedCodes := map[string]string{
		"ErrCodeAuthRequired":   "AUTH_REQUIRED",
		"ErrCodeAuthExpired":    "AUTH_EXPIRED",
		"ErrCodeNotFound":       "NOT_FOUND",
		"ErrCodeRateLimited":    "RATE_LIMITED",
		"ErrCodeInvalidInput":   "INVALID_INPUT",
		"ErrCodePurchaseFailed": "PURCHASE_FAILED",
		"ErrCodeNetworkError":   "NETWORK_ERROR",
		"ErrCodeAmazonError":    "AMAZON_ERROR",
	}

	actualCodes := map[string]string{
		"ErrCodeAuthRequired":   ErrCodeAuthRequired,
		"ErrCodeAuthExpired":    ErrCodeAuthExpired,
		"ErrCodeNotFound":       ErrCodeNotFound,
		"ErrCodeRateLimited":    ErrCodeRateLimited,
		"ErrCodeInvalidInput":   ErrCodeInvalidInput,
		"ErrCodePurchaseFailed": ErrCodePurchaseFailed,
		"ErrCodeNetworkError":   ErrCodeNetworkError,
		"ErrCodeAmazonError":    ErrCodeAmazonError,
	}

	for name, expectedValue := range expectedCodes {
		if actualValue := actualCodes[name]; actualValue != expectedValue {
			t.Errorf("Expected %s to be %s, got %s", name, expectedValue, actualValue)
		}
	}
}

func TestCLIError_AsError(t *testing.T) {
	// Verify CLIError implements error interface
	var err error = NewCLIError(ErrCodeNotFound, "Resource not found")

	if err == nil {
		t.Errorf("Expected error to be non-nil")
	}

	if err.Error() == "" {
		t.Errorf("Expected error message to be non-empty")
	}
}

func TestCLIError_EmptyDetails(t *testing.T) {
	err := NewCLIError(ErrCodeAmazonError, "Service unavailable")

	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal CLIError: %v", marshalErr)
	}

	// Verify empty details are marshaled as empty object
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	details := result["details"].(map[string]interface{})
	if len(details) != 0 {
		t.Errorf("Expected empty details object, got %d items", len(details))
	}
}
