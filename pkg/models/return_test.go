package models

import (
	"encoding/json"
	"testing"
)

func TestReturnableItemJSON(t *testing.T) {
	item := ReturnableItem{
		OrderID:      "123-4567890-1234567",
		ItemID:       "ITEM001",
		ASIN:         "B08N5WRWNW",
		Title:        "Test Product",
		Price:        29.99,
		PurchaseDate: "2024-01-15",
		ReturnWindow: "2024-02-14",
	}

	jsonBytes, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded ReturnableItem
	err = json.Unmarshal(jsonBytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.OrderID != item.OrderID {
		t.Errorf("OrderID mismatch: expected %s, got %s", item.OrderID, decoded.OrderID)
	}
	if decoded.ASIN != item.ASIN {
		t.Errorf("ASIN mismatch: expected %s, got %s", item.ASIN, decoded.ASIN)
	}
	if decoded.Price != item.Price {
		t.Errorf("Price mismatch: expected %f, got %f", item.Price, decoded.Price)
	}
}

func TestReturnableItemsResponseJSON(t *testing.T) {
	response := ReturnableItemsResponse{
		Items: []ReturnableItem{
			{
				OrderID:      "123-4567890-1234567",
				ItemID:       "ITEM001",
				ASIN:         "B08N5WRWNW",
				Title:        "Product 1",
				Price:        29.99,
				PurchaseDate: "2024-01-15",
				ReturnWindow: "2024-02-14",
			},
			{
				OrderID:      "123-4567890-1234568",
				ItemID:       "ITEM002",
				ASIN:         "B07XJ8C8F5",
				Title:        "Product 2",
				Price:        12.99,
				PurchaseDate: "2024-01-20",
				ReturnWindow: "2024-02-19",
			},
		},
		TotalCount: 2,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded ReturnableItemsResponse
	err = json.Unmarshal(jsonBytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.TotalCount != response.TotalCount {
		t.Errorf("TotalCount mismatch: expected %d, got %d", response.TotalCount, decoded.TotalCount)
	}

	if len(decoded.Items) != len(response.Items) {
		t.Errorf("Items length mismatch: expected %d, got %d", len(response.Items), len(decoded.Items))
	}
}

func TestReturnOptionJSON(t *testing.T) {
	option := ReturnOption{
		Method:          "UPS",
		Label:           "UPS Drop Off",
		DropoffLocation: "123 Main St",
		Fee:             0.0,
	}

	jsonBytes, err := json.Marshal(option)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded ReturnOption
	err = json.Unmarshal(jsonBytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.Method != option.Method {
		t.Errorf("Method mismatch: expected %s, got %s", option.Method, decoded.Method)
	}
}

func TestReturnJSON(t *testing.T) {
	returnItem := Return{
		ReturnID:  "R123456",
		OrderID:   "123-4567890-1234567",
		ItemID:    "ITEM001",
		Status:    "initiated",
		Reason:    "defective",
		CreatedAt: "2024-01-25T10:00:00Z",
	}

	jsonBytes, err := json.Marshal(returnItem)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var decoded Return
	err = json.Unmarshal(jsonBytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.ReturnID != returnItem.ReturnID {
		t.Errorf("ReturnID mismatch: expected %s, got %s", returnItem.ReturnID, decoded.ReturnID)
	}
	if decoded.Reason != returnItem.Reason {
		t.Errorf("Reason mismatch: expected %s, got %s", returnItem.Reason, decoded.Reason)
	}
}
