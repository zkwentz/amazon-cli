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

// TestQuickBuy tests the quick buy flow which combines add to cart and checkout
func TestQuickBuy(t *testing.T) {
	tests := []struct {
		name        string
		asin        string
		quantity    int
		addressID   string
		paymentID   string
		wantErr     bool
		errContains string
	}{
		{
			name:      "valid quick buy single item",
			asin:      "B08N5WRWNW",
			quantity:  1,
			addressID: "addr123",
			paymentID: "pay123",
			wantErr:   false,
		},
		{
			name:      "valid quick buy multiple quantity",
			asin:      "B08N5WRWNW",
			quantity:  3,
			addressID: "addr456",
			paymentID: "pay456",
			wantErr:   false,
		},
		{
			name:        "empty ASIN should fail",
			asin:        "",
			quantity:    1,
			addressID:   "addr123",
			paymentID:   "pay123",
			wantErr:     true,
			errContains: "ASIN cannot be empty",
		},
		{
			name:        "zero quantity should fail",
			asin:        "B08N5WRWNW",
			quantity:    0,
			addressID:   "addr123",
			paymentID:   "pay123",
			wantErr:     true,
			errContains: "quantity must be positive",
		},
		{
			name:        "negative quantity should fail",
			asin:        "B08N5WRWNW",
			quantity:    -1,
			addressID:   "addr123",
			paymentID:   "pay123",
			wantErr:     true,
			errContains: "quantity must be positive",
		},
		{
			name:        "empty addressID should fail",
			asin:        "B08N5WRWNW",
			quantity:    1,
			addressID:   "",
			paymentID:   "pay123",
			wantErr:     true,
			errContains: "addressID cannot be empty",
		},
		{
			name:        "empty paymentID should fail",
			asin:        "B08N5WRWNW",
			quantity:    1,
			addressID:   "addr123",
			paymentID:   "",
			wantErr:     true,
			errContains: "paymentID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			// Execute QuickBuy
			confirmation, err := client.QuickBuy(tt.asin, tt.quantity, tt.addressID, tt.paymentID)

			// Check error expectations
			if tt.wantErr {
				if err == nil {
					t.Errorf("QuickBuy() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("QuickBuy() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			// Check success expectations
			if err != nil {
				t.Errorf("QuickBuy() unexpected error: %v", err)
				return
			}

			if confirmation == nil {
				t.Error("QuickBuy() returned nil confirmation")
				return
			}

			// Verify confirmation fields
			if confirmation.OrderID == "" {
				t.Error("QuickBuy() OrderID is empty")
			}

			if confirmation.Total <= 0 {
				t.Error("QuickBuy() Total should be greater than 0")
			}

			if confirmation.EstimatedDelivery == "" {
				t.Error("QuickBuy() EstimatedDelivery is empty")
			}
		})
	}
}

// TestQuickBuy_OrderConfirmation verifies the order confirmation structure
func TestQuickBuy_OrderConfirmation(t *testing.T) {
	client := NewClient()

	// Execute quick buy
	confirmation, err := client.QuickBuy("B08N5WRWNW", 2, "addr123", "pay123")
	if err != nil {
		t.Fatalf("QuickBuy() error = %v", err)
	}

	// Verify order confirmation structure
	if confirmation.OrderID == "" {
		t.Error("OrderID should not be empty")
	}

	// Check total is approximately correct (29.99 * 2 = 59.98 + 8% tax = 64.7784)
	expectedTotal := 64.78
	if confirmation.Total < expectedTotal-0.01 || confirmation.Total > expectedTotal+0.01 {
		t.Errorf("Total = %v, want approximately %v", confirmation.Total, expectedTotal)
	}

	if confirmation.EstimatedDelivery == "" {
		t.Error("EstimatedDelivery should not be empty")
	}
}

// TestQuickBuy_CartState verifies cart state after quick buy
func TestQuickBuy_CartState(t *testing.T) {
	client := NewClient()

	// Check initial cart is empty
	initialCart, err := client.GetCart()
	if err != nil {
		t.Fatalf("GetCart() error = %v", err)
	}
	if initialCart.ItemCount != 0 {
		t.Fatalf("Initial cart should be empty, got %d items", initialCart.ItemCount)
	}

	// Execute quick buy
	_, err = client.QuickBuy("B08N5WRWNW", 1, "addr123", "pay123")
	if err != nil {
		t.Fatalf("QuickBuy() error = %v", err)
	}

	// Verify cart now has items (quick buy adds to cart)
	finalCart, err := client.GetCart()
	if err != nil {
		t.Fatalf("GetCart() error = %v", err)
	}
	if finalCart.ItemCount == 0 {
		t.Error("Cart should have items after quick buy")
	}
	if len(finalCart.Items) == 0 {
		t.Error("Cart items should not be empty after quick buy")
	}
}

// TestQuickBuy_MultipleItems tests quick buy with different ASINs
func TestQuickBuy_MultipleItems(t *testing.T) {
	testCases := []struct {
		name     string
		asin     string
		quantity int
	}{
		{
			name:     "single item",
			asin:     "B08N5WRWNW",
			quantity: 1,
		},
		{
			name:     "multiple of same item",
			asin:     "B08N5WRWNW",
			quantity: 5,
		},
		{
			name:     "different ASIN",
			asin:     "B07XJ8C8F5",
			quantity: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient()

			confirmation, err := client.QuickBuy(tc.asin, tc.quantity, "addr123", "pay123")
			if err != nil {
				t.Fatalf("QuickBuy() error = %v", err)
			}

			if confirmation == nil {
				t.Fatal("QuickBuy() returned nil confirmation")
			}

			// Verify ASIN is in cart
			cart, err := client.GetCart()
			if err != nil {
				t.Fatalf("GetCart() error = %v", err)
			}

			found := false
			for _, item := range cart.Items {
				if item.ASIN == tc.asin && item.Quantity == tc.quantity {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected ASIN %s with quantity %d not found in cart", tc.asin, tc.quantity)
			}
		})
	}
}

// TestQuickBuy_Integration tests the complete quick buy flow
func TestQuickBuy_Integration(t *testing.T) {
	client := NewClient()

	// Test data
	asin := "B08N5WRWNW"
	quantity := 1
	addressID := "addr123"
	paymentID := "pay123"

	// Execute quick buy
	confirmation, err := client.QuickBuy(asin, quantity, addressID, paymentID)
	if err != nil {
		t.Fatalf("QuickBuy() error = %v", err)
	}

	// Verify all confirmation fields are populated
	if confirmation.OrderID == "" {
		t.Error("OrderID should be populated")
	}

	if confirmation.Total <= 0 {
		t.Error("Total should be greater than 0")
	}

	if confirmation.EstimatedDelivery == "" {
		t.Error("EstimatedDelivery should be populated")
	}

	// Verify cart reflects the purchase
	cart, err := client.GetCart()
	if err != nil {
		t.Fatalf("GetCart() error = %v", err)
	}

	if len(cart.Items) == 0 {
		t.Error("Cart should contain items after quick buy")
	}

	// Verify the correct item is in cart
	foundItem := false
	for _, item := range cart.Items {
		if item.ASIN == asin {
			foundItem = true
			if item.Quantity != quantity {
				t.Errorf("Item quantity = %d, want %d", item.Quantity, quantity)
			}
		}
	}
	if !foundItem {
		t.Errorf("ASIN %s not found in cart after quick buy", asin)
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
