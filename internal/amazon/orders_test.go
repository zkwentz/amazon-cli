package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetOrderTracking(t *testing.T) {
	// Create a client with default config
	cfg := config.DefaultConfig()
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		name      string
		orderID   string
		wantError bool
		errorCode string
	}{
		{
			name:      "valid order ID",
			orderID:   "123-4567890-1234567",
			wantError: false,
		},
		{
			name:      "empty order ID",
			orderID:   "",
			wantError: true,
			errorCode: models.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracking, err := client.GetOrderTracking(tt.orderID)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}

				if cliErr, ok := err.(*models.CLIError); ok {
					if cliErr.Code != tt.errorCode {
						t.Errorf("expected error code %s, got %s", tt.errorCode, cliErr.Code)
					}
				} else {
					t.Error("expected CLIError type")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tracking == nil {
				t.Error("expected tracking data but got nil")
				return
			}

			// Verify tracking has required fields
			if tracking.Carrier == "" {
				t.Error("tracking carrier is empty")
			}
			if tracking.TrackingNumber == "" {
				t.Error("tracking number is empty")
			}
			if tracking.Status == "" {
				t.Error("tracking status is empty")
			}
		})
	}
}

func TestGetOrderTracking_Structure(t *testing.T) {
	cfg := config.DefaultConfig()
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tracking, err := client.GetOrderTracking("123-4567890-1234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the structure matches the expected schema
	expectedFields := map[string]bool{
		"Carrier":        tracking.Carrier != "",
		"TrackingNumber": tracking.TrackingNumber != "",
		"Status":         tracking.Status != "",
		"DeliveryDate":   tracking.DeliveryDate != "",
	}

	for field, hasValue := range expectedFields {
		if !hasValue {
			t.Errorf("field %s is missing or empty", field)
		}
	}
}
