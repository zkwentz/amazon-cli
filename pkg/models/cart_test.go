package models

import (
	"encoding/json"
	"testing"
)

func TestCartItemJSON(t *testing.T) {
	item := CartItem{
		ASIN:     "B08N5WRWNW",
		Title:    "Test Product",
		Price:    29.99,
		Quantity: 2,
		Subtotal: 59.98,
		Prime:    true,
		InStock:  true,
	}

	// Test marshaling
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal CartItem: %v", err)
	}

	// Test unmarshaling
	var decoded CartItem
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal CartItem: %v", err)
	}

	// Verify fields
	if decoded.ASIN != item.ASIN {
		t.Errorf("ASIN mismatch: got %s, want %s", decoded.ASIN, item.ASIN)
	}
	if decoded.Title != item.Title {
		t.Errorf("Title mismatch: got %s, want %s", decoded.Title, item.Title)
	}
	if decoded.Price != item.Price {
		t.Errorf("Price mismatch: got %.2f, want %.2f", decoded.Price, item.Price)
	}
	if decoded.Quantity != item.Quantity {
		t.Errorf("Quantity mismatch: got %d, want %d", decoded.Quantity, item.Quantity)
	}
	if decoded.Subtotal != item.Subtotal {
		t.Errorf("Subtotal mismatch: got %.2f, want %.2f", decoded.Subtotal, item.Subtotal)
	}
	if decoded.Prime != item.Prime {
		t.Errorf("Prime mismatch: got %v, want %v", decoded.Prime, item.Prime)
	}
	if decoded.InStock != item.InStock {
		t.Errorf("InStock mismatch: got %v, want %v", decoded.InStock, item.InStock)
	}
}

func TestCartJSON(t *testing.T) {
	cart := Cart{
		Items: []CartItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Product 1",
				Price:    29.99,
				Quantity: 1,
				Subtotal: 29.99,
				Prime:    true,
				InStock:  true,
			},
			{
				ASIN:     "B08N5WRWXX",
				Title:    "Product 2",
				Price:    49.99,
				Quantity: 2,
				Subtotal: 99.98,
				Prime:    false,
				InStock:  true,
			},
		},
		Subtotal:     129.97,
		EstimatedTax: 10.40,
		Total:        140.37,
		ItemCount:    2,
	}

	// Test marshaling
	data, err := json.Marshal(cart)
	if err != nil {
		t.Fatalf("Failed to marshal Cart: %v", err)
	}

	// Test unmarshaling
	var decoded Cart
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Cart: %v", err)
	}

	// Verify fields
	if len(decoded.Items) != len(cart.Items) {
		t.Errorf("Items count mismatch: got %d, want %d", len(decoded.Items), len(cart.Items))
	}
	if decoded.Subtotal != cart.Subtotal {
		t.Errorf("Subtotal mismatch: got %.2f, want %.2f", decoded.Subtotal, cart.Subtotal)
	}
	if decoded.EstimatedTax != cart.EstimatedTax {
		t.Errorf("EstimatedTax mismatch: got %.2f, want %.2f", decoded.EstimatedTax, cart.EstimatedTax)
	}
	if decoded.Total != cart.Total {
		t.Errorf("Total mismatch: got %.2f, want %.2f", decoded.Total, cart.Total)
	}
	if decoded.ItemCount != cart.ItemCount {
		t.Errorf("ItemCount mismatch: got %d, want %d", decoded.ItemCount, cart.ItemCount)
	}
}

func TestAddressJSON(t *testing.T) {
	address := Address{
		ID:      "addr_123",
		Name:    "John Doe",
		Street:  "123 Main St",
		City:    "Springfield",
		State:   "IL",
		Zip:     "62701",
		Country: "US",
		Default: true,
	}

	// Test marshaling
	data, err := json.Marshal(address)
	if err != nil {
		t.Fatalf("Failed to marshal Address: %v", err)
	}

	// Test unmarshaling
	var decoded Address
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Address: %v", err)
	}

	// Verify fields
	if decoded.ID != address.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, address.ID)
	}
	if decoded.Name != address.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, address.Name)
	}
	if decoded.Street != address.Street {
		t.Errorf("Street mismatch: got %s, want %s", decoded.Street, address.Street)
	}
	if decoded.City != address.City {
		t.Errorf("City mismatch: got %s, want %s", decoded.City, address.City)
	}
	if decoded.State != address.State {
		t.Errorf("State mismatch: got %s, want %s", decoded.State, address.State)
	}
	if decoded.Zip != address.Zip {
		t.Errorf("Zip mismatch: got %s, want %s", decoded.Zip, address.Zip)
	}
	if decoded.Country != address.Country {
		t.Errorf("Country mismatch: got %s, want %s", decoded.Country, address.Country)
	}
	if decoded.Default != address.Default {
		t.Errorf("Default mismatch: got %v, want %v", decoded.Default, address.Default)
	}
}

func TestPaymentMethodJSON(t *testing.T) {
	payment := PaymentMethod{
		ID:      "pay_123",
		Type:    "Visa",
		Last4:   "4242",
		Default: true,
	}

	// Test marshaling
	data, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal PaymentMethod: %v", err)
	}

	// Test unmarshaling
	var decoded PaymentMethod
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal PaymentMethod: %v", err)
	}

	// Verify fields
	if decoded.ID != payment.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, payment.ID)
	}
	if decoded.Type != payment.Type {
		t.Errorf("Type mismatch: got %s, want %s", decoded.Type, payment.Type)
	}
	if decoded.Last4 != payment.Last4 {
		t.Errorf("Last4 mismatch: got %s, want %s", decoded.Last4, payment.Last4)
	}
	if decoded.Default != payment.Default {
		t.Errorf("Default mismatch: got %v, want %v", decoded.Default, payment.Default)
	}
}

func TestCheckoutPreviewJSON(t *testing.T) {
	preview := CheckoutPreview{
		Cart: &Cart{
			Items: []CartItem{
				{
					ASIN:     "B08N5WRWNW",
					Title:    "Product 1",
					Price:    29.99,
					Quantity: 1,
					Subtotal: 29.99,
					Prime:    true,
					InStock:  true,
				},
			},
			Subtotal:     29.99,
			EstimatedTax: 2.40,
			Total:        32.39,
			ItemCount:    1,
		},
		Address: &Address{
			ID:      "addr_123",
			Name:    "John Doe",
			Street:  "123 Main St",
			City:    "Springfield",
			State:   "IL",
			Zip:     "62701",
			Country: "US",
			Default: true,
		},
		PaymentMethod: &PaymentMethod{
			ID:      "pay_123",
			Type:    "Visa",
			Last4:   "4242",
			Default: true,
		},
		DeliveryOptions: []string{"Standard", "Express", "Same-Day"},
	}

	// Test marshaling
	data, err := json.Marshal(preview)
	if err != nil {
		t.Fatalf("Failed to marshal CheckoutPreview: %v", err)
	}

	// Test unmarshaling
	var decoded CheckoutPreview
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal CheckoutPreview: %v", err)
	}

	// Verify nested structures exist
	if decoded.Cart == nil {
		t.Error("Cart is nil")
	}
	if decoded.Address == nil {
		t.Error("Address is nil")
	}
	if decoded.PaymentMethod == nil {
		t.Error("PaymentMethod is nil")
	}
	if len(decoded.DeliveryOptions) != len(preview.DeliveryOptions) {
		t.Errorf("DeliveryOptions count mismatch: got %d, want %d", len(decoded.DeliveryOptions), len(preview.DeliveryOptions))
	}
}

func TestOrderConfirmationJSON(t *testing.T) {
	confirmation := OrderConfirmation{
		OrderID:           "123-4567890-1234567",
		Total:             140.37,
		EstimatedDelivery: "2024-01-20",
	}

	// Test marshaling
	data, err := json.Marshal(confirmation)
	if err != nil {
		t.Fatalf("Failed to marshal OrderConfirmation: %v", err)
	}

	// Test unmarshaling
	var decoded OrderConfirmation
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal OrderConfirmation: %v", err)
	}

	// Verify fields
	if decoded.OrderID != confirmation.OrderID {
		t.Errorf("OrderID mismatch: got %s, want %s", decoded.OrderID, confirmation.OrderID)
	}
	if decoded.Total != confirmation.Total {
		t.Errorf("Total mismatch: got %.2f, want %.2f", decoded.Total, confirmation.Total)
	}
	if decoded.EstimatedDelivery != confirmation.EstimatedDelivery {
		t.Errorf("EstimatedDelivery mismatch: got %s, want %s", decoded.EstimatedDelivery, confirmation.EstimatedDelivery)
	}
}

func TestEmptyCart(t *testing.T) {
	cart := Cart{
		Items:        []CartItem{},
		Subtotal:     0.0,
		EstimatedTax: 0.0,
		Total:        0.0,
		ItemCount:    0,
	}

	// Test marshaling empty cart
	data, err := json.Marshal(cart)
	if err != nil {
		t.Fatalf("Failed to marshal empty Cart: %v", err)
	}

	// Test unmarshaling
	var decoded Cart
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty Cart: %v", err)
	}

	if decoded.Items == nil {
		t.Error("Items should be empty slice, not nil")
	}
	if len(decoded.Items) != 0 {
		t.Errorf("Empty cart should have 0 items, got %d", len(decoded.Items))
	}
}
