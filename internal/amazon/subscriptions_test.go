package amazon

import (
	"testing"
)

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()
	response, err := client.GetSubscriptions()

	if err != nil {
		t.Errorf("GetSubscriptions() unexpected error: %v", err)
	}

	if response == nil {
		t.Error("GetSubscriptions() returned nil response")
		return
	}

	if response.Subscriptions == nil {
		t.Error("GetSubscriptions() Subscriptions is nil")
	}
}

func TestGetSubscription(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
		errString      string
	}{
		{
			name:           "valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
			errString:      "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errString:      "subscription ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.GetSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Error("GetSubscription() expected error but got none")
				} else if err.Error() != tt.errString {
					t.Errorf("GetSubscription() error = %v, want %v", err.Error(), tt.errString)
				}
				return
			}

			if err != nil {
				t.Errorf("GetSubscription() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Error("GetSubscription() returned nil subscription")
				return
			}

			// Verify subscription fields
			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("GetSubscription() SubscriptionID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}
		})
	}
}

func TestSkipDelivery(t *testing.T) {
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
			errContains:    "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscription ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.SkipDelivery(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Error("SkipDelivery() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("SkipDelivery() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("SkipDelivery() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Error("SkipDelivery() returned nil subscription")
				return
			}

			// Verify next delivery was updated
			if subscription.NextDelivery == "" {
				t.Error("SkipDelivery() NextDelivery should not be empty")
			}
		})
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
			weeks:          8,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "minimum valid frequency (1 week)",
			subscriptionID: "S01-1234567-8901234",
			weeks:          1,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "maximum valid frequency (26 weeks)",
			subscriptionID: "S01-1234567-8901234",
			weeks:          26,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			weeks:          4,
			wantErr:        true,
			errContains:    "subscription ID cannot be empty",
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
			name:           "weeks above 26 should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "weeks above 26 (extreme) should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          52,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.wantErr {
				if err == nil {
					t.Error("UpdateFrequency() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateFrequency() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateFrequency() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Error("UpdateFrequency() returned nil subscription")
				return
			}

			// Verify frequency was updated
			if subscription.FrequencyWeeks != tt.weeks {
				t.Errorf("UpdateFrequency() FrequencyWeeks = %v, want %v", subscription.FrequencyWeeks, tt.weeks)
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
			errContains:    "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscription ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.CancelSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Error("CancelSubscription() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("CancelSubscription() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CancelSubscription() unexpected error: %v", err)
				return
			}

			if subscription == nil {
				t.Error("CancelSubscription() returned nil subscription")
				return
			}

			// Verify status was updated to cancelled
			if subscription.Status != "cancelled" {
				t.Errorf("CancelSubscription() Status = %v, want cancelled", subscription.Status)
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
