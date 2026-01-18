package amazon

import (
	"testing"
)

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()

	response, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions returned error: %v", err)
	}

	if response == nil {
		t.Fatal("GetSubscriptions returned nil response")
	}

	if len(response.Subscriptions) == 0 {
		t.Fatal("GetSubscriptions returned empty subscriptions list")
	}

	// Verify first subscription has required fields
	sub := response.Subscriptions[0]
	if sub.SubscriptionID == "" {
		t.Error("Subscription missing SubscriptionID")
	}
	if sub.ASIN == "" {
		t.Error("Subscription missing ASIN")
	}
	if sub.Title == "" {
		t.Error("Subscription missing Title")
	}
	if sub.Price <= 0 {
		t.Error("Subscription has invalid Price")
	}
	if sub.FrequencyWeeks <= 0 {
		t.Error("Subscription has invalid FrequencyWeeks")
	}
	if sub.Status == "" {
		t.Error("Subscription missing Status")
	}
}

func TestGetSubscription(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	subscription, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("GetSubscription returned error: %v", err)
	}

	if subscription == nil {
		t.Fatal("GetSubscription returned nil subscription")
	}

	if subscription.SubscriptionID != subscriptionID {
		t.Errorf("Expected subscription ID %s, got %s", subscriptionID, subscription.SubscriptionID)
	}

	if subscription.Title == "" {
		t.Error("Subscription missing Title")
	}
}

func TestSkipDelivery(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	subscription, err := client.SkipDelivery(subscriptionID)
	if err != nil {
		t.Fatalf("SkipDelivery returned error: %v", err)
	}

	if subscription == nil {
		t.Fatal("SkipDelivery returned nil subscription")
	}

	if subscription.SubscriptionID != subscriptionID {
		t.Errorf("Expected subscription ID %s, got %s", subscriptionID, subscription.SubscriptionID)
	}

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
		t.Fatalf("UpdateFrequency returned error: %v", err)
	}

	if subscription == nil {
		t.Fatal("UpdateFrequency returned nil subscription")
	}

	if subscription.FrequencyWeeks != newFrequency {
		t.Errorf("Expected frequency %d weeks, got %d", newFrequency, subscription.FrequencyWeeks)
	}
}

func TestCancelSubscription(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	subscription, err := client.CancelSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("CancelSubscription returned error: %v", err)
	}

	if subscription == nil {
		t.Fatal("CancelSubscription returned nil subscription")
	}

	if subscription.Status != "cancelled" {
		t.Errorf("Expected status 'cancelled', got '%s'", subscription.Status)
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries returned error: %v", err)
	}

	if deliveries == nil {
		t.Fatal("GetUpcomingDeliveries returned nil")
	}

	if len(deliveries) == 0 {
		t.Fatal("GetUpcomingDeliveries returned empty list")
	}

	// Verify first delivery has required fields
	delivery := deliveries[0]
	if delivery.SubscriptionID == "" {
		t.Error("Delivery missing SubscriptionID")
	}
	if delivery.ASIN == "" {
		t.Error("Delivery missing ASIN")
	}
	if delivery.Title == "" {
		t.Error("Delivery missing Title")
	}
	if delivery.DeliveryDate == "" {
		t.Error("Delivery missing DeliveryDate")
	}
	if delivery.Quantity <= 0 {
		t.Error("Delivery has invalid Quantity")
	}
}

func TestGetUpcomingDeliveriesSorting(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries returned error: %v", err)
	}

	if len(deliveries) < 2 {
		t.Skip("Not enough deliveries to test sorting")
	}

	// Verify deliveries are sorted by date (in test data they should be)
	// In a real implementation, we'd verify dates are in ascending order
	for i := 0; i < len(deliveries)-1; i++ {
		if deliveries[i].DeliveryDate > deliveries[i+1].DeliveryDate {
			t.Errorf("Deliveries not sorted: %s comes after %s",
				deliveries[i].DeliveryDate, deliveries[i+1].DeliveryDate)
		}
	}
}
