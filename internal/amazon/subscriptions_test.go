package amazon

import (
	"testing"
	"time"
)

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()

	resp, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions failed: %v", err)
	}

	if resp == nil {
		t.Fatal("GetSubscriptions returned nil response")
	}

	if len(resp.Subscriptions) == 0 {
		t.Error("Expected at least one subscription, got 0")
	}

	// Verify subscription structure
	for _, sub := range resp.Subscriptions {
		if sub.SubscriptionID == "" {
			t.Error("Subscription has empty ID")
		}
		if sub.ASIN == "" {
			t.Error("Subscription has empty ASIN")
		}
		if sub.Title == "" {
			t.Error("Subscription has empty title")
		}
		if sub.Price <= 0 {
			t.Errorf("Subscription has invalid price: %f", sub.Price)
		}
		if sub.FrequencyWeeks <= 0 {
			t.Errorf("Subscription has invalid frequency: %d", sub.FrequencyWeeks)
		}
		if sub.NextDelivery == "" {
			t.Error("Subscription has empty next delivery date")
		}
		if sub.Status == "" {
			t.Error("Subscription has empty status")
		}
		if sub.Quantity <= 0 {
			t.Errorf("Subscription has invalid quantity: %d", sub.Quantity)
		}
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
			sub, err := client.GetSubscription(tt.subscriptionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetSubscription failed: %v", err)
			}

			if sub == nil {
				t.Fatal("GetSubscription returned nil subscription")
			}

			if sub.SubscriptionID != tt.subscriptionID {
				t.Errorf("Expected subscription ID %s, got %s", tt.subscriptionID, sub.SubscriptionID)
			}

			if sub.ASIN == "" {
				t.Error("Subscription has empty ASIN")
			}
			if sub.Title == "" {
				t.Error("Subscription has empty title")
			}
			if sub.Price <= 0 {
				t.Errorf("Subscription has invalid price: %f", sub.Price)
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
			originalSub, err := client.GetSubscription(tt.subscriptionID)
			if err != nil && !tt.expectError {
				t.Fatalf("Failed to get original subscription: %v", err)
			}

			// Skip delivery
			sub, err := client.SkipDelivery(tt.subscriptionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("SkipDelivery failed: %v", err)
			}

			if sub == nil {
				t.Fatal("SkipDelivery returned nil subscription")
			}

			// Verify the next delivery date was pushed forward
			originalDate, err := time.Parse("2006-01-02", originalSub.NextDelivery)
			if err != nil {
				t.Fatalf("Failed to parse original date: %v", err)
			}

			newDate, err := time.Parse("2006-01-02", sub.NextDelivery)
			if err != nil {
				t.Fatalf("Failed to parse new date: %v", err)
			}

			expectedDate := originalDate.AddDate(0, 0, originalSub.FrequencyWeeks*7)
			if !newDate.Equal(expectedDate) {
				t.Errorf("Expected next delivery %s, got %s", expectedDate.Format("2006-01-02"), sub.NextDelivery)
			}

			if newDate.Before(originalDate) || newDate.Equal(originalDate) {
				t.Errorf("Next delivery date should be after original date. Original: %s, New: %s", originalSub.NextDelivery, sub.NextDelivery)
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
			name:           "Valid frequency - 4 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          4,
			expectError:    false,
		},
		{
			name:           "Valid frequency - 8 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          8,
			expectError:    false,
		},
		{
			name:           "Valid frequency - 1 week (minimum)",
			subscriptionID: "S01-1234567-8901234",
			weeks:          1,
			expectError:    false,
		},
		{
			name:           "Valid frequency - 26 weeks (maximum)",
			subscriptionID: "S01-1234567-8901234",
			weeks:          26,
			expectError:    false,
		},
		{
			name:           "Invalid frequency - 0 weeks",
			subscriptionID: "S01-1234567-8901234",
			weeks:          0,
			expectError:    true,
		},
		{
			name:           "Invalid frequency - negative",
			subscriptionID: "S01-1234567-8901234",
			weeks:          -1,
			expectError:    true,
		},
		{
			name:           "Invalid frequency - too high",
			subscriptionID: "S01-1234567-8901234",
			weeks:          27,
			expectError:    true,
		},
		{
			name:           "Empty subscription ID",
			subscriptionID: "",
			weeks:          4,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := client.UpdateFrequency(tt.subscriptionID, tt.weeks)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("UpdateFrequency failed: %v", err)
			}

			if sub == nil {
				t.Fatal("UpdateFrequency returned nil subscription")
			}

			if sub.FrequencyWeeks != tt.weeks {
				t.Errorf("Expected frequency %d weeks, got %d", tt.weeks, sub.FrequencyWeeks)
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
			sub, err := client.CancelSubscription(tt.subscriptionID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("CancelSubscription failed: %v", err)
			}

			if sub == nil {
				t.Fatal("CancelSubscription returned nil subscription")
			}

			if sub.Status != "cancelled" {
				t.Errorf("Expected status 'cancelled', got '%s'", sub.Status)
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
		t.Fatal("GetUpcomingDeliveries returned nil")
	}

	// Verify delivery structure
	for _, delivery := range deliveries {
		if delivery.SubscriptionID == "" {
			t.Error("Delivery has empty subscription ID")
		}
		if delivery.ASIN == "" {
			t.Error("Delivery has empty ASIN")
		}
		if delivery.Title == "" {
			t.Error("Delivery has empty title")
		}
		if delivery.DeliveryDate == "" {
			t.Error("Delivery has empty delivery date")
		}
		if delivery.Quantity <= 0 {
			t.Errorf("Delivery has invalid quantity: %d", delivery.Quantity)
		}

		// Verify date format
		_, err := time.Parse("2006-01-02", delivery.DeliveryDate)
		if err != nil {
			t.Errorf("Delivery has invalid date format: %v", err)
		}
	}

	// Verify deliveries are sorted by date (earliest first)
	if len(deliveries) > 1 {
		for i := 0; i < len(deliveries)-1; i++ {
			date1, err1 := time.Parse("2006-01-02", deliveries[i].DeliveryDate)
			date2, err2 := time.Parse("2006-01-02", deliveries[i+1].DeliveryDate)

			if err1 != nil || err2 != nil {
				continue
			}

			if date1.After(date2) {
				t.Errorf("Deliveries not sorted correctly: %s comes before %s", deliveries[i].DeliveryDate, deliveries[i+1].DeliveryDate)
			}
		}
	}
}

func TestSubscriptionWorkflow(t *testing.T) {
	client := NewClient()

	// 1. Get all subscriptions
	subsResp, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions failed: %v", err)
	}

	if len(subsResp.Subscriptions) == 0 {
		t.Fatal("No subscriptions available for workflow test")
	}

	subscriptionID := subsResp.Subscriptions[0].SubscriptionID

	// 2. Get specific subscription
	sub, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("GetSubscription failed: %v", err)
	}

	originalFrequency := sub.FrequencyWeeks
	originalNextDelivery := sub.NextDelivery

	// 3. Update frequency
	newFrequency := 8
	updatedSub, err := client.UpdateFrequency(subscriptionID, newFrequency)
	if err != nil {
		t.Fatalf("UpdateFrequency failed: %v", err)
	}

	if updatedSub.FrequencyWeeks != newFrequency {
		t.Errorf("Frequency not updated. Expected %d, got %d", newFrequency, updatedSub.FrequencyWeeks)
	}

	// 4. Skip delivery
	skippedSub, err := client.SkipDelivery(subscriptionID)
	if err != nil {
		t.Fatalf("SkipDelivery failed: %v", err)
	}

	if skippedSub.NextDelivery == originalNextDelivery {
		t.Error("Next delivery date should have changed after skip")
	}

	// 5. Get upcoming deliveries
	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries failed: %v", err)
	}

	if len(deliveries) == 0 {
		t.Error("Expected upcoming deliveries, got none")
	}

	// 6. Cancel subscription
	cancelledSub, err := client.CancelSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("CancelSubscription failed: %v", err)
	}

	if cancelledSub.Status != "cancelled" {
		t.Errorf("Expected status 'cancelled', got '%s'", cancelledSub.Status)
	}

	// Restore original frequency for cleanup
	_, _ = client.UpdateFrequency(subscriptionID, originalFrequency)
}
