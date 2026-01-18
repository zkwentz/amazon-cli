package amazon

import (
	"testing"
	"time"
)

func TestGetSubscriptions(t *testing.T) {
	client := NewClient()
	resp, err := client.GetSubscriptions()

	if err != nil {
		t.Fatalf("GetSubscriptions() error = %v", err)
	}

	if resp == nil {
		t.Fatal("GetSubscriptions() returned nil response")
	}

	if len(resp.Subscriptions) == 0 {
		t.Error("GetSubscriptions() returned empty subscriptions list")
	}

	// Verify subscription structure
	for _, sub := range resp.Subscriptions {
		if sub.SubscriptionID == "" {
			t.Error("Subscription has empty SubscriptionID")
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
	client := NewClient()
	sub, err := client.GetSubscription("S01-1234567-8901234")

	if err != nil {
		t.Fatalf("GetSubscription() error = %v", err)
	}

	if sub == nil {
		t.Fatal("GetSubscription() returned nil")
	}

	if sub.SubscriptionID == "" {
		t.Error("SubscriptionID is empty")
	}
}

func TestSkipDelivery(t *testing.T) {
	client := NewClient()
	sub, err := client.SkipDelivery("S01-1234567-8901234")

	if err != nil {
		t.Fatalf("SkipDelivery() error = %v", err)
	}

	if sub == nil {
		t.Fatal("SkipDelivery() returned nil")
	}

	if sub.Status != "active" {
		t.Errorf("SkipDelivery() status = %v, want active", sub.Status)
	}
}

func TestUpdateFrequency(t *testing.T) {
	tests := []struct {
		name  string
		weeks int
	}{
		{"4 weeks", 4},
		{"8 weeks", 8},
		{"12 weeks", 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			sub, err := client.UpdateFrequency("S01-1234567-8901234", tt.weeks)

			if err != nil {
				t.Fatalf("UpdateFrequency() error = %v", err)
			}

			if sub == nil {
				t.Fatal("UpdateFrequency() returned nil")
			}

			if sub.FrequencyWeeks != tt.weeks {
				t.Errorf("FrequencyWeeks = %v, want %v", sub.FrequencyWeeks, tt.weeks)
			}
		})
	}
}

func TestCancelSubscription(t *testing.T) {
	client := NewClient()
	sub, err := client.CancelSubscription("S01-1234567-8901234")

	if err != nil {
		t.Fatalf("CancelSubscription() error = %v", err)
	}

	if sub == nil {
		t.Fatal("CancelSubscription() returned nil")
	}

	if sub.Status != "cancelled" {
		t.Errorf("Status = %v, want cancelled", sub.Status)
	}
}

func TestGetUpcomingDeliveries(t *testing.T) {
	client := NewClient()
	deliveries, err := client.GetUpcomingDeliveries()

	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() error = %v", err)
	}

	if deliveries == nil {
		t.Fatal("GetUpcomingDeliveries() returned nil")
	}

	// Should have at least one delivery
	if len(deliveries) == 0 {
		t.Fatal("GetUpcomingDeliveries() returned empty list")
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
}

func TestGetUpcomingDeliveries_Sorting(t *testing.T) {
	client := NewClient()
	deliveries, err := client.GetUpcomingDeliveries()

	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() error = %v", err)
	}

	if len(deliveries) < 2 {
		t.Skip("Need at least 2 deliveries to test sorting")
	}

	// Verify deliveries are sorted by date (earliest first)
	for i := 0; i < len(deliveries)-1; i++ {
		currentDate, err1 := time.Parse("2006-01-02", deliveries[i].DeliveryDate)
		nextDate, err2 := time.Parse("2006-01-02", deliveries[i+1].DeliveryDate)

		if err1 != nil {
			t.Fatalf("Failed to parse date %s: %v", deliveries[i].DeliveryDate, err1)
		}
		if err2 != nil {
			t.Fatalf("Failed to parse date %s: %v", deliveries[i+1].DeliveryDate, err2)
		}

		if currentDate.After(nextDate) {
			t.Errorf("Deliveries not sorted correctly: %s comes after %s",
				deliveries[i].DeliveryDate, deliveries[i+1].DeliveryDate)
		}
	}
}

func TestGetUpcomingDeliveries_DateOrdering(t *testing.T) {
	client := NewClient()
	deliveries, err := client.GetUpcomingDeliveries()

	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() error = %v", err)
	}

	if len(deliveries) < 2 {
		t.Skip("Need at least 2 deliveries to test date ordering")
	}

	// Verify that the first delivery is earlier than the last delivery
	firstDate, err1 := time.Parse("2006-01-02", deliveries[0].DeliveryDate)
	lastDate, err2 := time.Parse("2006-01-02", deliveries[len(deliveries)-1].DeliveryDate)

	if err1 != nil {
		t.Fatalf("Failed to parse first date: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("Failed to parse last date: %v", err2)
	}

	if firstDate.After(lastDate) {
		t.Errorf("First delivery (%s) should be before or equal to last delivery (%s)",
			deliveries[0].DeliveryDate, deliveries[len(deliveries)-1].DeliveryDate)
	}
}

func TestGetUpcomingDeliveries_OnlyActiveSubscriptions(t *testing.T) {
	client := NewClient()

	// Get all subscriptions first
	resp, err := client.GetSubscriptions()
	if err != nil {
		t.Fatalf("GetSubscriptions() error = %v", err)
	}

	// Get upcoming deliveries
	deliveries, err := client.GetUpcomingDeliveries()
	if err != nil {
		t.Fatalf("GetUpcomingDeliveries() error = %v", err)
	}

	// Count active subscriptions with next delivery dates
	activeWithDelivery := 0
	for _, sub := range resp.Subscriptions {
		if sub.Status == "active" && sub.NextDelivery != "" {
			activeWithDelivery++
		}
	}

	// Verify that deliveries count matches active subscriptions with delivery dates
	if len(deliveries) != activeWithDelivery {
		t.Errorf("GetUpcomingDeliveries() returned %d deliveries, expected %d (active with dates)",
			len(deliveries), activeWithDelivery)
	}
}
