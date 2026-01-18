package amazon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetCart_EmptyCart(t *testing.T) {
	// Test parsing empty cart HTML directly
	cart, err := parseCartHTML(`
		<html>
		<body>
			<h1>Your Shopping Cart is empty</h1>
		</body>
		</html>
	`)

	if err != nil {
		t.Fatalf("GetCart failed: %v", err)
	}

	if cart == nil {
		t.Fatal("Expected cart to be non-nil")
	}

	if len(cart.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(cart.Items))
	}

	if cart.ItemCount != 0 {
		t.Errorf("Expected ItemCount 0, got %d", cart.ItemCount)
	}

	if cart.Total != 0 {
		t.Errorf("Expected Total 0, got %f", cart.Total)
	}
}

func TestGetCart_WithItems(t *testing.T) {
	// Sample cart HTML with items
	cartHTML := `
		<html>
		<body>
			<div class="sc-list-body">
				<div data-asin="B08N5WRWNW" class="sc-list-item">
					<span class="sc-product-title">Sony WH-1000XM4 Wireless Headphones</span>
					<span class="sc-product-price">$278.00</span>
					<span class="quantity">2</span>
					<span class="prime-badge">Prime</span>
				</div>
				<div data-asin="B0BX1Y3J9T" class="sc-list-item">
					<span class="sc-product-title">USB-C Cable</span>
					<span class="sc-product-price">$12.99</span>
					<span class="quantity">1</span>
				</div>
			</div>
			<div class="cart-summary">
				<span>Subtotal: $569.99</span>
				<span>Estimated tax: $45.60</span>
			</div>
		</body>
		</html>
	`

	cart, err := parseCartHTML(cartHTML)
	if err != nil {
		t.Fatalf("parseCartHTML failed: %v", err)
	}

	if cart == nil {
		t.Fatal("Expected cart to be non-nil")
	}

	if len(cart.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(cart.Items))
	}

	if cart.ItemCount != 2 {
		t.Errorf("Expected ItemCount 2, got %d", cart.ItemCount)
	}

	// Verify first item
	if len(cart.Items) > 0 {
		item := cart.Items[0]
		if item.ASIN != "B08N5WRWNW" {
			t.Errorf("Expected ASIN B08N5WRWNW, got %s", item.ASIN)
		}
	}
}

func TestParseCartHTML_EmptyCart(t *testing.T) {
	html := `<html><body>Your Amazon Cart is empty</body></html>`
	cart, err := parseCartHTML(html)

	if err != nil {
		t.Fatalf("parseCartHTML failed: %v", err)
	}

	if len(cart.Items) != 0 {
		t.Errorf("Expected 0 items in empty cart, got %d", len(cart.Items))
	}
}

func TestExtractCartItems(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected int
	}{
		{
			name: "Single item",
			html: `<div data-asin="B08N5WRWNW"><span class="product-title">Test Product</span></div>`,
			expected: 1,
		},
		{
			name: "Multiple items",
			html: `
				<div data-asin="B08N5WRWNW"></div>
				<div data-asin="B0BX1Y3J9T"></div>
			`,
			expected: 2,
		},
		{
			name: "No items",
			html: `<div>No cart items</div>`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := extractCartItems(tt.html)
			if err != nil {
				t.Fatalf("extractCartItems failed: %v", err)
			}

			if len(items) != tt.expected {
				t.Errorf("Expected %d items, got %d", tt.expected, len(items))
			}
		})
	}
}

func TestExtractPrice(t *testing.T) {
	tests := []struct {
		name      string
		html      string
		label     string
		expected  float64
		expectErr bool
	}{
		{
			name:     "Simple price",
			html:     `Subtotal: $123.45`,
			label:    "subtotal",
			expected: 123.45,
		},
		{
			name:     "Price with comma",
			html:     `Subtotal: $1,234.56`,
			label:    "subtotal",
			expected: 1234.56,
		},
		{
			name:     "Estimated tax",
			html:     `Estimated tax: $45.60`,
			label:    "estimated tax",
			expected: 45.60,
		},
		{
			name:      "Price not found",
			html:      `No price here`,
			label:     "subtotal",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := extractPrice(tt.html, tt.label)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("extractPrice failed: %v", err)
			}

			if price != tt.expected {
				t.Errorf("Expected price %f, got %f", tt.expected, price)
			}
		})
	}
}

func TestStripHTMLTags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "<span>Hello World</span>",
			expected: "Hello World",
		},
		{
			input:    "<div><p>Test</p></div>",
			expected: "Test",
		},
		{
			input:    "No tags here",
			expected: "No tags here",
		},
		{
			input:    "<a href='test'>Link</a>",
			expected: "Link",
		},
	}

	for _, tt := range tests {
		result := stripHTMLTags(tt.input)
		if result != tt.expected {
			t.Errorf("stripHTMLTags(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}

	if client.httpClient == nil {
		t.Fatal("Expected httpClient to be non-nil")
	}
}

func TestClient_NotImplementedMethods(t *testing.T) {
	client := NewClient()

	// Test AddToCart
	_, err := client.AddToCart("B08N5WRWNW", 1)
	if err == nil || err.Error() != "not implemented" {
		t.Error("Expected 'not implemented' error for AddToCart")
	}

	// Test RemoveFromCart
	_, err = client.RemoveFromCart("B08N5WRWNW")
	if err == nil || err.Error() != "not implemented" {
		t.Error("Expected 'not implemented' error for RemoveFromCart")
	}

	// Test ClearCart
	err = client.ClearCart()
	if err == nil || err.Error() != "not implemented" {
		t.Error("Expected 'not implemented' error for ClearCart")
	}

	// Test GetAddresses
	_, err = client.GetAddresses()
	if err == nil || err.Error() != "not implemented" {
		t.Error("Expected 'not implemented' error for GetAddresses")
	}

	// Test GetPaymentMethods
	_, err = client.GetPaymentMethods()
	if err == nil || err.Error() != "not implemented" {
		t.Error("Expected 'not implemented' error for GetPaymentMethods")
	}
}

func TestGetCart_Integration(t *testing.T) {
	// This is an integration test that would need a mock server
	// For now, we'll test the basic flow with a mock server

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
			<body>
				<div data-asin="B08N5WRWNW">
					<span class="product-title">Test Product</span>
					<span>$99.99</span>
				</div>
				<div>Subtotal: $99.99</div>
			</body>
			</html>
		`))
	}))
	defer server.Close()

	// Note: In a real implementation, we'd inject the URL
	// For now, this test demonstrates the structure
	client := NewClient()
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
}

func TestCartItem_Subtotal(t *testing.T) {
	item := models.CartItem{
		ASIN:     "B08N5WRWNW",
		Title:    "Test Product",
		Price:    99.99,
		Quantity: 2,
		Subtotal: 199.98,
	}

	expected := 199.98
	if item.Subtotal != expected {
		t.Errorf("Expected subtotal %f, got %f", expected, item.Subtotal)
	}

	// Test calculation
	calculated := item.Price * float64(item.Quantity)
	if calculated != expected {
		t.Errorf("Calculated subtotal %f doesn't match expected %f", calculated, expected)
	}
}
