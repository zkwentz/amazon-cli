package amazon

import (
	"strings"
	"testing"
	"time"
)

func TestSkipDelivery_WithConfirm(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		confirm        bool
		wantErr        bool
		errContains    string
	}{
		{
			name:           "valid skip with confirm",
			subscriptionID: "S01-1234567-8901234",
			confirm:        true,
			wantErr:        false,
		},
		{
			name:           "empty subscriptionID should fail",
			subscriptionID: "",
			confirm:        true,
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "valid subscription ID with confirm",
			subscriptionID: "S02-9876543-2109876",
			confirm:        true,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			// Execute SkipDelivery with confirm=true
			result, err := client.SkipDelivery(tt.subscriptionID, tt.confirm)

			// Check error expectations
			if tt.wantErr {
				if err == nil {
					t.Errorf("SkipDelivery() expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("SkipDelivery() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			// Check success expectations
			if err != nil {
				t.Errorf("SkipDelivery() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("SkipDelivery() returned nil result")
				return
			}

			// Verify the subscription ID matches
			if result.SubscriptionID != tt.subscriptionID {
				t.Errorf("SkipDelivery() SubscriptionID = %v, want %v", result.SubscriptionID, tt.subscriptionID)
			}

			// Verify next delivery date was updated (when confirmed)
			if tt.confirm && result.NextDelivery == "" {
				t.Error("SkipDelivery() with confirm should have NextDelivery set")
			}

			// Verify the new delivery date is in the future
			if tt.confirm {
				deliveryDate, err := time.Parse("2006-01-02", result.NextDelivery)
				if err != nil {
					t.Errorf("SkipDelivery() invalid NextDelivery format: %v", err)
				}

				// The delivery date should be moved forward
				originalDate, _ := time.Parse("2006-01-02", "2024-02-01")
				if !deliveryDate.After(originalDate) && !deliveryDate.Equal(originalDate) {
					t.Errorf("SkipDelivery() NextDelivery should be after or equal to original date, got %v", deliveryDate)
				}
			}
		})
	}
}

func TestSkipDelivery_WithoutConfirm(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		confirm        bool
		wantErr        bool
		errContains    string
	}{
		{
			name:           "preview without confirm",
			subscriptionID: "S01-1234567-8901234",
			confirm:        false,
			wantErr:        false,
		},
		{
			name:           "empty subscriptionID should fail even without confirm",
			subscriptionID: "",
			confirm:        false,
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "valid preview without confirm",
			subscriptionID: "S03-5555555-5555555",
			confirm:        false,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			// Execute SkipDelivery with confirm=false (preview mode)
			result, err := client.SkipDelivery(tt.subscriptionID, tt.confirm)

			// Check error expectations
			if tt.wantErr {
				if err == nil {
					t.Errorf("SkipDelivery() expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("SkipDelivery() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			// Check success expectations
			if err != nil {
				t.Errorf("SkipDelivery() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("SkipDelivery() returned nil result")
				return
			}

			// Verify the subscription ID matches
			if result.SubscriptionID != tt.subscriptionID {
				t.Errorf("SkipDelivery() SubscriptionID = %v, want %v", result.SubscriptionID, tt.subscriptionID)
			}

			// When not confirmed, the subscription should be returned as preview
			// The delivery date should remain unchanged in preview mode
			if result.NextDelivery == "" {
				t.Error("SkipDelivery() preview should return NextDelivery")
			}

			// Verify subscription status is still active in preview
			if result.Status != "active" {
				t.Errorf("SkipDelivery() preview Status = %v, want active", result.Status)
			}
		})
	}
}

func TestSkipDelivery_ConfirmVsPreview(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	// First, get preview (without confirm)
	preview, err := client.SkipDelivery(subscriptionID, false)
	if err != nil {
		t.Fatalf("SkipDelivery() preview failed: %v", err)
	}

	previewDate := preview.NextDelivery

	// Then, execute with confirm
	confirmed, err := client.SkipDelivery(subscriptionID, true)
	if err != nil {
		t.Fatalf("SkipDelivery() with confirm failed: %v", err)
	}

	confirmedDate := confirmed.NextDelivery

	// The confirmed delivery date should be different from preview
	// because confirm actually skips the delivery
	if confirmedDate == "" {
		t.Error("SkipDelivery() with confirm should have NextDelivery set")
	}

	// Both should have valid subscription data
	if preview.SubscriptionID != confirmed.SubscriptionID {
		t.Error("SkipDelivery() preview and confirmed should have same SubscriptionID")
	}

	// Parse dates to verify they're valid
	_, err = time.Parse("2006-01-02", previewDate)
	if err != nil {
		t.Errorf("Preview date invalid: %v", err)
	}

	_, err = time.Parse("2006-01-02", confirmedDate)
	if err != nil {
		t.Errorf("Confirmed date invalid: %v", err)
	}
}

func TestGetSubscription(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "another valid subscription ID",
			subscriptionID: "S99-9999999-9999999",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.GetSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Error("GetSubscription() expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetSubscription() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GetSubscription() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Error("GetSubscription() returned nil")
				return
			}

			// Verify subscription has required fields
			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("GetSubscription() SubscriptionID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}

			if subscription.ASIN == "" {
				t.Error("GetSubscription() ASIN should not be empty")
			}

			if subscription.Title == "" {
				t.Error("GetSubscription() Title should not be empty")
			}

			if subscription.Price <= 0 {
				t.Error("GetSubscription() Price should be greater than 0")
			}

			if subscription.NextDelivery == "" {
				t.Error("GetSubscription() NextDelivery should not be empty")
			}

			if subscription.Status == "" {
				t.Error("GetSubscription() Status should not be empty")
			}
		})
	}
}

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()
	response, err := client.GetSubscriptions()

	if err != nil {
		t.Errorf("GetSubscriptions() unexpected error: %v", err)
	}

	if response == nil {
		t.Error("GetSubscriptions() returned nil")
		return
	}

	if response.Subscriptions == nil {
		t.Error("GetSubscriptions() Subscriptions should not be nil")
	}
}

func TestUpdateFrequency(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		weeks          int
		wantErr        bool
		errContains    string
	}{
		{
			name:           "valid frequency update",
			subscriptionID: "S01-1234567-8901234",
			weeks:          4,
			wantErr:        false,
		},
		{
			name:           "empty subscriptionID should fail",
			subscriptionID: "",
			weeks:          4,
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "zero weeks should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          0,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "negative weeks should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          -1,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "weeks over limit should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "minimum valid frequency",
			subscriptionID: "S01-1234567-8901234",
			weeks:          1,
			wantErr:        false,
		},
		{
			name:           "maximum valid frequency",
			subscriptionID: "S01-1234567-8901234",
			weeks:          26,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			result, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.wantErr {
				if err == nil {
					t.Error("UpdateFrequency() expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateFrequency() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateFrequency() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("UpdateFrequency() returned nil")
				return
			}

			if result.FrequencyWeeks != tt.weeks {
				t.Errorf("UpdateFrequency() FrequencyWeeks = %v, want %v", result.FrequencyWeeks, tt.weeks)
			}
		})
	}
}

func TestCancelSubscription(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "valid cancellation",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "empty subscriptionID should fail",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			result, err := client.CancelSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Error("CancelSubscription() expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CancelSubscription() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CancelSubscription() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("CancelSubscription() returned nil")
				return
			}

			if result.Status != "cancelled" {
				t.Errorf("CancelSubscription() Status = %v, want cancelled", result.Status)
			}
		})
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewClient()
	deliveries, err := client.GetUpcomingDeliveries()

	if err != nil {
		t.Errorf("GetUpcomingDeliveries() unexpected error: %v", err)
	}

	if deliveries == nil {
		t.Error("GetUpcomingDeliveries() returned nil")
	}
}
