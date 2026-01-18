package amazon

import (
	"testing"
)

func TestCancelSubscription(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name          string
		subscriptionID string
		wantErr       bool
		errContains   string
	}{
		{
			name:           "Valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "Another valid subscription ID",
			subscriptionID: "S99-9999999-9999999",
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
					t.Errorf("CancelSubscription() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("CancelSubscription() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CancelSubscription() unexpected error = %v", err)
				return
			}

			if subscription == nil {
				t.Error("CancelSubscription() returned nil subscription")
				return
			}

			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("CancelSubscription() subscription ID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}

			if subscription.Status != "cancelled" {
				t.Errorf("CancelSubscription() status = %v, want 'cancelled'", subscription.Status)
			}
		})
	}
}

func TestCancelSubscription_AlreadyCancelled(t *testing.T) {
	client := NewClient()

	// First cancellation should succeed
	subscriptionID := "S01-1234567-8901234"
	subscription, err := client.CancelSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("First CancelSubscription() failed: %v", err)
	}

	if subscription.Status != "cancelled" {
		t.Errorf("First CancelSubscription() status = %v, want 'cancelled'", subscription.Status)
	}

	// Note: In the current implementation, GetSubscription returns a fresh subscription
	// so we can't test double-cancellation without a more sophisticated mock.
	// This is a limitation of the current test setup. In a real implementation with
	// actual state management, we would test that cancelling an already-cancelled
	// subscription returns an error.
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
					t.Errorf("GetSubscription() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetSubscription() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GetSubscription() unexpected error = %v", err)
				return
			}

			if subscription == nil {
				t.Error("GetSubscription() returned nil subscription")
				return
			}

			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("GetSubscription() subscription ID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}

			// Verify subscription has required fields
			if subscription.ASIN == "" {
				t.Error("GetSubscription() subscription ASIN is empty")
			}
			if subscription.Title == "" {
				t.Error("GetSubscription() subscription Title is empty")
			}
			if subscription.Status == "" {
				t.Error("GetSubscription() subscription Status is empty")
			}
		})
	}
}

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()

	response, err := client.GetSubscriptions()
	if err != nil {
		t.Errorf("GetSubscriptions() unexpected error = %v", err)
		return
	}

	if response == nil {
		t.Error("GetSubscriptions() returned nil response")
		return
	}

	// In the mock implementation, this returns an empty slice
	if response.Subscriptions == nil {
		t.Error("GetSubscriptions() Subscriptions field is nil, expected empty slice")
	}
}

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
					t.Errorf("SkipDelivery() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("SkipDelivery() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("SkipDelivery() unexpected error = %v", err)
				return
			}

			if subscription == nil {
				t.Error("SkipDelivery() returned nil subscription")
			}
		})
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
			name:           "Valid frequency - 4 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          4,
			wantErr:        false,
		},
		{
			name:           "Valid frequency - 1 week",
			subscriptionID: "S01-1234567-8901234",
			weeks:          1,
			wantErr:        false,
		},
		{
			name:           "Valid frequency - 26 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          26,
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
			name:           "Invalid frequency - 0 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          0,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "Invalid frequency - negative",
			subscriptionID: "S01-1234567-8901234",
			weeks:          -1,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "Invalid frequency - too high",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateFrequency() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateFrequency() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateFrequency() unexpected error = %v", err)
				return
			}

			if subscription == nil {
				t.Error("UpdateFrequency() returned nil subscription")
				return
			}

			if subscription.FrequencyWeeks != tt.weeks {
				t.Errorf("UpdateFrequency() frequency = %v, want %v", subscription.FrequencyWeeks, tt.weeks)
			}
		})
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Errorf("GetUpcomingDeliveries() unexpected error = %v", err)
		return
	}

	// In the mock implementation, this returns an empty slice
	if deliveries == nil {
		t.Error("GetUpcomingDeliveries() returned nil, expected empty slice")
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
