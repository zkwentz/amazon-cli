package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestPrintError_WithCLIError(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	printer := NewPrinter("json", false)
	testErr := models.NewCLIError(models.ErrCodeAuthExpired, "Authentication token has expired")
	testErr.WithDetails(map[string]interface{}{
		"expired_at": "2024-01-20T12:00:00Z",
	})

	returnedErr := printer.PrintError(testErr)

	// Close writer and restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify the original error is returned
	if returnedErr != testErr {
		t.Errorf("Expected PrintError to return the original error")
	}

	// Parse the JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Verify structure
	errorObj, ok := result["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'error' field in output")
	}

	// Verify code
	if code := errorObj["code"].(string); code != models.ErrCodeAuthExpired {
		t.Errorf("Expected code %s, got %s", models.ErrCodeAuthExpired, code)
	}

	// Verify message
	if msg := errorObj["message"].(string); msg != "Authentication token has expired" {
		t.Errorf("Expected message 'Authentication token has expired', got %s", msg)
	}

	// Verify details
	details, ok := errorObj["details"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'details' field in error object")
	}

	if expiredAt := details["expired_at"].(string); expiredAt != "2024-01-20T12:00:00Z" {
		t.Errorf("Expected expired_at '2024-01-20T12:00:00Z', got %s", expiredAt)
	}
}

func TestPrintError_WithStandardError(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	printer := NewPrinter("json", false)
	testErr := errors.New("something went wrong")

	returnedErr := printer.PrintError(testErr)

	// Close writer and restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify the original error is returned
	if returnedErr != testErr {
		t.Errorf("Expected PrintError to return the original error")
	}

	// Parse the JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Verify structure
	errorObj, ok := result["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'error' field in output")
	}

	// Standard errors should be wrapped with AMAZON_ERROR code
	if code := errorObj["code"].(string); code != models.ErrCodeAmazonError {
		t.Errorf("Expected code %s for standard error, got %s", models.ErrCodeAmazonError, code)
	}

	// Verify message
	if msg := errorObj["message"].(string); msg != "something went wrong" {
		t.Errorf("Expected message 'something went wrong', got %s", msg)
	}

	// Verify details exists (even if empty)
	if _, ok := errorObj["details"]; !ok {
		t.Errorf("Expected 'details' field in error object")
	}
}

func TestPrintError_WithNilError(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	printer := NewPrinter("json", false)
	returnedErr := printer.PrintError(nil)

	// Close writer and restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify nil is returned
	if returnedErr != nil {
		t.Errorf("Expected PrintError to return nil for nil error")
	}

	// Verify no output was produced
	if strings.TrimSpace(output) != "" {
		t.Errorf("Expected no output for nil error, got: %s", output)
	}
}

func TestPrintError_AllErrorCodes(t *testing.T) {
	errorCodes := []struct {
		code    string
		message string
	}{
		{models.ErrCodeAuthRequired, "Not logged in"},
		{models.ErrCodeAuthExpired, "Token expired"},
		{models.ErrCodeNotFound, "Resource not found"},
		{models.ErrCodeRateLimited, "Too many requests"},
		{models.ErrCodeInvalidInput, "Invalid command input"},
		{models.ErrCodePurchaseFailed, "Purchase could not be completed"},
		{models.ErrCodeNetworkError, "Network connectivity issue"},
		{models.ErrCodeAmazonError, "Amazon returned an error"},
	}

	for _, tc := range errorCodes {
		t.Run(tc.code, func(t *testing.T) {
			// Capture stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			printer := NewPrinter("json", false)
			testErr := models.NewCLIError(tc.code, tc.message)
			printer.PrintError(testErr)

			// Close writer and restore stderr
			w.Close()
			os.Stderr = oldStderr

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Parse the JSON output
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(output), &result); err != nil {
				t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
			}

			// Verify error code
			errorObj := result["error"].(map[string]interface{})
			if code := errorObj["code"].(string); code != tc.code {
				t.Errorf("Expected code %s, got %s", tc.code, code)
			}

			// Verify message
			if msg := errorObj["message"].(string); msg != tc.message {
				t.Errorf("Expected message %s, got %s", tc.message, msg)
			}
		})
	}
}

func TestPrintError_JSONFormatting(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	printer := NewPrinter("json", false)
	testErr := models.NewCLIError(models.ErrCodeInvalidInput, "Invalid ASIN format")

	printer.PrintError(testErr)

	// Close writer and restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify the output is properly indented JSON
	if !strings.Contains(output, "{\n") {
		t.Errorf("Expected indented JSON output")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, output)
	}
}

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		format   string
		expected OutputFormat
	}{
		{"json", FormatJSON},
		{"table", FormatTable},
		{"raw", FormatRaw},
		{"invalid", FormatJSON}, // defaults to JSON
		{"", FormatJSON},         // defaults to JSON
	}

	for _, tc := range tests {
		t.Run(tc.format, func(t *testing.T) {
			printer := NewPrinter(tc.format, false)
			if printer.format != tc.expected {
				t.Errorf("Expected format %s, got %s", tc.expected, printer.format)
			}
		})
	}
}

func TestPrintError_OutputToStderr(t *testing.T) {
	// Capture both stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout = wOut
	os.Stderr = wErr

	printer := NewPrinter("json", false)
	testErr := models.NewCLIError(models.ErrCodeNetworkError, "Connection timeout")
	printer.PrintError(testErr)

	// Close writers and restore
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read captured output
	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

	stdoutOutput := bufOut.String()
	stderrOutput := bufErr.String()

	// Verify error output goes to stderr, not stdout
	if stdoutOutput != "" {
		t.Errorf("Expected no output to stdout, got: %s", stdoutOutput)
	}

	if stderrOutput == "" {
		t.Errorf("Expected output to stderr, got none")
	}

	// Verify stderr contains the error
	if !strings.Contains(stderrOutput, "NETWORK_ERROR") {
		t.Errorf("Expected stderr to contain error code")
	}
}
