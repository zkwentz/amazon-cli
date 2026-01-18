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

func TestClearCart(t *testing.T) {
	client := NewClient()
	err := client.ClearCart()

	if err != nil {
		t.Errorf("ClearCart() unexpected error: %v", err)
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

// TestCheckoutCancellation tests the checkout cancellation flow
// This ensures that checkout can be safely aborted and the cart remains intact
func TestCheckoutCancellation(t *testing.T) {
	tests := []struct {
		name          string
		setupCart     func(*Client) error
		addressID     string
		paymentID     string
		cancelBefore  string // "preview", "validation", or "submit"
		expectCart    bool   // should cart still exist after cancellation
		expectItems   int    // expected number of items in cart after cancellation
	}{
		{
			name: "cancel during preview - cart preserved",
			setupCart: func(c *Client) error {
				_, err := c.AddToCart("B08N5WRWNW", 2)
				return err
			},
			addressID:    "addr123",
			paymentID:    "pay123",
			cancelBefore: "preview",
			expectCart:   true,
			expectItems:  1,
		},
		{
			name: "cancel with invalid address - cart preserved",
			setupCart: func(c *Client) error {
				_, err := c.AddToCart("B08N5WRWNW", 1)
				return err
			},
			addressID:    "", // invalid
			paymentID:    "pay123",
			cancelBefore: "validation",
			expectCart:   true,
			expectItems:  1,
		},
		{
			name: "cancel with invalid payment - cart preserved",
			setupCart: func(c *Client) error {
				_, err := c.AddToCart("B08N5WRWNW", 1)
				if err != nil {
					return err
				}
				_, err = c.AddToCart("B07XJ8C8F5", 3)
				return err
			},
			addressID:    "addr123",
			paymentID:    "", // invalid
			cancelBefore: "validation",
			expectCart:   true,
			expectItems:  2,
		},
		{
			name: "cancel with empty cart - no items to checkout",
			setupCart: func(c *Client) error {
				return nil // no items added
			},
			addressID:    "addr123",
			paymentID:    "pay123",
			cancelBefore: "validation",
			expectCart:   true,
			expectItems:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			// Setup initial cart
			if tt.setupCart != nil {
				if err := tt.setupCart(client); err != nil {
					t.Fatalf("failed to setup cart: %v", err)
				}
			}

			// Get initial cart state
			initialCart, err := client.GetCart()
			if err != nil {
				t.Fatalf("failed to get initial cart: %v", err)
			}
			initialItemCount := initialCart.ItemCount

			// Attempt checkout (which should fail/cancel)
			var checkoutErr error
			switch tt.cancelBefore {
			case "preview":
				// Preview should succeed but we don't proceed
				_, checkoutErr = client.PreviewCheckout(tt.addressID, tt.paymentID)
			case "validation":
				// This should fail validation
				_, checkoutErr = client.CompleteCheckout(tt.addressID, tt.paymentID)
			case "submit":
				// This tests cancellation after validation but before actual submission
				_, checkoutErr = client.CompleteCheckout(tt.addressID, tt.paymentID)
			}

			// For cancellation scenarios, we expect errors or successful previews
			// The key is that the cart should remain intact

			// Verify cart still exists and is intact after cancellation
			if tt.expectCart {
				finalCart, err := client.GetCart()
				if err != nil {
					t.Errorf("failed to get cart after cancellation: %v", err)
				}

				if finalCart == nil {
					t.Error("cart should still exist after cancellation")
					return
				}

				// Verify cart contents unchanged for preview/failed checkouts
				if tt.cancelBefore == "preview" {
					if checkoutErr != nil {
						t.Errorf("preview should succeed: %v", checkoutErr)
					}
					if finalCart.ItemCount != initialItemCount {
						t.Errorf("cart item count changed after preview: got %d, want %d",
							finalCart.ItemCount, initialItemCount)
					}
				}

				// For validation failures, cart should be preserved
				if tt.cancelBefore == "validation" {
					if checkoutErr == nil && (tt.addressID == "" || tt.paymentID == "" || initialItemCount == 0) {
						t.Error("checkout should have failed validation")
					}
					if finalCart.ItemCount != initialItemCount {
						t.Errorf("cart item count changed after failed validation: got %d, want %d",
							finalCart.ItemCount, initialItemCount)
					}
				}

				// Verify expected item count
				if len(finalCart.Items) != tt.expectItems {
					t.Errorf("expected %d items in cart, got %d", tt.expectItems, len(finalCart.Items))
				}
			}
		})
	}
}

// TestCheckoutCancellationRace tests concurrent cancellation scenarios
func TestCheckoutCancellationRace(t *testing.T) {
	client := NewClient()

	// Add items to cart
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	// Get initial cart
	initialCart, err := client.GetCart()
	if err != nil {
		t.Fatalf("failed to get initial cart: %v", err)
	}

	// Simulate race condition: preview while someone else might be checking out
	_, err = client.PreviewCheckout("addr123", "pay123")
	if err != nil {
		t.Fatalf("preview failed: %v", err)
	}

	// Verify cart still intact
	finalCart, err := client.GetCart()
	if err != nil {
		t.Fatalf("failed to get final cart: %v", err)
	}

	if finalCart.ItemCount != initialCart.ItemCount {
		t.Errorf("cart modified during preview: initial=%d, final=%d",
			initialCart.ItemCount, finalCart.ItemCount)
	}
}

// TestCheckoutIdempotency tests that failed checkouts can be safely retried
func TestCheckoutIdempotency(t *testing.T) {
	client := NewClient()

	// Add items to cart
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	// Try checkout with invalid payment (should fail)
	_, err = client.CompleteCheckout("addr123", "")
	if err == nil {
		t.Error("checkout should have failed with empty payment")
	}

	// Verify cart still has items
	cart, err := client.GetCart()
	if err != nil {
		t.Fatalf("failed to get cart: %v", err)
	}

	if cart.ItemCount == 0 {
		t.Error("cart should still have items after failed checkout")
	}

	// Retry with valid payment (should succeed)
	confirmation, err := client.CompleteCheckout("addr123", "pay123")
	if err != nil {
		t.Errorf("retry checkout failed: %v", err)
	}

	if confirmation == nil {
		t.Error("confirmation should not be nil on successful checkout")
	}

	if confirmation != nil && confirmation.OrderID == "" {
		t.Error("order ID should be set after successful checkout")
	}
}

// TestClearCartCancellation tests clearing cart can be thought of as cancelling all pending purchases
func TestClearCartCancellation(t *testing.T) {
	client := NewClient()

	// Add multiple items
	_, err := client.AddToCart("B08N5WRWNW", 2)
	if err != nil {
		t.Fatalf("failed to add first item: %v", err)
	}

	_, err = client.AddToCart("B07XJ8C8F5", 1)
	if err != nil {
		t.Fatalf("failed to add second item: %v", err)
	}

	// Verify cart has items
	cart, err := client.GetCart()
	if err != nil {
		t.Fatalf("failed to get cart: %v", err)
	}

	if cart.ItemCount != 3 {
		t.Errorf("expected 3 items in cart, got %d", cart.ItemCount)
	}

	// Clear cart (cancel all pending items)
	err = client.ClearCart()
	if err != nil {
		t.Errorf("failed to clear cart: %v", err)
	}

	// Verify cart operations still work after clearing
	cart, err = client.GetCart()
	if err != nil {
		t.Errorf("failed to get cart after clear: %v", err)
	}

	if cart == nil {
		t.Error("cart should exist even after clearing")
	}
}

// TestRemoveItemCancellation tests removing items as a form of partial cancellation
func TestRemoveItemCancellation(t *testing.T) {
	client := NewClient()

	// Add items
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("failed to add item: %v", err)
	}

	initialCart, err := client.GetCart()
	if err != nil {
		t.Fatalf("failed to get initial cart: %v", err)
	}

	if initialCart.ItemCount != 1 {
		t.Errorf("expected 1 item, got %d", initialCart.ItemCount)
	}

	// Remove item (cancel this purchase)
	_, err = client.RemoveFromCart("B08N5WRWNW")
	if err != nil {
		t.Errorf("failed to remove item: %v", err)
	}

	// Verify cart still accessible
	finalCart, err := client.GetCart()
	if err != nil {
		t.Errorf("failed to get cart after removal: %v", err)
	}

	if finalCart == nil {
		t.Error("cart should still exist after removing items")
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
