package validation

import (
	"fmt"
	"regexp"
)

// asinPattern matches exactly 10 alphanumeric characters
var asinPattern = regexp.MustCompile(`^[A-Z0-9]{10}$`)

// ValidateASIN validates that an ASIN is exactly 10 alphanumeric characters
// ASIN (Amazon Standard Identification Number) format:
// - Must be exactly 10 characters long
// - Must contain only uppercase letters (A-Z) and digits (0-9)
//
// Returns an error if the ASIN is invalid, nil if valid.
func ValidateASIN(asin string) error {
	if asin == "" {
		return fmt.Errorf("ASIN cannot be empty")
	}

	if len(asin) != 10 {
		return fmt.Errorf("ASIN must be exactly 10 characters long, got %d", len(asin))
	}

	if !asinPattern.MatchString(asin) {
		return fmt.Errorf("ASIN must contain only uppercase letters and digits")
	}

	return nil
}

// IsValidASIN checks if an ASIN is valid without returning an error
// Returns true if the ASIN is valid, false otherwise.
func IsValidASIN(asin string) bool {
	return ValidateASIN(asin) == nil
}
