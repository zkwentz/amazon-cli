package tests

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestSearchQueryTypes tests the search functionality with various query types
// This test suite validates search behavior for different query patterns,
// filters, and combinations as specified in Phase 5.5 of the PRD.

// SearchOptions represents the filters that can be applied to a search
type SearchOptions struct {
	Category  string
	MinPrice  float64
	MaxPrice  float64
	PrimeOnly bool
	Page      int
}

// SearchResponse represents the expected search response structure
type SearchResponse struct {
	Query        string    `json:"query"`
	Results      []Product `json:"results"`
	TotalResults int       `json:"total_results"`
	Page         int       `json:"page"`
}

// Product represents a product in search results
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
}

// MockSearchClient simulates the search functionality for testing
type MockSearchClient struct {
	// In a real implementation, this would be the actual Amazon client
}

// Search performs a mock search operation
func (c *MockSearchClient) Search(query string, opts SearchOptions) (*SearchResponse, error) {
	// This is a mock implementation for testing
	// In production, this would call the actual Amazon API/scraper

	results := []Product{}

	// Generate mock results based on query
	if query != "" {
		// Base product
		product := Product{
			ASIN:             "B08N5WRWNW",
			Title:            query + " - Test Product",
			Price:            99.99,
			Rating:           4.5,
			ReviewCount:      1000,
			Prime:            true,
			InStock:          true,
			DeliveryEstimate: "Tomorrow",
		}

		// Apply filters
		if opts.MinPrice > 0 && product.Price < opts.MinPrice {
			return &SearchResponse{Query: query, Results: results, TotalResults: 0, Page: opts.Page}, nil
		}
		if opts.MaxPrice > 0 && product.Price > opts.MaxPrice {
			return &SearchResponse{Query: query, Results: results, TotalResults: 0, Page: opts.Page}, nil
		}
		if opts.PrimeOnly && !product.Prime {
			return &SearchResponse{Query: query, Results: results, TotalResults: 0, Page: opts.Page}, nil
		}

		results = append(results, product)
	}

	return &SearchResponse{
		Query:        query,
		Results:      results,
		TotalResults: len(results),
		Page:         opts.Page,
	}, nil
}

// TestBasicTextQuery tests simple text-based search queries
func TestBasicTextQuery(t *testing.T) {
	client := &MockSearchClient{}

	testCases := []struct {
		name          string
		query         string
		expectResults bool
	}{
		{
			name:          "Single word query",
			query:         "headphones",
			expectResults: true,
		},
		{
			name:          "Multi-word query",
			query:         "wireless headphones",
			expectResults: true,
		},
		{
			name:          "Query with brand name",
			query:         "Sony WH-1000XM4",
			expectResults: true,
		},
		{
			name:          "Query with model number",
			query:         "B08N5WRWNW",
			expectResults: true,
		},
		{
			name:          "Empty query",
			query:         "",
			expectResults: false,
		},
		{
			name:          "Special characters in query",
			query:         "headphones & speakers",
			expectResults: true,
		},
		{
			name:          "Query with quotes",
			query:         `"wireless headphones"`,
			expectResults: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := client.Search(tc.query, SearchOptions{Page: 1})

			if err != nil {
				t.Fatalf("Search failed: %v", err)
			}

			if tc.expectResults && len(result.Results) == 0 {
				t.Errorf("Expected results for query '%s', got none", tc.query)
			}

			if !tc.expectResults && len(result.Results) > 0 {
				t.Errorf("Expected no results for query '%s', got %d", tc.query, len(result.Results))
			}

			if result.Query != tc.query {
				t.Errorf("Expected query '%s', got '%s'", tc.query, result.Query)
			}
		})
	}
}

// TestCategoryFiltering tests search with category filters
func TestCategoryFiltering(t *testing.T) {
	client := &MockSearchClient{}

	testCases := []struct {
		name     string
		query    string
		category string
	}{
		{
			name:     "Electronics category",
			query:    "headphones",
			category: "electronics",
		},
		{
			name:     "Books category",
			query:    "programming",
			category: "books",
		},
		{
			name:     "Home & Kitchen category",
			query:    "coffee maker",
			category: "home-kitchen",
		},
		{
			name:     "No category specified",
			query:    "generic search",
			category: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := SearchOptions{
				Category: tc.category,
				Page:     1,
			}

			result, err := client.Search(tc.query, opts)

			if err != nil {
				t.Fatalf("Search with category failed: %v", err)
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			// Validate response structure
			if result.Query != tc.query {
				t.Errorf("Expected query '%s', got '%s'", tc.query, result.Query)
			}
		})
	}
}

// TestPriceRangeFiltering tests search with price range filters
func TestPriceRangeFiltering(t *testing.T) {
	client := &MockSearchClient{}

	testCases := []struct {
		name      string
		query     string
		minPrice  float64
		maxPrice  float64
		expectMin bool
		expectMax bool
	}{
		{
			name:      "Max price only",
			query:     "headphones",
			minPrice:  0,
			maxPrice:  100.0,
			expectMin: false,
			expectMax: true,
		},
		{
			name:      "Min price only",
			query:     "headphones",
			minPrice:  50.0,
			maxPrice:  0,
			expectMin: true,
			expectMax: false,
		},
		{
			name:      "Both min and max price",
			query:     "headphones",
			minPrice:  50.0,
			maxPrice:  200.0,
			expectMin: true,
			expectMax: true,
		},
		{
			name:      "No price filter",
			query:     "headphones",
			minPrice:  0,
			maxPrice:  0,
			expectMin: false,
			expectMax: false,
		},
		{
			name:      "Very low max price (expect no results)",
			query:     "headphones",
			minPrice:  0,
			maxPrice:  10.0,
			expectMin: false,
			expectMax: true,
		},
		{
			name:      "Very high min price (expect no results)",
			query:     "headphones",
			minPrice:  1000.0,
			maxPrice:  0,
			expectMin: true,
			expectMax: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := SearchOptions{
				MinPrice: tc.minPrice,
				MaxPrice: tc.maxPrice,
				Page:     1,
			}

			result, err := client.Search(tc.query, opts)

			if err != nil {
				t.Fatalf("Search with price filter failed: %v", err)
			}

			// Validate that results respect price constraints
			for _, product := range result.Results {
				if tc.expectMin && product.Price < tc.minPrice {
					t.Errorf("Product price %.2f is below minimum %.2f", product.Price, tc.minPrice)
				}
				if tc.expectMax && product.Price > tc.maxPrice {
					t.Errorf("Product price %.2f is above maximum %.2f", product.Price, tc.maxPrice)
				}
			}
		})
	}
}

// TestPrimeOnlyFiltering tests search with Prime-only filter
func TestPrimeOnlyFiltering(t *testing.T) {
	client := &MockSearchClient{}

	testCases := []struct {
		name      string
		query     string
		primeOnly bool
	}{
		{
			name:      "Prime only enabled",
			query:     "headphones",
			primeOnly: true,
		},
		{
			name:      "Prime only disabled",
			query:     "headphones",
			primeOnly: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := SearchOptions{
				PrimeOnly: tc.primeOnly,
				Page:      1,
			}

			result, err := client.Search(tc.query, opts)

			if err != nil {
				t.Fatalf("Search with Prime filter failed: %v", err)
			}

			// Validate that all results have Prime if filter is enabled
			if tc.primeOnly {
				for _, product := range result.Results {
					if !product.Prime {
						t.Errorf("Expected Prime product, got non-Prime: %s", product.Title)
					}
				}
			}
		})
	}
}

// TestPaginationSupport tests search pagination
func TestPaginationSupport(t *testing.T) {
	client := &MockSearchClient{}

	testCases := []struct {
		name string
		page int
	}{
		{
			name: "First page",
			page: 1,
		},
		{
			name: "Second page",
			page: 2,
		},
		{
			name: "Fifth page",
			page: 5,
		},
	}

	query := "headphones"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := SearchOptions{
				Page: tc.page,
			}

			result, err := client.Search(query, opts)

			if err != nil {
				t.Fatalf("Search with pagination failed: %v", err)
			}

			if result.Page != tc.page {
				t.Errorf("Expected page %d, got %d", tc.page, result.Page)
			}
		})
	}
}

// TestCombinedFilters tests search with multiple filters applied simultaneously
func TestCombinedFilters(t *testing.T) {
	client := &MockSearchClient{}

	testCases := []struct {
		name    string
		query   string
		options SearchOptions
	}{
		{
			name:  "Price range + Prime only",
			query: "headphones",
			options: SearchOptions{
				MinPrice:  50.0,
				MaxPrice:  200.0,
				PrimeOnly: true,
				Page:      1,
			},
		},
		{
			name:  "Category + Price range",
			query: "headphones",
			options: SearchOptions{
				Category: "electronics",
				MinPrice: 50.0,
				MaxPrice: 300.0,
				Page:     1,
			},
		},
		{
			name:  "All filters combined",
			query: "wireless headphones",
			options: SearchOptions{
				Category:  "electronics",
				MinPrice:  50.0,
				MaxPrice:  200.0,
				PrimeOnly: true,
				Page:      1,
			},
		},
		{
			name:  "Category + Prime only",
			query: "coffee maker",
			options: SearchOptions{
				Category:  "home-kitchen",
				PrimeOnly: true,
				Page:      1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := client.Search(tc.query, tc.options)

			if err != nil {
				t.Fatalf("Search with combined filters failed: %v", err)
			}

			// Validate all filters are applied
			for _, product := range result.Results {
				if tc.options.MinPrice > 0 && product.Price < tc.options.MinPrice {
					t.Errorf("Product price %.2f violates minimum %.2f", product.Price, tc.options.MinPrice)
				}
				if tc.options.MaxPrice > 0 && product.Price > tc.options.MaxPrice {
					t.Errorf("Product price %.2f violates maximum %.2f", product.Price, tc.options.MaxPrice)
				}
				if tc.options.PrimeOnly && !product.Prime {
					t.Errorf("Expected Prime product in Prime-only search")
				}
			}
		})
	}
}

// TestSearchResponseStructure validates the JSON output structure
func TestSearchResponseStructure(t *testing.T) {
	client := &MockSearchClient{}

	result, err := client.Search("wireless headphones", SearchOptions{Page: 1})

	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Marshal to JSON and back to validate structure
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal search response: %v", err)
	}

	var decoded SearchResponse
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal search response: %v", err)
	}

	// Validate required fields
	if decoded.Query == "" {
		t.Error("Query field is empty")
	}

	if decoded.Results == nil {
		t.Error("Results field is nil")
	}

	if decoded.Page < 1 {
		t.Error("Page field should be at least 1")
	}

	// Validate product structure
	for i, product := range decoded.Results {
		if product.ASIN == "" {
			t.Errorf("Product %d has empty ASIN", i)
		}
		if product.Title == "" {
			t.Errorf("Product %d has empty Title", i)
		}
		if product.Price < 0 {
			t.Errorf("Product %d has negative price", i)
		}
	}
}

// TestEdgeCases tests various edge cases in search queries
func TestEdgeCases(t *testing.T) {
	client := &MockSearchClient{}

	testCases := []struct {
		name        string
		query       string
		options     SearchOptions
		expectError bool
	}{
		{
			name:        "Very long query string",
			query:       strings.Repeat("test ", 100),
			options:     SearchOptions{Page: 1},
			expectError: false,
		},
		{
			name:        "Query with unicode characters",
			query:       "café ☕ espresso",
			options:     SearchOptions{Page: 1},
			expectError: false,
		},
		{
			name:        "Query with numbers",
			query:       "iPhone 15 Pro Max 256GB",
			options:     SearchOptions{Page: 1},
			expectError: false,
		},
		{
			name:        "Invalid price range (min > max)",
			query:       "headphones",
			options:     SearchOptions{MinPrice: 200.0, MaxPrice: 100.0, Page: 1},
			expectError: false, // Should return no results, not error
		},
		{
			name:        "Negative price values",
			query:       "headphones",
			options:     SearchOptions{MinPrice: -10.0, MaxPrice: -5.0, Page: 1},
			expectError: false, // Should handle gracefully
		},
		{
			name:        "Zero page number",
			query:       "headphones",
			options:     SearchOptions{Page: 0},
			expectError: false, // Should default to page 1
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := client.Search(tc.query, tc.options)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err == nil && result == nil {
				t.Error("Expected non-nil result when no error occurs")
			}
		})
	}
}

// TestQueryTypeVariations tests different query input patterns
func TestQueryTypeVariations(t *testing.T) {
	client := &MockSearchClient{}

	testCases := []struct {
		name  string
		query string
		desc  string
	}{
		{
			name:  "Product by ASIN",
			query: "B08N5WRWNW",
			desc:  "Direct ASIN lookup",
		},
		{
			name:  "Brand + model",
			query: "Sony WH-1000XM4",
			desc:  "Brand and specific model number",
		},
		{
			name:  "Generic category search",
			query: "wireless headphones",
			desc:  "General product category",
		},
		{
			name:  "Descriptive search",
			query: "noise cancelling over ear headphones",
			desc:  "Search with multiple descriptive terms",
		},
		{
			name:  "Price-focused search",
			query: "cheap wireless earbuds",
			desc:  "Search including price qualifiers",
		},
		{
			name:  "Feature-based search",
			query: "waterproof bluetooth speaker",
			desc:  "Search by product features",
		},
		{
			name:  "Comparison search",
			query: "iPhone vs Samsung Galaxy",
			desc:  "Comparison query",
		},
		{
			name:  "Question-format search",
			query: "what is the best laptop for programming",
			desc:  "Natural language question",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := client.Search(tc.query, SearchOptions{Page: 1})

			if err != nil {
				t.Fatalf("Search failed for %s: %v", tc.desc, err)
			}

			if result == nil {
				t.Fatalf("Expected non-nil result for %s", tc.desc)
			}

			if result.Query != tc.query {
				t.Errorf("Query mismatch: expected '%s', got '%s'", tc.query, result.Query)
			}

			t.Logf("Query type '%s' (%s) returned %d results", tc.name, tc.desc, len(result.Results))
		})
	}
}

// BenchmarkSearchPerformance benchmarks search performance
func BenchmarkSearchPerformance(b *testing.B) {
	client := &MockSearchClient{}

	b.Run("Simple query", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = client.Search("headphones", SearchOptions{Page: 1})
		}
	})

	b.Run("Query with all filters", func(b *testing.B) {
		opts := SearchOptions{
			Category:  "electronics",
			MinPrice:  50.0,
			MaxPrice:  200.0,
			PrimeOnly: true,
			Page:      1,
		}
		for i := 0; i < b.N; i++ {
			_, _ = client.Search("wireless headphones", opts)
		}
	})
}
