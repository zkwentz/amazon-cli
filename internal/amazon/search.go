package amazon

import (
	"fmt"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// Search performs a product search with the given query and options
func Search(query string, opts models.SearchOptions) (*models.SearchResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// This is a placeholder implementation that would normally fetch from Amazon
	// For now, it returns mock data filtered by price range
	mockProducts := []models.Product{
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
			ASIN:             "B07PXGQC1Q",
			Title:            "Bose QuietComfort 35 II",
			Price:            299.00,
			Rating:           4.5,
			ReviewCount:      35000,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "Tomorrow",
		},
		{
			ASIN:             "B0863TXGM3",
			Title:            "JBL Tune 510BT Wireless Headphones",
			Price:            29.95,
			Rating:           4.4,
			ReviewCount:      15000,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "2 days",
		},
		{
			ASIN:             "B09JB3PPMK",
			Title:            "Anker Soundcore Q30 Wireless Headphones",
			Price:            79.99,
			Rating:           4.6,
			ReviewCount:      25000,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "Tomorrow",
		},
		{
			ASIN:             "B08PZHYWJS",
			Title:            "Beats Studio Buds",
			Price:            149.95,
			Rating:           4.4,
			ReviewCount:      18000,
			Prime:            false,
			InStock:          true,
			DeliveryEstimate: "3 days",
		},
	}

	// Apply price range filtering
	filteredProducts := filterByPriceRange(mockProducts, opts.MinPrice, opts.MaxPrice)

	// Apply Prime-only filtering if requested
	if opts.PrimeOnly {
		filteredProducts = filterByPrime(filteredProducts)
	}

	return &models.SearchResponse{
		Query:        query,
		Results:      filteredProducts,
		TotalResults: len(filteredProducts),
		Page:         opts.Page,
	}, nil
}

// filterByPriceRange filters products based on min and max price
func filterByPriceRange(products []models.Product, minPrice, maxPrice float64) []models.Product {
	if minPrice <= 0 && maxPrice <= 0 {
		return products
	}

	filtered := make([]models.Product, 0)
	for _, product := range products {
		// If only minPrice is set
		if minPrice > 0 && maxPrice <= 0 {
			if product.Price >= minPrice {
				filtered = append(filtered, product)
			}
		} else if minPrice <= 0 && maxPrice > 0 {
			// If only maxPrice is set
			if product.Price <= maxPrice {
				filtered = append(filtered, product)
			}
		} else {
			// Both minPrice and maxPrice are set
			if product.Price >= minPrice && product.Price <= maxPrice {
				filtered = append(filtered, product)
			}
		}
	}
	return filtered
}

// filterByPrime filters products to only include Prime-eligible items
func filterByPrime(products []models.Product) []models.Product {
	filtered := make([]models.Product, 0)
	for _, product := range products {
		if product.Prime {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

// ptrFloat64 is a helper function to create a pointer to a float64
func ptrFloat64(f float64) *float64 {
	return &f
}
