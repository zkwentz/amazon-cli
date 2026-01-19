package amazon

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestParseOrdersHTML_ReturnsCorrectCount(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Parse the HTML
	orders, err := parseOrdersHTML(fixtureData)
	if err != nil {
		t.Fatalf("parseOrdersHTML failed: %v", err)
	}

	// Verify we got 3 orders
	expectedCount := 3
	if len(orders) != expectedCount {
		t.Errorf("Expected %d orders, got %d", expectedCount, len(orders))
	}
}

func TestParseOrdersHTML_ExtractsAllFields(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Parse the HTML
	orders, err := parseOrdersHTML(fixtureData)
	if err != nil {
		t.Fatalf("parseOrdersHTML failed: %v", err)
	}

	// Verify all fields are extracted for each order
	for i, order := range orders {
		if order.OrderID == "" {
			t.Errorf("Order %d: OrderID is empty", i)
		}
		if order.Date == "" {
			t.Errorf("Order %d: Date is empty", i)
		}
		if order.Total == 0 {
			t.Errorf("Order %d: Total is 0", i)
		}
		if order.Status == "" {
			t.Errorf("Order %d: Status is empty", i)
		}
	}
}

func TestParseOrdersHTML_ExtractsOrderIDs(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Parse the HTML
	orders, err := parseOrdersHTML(fixtureData)
	if err != nil {
		t.Fatalf("parseOrdersHTML failed: %v", err)
	}

	// Expected order IDs from the fixture
	expectedOrderIDs := []string{
		"111-2222222-3333333",
		"111-4444444-5555555",
		"111-6666666-7777777",
	}

	// Verify order IDs match
	for i, expectedID := range expectedOrderIDs {
		if i >= len(orders) {
			t.Errorf("Expected order %d with ID %s, but only got %d orders", i, expectedID, len(orders))
			continue
		}
		if orders[i].OrderID != expectedID {
			t.Errorf("Order %d: expected OrderID %s, got %s", i, expectedID, orders[i].OrderID)
		}
	}
}

func TestParseOrdersHTML_ExtractsTotals(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Parse the HTML
	orders, err := parseOrdersHTML(fixtureData)
	if err != nil {
		t.Fatalf("parseOrdersHTML failed: %v", err)
	}

	// Expected totals from the fixture
	expectedTotals := []float64{29.99, 54.99, 149.99}

	// Verify totals match
	for i, expectedTotal := range expectedTotals {
		if i >= len(orders) {
			t.Errorf("Expected order %d with total %.2f, but only got %d orders", i, expectedTotal, len(orders))
			continue
		}
		if orders[i].Total != expectedTotal {
			t.Errorf("Order %d: expected Total %.2f, got %.2f", i, expectedTotal, orders[i].Total)
		}
	}
}

func TestParseOrdersHTML_ExtractsStatuses(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Parse the HTML
	orders, err := parseOrdersHTML(fixtureData)
	if err != nil {
		t.Fatalf("parseOrdersHTML failed: %v", err)
	}

	// Expected statuses from the fixture
	expectedStatuses := []string{"delivered", "pending", "cancelled"}

	// Verify statuses match
	for i, expectedStatus := range expectedStatuses {
		if i >= len(orders) {
			t.Errorf("Expected order %d with status %s, but only got %d orders", i, expectedStatus, len(orders))
			continue
		}
		if orders[i].Status != expectedStatus {
			t.Errorf("Order %d: expected Status %s, got %s", i, expectedStatus, orders[i].Status)
		}
	}
}

func TestParseOrdersHTML_EmptyHTML(t *testing.T) {
	html := []byte(`<html><body><div id="ordersContainer"></div></body></html>`)

	orders, err := parseOrdersHTML(html)
	if err != nil {
		t.Fatalf("parseOrdersHTML failed: %v", err)
	}

	if len(orders) != 0 {
		t.Errorf("Expected 0 orders for empty HTML, got %d", len(orders))
	}
}

func TestParseOrdersHTML_InvalidHTML(t *testing.T) {
	html := []byte(`not valid html`)

	// This should still parse without error (goquery is lenient)
	// but return no orders
	orders, err := parseOrdersHTML(html)
	if err != nil {
		t.Fatalf("parseOrdersHTML failed: %v", err)
	}

	if len(orders) != 0 {
		t.Errorf("Expected 0 orders for invalid HTML, got %d", len(orders))
	}
}

func TestParseOrdersHTML_PartialData(t *testing.T) {
	// HTML with an order that has only some fields
	html := []byte(`
		<html>
		<body>
			<div class="order" data-order-id="111-2222222-3333333">
				<div class="order-header">
					<div class="order-info">
						<span class="order-date">January 15, 2026</span>
					</div>
				</div>
			</div>
		</body>
		</html>
	`)

	orders, err := parseOrdersHTML(html)
	if err != nil {
		t.Fatalf("parseOrdersHTML failed: %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(orders))
	}

	// Verify the order ID was extracted
	if orders[0].OrderID != "111-2222222-3333333" {
		t.Errorf("Expected OrderID 111-2222222-3333333, got %s", orders[0].OrderID)
	}

	// Verify the date was extracted
	if orders[0].Date == "" {
		t.Errorf("Expected Date to be extracted")
	}

	// Verify missing fields are zero values
	if orders[0].Total != 0 {
		t.Errorf("Expected Total to be 0 for missing data, got %.2f", orders[0].Total)
	}
}

func TestParsePrice(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"$29.99", 29.99},
		{"$1,234.56", 1234.56},
		{"54.99", 54.99},
		{"$149.99", 149.99},
		{"$0.99", 0.99},
		{"$1,000.00", 1000.00},
		{"invalid", 0.0},
		{"", 0.0},
		{"$", 0.0},
	}

	for _, tt := range tests {
		result := parsePrice(tt.input)
		if result != tt.expected {
			t.Errorf("parsePrice(%q) = %.2f, expected %.2f", tt.input, result, tt.expected)
		}
	}
}

func TestGetOrders_Success(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path
		if r.URL.Path != "/gp/your-account/order-history" {
			t.Errorf("Expected path /gp/your-account/order-history, got %s", r.URL.Path)
		}

		// Return the fixture data
		w.WriteHeader(http.StatusOK)
		w.Write(fixtureData)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrders
	response, err := client.GetOrders(10, "")
	if err != nil {
		t.Fatalf("GetOrders() error = %v", err)
	}

	// Verify response
	if response == nil {
		t.Fatal("Expected non-nil response")
	}

	if len(response.Orders) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(response.Orders))
	}

	if response.TotalCount != 3 {
		t.Errorf("Expected TotalCount 3, got %d", response.TotalCount)
	}
}

func TestGetOrders_WithLimit(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(fixtureData)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrders with limit of 2
	response, err := client.GetOrders(2, "")
	if err != nil {
		t.Fatalf("GetOrders() error = %v", err)
	}

	// Verify response
	if len(response.Orders) != 2 {
		t.Errorf("Expected 2 orders with limit, got %d", len(response.Orders))
	}

	if response.TotalCount != 2 {
		t.Errorf("Expected TotalCount 2, got %d", response.TotalCount)
	}
}

func TestGetOrders_WithStatusFilter(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(fixtureData)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrders with status filter
	response, err := client.GetOrders(10, "delivered")
	if err != nil {
		t.Fatalf("GetOrders() error = %v", err)
	}

	// Verify response
	if len(response.Orders) != 1 {
		t.Errorf("Expected 1 delivered order, got %d", len(response.Orders))
	}

	// Verify all returned orders have the correct status
	for _, order := range response.Orders {
		if order.Status != "delivered" {
			t.Errorf("Expected status 'delivered', got %s", order.Status)
		}
	}
}

func TestGetOrders_DefaultLimit(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_list_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(fixtureData)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrders with zero limit (should default to 10)
	response, err := client.GetOrders(0, "")
	if err != nil {
		t.Fatalf("GetOrders() error = %v", err)
	}

	// Verify response (fixture has 3 orders, all should be returned)
	if len(response.Orders) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(response.Orders))
	}
}

func TestGetOrders_HTTPError(t *testing.T) {
	// Create a mock HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrders
	_, err := client.GetOrders(10, "")
	if err == nil {
		t.Fatal("Expected error for HTTP 500, got nil")
	}
}

func TestGetOrders_InvalidHTML(t *testing.T) {
	// Create a mock HTTP server that returns invalid HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>No orders here</body></html>"))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrders
	response, err := client.GetOrders(10, "")
	if err != nil {
		t.Fatalf("GetOrders() error = %v", err)
	}

	// Should return empty orders list
	if len(response.Orders) != 0 {
		t.Errorf("Expected 0 orders for invalid HTML, got %d", len(response.Orders))
	}
}

func TestGetOrder_EmptyOrderID(t *testing.T) {
	client := NewClient()

	_, err := client.GetOrder("")
	if err == nil {
		t.Fatal("Expected error for empty order ID, got nil")
	}

	if err.Error() != "order ID cannot be empty" {
		t.Errorf("Expected 'order ID cannot be empty' error, got: %v", err)
	}
}

func TestGetOrder_InvalidOrderIDFormat(t *testing.T) {
	client := NewClient()

	tests := []struct {
		orderID string
		desc    string
	}{
		{"123", "too short"},
		{"123-456-789", "wrong segment lengths"},
		{"abc-1234567-1234567", "letters in first segment"},
		{"123-abcdefg-1234567", "letters in second segment"},
		{"123-1234567-abcdefg", "letters in third segment"},
		{"1234-1234567-1234567", "first segment too long"},
		{"123-12345678-1234567", "second segment too long"},
		{"123-1234567-12345678", "third segment too long"},
	}

	for _, tt := range tests {
		_, err := client.GetOrder(tt.orderID)
		if err == nil {
			t.Errorf("Expected error for invalid order ID (%s): %s, got nil", tt.desc, tt.orderID)
		}
	}
}

func TestGetOrder_Success(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "order_detail_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path
		if r.URL.Path != "/gp/your-account/order-details" {
			t.Errorf("Expected path /gp/your-account/order-details, got %s", r.URL.Path)
		}

		// Verify the orderID query parameter
		orderID := r.URL.Query().Get("orderID")
		if orderID != "123-4567890-1234567" {
			t.Errorf("Expected orderID 123-4567890-1234567, got %s", orderID)
		}

		// Return the fixture data
		w.WriteHeader(http.StatusOK)
		w.Write(fixtureData)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrder
	order, err := client.GetOrder("123-4567890-1234567")
	if err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}

	// Verify response
	if order == nil {
		t.Fatal("Expected non-nil order")
	}

	if order.OrderID != "123-4567890-1234567" {
		t.Errorf("Expected OrderID 123-4567890-1234567, got %s", order.OrderID)
	}

	if order.Date != "2026-01-15" {
		t.Errorf("Expected Date 2026-01-15, got %s", order.Date)
	}

	if order.Total != 84.98 {
		t.Errorf("Expected Total 84.98, got %.2f", order.Total)
	}

	if order.Status != "delivered" {
		t.Errorf("Expected Status delivered, got %s", order.Status)
	}

	if len(order.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(order.Items))
	}

	// Verify tracking information
	if order.Tracking == nil {
		t.Fatal("Expected Tracking to be set")
	}

	if order.Tracking.Carrier != "UPS" {
		t.Errorf("Expected Carrier UPS, got %s", order.Tracking.Carrier)
	}

	if order.Tracking.TrackingNumber != "1Z999AA10123456784" {
		t.Errorf("Expected TrackingNumber 1Z999AA10123456784, got %s", order.Tracking.TrackingNumber)
	}
}

func TestGetOrder_HTTPError(t *testing.T) {
	// Create a mock HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrder
	_, err := client.GetOrder("111-2222222-3333333")
	if err == nil {
		t.Fatal("Expected error for HTTP 500, got nil")
	}
}

func TestGetOrder_CAPTCHADetection(t *testing.T) {
	// Create a mock HTTP server that returns a CAPTCHA page
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><div>Sorry, we just need to make sure you're not a robot</div></body></html>`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrder
	_, err := client.GetOrder("111-2222222-3333333")
	if err == nil {
		t.Fatal("Expected error for CAPTCHA, got nil")
	}

	if err.Error() != "CAPTCHA detected - Amazon is blocking automated access" {
		t.Errorf("Expected CAPTCHA error, got: %v", err)
	}
}

func TestGetOrder_ParseError(t *testing.T) {
	// Create a mock HTTP server that returns HTML without an order ID
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><div>Order not found</div></body></html>`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrder
	_, err := client.GetOrder("111-2222222-3333333")
	if err == nil {
		t.Fatal("Expected error for missing order ID, got nil")
	}
}
