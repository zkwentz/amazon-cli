package amazon

import (
	"testing"
	"time"
)

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
		{
			name:           "empty subscription ID",
			subscriptionID: "",
			weeks:          4,
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "zero weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          0,
			wantErr:        true,
			errContains:    "weeks must be positive",
		},
		{
			name:           "negative weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          -5,
			wantErr:        true,
			errContains:    "weeks must be positive",
		},
		{
			name:           "weeks exceed maximum",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			wantErr:        true,
			errContains:    "weeks cannot exceed 26",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateFrequency() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateFrequency() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateFrequency() unexpected error = %v", err)
				return
			}

			if subscription == nil {
				t.Errorf("UpdateFrequency() returned nil subscription")
				return
			}

			// Verify the frequency was updated
			if subscription.FrequencyWeeks != tt.weeks {
				t.Errorf("UpdateFrequency() frequency = %d, want %d", subscription.FrequencyWeeks, tt.weeks)
			}

			// Verify next delivery date is in the future
			nextDelivery, err := time.Parse("2006-01-02", subscription.NextDelivery)
			if err != nil {
				t.Errorf("UpdateFrequency() invalid next delivery date format: %v", err)
				return
			}

			if nextDelivery.Before(time.Now()) {
				t.Errorf("UpdateFrequency() next delivery date %s is in the past", subscription.NextDelivery)
			}

			// Verify subscription ID matches
			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("UpdateFrequency() subscription ID = %s, want %s", subscription.SubscriptionID, tt.subscriptionID)
			}
		})
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
			name:           "valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "empty subscription ID",
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
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetSubscription() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GetSubscription() unexpected error = %v", err)
				return
			}

			if subscription == nil {
				t.Errorf("GetSubscription() returned nil subscription")
				return
			}

			// Verify subscription has required fields
			if subscription.SubscriptionID == "" {
				t.Errorf("GetSubscription() missing subscription ID")
			}
			if subscription.ASIN == "" {
				t.Errorf("GetSubscription() missing ASIN")
			}
			if subscription.Title == "" {
				t.Errorf("GetSubscription() missing title")
			}
			if subscription.FrequencyWeeks <= 0 {
				t.Errorf("GetSubscription() invalid frequency weeks: %d", subscription.FrequencyWeeks)
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
			errContains:    "subscriptionID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get original subscription to compare
			var originalNextDelivery string
			if tt.subscriptionID != "" {
				original, _ := client.GetSubscription(tt.subscriptionID)
				if original != nil {
					originalNextDelivery = original.NextDelivery
				}
			}

			subscription, err := client.SkipDelivery(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SkipDelivery() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("SkipDelivery() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("SkipDelivery() unexpected error = %v", err)
				return
			}

			if subscription == nil {
				t.Errorf("SkipDelivery() returned nil subscription")
				return
			}

			// Verify next delivery was pushed forward
			if originalNextDelivery != "" && subscription.NextDelivery == originalNextDelivery {
				t.Errorf("SkipDelivery() next delivery date did not change")
			}

			// Verify next delivery is in the future
			nextDelivery, err := time.Parse("2006-01-02", subscription.NextDelivery)
			if err != nil {
				t.Errorf("SkipDelivery() invalid next delivery date format: %v", err)
				return
			}

			if nextDelivery.Before(time.Now()) {
				t.Errorf("SkipDelivery() next delivery date %s is in the past", subscription.NextDelivery)
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
			errContains:    "subscriptionID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.CancelSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CancelSubscription() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("CancelSubscription() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CancelSubscription() unexpected error = %v", err)
				return
			}

			if subscription == nil {
				t.Errorf("CancelSubscription() returned nil subscription")
				return
			}

			// Verify status is cancelled
			if subscription.Status != "cancelled" {
				t.Errorf("CancelSubscription() status = %s, want cancelled", subscription.Status)
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
		t.Errorf("GetSubscriptions() returned nil response")
		return
	}

	// Subscriptions can be empty, just verify it's not nil
	if response.Subscriptions == nil {
		t.Errorf("GetSubscriptions() returned nil subscriptions slice")
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Errorf("GetUpcomingDeliveries() unexpected error = %v", err)
		return
	}

	// Deliveries can be empty, just verify it's not nil
	if deliveries == nil {
		t.Errorf("GetUpcomingDeliveries() returned nil deliveries slice")
	}
}
