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

func TestNewPrinter(t *testing.T) {
	testCases := []struct {
		format   string
		expected OutputFormat
	}{
		{"json", FormatJSON},
		{"JSON", FormatJSON},
		{"table", FormatTable},
		{"TABLE", FormatTable},
		{"raw", FormatRaw},
		{"RAW", FormatRaw},
		{"invalid", FormatJSON}, // defaults to JSON
		{"", FormatJSON},        // defaults to JSON
	}

	for _, tc := range testCases {
		p := NewPrinter(tc.format, false)
		if p.format != tc.expected {
			t.Errorf("Format '%s': expected %s, got %s", tc.format, tc.expected, p.format)
		}
	}
}

func TestNewPrinter_QuietMode(t *testing.T) {
	p := NewPrinter("json", true)
	if !p.quiet {
		t.Error("Expected quiet mode to be true")
	}

	p = NewPrinter("json", false)
	if p.quiet {
		t.Error("Expected quiet mode to be false")
	}
}

func TestPrinter_IsQuiet(t *testing.T) {
	p := NewPrinter("json", true)
	if !p.IsQuiet() {
		t.Error("Expected IsQuiet to return true")
	}

	p = NewPrinter("json", false)
	if p.IsQuiet() {
		t.Error("Expected IsQuiet to return false")
	}
}

func TestPrinter_GetFormat(t *testing.T) {
	p := NewPrinter("table", false)
	if p.GetFormat() != FormatTable {
		t.Errorf("Expected format %s, got %s", FormatTable, p.GetFormat())
	}
}

func TestPrinter_Print_JSON(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p := NewPrinter("json", false)
	data := map[string]interface{}{
		"status": "success",
		"count":  42,
	}

	err := p.Print(data)
	if err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got %v", result["status"])
	}

	// JSON numbers are float64
	if result["count"] != float64(42) {
		t.Errorf("Expected count 42, got %v", result["count"])
	}
}

func TestPrinter_Print_Quiet(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p := NewPrinter("json", true) // quiet mode
	data := map[string]string{"test": "data"}

	err := p.Print(data)
	if err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// In quiet mode, nothing should be printed
	if output != "" {
		t.Errorf("Expected no output in quiet mode, got: %s", output)
	}
}

func TestPrinter_PrintError_CLIError(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	p := NewPrinter("json", false)
	cliErr := models.NewCLIError(models.ErrCodeAuthExpired, "Token has expired").
		WithDetail("expiresAt", "2024-01-01")

	err := p.PrintError(cliErr)
	if err != nil {
		t.Fatalf("PrintError failed: %v", err)
	}

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify JSON error output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON error output: %v", err)
	}

	errorObj, ok := result["error"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected error object in output")
	}

	if errorObj["code"] != models.ErrCodeAuthExpired {
		t.Errorf("Expected code %s, got %v", models.ErrCodeAuthExpired, errorObj["code"])
	}

	if errorObj["message"] != "Token has expired" {
		t.Errorf("Expected message 'Token has expired', got %v", errorObj["message"])
	}
}

func TestPrinter_PrintError_StandardError(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	p := NewPrinter("json", false)
	stdErr := errors.New("standard error message")

	err := p.PrintError(stdErr)
	if err != nil {
		t.Fatalf("PrintError failed: %v", err)
	}

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify JSON error output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON error output: %v", err)
	}

	errorObj, ok := result["error"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected error object in output")
	}

	if errorObj["message"] != "standard error message" {
		t.Errorf("Expected message 'standard error message', got %v", errorObj["message"])
	}
}

func TestPrinter_PrintError_Nil(t *testing.T) {
	p := NewPrinter("json", false)
	err := p.PrintError(nil)
	if err != nil {
		t.Errorf("Expected nil error to be handled, got %v", err)
	}
}

func TestPrinter_PrintSuccess(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p := NewPrinter("json", false)
	err := p.PrintSuccess("Operation completed")
	if err != nil {
		t.Fatalf("PrintSuccess failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got %v", result["status"])
	}

	if result["message"] != "Operation completed" {
		t.Errorf("Expected message 'Operation completed', got %v", result["message"])
	}
}

func TestPrinter_PrintSuccess_Quiet(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p := NewPrinter("json", true) // quiet mode
	err := p.PrintSuccess("Should not print")
	if err != nil {
		t.Fatalf("PrintSuccess failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// In quiet mode, nothing should be printed
	if output != "" {
		t.Errorf("Expected no output in quiet mode, got: %s", output)
	}
}

func TestPrinter_Print_Raw(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	p := NewPrinter("raw", false)
	data := "raw string output"

	err := p.Print(data)
	if err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	if output != "raw string output" {
		t.Errorf("Expected 'raw string output', got '%s'", output)
	}
}

func TestPrinter_PrintError_TableFormat(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	p := NewPrinter("table", false)
	cliErr := models.NewCLIError(models.ErrCodeNotFound, "Resource not found")

	err := p.PrintError(cliErr)
	if err != nil {
		t.Fatalf("PrintError failed: %v", err)
	}

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// For table format, should print simplified error
	if !strings.Contains(output, "NOT_FOUND") {
		t.Errorf("Expected output to contain 'NOT_FOUND', got: %s", output)
	}

	if !strings.Contains(output, "Resource not found") {
		t.Errorf("Expected output to contain 'Resource not found', got: %s", output)
	}
}

func TestOutputFormat_Constants(t *testing.T) {
	if FormatJSON != "json" {
		t.Errorf("Expected FormatJSON to be 'json', got %s", FormatJSON)
	}

	if FormatTable != "table" {
		t.Errorf("Expected FormatTable to be 'table', got %s", FormatTable)
	}

	if FormatRaw != "raw" {
		t.Errorf("Expected FormatRaw to be 'raw', got %s", FormatRaw)
	}
}
