package validation

import "fmt"

// ValidatePriceRange validates that min price is less than max price and both are positive
// Returns an error if validation fails
func ValidatePriceRange(minPrice, maxPrice float64) error {
	// Check if both prices are positive
	if minPrice < 0 {
		return fmt.Errorf("min price must be positive, got %.2f", minPrice)
	}

	if maxPrice < 0 {
		return fmt.Errorf("max price must be positive, got %.2f", maxPrice)
	}

	// Check if min is less than max
	if minPrice >= maxPrice {
		return fmt.Errorf("min price (%.2f) must be less than max price (%.2f)", minPrice, maxPrice)
	}

	return nil
}
