package models

import "time"

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

// ReturnOption represents a return method option
type ReturnOption struct {
	Method          string  `json:"method"`
	Label           string  `json:"label"`
	DropoffLocation string  `json:"dropoff_location,omitempty"`
	Fee             float64 `json:"fee,omitempty"`
}

// Return represents a return request
type Return struct {
	ReturnID  string    `json:"return_id"`
	OrderID   string    `json:"order_id"`
	ItemID    string    `json:"item_id"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// ReturnLabel represents a return shipping label
type ReturnLabel struct {
	URL          string `json:"url"`
	Carrier      string `json:"carrier"`
	Instructions string `json:"instructions"`
}

// ValidReturnReasons contains all valid return reason codes
var ValidReturnReasons = map[string]bool{
	"defective":         true,
	"wrong_item":        true,
	"not_as_described":  true,
	"no_longer_needed":  true,
	"better_price":      true,
	"other":             true,
}

// IsValidReturnReason checks if a reason code is valid
func IsValidReturnReason(reason string) bool {
	return ValidReturnReasons[reason]
}
