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

func TestParseTrackingHTML_ValidHTML(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "tracking_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Parse the HTML
	tracking, err := parseTrackingHTML(fixtureData)
	if err != nil {
		t.Fatalf("parseTrackingHTML failed: %v", err)
	}

	// Verify carrier
	expectedCarrier := "UPS"
	if tracking.Carrier != expectedCarrier {
		t.Errorf("Expected Carrier %q, got %q", expectedCarrier, tracking.Carrier)
	}

	// Verify tracking number
	expectedTrackingNumber := "1Z999AA10123456784"
	if tracking.TrackingNumber != expectedTrackingNumber {
		t.Errorf("Expected TrackingNumber %q, got %q", expectedTrackingNumber, tracking.TrackingNumber)
	}

	// Verify status
	expectedStatus := "in transit"
	if tracking.Status != expectedStatus {
		t.Errorf("Expected Status %q, got %q", expectedStatus, tracking.Status)
	}

	// Verify delivery date
	expectedDeliveryDate := "2026-01-20"
	if tracking.DeliveryDate != expectedDeliveryDate {
		t.Errorf("Expected DeliveryDate %q, got %q", expectedDeliveryDate, tracking.DeliveryDate)
	}

	// Verify events
	if len(tracking.Events) != 3 {
		t.Errorf("Expected 3 tracking events, got %d", len(tracking.Events))
	}

	// Verify first event
	if len(tracking.Events) > 0 {
		firstEvent := tracking.Events[0]
		if firstEvent.Status != "Out for delivery" {
			t.Errorf("Expected first event status %q, got %q", "Out for delivery", firstEvent.Status)
		}
		if firstEvent.Location != "Local Distribution Center - Seattle, WA" {
			t.Errorf("Expected first event location %q, got %q", "Local Distribution Center - Seattle, WA", firstEvent.Location)
		}
		if firstEvent.Timestamp == "" {
			t.Error("Expected first event timestamp to be non-empty")
		}
	}
}

func TestParseTrackingHTML_MinimalData(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<div class="tracking-section">
				<div class="tracking-carrier">
					<span class="value">USPS</span>
				</div>
				<div class="tracking-number">
					<span class="value">9400111899561543123456</span>
				</div>
			</div>
		</body>
		</html>
	`)

	tracking, err := parseTrackingHTML(html)
	if err != nil {
		t.Fatalf("parseTrackingHTML failed: %v", err)
	}

	if tracking.Carrier != "USPS" {
		t.Errorf("Expected Carrier USPS, got %q", tracking.Carrier)
	}

	if tracking.TrackingNumber != "9400111899561543123456" {
		t.Errorf("Expected TrackingNumber 9400111899561543123456, got %q", tracking.TrackingNumber)
	}

	if tracking.Status != "" {
		t.Errorf("Expected empty Status, got %q", tracking.Status)
	}

	if tracking.DeliveryDate != "" {
		t.Errorf("Expected empty DeliveryDate, got %q", tracking.DeliveryDate)
	}
}

func TestParseTrackingHTML_NoTrackingInfo(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<div class="order-details">
				<h1>Order Details</h1>
			</div>
		</body>
		</html>
	`)

	_, err := parseTrackingHTML(html)
	if err == nil {
		t.Error("Expected error for HTML without tracking info, got nil")
	}
}

func TestParseTrackingHTML_WithEvents(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<div class="tracking-section">
				<div class="tracking-carrier">
					<span class="value">FedEx</span>
				</div>
				<div class="tracking-number">
					<span class="value">123456789012</span>
				</div>
				<div class="tracking-status">
					<span class="value">Delivered</span>
				</div>
			</div>
			<div class="tracking-events">
				<div class="event">
					<div class="event-timestamp">January 18, 2026 3:45 PM</div>
					<div class="event-location">Seattle, WA</div>
					<div class="event-status">Delivered</div>
				</div>
				<div class="event">
					<div class="event-timestamp">January 18, 2026 8:00 AM</div>
					<div class="event-location">Portland, OR</div>
					<div class="event-status">In transit</div>
				</div>
			</div>
		</body>
		</html>
	`)

	tracking, err := parseTrackingHTML(html)
	if err != nil {
		t.Fatalf("parseTrackingHTML failed: %v", err)
	}

	if len(tracking.Events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(tracking.Events))
	}

	// Check first event
	if tracking.Events[0].Status != "Delivered" {
		t.Errorf("Expected first event status Delivered, got %q", tracking.Events[0].Status)
	}
	if tracking.Events[0].Location != "Seattle, WA" {
		t.Errorf("Expected first event location Seattle, WA, got %q", tracking.Events[0].Location)
	}

	// Check second event
	if tracking.Events[1].Status != "In transit" {
		t.Errorf("Expected second event status In transit, got %q", tracking.Events[1].Status)
	}
}

func TestParseTrackingHTML_DateParsing(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<div class="tracking-section">
				<div class="tracking-carrier">
					<span class="value">Amazon Logistics</span>
				</div>
				<div class="tracking-number">
					<span class="value">TBA123456789</span>
				</div>
				<div class="delivery-date">
					<span class="value">January 25, 2026</span>
				</div>
			</div>
		</body>
		</html>
	`)

	tracking, err := parseTrackingHTML(html)
	if err != nil {
		t.Fatalf("parseTrackingHTML failed: %v", err)
	}

	expectedDate := "2026-01-25"
	if tracking.DeliveryDate != expectedDate {
		t.Errorf("Expected DeliveryDate %q, got %q", expectedDate, tracking.DeliveryDate)
	}
}

func TestParseTrackingHTML_StatusNormalization(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<div class="tracking-section">
				<div class="tracking-carrier">
					<span class="value">DHL</span>
				</div>
				<div class="tracking-number">
					<span class="value">1234567890</span>
				</div>
				<div class="tracking-status">
					<span class="value">  OUT FOR DELIVERY  </span>
				</div>
			</div>
		</body>
		</html>
	`)

	tracking, err := parseTrackingHTML(html)
	if err != nil {
		t.Fatalf("parseTrackingHTML failed: %v", err)
	}

	expectedStatus := "out for delivery"
	if tracking.Status != expectedStatus {
		t.Errorf("Expected Status %q (normalized), got %q", expectedStatus, tracking.Status)
	}
}

func TestParseTrackingHTML_EmptyEvents(t *testing.T) {
	html := []byte(`
		<html>
		<body>
			<div class="tracking-section">
				<div class="tracking-carrier">
					<span class="value">UPS</span>
				</div>
				<div class="tracking-number">
					<span class="value">1Z999AA10123456784</span>
				</div>
			</div>
			<div class="tracking-events">
				<div class="event">
					<div class="event-timestamp"></div>
					<div class="event-location"></div>
					<div class="event-status"></div>
				</div>
			</div>
		</body>
		</html>
	`)

	tracking, err := parseTrackingHTML(html)
	if err != nil {
		t.Fatalf("parseTrackingHTML failed: %v", err)
	}

	// Empty event should not be added
	if len(tracking.Events) != 0 {
		t.Errorf("Expected 0 events (empty events should be filtered), got %d", len(tracking.Events))
	}
}

func TestGetOrderTracking_EmptyOrderID(t *testing.T) {
	client := NewClient()

	_, err := client.GetOrderTracking("")
	if err == nil {
		t.Fatal("Expected error for empty order ID, got nil")
	}

	if err.Error() != "order ID cannot be empty" {
		t.Errorf("Expected 'order ID cannot be empty' error, got: %v", err)
	}
}

func TestGetOrderTracking_Success(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "tracking_sample.html"))
	if err != nil {
		t.Fatalf("Failed to read fixture file: %v", err)
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path
		if r.URL.Path != "/progress-tracker/package/ref=ppx_yo_dt_b_track_package" {
			t.Errorf("Expected path /progress-tracker/package/ref=ppx_yo_dt_b_track_package, got %s", r.URL.Path)
		}

		// Verify the orderId query parameter
		orderID := r.URL.Query().Get("orderId")
		if orderID != "111-2222222-3333333" {
			t.Errorf("Expected orderId 111-2222222-3333333, got %s", orderID)
		}

		// Return the fixture data
		w.WriteHeader(http.StatusOK)
		w.Write(fixtureData)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrderTracking
	tracking, err := client.GetOrderTracking("111-2222222-3333333")
	if err != nil {
		t.Fatalf("GetOrderTracking() error = %v", err)
	}

	// Verify response
	if tracking == nil {
		t.Fatal("Expected non-nil tracking")
	}

	if tracking.Carrier != "UPS" {
		t.Errorf("Expected Carrier UPS, got %s", tracking.Carrier)
	}

	if tracking.TrackingNumber != "1Z999AA10123456784" {
		t.Errorf("Expected TrackingNumber 1Z999AA10123456784, got %s", tracking.TrackingNumber)
	}

	if tracking.Status != "in transit" {
		t.Errorf("Expected Status 'in transit', got %s", tracking.Status)
	}

	if tracking.DeliveryDate != "2026-01-20" {
		t.Errorf("Expected DeliveryDate 2026-01-20, got %s", tracking.DeliveryDate)
	}

	if len(tracking.Events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(tracking.Events))
	}
}

func TestGetOrderTracking_HTTPError(t *testing.T) {
	// Create a mock HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrderTracking
	_, err := client.GetOrderTracking("111-2222222-3333333")
	if err == nil {
		t.Fatal("Expected error for HTTP 500, got nil")
	}
}

func TestGetOrderTracking_CAPTCHADetection(t *testing.T) {
	// Create a mock HTTP server that returns a CAPTCHA page
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><div>Sorry, we just need to make sure you're not a robot</div></body></html>`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrderTracking
	_, err := client.GetOrderTracking("111-2222222-3333333")
	if err == nil {
		t.Fatal("Expected error for CAPTCHA, got nil")
	}

	expectedErr := "CAPTCHA detected - please complete CAPTCHA in browser and try again"
	if err.Error() != expectedErr {
		t.Errorf("Expected CAPTCHA error, got: %v", err)
	}
}

func TestGetOrderTracking_ParseError(t *testing.T) {
	// Create a mock HTTP server that returns HTML without tracking info
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><div>Tracking not available</div></body></html>`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient()
	client.baseURL = server.URL

	// Test GetOrderTracking
	_, err := client.GetOrderTracking("111-2222222-3333333")
	if err == nil {
		t.Fatal("Expected error for missing tracking info, got nil")
	}
}

func TestGetOrderTracking_ValidOrderIDFormat(t *testing.T) {
	// Load the fixture file
	fixtureData, err := os.ReadFile(filepath.Join("..", "..", "testdata", "orders", "tracking_sample.html"))
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

	// Test with valid order ID formats
	validOrderIDs := []string{
		"111-2222222-3333333",
		"999-8888888-7777777",
		"123-4567890-1234567",
	}

	for _, orderID := range validOrderIDs {
		tracking, err := client.GetOrderTracking(orderID)
		if err != nil {
			t.Errorf("GetOrderTracking(%s) unexpected error: %v", orderID, err)
		}
		if tracking == nil {
			t.Errorf("GetOrderTracking(%s) returned nil tracking", orderID)
		}
	}
}
