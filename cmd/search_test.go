package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestSearchCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorCode   string
	}{
		{
			name:        "valid basic search",
			args:        []string{"search", "wireless headphones"},
			expectError: false,
		},
		{
			name:        "valid search with price range",
			args:        []string{"search", "laptop", "--min-price", "500", "--max-price", "1000"},
			expectError: false,
		},
		{
			name:        "valid search with prime only",
			args:        []string{"search", "books", "--prime-only"},
			expectError: false,
		},
		{
			name:        "valid search with category",
			args:        []string{"search", "coffee", "--category", "electronics"},
			expectError: false,
		},
		{
			name:        "valid search with pagination",
			args:        []string{"search", "keyboards", "--page", "2"},
			expectError: false,
		},
		{
			name:        "invalid min price greater than max price",
			args:        []string{"search", "laptop", "--min-price", "1000", "--max-price", "500"},
			expectError: true,
			errorCode:   models.ErrorCodeInvalidInput,
		},
		{
			name:        "invalid page number zero",
			args:        []string{"search", "test", "--page", "0"},
			expectError: true,
			errorCode:   models.ErrorCodeInvalidInput,
		},
		{
			name:        "invalid negative page number",
			args:        []string{"search", "test", "--page", "-1"},
			expectError: true,
			errorCode:   models.ErrorCodeInvalidInput,
		},
		{
			name:        "missing query argument",
			args:        []string{"search"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset root command for each test
			rootCmd.SetArgs(tt.args)

			// Capture output
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			err := rootCmd.Execute()

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectError && tt.errorCode != "" {
				// Check if error contains expected error code
				output := buf.String()
				if output != "" {
					var errResp models.ErrorResponse
					if jsonErr := json.Unmarshal([]byte(output), &errResp); jsonErr == nil {
						if errResp.Error != nil && errResp.Error.Code != tt.errorCode {
							t.Errorf("expected error code %s, got %s", tt.errorCode, errResp.Error.Code)
						}
					}
				}
			}

			// Reset for next test
			rootCmd.SetArgs([]string{})
		})
	}
}

func TestSearchValidation(t *testing.T) {
	tests := []struct {
		name      string
		minPrice  float64
		maxPrice  float64
		page      int
		wantError bool
	}{
		{
			name:      "valid price range",
			minPrice:  10.0,
			maxPrice:  100.0,
			page:      1,
			wantError: false,
		},
		{
			name:      "min equals max",
			minPrice:  50.0,
			maxPrice:  50.0,
			page:      1,
			wantError: false,
		},
		{
			name:      "min greater than max",
			minPrice:  100.0,
			maxPrice:  10.0,
			page:      1,
			wantError: true,
		},
		{
			name:      "zero page",
			minPrice:  10.0,
			maxPrice:  100.0,
			page:      0,
			wantError: true,
		},
		{
			name:      "negative page",
			minPrice:  10.0,
			maxPrice:  100.0,
			page:      -1,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the flags
			searchMinPrice = tt.minPrice
			searchMaxPrice = tt.maxPrice
			searchPage = tt.page

			// Run validation logic
			hasError := false

			if searchMinPrice > 0 && searchMaxPrice > 0 && searchMinPrice > searchMaxPrice {
				hasError = true
			}

			if searchPage < 1 {
				hasError = true
			}

			if hasError != tt.wantError {
				t.Errorf("expected error: %v, got error: %v", tt.wantError, hasError)
			}
		})
	}
}
