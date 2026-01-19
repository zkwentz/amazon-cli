package validation

import (
	"fmt"
	"regexp"
	"unicode"
)

// ValidateASIN validates that an ASIN is exactly 10 alphanumeric characters
func ValidateASIN(asin string) error {
	if len(asin) != 10 {
		return fmt.Errorf("ASIN must be exactly 10 characters, got %d", len(asin))
	}

	for _, char := range asin {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return fmt.Errorf("ASIN must contain only alphanumeric characters")
		}
	}

	return nil
}

// ValidateOrderID validates that an order ID follows the format XXX-XXXXXXX-XXXXXXX
func ValidateOrderID(id string) error {
	// Pattern: XXX-XXXXXXX-XXXXXXX (3 digits, dash, 7 digits, dash, 7 digits)
	pattern := `^\d{3}-\d{7}-\d{7}$`
	matched, err := regexp.MatchString(pattern, id)
	if err != nil {
		return fmt.Errorf("error validating order ID: %w", err)
	}

	if !matched {
		return fmt.Errorf("order ID must follow format XXX-XXXXXXX-XXXXXXX (digits only)")
	}

	return nil
}

// ValidateQuantity validates that quantity is between 1 and 999
func ValidateQuantity(qty int) error {
	if qty < 1 {
		return fmt.Errorf("quantity must be at least 1, got %d", qty)
	}

	if qty > 999 {
		return fmt.Errorf("quantity must not exceed 999, got %d", qty)
	}

	return nil
}

// ValidatePriceRange validates that min >= 0 and max > min
func ValidatePriceRange(min, max float64) error {
	if min < 0 {
		return fmt.Errorf("minimum price must be non-negative, got %.2f", min)
	}

	if max <= min {
		return fmt.Errorf("maximum price (%.2f) must be greater than minimum price (%.2f)", max, min)
	}

	return nil
}
