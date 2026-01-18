package amazon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// TestAddToCart tests the AddToCart functionality
func TestAddToCart(t *testing.T) {
	tests := []struct {
		name        string
		asin        string
		quantity    int
		expectError bool
	}{
		{
			name:        "Valid add to cart",
			asin:        "B08N5WRWNW",
			quantity:    1,
			expectError: false,
		},
		{
			name:        "Empty ASIN",
			asin:        "",
			quantity:    1,
			expectError: true,
		},
		{
			name:        "Zero quantity",
			asin:        "B08N5WRWNW",
			quantity:    0,
			expectError: true,
		},
		{
			name:        "Negative quantity",
			asin:        "B08N5WRWNW",
			quantity:    -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<html><body>Cart</body></html>`))
			}))
			defer server.Close()

			client := &Client{
				httpClient: server.Client(),
			}
			service := NewCartService(client)

			_, err := service.AddToCart(tt.asin, tt.quantity)
			if (err != nil) != tt.expectError {
				t.Errorf("AddToCart() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// TestRemoveFromCart tests the RemoveFromCart functionality
func TestRemoveFromCart(t *testing.T) {
	tests := []struct {
		name        string
		asin        string
		expectError bool
	}{
		{
			name:        "Valid removal",
			asin:        "B08N5WRWNW",
			expectError: false,
		},
		{
			name:        "Empty ASIN",
			asin:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				// Return a cart with the item for GET requests
				if r.Method == "GET" {
					w.Write([]byte(`<html><body>
						<div data-asin="B08N5WRWNW">Item</div>
					</body></html>`))
				} else {
					w.Write([]byte(`<html><body>Updated cart</body></html>`))
				}
			}))
			defer server.Close()

			client := &Client{
				httpClient: server.Client(),
			}
			service := NewCartService(client)

			_, err := service.RemoveFromCart(tt.asin)
			if (err != nil) != tt.expectError {
				t.Errorf("RemoveFromCart() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// TestGetCart tests retrieving cart contents
func TestGetCart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body>Cart contents</body></html>`))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
	}
	service := NewCartService(client)

	cart, err := service.GetCart()
	if err != nil {
		t.Errorf("GetCart() error = %v", err)
	}
	if cart == nil {
		t.Error("GetCart() returned nil cart")
	}
}

// TestClearCart tests clearing all items from cart
func TestClearCart(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// First call returns cart with items, subsequent calls return empty cart
		if callCount == 0 {
			w.Write([]byte(`<html><body>
				<div data-asin="B08N5WRWNW">Item 1</div>
				<div data-asin="B08N5WRWNY">Item 2</div>
			</body></html>`))
		} else {
			w.Write([]byte(`<html><body>Empty cart</body></html>`))
		}
		callCount++
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
	}
	service := NewCartService(client)

	err := service.ClearCart()
	if err != nil {
		t.Errorf("ClearCart() error = %v", err)
	}
}

// TestGetAddresses tests retrieving saved addresses
func TestGetAddresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body>Addresses list</body></html>`))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
	}
	service := NewCartService(client)

	addresses, err := service.GetAddresses()
	if err != nil {
		t.Errorf("GetAddresses() error = %v", err)
	}
	if addresses == nil {
		t.Error("GetAddresses() returned nil")
	}
}

// TestGetPaymentMethods tests retrieving saved payment methods
func TestGetPaymentMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body>Payment methods list</body></html>`))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
	}
	service := NewCartService(client)

	methods, err := service.GetPaymentMethods()
	if err != nil {
		t.Errorf("GetPaymentMethods() error = %v", err)
	}
	if methods == nil {
		t.Error("GetPaymentMethods() returned nil")
	}
}

// TestPreviewCheckout tests checkout preview functionality
func TestPreviewCheckout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body>Checkout preview</body></html>`))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
	}
	service := NewCartService(client)

	preview, err := service.PreviewCheckout("addr_123", "pay_456")
	if err != nil {
		t.Errorf("PreviewCheckout() error = %v", err)
	}
	if preview == nil {
		t.Error("PreviewCheckout() returned nil")
	}
}

// TestCompleteCheckout tests completing a purchase
func TestCompleteCheckout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == "POST" {
			w.Write([]byte(`<html><body>Order confirmed</body></html>`))
		} else {
			w.Write([]byte(`<html><body>Checkout preview</body></html>`))
		}
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
	}
	service := NewCartService(client)

	confirmation, err := service.CompleteCheckout("addr_123", "pay_456")
	if err != nil {
		t.Errorf("CompleteCheckout() error = %v", err)
	}
	if confirmation == nil {
		t.Error("CompleteCheckout() returned nil")
	}
}

// TestNewCartService tests cart service creation
func TestNewCartService(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{},
	}
	service := NewCartService(client)
	if service == nil {
		t.Error("NewCartService() returned nil")
	}
	if service.client != client {
		t.Error("NewCartService() did not set client properly")
	}
}

// TestCartItemStructure tests the CartItem model structure
func TestCartItemStructure(t *testing.T) {
	item := models.CartItem{
		ASIN:     "B08N5WRWNW",
		Title:    "Test Product",
		Price:    29.99,
		Quantity: 2,
		Subtotal: 59.98,
		Prime:    true,
		InStock:  true,
	}

	if item.ASIN != "B08N5WRWNW" {
		t.Errorf("Expected ASIN B08N5WRWNW, got %s", item.ASIN)
	}
	if item.Quantity != 2 {
		t.Errorf("Expected quantity 2, got %d", item.Quantity)
	}
	if item.Subtotal != 59.98 {
		t.Errorf("Expected subtotal 59.98, got %f", item.Subtotal)
	}
}

// TestCartStructure tests the Cart model structure
func TestCartStructure(t *testing.T) {
	cart := models.Cart{
		Items: []models.CartItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Product 1",
				Price:    29.99,
				Quantity: 1,
				Subtotal: 29.99,
				Prime:    true,
				InStock:  true,
			},
		},
		Subtotal:     29.99,
		EstimatedTax: 2.40,
		Total:        32.39,
		ItemCount:    1,
	}

	if len(cart.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(cart.Items))
	}
	if cart.Total != 32.39 {
		t.Errorf("Expected total 32.39, got %f", cart.Total)
	}
}
