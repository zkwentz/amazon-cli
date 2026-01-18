package amazon

import (
	"errors"
	"fmt"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// Client represents the Amazon API client
type Client struct {
	// This would contain http.Client, rate limiter, config, etc.
	// For now, we'll keep it minimal for the AddToCart implementation
}

// AddToCart adds an item to the shopping cart
// Parameters:
//   - asin: Amazon Standard Identification Number (10 alphanumeric characters)
//   - quantity: Number of items to add (must be positive)
// Returns:
//   - *Cart: Updated cart with all items and totals
//   - error: Any error that occurred during the operation
func (c *Client) AddToCart(asin string, quantity int) (*models.Cart, error) {
	// Validate ASIN format
	if err := validateASIN(asin); err != nil {
		return nil, err
	}

	// Validate quantity
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	// In a real implementation, this would:
	// 1. Call rate limiter Wait() before request
	// 2. Build the add-to-cart request with ASIN and quantity
	// 3. Submit POST request to Amazon's add-to-cart endpoint
	// 4. Parse the response to extract updated cart information
	// 5. Handle errors (item out of stock, quantity limits, etc.)
	// 6. Return the updated Cart struct

	// For now, return a mock cart response
	// This simulates a successful add-to-cart operation
	cart := &models.Cart{
		Items: []models.CartItem{
			{
				ASIN:     asin,
				Title:    "Sample Product",
				Price:    29.99,
				Quantity: quantity,
				Subtotal: 29.99 * float64(quantity),
				Prime:    true,
				InStock:  true,
			},
		},
		Subtotal:     29.99 * float64(quantity),
		EstimatedTax: 29.99 * float64(quantity) * 0.08, // 8% tax estimate
		ItemCount:    quantity,
	}
	cart.Total = cart.Subtotal + cart.EstimatedTax

	return cart, nil
}

// validateASIN checks if the ASIN is in valid format
// Amazon ASINs are 10 alphanumeric characters
func validateASIN(asin string) error {
	if len(asin) != 10 {
		return fmt.Errorf("invalid ASIN format: must be 10 characters, got %d", len(asin))
	}

	// Check if all characters are alphanumeric
	for _, char := range asin {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return fmt.Errorf("invalid ASIN format: must contain only uppercase letters and numbers")
		}
	}

	return nil
}
