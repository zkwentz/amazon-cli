package amazon

import (
	"testing"
	"time"
)

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()

	resp, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions() returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("GetSubscriptions() returned nil response")
	}

	if len(resp.Subscriptions) == 0 {
		t.Error("Expected at least one subscription, got 0")
	}

	// Verify structure of first subscription
	if len(resp.Subscriptions) > 0 {
		sub := resp.Subscriptions[0]
		if sub.SubscriptionID == "" {
			t.Error("Subscription ID should not be empty")
		}
		if sub.ASIN == "" {
			t.Error("Subscription ASIN should not be empty")
		}
		if sub.Title == "" {
			t.Error("Subscription Title should not be empty")
		}
		if sub.Price <= 0 {
			t.Errorf("Expected positive price, got %f", sub.Price)
		}
		if sub.FrequencyWeeks <= 0 {
			t.Errorf("Expected positive frequency weeks, got %d", sub.FrequencyWeeks)
		}
		if sub.NextDelivery == "" {
			t.Error("NextDelivery should not be empty")
		}
		if sub.Status == "" {
			t.Error("Status should not be empty")
		}
		if sub.Quantity <= 0 {
			t.Errorf("Expected positive quantity, got %d", sub.Quantity)
		}
	}
}

func TestGetSubscription(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	sub, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("GetSubscription() returned error: %v", err)
	}

	if sub == nil {
		t.Fatal("GetSubscription() returned nil")
	}

	if sub.SubscriptionID != subscriptionID {
		t.Errorf("Expected subscription ID %s, got %s", subscriptionID, sub.SubscriptionID)
	}

	if sub.ASIN == "" {
		t.Error("Subscription ASIN should not be empty")
	}

	if sub.Title == "" {
		t.Error("Subscription Title should not be empty")
	}
}

func TestSkipDelivery(t *testing.T) {
	client := NewClient()
	subscriptionID := "S01-1234567-8901234"

	// Get original subscription
	originalSub, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("GetSubscription() returned error: %v", err)
	}

	originalDate, err := time.Parse("2006-01-02", originalSub.NextDelivery)
	if err != nil {
		t.Fatalf("Failed to parse original delivery date: %v", err)
	}

	// Skip delivery
	updatedSub, err := client.SkipDelivery(subscriptionID)
	if err != nil {
		t.Fatalf("SkipDelivery() returned error: %v", err)
	}

	if updatedSub == nil {
		t.Fatal("SkipDelivery() returned nil")
	}

	updatedDate, err := time.Parse("2006-01-02", updatedSub.NextDelivery)
	if err != nil {
		t.Fatalf("Failed to parse updated delivery date: %v", err)
	}

	// Verify the new date is later than the original
	if !updatedDate.After(originalDate) {
		t.Errorf("Expected updated delivery date to be after original date. Original: %s, Updated: %s",
			originalDate, updatedDate)
	}

	// Verify the difference is approximately the frequency weeks
	expectedDiff := time.Duration(updatedSub.FrequencyWeeks*7*24) * time.Hour
	actualDiff := updatedDate.Sub(originalDate)

	// Allow for some tolerance (within 1 day)
	tolerance := 24 * time.Hour
	if actualDiff < expectedDiff-tolerance || actualDiff > expectedDiff+tolerance {
		t.Errorf("Expected delivery to be skipped by approximately %v, got %v", expectedDiff, actualDiff)
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() returned error: %v", err)
	}

	if deliveries == nil {
		t.Fatal("GetUpcomingDeliveries() returned nil")
	}

	if len(deliveries) == 0 {
		t.Error("Expected at least one upcoming delivery, got 0")
	}

	// Verify structure of deliveries
	for i, delivery := range deliveries {
		if delivery.SubscriptionID == "" {
			t.Errorf("Delivery %d: SubscriptionID should not be empty", i)
		}
		if delivery.ASIN == "" {
			t.Errorf("Delivery %d: ASIN should not be empty", i)
		}
		if delivery.Title == "" {
			t.Errorf("Delivery %d: Title should not be empty", i)
		}
		if delivery.DeliveryDate == "" {
			t.Errorf("Delivery %d: DeliveryDate should not be empty", i)
		}
		if delivery.Quantity <= 0 {
			t.Errorf("Delivery %d: Expected positive quantity, got %d", i, delivery.Quantity)
		}

		// Verify date format is valid
		_, err := time.Parse("2006-01-02", delivery.DeliveryDate)
		if err != nil {
			t.Errorf("Delivery %d: Invalid date format %s: %v", i, delivery.DeliveryDate, err)
		}
	}
}

func TestGetUpcomingDeliveriesSortedByDate(t *testing.T) {
	client := NewClient()

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() returned error: %v", err)
	}

	if len(deliveries) < 2 {
		t.Skip("Need at least 2 deliveries to test sorting")
	}

	// Verify deliveries are sorted by date (earliest first)
	for i := 0; i < len(deliveries)-1; i++ {
		date1, err1 := time.Parse("2006-01-02", deliveries[i].DeliveryDate)
		date2, err2 := time.Parse("2006-01-02", deliveries[i+1].DeliveryDate)

		if err1 != nil || err2 != nil {
			t.Fatalf("Failed to parse dates: %v, %v", err1, err2)
		}

		if date1.After(date2) {
			t.Errorf("Deliveries not sorted correctly: delivery[%d] date %s is after delivery[%d] date %s",
				i, deliveries[i].DeliveryDate, i+1, deliveries[i+1].DeliveryDate)
		}
	}
}

func TestGetUpcomingDeliveriesOnlyActiveSubscriptions(t *testing.T) {
	client := NewClient()

	// Get all subscriptions to check status
	subscriptionsResp, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions() returned error: %v", err)
	}

	// Get upcoming deliveries
	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() returned error: %v", err)
	}

	// Count active subscriptions
	activeCount := 0
	for _, sub := range subscriptionsResp.Subscriptions {
		if sub.Status == "active" {
			activeCount++
		}
	}

	// Verify number of deliveries matches number of active subscriptions
	if len(deliveries) != activeCount {
		t.Errorf("Expected %d deliveries (matching active subscriptions), got %d",
			activeCount, len(deliveries))
	}

	// Verify all deliveries correspond to active subscriptions
	for _, delivery := range deliveries {
		found := false
		for _, sub := range subscriptionsResp.Subscriptions {
			if sub.SubscriptionID == delivery.SubscriptionID {
				if sub.Status != "active" {
					t.Errorf("Delivery for subscription %s is included but subscription status is %s, not active",
						delivery.SubscriptionID, sub.Status)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Delivery for subscription %s not found in subscriptions list", delivery.SubscriptionID)
		}
	}
}

func TestGetUpcomingDeliveriesFieldMapping(t *testing.T) {
	client := NewClient()

	// Get subscriptions and deliveries
	subscriptionsResp, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions() returned error: %v", err)
	}

	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() returned error: %v", err)
	}

	// Verify each delivery's fields match the corresponding subscription
	for _, delivery := range deliveries {
		var matchingSub *struct {
			SubscriptionID string
			ASIN           string
			Title          string
			NextDelivery   string
			Quantity       int
		}

		for _, sub := range subscriptionsResp.Subscriptions {
			if sub.SubscriptionID == delivery.SubscriptionID {
				matchingSub = &struct {
					SubscriptionID string
					ASIN           string
					Title          string
					NextDelivery   string
					Quantity       int
				}{
					SubscriptionID: sub.SubscriptionID,
					ASIN:           sub.ASIN,
					Title:          sub.Title,
					NextDelivery:   sub.NextDelivery,
					Quantity:       sub.Quantity,
				}
				break
			}
		}

		if matchingSub == nil {
			t.Errorf("Could not find subscription matching delivery %s", delivery.SubscriptionID)
			continue
		}

		// Verify fields match
		if delivery.ASIN != matchingSub.ASIN {
			t.Errorf("ASIN mismatch for subscription %s: expected %s, got %s",
				delivery.SubscriptionID, matchingSub.ASIN, delivery.ASIN)
		}

		if delivery.Title != matchingSub.Title {
			t.Errorf("Title mismatch for subscription %s: expected %s, got %s",
				delivery.SubscriptionID, matchingSub.Title, delivery.Title)
		}

		if delivery.DeliveryDate != matchingSub.NextDelivery {
			t.Errorf("DeliveryDate mismatch for subscription %s: expected %s, got %s",
				delivery.SubscriptionID, matchingSub.NextDelivery, delivery.DeliveryDate)
		}

		if delivery.Quantity != matchingSub.Quantity {
			t.Errorf("Quantity mismatch for subscription %s: expected %d, got %d",
				delivery.SubscriptionID, matchingSub.Quantity, delivery.Quantity)
		}
	}
}
