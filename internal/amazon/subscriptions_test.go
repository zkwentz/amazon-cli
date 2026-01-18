package amazon

import (
	"testing"
)

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()

	response, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions failed: %v", err)
	}

	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if len(response.Subscriptions) == 0 {
		t.Fatal("Expected at least one subscription in mock data")
	}

	// Verify structure of first subscription
	sub := response.Subscriptions[0]
	if sub.SubscriptionID == "" {
		t.Error("Expected non-empty subscription ID")
	}
	if sub.ASIN == "" {
		t.Error("Expected non-empty ASIN")
	}
	if sub.Title == "" {
		t.Error("Expected non-empty title")
	}
	if sub.Price <= 0 {
		t.Error("Expected positive price")
	}
	if sub.FrequencyWeeks <= 0 {
		t.Error("Expected positive frequency")
	}
	if sub.NextDelivery == "" {
		t.Error("Expected non-empty next delivery date")
	}
	if sub.Status == "" {
		t.Error("Expected non-empty status")
	}
}

func TestGetSubscription(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name           string
		subscriptionID string
		expectError    bool
	}{
		{
			name:           "Valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			expectError:    false,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.GetSubscription(tt.subscriptionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetSubscription failed: %v", err)
			}

			if subscription == nil {
				t.Fatal("Expected non-nil subscription")
			}

			if subscription.SubscriptionID != tt.subscriptionID {
				t.Errorf("Expected subscription ID %s, got %s", tt.subscriptionID, subscription.SubscriptionID)
			}
		})
	}
}

func TestSkipDelivery(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name           string
		subscriptionID string
		expectError    bool
	}{
		{
			name:           "Valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			expectError:    false,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get original subscription
			originalSub, _ := client.GetSubscription(tt.subscriptionID)

			subscription, err := client.SkipDelivery(tt.subscriptionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("SkipDelivery failed: %v", err)
			}

			if subscription == nil {
				t.Fatal("Expected non-nil subscription")
			}

			// Verify that the next delivery date was updated
			if originalSub != nil && subscription.NextDelivery == originalSub.NextDelivery {
				t.Error("Expected next delivery date to be updated after skip")
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
		expectError    bool
	}{
		{
			name:           "Valid frequency update",
			subscriptionID: "S01-1234567-8901234",
			weeks:          8,
			expectError:    false,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			weeks:          4,
			expectError:    true,
		},
		{
			name:           "Invalid frequency (too low)",
			subscriptionID: "S01-1234567-8901234",
			weeks:          0,
			expectError:    true,
		},
		{
			name:           "Invalid frequency (too high)",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			expectError:    true,
		},
		{
			name:           "Valid minimum frequency",
			subscriptionID: "S01-1234567-8901234",
			weeks:          1,
			expectError:    false,
		},
		{
			name:           "Valid maximum frequency",
			subscriptionID: "S01-1234567-8901234",
			weeks:          26,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("UpdateFrequency failed: %v", err)
			}

			if subscription == nil {
				t.Fatal("Expected non-nil subscription")
			}

			if subscription.FrequencyWeeks != tt.weeks {
				t.Errorf("Expected frequency %d, got %d", tt.weeks, subscription.FrequencyWeeks)
			}
		})
	}
}

func TestCancelSubscription(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name           string
		subscriptionID string
		expectError    bool
	}{
		{
			name:           "Valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			expectError:    false,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.CancelSubscription(tt.subscriptionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("CancelSubscription failed: %v", err)
			}

			if subscription == nil {
				t.Fatal("Expected non-nil subscription")
			}

			if subscription.Status != "cancelled" {
				t.Errorf("Expected status 'cancelled', got '%s'", subscription.Status)
			}
		})
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries failed: %v", err)
	}

	if deliveries == nil {
		t.Fatal("Expected non-nil deliveries")
	}

	if len(deliveries) == 0 {
		t.Fatal("Expected at least one upcoming delivery in mock data")
	}

	// Verify structure of first delivery
	delivery := deliveries[0]
	if delivery.SubscriptionID == "" {
		t.Error("Expected non-empty subscription ID")
	}
	if delivery.ASIN == "" {
		t.Error("Expected non-empty ASIN")
	}
	if delivery.Title == "" {
		t.Error("Expected non-empty title")
	}
	if delivery.DeliveryDate == "" {
		t.Error("Expected non-empty delivery date")
	}
	if delivery.Quantity <= 0 {
		t.Error("Expected positive quantity")
	}

	// Verify deliveries are sorted by date
	if len(deliveries) > 1 {
		for i := 0; i < len(deliveries)-1; i++ {
			if deliveries[i].DeliveryDate > deliveries[i+1].DeliveryDate {
				t.Error("Expected deliveries to be sorted by date")
			}
		}
	}
}
