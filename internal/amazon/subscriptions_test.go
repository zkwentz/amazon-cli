package amazon

import (
	"testing"
	"time"
)

func TestSkipDelivery(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "Valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.SkipDelivery(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SkipDelivery() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("SkipDelivery() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("SkipDelivery() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Errorf("SkipDelivery() returned nil subscription")
				return
			}

			// Verify the subscription was returned with the correct ID
			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("SkipDelivery() subscription ID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}

			// Verify the next delivery date was updated (should be in the future)
			if subscription.NextDelivery == "" {
				t.Errorf("SkipDelivery() next delivery date is empty")
			}

			// Parse and verify the delivery date is in the future
			deliveryDate, err := time.Parse("2006-01-02", subscription.NextDelivery)
			if err != nil {
				t.Errorf("SkipDelivery() invalid delivery date format: %v", err)
			}

			// The new delivery should be at least one day in the future
			// (using "now" as reference, though in real implementation it would be based on current delivery date)
			if !deliveryDate.After(time.Now().AddDate(0, 0, -1)) {
				t.Errorf("SkipDelivery() delivery date %v should be in the future", deliveryDate)
			}
		})
	}
}

func TestSkipDelivery_DateCalculation(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	// Skip delivery
	subscription, err := client.SkipDelivery(subscriptionID)
	if err != nil {
		t.Fatalf("SkipDelivery() unexpected error: %v", err)
	}

	// Parse the original delivery date (from GetSubscription)
	originalSub, _ := client.GetSubscription(subscriptionID)
	originalDate, err := time.Parse("2006-01-02", originalSub.NextDelivery)
	if err != nil {
		t.Fatalf("Failed to parse original delivery date: %v", err)
	}

	// Parse the new delivery date
	newDate, err := time.Parse("2006-01-02", subscription.NextDelivery)
	if err != nil {
		t.Fatalf("Failed to parse new delivery date: %v", err)
	}

	// Calculate expected date (original + frequency weeks)
	expectedDate := originalDate.AddDate(0, 0, subscription.FrequencyWeeks*7)

	// Verify the new date matches expected (allowing for same day)
	if !newDate.Equal(expectedDate) {
		t.Errorf("SkipDelivery() new delivery date = %v, want %v (original %v + %d weeks)",
			newDate, expectedDate, originalDate, subscription.FrequencyWeeks)
	}
}

func TestGetSubscription(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "Valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.GetSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetSubscription() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetSubscription() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GetSubscription() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Errorf("GetSubscription() returned nil subscription")
				return
			}

			// Verify subscription has required fields
			if subscription.SubscriptionID == "" {
				t.Errorf("GetSubscription() subscription ID is empty")
			}
			if subscription.ASIN == "" {
				t.Errorf("GetSubscription() ASIN is empty")
			}
			if subscription.Status == "" {
				t.Errorf("GetSubscription() status is empty")
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
		t.Errorf("GetSubscriptions() returned nil response")
	}

	// For placeholder implementation, this returns empty array
	if response.Subscriptions == nil {
		t.Errorf("GetSubscriptions() subscriptions array is nil")
	}
}

func TestUpdateFrequency(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name           string
		subscriptionID string
		weeks          int
		wantErr        bool
		errContains    string
	}{
		{
			name:           "Valid frequency update",
			subscriptionID: "S01-1234567-8901234",
			weeks:          8,
			wantErr:        false,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			weeks:          4,
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "Invalid weeks - zero",
			subscriptionID: "S01-1234567-8901234",
			weeks:          0,
			wantErr:        true,
			errContains:    "frequency weeks must be between 1 and 26",
		},
		{
			name:           "Invalid weeks - negative",
			subscriptionID: "S01-1234567-8901234",
			weeks:          -1,
			wantErr:        true,
			errContains:    "frequency weeks must be between 1 and 26",
		},
		{
			name:           "Invalid weeks - too high",
			subscriptionID: "S01-1234567-8901234",
			weeks:          30,
			wantErr:        true,
			errContains:    "frequency weeks must be between 1 and 26",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateFrequency() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateFrequency() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateFrequency() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Errorf("UpdateFrequency() returned nil subscription")
				return
			}

			if subscription.FrequencyWeeks != tt.weeks {
				t.Errorf("UpdateFrequency() frequency = %v, want %v", subscription.FrequencyWeeks, tt.weeks)
			}
		})
	}
}

func TestCancelSubscription(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "Valid cancellation",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.CancelSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CancelSubscription() expected error but got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("CancelSubscription() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CancelSubscription() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Errorf("CancelSubscription() returned nil subscription")
				return
			}

			if subscription.Status != "cancelled" {
				t.Errorf("CancelSubscription() status = %v, want cancelled", subscription.Status)
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
		t.Errorf("GetUpcomingDeliveries() returned nil")
	}

	// For placeholder implementation, this returns empty array
	if len(deliveries) != 0 {
		t.Logf("GetUpcomingDeliveries() returned %d deliveries", len(deliveries))
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
