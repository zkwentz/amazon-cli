package models

import (
	"fmt"
	"strings"
)

// ReturnReasonCode represents a valid reason for returning an item
type ReturnReasonCode string

const (
	ReasonDefective       ReturnReasonCode = "defective"
	ReasonWrongItem       ReturnReasonCode = "wrong_item"
	ReasonNotAsDescribed  ReturnReasonCode = "not_as_described"
	ReasonNoLongerNeeded  ReturnReasonCode = "no_longer_needed"
	ReasonBetterPrice     ReturnReasonCode = "better_price"
	ReasonOther           ReturnReasonCode = "other"
)

// AllReturnReasons contains all valid return reason codes
var AllReturnReasons = []ReturnReasonCode{
	ReasonDefective,
	ReasonWrongItem,
	ReasonNotAsDescribed,
	ReasonNoLongerNeeded,
	ReasonBetterPrice,
	ReasonOther,
}

// ReturnableItem represents an item that can be returned
type ReturnableItem struct {
	OrderID      string  `json:"order_id"`
	ItemID       string  `json:"item_id"`
	ASIN         string  `json:"asin"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	PurchaseDate string  `json:"purchase_date"`
	ReturnWindow string  `json:"return_window"`
}

// ReturnOption represents a method for returning an item
type ReturnOption struct {
	Method          string  `json:"method"`
	Label           string  `json:"label"`
	DropoffLocation string  `json:"dropoff_location,omitempty"`
	Fee             float64 `json:"fee"`
}

// Return represents an initiated return request
type Return struct {
	ReturnID  string           `json:"return_id"`
	OrderID   string           `json:"order_id"`
	ItemID    string           `json:"item_id"`
	Status    string           `json:"status"`
	Reason    ReturnReasonCode `json:"reason"`
	CreatedAt string           `json:"created_at"`
}

// ReturnLabel represents a return shipping label
type ReturnLabel struct {
	URL          string `json:"url"`
	Carrier      string `json:"carrier"`
	Instructions string `json:"instructions"`
}

// ValidateReturnReason checks if a return reason code is valid
func ValidateReturnReason(reason string) error {
	// Convert to lowercase for case-insensitive comparison
	normalizedReason := strings.ToLower(strings.TrimSpace(reason))

	for _, validReason := range AllReturnReasons {
		if normalizedReason == string(validReason) {
			return nil
		}
	}

	return fmt.Errorf("invalid return reason '%s': must be one of [defective, wrong_item, not_as_described, no_longer_needed, better_price, other]", reason)
}

// IsValidReturnReason checks if a return reason code is valid (returns bool)
func IsValidReturnReason(reason string) bool {
	return ValidateReturnReason(reason) == nil
}

// GetReturnReasonDescription returns a human-readable description for a reason code
func GetReturnReasonDescription(reason ReturnReasonCode) string {
	descriptions := map[ReturnReasonCode]string{
		ReasonDefective:      "Item is defective or doesn't work",
		ReasonWrongItem:      "Received wrong item",
		ReasonNotAsDescribed: "Item not as described",
		ReasonNoLongerNeeded: "No longer needed",
		ReasonBetterPrice:    "Found better price elsewhere",
		ReasonOther:          "Other reason",
	}

	if desc, ok := descriptions[reason]; ok {
		return desc
	}
	return "Unknown reason"
}

// NormalizeReturnReason converts a string to a valid ReturnReasonCode
func NormalizeReturnReason(reason string) (ReturnReasonCode, error) {
	if err := ValidateReturnReason(reason); err != nil {
		return "", err
	}

	return ReturnReasonCode(strings.ToLower(strings.TrimSpace(reason))), nil
}
