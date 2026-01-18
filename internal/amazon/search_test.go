package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		opts          models.SearchOptions
		expectError   bool
		expectedCount int
		description   string
	}{
		{
			name:          "Empty query should return error",
			query:         "",
			opts:          models.SearchOptions{},
			expectError:   true,
			expectedCount: 0,
			description:   "Search with empty query string should fail",
		},
		{
			name:          "Search without price filter returns all products",
			query:         "headphones",
			opts:          models.SearchOptions{},
			expectError:   false,
			expectedCount: 5,
			description:   "Search without price constraints should return all 5 mock products",
		},
		{
			name:  "Search with minPrice only",
			query: "headphones",
			opts: models.SearchOptions{
				MinPrice: 100.0,
			},
			expectError:   false,
			expectedCount: 3,
			description:   "Should return products >= $100 (Sony $278, Bose $299, Beats $149.95)",
		},
		{
			name:  "Search with maxPrice only",
			query: "headphones",
			opts: models.SearchOptions{
				MaxPrice: 100.0,
			},
			expectError:   false,
			expectedCount: 2,
			description:   "Should return products <= $100 (JBL $29.95, Anker $79.99)",
		},
		{
			name:  "Search with both minPrice and maxPrice",
			query: "headphones",
			opts: models.SearchOptions{
				MinPrice: 50.0,
				MaxPrice: 200.0,
			},
			expectError:   false,
			expectedCount: 2,
			description:   "Should return products between $50-$200 (Anker $79.99, Beats $149.95)",
		},
		{
			name:  "Search with narrow price range",
			query: "headphones",
			opts: models.SearchOptions{
				MinPrice: 70.0,
				MaxPrice: 80.0,
			},
			expectError:   false,
			expectedCount: 1,
			description:   "Should return only Anker $79.99",
		},
		{
			name:  "Search with price range excluding all products",
			query: "headphones",
			opts: models.SearchOptions{
				MinPrice: 400.0,
				MaxPrice: 500.0,
			},
			expectError:   false,
			expectedCount: 0,
			description:   "Should return no products when range excludes all",
		},
		{
			name:  "Search with exact price match",
			query: "headphones",
			opts: models.SearchOptions{
				MinPrice: 29.95,
				MaxPrice: 29.95,
			},
			expectError:   false,
			expectedCount: 1,
			description:   "Should return JBL at exact price $29.95",
		},
		{
			name:  "Search with Prime filter only",
			query: "headphones",
			opts: models.SearchOptions{
				PrimeOnly: true,
			},
			expectError:   false,
			expectedCount: 4,
			description:   "Should return only Prime-eligible products (excludes Beats)",
		},
		{
			name:  "Search with price range and Prime filter",
			query: "headphones",
			opts: models.SearchOptions{
				MinPrice:  50.0,
				MaxPrice:  200.0,
				PrimeOnly: true,
			},
			expectError:   false,
			expectedCount: 1,
			description:   "Should return only Anker $79.99 (Beats excluded by Prime filter)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Search(tt.query, tt.opts)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Fatal("Expected result but got nil")
			}

			if len(result.Results) != tt.expectedCount {
				t.Errorf("Expected %d products, got %d. %s", tt.expectedCount, len(result.Results), tt.description)
				for _, p := range result.Results {
					t.Logf("  - %s: $%.2f (Prime: %v)", p.Title, p.Price, p.Prime)
				}
			}

			// Verify all returned products match the price criteria
			if tt.opts.MinPrice > 0 || tt.opts.MaxPrice > 0 {
				for _, product := range result.Results {
					if tt.opts.MinPrice > 0 && product.Price < tt.opts.MinPrice {
						t.Errorf("Product %s ($%.2f) is below minimum price $%.2f",
							product.Title, product.Price, tt.opts.MinPrice)
					}
					if tt.opts.MaxPrice > 0 && product.Price > tt.opts.MaxPrice {
						t.Errorf("Product %s ($%.2f) is above maximum price $%.2f",
							product.Title, product.Price, tt.opts.MaxPrice)
					}
				}
			}

			// Verify Prime filter if enabled
			if tt.opts.PrimeOnly {
				for _, product := range result.Results {
					if !product.Prime {
						t.Errorf("Product %s should be Prime-eligible but isn't", product.Title)
					}
				}
			}

			// Verify query is set in response
			if result.Query != tt.query {
				t.Errorf("Expected query %q, got %q", tt.query, result.Query)
			}
		})
	}
}

func TestFilterByPriceRange(t *testing.T) {
	products := []models.Product{
		{ASIN: "A1", Title: "Product A", Price: 10.0},
		{ASIN: "A2", Title: "Product B", Price: 50.0},
		{ASIN: "A3", Title: "Product C", Price: 100.0},
		{ASIN: "A4", Title: "Product D", Price: 150.0},
		{ASIN: "A5", Title: "Product E", Price: 200.0},
	}

	tests := []struct {
		name          string
		minPrice      float64
		maxPrice      float64
		expectedASINs []string
	}{
		{
			name:          "No filters returns all",
			minPrice:      0,
			maxPrice:      0,
			expectedASINs: []string{"A1", "A2", "A3", "A4", "A5"},
		},
		{
			name:          "Min price filter only",
			minPrice:      100.0,
			maxPrice:      0,
			expectedASINs: []string{"A3", "A4", "A5"},
		},
		{
			name:          "Max price filter only",
			minPrice:      0,
			maxPrice:      100.0,
			expectedASINs: []string{"A1", "A2", "A3"},
		},
		{
			name:          "Both min and max filters",
			minPrice:      50.0,
			maxPrice:      150.0,
			expectedASINs: []string{"A2", "A3", "A4"},
		},
		{
			name:          "Range with no matches",
			minPrice:      250.0,
			maxPrice:      300.0,
			expectedASINs: []string{},
		},
		{
			name:          "Exact price match",
			minPrice:      100.0,
			maxPrice:      100.0,
			expectedASINs: []string{"A3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByPriceRange(products, tt.minPrice, tt.maxPrice)

			if len(result) != len(tt.expectedASINs) {
				t.Errorf("Expected %d products, got %d", len(tt.expectedASINs), len(result))
			}

			resultASINs := make(map[string]bool)
			for _, p := range result {
				resultASINs[p.ASIN] = true
			}

			for _, expectedASIN := range tt.expectedASINs {
				if !resultASINs[expectedASIN] {
					t.Errorf("Expected ASIN %s not found in results", expectedASIN)
				}
			}
		})
	}
}

func TestFilterByPrime(t *testing.T) {
	products := []models.Product{
		{ASIN: "A1", Title: "Prime Product 1", Prime: true},
		{ASIN: "A2", Title: "Non-Prime Product", Prime: false},
		{ASIN: "A3", Title: "Prime Product 2", Prime: true},
		{ASIN: "A4", Title: "Non-Prime Product 2", Prime: false},
		{ASIN: "A5", Title: "Prime Product 3", Prime: true},
	}

	result := filterByPrime(products)

	if len(result) != 3 {
		t.Errorf("Expected 3 Prime products, got %d", len(result))
	}

	for _, product := range result {
		if !product.Prime {
			t.Errorf("Product %s should be Prime-eligible", product.Title)
		}
	}
}

func TestPriceRangeEdgeCases(t *testing.T) {
	products := []models.Product{
		{ASIN: "A1", Title: "Product A", Price: 0.01},
		{ASIN: "A2", Title: "Product B", Price: 0.99},
		{ASIN: "A3", Title: "Product C", Price: 1.00},
		{ASIN: "A4", Title: "Product D", Price: 9999.99},
	}

	tests := []struct {
		name          string
		minPrice      float64
		maxPrice      float64
		expectedCount int
		description   string
	}{
		{
			name:          "Very low price range",
			minPrice:      0.01,
			maxPrice:      0.99,
			expectedCount: 2,
			description:   "Should handle cents correctly",
		},
		{
			name:          "Very high max price",
			minPrice:      0,
			maxPrice:      10000.0,
			expectedCount: 4,
			description:   "Should include all products below high max",
		},
		{
			name:          "Min equals product price",
			minPrice:      1.00,
			maxPrice:      0,
			expectedCount: 2,
			description:   "Should include product at exact min price",
		},
		{
			name:          "Max equals product price",
			minPrice:      0,
			maxPrice:      1.00,
			expectedCount: 3,
			description:   "Should include product at exact max price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByPriceRange(products, tt.minPrice, tt.maxPrice)
			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d products, got %d. %s", tt.expectedCount, len(result), tt.description)
			}
		})
	}
}
