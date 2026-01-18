package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		quiet          bool
		expectedFormat OutputFormat
		expectedQuiet  bool
	}{
		{
			name:           "json format",
			format:         "json",
			quiet:          false,
			expectedFormat: FormatJSON,
			expectedQuiet:  false,
		},
		{
			name:           "table format",
			format:         "table",
			quiet:          false,
			expectedFormat: FormatTable,
			expectedQuiet:  false,
		},
		{
			name:           "raw format",
			format:         "raw",
			quiet:          false,
			expectedFormat: FormatRaw,
			expectedQuiet:  false,
		},
		{
			name:           "quiet mode enabled",
			format:         "json",
			quiet:          true,
			expectedFormat: FormatJSON,
			expectedQuiet:  true,
		},
		{
			name:           "invalid format defaults to json",
			format:         "invalid",
			quiet:          false,
			expectedFormat: FormatJSON,
			expectedQuiet:  false,
		},
		{
			name:           "empty format defaults to json",
			format:         "",
			quiet:          false,
			expectedFormat: FormatJSON,
			expectedQuiet:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := NewPrinter(tt.format, tt.quiet)
			if printer == nil {
				t.Fatal("NewPrinter returned nil")
			}
			if printer.format != tt.expectedFormat {
				t.Errorf("expected format %v, got %v", tt.expectedFormat, printer.format)
			}
			if printer.quiet != tt.expectedQuiet {
				t.Errorf("expected quiet %v, got %v", tt.expectedQuiet, printer.quiet)
			}
			if printer.writer == nil {
				t.Error("writer should not be nil")
			}
		})
	}
}

func TestPrinter_PrintJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		quiet    bool
		expected string
	}{
		{
			name: "simple map",
			data: map[string]string{
				"key": "value",
			},
			quiet:    false,
			expected: "{\n  \"key\": \"value\"\n}\n",
		},
		{
			name: "nested structure",
			data: map[string]interface{}{
				"orders": []map[string]interface{}{
					{
						"order_id": "123",
						"total":    29.99,
					},
				},
			},
			quiet:    false,
			expected: "{\n  \"orders\": [\n    {\n      \"order_id\": \"123\",\n      \"total\": 29.99\n    }\n  ]\n}\n",
		},
		{
			name:     "quiet mode suppresses output",
			data:     map[string]string{"key": "value"},
			quiet:    true,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			printer := NewPrinter("json", tt.quiet)
			printer.SetWriter(buf)

			err := printer.Print(tt.data)
			if err != nil {
				t.Errorf("Print() error = %v", err)
				return
			}

			got := buf.String()
			if got != tt.expected {
				t.Errorf("expected output:\n%s\ngot:\n%s", tt.expected, got)
			}
		})
	}
}

func TestPrinter_PrintRaw(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected string
	}{
		{
			name:     "string data",
			data:     "hello world",
			expected: "hello world\n",
		},
		{
			name:     "integer data",
			data:     42,
			expected: "42\n",
		},
		{
			name:     "boolean data",
			data:     true,
			expected: "true\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			printer := NewPrinter("raw", false)
			printer.SetWriter(buf)

			err := printer.Print(tt.data)
			if err != nil {
				t.Errorf("Print() error = %v", err)
				return
			}

			got := buf.String()
			if got != tt.expected {
				t.Errorf("expected output:\n%s\ngot:\n%s", tt.expected, got)
			}
		})
	}
}

func TestPrinter_PrintError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		quiet       bool
		shouldPrint bool
	}{
		{
			name:        "error in normal mode",
			err:         errors.New("test error"),
			quiet:       false,
			shouldPrint: true,
		},
		{
			name:        "error in quiet mode still prints",
			err:         errors.New("test error"),
			quiet:       true,
			shouldPrint: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			printer := NewPrinter("json", tt.quiet)
			printer.SetWriter(buf)

			err := printer.PrintError(tt.err)
			if err != nil {
				t.Errorf("PrintError() error = %v", err)
				return
			}

			output := buf.String()
			if tt.shouldPrint && output == "" {
				t.Error("expected error output, got empty string")
			}

			if tt.shouldPrint {
				// Verify it's valid JSON
				var result map[string]interface{}
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("output is not valid JSON: %v", err)
				}

				// Verify structure
				if _, ok := result["error"]; !ok {
					t.Error("output missing 'error' key")
				}

				// Verify error message is present
				if !strings.Contains(output, tt.err.Error()) {
					t.Errorf("output doesn't contain error message: %s", tt.err.Error())
				}
			}
		})
	}
}

func TestPrinter_QuietMode(t *testing.T) {
	buf := &bytes.Buffer{}
	printer := NewPrinter("json", true)
	printer.SetWriter(buf)

	data := map[string]string{"test": "data"}
	err := printer.Print(data)
	if err != nil {
		t.Errorf("Print() error = %v", err)
	}

	// Quiet mode should suppress normal output
	if buf.String() != "" {
		t.Errorf("expected no output in quiet mode, got: %s", buf.String())
	}

	// But errors should still print
	buf.Reset()
	testErr := errors.New("test error")
	err = printer.PrintError(testErr)
	if err != nil {
		t.Errorf("PrintError() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("expected error output even in quiet mode")
	}
}

func TestPrinter_SetWriter(t *testing.T) {
	printer := NewPrinter("json", false)
	buf := &bytes.Buffer{}
	printer.SetWriter(buf)

	data := map[string]string{"test": "data"}
	err := printer.Print(data)
	if err != nil {
		t.Errorf("Print() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("expected output to custom writer")
	}
}

func TestOutputFormatConstants(t *testing.T) {
	if FormatJSON != "json" {
		t.Errorf("FormatJSON constant incorrect: %s", FormatJSON)
	}
	if FormatTable != "table" {
		t.Errorf("FormatTable constant incorrect: %s", FormatTable)
	}
	if FormatRaw != "raw" {
		t.Errorf("FormatRaw constant incorrect: %s", FormatRaw)
	}
}
