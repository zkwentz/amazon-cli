package amazon

import (
	"testing"
	"time"
)

func TestGetSubscription(t *testing.T) {
	client := NewClient()

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
			name:           "empty subscription ID",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscription ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.GetSubscription(tt.subscriptionID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetSubscription() expected error but got none")
				}
				if tt.errContains != "" && err != nil {
					if !containsString(err.Error(), tt.errContains) {
						t.Errorf("GetSubscription() error = %v, want error containing %v", err, tt.errContains)
					}
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
			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("GetSubscription() subscriptionID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}
			if subscription.Status != "active" {
				t.Errorf("GetSubscription() status = %v, want active", subscription.Status)
			}
		})
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
			name:           "valid skip delivery",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "empty subscription ID",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscription ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get original subscription first
			var originalDelivery string
			if !tt.wantErr {
				original, err := client.GetSubscription(tt.subscriptionID)
				if err != nil {
					t.Fatalf("Failed to get original subscription: %v", err)
				}
				originalDelivery = original.NextDelivery
			}

			// Skip delivery
			subscription, err := client.SkipDelivery(tt.subscriptionID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SkipDelivery() expected error but got none")
				}
				if tt.errContains != "" && err != nil {
					if !containsString(err.Error(), tt.errContains) {
						t.Errorf("SkipDelivery() error = %v, want error containing %v", err, tt.errContains)
					}
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

			// Verify next delivery date was updated
			originalDate, err := time.Parse("2006-01-02", originalDelivery)
			if err != nil {
				t.Fatalf("Failed to parse original delivery date: %v", err)
			}
			newDate, err := time.Parse("2006-01-02", subscription.NextDelivery)
			if err != nil {
				t.Fatalf("Failed to parse new delivery date: %v", err)
			}

			// New date should be after original date
			if !newDate.After(originalDate) {
				t.Errorf("SkipDelivery() new delivery date %v should be after original date %v", newDate, originalDate)
			}

			// New date should be approximately frequency_weeks * 7 days later
			expectedDiff := time.Duration(subscription.FrequencyWeeks*7) * 24 * time.Hour
			actualDiff := newDate.Sub(originalDate)
			if actualDiff < expectedDiff-24*time.Hour || actualDiff > expectedDiff+24*time.Hour {
				t.Errorf("SkipDelivery() date difference %v should be approximately %v", actualDiff, expectedDiff)
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
			name:           "valid frequency update",
			subscriptionID: "S01-1234567-8901234",
			weeks:          8,
			wantErr:        false,
		},
		{
			name:           "empty subscription ID",
			subscriptionID: "",
			weeks:          4,
			wantErr:        true,
			errContains:    "subscription ID cannot be empty",
		},
		{
			name:           "zero weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          0,
			wantErr:        true,
			errContains:    "frequency must be positive",
		},
		{
			name:           "negative weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          -1,
			wantErr:        true,
			errContains:    "frequency must be positive",
		},
		{
			name:           "too many weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			wantErr:        true,
			errContains:    "frequency cannot exceed 26 weeks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateFrequency() expected error but got none")
				}
				if tt.errContains != "" && err != nil {
					if !containsString(err.Error(), tt.errContains) {
						t.Errorf("UpdateFrequency() error = %v, want error containing %v", err, tt.errContains)
					}
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
			name:           "valid cancellation",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "empty subscription ID",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscription ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.CancelSubscription(tt.subscriptionID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("CancelSubscription() expected error but got none")
				}
				if tt.errContains != "" && err != nil {
					if !containsString(err.Error(), tt.errContains) {
						t.Errorf("CancelSubscription() error = %v, want error containing %v", err, tt.errContains)
					}
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

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()

	response, err := client.GetSubscriptions()
	if err != nil {
		t.Errorf("GetSubscriptions() unexpected error: %v", err)
		return
	}
	if response == nil {
		t.Errorf("GetSubscriptions() returned nil response")
		return
	}
	if response.Subscriptions == nil {
		t.Errorf("GetSubscriptions() returned nil subscriptions slice")
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Errorf("GetUpcomingDeliveries() unexpected error: %v", err)
		return
	}
	if deliveries == nil {
		t.Errorf("GetUpcomingDeliveries() returned nil deliveries slice")
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
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
