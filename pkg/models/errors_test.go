package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewCLIError(t *testing.T) {
	err := NewCLIError(ErrCodeInvalidInput, "Invalid ASIN provided")

	if err.Code != ErrCodeInvalidInput {
		t.Errorf("Expected code %s, got %s", ErrCodeInvalidInput, err.Code)
	}

	if err.Message != "Invalid ASIN provided" {
		t.Errorf("Expected message 'Invalid ASIN provided', got '%s'", err.Message)
	}

	if err.Details == nil {
		t.Error("Expected Details map to be initialized, got nil")
	}
}

func TestNewCLIErrorWithDetails(t *testing.T) {
	details := map[string]interface{}{
		"field": "asin",
		"value": "invalid123",
	}

	err := NewCLIErrorWithDetails(ErrCodeInvalidInput, "Invalid ASIN provided", details)

	if err.Code != ErrCodeInvalidInput {
		t.Errorf("Expected code %s, got %s", ErrCodeInvalidInput, err.Code)
	}

	if err.Message != "Invalid ASIN provided" {
		t.Errorf("Expected message 'Invalid ASIN provided', got '%s'", err.Message)
	}

	if err.Details["field"] != "asin" {
		t.Errorf("Expected details field 'asin', got '%v'", err.Details["field"])
	}

	if err.Details["value"] != "invalid123" {
		t.Errorf("Expected details value 'invalid123', got '%v'", err.Details["value"])
	}
}

func TestCLIErrorError(t *testing.T) {
	err := NewCLIError(ErrCodeNetworkError, "Connection timeout")

	errorString := err.Error()
	expected := "NETWORK_ERROR: Connection timeout"

	if errorString != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, errorString)
	}
}

func TestCLIErrorToJSON(t *testing.T) {
	err := NewCLIError(ErrCodeAuthRequired, "Authentication required")

	jsonOutput := err.ToJSON()

	// Verify it's valid JSON
	var result ErrorResponse
	if parseErr := json.Unmarshal([]byte(jsonOutput), &result); parseErr != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", parseErr, jsonOutput)
	}

	// Verify structure
	if result.Error == nil {
		t.Fatal("Expected error field to be present in JSON output")
	}

	if result.Error.Code != ErrCodeAuthRequired {
		t.Errorf("Expected code %s, got %s", ErrCodeAuthRequired, result.Error.Code)
	}

	if result.Error.Message != "Authentication required" {
		t.Errorf("Expected message 'Authentication required', got '%s'", result.Error.Message)
	}
}

func TestCLIErrorToJSONWithDetails(t *testing.T) {
	details := map[string]interface{}{
		"order_id": "123-456-789",
		"status":   404,
	}

	err := NewCLIErrorWithDetails(ErrCodeNotFound, "Order not found", details)

	jsonOutput := err.ToJSON()

	// Verify it's valid JSON
	var result ErrorResponse
	if parseErr := json.Unmarshal([]byte(jsonOutput), &result); parseErr != nil {
		t.Fatalf("Failed to parse JSON output: %v", parseErr)
	}

	// Verify details are present
	if result.Error.Details == nil {
		t.Fatal("Expected details field to be present")
	}

	if result.Error.Details["order_id"] != "123-456-789" {
		t.Errorf("Expected order_id '123-456-789', got '%v'", result.Error.Details["order_id"])
	}

	// Note: JSON unmarshaling converts numbers to float64
	if result.Error.Details["status"] != float64(404) {
		t.Errorf("Expected status 404, got '%v'", result.Error.Details["status"])
	}
}

func TestCLIErrorWithDetail(t *testing.T) {
	err := NewCLIError(ErrCodePurchaseFailed, "Payment declined")

	err.WithDetail("reason", "insufficient_funds")
	err.WithDetail("card_last4", "1234")

	if err.Details["reason"] != "insufficient_funds" {
		t.Errorf("Expected reason 'insufficient_funds', got '%v'", err.Details["reason"])
	}

	if err.Details["card_last4"] != "1234" {
		t.Errorf("Expected card_last4 '1234', got '%v'", err.Details["card_last4"])
	}
}

func TestCLIErrorToJSONNeverPanics(t *testing.T) {
	// Test with circular reference that might cause marshal issues
	// This ensures we have a fallback that never panics
	err := NewCLIError(ErrCodeInternalError, "Test error")

	// Create a circular reference (this would normally cause json.Marshal to fail)
	circular := make(map[string]interface{})
	circular["self"] = circular
	err.Details["circular"] = circular

	// This should not panic - it should return a fallback JSON
	jsonOutput := err.ToJSON()

	// Verify we got some JSON output (either normal or fallback)
	if jsonOutput == "" {
		t.Error("Expected JSON output, got empty string")
	}

	// Fallback JSON should at least contain the error code and message
	if !strings.Contains(jsonOutput, ErrCodeInternalError) {
		t.Errorf("Expected output to contain error code, got: %s", jsonOutput)
	}
}

func TestAllErrorCodes(t *testing.T) {
	errorCodes := []string{
		ErrCodeAuthRequired,
		ErrCodeAuthExpired,
		ErrCodeNotFound,
		ErrCodeRateLimited,
		ErrCodeInvalidInput,
		ErrCodePurchaseFailed,
		ErrCodeNetworkError,
		ErrCodeAmazonError,
		ErrCodeInternalError,
	}

	for _, code := range errorCodes {
		err := NewCLIError(code, "Test message")
		jsonOutput := err.ToJSON()

		// Verify each error code produces valid JSON
		var result ErrorResponse
		if parseErr := json.Unmarshal([]byte(jsonOutput), &result); parseErr != nil {
			t.Errorf("Error code %s produced invalid JSON: %v", code, parseErr)
		}

		if result.Error.Code != code {
			t.Errorf("Expected code %s in JSON, got %s", code, result.Error.Code)
		}
	}
}

func TestErrorResponseStructure(t *testing.T) {
	// Verify the JSON structure matches PRD specification
	err := NewCLIErrorWithDetails(
		ErrCodeAuthExpired,
		"Authentication token has expired. Run 'amazon-cli auth login' to re-authenticate.",
		map[string]interface{}{
			"expires_at": "2024-01-15T10:00:00Z",
		},
	)

	jsonOutput := err.ToJSON()

	// Parse and verify structure
	var result map[string]interface{}
	if parseErr := json.Unmarshal([]byte(jsonOutput), &result); parseErr != nil {
		t.Fatalf("Failed to parse JSON: %v", parseErr)
	}

	// Verify top-level "error" field
	if _, ok := result["error"]; !ok {
		t.Error("Expected top-level 'error' field in JSON output")
	}

	errorObj, ok := result["error"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'error' to be an object")
	}

	// Verify required fields
	if _, ok := errorObj["code"]; !ok {
		t.Error("Expected 'code' field in error object")
	}

	if _, ok := errorObj["message"]; !ok {
		t.Error("Expected 'message' field in error object")
	}

	// Verify details field when present
	if _, ok := errorObj["details"]; !ok {
		t.Error("Expected 'details' field in error object")
	}
}
