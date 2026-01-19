package amazon

import (
	"strings"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// Search searches for products on Amazon
func (c *Client) Search(query string, opts models.SearchOptions) (*models.SearchResponse, error) {
	// TODO: Implement actual Amazon search
	// For now, return mock data

	if opts.Page <= 0 {
		opts.Page = 1
	}

	originalPrice := 349.99

	products := []models.Product{
		{
			ASIN:             "B08N5WRWNW",
			Title:            "Sony WH-1000XM4 Wireless Premium Noise Canceling Headphones",
			Price:            278.00,
			OriginalPrice:    &originalPrice,
			Rating:           4.7,
			ReviewCount:      52431,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "Tomorrow",
		},
		{
			ASIN:             "B0BXY1234Z",
			Title:            "Apple AirPods Pro (2nd Generation)",
			Price:            189.99,
			Rating:           4.8,
			ReviewCount:      89234,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "Tomorrow",
		},
		{
			ASIN:             "B09ABC7890",
			Title:            "Bose QuietComfort 45 Bluetooth Wireless Headphones",
			Price:            249.00,
			Rating:           4.6,
			ReviewCount:      31256,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "2 days",
		},
	}

	// Filter by query (simple substring match)
	if query != "" {
		filtered := []models.Product{}
		queryLower := strings.ToLower(query)
		for _, p := range products {
			if strings.Contains(strings.ToLower(p.Title), queryLower) {
				filtered = append(filtered, p)
			}
		}
		// If no matches, return all for demo purposes
		if len(filtered) > 0 {
			products = filtered
		}
	}

	// Filter by price range
	if opts.MinPrice > 0 || opts.MaxPrice > 0 {
		filtered := []models.Product{}
		for _, p := range products {
			if opts.MinPrice > 0 && p.Price < opts.MinPrice {
				continue
			}
			if opts.MaxPrice > 0 && p.Price > opts.MaxPrice {
				continue
			}
			filtered = append(filtered, p)
		}
		products = filtered
	}

	// Filter by Prime
	if opts.PrimeOnly {
		filtered := []models.Product{}
		for _, p := range products {
			if p.Prime {
				filtered = append(filtered, p)
			}
		}
		products = filtered
	}

	return &models.SearchResponse{
		Query:        query,
		Results:      products,
		TotalResults: len(products),
		Page:         opts.Page,
	}, nil
}
