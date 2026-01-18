package amazon

import (
	"testing"
)

func TestRemoveFromCart(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		asin        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid ASIN",
			asin:        "B08N5WRWNW",
			expectError: false,
		},
		{
			name:        "Empty ASIN",
			asin:        "",
			expectError: true,
			errorMsg:    "ASIN cannot be empty",
		},
		{
			name:        "Invalid ASIN - too short",
			asin:        "B08N5",
			expectError: true,
			errorMsg:    "invalid ASIN format",
		},
		{
			name:        "Invalid ASIN - too long",
			asin:        "B08N5WRWNW123",
			expectError: true,
			errorMsg:    "invalid ASIN format",
		},
		{
			name:        "Valid ASIN with different format",
			asin:        "1234567890",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, err := client.RemoveFromCart(tt.asin)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					// Check if error message contains expected text
					if len(err.Error()) < len(tt.errorMsg) || err.Error()[:len(tt.errorMsg)] != tt.errorMsg {
						t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if cart == nil {
					t.Error("Expected cart to be returned, got nil")
				}
				if cart != nil {
					// Verify cart structure
					if cart.Items == nil {
						t.Error("Expected cart.Items to be initialized")
					}
				}
			}
		})
	}
}

func TestRemoveFromCart_ReturnValue(t *testing.T) {
	client := NewClient()

	cart, err := client.RemoveFromCart("B08N5WRWNW")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify cart structure is properly initialized
	if cart == nil {
		t.Fatal("Cart should not be nil")
	}

	if cart.Items == nil {
		t.Error("Cart.Items should be initialized")
	}

	if cart.Subtotal < 0 {
		t.Error("Cart.Subtotal should not be negative")
	}

	if cart.EstimatedTax < 0 {
		t.Error("Cart.EstimatedTax should not be negative")
	}

	if cart.Total < 0 {
		t.Error("Cart.Total should not be negative")
	}

	if cart.ItemCount < 0 {
		t.Error("Cart.ItemCount should not be negative")
	}
}

func TestAddToCart(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		asin        string
		quantity    int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid add",
			asin:        "B08N5WRWNW",
			quantity:    1,
			expectError: false,
		},
		{
			name:        "Empty ASIN",
			asin:        "",
			quantity:    1,
			expectError: true,
			errorMsg:    "ASIN cannot be empty",
		},
		{
			name:        "Zero quantity",
			asin:        "B08N5WRWNW",
			quantity:    0,
			expectError: true,
			errorMsg:    "quantity must be greater than 0",
		},
		{
			name:        "Negative quantity",
			asin:        "B08N5WRWNW",
			quantity:    -1,
			expectError: true,
			errorMsg:    "quantity must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, err := client.AddToCart(tt.asin, tt.quantity)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if cart == nil {
					t.Error("Expected cart to be returned, got nil")
				}
			}
		})
	}
}

func TestGetCart(t *testing.T) {
	client := NewClient()

	cart, err := client.GetCart()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cart == nil {
		t.Fatal("Cart should not be nil")
	}

	// Verify cart structure
	if cart.Items == nil {
		t.Error("Cart.Items should be initialized")
	}
}
