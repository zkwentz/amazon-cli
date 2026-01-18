package models

import (
	"encoding/json"
	"testing"
)

func TestCLIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		cliError *CLIError
		want     string
	}{
		{
			name: "error without details",
			cliError: &CLIError{
				Code:    AuthRequired,
				Message: "Not logged in",
				Details: make(map[string]interface{}),
			},
			want: "AUTH_REQUIRED: Not logged in",
		},
		{
			name: "error with details",
			cliError: &CLIError{
				Code:    InvalidInput,
				Message: "Invalid ASIN format",
				Details: map[string]interface{}{
					"asin": "INVALID123",
				},
			},
			want: "INVALID_INPUT: Invalid ASIN format (details: map[asin:INVALID123])",
		},
		{
			name: "network error",
			cliError: &CLIError{
				Code:    NetworkError,
				Message: "Failed to connect to Amazon",
				Details: make(map[string]interface{}),
			},
			want: "NETWORK_ERROR: Failed to connect to Amazon",
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
	code := AuthExpired
	message := "Token has expired"

	err := NewCLIError(code, message)

	if err.Code != code {
		t.Errorf("NewCLIError() Code = %v, want %v", err.Code, code)
	}
	if err.Message != message {
		t.Errorf("NewCLIError() Message = %v, want %v", err.Message, message)
	}
	if err.Details == nil {
		t.Error("NewCLIError() Details should not be nil")
	}
	if len(err.Details) != 0 {
		t.Errorf("NewCLIError() Details length = %v, want 0", len(err.Details))
	}
}

func TestNewCLIErrorWithDetails(t *testing.T) {
	code := NotFound
	message := "Order not found"
	details := map[string]interface{}{
		"order_id": "123-4567890-1234567",
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
	if orderID, ok := err.Details["order_id"].(string); !ok || orderID != "123-4567890-1234567" {
		t.Errorf("NewCLIErrorWithDetails() Details[order_id] = %v, want 123-4567890-1234567", err.Details["order_id"])
	}
}

func TestCLIError_WithDetails(t *testing.T) {
	err := NewCLIError(RateLimited, "Too many requests")

	// Add single detail
	err.WithDetails("retry_after", 60)

	if retryAfter, ok := err.Details["retry_after"].(int); !ok || retryAfter != 60 {
		t.Errorf("WithDetails() retry_after = %v, want 60", err.Details["retry_after"])
	}

	// Add another detail
	err.WithDetails("request_id", "req_12345")

	if requestID, ok := err.Details["request_id"].(string); !ok || requestID != "req_12345" {
		t.Errorf("WithDetails() request_id = %v, want req_12345", err.Details["request_id"])
	}

	// Verify both details exist
	if len(err.Details) != 2 {
		t.Errorf("WithDetails() Details length = %v, want 2", len(err.Details))
	}
}

func TestCLIError_WithDetails_NilDetails(t *testing.T) {
	// Test that WithDetails initializes Details if nil
	err := &CLIError{
		Code:    PurchaseFailed,
		Message: "Payment failed",
		Details: nil,
	}

	err.WithDetails("reason", "card_declined")

	if err.Details == nil {
		t.Error("WithDetails() should initialize nil Details")
	}
	if reason, ok := err.Details["reason"].(string); !ok || reason != "card_declined" {
		t.Errorf("WithDetails() reason = %v, want card_declined", err.Details["reason"])
	}
}

func TestCLIError_JSONMarshaling(t *testing.T) {
	err := NewCLIErrorWithDetails(
		AmazonError,
		"Amazon returned an error",
		map[string]interface{}{
			"status_code": 500,
			"response":    "Internal Server Error",
		},
	)

	// Marshal to JSON
	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal CLIError to JSON: %v", marshalErr)
	}

	// Unmarshal back
	var unmarshaled CLIError
	if unmarshalErr := json.Unmarshal(jsonData, &unmarshaled); unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal CLIError from JSON: %v", unmarshalErr)
	}

	// Verify fields
	if unmarshaled.Code != err.Code {
		t.Errorf("Unmarshaled Code = %v, want %v", unmarshaled.Code, err.Code)
	}
	if unmarshaled.Message != err.Message {
		t.Errorf("Unmarshaled Message = %v, want %v", unmarshaled.Message, err.Message)
	}
	if statusCode, ok := unmarshaled.Details["status_code"].(float64); !ok || int(statusCode) != 500 {
		t.Errorf("Unmarshaled Details[status_code] = %v, want 500", unmarshaled.Details["status_code"])
	}
}

func TestErrorCodeConstants(t *testing.T) {
	// Test that all error codes are defined with expected values
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"AuthRequired", AuthRequired, "AUTH_REQUIRED"},
		{"AuthExpired", AuthExpired, "AUTH_EXPIRED"},
		{"NotFound", NotFound, "NOT_FOUND"},
		{"RateLimited", RateLimited, "RATE_LIMITED"},
		{"InvalidInput", InvalidInput, "INVALID_INPUT"},
		{"PurchaseFailed", PurchaseFailed, "PURCHASE_FAILED"},
		{"NetworkError", NetworkError, "NETWORK_ERROR"},
		{"AmazonError", AmazonError, "AMAZON_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Error code %s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestCLIError_ErrorInterface(t *testing.T) {
	// Verify CLIError implements error interface
	var _ error = &CLIError{}
	var _ error = (*CLIError)(nil)
}

func TestCLIError_ChainableWithDetails(t *testing.T) {
	// Test that WithDetails returns the error for chaining
	err := NewCLIError(InvalidInput, "Invalid parameters").
		WithDetails("field", "asin").
		WithDetails("value", "INVALID").
		WithDetails("expected_format", "10 alphanumeric characters")

	if len(err.Details) != 3 {
		t.Errorf("Chained WithDetails() Details length = %v, want 3", len(err.Details))
	}

	// Verify all details are present
	if field, ok := err.Details["field"].(string); !ok || field != "asin" {
		t.Errorf("Details[field] = %v, want asin", err.Details["field"])
	}
	if value, ok := err.Details["value"].(string); !ok || value != "INVALID" {
		t.Errorf("Details[value] = %v, want INVALID", err.Details["value"])
	}
	if format, ok := err.Details["expected_format"].(string); !ok || format != "10 alphanumeric characters" {
		t.Errorf("Details[expected_format] = %v, want '10 alphanumeric characters'", err.Details["expected_format"])
	}
}
