package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		quiet          bool
		expectedFormat OutputFormat
	}{
		{"JSON format", "json", false, JSON},
		{"Table format", "table", false, Table},
		{"Raw format", "raw", false, Raw},
		{"Default format", "unknown", false, JSON},
		{"Quiet mode", "json", true, JSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPrinter(tt.format, tt.quiet)
			if p.format != tt.expectedFormat {
				t.Errorf("expected format %v, got %v", tt.expectedFormat, p.format)
			}
			if p.quiet != tt.quiet {
				t.Errorf("expected quiet %v, got %v", tt.quiet, p.quiet)
			}
		})
	}
}

func TestPrint_JSON(t *testing.T) {
	p := NewPrinter("json", false)

	testData := map[string]interface{}{
		"name":  "Test Product",
		"price": 29.99,
		"id":    12345,
	}

	output := captureOutput(func() {
		err := p.Print(testData)
		if err != nil {
			t.Errorf("Print returned error: %v", err)
		}
	})

	// Verify output is valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	// Verify data matches
	if parsed["name"] != "Test Product" {
		t.Errorf("expected name 'Test Product', got %v", parsed["name"])
	}
	if parsed["price"] != 29.99 {
		t.Errorf("expected price 29.99, got %v", parsed["price"])
	}
}

func TestPrint_Quiet(t *testing.T) {
	p := NewPrinter("json", true)

	testData := map[string]string{"test": "data"}

	output := captureOutput(func() {
		err := p.Print(testData)
		if err != nil {
			t.Errorf("Print returned error: %v", err)
		}
	})

	// Quiet mode should produce no output
	if output != "" {
		t.Errorf("expected no output in quiet mode, got: %s", output)
	}
}

func TestPrint_Raw(t *testing.T) {
	p := NewPrinter("raw", false)

	testData := "Simple string data"

	output := captureOutput(func() {
		err := p.Print(testData)
		if err != nil {
			t.Errorf("Print returned error: %v", err)
		}
	})

	// Raw mode should output the data as-is with newline
	expected := "Simple string data\n"
	if output != expected {
		t.Errorf("expected output '%s', got '%s'", expected, output)
	}
}

func TestPrintError(t *testing.T) {
	p := NewPrinter("json", false)

	tests := []struct {
		name         string
		err          error
		expectedCode string
	}{
		{
			name:         "CLIError",
			err:          models.NewCLIError(models.AuthRequired, "Not authenticated"),
			expectedCode: models.AuthRequired,
		},
		{
			name:         "Generic error",
			err:          io.EOF,
			expectedCode: models.AmazonError,
		},
		{
			name:         "Nil error",
			err:          nil,
			expectedCode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				err := p.PrintError(tt.err)
				if err != nil {
					t.Errorf("PrintError returned error: %v", err)
				}
			})

			if tt.err == nil {
				if output != "" {
					t.Errorf("expected no output for nil error, got: %s", output)
				}
				return
			}

			// Verify output is valid JSON
			var errResponse models.ErrorResponse
			if err := json.Unmarshal([]byte(output), &errResponse); err != nil {
				t.Errorf("Output is not valid JSON: %v\nOutput: %s", err, output)
			}

			// Verify error code
			if errResponse.Error.Code != tt.expectedCode {
				t.Errorf("expected error code '%s', got '%s'", tt.expectedCode, errResponse.Error.Code)
			}
		})
	}
}

func TestPrintError_WithDetails(t *testing.T) {
	p := NewPrinter("json", false)

	details := map[string]interface{}{
		"status_code": 404,
		"url":         "https://example.com",
	}
	err := models.NewCLIErrorWithDetails(models.NotFound, "Resource not found", details)

	output := captureOutput(func() {
		printErr := p.PrintError(err)
		if printErr != nil {
			t.Errorf("PrintError returned error: %v", printErr)
		}
	})

	// Verify output is valid JSON
	var errResponse models.ErrorResponse
	if err := json.Unmarshal([]byte(output), &errResponse); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	// Verify error details
	if errResponse.Error.Details["status_code"] != float64(404) {
		t.Errorf("expected status_code 404, got %v", errResponse.Error.Details["status_code"])
	}
}

func TestPrintJSON_InvalidData(t *testing.T) {
	p := NewPrinter("json", false)

	// Create data that cannot be marshaled to JSON (channel)
	invalidData := make(chan int)

	err := p.Print(invalidData)
	if err == nil {
		t.Error("expected error for invalid JSON data, got nil")
	}
}

func TestPrint_Table(t *testing.T) {
	p := NewPrinter("table", false)

	testData := []map[string]interface{}{
		{"name": "Product 1", "price": 10.99},
		{"name": "Product 2", "price": 20.99},
	}

	output := captureOutput(func() {
		err := p.Print(testData)
		if err != nil {
			t.Errorf("Print returned error: %v", err)
		}
	})

	// Basic check that output contains the data
	if output == "" {
		t.Error("expected table output, got empty string")
	}
}

func TestPrint_TableSingleObject(t *testing.T) {
	p := NewPrinter("table", false)

	testData := map[string]interface{}{
		"name":  "Test Product",
		"price": 29.99,
	}

	output := captureOutput(func() {
		err := p.Print(testData)
		if err != nil {
			t.Errorf("Print returned error: %v", err)
		}
	})

	// Basic check that output contains the data
	if output == "" {
		t.Error("expected table output, got empty string")
	}
}
