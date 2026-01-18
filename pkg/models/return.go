package models

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
	ReturnID  string `json:"return_id"`
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
}

// ReturnLabel represents a return shipping label
type ReturnLabel struct {
	URL          string `json:"url"`
	Carrier      string `json:"carrier"`
	Instructions string `json:"instructions"`
}

// ReturnCreateRequest represents the request to create a return
type ReturnCreateRequest struct {
	OrderID string `json:"order_id"`
	ItemID  string `json:"item_id"`
	Reason  string `json:"reason"`
	Confirm bool   `json:"confirm"`
}

// ReturnCreateResponse represents the response from creating a return
type ReturnCreateResponse struct {
	DryRun      bool    `json:"dry_run,omitempty"`
	WouldReturn *Return `json:"would_return,omitempty"`
	Message     string  `json:"message,omitempty"`
	Return      *Return `json:"return,omitempty"`
}

// ValidReturnReasons contains all valid reason codes
var ValidReturnReasons = map[string]string{
	"defective":        "Item is defective or doesn't work",
	"wrong_item":       "Received wrong item",
	"not_as_described": "Item not as described",
	"no_longer_needed": "No longer needed",
	"better_price":     "Found better price elsewhere",
	"other":            "Other reason",
}

// IsValidReturnReason checks if a reason code is valid
func IsValidReturnReason(reason string) bool {
	_, ok := ValidReturnReasons[reason]
	return ok
}
