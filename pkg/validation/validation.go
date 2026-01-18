package validation

import (
	"fmt"
	"regexp"
)

// ASIN validation pattern (10 alphanumeric characters)
var asinPattern = regexp.MustCompile(`^[A-Z0-9]{10}$`)

// Order ID pattern (Amazon order IDs are in format: 123-1234567-1234567)
var orderIDPattern = regexp.MustCompile(`^\d{3}-\d{7}-\d{7}$`)

// Subscription ID pattern (Amazon subscription IDs are in format: S01-1234567-1234567)
var subscriptionIDPattern = regexp.MustCompile(`^S\d{2}-\d{7}-\d{7}$`)

// Constants for validation limits
const (
	MinQuantity = 1
	MaxQuantity = 999
	MinPrice    = 0.01
	MaxPrice    = 999999.99
)

// ValidateASIN validates that an ASIN is in the correct format
func ValidateASIN(asin string) error {
	if asin == "" {
		return fmt.Errorf("ASIN cannot be empty")
	}
	if !asinPattern.MatchString(asin) {
		return fmt.Errorf("invalid ASIN format: must be 10 alphanumeric characters")
	}
	return nil
}

// ValidateQuantity validates that a quantity is within acceptable range
func ValidateQuantity(quantity int) error {
	if quantity < MinQuantity {
		return fmt.Errorf("quantity must be at least %d", MinQuantity)
	}
	if quantity > MaxQuantity {
		return fmt.Errorf("quantity cannot exceed %d", MaxQuantity)
	}
	return nil
}

// ValidateOrderID validates that an order ID is in the correct format
func ValidateOrderID(orderID string) error {
	if orderID == "" {
		return fmt.Errorf("order ID cannot be empty")
	}
	if !orderIDPattern.MatchString(orderID) {
		return fmt.Errorf("invalid order ID format: must be in format XXX-XXXXXXX-XXXXXXX")
	}
	return nil
}

// ValidateSubscriptionID validates that a subscription ID is in the correct format
func ValidateSubscriptionID(subscriptionID string) error {
	if subscriptionID == "" {
		return fmt.Errorf("subscription ID cannot be empty")
	}
	if !subscriptionIDPattern.MatchString(subscriptionID) {
		return fmt.Errorf("invalid subscription ID format: must be in format SXX-XXXXXXX-XXXXXXX")
	}
	return nil
}

// ValidateAddressID validates that an address ID is not empty
func ValidateAddressID(addressID string) error {
	if addressID == "" {
		return fmt.Errorf("address ID cannot be empty")
	}
	return nil
}

// ValidatePaymentID validates that a payment ID is not empty
func ValidatePaymentID(paymentID string) error {
	if paymentID == "" {
		return fmt.Errorf("payment ID cannot be empty")
	}
	return nil
}

// ValidatePriceRange validates min and max price values
func ValidatePriceRange(minPrice, maxPrice float64) error {
	if minPrice < 0 {
		return fmt.Errorf("minimum price cannot be negative")
	}
	if maxPrice < 0 {
		return fmt.Errorf("maximum price cannot be negative")
	}
	if minPrice > 0 && maxPrice > 0 && minPrice > maxPrice {
		return fmt.Errorf("minimum price cannot be greater than maximum price")
	}
	if minPrice > MaxPrice {
		return fmt.Errorf("minimum price exceeds maximum allowed value of %.2f", MaxPrice)
	}
	if maxPrice > MaxPrice {
		return fmt.Errorf("maximum price exceeds maximum allowed value of %.2f", MaxPrice)
	}
	return nil
}

// ValidateReturnReason validates that a return reason is valid
func ValidateReturnReason(reason string) error {
	validReasons := map[string]bool{
		"defective":        true,
		"wrong_item":       true,
		"not_as_described": true,
		"no_longer_needed": true,
		"better_price":     true,
		"other":            true,
	}

	if reason == "" {
		return fmt.Errorf("return reason cannot be empty")
	}

	if !validReasons[reason] {
		return fmt.Errorf("invalid return reason: must be one of defective, wrong_item, not_as_described, no_longer_needed, better_price, or other")
	}

	return nil
}

// ValidateFrequencyWeeks validates subscription delivery frequency in weeks
func ValidateFrequencyWeeks(weeks int) error {
	if weeks < 1 {
		return fmt.Errorf("frequency must be at least 1 week")
	}
	if weeks > 26 {
		return fmt.Errorf("frequency cannot exceed 26 weeks")
	}
	return nil
}

// ValidateItemID validates that an item ID is not empty
func ValidateItemID(itemID string) error {
	if itemID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}
	return nil
}

// ValidateReturnID validates that a return ID is not empty
func ValidateReturnID(returnID string) error {
	if returnID == "" {
		return fmt.Errorf("return ID cannot be empty")
	}
	return nil
}

// ValidateSearchQuery validates that a search query is not empty
func ValidateSearchQuery(query string) error {
	if query == "" {
		return fmt.Errorf("search query cannot be empty")
	}
	if len(query) > 500 {
		return fmt.Errorf("search query too long: maximum 500 characters")
	}
	return nil
}
