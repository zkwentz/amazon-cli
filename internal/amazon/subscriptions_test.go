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
		t.Error("GetSubscriptions() Subscriptions slice is nil")
		return
	}

	// Mock implementation returns 2 subscriptions
	if len(response.Subscriptions) != 2 {
		t.Errorf("GetSubscriptions() expected 2 subscriptions, got %d", len(response.Subscriptions))
	}

	// Verify first subscription structure
	if len(response.Subscriptions) > 0 {
		sub := response.Subscriptions[0]
		if sub.SubscriptionID == "" {
			t.Error("GetSubscriptions() subscription SubscriptionID is empty")
		}
		if sub.ASIN == "" {
			t.Error("GetSubscriptions() subscription ASIN is empty")
		}
		if sub.Title == "" {
			t.Error("GetSubscriptions() subscription Title is empty")
		}
		if sub.Price <= 0 {
			t.Error("GetSubscriptions() subscription Price should be greater than 0")
		}
		if sub.FrequencyWeeks <= 0 {
			t.Error("GetSubscriptions() subscription FrequencyWeeks should be greater than 0")
		}
		if sub.Status == "" {
			t.Error("GetSubscriptions() subscription Status is empty")
		}
		if sub.Quantity <= 0 {
			t.Error("GetSubscriptions() subscription Quantity should be greater than 0")
		}
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
			errString:      "subscriptionID cannot be empty",
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
			if subscription.ASIN == "" {
				t.Error("GetSubscription() ASIN is empty")
			}
			if subscription.Title == "" {
				t.Error("GetSubscription() Title is empty")
			}
			if subscription.Price <= 0 {
				t.Error("GetSubscription() Price should be greater than 0")
			}
		})
	}
}

func TestSkipDelivery(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
		errString      string
	}{
		{
			name:           "valid skip delivery",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
			errString:      "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errString:      "subscriptionID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.SkipDelivery(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Error("SkipDelivery() expected error but got none")
				} else if err.Error() != tt.errString {
					t.Errorf("SkipDelivery() error = %v, want %v", err.Error(), tt.errString)
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

			// Verify subscription has updated next delivery
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
		errString      string
	}{
		{
			name:           "valid frequency update",
			subscriptionID: "S01-1234567-8901234",
			weeks:          4,
			wantErr:        false,
			errString:      "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			weeks:          4,
			wantErr:        true,
			errString:      "subscriptionID cannot be empty",
		},
		{
			name:           "zero weeks should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          0,
			wantErr:        true,
			errString:      "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "negative weeks should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          -1,
			wantErr:        true,
			errString:      "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "weeks exceeding maximum should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			wantErr:        true,
			errString:      "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "maximum valid weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          26,
			wantErr:        false,
			errString:      "",
		},
		{
			name:           "minimum valid weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          1,
			wantErr:        false,
			errString:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.wantErr {
				if err == nil {
					t.Error("UpdateFrequency() expected error but got none")
				} else if err.Error() != tt.errString {
					t.Errorf("UpdateFrequency() error = %v, want %v", err.Error(), tt.errString)
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
		errString      string
	}{
		{
			name:           "valid cancellation",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
			errString:      "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errString:      "subscriptionID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.CancelSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Error("CancelSubscription() expected error but got none")
				} else if err.Error() != tt.errString {
					t.Errorf("CancelSubscription() error = %v, want %v", err.Error(), tt.errString)
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

			// Verify subscription is cancelled
			if subscription.Status != "cancelled" {
				t.Errorf("CancelSubscription() Status = %v, want 'cancelled'", subscription.Status)
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
		return
	}

	// Mock implementation returns 2 upcoming deliveries
	if len(deliveries) != 2 {
		t.Errorf("GetUpcomingDeliveries() expected 2 deliveries, got %d", len(deliveries))
	}

	// Verify first delivery structure
	if len(deliveries) > 0 {
		delivery := deliveries[0]
		if delivery.SubscriptionID == "" {
			t.Error("GetUpcomingDeliveries() delivery SubscriptionID is empty")
		}
		if delivery.ASIN == "" {
			t.Error("GetUpcomingDeliveries() delivery ASIN is empty")
		}
		if delivery.Title == "" {
			t.Error("GetUpcomingDeliveries() delivery Title is empty")
		}
		if delivery.DeliveryDate == "" {
			t.Error("GetUpcomingDeliveries() delivery DeliveryDate is empty")
		}
		if delivery.Quantity <= 0 {
			t.Error("GetUpcomingDeliveries() delivery Quantity should be greater than 0")
		}
	}
}
