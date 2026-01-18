package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

func TestNewPrinter(t *testing.T) {
	printer := NewPrinter("json", false, true)

	if printer.format != FormatJSON {
		t.Errorf("Expected format JSON, got %s", printer.format)
	}

	if printer.quiet != false {
		t.Error("Expected quiet to be false")
	}

	if printer.verbose != true {
		t.Error("Expected verbose to be true")
	}
}

func TestPrinter_Print_JSON(t *testing.T) {
	var buf bytes.Buffer
	printer := NewPrinterWithWriter("json", false, false, &buf)

	data := map[string]string{
		"test": "value",
	}

	err := printer.Print(data)
	if err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["test"] != "value" {
		t.Errorf("Expected test='value', got '%s'", result["test"])
	}
}

func TestPrinter_Print_Quiet(t *testing.T) {
	var buf bytes.Buffer
	printer := NewPrinterWithWriter("json", true, false, &buf)

	data := map[string]string{
		"test": "value",
	}

	err := printer.Print(data)
	if err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	if buf.Len() > 0 {
		t.Error("Expected no output in quiet mode, got output")
	}
}

func TestPrinter_PrintError_WithoutVerbose(t *testing.T) {
	var buf bytes.Buffer
	printer := NewPrinterWithWriter("json", false, false, &buf)

	cliErr := models.NewCLIError(models.ErrInvalidInput, "invalid ASIN format")
	cliErr.WithDetails("asin", "BADFORMAT")

	err := printer.PrintError(cliErr)
	if err != nil {
		t.Fatalf("PrintError failed: %v", err)
	}

	// Trim the trailing newline that the encoder adds
	output := bytes.TrimSpace(buf.Bytes())

	var result models.ErrorResponse
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result.Error == nil {
		t.Fatal("Expected error in output")
	}

	if result.Error.Code != models.ErrInvalidInput {
		t.Errorf("Expected code %s, got %s", models.ErrInvalidInput, result.Error.Code)
	}

	if result.Error.Message != "invalid ASIN format" {
		t.Errorf("Expected message 'invalid ASIN format', got '%s'", result.Error.Message)
	}

	if result.Error.Details["asin"] != "BADFORMAT" {
		t.Errorf("Expected asin detail 'BADFORMAT', got '%v'", result.Error.Details["asin"])
	}

	// Should not have debug info when verbose is false
	if result.Error.Debug != nil {
		t.Error("Expected no debug info when verbose is false")
	}
}

func TestPrinter_PrintError_WithVerbose(t *testing.T) {
	var buf bytes.Buffer
	printer := NewPrinterWithWriter("json", false, true, &buf)

	cause := errors.New("connection timeout")
	cliErr := models.NewCLIErrorWithCause(models.ErrNetworkError, "network request failed", cause)
	cliErr.WithStackTrace("goroutine 1 [running]:\nmain.test()")

	err := printer.PrintError(cliErr)
	if err != nil {
		t.Fatalf("PrintError failed: %v", err)
	}

	// Trim the trailing newline that the encoder adds
	output := bytes.TrimSpace(buf.Bytes())

	var result models.ErrorResponse
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result.Error == nil {
		t.Fatal("Expected error in output")
	}

	if result.Error.Code != models.ErrNetworkError {
		t.Errorf("Expected code %s, got %s", models.ErrNetworkError, result.Error.Code)
	}

	if result.Error.Message != "network request failed" {
		t.Errorf("Expected message 'network request failed', got '%s'", result.Error.Message)
	}

	// Should have debug info when verbose is true
	if result.Error.Debug == nil {
		t.Fatal("Expected debug info when verbose is true")
	}

	if result.Error.Debug.Cause != "connection timeout" {
		t.Errorf("Expected cause 'connection timeout', got '%s'", result.Error.Debug.Cause)
	}

	if !strings.Contains(result.Error.Debug.StackTrace, "goroutine 1") {
		t.Errorf("Expected stack trace to contain 'goroutine 1', got '%s'", result.Error.Debug.StackTrace)
	}
}

func TestPrinter_PrintError_StandardError(t *testing.T) {
	var buf bytes.Buffer
	printer := NewPrinterWithWriter("json", false, false, &buf)

	// Test with a standard error (not CLIError)
	standardErr := errors.New("standard error message")

	err := printer.PrintError(standardErr)
	if err != nil {
		t.Fatalf("PrintError failed: %v", err)
	}

	// Trim the trailing newline that the encoder adds
	output := bytes.TrimSpace(buf.Bytes())

	var result models.ErrorResponse
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result.Error == nil {
		t.Fatal("Expected error in output")
	}

	// Standard errors should be wrapped as AMAZON_ERROR
	if result.Error.Code != models.ErrAmazonError {
		t.Errorf("Expected code %s, got %s", models.ErrAmazonError, result.Error.Code)
	}

	if result.Error.Message != "standard error message" {
		t.Errorf("Expected message 'standard error message', got '%s'", result.Error.Message)
	}
}

func TestPrinter_PrintError_WithVerbose_NoCause(t *testing.T) {
	var buf bytes.Buffer
	printer := NewPrinterWithWriter("json", false, true, &buf)

	// Error without cause or stack trace
	cliErr := models.NewCLIError(models.ErrNotFound, "resource not found")

	err := printer.PrintError(cliErr)
	if err != nil {
		t.Fatalf("PrintError failed: %v", err)
	}

	// Trim the trailing newline that the encoder adds
	output := bytes.TrimSpace(buf.Bytes())

	var result models.ErrorResponse
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Should not have debug info if there's no cause or stack trace
	if result.Error.Debug != nil {
		t.Error("Expected no debug info when there's no cause or stack trace")
	}
}

func TestPrinter_GetVerbose(t *testing.T) {
	printer := NewPrinter("json", false, true)
	if !printer.GetVerbose() {
		t.Error("Expected GetVerbose() to return true")
	}

	printer2 := NewPrinter("json", false, false)
	if printer2.GetVerbose() {
		t.Error("Expected GetVerbose() to return false")
	}
}

func TestPrinter_GetQuiet(t *testing.T) {
	printer := NewPrinter("json", true, false)
	if !printer.GetQuiet() {
		t.Error("Expected GetQuiet() to return true")
	}

	printer2 := NewPrinter("json", false, false)
	if printer2.GetQuiet() {
		t.Error("Expected GetQuiet() to return false")
	}
}

func TestPrinter_PrintError_Nil(t *testing.T) {
	var buf bytes.Buffer
	printer := NewPrinterWithWriter("json", false, false, &buf)

	err := printer.PrintError(nil)
	if err != nil {
		t.Fatalf("PrintError with nil should not fail: %v", err)
	}

	if buf.Len() > 0 {
		t.Error("Expected no output for nil error")
	}
}
