package models

import (
	"encoding/json"
	"testing"
)

func TestReturnableItem_JSON(t *testing.T) {
	item := ReturnableItem{
		OrderID:      "111-2222222-3333333",
		ItemID:       "item-123",
		ASIN:         "B08N5WRWNW",
		Title:        "Test Product",
		Price:        29.99,
		PurchaseDate: "2026-01-15",
		ReturnWindow: "2026-02-15",
	}

	// Test marshaling
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal ReturnableItem: %v", err)
	}

	// Test unmarshaling
	var decoded ReturnableItem
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ReturnableItem: %v", err)
	}

	// Verify fields
	if decoded.OrderID != item.OrderID {
		t.Errorf("OrderID mismatch: got %s, want %s", decoded.OrderID, item.OrderID)
	}
	if decoded.ItemID != item.ItemID {
		t.Errorf("ItemID mismatch: got %s, want %s", decoded.ItemID, item.ItemID)
	}
	if decoded.ASIN != item.ASIN {
		t.Errorf("ASIN mismatch: got %s, want %s", decoded.ASIN, item.ASIN)
	}
	if decoded.Title != item.Title {
		t.Errorf("Title mismatch: got %s, want %s", decoded.Title, item.Title)
	}
	if decoded.Price != item.Price {
		t.Errorf("Price mismatch: got %f, want %f", decoded.Price, item.Price)
	}
	if decoded.PurchaseDate != item.PurchaseDate {
		t.Errorf("PurchaseDate mismatch: got %s, want %s", decoded.PurchaseDate, item.PurchaseDate)
	}
	if decoded.ReturnWindow != item.ReturnWindow {
		t.Errorf("ReturnWindow mismatch: got %s, want %s", decoded.ReturnWindow, item.ReturnWindow)
	}
}

func TestReturnOption_JSON(t *testing.T) {
	option := ReturnOption{
		Method:          "UPS Drop-off",
		Label:           "Free UPS Return",
		DropoffLocation: "UPS Store - 123 Main St",
		Fee:             0.0,
	}

	// Test marshaling
	data, err := json.Marshal(option)
	if err != nil {
		t.Fatalf("Failed to marshal ReturnOption: %v", err)
	}

	// Test unmarshaling
	var decoded ReturnOption
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ReturnOption: %v", err)
	}

	// Verify fields
	if decoded.Method != option.Method {
		t.Errorf("Method mismatch: got %s, want %s", decoded.Method, option.Method)
	}
	if decoded.Label != option.Label {
		t.Errorf("Label mismatch: got %s, want %s", decoded.Label, option.Label)
	}
	if decoded.DropoffLocation != option.DropoffLocation {
		t.Errorf("DropoffLocation mismatch: got %s, want %s", decoded.DropoffLocation, option.DropoffLocation)
	}
	if decoded.Fee != option.Fee {
		t.Errorf("Fee mismatch: got %f, want %f", decoded.Fee, option.Fee)
	}
}

func TestReturn_JSON(t *testing.T) {
	returnItem := Return{
		ReturnID:  "return-123",
		OrderID:   "111-2222222-3333333",
		ItemID:    "item-123",
		Status:    "Pending",
		Reason:    "Defective item",
		CreatedAt: "2026-01-18T10:00:00Z",
	}

	// Test marshaling
	data, err := json.Marshal(returnItem)
	if err != nil {
		t.Fatalf("Failed to marshal Return: %v", err)
	}

	// Test unmarshaling
	var decoded Return
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Return: %v", err)
	}

	// Verify fields
	if decoded.ReturnID != returnItem.ReturnID {
		t.Errorf("ReturnID mismatch: got %s, want %s", decoded.ReturnID, returnItem.ReturnID)
	}
	if decoded.OrderID != returnItem.OrderID {
		t.Errorf("OrderID mismatch: got %s, want %s", decoded.OrderID, returnItem.OrderID)
	}
	if decoded.ItemID != returnItem.ItemID {
		t.Errorf("ItemID mismatch: got %s, want %s", decoded.ItemID, returnItem.ItemID)
	}
	if decoded.Status != returnItem.Status {
		t.Errorf("Status mismatch: got %s, want %s", decoded.Status, returnItem.Status)
	}
	if decoded.Reason != returnItem.Reason {
		t.Errorf("Reason mismatch: got %s, want %s", decoded.Reason, returnItem.Reason)
	}
	if decoded.CreatedAt != returnItem.CreatedAt {
		t.Errorf("CreatedAt mismatch: got %s, want %s", decoded.CreatedAt, returnItem.CreatedAt)
	}
}

func TestReturnLabel_JSON(t *testing.T) {
	label := ReturnLabel{
		URL:          "https://example.com/label.pdf",
		Carrier:      "UPS",
		Instructions: "Drop off at any UPS location",
	}

	// Test marshaling
	data, err := json.Marshal(label)
	if err != nil {
		t.Fatalf("Failed to marshal ReturnLabel: %v", err)
	}

	// Test unmarshaling
	var decoded ReturnLabel
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ReturnLabel: %v", err)
	}

	// Verify fields
	if decoded.URL != label.URL {
		t.Errorf("URL mismatch: got %s, want %s", decoded.URL, label.URL)
	}
	if decoded.Carrier != label.Carrier {
		t.Errorf("Carrier mismatch: got %s, want %s", decoded.Carrier, label.Carrier)
	}
	if decoded.Instructions != label.Instructions {
		t.Errorf("Instructions mismatch: got %s, want %s", decoded.Instructions, label.Instructions)
	}
}
