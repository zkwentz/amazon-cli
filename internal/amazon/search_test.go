package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestSearch(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name  string
		query string
		opts  models.SearchOptions
	}{
		{
			name:  "basic search",
			query: "wireless headphones",
			opts:  models.SearchOptions{Page: 1},
		},
		{
			name:  "search with category",
			query: "coffee",
			opts: models.SearchOptions{
				Category: "electronics",
				Page:     1,
			},
		},
		{
			name:  "search with price range",
			query: "laptop",
			opts: models.SearchOptions{
				MinPrice: 500.0,
				MaxPrice: 1000.0,
				Page:     1,
			},
		},
		{
			name:  "search with prime only",
			query: "books",
			opts: models.SearchOptions{
				PrimeOnly: true,
				Page:      1,
			},
		},
		{
			name:  "search with pagination",
			query: "keyboards",
			opts: models.SearchOptions{
				Page: 2,
			},
		},
		{
			name:  "search with all options",
			query: "monitor",
			opts: models.SearchOptions{
				Category:  "electronics",
				MinPrice:  200.0,
				MaxPrice:  500.0,
				PrimeOnly: true,
				Page:      3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.Search(tt.query, tt.opts)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if response == nil {
				t.Error("expected response, got nil")
				return
			}

			if response.Query != tt.query {
				t.Errorf("expected query %s, got %s", tt.query, response.Query)
			}

			if response.Page != tt.opts.Page {
				t.Errorf("expected page %d, got %d", tt.opts.Page, response.Page)
			}

			if response.Results == nil {
				t.Error("expected results slice, got nil")
			}
		})
	}
}

func TestSearchOptions(t *testing.T) {
	client := NewClient()

	// Test that search options are properly constructed
	opts := models.SearchOptions{
		Category:  "electronics",
		MinPrice:  10.0,
		MaxPrice:  100.0,
		PrimeOnly: true,
		Page:      2,
	}

	response, err := client.Search("test query", opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Query != "test query" {
		t.Errorf("expected query 'test query', got %s", response.Query)
	}

	if response.Page != 2 {
		t.Errorf("expected page 2, got %d", response.Page)
	}
}
