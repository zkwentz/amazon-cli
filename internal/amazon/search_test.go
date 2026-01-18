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
		wantErr       bool
		wantMinCount  int
		wantMaxCount  int
		validateFunc  func(*testing.T, *models.SearchResponse)
	}{
		{
			name:         "empty query returns error",
			query:        "",
			opts:         models.SearchOptions{},
			wantErr:      true,
			wantMinCount: 0,
			wantMaxCount: 0,
		},
		{
			name:         "basic search returns all products",
			query:        "headphones",
			opts:         models.SearchOptions{},
			wantErr:      false,
			wantMinCount: 4,
			wantMaxCount: 4,
			validateFunc: func(t *testing.T, resp *models.SearchResponse) {
				if resp.Query != "headphones" {
					t.Errorf("expected query 'headphones', got '%s'", resp.Query)
				}
				if resp.Page != 1 {
					t.Errorf("expected page 1, got %d", resp.Page)
				}
			},
		},
		{
			name:         "prime-only filter returns only Prime products",
			query:        "headphones",
			opts:         models.SearchOptions{PrimeOnly: true},
			wantErr:      false,
			wantMinCount: 2,
			wantMaxCount: 2,
			validateFunc: func(t *testing.T, resp *models.SearchResponse) {
				for _, product := range resp.Results {
					if !product.Prime {
						t.Errorf("prime-only filter failed: product %s (ASIN: %s) is not Prime eligible",
							product.Title, product.ASIN)
					}
				}
			},
		},
		{
			name:         "prime-only with min price filter",
			query:        "headphones",
			opts:         models.SearchOptions{PrimeOnly: true, MinPrice: 200.00},
			wantErr:      false,
			wantMinCount: 1,
			wantMaxCount: 1,
			validateFunc: func(t *testing.T, resp *models.SearchResponse) {
				for _, product := range resp.Results {
					if !product.Prime {
						t.Errorf("product %s is not Prime eligible", product.Title)
					}
					if product.Price < 200.00 {
						t.Errorf("product %s price %.2f is below minimum 200.00", product.Title, product.Price)
					}
				}
			},
		},
		{
			name:         "prime-only with max price filter",
			query:        "headphones",
			opts:         models.SearchOptions{PrimeOnly: true, MaxPrice: 250.00},
			wantErr:      false,
			wantMinCount: 1,
			wantMaxCount: 1,
			validateFunc: func(t *testing.T, resp *models.SearchResponse) {
				for _, product := range resp.Results {
					if !product.Prime {
						t.Errorf("product %s is not Prime eligible", product.Title)
					}
					if product.Price > 250.00 {
						t.Errorf("product %s price %.2f exceeds maximum 250.00", product.Title, product.Price)
					}
				}
			},
		},
		{
			name:         "prime-only with price range",
			query:        "headphones",
			opts:         models.SearchOptions{PrimeOnly: true, MinPrice: 100.00, MaxPrice: 300.00},
			wantErr:      false,
			wantMinCount: 2,
			wantMaxCount: 2,
			validateFunc: func(t *testing.T, resp *models.SearchResponse) {
				for _, product := range resp.Results {
					if !product.Prime {
						t.Errorf("product %s is not Prime eligible", product.Title)
					}
					if product.Price < 100.00 || product.Price > 300.00 {
						t.Errorf("product %s price %.2f is outside range 100.00-300.00", product.Title, product.Price)
					}
				}
			},
		},
		{
			name:         "non-prime with price filter",
			query:        "headphones",
			opts:         models.SearchOptions{PrimeOnly: false, MaxPrice: 100.00},
			wantErr:      false,
			wantMinCount: 2,
			wantMaxCount: 2,
			validateFunc: func(t *testing.T, resp *models.SearchResponse) {
				nonPrimeCount := 0
				for _, product := range resp.Results {
					if product.Price > 100.00 {
						t.Errorf("product %s price %.2f exceeds maximum 100.00", product.Title, product.Price)
					}
					if !product.Prime {
						nonPrimeCount++
					}
				}
				// Should include non-Prime products when PrimeOnly is false
				if nonPrimeCount == 0 {
					t.Error("expected at least one non-Prime product in results")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Search(tt.query, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if resp == nil {
				t.Fatal("expected response, got nil")
			}

			resultCount := len(resp.Results)
			if resultCount < tt.wantMinCount || resultCount > tt.wantMaxCount {
				t.Errorf("Search() returned %d results, want between %d and %d",
					resultCount, tt.wantMinCount, tt.wantMaxCount)
			}

			if resp.TotalResults != resultCount {
				t.Errorf("TotalResults = %d, but Results length = %d",
					resp.TotalResults, resultCount)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, resp)
			}
		})
	}
}

func TestApplyFilters(t *testing.T) {
	testProducts := []models.Product{
		{
			ASIN:  "PRIME001",
			Title: "Prime Product 1",
			Price: 50.00,
			Prime: true,
		},
		{
			ASIN:  "PRIME002",
			Title: "Prime Product 2",
			Price: 150.00,
			Prime: true,
		},
		{
			ASIN:  "NONPRIME001",
			Title: "Non-Prime Product 1",
			Price: 75.00,
			Prime: false,
		},
		{
			ASIN:  "NONPRIME002",
			Title: "Non-Prime Product 2",
			Price: 200.00,
			Prime: false,
		},
	}

	tests := []struct {
		name      string
		opts      models.SearchOptions
		wantCount int
		wantASINs []string
	}{
		{
			name:      "no filters returns all products",
			opts:      models.SearchOptions{},
			wantCount: 4,
			wantASINs: []string{"PRIME001", "PRIME002", "NONPRIME001", "NONPRIME002"},
		},
		{
			name:      "prime-only filter",
			opts:      models.SearchOptions{PrimeOnly: true},
			wantCount: 2,
			wantASINs: []string{"PRIME001", "PRIME002"},
		},
		{
			name:      "min price filter",
			opts:      models.SearchOptions{MinPrice: 100.00},
			wantCount: 2,
			wantASINs: []string{"PRIME002", "NONPRIME002"},
		},
		{
			name:      "max price filter",
			opts:      models.SearchOptions{MaxPrice: 100.00},
			wantCount: 2,
			wantASINs: []string{"PRIME001", "NONPRIME001"},
		},
		{
			name:      "price range filter",
			opts:      models.SearchOptions{MinPrice: 60.00, MaxPrice: 160.00},
			wantCount: 2,
			wantASINs: []string{"PRIME002", "NONPRIME001"},
		},
		{
			name:      "prime-only with min price",
			opts:      models.SearchOptions{PrimeOnly: true, MinPrice: 100.00},
			wantCount: 1,
			wantASINs: []string{"PRIME002"},
		},
		{
			name:      "prime-only with max price",
			opts:      models.SearchOptions{PrimeOnly: true, MaxPrice: 100.00},
			wantCount: 1,
			wantASINs: []string{"PRIME001"},
		},
		{
			name:      "prime-only with price range",
			opts:      models.SearchOptions{PrimeOnly: true, MinPrice: 40.00, MaxPrice: 160.00},
			wantCount: 2,
			wantASINs: []string{"PRIME001", "PRIME002"},
		},
		{
			name:      "filters excluding all products",
			opts:      models.SearchOptions{PrimeOnly: true, MinPrice: 300.00},
			wantCount: 0,
			wantASINs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyFilters(testProducts, tt.opts)

			if len(result) != tt.wantCount {
				t.Errorf("applyFilters() returned %d products, want %d", len(result), tt.wantCount)
			}

			// Check that the correct ASINs are present
			resultASINs := make(map[string]bool)
			for _, product := range result {
				resultASINs[product.ASIN] = true
			}

			for _, wantASIN := range tt.wantASINs {
				if !resultASINs[wantASIN] {
					t.Errorf("expected ASIN %s in results, but not found", wantASIN)
				}
			}

			// Verify Prime filter
			if tt.opts.PrimeOnly {
				for _, product := range result {
					if !product.Prime {
						t.Errorf("prime-only filter failed: product %s is not Prime", product.ASIN)
					}
				}
			}

			// Verify price filters
			for _, product := range result {
				if tt.opts.MinPrice > 0 && product.Price < tt.opts.MinPrice {
					t.Errorf("product %s price %.2f is below minimum %.2f",
						product.ASIN, product.Price, tt.opts.MinPrice)
				}
				if tt.opts.MaxPrice > 0 && product.Price > tt.opts.MaxPrice {
					t.Errorf("product %s price %.2f exceeds maximum %.2f",
						product.ASIN, product.Price, tt.opts.MaxPrice)
				}
			}
		})
	}
}

func TestSearchPageNumber(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		page     int
		wantPage int
	}{
		{
			name:     "page 0 defaults to 1",
			query:    "test",
			page:     0,
			wantPage: 1,
		},
		{
			name:     "page 2 is preserved",
			query:    "test",
			page:     2,
			wantPage: 2,
		},
		{
			name:     "page 10 is preserved",
			query:    "test",
			page:     10,
			wantPage: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := models.SearchOptions{Page: tt.page}
			resp, err := Search(tt.query, opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.Page != tt.wantPage {
				t.Errorf("Search() page = %d, want %d", resp.Page, tt.wantPage)
			}
		})
	}
}
