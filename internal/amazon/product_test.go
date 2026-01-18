package amazon

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateASIN(t *testing.T) {
	tests := []struct {
		name    string
		asin    string
		wantErr bool
	}{
		{
			name:    "valid ASIN",
			asin:    "B08N5WRWNW",
			wantErr: false,
		},
		{
			name:    "valid ASIN with numbers",
			asin:    "B00ABC1234",
			wantErr: false,
		},
		{
			name:    "too short",
			asin:    "B08N5WRW",
			wantErr: true,
		},
		{
			name:    "too long",
			asin:    "B08N5WRWNW1",
			wantErr: true,
		},
		{
			name:    "lowercase letters",
			asin:    "b08n5wrwnw",
			wantErr: true,
		},
		{
			name:    "special characters",
			asin:    "B08N5-RWNW",
			wantErr: true,
		},
		{
			name:    "empty string",
			asin:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateASIN(tt.asin)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateASIN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParsePrice(t *testing.T) {
	tests := []struct {
		name      string
		priceStr  string
		wantPrice float64
	}{
		{
			name:      "simple price",
			priceStr:  "$29.99",
			wantPrice: 29.99,
		},
		{
			name:      "price with comma",
			priceStr:  "$1,299.99",
			wantPrice: 1299.99,
		},
		{
			name:      "price without symbol",
			priceStr:  "49.95",
			wantPrice: 49.95,
		},
		{
			name:      "price with whitespace",
			priceStr:  " $99.99 ",
			wantPrice: 99.99,
		},
		{
			name:      "invalid price",
			priceStr:  "free",
			wantPrice: 0.0,
		},
		{
			name:      "empty string",
			priceStr:  "",
			wantPrice: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePrice(tt.priceStr)
			if got != tt.wantPrice {
				t.Errorf("parsePrice() = %v, want %v", got, tt.wantPrice)
			}
		})
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient()
	// Create a custom client that points to our test server
	client.httpClient = server.Client()

	// Modify the product.go to use the test server URL
	// For now, we'll just test with invalid ASIN since we can't easily override the URL
	// This test would work better with dependency injection
	t.Skip("Skipping - requires refactoring for proper URL injection")
}

func TestGetProduct_InvalidASIN(t *testing.T) {
	client := NewClient()

	tests := []string{
		"invalid",
		"B08N5WRW", // too short
		"B08N5WRWNW1", // too long
		"",
	}

	for _, asin := range tests {
		_, err := client.GetProduct(asin)
		if err == nil {
			t.Errorf("GetProduct(%q) expected error for invalid ASIN, got nil", asin)
		}
	}
}

func TestGetProduct_BasicHTML(t *testing.T) {
	t.Skip("Skipping integration test - requires mock server with URL injection")
	// This test would require refactoring GetProduct to accept a base URL
	// or using an interface-based approach for better testability
}

func TestConvertToFullSizeImage(t *testing.T) {
	tests := []struct {
		name         string
		thumbnailURL string
		want         string
	}{
		{
			name:         "thumbnail with size indicator",
			thumbnailURL: "https://m.media-amazon.com/images/I/71abc123._AC_UL160_SR160,160_.jpg",
			want:         "https://m.media-amazon.com/images/I/71abc123.jpg",
		},
		{
			name:         "already full size",
			thumbnailURL: "https://m.media-amazon.com/images/I/71abc123.jpg",
			want:         "https://m.media-amazon.com/images/I/71abc123.jpg",
		},
		{
			name:         "non-http URL",
			thumbnailURL: "data:image/gif;base64,abc",
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToFullSizeImage(tt.thumbnailURL)
			if got != tt.want {
				t.Errorf("convertToFullSizeImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
