package models

import (
	"encoding/json"
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

	// Verify the error string is valid JSON
	if !json.Valid([]byte(errorStr)) {
		t.Errorf("Expected valid JSON from Error(), got: %s", errorStr)
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

// TestCLIError_ErrorReturnsJSON verifies that Error() returns valid JSON
func TestCLIError_ErrorReturnsJSON(t *testing.T) {
	tests := []struct {
		name     string
		cliError *CLIError
		wantCode string
		wantMsg  string
	}{
		{
			name: "simple error without details",
			cliError: &CLIError{
				Code:    ErrNotFound,
				Message: "Item not found",
				Details: map[string]interface{}{},
			},
			wantCode: ErrNotFound,
			wantMsg:  "Item not found",
		},
		{
			name: "error with string details",
			cliError: &CLIError{
				Code:    ErrInvalidInput,
				Message: "Invalid ASIN provided",
				Details: map[string]interface{}{
					"asin":  "INVALID123",
					"field": "product_id",
				},
			},
			wantCode: ErrInvalidInput,
			wantMsg:  "Invalid ASIN provided",
		},
		{
			name: "error with mixed type details",
			cliError: &CLIError{
				Code:    ErrRateLimited,
				Message: "Too many requests",
				Details: map[string]interface{}{
					"retry_after": 60,
					"limit":       100,
					"endpoint":    "/api/cart",
				},
			},
			wantCode: ErrRateLimited,
			wantMsg:  "Too many requests",
		},
		{
			name: "auth error",
			cliError: &CLIError{
				Code:    ErrAuthRequired,
				Message: "Authentication required to proceed",
				Details: map[string]interface{}{},
			},
			wantCode: ErrAuthRequired,
			wantMsg:  "Authentication required to proceed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the JSON error string
			errorString := tt.cliError.Error()

			// Verify it's valid JSON
			if !json.Valid([]byte(errorString)) {
				t.Fatalf("Error() did not return valid JSON: %s", errorString)
			}

			// Unmarshal and verify fields
			var unmarshaled CLIError
			err := json.Unmarshal([]byte(errorString), &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON from Error(): %v\nGot: %s", err, errorString)
			}

			// Verify the fields match
			if unmarshaled.Code != tt.wantCode {
				t.Errorf("Code mismatch: got %q, want %q", unmarshaled.Code, tt.wantCode)
			}
			if unmarshaled.Message != tt.wantMsg {
				t.Errorf("Message mismatch: got %q, want %q", unmarshaled.Message, tt.wantMsg)
			}

			// Verify details are preserved
			if len(tt.cliError.Details) > 0 {
				if len(unmarshaled.Details) != len(tt.cliError.Details) {
					t.Errorf("Details length mismatch: got %d, want %d",
						len(unmarshaled.Details), len(tt.cliError.Details))
				}
				for key := range tt.cliError.Details {
					if _, exists := unmarshaled.Details[key]; !exists {
						t.Errorf("Detail key %q not found in unmarshaled error", key)
					}
				}
			}
		})
	}
}

// TestCLIError_StructFieldsComplete verifies all required fields exist
func TestCLIError_StructFieldsComplete(t *testing.T) {
	// Test that CLIError has all three required fields
	err := &CLIError{
		Code:    "TEST_CODE",
		Message: "Test message",
		Details: map[string]interface{}{
			"key":   "value",
			"count": 42,
		},
	}

	if err.Code != "TEST_CODE" {
		t.Errorf("Code field not working: got %q, want %q", err.Code, "TEST_CODE")
	}
	if err.Message != "Test message" {
		t.Errorf("Message field not working: got %q, want %q", err.Message, "Test message")
	}
	if err.Details["key"] != "value" {
		t.Errorf("Details field not working: got %v, want %q", err.Details["key"], "value")
	}
	if err.Details["count"] != 42 {
		t.Errorf("Details field not working: got %v, want %d", err.Details["count"], 42)
	}
}

// TestCLIError_ImplementsErrorInterface verifies CLIError implements error interface
func TestCLIError_ImplementsErrorInterface(t *testing.T) {
	// Compile-time check that CLIError implements error interface
	var _ error = &CLIError{}

	// Runtime check
	err := NewCLIError(ErrNetworkError, "Connection timeout")
	errorString := err.Error()

	if errorString == "" {
		t.Error("Error() returned empty string")
	}

	// Verify the error string is valid JSON
	if !json.Valid([]byte(errorString)) {
		t.Errorf("Error() did not return valid JSON: %s", errorString)
	}
}

// TestCLIError_JSONMarshaling tests direct JSON marshaling
func TestCLIError_JSONMarshaling(t *testing.T) {
	err := &CLIError{
		Code:    ErrAmazonError,
		Message: "Service unavailable",
		Details: map[string]interface{}{
			"service": "cart",
			"status":  503,
		},
	}

	// Test marshaling
	jsonBytes, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal CLIError: %v", marshalErr)
	}

	// Test unmarshaling
	var unmarshaled CLIError
	unmarshalErr := json.Unmarshal(jsonBytes, &unmarshaled)
	if unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal CLIError: %v", unmarshalErr)
	}

	// Verify fields
	if unmarshaled.Code != err.Code {
		t.Errorf("Code after marshal/unmarshal: got %q, want %q", unmarshaled.Code, err.Code)
	}
	if unmarshaled.Message != err.Message {
		t.Errorf("Message after marshal/unmarshal: got %q, want %q", unmarshaled.Message, err.Message)
	}
	if unmarshaled.Details["service"] != err.Details["service"] {
		t.Errorf("Details[service] after marshal/unmarshal: got %v, want %v",
			unmarshaled.Details["service"], err.Details["service"])
	}
}
