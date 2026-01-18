package amazon

import (
	"testing"
)

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
			errContains:    "",
		},
		{
			name:           "another valid subscription ID",
			subscriptionID: "S01-9876543-2109876",
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "non-existent subscription ID should fail",
			subscriptionID: "S01-0000000-0000000",
			wantErr:        true,
			errContains:    "subscription not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.GetSubscription(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Error("GetSubscription() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetSubscription() error = %v, want error containing %q", err, tt.errContains)
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

			// Verify subscription fields are populated
			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("SubscriptionID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}

			if subscription.ASIN == "" {
				t.Error("ASIN should not be empty")
			}

			if subscription.Title == "" {
				t.Error("Title should not be empty")
			}

			if subscription.Price <= 0 {
				t.Error("Price should be greater than 0")
			}

			if subscription.DiscountPercent < 0 {
				t.Error("DiscountPercent should not be negative")
			}

			if subscription.FrequencyWeeks <= 0 {
				t.Error("FrequencyWeeks should be greater than 0")
			}

			if subscription.Status == "" {
				t.Error("Status should not be empty")
			}

			if subscription.Quantity <= 0 {
				t.Error("Quantity should be greater than 0")
			}
		})
	}
}

func TestGetSubscription_ReturnsCopy(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	// Get subscription twice
	sub1, err1 := client.GetSubscription(subscriptionID)
	if err1 != nil {
		t.Fatalf("GetSubscription() error = %v", err1)
	}

	sub2, err2 := client.GetSubscription(subscriptionID)
	if err2 != nil {
		t.Fatalf("GetSubscription() error = %v", err2)
	}

	// Verify they're separate copies
	if sub1 == sub2 {
		t.Error("GetSubscription() should return a copy, not the same pointer")
	}

	// Modify one and verify the other is unchanged
	originalTitle := sub1.Title
	sub1.Title = "Modified Title"

	if sub2.Title != originalTitle {
		t.Error("Modifying one subscription should not affect another copy")
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
		t.Error("GetSubscriptions() returned nil response")
		return
	}

	if response.Subscriptions == nil {
		t.Error("GetSubscriptions() Subscriptions slice is nil")
		return
	}

	// Should have at least the test data
	if len(response.Subscriptions) == 0 {
		t.Error("GetSubscriptions() should return at least one subscription")
	}

	// Verify all subscriptions have required fields
	for i, sub := range response.Subscriptions {
		if sub.SubscriptionID == "" {
			t.Errorf("Subscription[%d] SubscriptionID is empty", i)
		}
		if sub.ASIN == "" {
			t.Errorf("Subscription[%d] ASIN is empty", i)
		}
		if sub.Title == "" {
			t.Errorf("Subscription[%d] Title is empty", i)
		}
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
			name:           "valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "non-existent subscription ID should fail",
			subscriptionID: "S01-0000000-0000000",
			wantErr:        true,
			errContains:    "subscription not found",
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
			weeks:          2,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "max frequency (26 weeks)",
			subscriptionID: "S01-1234567-8901234",
			weeks:          26,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "min frequency (1 week)",
			subscriptionID: "S01-1234567-8901234",
			weeks:          1,
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "empty subscription ID should fail",
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
			name:           "too many weeks should fail",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			wantErr:        true,
			errContains:    "frequency must be between 1 and 26 weeks",
		},
		{
			name:           "non-existent subscription ID should fail",
			subscriptionID: "S01-0000000-0000000",
			weeks:          4,
			wantErr:        true,
			errContains:    "subscription not found",
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
			name:           "valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
			errContains:    "",
		},
		{
			name:           "empty subscription ID should fail",
			subscriptionID: "",
			wantErr:        true,
			errContains:    "subscriptionID cannot be empty",
		},
		{
			name:           "non-existent subscription ID should fail",
			subscriptionID: "S01-0000000-0000000",
			wantErr:        true,
			errContains:    "subscription not found",
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

			// Verify status is set to cancelled
			if subscription.Status != "cancelled" {
				t.Errorf("Status = %v, want 'cancelled'", subscription.Status)
			}
		})
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

	// Should have at least one active subscription in test data
	if len(deliveries) == 0 {
		t.Error("GetUpcomingDeliveries() should return at least one delivery")
	}

	// Verify all deliveries have required fields
	for i, delivery := range deliveries {
		if delivery.SubscriptionID == "" {
			t.Errorf("Delivery[%d] SubscriptionID is empty", i)
		}
		if delivery.ASIN == "" {
			t.Errorf("Delivery[%d] ASIN is empty", i)
		}
		if delivery.Title == "" {
			t.Errorf("Delivery[%d] Title is empty", i)
		}
		if delivery.DeliveryDate == "" {
			t.Errorf("Delivery[%d] DeliveryDate is empty", i)
		}
		if delivery.Quantity <= 0 {
			t.Errorf("Delivery[%d] Quantity should be greater than 0", i)
		}
	}
}
