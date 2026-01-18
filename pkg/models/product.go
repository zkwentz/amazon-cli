package models

// Product represents an Amazon product with all its details
type Product struct {
	ASIN             string    `json:"asin"`
	Title            string    `json:"title"`
	Price            float64   `json:"price"`
	OriginalPrice    *float64  `json:"original_price,omitempty"`
	Rating           float64   `json:"rating"`
	ReviewCount      int       `json:"review_count"`
	Prime            bool      `json:"prime"`
	InStock          bool      `json:"in_stock"`
	DeliveryEstimate string    `json:"delivery_estimate"`
	Description      string    `json:"description,omitempty"`
	Features         []string  `json:"features,omitempty"`
	Images           []string  `json:"images,omitempty"`
}

// SearchOptions contains all filtering options for product search
type SearchOptions struct {
	Category  string
	MinPrice  float64
	MaxPrice  float64
	PrimeOnly bool
	Page      int
}

// SearchResponse represents the result of a product search
type SearchResponse struct {
	Query        string    `json:"query"`
	Results      []Product `json:"results"`
	TotalResults int       `json:"total_results"`
	Page         int       `json:"page"`
}

// Review represents a product review
type Review struct {
	Rating   float64 `json:"rating"`
	Title    string  `json:"title"`
	Body     string  `json:"body"`
	Author   string  `json:"author"`
	Date     string  `json:"date"`
	Verified bool    `json:"verified"`
}

// ReviewsResponse represents a collection of reviews for a product
type ReviewsResponse struct {
	ASIN          string   `json:"asin"`
	Reviews       []Review `json:"reviews"`
	AverageRating float64  `json:"average_rating"`
	TotalReviews  int      `json:"total_reviews"`
}
