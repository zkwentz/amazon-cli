package models

// CartItem represents an item in the shopping cart
type CartItem struct {
	ASIN     string  `json:"asin"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
	Prime    bool    `json:"prime"`
	InStock  bool    `json:"in_stock"`
}

// Cart represents the shopping cart with all items and totals
type Cart struct {
	Items        []CartItem `json:"items"`
	Subtotal     float64    `json:"subtotal"`
	EstimatedTax float64    `json:"estimated_tax"`
	Total        float64    `json:"total"`
	ItemCount    int        `json:"item_count"`
}
