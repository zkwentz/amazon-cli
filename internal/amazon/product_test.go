package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestIsValidASIN(t *testing.T) {
	tests := []struct {
		name  string
		asin  string
		valid bool
	}{
		{"Valid ASIN", "B08N5WRWNW", true},
		{"Valid numeric ASIN", "0123456789", true},
		{"Too short", "B08N5WRW", false},
		{"Too long", "B08N5WRWNW1", false},
		{"Contains lowercase", "b08n5wrwnw", false},
		{"Contains special chars", "B08N5WRW-W", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidASIN(tt.asin)
			if result != tt.valid {
				t.Errorf("isValidASIN(%s) = %v, want %v", tt.asin, result, tt.valid)
			}
		})
	}
}

func TestParsePrice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Simple price", "$29.99", 29.99},
		{"Price with comma", "$1,234.56", 1234.56},
		{"Price without dollar sign", "99.99", 99.99},
		{"Price with spaces", " $49.99 ", 49.99},
		{"Invalid price", "N/A", 0},
		{"Empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePrice(tt.input)
			if result != tt.expected {
				t.Errorf("parsePrice(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseRating(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Valid rating", "4.5 out of 5 stars", 4.5},
		{"Perfect rating", "5.0 out of 5 stars", 5.0},
		{"Low rating", "2.3 out of 5 stars", 2.3},
		{"Invalid format", "No rating", 0},
		{"Empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseRating(tt.input)
			if result != tt.expected {
				t.Errorf("parseRating(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseReviewCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"With comma", "1,234 ratings", 1234},
		{"Without comma", "999 ratings", 999},
		{"Large number", "52,431 global ratings", 52431},
		{"Invalid format", "No reviews", 0},
		{"Empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseReviewCount(tt.input)
			if result != tt.expected {
				t.Errorf("parseReviewCount(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetProduct_InvalidASIN(t *testing.T) {
	cfg := &config.Config{
		RateLimiting: config.RateLimitConfig{
			MinDelayMs: 0,
			MaxDelayMs: 0,
			MaxRetries: 0,
		},
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetProduct("invalid")
	if err == nil {
		t.Error("Expected error for invalid ASIN, got nil")
	}

	if cliErr, ok := err.(*models.CLIError); ok {
		if cliErr.Code != models.ErrCodeInvalidInput {
			t.Errorf("Expected error code %s, got %s", models.ErrCodeInvalidInput, cliErr.Code)
		}
	} else {
		t.Error("Expected CLIError type")
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	t.Skip("Skipping integration test - would need HTTP client mocking")
	// This test would require mocking the HTTP client completely
	// In a production implementation, we'd inject an HTTP client interface
	// and provide a mock implementation for testing
}

func TestParseProductPage(t *testing.T) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Test Product</title></head>
<body>
	<h1 id="productTitle">Test Product Title</h1>
	<span class="a-price"><span class="a-offscreen">$29.99</span></span>
	<span class="a-icon-alt">4.5 out of 5 stars</span>
	<span id="acrCustomerReviewText">1,234 ratings</span>
	<div id="priceBadging_feature_div"><i class="a-icon-prime"></i></div>
	<div id="availability"><span>In Stock</span></div>
	<div id="feature-bullets">
		<ul>
			<li><span class="a-list-item">Feature 1</span></li>
			<li><span class="a-list-item">Feature 2</span></li>
		</ul>
	</div>
</body>
</html>
`

	product, err := parseProductPage("B08N5WRWNW", []byte(html))
	if err != nil {
		t.Fatalf("Failed to parse product page: %v", err)
	}

	if product.ASIN != "B08N5WRWNW" {
		t.Errorf("Expected ASIN B08N5WRWNW, got %s", product.ASIN)
	}

	if product.Title != "Test Product Title" {
		t.Errorf("Expected title 'Test Product Title', got '%s'", product.Title)
	}

	if product.Price != 29.99 {
		t.Errorf("Expected price 29.99, got %v", product.Price)
	}

	if product.Rating != 4.5 {
		t.Errorf("Expected rating 4.5, got %v", product.Rating)
	}

	if product.ReviewCount != 1234 {
		t.Errorf("Expected review count 1234, got %v", product.ReviewCount)
	}

	if !product.Prime {
		t.Error("Expected Prime to be true")
	}

	if !product.InStock {
		t.Error("Expected InStock to be true")
	}

	if len(product.Features) != 2 {
		t.Errorf("Expected 2 features, got %d", len(product.Features))
	}
}
