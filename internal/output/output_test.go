package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		name   string
		format string
		quiet  bool
	}{
		{
			name:   "JSON format, not quiet",
			format: "json",
			quiet:  false,
		},
		{
			name:   "Table format, quiet",
			format: "table",
			quiet:  true,
		},
		{
			name:   "Raw format, not quiet",
			format: "raw",
			quiet:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPrinter(tt.format, tt.quiet)
			if p == nil {
				t.Fatal("NewPrinter returned nil")
			}
			if p.format != OutputFormat(tt.format) {
				t.Errorf("expected format %s, got %s", tt.format, p.format)
			}
			if p.quiet != tt.quiet {
				t.Errorf("expected quiet %v, got %v", tt.quiet, p.quiet)
			}
		})
	}
}

func TestPrintJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected string
	}{
		{
			name: "Simple struct",
			data: struct {
				Name  string `json:"name"`
				Value int    `json:"value"`
			}{
				Name:  "test",
				Value: 42,
			},
			expected: `{
  "name": "test",
  "value": 42
}`,
		},
		{
			name: "Array of items",
			data: []string{"item1", "item2", "item3"},
			expected: `[
  "item1",
  "item2",
  "item3"
]`,
		},
		{
			name: "Map",
			data: map[string]interface{}{
				"status": "success",
				"count":  5,
			},
			expected: `{
  "count": 5,
  "status": "success"
}`,
		},
		{
			name:     "String",
			data:     "test string",
			expected: `"test string"`,
		},
		{
			name:     "Number",
			data:     123,
			expected: `123`,
		},
		{
			name:     "Boolean",
			data:     true,
			expected: `true`,
		},
		{
			name:     "Null",
			data:     nil,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			p := NewPrinterWithWriter("json", false, &buf)

			err := p.Print(tt.data)
			if err != nil {
				t.Fatalf("Print returned error: %v", err)
			}

			output := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if output != expected {
				t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", expected, output)
			}

			// Validate that output is valid JSON
			var parsed interface{}
			if err := json.Unmarshal([]byte(output), &parsed); err != nil {
				t.Errorf("Output is not valid JSON: %v", err)
			}
		})
	}
}

func TestPrintQuietMode(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter("json", true, &buf)

	data := map[string]string{"test": "value"}
	err := p.Print(data)
	if err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("Expected no output in quiet mode, got: %s", buf.String())
	}
}

func TestPrintError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedCode  string
		expectedMsg   string
		hasDetails    bool
		detailsLength int
	}{
		{
			name:         "CLIError with details",
			err:          models.NewCLIError(models.ErrAuthRequired, "Authentication required", map[string]interface{}{"hint": "Run 'amazon-cli auth login'"}),
			expectedCode: models.ErrAuthRequired,
			expectedMsg:  "Authentication required",
			hasDetails:   true,
		},
		{
			name:         "CLIError without details",
			err:          models.NewCLIError(models.ErrNotFound, "Resource not found", nil),
			expectedCode: models.ErrNotFound,
			expectedMsg:  "Resource not found",
			hasDetails:   false,
		},
		{
			name:         "Generic error",
			err:          errors.New("something went wrong"),
			expectedCode: "GENERAL_ERROR",
			expectedMsg:  "something went wrong",
			hasDetails:   false,
		},
		{
			name:         "Auth expired error",
			err:          models.NewCLIError(models.ErrAuthExpired, "Token has expired", nil),
			expectedCode: models.ErrAuthExpired,
			expectedMsg:  "Token has expired",
			hasDetails:   false,
		},
		{
			name:         "Rate limited error",
			err:          models.NewCLIError(models.ErrRateLimited, "Too many requests", map[string]interface{}{"retry_after": 60}),
			expectedCode: models.ErrRateLimited,
			expectedMsg:  "Too many requests",
			hasDetails:   true,
		},
		{
			name:         "Network error",
			err:          models.NewCLIError(models.ErrNetworkError, "Connection failed", map[string]interface{}{"url": "https://amazon.com"}),
			expectedCode: models.ErrNetworkError,
			expectedMsg:  "Connection failed",
			hasDetails:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			p := NewPrinterWithWriter("json", false, &buf)

			err := p.PrintError(tt.err)
			if err != nil {
				t.Fatalf("PrintError returned error: %v", err)
			}

			output := strings.TrimSpace(buf.String())

			// Parse the error response
			var response models.ErrorResponse
			if err := json.Unmarshal([]byte(output), &response); err != nil {
				t.Fatalf("Failed to parse error JSON: %v\nOutput: %s", err, output)
			}

			if response.Error == nil {
				t.Fatal("Error field is nil in response")
			}

			if response.Error.Code != tt.expectedCode {
				t.Errorf("Expected error code %s, got %s", tt.expectedCode, response.Error.Code)
			}

			if response.Error.Message != tt.expectedMsg {
				t.Errorf("Expected error message %s, got %s", tt.expectedMsg, response.Error.Message)
			}

			if tt.hasDetails {
				if len(response.Error.Details) == 0 {
					t.Error("Expected error details but got none")
				}
			} else {
				if len(response.Error.Details) != 0 {
					t.Errorf("Expected no error details but got %v", response.Error.Details)
				}
			}

			// Validate JSON structure
			var jsonMap map[string]interface{}
			if err := json.Unmarshal([]byte(output), &jsonMap); err != nil {
				t.Fatalf("Output is not valid JSON: %v", err)
			}

			// Check that "error" key exists at top level
			if _, ok := jsonMap["error"]; !ok {
				t.Error("Output JSON missing 'error' key")
			}
		})
	}
}

func TestPrintErrorQuietMode(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter("json", true, &buf)

	err := models.NewCLIError(models.ErrNotFound, "Not found", nil)
	if err2 := p.PrintError(err); err2 != nil {
		t.Fatalf("PrintError returned error: %v", err2)
	}

	if buf.Len() != 0 {
		t.Errorf("Expected no output in quiet mode, got: %s", buf.String())
	}
}

func TestPrintRawFormat(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter("raw", false, &buf)

	data := "test string"
	err := p.Print(data)
	if err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != data {
		t.Errorf("Expected raw output '%s', got '%s'", data, output)
	}
}

func TestPrintTableFormatNotImplemented(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter("table", false, &buf)

	data := map[string]string{"test": "value"}
	err := p.Print(data)
	if err == nil {
		t.Fatal("Expected error for table format, got nil")
	}

	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("Expected 'not yet implemented' error, got: %v", err)
	}
}

func TestPrintUnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter("unknown", false, &buf)

	data := map[string]string{"test": "value"}
	err := p.Print(data)
	if err == nil {
		t.Fatal("Expected error for unknown format, got nil")
	}

	if !strings.Contains(err.Error(), "unknown output format") {
		t.Errorf("Expected 'unknown output format' error, got: %v", err)
	}
}

func TestJSONOutputIndentation(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter("json", false, &buf)

	data := map[string]interface{}{
		"parent": map[string]interface{}{
			"child": "value",
		},
	}

	err := p.Print(data)
	if err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	output := buf.String()

	// Check that output contains proper indentation (2 spaces)
	if !strings.Contains(output, "  ") {
		t.Error("Output does not appear to be indented")
	}

	// Verify it's valid JSON
	var parsed interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}
}

func TestComplexDataStructures(t *testing.T) {
	// Test with a complex nested structure similar to what might be in the PRD
	type OrderItem struct {
		ASIN     string  `json:"asin"`
		Title    string  `json:"title"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}

	type Order struct {
		OrderID string      `json:"order_id"`
		Date    string      `json:"date"`
		Total   float64     `json:"total"`
		Status  string      `json:"status"`
		Items   []OrderItem `json:"items"`
	}

	type OrdersResponse struct {
		Orders     []Order `json:"orders"`
		TotalCount int     `json:"total_count"`
	}

	data := OrdersResponse{
		Orders: []Order{
			{
				OrderID: "123-4567890-1234567",
				Date:    "2024-01-15",
				Total:   29.99,
				Status:  "delivered",
				Items: []OrderItem{
					{
						ASIN:     "B08N5WRWNW",
						Title:    "Product Name",
						Quantity: 1,
						Price:    29.99,
					},
				},
			},
		},
		TotalCount: 1,
	}

	var buf bytes.Buffer
	p := NewPrinterWithWriter("json", false, &buf)

	err := p.Print(data)
	if err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	output := buf.String()

	// Parse and verify structure
	var parsed OrdersResponse
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	if parsed.TotalCount != 1 {
		t.Errorf("Expected total_count 1, got %d", parsed.TotalCount)
	}

	if len(parsed.Orders) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(parsed.Orders))
	}

	if parsed.Orders[0].OrderID != "123-4567890-1234567" {
		t.Errorf("Expected order ID '123-4567890-1234567', got '%s'", parsed.Orders[0].OrderID)
	}

	if parsed.Orders[0].Total != 29.99 {
		t.Errorf("Expected total 29.99, got %f", parsed.Orders[0].Total)
	}
}

func TestAllErrorCodes(t *testing.T) {
	errorCodes := []struct {
		code    string
		message string
	}{
		{models.ErrAuthRequired, "Not logged in"},
		{models.ErrAuthExpired, "Token expired"},
		{models.ErrNotFound, "Resource not found"},
		{models.ErrRateLimited, "Too many requests"},
		{models.ErrInvalidInput, "Invalid command input"},
		{models.ErrPurchaseFailed, "Purchase could not be completed"},
		{models.ErrNetworkError, "Network connectivity issue"},
		{models.ErrAmazonError, "Amazon returned an error"},
	}

	for _, tc := range errorCodes {
		t.Run(tc.code, func(t *testing.T) {
			var buf bytes.Buffer
			p := NewPrinterWithWriter("json", false, &buf)

			err := models.NewCLIError(tc.code, tc.message, nil)
			if err2 := p.PrintError(err); err2 != nil {
				t.Fatalf("PrintError returned error: %v", err2)
			}

			var response models.ErrorResponse
			if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse error JSON: %v", err)
			}

			if response.Error.Code != tc.code {
				t.Errorf("Expected code %s, got %s", tc.code, response.Error.Code)
			}

			if response.Error.Message != tc.message {
				t.Errorf("Expected message %s, got %s", tc.message, response.Error.Message)
			}
		})
	}
}
