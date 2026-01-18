package amazon

import (
	"testing"
)

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()

	subscriptions, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions failed: %v", err)
	}

	if subscriptions == nil {
		t.Fatal("GetSubscriptions returned nil")
	}

	if len(subscriptions.Subscriptions) == 0 {
		t.Error("GetSubscriptions returned empty subscriptions list")
	}

	// Verify first subscription has required fields
	if len(subscriptions.Subscriptions) > 0 {
		sub := subscriptions.Subscriptions[0]
		if sub.SubscriptionID == "" {
			t.Error("Subscription is missing SubscriptionID")
		}
		if sub.ASIN == "" {
			t.Error("Subscription is missing ASIN")
		}
		if sub.Title == "" {
			t.Error("Subscription is missing Title")
		}
		if sub.Price <= 0 {
			t.Error("Subscription has invalid Price")
		}
		if sub.FrequencyWeeks <= 0 {
			t.Error("Subscription has invalid FrequencyWeeks")
		}
		if sub.Status == "" {
			t.Error("Subscription is missing Status")
		}
		if sub.Quantity <= 0 {
			t.Error("Subscription has invalid Quantity")
		}
	}
}

func TestGetSubscription(t *testing.T) {
	client := NewClient()

	subscriptionID := "S01-1234567-8901234"
	subscription, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("GetSubscription failed: %v", err)
	}

	if subscription == nil {
		t.Fatal("GetSubscription returned nil")
	}

	if subscription.SubscriptionID != subscriptionID {
		t.Errorf("Expected subscription ID %s, got %s", subscriptionID, subscription.SubscriptionID)
	}

	if subscription.ASIN == "" {
		t.Error("Subscription is missing ASIN")
	}
	if subscription.Title == "" {
		t.Error("Subscription is missing Title")
	}
}

func TestSkipDelivery(t *testing.T) {
	client := NewClient()

	subscriptionID := "S01-1234567-8901234"
	subscription, err := client.SkipDelivery(subscriptionID)
	if err != nil {
		t.Fatalf("SkipDelivery failed: %v", err)
	}

	if subscription == nil {
		t.Fatal("SkipDelivery returned nil")
	}

	if subscription.SubscriptionID != subscriptionID {
		t.Errorf("Expected subscription ID %s, got %s", subscriptionID, subscription.SubscriptionID)
	}

	// Verify the subscription is still active after skipping
	if subscription.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", subscription.Status)
	}
}

func TestUpdateFrequency(t *testing.T) {
	client := NewClient()

	subscriptionID := "S01-1234567-8901234"
	newFrequency := 8

	subscription, err := client.UpdateFrequency(subscriptionID, newFrequency)
	if err != nil {
		t.Fatalf("UpdateFrequency failed: %v", err)
	}

	if subscription == nil {
		t.Fatal("UpdateFrequency returned nil")
	}

	if subscription.FrequencyWeeks != newFrequency {
		t.Errorf("Expected frequency %d, got %d", newFrequency, subscription.FrequencyWeeks)
	}
}

func TestCancelSubscription(t *testing.T) {
	client := NewClient()

	subscriptionID := "S01-1234567-8901234"
	subscription, err := client.CancelSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("CancelSubscription failed: %v", err)
	}

	if subscription == nil {
		t.Fatal("CancelSubscription returned nil")
	}

	if subscription.Status != "cancelled" {
		t.Errorf("Expected status 'cancelled', got '%s'", subscription.Status)
	}

	// Next delivery should be empty for cancelled subscription
	if subscription.NextDelivery != "" {
		t.Errorf("Expected empty NextDelivery for cancelled subscription, got '%s'", subscription.NextDelivery)
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

	if len(deliveries) == 0 {
		t.Error("GetUpcomingDeliveries returned empty deliveries list")
	}

	// Verify first delivery has required fields
	if len(deliveries) > 0 {
		delivery := deliveries[0]
		if delivery.SubscriptionID == "" {
			t.Error("Delivery is missing SubscriptionID")
		}
		if delivery.ASIN == "" {
			t.Error("Delivery is missing ASIN")
		}
		if delivery.Title == "" {
			t.Error("Delivery is missing Title")
		}
		if delivery.DeliveryDate == "" {
			t.Error("Delivery is missing DeliveryDate")
		}
		if delivery.Quantity <= 0 {
			t.Error("Delivery has invalid Quantity")
		}
	}
}
