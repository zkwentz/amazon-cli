package amazon

import (
	"testing"
)

func TestAddToCart(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name        string
		asin        string
		quantity    int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid ASIN and quantity",
			asin:        "B08N5WRWNW",
			quantity:    1,
			expectError: false,
		},
		{
			name:        "valid ASIN with multiple quantity",
			asin:        "B07XJ8C8F5",
			quantity:    3,
			expectError: false,
		},
		{
			name:        "invalid ASIN - too short",
			asin:        "B08N5WRW",
			quantity:    1,
			expectError: true,
			errorMsg:    "invalid ASIN format: must be 10 characters, got 8",
		},
		{
			name:        "invalid ASIN - too long",
			asin:        "B08N5WRWNWX",
			quantity:    1,
			expectError: true,
			errorMsg:    "invalid ASIN format: must be 10 characters, got 11",
		},
		{
			name:        "invalid ASIN - lowercase letters",
			asin:        "b08n5wrwnw",
			quantity:    1,
			expectError: true,
			errorMsg:    "invalid ASIN format: must contain only uppercase letters and numbers",
		},
		{
			name:        "invalid ASIN - special characters",
			asin:        "B08N5WRW-W",
			quantity:    1,
			expectError: true,
			errorMsg:    "invalid ASIN format: must contain only uppercase letters and numbers",
		},
		{
			name:        "invalid quantity - zero",
			asin:        "B08N5WRWNW",
			quantity:    0,
			expectError: true,
			errorMsg:    "quantity must be positive",
		},
		{
			name:        "invalid quantity - negative",
			asin:        "B08N5WRWNW",
			quantity:    -1,
			expectError: true,
			errorMsg:    "quantity must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, err := client.AddToCart(tt.asin, tt.quantity)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cart == nil {
					t.Errorf("expected cart to be non-nil")
				}
				if cart != nil {
					// Verify cart structure
					if len(cart.Items) == 0 {
						t.Errorf("expected cart to have items")
					}
					if cart.ItemCount != tt.quantity {
						t.Errorf("expected item count %d, got %d", tt.quantity, cart.ItemCount)
					}
					if cart.Items[0].ASIN != tt.asin {
						t.Errorf("expected ASIN %s, got %s", tt.asin, cart.Items[0].ASIN)
					}
					if cart.Items[0].Quantity != tt.quantity {
						t.Errorf("expected quantity %d, got %d", tt.quantity, cart.Items[0].Quantity)
					}
					// Verify totals are calculated
					if cart.Total <= 0 {
						t.Errorf("expected positive total, got %f", cart.Total)
					}
					if cart.Subtotal <= 0 {
						t.Errorf("expected positive subtotal, got %f", cart.Subtotal)
					}
				}
			}
		})
	}
}

func TestValidateASIN(t *testing.T) {
	tests := []struct {
		name        string
		asin        string
		expectError bool
	}{
		{
			name:        "valid ASIN",
			asin:        "B08N5WRWNW",
			expectError: false,
		},
		{
			name:        "valid ASIN with numbers only",
			asin:        "1234567890",
			expectError: false,
		},
		{
			name:        "valid ASIN with letters only",
			asin:        "ABCDEFGHIJ",
			expectError: false,
		},
		{
			name:        "invalid ASIN - too short",
			asin:        "B08N5",
			expectError: true,
		},
		{
			name:        "invalid ASIN - too long",
			asin:        "B08N5WRWNWX",
			expectError: true,
		},
		{
			name:        "invalid ASIN - lowercase",
			asin:        "b08n5wrwnw",
			expectError: true,
		},
		{
			name:        "invalid ASIN - special character",
			asin:        "B08N5WRW@W",
			expectError: true,
		},
		{
			name:        "invalid ASIN - space",
			asin:        "B08N5WRW W",
			expectError: true,
		},
		{
			name:        "empty ASIN",
			asin:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateASIN(tt.asin)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
