package models

import (
	"encoding/json"
	"testing"
)

func TestOrderJSONSerialization(t *testing.T) {
	order := Order{
		OrderID: "123-4567890-1234567",
		Date:    "2024-01-15",
		Total:   29.99,
		Status:  "delivered",
		Items: []OrderItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Sample Product",
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

	// Test marshaling to JSON
	jsonData, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	// Test unmarshaling from JSON
	var unmarshaledOrder Order
	err = json.Unmarshal(jsonData, &unmarshaledOrder)
	if err != nil {
		t.Fatalf("Failed to unmarshal order: %v", err)
	}

	// Verify fields
	if unmarshaledOrder.OrderID != order.OrderID {
		t.Errorf("OrderID mismatch: got %s, want %s", unmarshaledOrder.OrderID, order.OrderID)
	}

	if unmarshaledOrder.Date != order.Date {
		t.Errorf("Date mismatch: got %s, want %s", unmarshaledOrder.Date, order.Date)
	}

	if unmarshaledOrder.Total != order.Total {
		t.Errorf("Total mismatch: got %f, want %f", unmarshaledOrder.Total, order.Total)
	}

	if unmarshaledOrder.Status != order.Status {
		t.Errorf("Status mismatch: got %s, want %s", unmarshaledOrder.Status, order.Status)
	}

	if len(unmarshaledOrder.Items) != len(order.Items) {
		t.Errorf("Items count mismatch: got %d, want %d", len(unmarshaledOrder.Items), len(order.Items))
	}

	if unmarshaledOrder.Tracking == nil {
		t.Error("Tracking is nil after unmarshal")
	}
}

func TestOrdersResponseJSONSerialization(t *testing.T) {
	response := OrdersResponse{
		Orders: []Order{
			{
				OrderID: "123-4567890-1234567",
				Date:    "2024-01-15",
				Total:   29.99,
				Status:  "delivered",
				Items:   []OrderItem{},
			},
		},
		TotalCount: 1,
	}

	// Test marshaling to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal OrdersResponse: %v", err)
	}

	// Test unmarshaling from JSON
	var unmarshaledResponse OrdersResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal OrdersResponse: %v", err)
	}

	if unmarshaledResponse.TotalCount != response.TotalCount {
		t.Errorf("TotalCount mismatch: got %d, want %d", unmarshaledResponse.TotalCount, response.TotalCount)
	}

	if len(unmarshaledResponse.Orders) != len(response.Orders) {
		t.Errorf("Orders count mismatch: got %d, want %d", len(unmarshaledResponse.Orders), len(response.Orders))
	}
}

func TestOrderWithoutTracking(t *testing.T) {
	order := Order{
		OrderID: "123-4567890-1234567",
		Date:    "2024-01-15",
		Total:   29.99,
		Status:  "pending",
		Items:   []OrderItem{},
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order without tracking: %v", err)
	}

	// Verify tracking field is omitted when nil
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if _, exists := jsonMap["tracking"]; exists {
		t.Error("Tracking field should be omitted when nil")
	}
}
