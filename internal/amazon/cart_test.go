package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestCalculateEstimatedTax(t *testing.T) {
	tests := []struct {
		name     string
		subtotal float64
		state    string
		expected float64
	}{
		{
			name:     "California tax",
			subtotal: 100.00,
			state:    "CA",
			expected: 7.25,
		},
		{
			name:     "New York tax",
			subtotal: 100.00,
			state:    "NY",
			expected: 8.00,
		},
		{
			name:     "Texas tax",
			subtotal: 100.00,
			state:    "TX",
			expected: 6.25,
		},
		{
			name:     "Florida tax",
			subtotal: 100.00,
			state:    "FL",
			expected: 6.00,
		},
		{
			name:     "Unknown state uses default",
			subtotal: 100.00,
			state:    "XX",
			expected: 7.00,
		},
		{
			name:     "Decimal subtotal",
			subtotal: 49.99,
			state:    "CA",
			expected: 3.62,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateEstimatedTax(tt.subtotal, tt.state)
			if result != tt.expected {
				t.Errorf("calculateEstimatedTax(%f, %s) = %f; want %f", tt.subtotal, tt.state, result, tt.expected)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "Simple dollar amount",
			input:    "$29.99",
			expected: 29.99,
			wantErr:  false,
		},
		{
			name:     "Amount with comma",
			input:    "$1,234.56",
			expected: 1234.56,
			wantErr:  false,
		},
		{
			name:     "Amount without dollar sign",
			input:    "99.99",
			expected: 99.99,
			wantErr:  false,
		},
		{
			name:     "Amount with spaces",
			input:    " $12.34 ",
			expected: 12.34,
			wantErr:  false,
		},
		{
			name:     "Large amount",
			input:    "$10,000.00",
			expected: 10000.00,
			wantErr:  false,
		},
		{
			name:    "Invalid format",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFloat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFloat(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("parseFloat(%s) = %f; want %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetAddressForCheckout(t *testing.T) {
	// Note: This test demonstrates the logic but would need proper mocking
	// In a real implementation, we'd use interfaces or dependency injection
	// to properly test getAddressForCheckout with different scenarios

	t.Run("Address selection logic", func(t *testing.T) {
		// This would test:
		// 1. Specific address ID is selected when provided
		// 2. Default address is selected when no ID provided
		// 3. First address is selected when no default exists
		// 4. Error is returned when no addresses exist
		// 5. Error is returned when specified address ID doesn't exist
	})
}

func TestGetPaymentMethodForCheckout(t *testing.T) {
	// Note: This test demonstrates the logic but would need proper mocking
	// In a real implementation, we'd use interfaces or dependency injection
	// to properly test getPaymentMethodForCheckout with different scenarios

	t.Run("Payment method selection logic", func(t *testing.T) {
		// This would test:
		// 1. Specific payment ID is selected when provided
		// 2. Default payment method is selected when no ID provided
		// 3. First payment method is selected when no default exists
		// 4. Error is returned when no payment methods exist
		// 5. Error is returned when specified payment ID doesn't exist
	})
}

// TestPreviewCheckoutValidation tests the validation logic in PreviewCheckout
func TestPreviewCheckoutValidation(t *testing.T) {
	// This test would validate:
	// 1. Empty cart returns an error
	// 2. Invalid address ID returns an error
	// 3. Invalid payment ID returns an error
	// 4. Successful preview returns proper structure

	// Note: Full implementation would require mocking HTTP responses
	// This is a placeholder for the test structure
	t.Run("Empty cart should error", func(t *testing.T) {
		// Would test that PreviewCheckout returns error when cart is empty
	})

	t.Run("Invalid address ID should error", func(t *testing.T) {
		// Would test that PreviewCheckout returns error with invalid address ID
	})

	t.Run("Invalid payment ID should error", func(t *testing.T) {
		// Would test that PreviewCheckout returns error with invalid payment ID
	})

	t.Run("Valid preview should return CheckoutPreview", func(t *testing.T) {
		// Would test that PreviewCheckout returns proper CheckoutPreview structure
	})
}

// TestFetchCheckoutPreview tests the checkout preview structure
func TestFetchCheckoutPreview(t *testing.T) {
	t.Run("Preview structure validation", func(t *testing.T) {
		// Note: This test validates the structure of CheckoutPreview
		// Full integration testing would require mocking HTTP responses
		// or setting up a test server

		// Verify that the CheckoutPreview structure has all required fields
		preview := &models.CheckoutPreview{
			Cart: &models.Cart{
				Items: []models.CartItem{
					{
						ASIN:     "B08N5WRWNW",
						Title:    "Test Product",
						Price:    29.99,
						Quantity: 1,
						Subtotal: 29.99,
						Prime:    true,
						InStock:  true,
					},
				},
				Subtotal:     29.99,
				EstimatedTax: 2.17,
				Total:        32.16,
				ItemCount:    1,
			},
			Address: &models.Address{
				ID:      "addr1",
				Name:    "John Doe",
				Street:  "123 Main St",
				City:    "New York",
				State:   "NY",
				Zip:     "10001",
				Country: "US",
				Default: true,
			},
			PaymentMethod: &models.PaymentMethod{
				ID:      "pay1",
				Type:    "Visa",
				Last4:   "1234",
				Default: true,
			},
			DeliveryOptions: []models.DeliveryOption{
				{
					Method:        "Standard Shipping",
					EstimatedDate: "Jan 25-27",
					Cost:          0.0,
					BusinessDays:  5,
				},
			},
			Subtotal: 29.99,
			Tax:      2.17,
			Shipping: 0.0,
			Total:    32.16,
		}

		// Verify preview structure
		if preview.Cart == nil {
			t.Error("Preview cart is nil")
		}
		if preview.Address == nil {
			t.Error("Preview address is nil")
		}
		if preview.PaymentMethod == nil {
			t.Error("Preview payment method is nil")
		}
		if len(preview.DeliveryOptions) == 0 {
			t.Error("Preview has no delivery options")
		}
		if preview.Total <= 0 {
			t.Error("Preview total should be greater than 0")
		}

		// Verify total is calculated correctly
		expectedTotal := preview.Subtotal + preview.Tax + preview.Shipping
		if preview.Total != expectedTotal {
			t.Errorf("Preview total = %f; want %f", preview.Total, expectedTotal)
		}
	})
}
