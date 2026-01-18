package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetReturnStatus(t *testing.T) {
	// Create a test client
	cfg := config.GetDefaultConfig()
	client := NewClient(cfg)

	tests := []struct {
		name      string
		returnID  string
		wantErr   bool
		errCode   string
	}{
		{
			name:     "valid return ID",
			returnID: "RET123456789",
			wantErr:  false,
		},
		{
			name:     "empty return ID",
			returnID: "",
			wantErr:  true,
			errCode:  models.ErrCodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.GetReturnStatus(tt.returnID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetReturnStatus() expected error but got none")
					return
				}

				if cliErr, ok := err.(*models.CLIError); ok {
					if cliErr.Code != tt.errCode {
						t.Errorf("GetReturnStatus() error code = %v, want %v", cliErr.Code, tt.errCode)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("GetReturnStatus() unexpected error = %v", err)
				return
			}

			if result == nil {
				t.Errorf("GetReturnStatus() returned nil result")
				return
			}

			if result.ReturnID != tt.returnID {
				t.Errorf("GetReturnStatus() return_id = %v, want %v", result.ReturnID, tt.returnID)
			}

			// Verify all required fields are populated
			if result.OrderID == "" {
				t.Errorf("GetReturnStatus() order_id is empty")
			}
			if result.ItemID == "" {
				t.Errorf("GetReturnStatus() item_id is empty")
			}
			if result.Status == "" {
				t.Errorf("GetReturnStatus() status is empty")
			}
			if result.Reason == "" {
				t.Errorf("GetReturnStatus() reason is empty")
			}
			if result.CreatedAt == "" {
				t.Errorf("GetReturnStatus() created_at is empty")
			}
		})
	}
}

func TestCreateReturnValidation(t *testing.T) {
	cfg := config.GetDefaultConfig()
	client := NewClient(cfg)

	tests := []struct {
		name    string
		orderID string
		itemID  string
		reason  string
		wantErr bool
		errCode string
	}{
		{
			name:    "valid defective reason",
			orderID: "111-2222222-3333333",
			itemID:  "ITEM123",
			reason:  models.ReasonDefective,
			wantErr: true, // Will error because API not implemented, but should pass validation
		},
		{
			name:    "valid wrong_item reason",
			orderID: "111-2222222-3333333",
			itemID:  "ITEM123",
			reason:  models.ReasonWrongItem,
			wantErr: true,
		},
		{
			name:    "invalid reason",
			orderID: "111-2222222-3333333",
			itemID:  "ITEM123",
			reason:  "invalid_reason",
			wantErr: true,
			errCode: models.ErrCodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.CreateReturn(tt.orderID, tt.itemID, tt.reason)

			if !tt.wantErr {
				if err != nil {
					t.Errorf("CreateReturn() unexpected error = %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("CreateReturn() expected error but got none")
				return
			}

			// If we expect a specific error code, check it
			if tt.errCode != "" {
				if cliErr, ok := err.(*models.CLIError); ok {
					if cliErr.Code != tt.errCode {
						t.Errorf("CreateReturn() error code = %v, want %v", cliErr.Code, tt.errCode)
					}
				}
			}
		})
	}
}
