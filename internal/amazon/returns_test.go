package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestCreateReturnValidation(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		orderID     string
		itemID      string
		reason      string
		expectError bool
	}{
		{
			name:        "Valid return request",
			orderID:     "123-4567890-1234567",
			itemID:      "ITEM123",
			reason:      "defective",
			expectError: false,
		},
		{
			name:        "Invalid reason",
			orderID:     "123-4567890-1234567",
			itemID:      "ITEM123",
			reason:      "invalid_reason",
			expectError: true,
		},
		{
			name:        "Empty reason",
			orderID:     "123-4567890-1234567",
			itemID:      "ITEM123",
			reason:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.CreateReturn(tt.orderID, tt.itemID, tt.reason)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if _, ok := err.(*models.CLIError); !ok {
					t.Errorf("Expected CLIError but got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Expected return object but got nil")
				}
				if result.OrderID != tt.orderID {
					t.Errorf("OrderID = %s, want %s", result.OrderID, tt.orderID)
				}
				if result.ItemID != tt.itemID {
					t.Errorf("ItemID = %s, want %s", result.ItemID, tt.itemID)
				}
				if result.Reason != tt.reason {
					t.Errorf("Reason = %s, want %s", result.Reason, tt.reason)
				}
				if result.Status != "initiated" {
					t.Errorf("Status = %s, want initiated", result.Status)
				}
			}
		})
	}
}

func TestGetReturnableItems(t *testing.T) {
	client := NewClient()
	items, err := client.GetReturnableItems()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if items == nil {
		t.Errorf("Expected non-nil items slice")
	}
}

func TestGetReturnOptions(t *testing.T) {
	client := NewClient()
	options, err := client.GetReturnOptions("123-4567890-1234567", "ITEM123")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if options == nil {
		t.Errorf("Expected non-nil options slice")
	}

	if len(options) == 0 {
		t.Errorf("Expected at least one return option")
	}
}

func TestGetReturnLabel(t *testing.T) {
	client := NewClient()
	label, err := client.GetReturnLabel("R123-456-789")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if label == nil {
		t.Errorf("Expected non-nil label")
	}

	if label.URL == "" {
		t.Errorf("Expected non-empty URL")
	}
}

func TestGetReturnStatus(t *testing.T) {
	client := NewClient()
	returnObj, err := client.GetReturnStatus("R123-456-789")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if returnObj == nil {
		t.Errorf("Expected non-nil return object")
	}
}
