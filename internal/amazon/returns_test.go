package amazon

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestCreateReturn(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		orderID     string
		itemID      string
		reason      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid return creation",
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  "defective",
			wantErr: false,
		},
		{
			name:        "missing orderID",
			orderID:     "",
			itemID:      "ITEM123",
			reason:      "defective",
			wantErr:     true,
			errContains: "orderID is required",
		},
		{
			name:        "missing itemID",
			orderID:     "123-4567890-1234567",
			itemID:      "",
			reason:      "defective",
			wantErr:     true,
			errContains: "itemID is required",
		},
		{
			name:        "missing reason",
			orderID:     "123-4567890-1234567",
			itemID:      "ITEM123",
			reason:      "",
			wantErr:     true,
			errContains: "reason is required",
		},
		{
			name:        "invalid reason code",
			orderID:     "123-4567890-1234567",
			itemID:      "ITEM123",
			reason:      "invalid_reason",
			wantErr:     true,
			errContains: "invalid return reason",
		},
		{
			name:    "valid reason - wrong_item",
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  "wrong_item",
			wantErr: false,
		},
		{
			name:    "valid reason - not_as_described",
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  "not_as_described",
			wantErr: false,
		},
		{
			name:    "valid reason - no_longer_needed",
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  "no_longer_needed",
			wantErr: false,
		},
		{
			name:    "valid reason - better_price",
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  "better_price",
			wantErr: false,
		},
		{
			name:    "valid reason - other",
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			reason:  "other",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.CreateReturn(tt.orderID, tt.itemID, tt.reason)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateReturn() expected error but got none")
					return
				}
				if tt.errContains != "" {
					if err.Error() == "" || !contains(err.Error(), tt.errContains) {
						t.Errorf("CreateReturn() error = %v, want error containing %q", err, tt.errContains)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("CreateReturn() unexpected error = %v", err)
				return
			}

			// Validate the result
			if result == nil {
				t.Error("CreateReturn() returned nil result")
				return
			}

			if result.ReturnID == "" {
				t.Error("CreateReturn() returned empty ReturnID")
			}

			if result.OrderID != tt.orderID {
				t.Errorf("CreateReturn() OrderID = %v, want %v", result.OrderID, tt.orderID)
			}

			if result.ItemID != tt.itemID {
				t.Errorf("CreateReturn() ItemID = %v, want %v", result.ItemID, tt.itemID)
			}

			if result.Reason != tt.reason {
				t.Errorf("CreateReturn() Reason = %v, want %v", result.Reason, tt.reason)
			}

			if result.Status != "initiated" {
				t.Errorf("CreateReturn() Status = %v, want %v", result.Status, "initiated")
			}

			if result.CreatedAt.IsZero() {
				t.Error("CreateReturn() CreatedAt is zero")
			}

			// Check that CreatedAt is recent (within last minute)
			if time.Since(result.CreatedAt) > time.Minute {
				t.Errorf("CreateReturn() CreatedAt is not recent: %v", result.CreatedAt)
			}
		})
	}
}

func TestIsValidReturnReason(t *testing.T) {
	tests := []struct {
		reason string
		valid  bool
	}{
		{"defective", true},
		{"wrong_item", true},
		{"not_as_described", true},
		{"no_longer_needed", true},
		{"better_price", true},
		{"other", true},
		{"invalid", false},
		{"", false},
		{"DEFECTIVE", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.reason, func(t *testing.T) {
			if got := models.IsValidReturnReason(tt.reason); got != tt.valid {
				t.Errorf("IsValidReturnReason(%q) = %v, want %v", tt.reason, got, tt.valid)
			}
		})
	}
}

func TestGetReturnOptions(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		orderID     string
		itemID      string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing orderID",
			orderID:     "",
			itemID:      "ITEM123",
			wantErr:     true,
			errContains: "orderID is required",
		},
		{
			name:        "missing itemID",
			orderID:     "123-4567890-1234567",
			itemID:      "",
			wantErr:     true,
			errContains: "itemID is required",
		},
		{
			name:        "valid but not implemented",
			orderID:     "123-4567890-1234567",
			itemID:      "ITEM123",
			wantErr:     true,
			errContains: "Not implemented yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetReturnOptions(tt.orderID, tt.itemID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetReturnOptions() expected error but got none")
					return
				}
				if tt.errContains != "" {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("GetReturnOptions() error = %v, want error containing %q", err, tt.errContains)
					}
				}
			}
		})
	}
}

func TestGetReturnLabel(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		returnID    string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing returnID",
			returnID:    "",
			wantErr:     true,
			errContains: "returnID is required",
		},
		{
			name:        "valid but not implemented",
			returnID:    "RET123",
			wantErr:     true,
			errContains: "Not implemented yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetReturnLabel(tt.returnID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetReturnLabel() expected error but got none")
					return
				}
				if tt.errContains != "" {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("GetReturnLabel() error = %v, want error containing %q", err, tt.errContains)
					}
				}
			}
		})
	}
}

func TestGetReturnStatus(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		returnID    string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing returnID",
			returnID:    "",
			wantErr:     true,
			errContains: "returnID is required",
		},
		{
			name:        "valid but not implemented",
			returnID:    "RET123",
			wantErr:     true,
			errContains: "Not implemented yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetReturnStatus(tt.returnID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetReturnStatus() expected error but got none")
					return
				}
				if tt.errContains != "" {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("GetReturnStatus() error = %v, want error containing %q", err, tt.errContains)
					}
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
