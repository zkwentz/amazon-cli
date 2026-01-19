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
			_, _ = io.Copy(&buf, r)

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
	_, _ = io.Copy(&buf, r)

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
	_, _ = io.Copy(&buf, r)

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
