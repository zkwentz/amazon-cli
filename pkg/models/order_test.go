package models

import (
	"encoding/json"
	"testing"
)

func TestTrackingJSONSerialization(t *testing.T) {
	tracking := &Tracking{
		Carrier:        "UPS",
		TrackingNumber: "1Z999AA10123456784",
		Status:         "delivered",
		DeliveryDate:   "2024-01-17",
	}

	// Marshal to JSON
	data, err := json.Marshal(tracking)
	if err != nil {
		t.Fatalf("failed to marshal tracking: %v", err)
	}

	// Verify JSON structure
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Check required fields exist
	requiredFields := []string{"carrier", "tracking_number", "status", "delivery_date"}
	for _, field := range requiredFields {
		if _, ok := decoded[field]; !ok {
			t.Errorf("missing required field: %s", field)
		}
	}

	// Unmarshal back to struct
	var newTracking Tracking
	if err := json.Unmarshal(data, &newTracking); err != nil {
		t.Fatalf("failed to unmarshal to struct: %v", err)
	}

	// Verify values match
	if newTracking.Carrier != tracking.Carrier {
		t.Errorf("carrier mismatch: got %s, want %s", newTracking.Carrier, tracking.Carrier)
	}
	if newTracking.TrackingNumber != tracking.TrackingNumber {
		t.Errorf("tracking number mismatch: got %s, want %s", newTracking.TrackingNumber, tracking.TrackingNumber)
	}
	if newTracking.Status != tracking.Status {
		t.Errorf("status mismatch: got %s, want %s", newTracking.Status, tracking.Status)
	}
	if newTracking.DeliveryDate != tracking.DeliveryDate {
		t.Errorf("delivery date mismatch: got %s, want %s", newTracking.DeliveryDate, tracking.DeliveryDate)
	}
}

func TestOrderWithTracking(t *testing.T) {
	order := &Order{
		OrderID: "123-4567890-1234567",
		Date:    "2024-01-15",
		Total:   29.99,
		Status:  "delivered",
		Items: []OrderItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Product Name",
				Quantity: 1,
				Price:    29.99,
			},
		},
		Tracking: &Tracking{
			Carrier:        "UPS",
			TrackingNumber: "1Z999AA10123456784",
			Status:         "delivered",
			DeliveryDate:   "2024-01-17",
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal order: %v", err)
	}

	// Unmarshal back
	var decoded Order
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal order: %v", err)
	}

	// Verify tracking is preserved
	if decoded.Tracking == nil {
		t.Fatal("tracking is nil after unmarshal")
	}

	if decoded.Tracking.Carrier != order.Tracking.Carrier {
		t.Errorf("carrier mismatch after unmarshal")
	}
}
