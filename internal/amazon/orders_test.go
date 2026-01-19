package amazon

import (
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
