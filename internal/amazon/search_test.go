package amazon

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestParseSearchResultsHTML(t *testing.T) {
	// Sample HTML that mimics Amazon search results structure
	html := []byte(`
		<html>
		<body>
			<div data-asin="B08N5WRWNW" class="s-result-item">
				<h2 class="s-line-clamp-2">
					<a href="/dp/B08N5WRWNW">
						<span>Sony WH-1000XM4 Wireless Premium Noise Canceling Headphones</span>
					</a>
				</h2>
				<div class="a-price">
					<span class="a-offscreen">$278.00</span>
					<span class="a-price-whole">278</span>
				</div>
				<div class="a-price a-text-price">
					<span class="a-offscreen">$349.99</span>
				</div>
				<span aria-label="4.7 out of 5 stars">4.7 out of 5 stars</span>
				<span aria-label="52,431 ratings">52,431</span>
				<i class="a-icon-prime"></i>
			</div>
			<div data-asin="B0BXY1234Z" class="s-result-item">
				<h2>
					<a href="/dp/B0BXY1234Z">
						<span>Apple AirPods Pro (2nd Generation)</span>
					</a>
				</h2>
				<div class="a-price">
					<span class="a-offscreen">$189.99</span>
				</div>
				<span aria-label="4.8 out of 5 stars">4.8</span>
				<span class="a-size-base s-underline-text">89234</span>
				<span aria-label="Prime">Prime</span>
			</div>
			<div data-asin="B09ABC7890" class="s-result-item">
				<h2>
					<span class="s-title-instructions-style">Bose QuietComfort 45 Bluetooth Wireless Headphones</span>
				</h2>
				<div class="a-price">
					<span class="a-price-whole">249</span>
				</div>
				<span class="a-icon-alt">4.6 out of 5 stars</span>
				<span>31,256</span>
				<div class="a-size-base a-color-price">Out of Stock</div>
			</div>
		</body>
		</html>
	`)

	products, err := parseSearchResultsHTML(html)
	if err != nil {
		t.Fatalf("parseSearchResultsHTML failed: %v", err)
	}

	// Should find 3 products
	if len(products) != 3 {
		t.Errorf("Expected 3 products, got %d", len(products))
	}

	// Test first product (Sony headphones)
	if len(products) > 0 {
		p := products[0]
		if p.ASIN != "B08N5WRWNW" {
			t.Errorf("Expected ASIN B08N5WRWNW, got %s", p.ASIN)
		}
		if p.Title != "Sony WH-1000XM4 Wireless Premium Noise Canceling Headphones" {
			t.Errorf("Expected Sony headphones title, got %s", p.Title)
		}
		if p.Price != 278.00 {
			t.Errorf("Expected price 278.00, got %f", p.Price)
		}
		if p.OriginalPrice == nil || *p.OriginalPrice != 349.99 {
			t.Errorf("Expected original price 349.99, got %v", p.OriginalPrice)
		}
		if p.Rating != 4.7 {
			t.Errorf("Expected rating 4.7, got %f", p.Rating)
		}
		if p.ReviewCount != 52431 {
			t.Errorf("Expected review count 52431, got %d", p.ReviewCount)
		}
		if !p.Prime {
			t.Error("Expected Prime to be true")
		}
		if !p.InStock {
			t.Error("Expected InStock to be true")
		}
	}

	// Test second product (AirPods)
	if len(products) > 1 {
		p := products[1]
		if p.ASIN != "B0BXY1234Z" {
			t.Errorf("Expected ASIN B0BXY1234Z, got %s", p.ASIN)
		}
		if p.Price != 189.99 {
			t.Errorf("Expected price 189.99, got %f", p.Price)
		}
		if p.Rating != 4.8 {
			t.Errorf("Expected rating 4.8, got %f", p.Rating)
		}
		if p.ReviewCount != 89234 {
			t.Errorf("Expected review count 89234, got %d", p.ReviewCount)
		}
		if !p.Prime {
			t.Error("Expected Prime to be true")
		}
	}

	// Test third product (Bose - out of stock)
	if len(products) > 2 {
		p := products[2]
		if p.ASIN != "B09ABC7890" {
			t.Errorf("Expected ASIN B09ABC7890, got %s", p.ASIN)
		}
		if p.InStock {
			t.Error("Expected InStock to be false for out of stock product")
		}
	}
}

func TestParseSearchResultsHTML_EmptyHTML(t *testing.T) {
	html := []byte(`<html><body></body></html>`)

	products, err := parseSearchResultsHTML(html)
	if err != nil {
		t.Fatalf("parseSearchResultsHTML failed: %v", err)
	}

	if len(products) != 0 {
		t.Errorf("Expected 0 products for empty HTML, got %d", len(products))
	}
}

func TestParseSearchResultsHTML_InvalidHTML(t *testing.T) {
	html := []byte(`not valid html`)

	// Should still parse but return empty results
	products, err := parseSearchResultsHTML(html)
	if err != nil {
		t.Fatalf("parseSearchResultsHTML failed on invalid HTML: %v", err)
	}

	if len(products) != 0 {
		t.Errorf("Expected 0 products for invalid HTML, got %d", len(products))
	}
}

func TestParseSearchResultsHTML_MissingFields(t *testing.T) {
	// Product with ASIN but missing title and price
	html := []byte(`
		<html>
		<body>
			<div data-asin="B12345"></div>
		</body>
		</html>
	`)

	products, err := parseSearchResultsHTML(html)
	if err != nil {
		t.Fatalf("parseSearchResultsHTML failed: %v", err)
	}

	// Should not include products without title and price
	if len(products) != 0 {
		t.Errorf("Expected 0 products when missing required fields, got %d", len(products))
	}
}

func TestParseSearchResultsHTML_NoASIN(t *testing.T) {
	// Product without data-asin attribute
	html := []byte(`
		<html>
		<body>
			<div class="s-result-item">
				<h2><span>Some Product</span></h2>
				<div class="a-price"><span class="a-offscreen">$99.99</span></div>
			</div>
		</body>
		</html>
	`)

	products, err := parseSearchResultsHTML(html)
	if err != nil {
		t.Fatalf("parseSearchResultsHTML failed: %v", err)
	}

	// Should not include products without ASIN
	if len(products) != 0 {
		t.Errorf("Expected 0 products without ASIN, got %d", len(products))
	}
}

func TestParsePriceFromText(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"$29.99", 29.99},
		{"$1,299.99", 1299.99},
		{"$10", 10.0},
		{"99.95", 99.95},
		{"1234", 1234.0},
		{"$1,234,567.89", 1234567.89},
		{"", 0.0},
		{"no price here", 0.0},
		{"  $49.99  ", 49.99},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parsePriceFromText(tt.input)
			if result != tt.expected {
				t.Errorf("parsePriceFromText(%q) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseRating(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"4.5 out of 5 stars", 4.5},
		{"4.7 out of 5", 4.7},
		{"5.0 out of 5 stars", 5.0},
		{"3.2", 3.2},
		{"4.8", 4.8},
		{"", 0.0},
		{"no rating", 0.0},
		{"  4.6 out of 5 stars  ", 4.6},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseRating(tt.input)
			if result != tt.expected {
				t.Errorf("parseRating(%q) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseReviewCount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1,234 ratings", 1234},
		{"52,431", 52431},
		{"100", 100},
		{"1234567", 1234567},
		{"10,000 ratings", 10000},
		{"", 0},
		{"no count", 0},
		{"  5,678  ", 5678},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseReviewCount(tt.input)
			if result != tt.expected {
				t.Errorf("parseReviewCount(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseSearchResultsHTML_VariousPrimeIndicators(t *testing.T) {
	htmlWithIconPrime := []byte(`
		<html><body>
		<div data-asin="TEST1">
			<h2><span>Test Product 1</span></h2>
			<div class="a-price"><span class="a-offscreen">$99.99</span></div>
			<i class="a-icon-prime"></i>
		</div>
		</body></html>
	`)

	htmlWithClassPrime := []byte(`
		<html><body>
		<div data-asin="TEST2">
			<h2><span>Test Product 2</span></h2>
			<div class="a-price"><span class="a-offscreen">$99.99</span></div>
			<div class="s-prime">Prime</div>
		</div>
		</body></html>
	`)

	htmlWithAriaLabelPrime := []byte(`
		<html><body>
		<div data-asin="TEST3">
			<h2><span>Test Product 3</span></h2>
			<div class="a-price"><span class="a-offscreen">$99.99</span></div>
			<span aria-label="Prime">Prime eligible</span>
		</div>
		</body></html>
	`)

	tests := []struct {
		name string
		html []byte
	}{
		{"Icon Prime", htmlWithIconPrime},
		{"Class Prime", htmlWithClassPrime},
		{"Aria Label Prime", htmlWithAriaLabelPrime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, err := parseSearchResultsHTML(tt.html)
			if err != nil {
				t.Fatalf("parseSearchResultsHTML failed: %v", err)
			}

			if len(products) != 1 {
				t.Fatalf("Expected 1 product, got %d", len(products))
			}

			if !products[0].Prime {
				t.Error("Expected Prime to be true")
			}
		})
	}
}

func TestParseSearchResultsHTML_StockStatus(t *testing.T) {
	tests := []struct {
		name          string
		stockText     string
		expectedStock bool
	}{
		{"In Stock", "", true},
		{"Unavailable", "Currently unavailable", false},
		{"Out of Stock", "Out of Stock", false},
		{"Unavailable Caps", "UNAVAILABLE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := []byte(`
				<html><body>
				<div data-asin="TEST123">
					<h2><span>Test Product</span></h2>
					<div class="a-price"><span class="a-offscreen">$99.99</span></div>
					<div class="a-size-base a-color-secondary">` + tt.stockText + `</div>
				</div>
				</body></html>
			`)

			products, err := parseSearchResultsHTML(html)
			if err != nil {
				t.Fatalf("parseSearchResultsHTML failed: %v", err)
			}

			if len(products) != 1 {
				t.Fatalf("Expected 1 product, got %d", len(products))
			}

			if products[0].InStock != tt.expectedStock {
				t.Errorf("Expected InStock=%v for %q, got %v", tt.expectedStock, tt.stockText, products[0].InStock)
			}
		})
	}
}

func TestParseSearchResultsHTML_MultipleProducts(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<div data-asin="PROD1">
				<h2><span>Product 1</span></h2>
				<div class="a-price"><span class="a-offscreen">$10.00</span></div>
			</div>
			<div data-asin="PROD2">
				<h2><span>Product 2</span></h2>
				<div class="a-price"><span class="a-offscreen">$20.00</span></div>
			</div>
			<div data-asin="PROD3">
				<h2><span>Product 3</span></h2>
				<div class="a-price"><span class="a-offscreen">$30.00</span></div>
			</div>
			<div data-asin="PROD4">
				<h2><span>Product 4</span></h2>
				<div class="a-price"><span class="a-offscreen">$40.00</span></div>
			</div>
			<div data-asin="PROD5">
				<h2><span>Product 5</span></h2>
				<div class="a-price"><span class="a-offscreen">$50.00</span></div>
			</div>
		</body>
		</html>
	`)

	products, err := parseSearchResultsHTML(html)
	if err != nil {
		t.Fatalf("parseSearchResultsHTML failed: %v", err)
	}

	if len(products) != 5 {
		t.Errorf("Expected 5 products, got %d", len(products))
	}

	// Verify each product was parsed correctly
	expectedPrices := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	for i, p := range products {
		if p.Price != expectedPrices[i] {
			t.Errorf("Product %d: expected price %f, got %f", i, expectedPrices[i], p.Price)
		}
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	client := NewClient()

	_, err := client.Search("", models.SearchOptions{})
	if err == nil {
		t.Error("Expected error for empty query, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("Expected 'cannot be empty' error, got: %v", err)
	}
}

func TestSearch_URLConstruction(t *testing.T) {
	// Create a test server to intercept the request
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body></body></html>`))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	tests := []struct {
		name           string
		query          string
		opts           models.SearchOptions
		expectedParams map[string]string
	}{
		{
			name:  "Basic query",
			query: "headphones",
			opts:  models.SearchOptions{},
			expectedParams: map[string]string{
				"k": "headphones",
			},
		},
		{
			name:  "Query with category",
			query: "laptop",
			opts: models.SearchOptions{
				Category: "electronics",
			},
			expectedParams: map[string]string{
				"k": "laptop",
				"i": "electronics",
			},
		},
		{
			name:  "Query with price range",
			query: "phone",
			opts: models.SearchOptions{
				MinPrice: 100.50,
				MaxPrice: 500.75,
			},
			expectedParams: map[string]string{
				"k":          "phone",
				"low-price":  "10050",
				"high-price": "50075",
			},
		},
		{
			name:  "Query with Prime filter",
			query: "tablet",
			opts: models.SearchOptions{
				PrimeOnly: true,
			},
			expectedParams: map[string]string{
				"k":     "tablet",
				"prime": "true",
			},
		},
		{
			name:  "Query with page number",
			query: "mouse",
			opts: models.SearchOptions{
				Page: 3,
			},
			expectedParams: map[string]string{
				"k":    "mouse",
				"page": "3",
			},
		},
		{
			name:  "Query with all options",
			query: "keyboard",
			opts: models.SearchOptions{
				Category:  "computers",
				MinPrice:  50.0,
				MaxPrice:  150.0,
				PrimeOnly: true,
				Page:      2,
			},
			expectedParams: map[string]string{
				"k":          "keyboard",
				"i":          "computers",
				"low-price":  "5000",
				"high-price": "15000",
				"prime":      "true",
				"page":       "2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capturedURL = ""

			_, err := client.Search(tt.query, tt.opts)
			if err != nil {
				t.Fatalf("Search failed: %v", err)
			}

			// Parse the captured URL
			parsedURL, err := url.Parse(capturedURL)
			if err != nil {
				t.Fatalf("Failed to parse captured URL: %v", err)
			}

			// Check path
			if !strings.HasPrefix(parsedURL.Path, "/s") {
				t.Errorf("Expected path to start with /s, got %s", parsedURL.Path)
			}

			// Check query parameters
			queryParams := parsedURL.Query()
			for key, expectedValue := range tt.expectedParams {
				actualValue := queryParams.Get(key)
				if actualValue != expectedValue {
					t.Errorf("Expected param %s=%s, got %s", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestSearch_DefaultPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that page parameter is not set for page 1 (default)
		if r.URL.Query().Get("page") != "" {
			t.Error("Expected no page parameter for default page")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body></body></html>`))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	// Test with Page = 0 (should default to 1)
	_, err := client.Search("test", models.SearchOptions{Page: 0})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Test with Page = 1 (should not include page param)
	_, err = client.Search("test", models.SearchOptions{Page: 1})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
}

func TestSearch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	_, err := client.Search("test", models.SearchOptions{})
	if err == nil {
		t.Error("Expected error for non-200 status code, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "unexpected status code") {
		t.Errorf("Expected 'unexpected status code' error, got: %v", err)
	}
}

func TestSearch_CAPTCHADetection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
			<body>
				<p>Sorry, we just need to make sure you're not a robot.</p>
				<form action="/captcha">
					<input type="text" name="captcha">
				</form>
			</body>
			</html>
		`))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	_, err := client.Search("test", models.SearchOptions{})
	if err == nil {
		t.Error("Expected error for CAPTCHA detection, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "CAPTCHA") {
		t.Errorf("Expected 'CAPTCHA' error, got: %v", err)
	}
}

func TestSearch_SuccessfulResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
			<body>
				<div data-asin="B08N5WRWNW">
					<h2><span>Test Product 1</span></h2>
					<div class="a-price"><span class="a-offscreen">$99.99</span></div>
					<span aria-label="4.5 out of 5 stars">4.5</span>
					<i class="a-icon-prime"></i>
				</div>
				<div data-asin="B0BXY1234Z">
					<h2><span>Test Product 2</span></h2>
					<div class="a-price"><span class="a-offscreen">$149.99</span></div>
					<span aria-label="4.7 out of 5 stars">4.7</span>
				</div>
			</body>
			</html>
		`))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	resp, err := client.Search("test query", models.SearchOptions{})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.Query != "test query" {
		t.Errorf("Expected query 'test query', got %s", resp.Query)
	}

	if len(resp.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(resp.Results))
	}

	if resp.TotalResults != 2 {
		t.Errorf("Expected TotalResults 2, got %d", resp.TotalResults)
	}

	if resp.Page != 1 {
		t.Errorf("Expected Page 1, got %d", resp.Page)
	}

	// Verify first product
	if len(resp.Results) > 0 {
		p := resp.Results[0]
		if p.ASIN != "B08N5WRWNW" {
			t.Errorf("Expected ASIN B08N5WRWNW, got %s", p.ASIN)
		}
		if p.Title != "Test Product 1" {
			t.Errorf("Expected title 'Test Product 1', got %s", p.Title)
		}
		if p.Price != 99.99 {
			t.Errorf("Expected price 99.99, got %f", p.Price)
		}
		if !p.Prime {
			t.Error("Expected Prime to be true for first product")
		}
	}
}
