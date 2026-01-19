package amazon

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
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

func TestCompleteCheckout_NeverMakesRealHTTPPost(t *testing.T) {
	client := NewClient()

	// Setup cart with an item
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	// Replace the HTTP client with one that will fail if any requests are made
	// This ensures the mock implementation doesn't make external HTTP calls
	client.httpClient = &http.Client{
		Transport: &failingRoundTripper{},
	}

	// Complete checkout - should succeed without making HTTP calls
	confirmation, err := client.CompleteCheckout("addr123", "pay123")
	if err != nil {
		t.Fatalf("CompleteCheckout() error = %v, want no error (no HTTP calls should be made)", err)
	}

	// Verify we got a valid confirmation without making HTTP requests
	if confirmation == nil {
		t.Fatal("CompleteCheckout() returned nil confirmation")
	}

	if confirmation.OrderID == "" {
		t.Error("OrderID should not be empty even in mock implementation")
	}

	if confirmation.Total <= 0 {
		t.Error("Total should be greater than 0")
	}
}

// failingRoundTripper is an http.RoundTripper that fails if any HTTP request is attempted
type failingRoundTripper struct{}

func (f *failingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("HTTP request attempted but mock implementation should not make external calls: %s %s", req.Method, req.URL)
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
		return
	}

	// Empty cart should have zero items
	if cart.ItemCount != 0 {
		t.Errorf("GetCart() ItemCount = %v, want 0", cart.ItemCount)
	}
}

func TestRemoveFromCart_ActuallyRemovesItem(t *testing.T) {
	client := NewClient()

	// Add item to cart
	cart, err := client.AddToCart("B08N5WRWNW", 1)
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	// Verify count is 1
	if cart.ItemCount != 1 {
		t.Errorf("after adding item, ItemCount = %d, want 1", cart.ItemCount)
	}

	// Remove item from cart
	cart, err = client.RemoveFromCart("B08N5WRWNW")
	if err != nil {
		t.Fatalf("failed to remove from cart: %v", err)
	}

	// Verify count is 0
	if cart.ItemCount != 0 {
		t.Errorf("after removing item, ItemCount = %d, want 0", cart.ItemCount)
	}
}

func TestRemoveFromCart(t *testing.T) {
	tests := []struct {
		name        string
		setupCart   func(*Client) error
		asin        string
		wantErr     bool
		errContains string
		checkCart   func(*testing.T, *models.Cart)
	}{
		{
			name: "remove item from cart with single item",
			setupCart: func(c *Client) error {
				_, err := c.AddToCart("B08N5WRWNW", 1)
				return err
			},
			asin:        "B08N5WRWNW",
			wantErr:     false,
			errContains: "",
			checkCart: func(t *testing.T, cart *models.Cart) {
				if len(cart.Items) != 0 {
					t.Errorf("Expected 0 items in cart, got %d", len(cart.Items))
				}
				if cart.ItemCount != 0 {
					t.Errorf("Expected ItemCount 0, got %d", cart.ItemCount)
				}
				if cart.Subtotal != 0 {
					t.Errorf("Expected Subtotal 0, got %f", cart.Subtotal)
				}
				if cart.EstimatedTax != 0 {
					t.Errorf("Expected EstimatedTax 0, got %f", cart.EstimatedTax)
				}
				if cart.Total != 0 {
					t.Errorf("Expected Total 0, got %f", cart.Total)
				}
			},
		},
		{
			name: "remove one item from cart with multiple items",
			setupCart: func(c *Client) error {
				_, err := c.AddToCart("B08N5WRWNW", 1)
				if err != nil {
					return err
				}
				_, err = c.AddToCart("B07XJ8C8F5", 2)
				return err
			},
			asin:        "B08N5WRWNW",
			wantErr:     false,
			errContains: "",
			checkCart: func(t *testing.T, cart *models.Cart) {
				if len(cart.Items) != 1 {
					t.Errorf("Expected 1 item in cart, got %d", len(cart.Items))
				}
				if len(cart.Items) > 0 && cart.Items[0].ASIN != "B07XJ8C8F5" {
					t.Errorf("Expected remaining item ASIN B07XJ8C8F5, got %s", cart.Items[0].ASIN)
				}
				if cart.ItemCount != 2 {
					t.Errorf("Expected ItemCount 2, got %d", cart.ItemCount)
				}
				// Verify totals are recalculated correctly
				expectedSubtotal := 29.99 * 2
				if cart.Subtotal != expectedSubtotal {
					t.Errorf("Expected Subtotal %f, got %f", expectedSubtotal, cart.Subtotal)
				}
				expectedTax := expectedSubtotal * 0.08
				if cart.EstimatedTax != expectedTax {
					t.Errorf("Expected EstimatedTax %f, got %f", expectedTax, cart.EstimatedTax)
				}
				expectedTotal := expectedSubtotal + expectedTax
				if cart.Total != expectedTotal {
					t.Errorf("Expected Total %f, got %f", expectedTotal, cart.Total)
				}
			},
		},
		{
			name: "remove item with multiple quantity",
			setupCart: func(c *Client) error {
				_, err := c.AddToCart("B08N5WRWNW", 3)
				return err
			},
			asin:        "B08N5WRWNW",
			wantErr:     false,
			errContains: "",
			checkCart: func(t *testing.T, cart *models.Cart) {
				if len(cart.Items) != 0 {
					t.Errorf("Expected 0 items in cart, got %d", len(cart.Items))
				}
				if cart.ItemCount != 0 {
					t.Errorf("Expected ItemCount 0, got %d", cart.ItemCount)
				}
			},
		},
		{
			name:        "empty ASIN should fail",
			setupCart:   nil,
			asin:        "",
			wantErr:     true,
			errContains: "ASIN cannot be empty",
			checkCart:   nil,
		},
		{
			name:        "invalid ASIN format should fail",
			setupCart:   nil,
			asin:        "INVALID",
			wantErr:     true,
			errContains: "invalid ASIN format",
			checkCart:   nil,
		},
		{
			name: "removing non-existent item should fail",
			setupCart: func(c *Client) error {
				_, err := c.AddToCart("B08N5WRWNW", 1)
				return err
			},
			asin:        "B07XJ8C8F5",
			wantErr:     true,
			errContains: "not found in cart",
			checkCart:   nil,
		},
		{
			name:        "removing from empty cart should fail",
			setupCart:   nil,
			asin:        "B08N5WRWNW",
			wantErr:     true,
			errContains: "not found in cart",
			checkCart:   nil,
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

			// Execute RemoveFromCart
			cart, err := client.RemoveFromCart(tt.asin)

			// Check error expectations
			if tt.wantErr {
				if err == nil {
					t.Error("RemoveFromCart() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("RemoveFromCart() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			// Check success expectations
			if err != nil {
				t.Errorf("RemoveFromCart() unexpected error: %v", err)
				return
			}

			if cart == nil {
				t.Error("RemoveFromCart() returned nil cart")
				return
			}

			// Run additional cart checks if provided
			if tt.checkCart != nil {
				tt.checkCart(t, cart)
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

func TestClearCart_ResetsAllTotals(t *testing.T) {
	client := NewClient()

	// Add multiple items to the cart
	_, err := client.AddToCart("B08N5WRWNW", 3)
	if err != nil {
		t.Fatalf("failed to add first item to cart: %v", err)
	}

	_, err = client.AddToCart("B07XJ8C8F5", 2)
	if err != nil {
		t.Fatalf("failed to add second item to cart: %v", err)
	}

	_, err = client.AddToCart("B09ABCD123", 1)
	if err != nil {
		t.Fatalf("failed to add third item to cart: %v", err)
	}

	// Verify cart has items and non-zero totals before clearing
	cart, _ := client.GetCart()
	if cart.ItemCount == 0 {
		t.Fatal("cart.ItemCount should not be 0 before clearing")
	}
	if cart.Subtotal == 0 {
		t.Fatal("cart.Subtotal should not be 0 before clearing")
	}
	if cart.Total == 0 {
		t.Fatal("cart.Total should not be 0 before clearing")
	}

	// Clear the cart
	err = client.ClearCart()
	if err != nil {
		t.Fatalf("ClearCart() unexpected error: %v", err)
	}

	// Verify ItemCount is reset to 0
	cart, _ = client.GetCart()
	if cart.ItemCount != 0 {
		t.Errorf("After ClearCart(), ItemCount = %v, want 0", cart.ItemCount)
	}

	// Verify Subtotal is reset to 0
	if cart.Subtotal != 0 {
		t.Errorf("After ClearCart(), Subtotal = %v, want 0", cart.Subtotal)
	}

	// Verify Total is reset to 0
	if cart.Total != 0 {
		t.Errorf("After ClearCart(), Total = %v, want 0", cart.Total)
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

func TestPreviewCheckout_RealisticData(t *testing.T) {
	client := NewClient()

	// Add items to cart to test with realistic data
	_, err := client.AddToCart("B08N5WRWNW", 2)
	if err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	preview, err := client.PreviewCheckout("addr123", "pay123")
	if err != nil {
		t.Fatalf("PreviewCheckout() error = %v", err)
	}

	// Verify cart is returned with current contents
	if preview.Cart == nil {
		t.Fatal("Cart should not be nil")
	}
	if preview.Cart.ItemCount != 2 {
		t.Errorf("Cart.ItemCount = %d, want 2", preview.Cart.ItemCount)
	}
	if len(preview.Cart.Items) != 1 {
		t.Errorf("Cart.Items length = %d, want 1", len(preview.Cart.Items))
	}

	// Verify address has all fields populated
	if preview.Address == nil {
		t.Fatal("Address should not be nil")
	}
	if preview.Address.ID != "addr123" {
		t.Errorf("Address.ID = %q, want %q", preview.Address.ID, "addr123")
	}
	if preview.Address.Name == "" {
		t.Error("Address.Name should not be empty")
	}
	if preview.Address.Street == "" {
		t.Error("Address.Street should not be empty")
	}
	if preview.Address.City == "" {
		t.Error("Address.City should not be empty")
	}
	if preview.Address.State == "" {
		t.Error("Address.State should not be empty")
	}
	if preview.Address.Zip == "" {
		t.Error("Address.Zip should not be empty")
	}
	if preview.Address.Country == "" {
		t.Error("Address.Country should not be empty")
	}

	// Verify payment method has all fields populated
	if preview.PaymentMethod == nil {
		t.Fatal("PaymentMethod should not be nil")
	}
	if preview.PaymentMethod.ID != "pay123" {
		t.Errorf("PaymentMethod.ID = %q, want %q", preview.PaymentMethod.ID, "pay123")
	}
	if preview.PaymentMethod.Type == "" {
		t.Error("PaymentMethod.Type should not be empty")
	}
	if preview.PaymentMethod.Last4 == "" {
		t.Error("PaymentMethod.Last4 should not be empty")
	}
	if len(preview.PaymentMethod.Last4) != 4 {
		t.Errorf("PaymentMethod.Last4 length = %d, want 4", len(preview.PaymentMethod.Last4))
	}

	// Verify delivery options array has multiple options
	if len(preview.DeliveryOptions) == 0 {
		t.Error("DeliveryOptions should not be empty")
	}
	if len(preview.DeliveryOptions) < 2 {
		t.Errorf("DeliveryOptions length = %d, want at least 2 options", len(preview.DeliveryOptions))
	}
	for i, option := range preview.DeliveryOptions {
		if option == "" {
			t.Errorf("DeliveryOptions[%d] should not be empty", i)
		}
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
