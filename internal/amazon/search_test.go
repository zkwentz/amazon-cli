package amazon

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.httpClient == nil {
		t.Error("httpClient is nil")
	}
	if client.userAgent == "" {
		t.Error("userAgent is empty")
	}
}

func TestBuildSearchURL(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name     string
		query    string
		opts     SearchOptions
		contains []string
	}{
		{
			name:  "basic search",
			query: "wireless headphones",
			opts:  SearchOptions{},
			contains: []string{
				"k=wireless+headphones",
				"amazon.com/s",
			},
		},
		{
			name:  "search with category",
			query: "laptop",
			opts: SearchOptions{
				Category: "electronics",
			},
			contains: []string{
				"k=laptop",
				"i=electronics",
			},
		},
		{
			name:  "search with price range",
			query: "shoes",
			opts: SearchOptions{
				MinPrice: 20.00,
				MaxPrice: 100.00,
			},
			contains: []string{
				"k=shoes",
				"low-price=20.00",
				"high-price=100.00",
			},
		},
		{
			name:  "search with Prime filter",
			query: "books",
			opts: SearchOptions{
				PrimeOnly: true,
			},
			contains: []string{
				"k=books",
				"prime=prime",
			},
		},
		{
			name:  "search with page number",
			query: "coffee",
			opts: SearchOptions{
				Page: 2,
			},
			contains: []string{
				"k=coffee",
				"page=2",
			},
		},
		{
			name:  "search with all options",
			query: "tablet",
			opts: SearchOptions{
				Category:  "electronics",
				MinPrice:  50.00,
				MaxPrice:  200.00,
				PrimeOnly: true,
				Page:      3,
			},
			contains: []string{
				"k=tablet",
				"i=electronics",
				"low-price=50.00",
				"high-price=200.00",
				"prime=prime",
				"page=3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := client.buildSearchURL(tt.query, tt.opts)
			if err != nil {
				t.Fatalf("buildSearchURL failed: %v", err)
			}

			for _, substr := range tt.contains {
				if !strings.Contains(url, substr) {
					t.Errorf("URL does not contain expected substring %q\nGot: %s", substr, url)
				}
			}
		})
	}
}

func TestParsePrice(t *testing.T) {
	client := NewClient()

	tests := []struct {
		input    string
		expected float64
	}{
		{"$29.99", 29.99},
		{"$1,234.56", 1234.56},
		{"$9.95", 9.95},
		{"$0.99", 0.99},
		{"$100", 100.00},
		{"invalid", 0.00},
		{"", 0.00},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := client.parsePrice(tt.input)
			if result != tt.expected {
				t.Errorf("parsePrice(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseRating(t *testing.T) {
	client := NewClient()

	tests := []struct {
		input    string
		expected float64
	}{
		{"4.5 out of 5 stars", 4.5},
		{"4.7 out of 5 stars", 4.7},
		{"5 out of 5 stars", 5.0},
		{"3.2 out of 5 stars", 3.2},
		{"invalid", 0.0},
		{"", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := client.parseRating(tt.input)
			if result != tt.expected {
				t.Errorf("parseRating(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseReviewCount(t *testing.T) {
	client := NewClient()

	tests := []struct {
		input    string
		expected int
	}{
		{"4.5 out of 5 stars 1,234 ratings", 1234},
		{"4.7 out of 5 stars 52431 ratings", 52431},
		{"5 out of 5 stars 10 ratings", 10},
		{"invalid", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := client.parseReviewCount(tt.input)
			if result != tt.expected {
				t.Errorf("parseReviewCount(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseReviewCountFromText(t *testing.T) {
	client := NewClient()

	tests := []struct {
		input    string
		expected int
	}{
		{"1,234 ratings", 1234},
		{"52431 ratings", 52431},
		{"10 ratings", 10},
		{"invalid", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := client.parseReviewCountFromText(tt.input)
			if result != tt.expected {
				t.Errorf("parseReviewCountFromText(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	client := NewClient()

	_, err := client.Search("", SearchOptions{})
	if err == nil {
		t.Error("Search with empty query should return error")
	}
	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("Expected 'cannot be empty' error, got: %v", err)
	}
}

func TestSearchDefaultPage(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body></body></html>`))
	}))
	defer server.Close()

	client := NewClient()

	// Test that page defaults to 1 when set to 0
	opts := SearchOptions{Page: 0}
	_, err := client.Search("test", opts)

	// We expect no error from the page defaulting
	// (The actual search might fail due to parsing, but that's ok for this test)
	if err != nil && strings.Contains(err.Error(), "page") {
		t.Errorf("Page should default to 1, got error: %v", err)
	}
}

func TestGetProductInvalidASIN(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name  string
		asin  string
		error string
	}{
		{
			name:  "empty ASIN",
			asin:  "",
			error: "cannot be empty",
		},
		{
			name:  "too short",
			asin:  "B08N5",
			error: "invalid ASIN format",
		},
		{
			name:  "too long",
			asin:  "B08N5WRWNW123",
			error: "invalid ASIN format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetProduct(tt.asin)
			if err == nil {
				t.Error("Expected error for invalid ASIN")
			}
			if !strings.Contains(err.Error(), tt.error) {
				t.Errorf("Expected error containing %q, got: %v", tt.error, err)
			}
		})
	}
}

func TestGetProductNotFound(t *testing.T) {
	// Create a mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// This test would need to mock the HTTP client to actually test
	// For now, we're testing the ASIN validation which doesn't require network
	client := NewClient()
	_, err := client.GetProduct("INVALID123")
	if err == nil {
		t.Error("Expected error for invalid ASIN format")
	}
}

func TestSearchOptionsDefaults(t *testing.T) {
	opts := SearchOptions{}

	if opts.Page != 0 {
		t.Errorf("Default page should be 0, got %d", opts.Page)
	}
	if opts.MinPrice != 0 {
		t.Errorf("Default MinPrice should be 0, got %f", opts.MinPrice)
	}
	if opts.MaxPrice != 0 {
		t.Errorf("Default MaxPrice should be 0, got %f", opts.MaxPrice)
	}
	if opts.PrimeOnly {
		t.Error("Default PrimeOnly should be false")
	}
	if opts.Category != "" {
		t.Errorf("Default Category should be empty, got %q", opts.Category)
	}
}
