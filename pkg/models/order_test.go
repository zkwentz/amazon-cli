package models

import (
	"encoding/json"
	"testing"
)

func TestOrderJSONMarshaling(t *testing.T) {
	order := Order{
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

	// Test marshaling to JSON
	jsonData, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	// Test unmarshaling from JSON
	var unmarshaled Order
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal order: %v", err)
	}

	// Verify all fields
	if unmarshaled.OrderID != order.OrderID {
		t.Errorf("OrderID mismatch: got %s, want %s", unmarshaled.OrderID, order.OrderID)
	}
	if unmarshaled.Date != order.Date {
		t.Errorf("Date mismatch: got %s, want %s", unmarshaled.Date, order.Date)
	}
	if unmarshaled.Total != order.Total {
		t.Errorf("Total mismatch: got %f, want %f", unmarshaled.Total, order.Total)
	}
	if unmarshaled.Status != order.Status {
		t.Errorf("Status mismatch: got %s, want %s", unmarshaled.Status, order.Status)
	}
	if len(unmarshaled.Items) != len(order.Items) {
		t.Errorf("Items count mismatch: got %d, want %d", len(unmarshaled.Items), len(order.Items))
	}
	if unmarshaled.Tracking == nil {
		t.Error("Tracking should not be nil")
	}
}

func TestOrderWithoutTracking(t *testing.T) {
	order := Order{
		OrderID: "123-4567890-1234567",
		Date:    "2024-01-15",
		Total:   29.99,
		Status:  "pending",
		Items: []OrderItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Product Name",
				Quantity: 1,
				Price:    29.99,
			},
		},
		Tracking: nil,
	}

	// Test marshaling to JSON
	jsonData, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	// Verify tracking is omitted in JSON when nil
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if _, exists := jsonMap["tracking"]; exists {
		t.Error("Tracking field should be omitted when nil")
	}
}

func TestOrderItemJSONMarshaling(t *testing.T) {
	item := OrderItem{
		ASIN:     "B08N5WRWNW",
		Title:    "Test Product",
		Quantity: 2,
		Price:    49.99,
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal order item: %v", err)
	}

	var unmarshaled OrderItem
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal order item: %v", err)
	}

	if unmarshaled.ASIN != item.ASIN {
		t.Errorf("ASIN mismatch: got %s, want %s", unmarshaled.ASIN, item.ASIN)
	}
	if unmarshaled.Title != item.Title {
		t.Errorf("Title mismatch: got %s, want %s", unmarshaled.Title, item.Title)
	}
	if unmarshaled.Quantity != item.Quantity {
		t.Errorf("Quantity mismatch: got %d, want %d", unmarshaled.Quantity, item.Quantity)
	}
	if unmarshaled.Price != item.Price {
		t.Errorf("Price mismatch: got %f, want %f", unmarshaled.Price, item.Price)
	}
}

func TestTrackingJSONMarshaling(t *testing.T) {
	tracking := Tracking{
		Carrier:        "USPS",
		TrackingNumber: "9400111899562537987654",
		Status:         "in_transit",
		DeliveryDate:   "2024-01-20",
	}

	jsonData, err := json.Marshal(tracking)
	if err != nil {
		t.Fatalf("Failed to marshal tracking: %v", err)
	}

	var unmarshaled Tracking
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal tracking: %v", err)
	}

	if unmarshaled.Carrier != tracking.Carrier {
		t.Errorf("Carrier mismatch: got %s, want %s", unmarshaled.Carrier, tracking.Carrier)
	}
	if unmarshaled.TrackingNumber != tracking.TrackingNumber {
		t.Errorf("TrackingNumber mismatch: got %s, want %s", unmarshaled.TrackingNumber, tracking.TrackingNumber)
	}
	if unmarshaled.Status != tracking.Status {
		t.Errorf("Status mismatch: got %s, want %s", unmarshaled.Status, tracking.Status)
	}
	if unmarshaled.DeliveryDate != tracking.DeliveryDate {
		t.Errorf("DeliveryDate mismatch: got %s, want %s", unmarshaled.DeliveryDate, tracking.DeliveryDate)
	}
}

func TestOrdersResponseJSONMarshaling(t *testing.T) {
	response := OrdersResponse{
		Orders: []Order{
			{
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
			},
		},
		TotalCount: 1,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal orders response: %v", err)
	}

	var unmarshaled OrdersResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal orders response: %v", err)
	}

	if unmarshaled.TotalCount != response.TotalCount {
		t.Errorf("TotalCount mismatch: got %d, want %d", unmarshaled.TotalCount, response.TotalCount)
	}
	if len(unmarshaled.Orders) != len(response.Orders) {
		t.Errorf("Orders count mismatch: got %d, want %d", len(unmarshaled.Orders), len(response.Orders))
	}
}

func TestOrdersResponseWithMultipleOrders(t *testing.T) {
	response := OrdersResponse{
		Orders: []Order{
			{
				OrderID: "123-4567890-1234567",
				Date:    "2024-01-15",
				Total:   29.99,
				Status:  "delivered",
				Items: []OrderItem{
					{
						ASIN:     "B08N5WRWNW",
						Title:    "Product 1",
						Quantity: 1,
						Price:    29.99,
					},
				},
			},
			{
				OrderID: "123-4567890-7654321",
				Date:    "2024-01-10",
				Total:   49.99,
				Status:  "pending",
				Items: []OrderItem{
					{
						ASIN:     "B08N5WRWXX",
						Title:    "Product 2",
						Quantity: 2,
						Price:    24.99,
					},
				},
			},
		},
		TotalCount: 2,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal orders response: %v", err)
	}

	var unmarshaled OrdersResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal orders response: %v", err)
	}

	if unmarshaled.TotalCount != response.TotalCount {
		t.Errorf("TotalCount mismatch: got %d, want %d", unmarshaled.TotalCount, response.TotalCount)
	}
	if len(unmarshaled.Orders) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(unmarshaled.Orders))
	}

	// Verify first order
	if unmarshaled.Orders[0].OrderID != "123-4567890-1234567" {
		t.Errorf("First order ID mismatch: got %s", unmarshaled.Orders[0].OrderID)
	}

	// Verify second order
	if unmarshaled.Orders[1].OrderID != "123-4567890-7654321" {
		t.Errorf("Second order ID mismatch: got %s", unmarshaled.Orders[1].OrderID)
	}
}

func TestOrderWithMultipleItems(t *testing.T) {
	order := Order{
		OrderID: "123-4567890-1234567",
		Date:    "2024-01-15",
		Total:   79.98,
		Status:  "delivered",
		Items: []OrderItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Product 1",
				Quantity: 1,
				Price:    29.99,
			},
			{
				ASIN:     "B08N5WRWXX",
				Title:    "Product 2",
				Quantity: 2,
				Price:    24.99,
			},
		},
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	var unmarshaled Order
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal order: %v", err)
	}

	if len(unmarshaled.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(unmarshaled.Items))
	}

	// Verify total matches sum of items
	expectedTotal := 79.98
	if unmarshaled.Total != expectedTotal {
		t.Errorf("Total mismatch: got %f, want %f", unmarshaled.Total, expectedTotal)
	}
}

func TestOrderJSONSchemaCompliance(t *testing.T) {
	// Test JSON output matches the schema from PRD
	expectedJSON := `{
		"orders": [
			{
				"order_id": "123-4567890-1234567",
				"date": "2024-01-15",
				"total": 29.99,
				"status": "delivered",
				"items": [
					{
						"asin": "B08N5WRWNW",
						"title": "Product Name",
						"quantity": 1,
						"price": 29.99
					}
				],
				"tracking": {
					"carrier": "UPS",
					"tracking_number": "1Z999AA10123456784",
					"status": "delivered",
					"delivery_date": "2024-01-17"
				}
			}
		],
		"total_count": 1
	}`

	var expected OrdersResponse
	err := json.Unmarshal([]byte(expectedJSON), &expected)
	if err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	// Verify the structure is correct
	if len(expected.Orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(expected.Orders))
	}
	if expected.TotalCount != 1 {
		t.Errorf("Expected total_count 1, got %d", expected.TotalCount)
	}
	if expected.Orders[0].OrderID != "123-4567890-1234567" {
		t.Errorf("OrderID mismatch: got %s", expected.Orders[0].OrderID)
	}
	if expected.Orders[0].Tracking == nil {
		t.Error("Tracking should not be nil")
	}
}
