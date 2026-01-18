package models

// Product represents a product with all its details
type Product struct {
	ASIN             string   `json:"asin"`
	Title            string   `json:"title"`
	Price            float64  `json:"price"`
	OriginalPrice    *float64 `json:"original_price,omitempty"`
	Rating           float64  `json:"rating"`
	ReviewCount      int      `json:"review_count"`
	Prime            bool     `json:"prime"`
	InStock          bool     `json:"in_stock"`
	DeliveryEstimate string   `json:"delivery_estimate"`
	Description      string   `json:"description,omitempty"`
	Features         []string `json:"features,omitempty"`
	Images           []string `json:"images,omitempty"`
}

// BuyPreview represents a preview of a purchase before completion
type BuyPreview struct {
	DryRun   bool     `json:"dry_run"`
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
	Total    float64  `json:"total"`
	Message  string   `json:"message,omitempty"`
}
