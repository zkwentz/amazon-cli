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
	Returnable   bool    `json:"returnable"`
}

// ReturnOption represents a method for returning an item
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

// ReturnLabel represents a shipping label for a return
type ReturnLabel struct {
	URL          string `json:"url"`
	Carrier      string `json:"carrier"`
	Instructions string `json:"instructions"`
}

// IsReturnWindowExpired checks if the return window has expired
func (r *ReturnableItem) IsReturnWindowExpired() bool {
	if r.ReturnWindow == "" {
		return true
	}

	expiryDate, err := time.Parse("2006-01-02", r.ReturnWindow)
	if err != nil {
		return true
	}

	return time.Now().After(expiryDate)
}

// IsReturnable checks if an item is eligible for return
func (r *ReturnableItem) IsReturnable() bool {
	return r.Returnable && !r.IsReturnWindowExpired()
}
