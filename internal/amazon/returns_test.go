package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetReturnOptions(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name      string
		orderID   string
		itemID    string
		wantErr   bool
		errCode   string
	}{
		{
			name:    "valid order and item",
			orderID: "123-4567890-1234567",
			itemID:  "ITEM123",
			wantErr: false,
		},
		{
			name:    "empty order ID",
			orderID: "",
			itemID:  "ITEM123",
			wantErr: true,
			errCode: models.ErrCodeInvalidInput,
		},
		{
			name:    "empty item ID",
			orderID: "123-4567890-1234567",
			itemID:  "",
			wantErr: true,
			errCode: models.ErrCodeInvalidInput,
		},
		{
			name:    "both empty",
			orderID: "",
			itemID:  "",
			wantErr: true,
			errCode: models.ErrCodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options, err := client.GetReturnOptions(tt.orderID, tt.itemID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetReturnOptions() expected error but got none")
					return
				}
				if cliErr, ok := err.(*models.CLIError); ok {
					if cliErr.Code != tt.errCode {
						t.Errorf("GetReturnOptions() error code = %v, want %v", cliErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("GetReturnOptions() error is not CLIError type")
				}
			} else {
				if err != nil {
					t.Errorf("GetReturnOptions() unexpected error = %v", err)
					return
				}
				if len(options) == 0 {
					t.Errorf("GetReturnOptions() returned empty options")
				}
				// Verify structure of returned options
				for _, option := range options {
					if option.Method == "" {
						t.Errorf("GetReturnOptions() option missing method")
					}
					if option.Label == "" {
						t.Errorf("GetReturnOptions() option missing label")
					}
				}
			}
		})
	}
}

func TestGetReturnableItems(t *testing.T) {
	client := NewClient()

	items, err := client.GetReturnableItems()
	if err != nil {
		t.Errorf("GetReturnableItems() unexpected error = %v", err)
	}

	// Currently returns empty list, but should not error
	if items == nil {
		t.Errorf("GetReturnableItems() returned nil instead of empty slice")
	}
}
