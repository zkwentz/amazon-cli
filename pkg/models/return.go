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
	Status    string `json:"status"` // initiated, shipped, received, refunded
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ReturnLabel represents a return shipping label
type ReturnLabel struct {
	URL          string `json:"url"`
	Carrier      string `json:"carrier"`
	Instructions string `json:"instructions"`
}

// ReturnReason constants for validation
const (
	ReasonDefective       = "defective"
	ReasonWrongItem       = "wrong_item"
	ReasonNotAsDescribed  = "not_as_described"
	ReasonNoLongerNeeded  = "no_longer_needed"
	ReasonBetterPrice     = "better_price"
	ReasonOther           = "other"
)

// ReturnStatus constants
const (
	StatusInitiated = "initiated"
	StatusShipped   = "shipped"
	StatusReceived  = "received"
	StatusRefunded  = "refunded"
)
