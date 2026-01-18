package validation

import (
	"fmt"
	"regexp"
)

// Order ID format: XXX-XXXXXXX-XXXXXXX (3 digits - 7 digits - 7 digits)
var orderIDRegex = regexp.MustCompile(`^\d{3}-\d{7}-\d{7}$`)

// ValidateOrderID validates an Amazon order ID format
// Amazon order IDs follow the format: XXX-XXXXXXX-XXXXXXX
// Example: 123-4567890-1234567
func ValidateOrderID(orderID string) error {
	if orderID == "" {
		return fmt.Errorf("order ID cannot be empty")
	}

	if !orderIDRegex.MatchString(orderID) {
		return fmt.Errorf("invalid order ID format: expected XXX-XXXXXXX-XXXXXXX (e.g., 123-4567890-1234567), got %s", orderID)
	}

	return nil
}
