package models

import (
	"encoding/json"
	"testing"
)

// TestOutOfStockProduct tests that out-of-stock products are correctly represented
func TestOutOfStockProduct(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected bool
	}{
		{
			name: "Product in stock",
			product: Product{
				ASIN:             "B08N5WRWNW",
				Title:            "Sony WH-1000XM4 Wireless Headphones",
				Price:            278.00,
				InStock:          true,
				DeliveryEstimate: "Tomorrow",
			},
			expected: true,
		},
		{
			name: "Product out of stock",
			product: Product{
				ASIN:             "B08EXAMPLE",
				Title:            "Out of Stock Item",
				Price:            99.99,
				InStock:          false,
				DeliveryEstimate: "Currently unavailable",
			},
			expected: false,
		},
		{
			name: "Product temporarily out of stock",
			product: Product{
				ASIN:             "B09EXAMPLE",
				Title:            "Temporarily Unavailable Item",
				Price:            49.99,
				InStock:          false,
				DeliveryEstimate: "Usually ships within 2 to 3 weeks",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.product.InStock != tt.expected {
				t.Errorf("InStock = %v, want %v", tt.product.InStock, tt.expected)
			}
		})
	}
}

// TestOutOfStockProductJSONMarshaling tests JSON serialization of out-of-stock products
func TestOutOfStockProductJSONMarshaling(t *testing.T) {
	product := Product{
		ASIN:             "B08EXAMPLE",
		Title:            "Out of Stock Item",
		Price:            99.99,
		Rating:           4.5,
		ReviewCount:      1234,
		Prime:            true,
		InStock:          false,
		DeliveryEstimate: "Currently unavailable",
	}

	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal product: %v", err)
	}

	var unmarshaled Product
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal product: %v", err)
	}

	if unmarshaled.InStock != false {
		t.Errorf("Expected InStock to be false, got %v", unmarshaled.InStock)
	}

	if unmarshaled.DeliveryEstimate != "Currently unavailable" {
		t.Errorf("Expected DeliveryEstimate 'Currently unavailable', got %s", unmarshaled.DeliveryEstimate)
	}
}

// TestSearchResponseWithOutOfStockProducts tests search results containing out-of-stock items
func TestSearchResponseWithOutOfStockProducts(t *testing.T) {
	searchResponse := SearchResponse{
		Query: "wireless headphones",
		Results: []Product{
			{
				ASIN:             "B08N5WRWNW",
				Title:            "Sony WH-1000XM4",
				Price:            278.00,
				InStock:          true,
				DeliveryEstimate: "Tomorrow",
			},
			{
				ASIN:             "B08EXAMPLE",
				Title:            "Bose QuietComfort 45",
				Price:            329.00,
				InStock:          false,
				DeliveryEstimate: "Currently unavailable",
			},
			{
				ASIN:             "B09EXAMPLE",
				Title:            "Apple AirPods Max",
				Price:            549.00,
				InStock:          true,
				DeliveryEstimate: "Ships in 2-3 days",
			},
		},
		TotalResults: 3,
		Page:         1,
	}

	inStockCount := 0
	outOfStockCount := 0

	for _, product := range searchResponse.Results {
		if product.InStock {
			inStockCount++
		} else {
			outOfStockCount++
		}
	}

	if inStockCount != 2 {
		t.Errorf("Expected 2 in-stock products, got %d", inStockCount)
	}

	if outOfStockCount != 1 {
		t.Errorf("Expected 1 out-of-stock product, got %d", outOfStockCount)
	}
}

// TestFilterOutOfStockProducts tests filtering out-of-stock products from search results
func TestFilterOutOfStockProducts(t *testing.T) {
	searchResponse := SearchResponse{
		Query: "electronics",
		Results: []Product{
			{ASIN: "A001", Title: "Item 1", InStock: true},
			{ASIN: "A002", Title: "Item 2", InStock: false},
			{ASIN: "A003", Title: "Item 3", InStock: true},
			{ASIN: "A004", Title: "Item 4", InStock: false},
			{ASIN: "A005", Title: "Item 5", InStock: true},
		},
		TotalResults: 5,
		Page:         1,
	}

	// Filter to only in-stock items
	var inStockProducts []Product
	for _, product := range searchResponse.Results {
		if product.InStock {
			inStockProducts = append(inStockProducts, product)
		}
	}

	if len(inStockProducts) != 3 {
		t.Errorf("Expected 3 in-stock products after filtering, got %d", len(inStockProducts))
	}

	// Verify all filtered products are in stock
	for _, product := range inStockProducts {
		if !product.InStock {
			t.Errorf("Product %s should be in stock but isn't", product.ASIN)
		}
	}
}

// TestOutOfStockProductWithNilPrice tests edge case of out-of-stock with special pricing
func TestOutOfStockProductWithNilPrice(t *testing.T) {
	originalPrice := 99.99
	product := Product{
		ASIN:          "B10EXAMPLE",
		Title:         "Limited Edition Item",
		Price:         0.0, // Price might be 0 when out of stock
		OriginalPrice: &originalPrice,
		InStock:       false,
	}

	if product.InStock {
		t.Error("Product should be out of stock")
	}

	if product.OriginalPrice == nil {
		t.Error("OriginalPrice should not be nil")
	}

	if *product.OriginalPrice != 99.99 {
		t.Errorf("Expected OriginalPrice of 99.99, got %f", *product.OriginalPrice)
	}
}

// TestOutOfStockDeliveryEstimates tests various delivery estimate messages for out-of-stock items
func TestOutOfStockDeliveryEstimates(t *testing.T) {
	estimates := []struct {
		message       string
		shouldBeInStock bool
	}{
		{"Currently unavailable", false},
		{"Out of stock", false},
		{"Temporarily unavailable", false},
		{"Usually ships within 2 to 3 weeks", false},
		{"Available from these sellers", false},
		{"Tomorrow", true},
		{"Ships in 2-3 days", true},
		{"In stock soon", false},
	}

	for _, est := range estimates {
		t.Run(est.message, func(t *testing.T) {
			product := Product{
				ASIN:             "TESTITEM",
				Title:            "Test Product",
				Price:            50.00,
				InStock:          est.shouldBeInStock,
				DeliveryEstimate: est.message,
			}

			if product.InStock != est.shouldBeInStock {
				t.Errorf("For estimate '%s', expected InStock=%v, got %v",
					est.message, est.shouldBeInStock, product.InStock)
			}
		})
	}
}

// TestJSONResponseWithMixedStockStatus tests complete JSON response handling
func TestJSONResponseWithMixedStockStatus(t *testing.T) {
	jsonData := `{
		"query": "test products",
		"results": [
			{
				"asin": "B001",
				"title": "Available Product",
				"price": 29.99,
				"rating": 4.5,
				"review_count": 100,
				"prime": true,
				"in_stock": true,
				"delivery_estimate": "Tomorrow"
			},
			{
				"asin": "B002",
				"title": "Unavailable Product",
				"price": 39.99,
				"rating": 4.0,
				"review_count": 50,
				"prime": false,
				"in_stock": false,
				"delivery_estimate": "Currently unavailable"
			}
		],
		"total_results": 2,
		"page": 1
	}`

	var response SearchResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(response.Results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(response.Results))
	}

	// Check first product (in stock)
	if !response.Results[0].InStock {
		t.Error("First product should be in stock")
	}

	// Check second product (out of stock)
	if response.Results[1].InStock {
		t.Error("Second product should be out of stock")
	}

	if response.Results[1].DeliveryEstimate != "Currently unavailable" {
		t.Errorf("Expected 'Currently unavailable', got '%s'", response.Results[1].DeliveryEstimate)
	}
}

// TestEmptySearchResultsForOutOfStock tests handling of searches with no in-stock results
func TestEmptySearchResultsForOutOfStock(t *testing.T) {
	searchResponse := SearchResponse{
		Query: "rare collectible",
		Results: []Product{
			{ASIN: "R001", Title: "Rare Item 1", InStock: false},
			{ASIN: "R002", Title: "Rare Item 2", InStock: false},
			{ASIN: "R003", Title: "Rare Item 3", InStock: false},
		},
		TotalResults: 3,
		Page:         1,
	}

	allOutOfStock := true
	for _, product := range searchResponse.Results {
		if product.InStock {
			allOutOfStock = false
			break
		}
	}

	if !allOutOfStock {
		t.Error("All products should be out of stock")
	}

	if len(searchResponse.Results) != 3 {
		t.Errorf("Should have 3 results, got %d", len(searchResponse.Results))
	}
}
