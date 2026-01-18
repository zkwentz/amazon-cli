package amazon

import (
	"testing"
)

func TestCompleteCheckout(t *testing.T) {
	tests := []struct {
		name        string
		addressID   string
		paymentID   string
		setupCart   func(*Client) error
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty addressID should fail",
			addressID:   "",
			paymentID:   "pay123",
			setupCart:   nil,
			wantErr:     true,
			errContains: "addressID cannot be empty",
		},
		{
			name:        "empty paymentID should fail",
			addressID:   "addr123",
			paymentID:   "",
			setupCart:   nil,
			wantErr:     true,
			errContains: "paymentID cannot be empty",
		},
		{
			name:      "empty cart should fail",
			addressID: "addr123",
			paymentID: "pay123",
			setupCart: func(c *Client) error {
				// Cart is already empty by default
				return nil
			},
			wantErr:     true,
			errContains: "cart is empty",
		},
		{
			name:      "valid checkout should succeed",
			addressID: "addr123",
			paymentID: "pay123",
			setupCart: func(c *Client) error {
				// Add item to cart
				_, err := c.AddToCart("B08N5WRWNW", 1)
				return err
			},
			wantErr: false,
		},
		{
			name:      "multiple items in cart should succeed",
			addressID: "addr456",
			paymentID: "pay456",
			setupCart: func(c *Client) error {
				// Add multiple items to cart
				_, err := c.AddToCart("B08N5WRWNW", 2)
				if err != nil {
					return err
				}
				_, err = c.AddToCart("B07XJ8C8F5", 1)
				return err
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			// Setup cart if needed
			if tt.setupCart != nil {
				if err := tt.setupCart(client); err != nil {
					t.Fatalf("failed to setup cart: %v", err)
				}
			}

			// Execute CompleteCheckout
			confirmation, err := client.CompleteCheckout(tt.addressID, tt.paymentID)

			// Check error expectations
			if tt.wantErr {
				if err == nil {
					t.Errorf("CompleteCheckout() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("CompleteCheckout() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			// Check success expectations
			if err != nil {
				t.Errorf("CompleteCheckout() unexpected error: %v", err)
				return
			}

			if confirmation == nil {
				t.Error("CompleteCheckout() returned nil confirmation")
				return
			}

			// Verify confirmation fields
			if confirmation.OrderID == "" {
				t.Error("CompleteCheckout() OrderID is empty")
			}

			if confirmation.Total <= 0 {
				t.Error("CompleteCheckout() Total should be greater than 0")
			}

			if confirmation.EstimatedDelivery == "" {
				t.Error("CompleteCheckout() EstimatedDelivery is empty")
			}
		})
	}
}

func TestCompleteCheckout_OrderConfirmation(t *testing.T) {
	client := NewClient()

	// Setup cart
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	// Complete checkout
	confirmation, err := client.CompleteCheckout("addr123", "pay123")
	if err != nil {
		t.Fatalf("CompleteCheckout() error = %v", err)
	}

	// Verify order confirmation structure
	if confirmation.OrderID == "" {
		t.Error("OrderID should not be empty")
	}

	// Check total is approximately correct (29.99 + 8% tax = 32.3892)
	expectedTotal := 32.39
	if confirmation.Total < expectedTotal-0.01 || confirmation.Total > expectedTotal+0.01 {
		t.Errorf("Total = %v, want approximately %v", confirmation.Total, expectedTotal)
	}

	if confirmation.EstimatedDelivery == "" {
		t.Error("EstimatedDelivery should not be empty")
	}
}

func TestAddToCart(t *testing.T) {
	tests := []struct {
		name      string
		asin      string
		quantity  int
		wantErr   bool
		errString string
	}{
		{
			name:      "valid add to cart",
			asin:      "B08N5WRWNW",
			quantity:  1,
			wantErr:   false,
			errString: "",
		},
		{
			name:      "empty ASIN should fail",
			asin:      "",
			quantity:  1,
			wantErr:   true,
			errString: "ASIN cannot be empty",
		},
		{
			name:      "zero quantity should fail",
			asin:      "B08N5WRWNW",
			quantity:  0,
			wantErr:   true,
			errString: "quantity must be positive",
		},
		{
			name:      "negative quantity should fail",
			asin:      "B08N5WRWNW",
			quantity:  -1,
			wantErr:   true,
			errString: "quantity must be positive",
		},
		{
			name:      "multiple quantity",
			asin:      "B08N5WRWNW",
			quantity:  5,
			wantErr:   false,
			errString: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			cart, err := client.AddToCart(tt.asin, tt.quantity)

			if tt.wantErr {
				if err == nil {
					t.Error("AddToCart() expected error but got none")
				} else if err.Error() != tt.errString {
					t.Errorf("AddToCart() error = %v, want %v", err.Error(), tt.errString)
				}
				return
			}

			if err != nil {
				t.Errorf("AddToCart() unexpected error: %v", err)
				return
			}

			if cart == nil {
				t.Error("AddToCart() returned nil cart")
				return
			}

			// Verify cart has items
			if len(cart.Items) == 0 {
				t.Error("AddToCart() cart should have items")
			}

			// Verify quantity matches
			if cart.Items[0].Quantity != tt.quantity {
				t.Errorf("Cart item quantity = %v, want %v", cart.Items[0].Quantity, tt.quantity)
			}
		})
	}
}

// TestAddToCart_CalculatesTotals verifies that adding items calculates subtotal, tax, and total correctly
func TestAddToCart_CalculatesTotals(t *testing.T) {
	client := NewClient()

	// Add first item
	cart, err := client.AddToCart("B08N5WRWNW", 2)
	if err != nil {
		t.Fatalf("AddToCart() failed: %v", err)
	}

	expectedSubtotal := 29.99 * 2 // 59.98
	if cart.Subtotal != expectedSubtotal {
		t.Errorf("Subtotal = %v, want %v", cart.Subtotal, expectedSubtotal)
	}

	expectedTax := expectedSubtotal * 0.08 // 4.7984
	if cart.EstimatedTax != expectedTax {
		t.Errorf("EstimatedTax = %v, want %v", cart.EstimatedTax, expectedTax)
	}

	expectedTotal := expectedSubtotal + expectedTax
	if cart.Total != expectedTotal {
		t.Errorf("Total = %v, want %v", cart.Total, expectedTotal)
	}

	if cart.ItemCount != 2 {
		t.Errorf("ItemCount = %v, want 2", cart.ItemCount)
	}
}

// TestAddToCart_MultipleItems verifies that multiple different items can be added
func TestAddToCart_MultipleItems(t *testing.T) {
	client := NewClient()

	// Add first item
	cart1, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("AddToCart() first item failed: %v", err)
	}

	if len(cart1.Items) != 1 {
		t.Errorf("After first add, cart has %v items, want 1", len(cart1.Items))
	}

	// Add second item
	cart2, err := client.AddToCart("B07XJ8C8F5", 2)
	if err != nil {
		t.Fatalf("AddToCart() second item failed: %v", err)
	}

	if len(cart2.Items) != 2 {
		t.Errorf("After second add, cart has %v items, want 2", len(cart2.Items))
	}

	// Verify total item count
	if cart2.ItemCount != 3 {
		t.Errorf("ItemCount = %v, want 3 (1 + 2)", cart2.ItemCount)
	}

	// Verify subtotal is sum of both items
	expectedSubtotal := 29.99*1 + 29.99*2 // 89.97
	if cart2.Subtotal != expectedSubtotal {
		t.Errorf("Subtotal = %v, want %v", cart2.Subtotal, expectedSubtotal)
	}
}

// TestAddToCart_ItemFields verifies all cart item fields are populated correctly
func TestAddToCart_ItemFields(t *testing.T) {
	client := NewClient()

	cart, err := client.AddToCart("B08N5WRWNW", 3)
	if err != nil {
		t.Fatalf("AddToCart() failed: %v", err)
	}

	if len(cart.Items) != 1 {
		t.Fatalf("Cart has %v items, want 1", len(cart.Items))
	}

	item := cart.Items[0]

	// Verify ASIN
	if item.ASIN != "B08N5WRWNW" {
		t.Errorf("ASIN = %v, want B08N5WRWNW", item.ASIN)
	}

	// Verify title is set
	if item.Title == "" {
		t.Error("Title should not be empty")
	}

	// Verify price is positive
	if item.Price <= 0 {
		t.Errorf("Price = %v, should be positive", item.Price)
	}

	// Verify quantity
	if item.Quantity != 3 {
		t.Errorf("Quantity = %v, want 3", item.Quantity)
	}

	// Verify subtotal
	expectedSubtotal := item.Price * float64(item.Quantity)
	if item.Subtotal != expectedSubtotal {
		t.Errorf("Subtotal = %v, want %v", item.Subtotal, expectedSubtotal)
	}

	// Verify Prime flag
	if !item.Prime {
		t.Error("Prime should be true for test items")
	}

	// Verify InStock flag
	if !item.InStock {
		t.Error("InStock should be true for test items")
	}
}

func TestGetCart(t *testing.T) {
	client := NewClient()
	cart, err := client.GetCart()

	if err != nil {
		t.Errorf("GetCart() unexpected error: %v", err)
	}

	if cart == nil {
		t.Error("GetCart() returned nil cart")
	}

	// Empty cart should have zero items
	if cart.ItemCount != 0 {
		t.Errorf("GetCart() ItemCount = %v, want 0", cart.ItemCount)
	}
}

func TestRemoveFromCart(t *testing.T) {
	tests := []struct {
		name      string
		asin      string
		wantErr   bool
		errString string
	}{
		{
			name:      "valid ASIN",
			asin:      "B08N5WRWNW",
			wantErr:   false,
			errString: "",
		},
		{
			name:      "empty ASIN should fail",
			asin:      "",
			wantErr:   true,
			errString: "ASIN cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			cart, err := client.RemoveFromCart(tt.asin)

			if tt.wantErr {
				if err == nil {
					t.Error("RemoveFromCart() expected error but got none")
				} else if err.Error() != tt.errString {
					t.Errorf("RemoveFromCart() error = %v, want %v", err.Error(), tt.errString)
				}
				return
			}

			if err != nil {
				t.Errorf("RemoveFromCart() unexpected error: %v", err)
			}

			if cart == nil {
				t.Error("RemoveFromCart() returned nil cart")
			}
		})
	}
}

// TestRemoveFromCart_RemovesItem verifies that removing an item actually removes it from cart
func TestRemoveFromCart_RemovesItem(t *testing.T) {
	client := NewClient()

	// Add item to cart
	_, err := client.AddToCart("B08N5WRWNW", 2)
	if err != nil {
		t.Fatalf("AddToCart() failed: %v", err)
	}

	// Verify item is in cart
	cart, err := client.GetCart()
	if err != nil {
		t.Fatalf("GetCart() failed: %v", err)
	}

	if len(cart.Items) != 1 {
		t.Fatalf("Cart should have 1 item, got %v", len(cart.Items))
	}

	// Remove the item
	cart, err = client.RemoveFromCart("B08N5WRWNW")
	if err != nil {
		t.Fatalf("RemoveFromCart() failed: %v", err)
	}

	// Verify item is removed
	if len(cart.Items) != 0 {
		t.Errorf("Cart should have 0 items after removal, got %v", len(cart.Items))
	}

	if cart.ItemCount != 0 {
		t.Errorf("ItemCount should be 0 after removal, got %v", cart.ItemCount)
	}
}

// TestRemoveFromCart_UpdatesTotals verifies that totals are recalculated after removal
func TestRemoveFromCart_UpdatesTotals(t *testing.T) {
	client := NewClient()

	// Add two different items
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("AddToCart() first item failed: %v", err)
	}

	_, err = client.AddToCart("B07XJ8C8F5", 1)
	if err != nil {
		t.Fatalf("AddToCart() second item failed: %v", err)
	}

	// Verify cart has 2 items
	cart, _ := client.GetCart()
	if len(cart.Items) != 2 {
		t.Fatalf("Cart should have 2 items, got %v", len(cart.Items))
	}

	initialTotal := cart.Total

	// Remove one item
	cart, err = client.RemoveFromCart("B08N5WRWNW")
	if err != nil {
		t.Fatalf("RemoveFromCart() failed: %v", err)
	}

	// Verify cart now has 1 item
	if len(cart.Items) != 1 {
		t.Errorf("Cart should have 1 item after removal, got %v", len(cart.Items))
	}

	// Verify remaining item is the correct one
	if cart.Items[0].ASIN != "B07XJ8C8F5" {
		t.Errorf("Remaining item ASIN = %v, want B07XJ8C8F5", cart.Items[0].ASIN)
	}

	// Verify subtotal decreased
	expectedSubtotal := 29.99
	if cart.Subtotal != expectedSubtotal {
		t.Errorf("Subtotal = %v, want %v", cart.Subtotal, expectedSubtotal)
	}

	// Verify total decreased
	if cart.Total >= initialTotal {
		t.Errorf("Total should decrease after removal, was %v, now %v", initialTotal, cart.Total)
	}

	// Verify tax recalculated
	expectedTax := expectedSubtotal * 0.08
	if cart.EstimatedTax != expectedTax {
		t.Errorf("EstimatedTax = %v, want %v", cart.EstimatedTax, expectedTax)
	}

	// Verify item count
	if cart.ItemCount != 1 {
		t.Errorf("ItemCount = %v, want 1", cart.ItemCount)
	}
}

// TestRemoveFromCart_NonExistentItem verifies removing a non-existent item doesn't error
func TestRemoveFromCart_NonExistentItem(t *testing.T) {
	client := NewClient()

	// Add one item
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("AddToCart() failed: %v", err)
	}

	// Try to remove a different item
	cart, err := client.RemoveFromCart("B07XJ8C8F5")
	if err != nil {
		t.Errorf("RemoveFromCart() should not error for non-existent item: %v", err)
	}

	// Verify original item is still in cart
	if len(cart.Items) != 1 {
		t.Errorf("Cart should still have 1 item, got %v", len(cart.Items))
	}

	if cart.Items[0].ASIN != "B08N5WRWNW" {
		t.Errorf("Remaining item ASIN = %v, want B08N5WRWNW", cart.Items[0].ASIN)
	}
}

// TestRemoveFromCart_MultipleQuantity verifies removing item with quantity > 1 updates count correctly
func TestRemoveFromCart_MultipleQuantity(t *testing.T) {
	client := NewClient()

	// Add item with quantity 5
	_, err := client.AddToCart("B08N5WRWNW", 5)
	if err != nil {
		t.Fatalf("AddToCart() failed: %v", err)
	}

	cart, _ := client.GetCart()
	if cart.ItemCount != 5 {
		t.Fatalf("ItemCount should be 5, got %v", cart.ItemCount)
	}

	// Remove the item (should remove all 5)
	cart, err = client.RemoveFromCart("B08N5WRWNW")
	if err != nil {
		t.Fatalf("RemoveFromCart() failed: %v", err)
	}

	// Verify all items removed
	if cart.ItemCount != 0 {
		t.Errorf("ItemCount should be 0 after removal, got %v", cart.ItemCount)
	}

	if len(cart.Items) != 0 {
		t.Errorf("Cart should be empty, got %v items", len(cart.Items))
	}
}

func TestClearCart(t *testing.T) {
	client := NewClient()
	err := client.ClearCart()

	if err != nil {
		t.Errorf("ClearCart() unexpected error: %v", err)
	}
}

// TestClearCart_EmptyCart verifies clearing an empty cart doesn't error
func TestClearCart_EmptyCart(t *testing.T) {
	client := NewClient()

	// Clear already empty cart
	err := client.ClearCart()
	if err != nil {
		t.Errorf("ClearCart() on empty cart should not error: %v", err)
	}

	// Verify cart is still empty
	cart, _ := client.GetCart()
	if len(cart.Items) != 0 {
		t.Errorf("Cart should be empty, got %v items", len(cart.Items))
	}

	if cart.ItemCount != 0 {
		t.Errorf("ItemCount should be 0, got %v", cart.ItemCount)
	}
}

// TestClearCart_WithItems verifies clearing cart removes all items
func TestClearCart_WithItems(t *testing.T) {
	client := NewClient()

	// Add multiple items
	_, err := client.AddToCart("B08N5WRWNW", 2)
	if err != nil {
		t.Fatalf("AddToCart() first item failed: %v", err)
	}

	_, err = client.AddToCart("B07XJ8C8F5", 1)
	if err != nil {
		t.Fatalf("AddToCart() second item failed: %v", err)
	}

	_, err = client.AddToCart("B09ABC123X", 3)
	if err != nil {
		t.Fatalf("AddToCart() third item failed: %v", err)
	}

	// Verify cart has items
	cart, _ := client.GetCart()
	if len(cart.Items) != 3 {
		t.Fatalf("Cart should have 3 items, got %v", len(cart.Items))
	}

	if cart.ItemCount != 6 {
		t.Fatalf("ItemCount should be 6, got %v", cart.ItemCount)
	}

	// Clear the cart
	err = client.ClearCart()
	if err != nil {
		t.Fatalf("ClearCart() failed: %v", err)
	}

	// Verify cart is empty
	cart, _ = client.GetCart()
	if len(cart.Items) != 0 {
		t.Errorf("Cart should be empty after clear, got %v items", len(cart.Items))
	}

	if cart.ItemCount != 0 {
		t.Errorf("ItemCount should be 0 after clear, got %v", cart.ItemCount)
	}
}

// TestClearCart_ResetsTotals verifies clearing cart resets all totals to zero
func TestClearCart_ResetsTotals(t *testing.T) {
	client := NewClient()

	// Add items to cart
	_, err := client.AddToCart("B08N5WRWNW", 5)
	if err != nil {
		t.Fatalf("AddToCart() failed: %v", err)
	}

	// Verify cart has non-zero totals
	cart, _ := client.GetCart()
	if cart.Subtotal == 0 {
		t.Fatal("Subtotal should be non-zero before clear")
	}
	if cart.EstimatedTax == 0 {
		t.Fatal("EstimatedTax should be non-zero before clear")
	}
	if cart.Total == 0 {
		t.Fatal("Total should be non-zero before clear")
	}

	// Clear the cart
	err = client.ClearCart()
	if err != nil {
		t.Fatalf("ClearCart() failed: %v", err)
	}

	// Verify all totals are zero
	cart, _ = client.GetCart()
	if cart.Subtotal != 0 {
		t.Errorf("Subtotal should be 0 after clear, got %v", cart.Subtotal)
	}

	if cart.EstimatedTax != 0 {
		t.Errorf("EstimatedTax should be 0 after clear, got %v", cart.EstimatedTax)
	}

	if cart.Total != 0 {
		t.Errorf("Total should be 0 after clear, got %v", cart.Total)
	}
}

// TestClearCart_ThenAddItem verifies cart works normally after being cleared
func TestClearCart_ThenAddItem(t *testing.T) {
	client := NewClient()

	// Add and clear
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("AddToCart() failed: %v", err)
	}

	err = client.ClearCart()
	if err != nil {
		t.Fatalf("ClearCart() failed: %v", err)
	}

	// Add item after clearing
	cart, err := client.AddToCart("B07XJ8C8F5", 2)
	if err != nil {
		t.Fatalf("AddToCart() after clear failed: %v", err)
	}

	// Verify cart has only the new item
	if len(cart.Items) != 1 {
		t.Errorf("Cart should have 1 item after clear and add, got %v", len(cart.Items))
	}

	if cart.Items[0].ASIN != "B07XJ8C8F5" {
		t.Errorf("Item ASIN = %v, want B07XJ8C8F5", cart.Items[0].ASIN)
	}

	if cart.ItemCount != 2 {
		t.Errorf("ItemCount should be 2, got %v", cart.ItemCount)
	}

	// Verify totals calculated correctly
	expectedSubtotal := 29.99 * 2
	if cart.Subtotal != expectedSubtotal {
		t.Errorf("Subtotal = %v, want %v", cart.Subtotal, expectedSubtotal)
	}
}

func TestGetAddresses(t *testing.T) {
	client := NewClient()
	addresses, err := client.GetAddresses()

	if err != nil {
		t.Errorf("GetAddresses() unexpected error: %v", err)
	}

	if addresses == nil {
		t.Error("GetAddresses() returned nil")
	}
}

func TestGetPaymentMethods(t *testing.T) {
	client := NewClient()
	methods, err := client.GetPaymentMethods()

	if err != nil {
		t.Errorf("GetPaymentMethods() unexpected error: %v", err)
	}

	if methods == nil {
		t.Error("GetPaymentMethods() returned nil")
	}
}

func TestPreviewCheckout(t *testing.T) {
	tests := []struct {
		name        string
		addressID   string
		paymentID   string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid preview",
			addressID:   "addr123",
			paymentID:   "pay123",
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "empty addressID should fail",
			addressID:   "",
			paymentID:   "pay123",
			wantErr:     true,
			errContains: "addressID cannot be empty",
		},
		{
			name:        "empty paymentID should fail",
			addressID:   "addr123",
			paymentID:   "",
			wantErr:     true,
			errContains: "paymentID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			preview, err := client.PreviewCheckout(tt.addressID, tt.paymentID)

			if tt.wantErr {
				if err == nil {
					t.Error("PreviewCheckout() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("PreviewCheckout() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("PreviewCheckout() unexpected error: %v", err)
			}

			if preview == nil {
				t.Error("PreviewCheckout() returned nil")
				return
			}

			if preview.Cart == nil {
				t.Error("PreviewCheckout() Cart is nil")
			}

			if preview.Address == nil {
				t.Error("PreviewCheckout() Address is nil")
			}

			if preview.PaymentMethod == nil {
				t.Error("PreviewCheckout() PaymentMethod is nil")
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
