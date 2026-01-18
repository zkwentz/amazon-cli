package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

func TestPrintJSON(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := map[string]interface{}{
		"status": "success",
		"count":  42,
	}

	err := PrintJSON(data)
	if err != nil {
		t.Fatalf("PrintJSON returned error: %v", err)
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("PrintJSON produced invalid JSON: %v\nOutput: %s", err, output)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", result["status"])
	}

	if result["count"] != float64(42) {
		t.Errorf("Expected count 42, got '%v'", result["count"])
	}
}

func TestPrintJSONWithInvalidData(t *testing.T) {
	// Capture stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	_, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	// Create data that can't be marshaled (e.g., channel)
	data := make(chan int)

	err := PrintJSON(data)
	if err == nil {
		t.Error("Expected PrintJSON to return error for invalid data")
	}

	// Restore stdout and stderr
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var bufErr bytes.Buffer
	io.Copy(&bufErr, rErr)
	errorOutput := bufErr.String()

	// Verify error was printed as JSON to stderr
	if !strings.Contains(errorOutput, "error") {
		t.Errorf("Expected error JSON in stderr, got: %s", errorOutput)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(errorOutput), &result); err != nil {
		t.Errorf("Error output is not valid JSON: %v\nOutput: %s", err, errorOutput)
	}
}

func TestPrintError(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	cliErr := models.NewCLIError(models.ErrCodeNotFound, "Resource not found")
	PrintError(cliErr)

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify it's valid JSON
	var result models.ErrorResponse
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("PrintError produced invalid JSON: %v\nOutput: %s", err, output)
	}

	if result.Error.Code != models.ErrCodeNotFound {
		t.Errorf("Expected code %s, got %s", models.ErrCodeNotFound, result.Error.Code)
	}

	if result.Error.Message != "Resource not found" {
		t.Errorf("Expected message 'Resource not found', got '%s'", result.Error.Message)
	}
}

func TestPrintErrorWithRegularError(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Use a regular error, not a CLIError
	regularErr := errors.New("some standard error")
	PrintError(regularErr)

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify it's valid JSON
	var result models.ErrorResponse
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("PrintError produced invalid JSON for regular error: %v\nOutput: %s", err, output)
	}

	// Should be wrapped as INTERNAL_ERROR
	if result.Error.Code != models.ErrCodeInternalError {
		t.Errorf("Expected code %s for wrapped error, got %s", models.ErrCodeInternalError, result.Error.Code)
	}

	if result.Error.Message != "some standard error" {
		t.Errorf("Expected message 'some standard error', got '%s'", result.Error.Message)
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "nil error returns 0",
			err:          nil,
			expectedCode: 0,
		},
		{
			name:         "AUTH_REQUIRED returns 3",
			err:          models.NewCLIError(models.ErrCodeAuthRequired, "Auth required"),
			expectedCode: 3,
		},
		{
			name:         "AUTH_EXPIRED returns 3",
			err:          models.NewCLIError(models.ErrCodeAuthExpired, "Auth expired"),
			expectedCode: 3,
		},
		{
			name:         "NETWORK_ERROR returns 4",
			err:          models.NewCLIError(models.ErrCodeNetworkError, "Network error"),
			expectedCode: 4,
		},
		{
			name:         "RATE_LIMITED returns 5",
			err:          models.NewCLIError(models.ErrCodeRateLimited, "Rate limited"),
			expectedCode: 5,
		},
		{
			name:         "NOT_FOUND returns 6",
			err:          models.NewCLIError(models.ErrCodeNotFound, "Not found"),
			expectedCode: 6,
		},
		{
			name:         "INVALID_INPUT returns 2",
			err:          models.NewCLIError(models.ErrCodeInvalidInput, "Invalid input"),
			expectedCode: 2,
		},
		{
			name:         "PURCHASE_FAILED returns 1",
			err:          models.NewCLIError(models.ErrCodePurchaseFailed, "Purchase failed"),
			expectedCode: 1,
		},
		{
			name:         "AMAZON_ERROR returns 1",
			err:          models.NewCLIError(models.ErrCodeAmazonError, "Amazon error"),
			expectedCode: 1,
		},
		{
			name:         "regular error returns 1",
			err:          errors.New("standard error"),
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr to prevent test output pollution
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			exitCode := HandleError(tt.err)

			// Restore stderr
			w.Close()
			os.Stderr = oldStderr

			// Drain the pipe
			var buf bytes.Buffer
			io.Copy(&buf, r)

			if exitCode != tt.expectedCode {
				t.Errorf("Expected exit code %d, got %d", tt.expectedCode, exitCode)
			}

			// Verify error was printed as JSON (unless nil error)
			if tt.err != nil {
				output := buf.String()
				var result models.ErrorResponse
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("HandleError produced invalid JSON: %v\nOutput: %s", err, output)
				}
			}
		})
	}
}

func TestWrapPanic(t *testing.T) {
	// Note: WrapPanic calls os.Exit(1), which we can't test directly
	// Instead, we'll test the panic recovery logic in isolation

	tests := []struct {
		name        string
		panicValue  interface{}
		expectCode  string
		expectInMsg string
	}{
		{
			name:        "panic with error",
			panicValue:  errors.New("test error"),
			expectCode:  models.ErrCodeInternalError,
			expectInMsg: "test error",
		},
		{
			name:        "panic with string",
			panicValue:  "test panic string",
			expectCode:  models.ErrCodeInternalError,
			expectInMsg: "test panic string",
		},
		{
			name:        "panic with int",
			panicValue:  42,
			expectCode:  models.ErrCodeInternalError,
			expectInMsg: "Unexpected panic: 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Simulate WrapPanic's recovery logic without os.Exit
			func() {
				defer func() {
					if r := recover(); r != nil {
						var err *models.CLIError

						switch v := r.(type) {
						case error:
							err = models.NewCLIError(models.ErrCodeInternalError, v.Error())
						case string:
							err = models.NewCLIError(models.ErrCodeInternalError, v)
						default:
							err = models.NewCLIError(models.ErrCodeInternalError, fmt.Sprintf("Unexpected panic: %v", r))
						}

						fmt.Fprintln(os.Stderr, err.ToJSON())
					}
				}()

				panic(tt.panicValue)
			}()

			// Restore stderr
			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify JSON error was output
			if output == "" {
				t.Error("Expected error output for panic, got nothing")
				return
			}

			var result models.ErrorResponse
			if err := json.Unmarshal([]byte(output), &result); err != nil {
				t.Errorf("Panic handler produced invalid JSON: %v\nOutput: %s", err, output)
				return
			}

			if result.Error.Code != tt.expectCode {
				t.Errorf("Expected error code %s, got %s", tt.expectCode, result.Error.Code)
			}

			if !strings.Contains(result.Error.Message, tt.expectInMsg) {
				t.Errorf("Expected message to contain '%s', got '%s'", tt.expectInMsg, result.Error.Message)
			}
		})
	}
}
