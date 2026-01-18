package models

import (
	"fmt"
	"regexp"
)

// Subscription ID format: S01-1234567-8901234
// Pattern: S + 2 digits + hyphen + 7 digits + hyphen + 7 digits
var subscriptionIDPattern = regexp.MustCompile(`^S\d{2}-\d{7}-\d{7}$`)

// ValidateSubscriptionID validates that a subscription ID matches the expected format
// Expected format: S01-1234567-8901234
func ValidateSubscriptionID(subscriptionID string) error {
	if subscriptionID == "" {
		return fmt.Errorf("subscription ID cannot be empty")
	}

	if !subscriptionIDPattern.MatchString(subscriptionID) {
		return fmt.Errorf("invalid subscription ID format: expected format S01-1234567-8901234, got %s", subscriptionID)
	}

	return nil
}
