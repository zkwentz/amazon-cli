package amazon

import (
	"testing"
)

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
		t.Error("GetSubscriptions() Subscriptions is nil")
		return
	}

	// Verify we have mock subscriptions
	if len(response.Subscriptions) == 0 {
		t.Error("GetSubscriptions() returned empty subscriptions list")
		return
	}

	// Verify first subscription has required fields
	firstSub := response.Subscriptions[0]

	if firstSub.SubscriptionID == "" {
		t.Error("Subscription SubscriptionID is empty")
	}

	if firstSub.ASIN == "" {
		t.Error("Subscription ASIN is empty")
	}

	if firstSub.Title == "" {
		t.Error("Subscription Title is empty")
	}

	if firstSub.Price <= 0 {
		t.Errorf("Subscription Price = %v, want > 0", firstSub.Price)
	}

	if firstSub.DiscountPercent < 0 {
		t.Errorf("Subscription DiscountPercent = %v, want >= 0", firstSub.DiscountPercent)
	}

	if firstSub.FrequencyWeeks <= 0 {
		t.Errorf("Subscription FrequencyWeeks = %v, want > 0", firstSub.FrequencyWeeks)
	}

	if firstSub.NextDelivery == "" {
		t.Error("Subscription NextDelivery is empty")
	}

	if firstSub.Status == "" {
		t.Error("Subscription Status is empty")
	}

	if firstSub.Quantity <= 0 {
		t.Errorf("Subscription Quantity = %v, want > 0", firstSub.Quantity)
	}
}

func TestGetSubscriptions_ValidStatuses(t *testing.T) {
	client := NewClient()

	response, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions() error: %v", err)
	}

	validStatuses := map[string]bool{
		"active":    true,
		"paused":    true,
		"cancelled": true,
	}

	for _, sub := range response.Subscriptions {
		if !validStatuses[sub.Status] {
			t.Errorf("Subscription has invalid status: %v, want one of: active, paused, cancelled", sub.Status)
		}
	}
}

func TestGetSubscriptions_DiscountValidation(t *testing.T) {
	client := NewClient()

	response, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions() error: %v", err)
	}

	for _, sub := range response.Subscriptions {
		if sub.DiscountPercent < 0 || sub.DiscountPercent > 100 {
			t.Errorf("Subscription %s has invalid discount: %v, want between 0 and 100",
				sub.SubscriptionID, sub.DiscountPercent)
		}
	}
}

func TestGetSubscription(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
	}{
		{
			name:           "valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "another valid subscription ID",
			subscriptionID: "S01-9999999-9999999",
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
				t.Errorf("Subscription ID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}

			if subscription.ASIN == "" {
				t.Error("Subscription ASIN is empty")
			}

			if subscription.Title == "" {
				t.Error("Subscription Title is empty")
			}

			if subscription.Price <= 0 {
				t.Error("Subscription Price should be greater than 0")
			}
		})
	}
}

func TestSkipDelivery(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	subscription, err := client.SkipDelivery(subscriptionID)

	if err != nil {
		t.Errorf("SkipDelivery() unexpected error: %v", err)
		return
	}

	if subscription == nil {
		t.Error("SkipDelivery() returned nil subscription")
		return
	}

	if subscription.SubscriptionID != subscriptionID {
		t.Errorf("Subscription ID = %v, want %v", subscription.SubscriptionID, subscriptionID)
	}
}

func TestUpdateFrequency(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		weeks          int
		wantErr        bool
	}{
		{
			name:           "valid frequency update - 4 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          4,
			wantErr:        false,
		},
		{
			name:           "valid frequency update - 8 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          8,
			wantErr:        false,
		},
		{
			name:           "valid frequency update - 12 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          12,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			subscription, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.wantErr {
				if err == nil {
					t.Error("UpdateFrequency() expected error but got none")
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

			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("Subscription ID = %v, want %v", subscription.SubscriptionID, tt.subscriptionID)
			}
		})
	}
}

func TestCancelSubscription(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	subscription, err := client.CancelSubscription(subscriptionID)

	if err != nil {
		t.Errorf("CancelSubscription() unexpected error: %v", err)
		return
	}

	if subscription == nil {
		t.Error("CancelSubscription() returned nil subscription")
		return
	}

	if subscription.SubscriptionID != subscriptionID {
		t.Errorf("Subscription ID = %v, want %v", subscription.SubscriptionID, subscriptionID)
	}

	if subscription.Status != "cancelled" {
		t.Errorf("Subscription Status = %v, want 'cancelled'", subscription.Status)
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

	// Verify we have mock deliveries
	if len(deliveries) == 0 {
		t.Error("GetUpcomingDeliveries() returned empty deliveries list")
		return
	}

	// Verify first delivery has required fields
	firstDelivery := deliveries[0]

	if firstDelivery.SubscriptionID == "" {
		t.Error("UpcomingDelivery SubscriptionID is empty")
	}

	if firstDelivery.ASIN == "" {
		t.Error("UpcomingDelivery ASIN is empty")
	}

	if firstDelivery.Title == "" {
		t.Error("UpcomingDelivery Title is empty")
	}

	if firstDelivery.DeliveryDate == "" {
		t.Error("UpcomingDelivery DeliveryDate is empty")
	}

	if firstDelivery.Quantity <= 0 {
		t.Errorf("UpcomingDelivery Quantity = %v, want > 0", firstDelivery.Quantity)
	}
}

func TestGetUpcomingDeliveries_AllHaveValidData(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() error: %v", err)
	}

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
			t.Errorf("Delivery[%d] Quantity = %v, want > 0", i, delivery.Quantity)
		}
	}
}
