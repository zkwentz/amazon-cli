package models

// Product represents an Amazon product
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

// SearchResponse represents the response from a product search
type SearchResponse struct {
	Query        string    `json:"query"`
	Results      []Product `json:"results"`
	TotalResults int       `json:"total_results"`
	Page         int       `json:"page"`
}

// SearchOptions contains options for product search
type SearchOptions struct {
	Category string
	MinPrice float64
	MaxPrice float64
	PrimeOnly bool
	Page     int
}

// Review represents a product review
type Review struct {
	Rating   int    `json:"rating"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Author   string `json:"author"`
	Date     string `json:"date"`
	Verified bool   `json:"verified"`
}

// ReviewsResponse represents the response for product reviews
type ReviewsResponse struct {
	ASIN          string   `json:"asin"`
	AverageRating float64  `json:"average_rating"`
	TotalReviews  int      `json:"total_reviews"`
	Reviews       []Review `json:"reviews"`
}
