package amazon

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

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
