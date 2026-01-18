package models

import (
	"encoding/json"
	"testing"
)

func TestSubscriptionJSONMarshaling(t *testing.T) {
	subscription := Subscription{
		SubscriptionID:  "S01-1234567-8901234",
		ASIN:            "B00EXAMPLE",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "2024-02-01",
		Status:          "active",
		Quantity:        1,
	}

	// Test marshaling
	data, err := json.Marshal(subscription)
	if err != nil {
		t.Fatalf("Failed to marshal subscription: %v", err)
	}

	// Test unmarshaling
	var decoded Subscription
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal subscription: %v", err)
	}

	// Verify all fields
	if decoded.SubscriptionID != subscription.SubscriptionID {
		t.Errorf("SubscriptionID mismatch: got %s, want %s", decoded.SubscriptionID, subscription.SubscriptionID)
	}
	if decoded.ASIN != subscription.ASIN {
		t.Errorf("ASIN mismatch: got %s, want %s", decoded.ASIN, subscription.ASIN)
	}
	if decoded.Title != subscription.Title {
		t.Errorf("Title mismatch: got %s, want %s", decoded.Title, subscription.Title)
	}
	if decoded.Price != subscription.Price {
		t.Errorf("Price mismatch: got %f, want %f", decoded.Price, subscription.Price)
	}
	if decoded.DiscountPercent != subscription.DiscountPercent {
		t.Errorf("DiscountPercent mismatch: got %d, want %d", decoded.DiscountPercent, subscription.DiscountPercent)
	}
	if decoded.FrequencyWeeks != subscription.FrequencyWeeks {
		t.Errorf("FrequencyWeeks mismatch: got %d, want %d", decoded.FrequencyWeeks, subscription.FrequencyWeeks)
	}
	if decoded.NextDelivery != subscription.NextDelivery {
		t.Errorf("NextDelivery mismatch: got %s, want %s", decoded.NextDelivery, subscription.NextDelivery)
	}
	if decoded.Status != subscription.Status {
		t.Errorf("Status mismatch: got %s, want %s", decoded.Status, subscription.Status)
	}
	if decoded.Quantity != subscription.Quantity {
		t.Errorf("Quantity mismatch: got %d, want %d", decoded.Quantity, subscription.Quantity)
	}
}

func TestSubscriptionsResponseJSONMarshaling(t *testing.T) {
	response := SubscriptionsResponse{
		Subscriptions: []Subscription{
			{
				SubscriptionID:  "S01-1234567-8901234",
				ASIN:            "B00EXAMPLE",
				Title:           "Coffee Pods 100 Count",
				Price:           45.99,
				DiscountPercent: 15,
				FrequencyWeeks:  4,
				NextDelivery:    "2024-02-01",
				Status:          "active",
				Quantity:        1,
			},
			{
				SubscriptionID:  "S01-9876543-2109876",
				ASIN:            "B00EXAMPLE2",
				Title:           "Paper Towels 12 Pack",
				Price:           29.99,
				DiscountPercent: 10,
				FrequencyWeeks:  8,
				NextDelivery:    "2024-03-15",
				Status:          "active",
				Quantity:        2,
			},
		},
	}

	// Test marshaling
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal subscriptions response: %v", err)
	}

	// Test unmarshaling
	var decoded SubscriptionsResponse
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal subscriptions response: %v", err)
	}

	// Verify subscriptions count
	if len(decoded.Subscriptions) != len(response.Subscriptions) {
		t.Errorf("Subscriptions count mismatch: got %d, want %d", len(decoded.Subscriptions), len(response.Subscriptions))
	}

	// Verify first subscription
	if len(decoded.Subscriptions) > 0 {
		if decoded.Subscriptions[0].SubscriptionID != response.Subscriptions[0].SubscriptionID {
			t.Errorf("First subscription ID mismatch: got %s, want %s",
				decoded.Subscriptions[0].SubscriptionID, response.Subscriptions[0].SubscriptionID)
		}
	}
}

func TestUpcomingDeliveryJSONMarshaling(t *testing.T) {
	delivery := UpcomingDelivery{
		SubscriptionID: "S01-1234567-8901234",
		ASIN:           "B00EXAMPLE",
		Title:          "Coffee Pods 100 Count",
		DeliveryDate:   "2024-02-01",
		Quantity:       1,
	}

	// Test marshaling
	data, err := json.Marshal(delivery)
	if err != nil {
		t.Fatalf("Failed to marshal upcoming delivery: %v", err)
	}

	// Test unmarshaling
	var decoded UpcomingDelivery
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal upcoming delivery: %v", err)
	}

	// Verify all fields
	if decoded.SubscriptionID != delivery.SubscriptionID {
		t.Errorf("SubscriptionID mismatch: got %s, want %s", decoded.SubscriptionID, delivery.SubscriptionID)
	}
	if decoded.ASIN != delivery.ASIN {
		t.Errorf("ASIN mismatch: got %s, want %s", decoded.ASIN, delivery.ASIN)
	}
	if decoded.Title != delivery.Title {
		t.Errorf("Title mismatch: got %s, want %s", decoded.Title, delivery.Title)
	}
	if decoded.DeliveryDate != delivery.DeliveryDate {
		t.Errorf("DeliveryDate mismatch: got %s, want %s", decoded.DeliveryDate, delivery.DeliveryDate)
	}
	if decoded.Quantity != delivery.Quantity {
		t.Errorf("Quantity mismatch: got %d, want %d", decoded.Quantity, delivery.Quantity)
	}
}

func TestEmptySubscriptionsResponse(t *testing.T) {
	response := SubscriptionsResponse{
		Subscriptions: []Subscription{},
	}

	// Test marshaling
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal empty subscriptions response: %v", err)
	}

	// Test unmarshaling
	var decoded SubscriptionsResponse
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty subscriptions response: %v", err)
	}

	// Verify empty subscriptions
	if decoded.Subscriptions == nil {
		t.Error("Subscriptions should not be nil")
	}
	if len(decoded.Subscriptions) != 0 {
		t.Errorf("Subscriptions count mismatch: got %d, want 0", len(decoded.Subscriptions))
	}
}

func TestSubscriptionJSONTags(t *testing.T) {
	subscription := Subscription{
		SubscriptionID:  "S01-1234567-8901234",
		ASIN:            "B00EXAMPLE",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "2024-02-01",
		Status:          "active",
		Quantity:        1,
	}

	data, err := json.Marshal(subscription)
	if err != nil {
		t.Fatalf("Failed to marshal subscription: %v", err)
	}

	// Check that JSON uses expected field names
	jsonStr := string(data)
	expectedFields := []string{
		"subscription_id",
		"asin",
		"title",
		"price",
		"discount_percent",
		"frequency_weeks",
		"next_delivery",
		"status",
		"quantity",
	}

	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON output missing expected field: %s", field)
		}
	}
}

func TestUpcomingDeliveryJSONTags(t *testing.T) {
	delivery := UpcomingDelivery{
		SubscriptionID: "S01-1234567-8901234",
		ASIN:           "B00EXAMPLE",
		Title:          "Coffee Pods 100 Count",
		DeliveryDate:   "2024-02-01",
		Quantity:       1,
	}

	data, err := json.Marshal(delivery)
	if err != nil {
		t.Fatalf("Failed to marshal upcoming delivery: %v", err)
	}

	// Check that JSON uses expected field names
	jsonStr := string(data)
	expectedFields := []string{
		"subscription_id",
		"asin",
		"title",
		"delivery_date",
		"quantity",
	}

	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON output missing expected field: %s", field)
		}
	}
}

func TestSubscriptionWithDifferentStatuses(t *testing.T) {
	statuses := []string{"active", "paused", "cancelled"}

	for _, status := range statuses {
		subscription := Subscription{
			SubscriptionID:  "S01-1234567-8901234",
			ASIN:            "B00EXAMPLE",
			Title:           "Test Product",
			Price:           29.99,
			DiscountPercent: 10,
			FrequencyWeeks:  4,
			NextDelivery:    "2024-02-01",
			Status:          status,
			Quantity:        1,
		}

		data, err := json.Marshal(subscription)
		if err != nil {
			t.Fatalf("Failed to marshal subscription with status %s: %v", status, err)
		}

		var decoded Subscription
		err = json.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("Failed to unmarshal subscription with status %s: %v", status, err)
		}

		if decoded.Status != status {
			t.Errorf("Status mismatch: got %s, want %s", decoded.Status, status)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
