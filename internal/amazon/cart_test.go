package amazon

import (
	"testing"
)

func TestValidateASIN(t *testing.T) {
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
		{
			name:      "ASIN too short should fail",
			asin:      "B08N5WRW",
			wantErr:   true,
			errString: "invalid ASIN format: must be 10 alphanumeric characters",
		},
		{
			name:      "ASIN too long should fail",
			asin:      "B08N5WRWNW1",
			wantErr:   true,
			errString: "invalid ASIN format: must be 10 alphanumeric characters",
		},
		{
			name:      "ASIN with lowercase should fail",
			asin:      "b08n5wrwnw",
			wantErr:   true,
			errString: "invalid ASIN format: must be 10 alphanumeric characters",
		},
		{
			name:      "ASIN with special characters should fail",
			asin:      "B08N5WRWN!",
			wantErr:   true,
			errString: "invalid ASIN format: must be 10 alphanumeric characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateASIN(tt.asin)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateASIN() expected error but got none")
				} else if err.Error() != tt.errString {
					t.Errorf("ValidateASIN() error = %v, want %v", err.Error(), tt.errString)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateASIN() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		name      string
		quantity  int
		wantErr   bool
		errString string
	}{
		{
			name:      "valid quantity 1",
			quantity:  1,
			wantErr:   false,
			errString: "",
		},
		{
			name:      "valid quantity 5",
			quantity:  5,
			wantErr:   false,
			errString: "",
		},
		{
			name:      "zero quantity should fail",
			quantity:  0,
			wantErr:   true,
			errString: "quantity must be positive",
		},
		{
			name:      "negative quantity should fail",
			quantity:  -1,
			wantErr:   true,
			errString: "quantity must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuantity(tt.quantity)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateQuantity() expected error but got none")
				} else if err.Error() != tt.errString {
					t.Errorf("ValidateQuantity() error = %v, want %v", err.Error(), tt.errString)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateQuantity() unexpected error: %v", err)
			}
		})
	}
}

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
			name:      "invalid ASIN format should fail",
			asin:      "INVALID",
			quantity:  1,
			wantErr:   true,
			errString: "invalid ASIN format: must be 10 alphanumeric characters",
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
			name:      "empty ASIN should fail",
			asin:      "",
			wantErr:   true,
			errString: "ASIN cannot be empty",
		},
		{
			name:      "invalid ASIN format should fail",
			asin:      "INVALID",
			wantErr:   true,
			errString: "invalid ASIN format: must be 10 alphanumeric characters",
		},
		{
			name:      "valid ASIN not in cart",
			asin:      "B08N5WRWNW",
			wantErr:   true,
			errString: "item with ASIN B08N5WRWNW not found in cart",
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

	// Add items to the cart first
	_, err := client.AddToCart("B08N5WRWNW", 2)
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	_, err = client.AddToCart("B07XJ8C8F5", 3)
	if err != nil {
		t.Fatalf("failed to add second item to cart: %v", err)
	}

	// Verify cart has items before clearing
	cart, _ := client.GetCart()
	if cart.ItemCount == 0 {
		t.Fatal("cart should have items before clearing")
	}
	if len(cart.Items) == 0 {
		t.Fatal("cart.Items should not be empty before clearing")
	}
	if cart.Subtotal == 0 {
		t.Fatal("cart.Subtotal should not be 0 before clearing")
	}
	if cart.EstimatedTax == 0 {
		t.Fatal("cart.EstimatedTax should not be 0 before clearing")
	}
	if cart.Total == 0 {
		t.Fatal("cart.Total should not be 0 before clearing")
	}

	// Clear the cart
	err = client.ClearCart()
	if err != nil {
		t.Errorf("ClearCart() unexpected error: %v", err)
	}

	// Verify cart is completely reset
	cart, _ = client.GetCart()
	if len(cart.Items) != 0 {
		t.Errorf("ClearCart() Items length = %v, want 0", len(cart.Items))
	}
	if cart.ItemCount != 0 {
		t.Errorf("ClearCart() ItemCount = %v, want 0", cart.ItemCount)
	}
	if cart.Subtotal != 0 {
		t.Errorf("ClearCart() Subtotal = %v, want 0", cart.Subtotal)
	}
	if cart.EstimatedTax != 0 {
		t.Errorf("ClearCart() EstimatedTax = %v, want 0", cart.EstimatedTax)
	}
	if cart.Total != 0 {
		t.Errorf("ClearCart() Total = %v, want 0", cart.Total)
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
