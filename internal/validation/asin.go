package validation

import (
	"fmt"
	"regexp"
)

// ASIN format: 10 alphanumeric characters
// Examples: B08N5WRWNW, B00EXAMPLE
var asinRegex = regexp.MustCompile(`^[A-Z0-9]{10}$`)

// ValidateASIN validates that a string is a valid Amazon Standard Identification Number (ASIN)
// ASINs are exactly 10 alphanumeric characters (uppercase letters and digits)
func ValidateASIN(asin string) error {
	if asin == "" {
		return fmt.Errorf("ASIN cannot be empty")
	}

	if len(asin) != 10 {
		return fmt.Errorf("ASIN must be exactly 10 characters long, got %d", len(asin))
	}

	if !asinRegex.MatchString(asin) {
		return fmt.Errorf("ASIN must contain only uppercase letters and digits, got: %s", asin)
	}

	return nil
}

// IsValidASIN returns true if the string is a valid ASIN, false otherwise
func IsValidASIN(asin string) bool {
	return ValidateASIN(asin) == nil
}
