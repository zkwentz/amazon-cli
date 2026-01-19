package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		message        string
		details        map[string]interface{}
		expectedCode   string
		expectedMsg    string
		expectedDetail map[string]interface{}
	}{
		{
			name:           "basic error with nil details",
			code:           "TEST_ERROR",
			message:        "This is a test error",
			details:        nil,
			expectedCode:   "TEST_ERROR",
			expectedMsg:    "This is a test error",
			expectedDetail: map[string]interface{}{},
		},
		{
			name:    "error with empty details",
			code:    "EMPTY_DETAILS",
			message: "Error with empty details",
			details: map[string]interface{}{},
			expectedCode:   "EMPTY_DETAILS",
			expectedMsg:    "Error with empty details",
			expectedDetail: map[string]interface{}{},
		},
		{
			name:    "error with details",
			code:    "VALIDATION_ERROR",
			message: "Validation failed",
			details: map[string]interface{}{
				"field": "email",
				"reason": "invalid format",
			},
			expectedCode: "VALIDATION_ERROR",
			expectedMsg:  "Validation failed",
			expectedDetail: map[string]interface{}{
				"field": "email",
				"reason": "invalid format",
			},
		},
		{
			name:    "error with nested details",
			code:    "COMPLEX_ERROR",
			message: "Complex error occurred",
			details: map[string]interface{}{
				"user": map[string]interface{}{
					"id": 123,
					"name": "test",
				},
				"metadata": map[string]interface{}{
					"timestamp": "2024-01-01T00:00:00Z",
				},
			},
			expectedCode: "COMPLEX_ERROR",
			expectedMsg:  "Complex error occurred",
			expectedDetail: map[string]interface{}{
				"user": map[string]interface{}{
					"id": float64(123), // JSON unmarshals numbers as float64
					"name": "test",
				},
				"metadata": map[string]interface{}{
					"timestamp": "2024-01-01T00:00:00Z",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Call Error function
			err := Error(tt.code, tt.message, tt.details)
			if err != nil {
				t.Fatalf("Error() returned error: %v", err)
			}

			// Restore stderr and read captured output
			w.Close()
			os.Stderr = old
			var buf bytes.Buffer
			io.Copy(&buf, r)

			// Parse JSON output
			var result map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
				t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
			}

			// Verify structure
			errorObj, ok := result["error"].(map[string]interface{})
			if !ok {
				t.Fatal("Output missing 'error' object")
			}

			// Verify code
			if code, ok := errorObj["code"].(string); !ok || code != tt.expectedCode {
				t.Errorf("Expected code %q, got %q", tt.expectedCode, code)
			}

			// Verify message
			if msg, ok := errorObj["message"].(string); !ok || msg != tt.expectedMsg {
				t.Errorf("Expected message %q, got %q", tt.expectedMsg, msg)
			}

			// Verify details
			details, ok := errorObj["details"].(map[string]interface{})
			if !ok {
				t.Fatal("Output missing 'details' object")
			}

			// Deep comparison of details
			detailsJSON, _ := json.Marshal(details)
			expectedJSON, _ := json.Marshal(tt.expectedDetail)
			if string(detailsJSON) != string(expectedJSON) {
				t.Errorf("Expected details %s, got %s", string(expectedJSON), string(detailsJSON))
			}
		})
	}
}

func TestError_OutputsToStderr(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Call Error function
	err := Error("TEST_CODE", "test message", nil)
	if err != nil {
		t.Fatalf("Error() returned error: %v", err)
	}

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify output was captured (meaning it went to stderr)
	if buf.Len() == 0 {
		t.Fatal("Expected output to stderr, got none")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}
}

func TestError_JSONFormat(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	details := map[string]interface{}{
		"field": "test",
		"value": 42,
	}

	err := Error("FORMAT_TEST", "Testing JSON format", details)
	if err != nil {
		t.Fatalf("Error() returned error: %v", err)
	}

	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify exact JSON structure
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify top-level structure has only "error" key
	if len(result) != 1 {
		t.Errorf("Expected 1 top-level key, got %d", len(result))
	}

	errorObj := result["error"].(map[string]interface{})
	if len(errorObj) != 3 {
		t.Errorf("Expected 3 keys in error object (code, message, details), got %d", len(errorObj))
	}

	// Verify all expected keys exist
	if _, ok := errorObj["code"]; !ok {
		t.Error("Missing 'code' field in error object")
	}
	if _, ok := errorObj["message"]; !ok {
		t.Error("Missing 'message' field in error object")
	}
	if _, ok := errorObj["details"]; !ok {
		t.Error("Missing 'details' field in error object")
	}
}

func TestJSON_OutputsValidJSON(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test data
	testData := map[string]interface{}{
		"id":     123,
		"name":   "Test Product",
		"price":  29.99,
		"active": true,
		"tags":   []string{"electronics", "gadgets"},
	}

	// Call JSON function
	err := JSON(testData)
	if err != nil {
		t.Fatalf("JSON() returned error: %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify output is not empty
	if buf.Len() == 0 {
		t.Fatal("Expected JSON output, got none")
	}

	// Verify JSON is parseable
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify the parsed data matches original
	if result["id"].(float64) != 123 {
		t.Errorf("Expected id 123, got %v", result["id"])
	}
	if result["name"].(string) != "Test Product" {
		t.Errorf("Expected name 'Test Product', got %v", result["name"])
	}
	if result["price"].(float64) != 29.99 {
		t.Errorf("Expected price 29.99, got %v", result["price"])
	}
	if result["active"].(bool) != true {
		t.Errorf("Expected active true, got %v", result["active"])
	}
}

func TestError_OutputsErrorSchema(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Call Error function
	err := Error("SCHEMA_TEST", "Testing error schema", map[string]interface{}{
		"field": "value",
	})
	if err != nil {
		t.Fatalf("Error() returned error: %v", err)
	}

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify JSON is parseable
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify the schema structure: top-level should have "error" key
	if _, ok := result["error"]; !ok {
		t.Fatal("Output missing top-level 'error' key")
	}

	// Verify error object is a map
	errorObj, ok := result["error"].(map[string]interface{})
	if !ok {
		t.Fatal("'error' key should contain an object")
	}

	// Verify error object has exactly 3 keys: code, message, details
	expectedKeys := map[string]bool{"code": false, "message": false, "details": false}
	for key := range errorObj {
		if _, exists := expectedKeys[key]; exists {
			expectedKeys[key] = true
		} else {
			t.Errorf("Unexpected key in error object: %s", key)
		}
	}

	for key, found := range expectedKeys {
		if !found {
			t.Errorf("Missing required key in error object: %s", key)
		}
	}

	// Verify types
	if _, ok := errorObj["code"].(string); !ok {
		t.Error("'code' field should be a string")
	}
	if _, ok := errorObj["message"].(string); !ok {
		t.Error("'message' field should be a string")
	}
	if _, ok := errorObj["details"].(map[string]interface{}); !ok {
		t.Error("'details' field should be an object")
	}
}

func TestError_IncludesAllFields(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		message string
		details map[string]interface{}
	}{
		{
			name:    "all fields with details",
			code:    "TEST_CODE",
			message: "Test message",
			details: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
		},
		{
			name:    "all fields with nil details",
			code:    "ANOTHER_CODE",
			message: "Another message",
			details: nil,
		},
		{
			name:    "all fields with empty details",
			code:    "EMPTY_DETAILS_CODE",
			message: "Empty details message",
			details: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Call Error function
			err := Error(tt.code, tt.message, tt.details)
			if err != nil {
				t.Fatalf("Error() returned error: %v", err)
			}

			// Restore stderr and read captured output
			w.Close()
			os.Stderr = old
			var buf bytes.Buffer
			io.Copy(&buf, r)

			// Parse JSON output
			var result map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
				t.Fatalf("Failed to parse JSON output: %v", err)
			}

			// Verify error object exists
			errorObj, ok := result["error"].(map[string]interface{})
			if !ok {
				t.Fatal("Output missing 'error' object")
			}

			// Verify 'code' field exists and matches
			code, ok := errorObj["code"].(string)
			if !ok {
				t.Fatal("Missing or invalid 'code' field")
			}
			if code != tt.code {
				t.Errorf("Expected code %q, got %q", tt.code, code)
			}

			// Verify 'message' field exists and matches
			message, ok := errorObj["message"].(string)
			if !ok {
				t.Fatal("Missing or invalid 'message' field")
			}
			if message != tt.message {
				t.Errorf("Expected message %q, got %q", tt.message, message)
			}

			// Verify 'details' field exists
			details, ok := errorObj["details"].(map[string]interface{})
			if !ok {
				t.Fatal("Missing or invalid 'details' field")
			}

			// If details were nil or empty, verify the output has empty object
			if tt.details == nil || len(tt.details) == 0 {
				if len(details) != 0 {
					t.Errorf("Expected empty details object, got %v", details)
				}
			} else {
				// Verify details match
				if len(details) != len(tt.details) {
					t.Errorf("Expected %d detail entries, got %d", len(tt.details), len(details))
				}
			}
		})
	}
}
