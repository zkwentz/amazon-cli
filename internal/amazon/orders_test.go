package amazon

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetOrderHistory(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	currentYear := time.Now().Year()

	tests := []struct {
		name        string
		year        int
		expectError bool
		errorCode   string
	}{
		{
			name:        "Valid current year",
			year:        currentYear,
			expectError: false,
		},
		{
			name:        "Valid past year",
			year:        2020,
			expectError: false,
		},
		{
			name:        "Invalid year - too old",
			year:        1990,
			expectError: true,
			errorCode:   models.ErrorCodeInvalidInput,
		},
		{
			name:        "Invalid year - future",
			year:        currentYear + 1,
			expectError: true,
			errorCode:   models.ErrorCodeInvalidInput,
		},
		{
			name:        "Edge case - Amazon founding year",
			year:        1995,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.GetOrderHistory(tt.year)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}

				cliErr, ok := err.(*models.CLIError)
				if !ok {
					t.Errorf("Expected CLIError but got: %T", err)
					return
				}

				if cliErr.Code != tt.errorCode {
					t.Errorf("Expected error code %s but got %s", tt.errorCode, cliErr.Code)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				if response == nil {
					t.Error("Expected response but got nil")
					return
				}

				// Verify response structure
				if response.Orders == nil {
					t.Error("Expected Orders slice but got nil")
				}

				if response.TotalCount != len(response.Orders) {
					t.Errorf("TotalCount (%d) doesn't match actual orders count (%d)",
						response.TotalCount, len(response.Orders))
				}

				// Verify order data contains the correct year
				for _, order := range response.Orders {
					if order.Date[:4] != string(rune(tt.year/1000)+'0')+string(rune((tt.year/100)%10)+'0')+string(rune((tt.year/10)%10)+'0')+string(rune(tt.year%10)+'0') {
						t.Errorf("Order date %s doesn't start with year %d", order.Date, tt.year)
					}
				}
			}
		})
	}
}

func TestGetOrderHistoryResponseStructure(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	response, err := client.GetOrderHistory(2024)
	if err != nil {
		t.Fatalf("Failed to get order history: %v", err)
	}

	// Test that response has expected fields
	if response.Orders == nil {
		t.Error("Orders slice is nil")
	}

	if len(response.Orders) == 0 {
		t.Error("Expected at least one order in mock response")
		return
	}

	// Test first order structure
	order := response.Orders[0]

	if order.OrderID == "" {
		t.Error("Order ID is empty")
	}

	if order.Date == "" {
		t.Error("Order date is empty")
	}

	if order.Total <= 0 {
		t.Error("Order total should be positive")
	}

	if order.Status == "" {
		t.Error("Order status is empty")
	}

	if len(order.Items) == 0 {
		t.Error("Order should have at least one item")
	}

	// Test order item structure
	item := order.Items[0]
	if item.ASIN == "" {
		t.Error("Item ASIN is empty")
	}

	if item.Title == "" {
		t.Error("Item title is empty")
	}

	if item.Quantity <= 0 {
		t.Error("Item quantity should be positive")
	}

	if item.Price <= 0 {
		t.Error("Item price should be positive")
	}

	// Test tracking structure (if present)
	if order.Tracking != nil {
		if order.Tracking.Carrier == "" {
			t.Error("Tracking carrier is empty")
		}

		if order.Tracking.TrackingNumber == "" {
			t.Error("Tracking number is empty")
		}

		if order.Tracking.Status == "" {
			t.Error("Tracking status is empty")
		}
	}
}
