package amazon

import (
	"testing"
	"time"
)

func TestGetOrderHistory(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		year        int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid current year",
			year:        time.Now().Year(),
			expectError: false,
		},
		{
			name:        "Valid past year",
			year:        2020,
			expectError: false,
		},
		{
			name:        "Year too old (before Amazon existed)",
			year:        1990,
			expectError: true,
			errorMsg:    "invalid year",
		},
		{
			name:        "Year in the future",
			year:        time.Now().Year() + 1,
			expectError: true,
			errorMsg:    "invalid year",
		},
		{
			name:        "Edge case - 1995 (Amazon founded)",
			year:        1995,
			expectError: false,
		},
		{
			name:        "Edge case - current year",
			year:        time.Now().Year(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.GetOrderHistory(tt.year)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorMsg != "" && err != nil {
					if len(err.Error()) == 0 || err.Error()[:len(tt.errorMsg)] != tt.errorMsg {
						t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected non-nil response")
				}
				if response != nil && response.Orders == nil {
					t.Errorf("Expected Orders slice to be initialized")
				}
			}
		})
	}
}

func TestGetOrderHistory_ValidResponse(t *testing.T) {
	client := NewClient()
	year := 2023

	response, err := client.GetOrderHistory(year)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if response.Orders == nil {
		t.Error("Expected Orders slice to be initialized")
	}

	if response.TotalCount < 0 {
		t.Error("Expected TotalCount to be non-negative")
	}

	// Verify that TotalCount matches the length of Orders slice
	if response.TotalCount != len(response.Orders) {
		t.Errorf("Expected TotalCount (%d) to match Orders length (%d)", response.TotalCount, len(response.Orders))
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.httpClient == nil {
		t.Error("Expected non-nil HTTP client")
	}

	if client.baseURL == "" {
		t.Error("Expected non-empty base URL")
	}

	if client.baseURL != "https://www.amazon.com" {
		t.Errorf("Expected base URL to be 'https://www.amazon.com', got '%s'", client.baseURL)
	}
}
