package amazon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestSearch_EmptyQuery(t *testing.T) {
	client := NewClient()

	_, err := client.Search("", models.SearchOptions{})
	if err == nil {
		t.Error("Expected error for empty query, got nil")
	}

	if err.Error() != "search query cannot be empty" {
		t.Errorf("Expected empty query error, got: %v", err)
	}
}

func TestSearch_BasicQuery(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if r.URL.Query().Get("k") != "wireless headphones" {
			t.Errorf("Expected query 'wireless headphones', got '%s'", r.URL.Query().Get("k"))
		}

		// Return mock HTML response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockSearchHTML))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	response, err := client.Search("wireless headphones", models.SearchOptions{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if response.Query != "wireless headphones" {
		t.Errorf("Expected query 'wireless headphones', got '%s'", response.Query)
	}

	if response.Page != 1 {
		t.Errorf("Expected page 1, got %d", response.Page)
	}
}

func TestSearch_WithCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("i") != "electronics" {
			t.Errorf("Expected category 'electronics', got '%s'", r.URL.Query().Get("i"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockSearchHTML))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	opts := models.SearchOptions{
		Category: "electronics",
	}

	_, err := client.Search("test", opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSearch_WithPriceRange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("low-price") != "10.00" {
			t.Errorf("Expected low-price '10.00', got '%s'", r.URL.Query().Get("low-price"))
		}

		if r.URL.Query().Get("high-price") != "100.00" {
			t.Errorf("Expected high-price '100.00', got '%s'", r.URL.Query().Get("high-price"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockSearchHTML))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	opts := models.SearchOptions{
		MinPrice: 10.00,
		MaxPrice: 100.00,
	}

	_, err := client.Search("test", opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSearch_WithPrimeOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("prime") != "true" {
			t.Errorf("Expected prime 'true', got '%s'", r.URL.Query().Get("prime"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockSearchHTML))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	opts := models.SearchOptions{
		PrimeOnly: true,
	}

	_, err := client.Search("test", opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestSearch_WithPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "3" {
			t.Errorf("Expected page '3', got '%s'", r.URL.Query().Get("page"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockSearchHTML))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	opts := models.SearchOptions{
		Page: 3,
	}

	response, err := client.Search("test", opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.Page != 3 {
		t.Errorf("Expected page 3, got %d", response.Page)
	}
}

func TestSearch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	_, err := client.Search("test", models.SearchOptions{})
	if err == nil {
		t.Error("Expected error for HTTP 500, got nil")
	}
}

func TestExtractTotalResults(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name     string
		html     string
		expected int
	}{
		{
			name:     "results with commas",
			html:     "1-48 of over 10,000 results",
			expected: 10000,
		},
		{
			name:     "results without commas",
			html:     "1-48 of 250 results",
			expected: 250,
		},
		{
			name:     "no results found",
			html:     "No results found",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.extractTotalResults(tt.html)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestExtractProducts(t *testing.T) {
	client := NewClient()

	products := client.extractProducts(mockSearchHTML)

	if len(products) == 0 {
		t.Error("Expected at least one product to be extracted")
	}

	// Check first product
	if len(products) > 0 {
		product := products[0]

		if product.ASIN != "B08N5WRWNW" {
			t.Errorf("Expected ASIN 'B08N5WRWNW', got '%s'", product.ASIN)
		}

		if product.Title == "" {
			t.Error("Expected non-empty title")
		}

		if product.Price != 278.00 {
			t.Errorf("Expected price 278.00, got %.2f", product.Price)
		}

		if product.Rating != 4.7 {
			t.Errorf("Expected rating 4.7, got %.1f", product.Rating)
		}

		if !product.Prime {
			t.Error("Expected Prime to be true")
		}
	}
}

func TestExtractPrice(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name     string
		html     string
		expected float64
	}{
		{
			name:     "price with dollar sign",
			html:     "$29.99",
			expected: 29.99,
		},
		{
			name:     "price in span",
			html:     "<span>149.99</span>",
			expected: 149.99,
		},
		{
			name:     "whole number price",
			html:     "$50",
			expected: 50.0,
		},
		{
			name:     "no price found",
			html:     "no price here",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.extractPrice(tt.html)
			if result != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestExtractRating(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name     string
		html     string
		expected float64
	}{
		{
			name:     "standard rating",
			html:     "4.7 out of 5",
			expected: 4.7,
		},
		{
			name:     "perfect rating",
			html:     "5.0 out of 5 stars",
			expected: 5.0,
		},
		{
			name:     "no rating",
			html:     "no rating here",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.extractRating(tt.html)
			if result != tt.expected {
				t.Errorf("Expected %.1f, got %.1f", tt.expected, result)
			}
		})
	}
}

func TestCheckPrime(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "has prime badge",
			html:     "<span class='prime-badge'>Prime</span>",
			expected: true,
		},
		{
			name:     "has FREE delivery",
			html:     "FREE delivery tomorrow",
			expected: true,
		},
		{
			name:     "no prime",
			html:     "Standard shipping available",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.checkPrime(tt.html)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBuildSearchURL(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name     string
		query    string
		opts     models.SearchOptions
		contains []string
	}{
		{
			name:  "basic query",
			query: "laptop",
			opts:  models.SearchOptions{},
			contains: []string{
				"/s?",
				"k=laptop",
			},
		},
		{
			name:  "with all options",
			query: "headphones",
			opts: models.SearchOptions{
				Category:  "electronics",
				MinPrice:  50.00,
				MaxPrice:  200.00,
				PrimeOnly: true,
				Page:      2,
			},
			contains: []string{
				"k=headphones",
				"i=electronics",
				"low-price=50.00",
				"high-price=200.00",
				"prime=true",
				"page=2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := client.buildSearchURL(tt.query, tt.opts)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			for _, substr := range tt.contains {
				if !contains(url, substr) {
					t.Errorf("Expected URL to contain '%s', got: %s", substr, url)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Mock HTML response for testing
const mockSearchHTML = `
<!DOCTYPE html>
<html>
<body>
<div class="s-result-list">
	<div data-asin="B08N5WRWNW" class="s-result-item">
		<h2>
			<span>Sony WH-1000XM4 Wireless Headphones</span>
		</h2>
		<span class="a-price">
			<span>$278.00</span>
		</span>
		<span class="a-icon-alt">4.7 out of 5 stars</span>
		<span>52,431 ratings</span>
		<i class="a-icon-prime"></i>
		<span>FREE delivery Tomorrow</span>
	</div>
	<div data-asin="B0TESTITEM" class="s-result-item">
		<h2>
			<span>Test Product Item for Testing</span>
		</h2>
		<span class="a-price">
			<span>$49.99</span>
		</span>
		<span class="a-icon-alt">4.5 out of 5 stars</span>
		<span>1,234 ratings</span>
	</div>
</div>
<span class="s-result-count">1-48 of over 10,000 results</span>
</body>
</html>
`
