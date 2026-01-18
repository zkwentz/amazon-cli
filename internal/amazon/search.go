package amazon

import (
	"fmt"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// Search performs a product search with the given query and options
// If opts.PrimeOnly is true, only Prime-eligible products are returned
func Search(query string, opts models.SearchOptions) (*models.SearchResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// In a real implementation, this would make HTTP requests to Amazon
	// For now, this is a mock implementation to demonstrate the filtering logic

	// Mock products for demonstration
	allProducts := []models.Product{
		{
			ASIN:             "B08N5WRWNW",
			Title:            "Sony WH-1000XM4 Wireless Headphones",
			Price:            278.00,
			OriginalPrice:    ptrFloat64(349.99),
			Rating:           4.7,
			ReviewCount:      52431,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "Tomorrow",
		},
		{
			ASIN:             "B07Q5Y7DJ1",
			Title:            "Budget Wireless Headphones",
			Price:            45.99,
			Rating:           4.2,
			ReviewCount:      1250,
			Prime:            false,
			InStock:          true,
			DeliveryEstimate: "5-7 days",
		},
		{
			ASIN:             "B09JQMJHXY",
			Title:            "Premium Audio Headphones",
			Price:            199.99,
			Rating:           4.8,
			ReviewCount:      8934,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "Tomorrow",
		},
		{
			ASIN:             "B08XYZABC1",
			Title:            "Standard Headphones",
			Price:            89.99,
			Rating:           4.0,
			ReviewCount:      523,
			Prime:            false,
			InStock:          true,
			DeliveryEstimate: "3-5 days",
		},
	}

	// Apply filters
	filteredProducts := applyFilters(allProducts, opts)

	response := &models.SearchResponse{
		Query:        query,
		Results:      filteredProducts,
		TotalResults: len(filteredProducts),
		Page:         opts.Page,
	}

	if response.Page == 0 {
		response.Page = 1
	}

	return response, nil
}

// applyFilters applies all search filters to the products
func applyFilters(products []models.Product, opts models.SearchOptions) []models.Product {
	filtered := make([]models.Product, 0, len(products))

	for _, product := range products {
		// Prime-only filter
		if opts.PrimeOnly && !product.Prime {
			continue
		}

		// Price range filters
		if opts.MinPrice > 0 && product.Price < opts.MinPrice {
			continue
		}
		if opts.MaxPrice > 0 && product.Price > opts.MaxPrice {
			continue
		}

		filtered = append(filtered, product)
	}

	return filtered
}

// ptrFloat64 is a helper to create a pointer to a float64
func ptrFloat64(f float64) *float64 {
	return &f
}
