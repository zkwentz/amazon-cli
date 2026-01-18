package amazon

import (
	"testing"
)

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()
	subscriptions, err := client.GetSubscriptions()

	if err != nil {
		t.Errorf("GetSubscriptions() unexpected error: %v", err)
		return
	}

	if subscriptions == nil {
		t.Error("GetSubscriptions() returned nil")
		return
	}

	if len(subscriptions.Subscriptions) == 0 {
		t.Error("GetSubscriptions() returned empty subscriptions list")
	}

	// Verify subscription structure
	for _, sub := range subscriptions.Subscriptions {
		if sub.SubscriptionID == "" {
			t.Error("Subscription has empty ID")
		}
		if sub.ASIN == "" {
			t.Error("Subscription has empty ASIN")
		}
		if sub.Title == "" {
			t.Error("Subscription has empty Title")
		}
		if sub.Status == "" {
			t.Error("Subscription has empty Status")
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
				t.Error("GetSubscription() returned nil")
				return
			}

			// Verify subscription fields
			if subscription.SubscriptionID == "" {
				t.Error("Subscription has empty ID")
			}
			if subscription.ASIN == "" {
				t.Error("Subscription has empty ASIN")
			}
			if subscription.Title == "" {
				t.Error("Subscription has empty Title")
			}
			if subscription.Status == "" {
				t.Error("Subscription has empty Status")
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
			errString:      "subscription ID cannot be empty",
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
				t.Error("SkipDelivery() returned nil")
				return
			}

			// Verify subscription has next delivery date
			if subscription.NextDelivery == "" {
				t.Error("Subscription should have next delivery date after skip")
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
			name:           "valid frequency update - 4 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          4,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "valid frequency update - 8 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          8,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "minimum valid frequency - 1 week",
			subscriptionID: "S01-1234567-8901234",
			weeks:          1,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "maximum valid frequency - 26 weeks",
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
			name:           "too many weeks should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
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
				t.Error("UpdateFrequency() returned nil")
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
			errString:      "subscription ID cannot be empty",
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
				t.Error("CancelSubscription() returned nil")
				return
			}

			// Verify subscription status is cancelled
			if subscription.Status != "cancelled" {
				t.Errorf("CancelSubscription() Status = %v, want cancelled", subscription.Status)
			}

			// Verify next delivery is cleared
			if subscription.NextDelivery != "" {
				t.Error("CancelSubscription() NextDelivery should be empty for cancelled subscription")
			}
		})
	}
}

func TestCancelSubscription_StatusChange(t *testing.T) {
	client := NewClient()

	// Get initial subscription
	subscription, err := client.GetSubscription("S01-1234567-8901234")
	if err != nil {
		t.Fatalf("GetSubscription() error = %v", err)
	}

	initialStatus := subscription.Status
	initialNextDelivery := subscription.NextDelivery

	// Cancel subscription
	cancelledSub, err := client.CancelSubscription("S01-1234567-8901234")
	if err != nil {
		t.Fatalf("CancelSubscription() error = %v", err)
	}

	// Verify status changed
	if cancelledSub.Status == initialStatus {
		t.Error("CancelSubscription() status should change from initial status")
	}

	if cancelledSub.Status != "cancelled" {
		t.Errorf("CancelSubscription() Status = %v, want cancelled", cancelledSub.Status)
	}

	// Verify next delivery is cleared
	if cancelledSub.NextDelivery != "" {
		t.Errorf("CancelSubscription() NextDelivery = %v, want empty string (was %v)", cancelledSub.NextDelivery, initialNextDelivery)
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
		t.Error("GetUpcomingDeliveries() returned nil")
		return
	}

	// Verify delivery structure
	for _, delivery := range deliveries {
		if delivery.SubscriptionID == "" {
			t.Error("Delivery has empty SubscriptionID")
		}
		if delivery.ASIN == "" {
			t.Error("Delivery has empty ASIN")
		}
		if delivery.Title == "" {
			t.Error("Delivery has empty Title")
		}
		if delivery.DeliveryDate == "" {
			t.Error("Delivery has empty DeliveryDate")
		}
		if delivery.Quantity <= 0 {
			t.Error("Delivery has invalid Quantity")
		}
	}

	// Verify deliveries are sorted by date (basic check)
	if len(deliveries) > 1 {
		for i := 1; i < len(deliveries); i++ {
			if deliveries[i].DeliveryDate < deliveries[i-1].DeliveryDate {
				t.Error("GetUpcomingDeliveries() deliveries are not sorted by date")
				break
			}
		}
	}
}
