package models

// ReturnableItem represents an item that is eligible for return
type ReturnableItem struct {
	OrderID      string  `json:"order_id"`
	ItemID       string  `json:"item_id"`
	ASIN         string  `json:"asin"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	PurchaseDate string  `json:"purchase_date"`
	ReturnWindow string  `json:"return_window"`
}

// ReturnOption represents a return method option for an item
type ReturnOption struct {
	Method          string  `json:"method"`
	Label           string  `json:"label"`
	DropoffLocation string  `json:"dropoff_location,omitempty"`
	Fee             float64 `json:"fee"`
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
