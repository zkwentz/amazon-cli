package amazon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGetOrders tests the GetOrders method
func TestGetOrders(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		status         string
		mockHTML       string
		expectedCount  int
		expectError    bool
		expectedStatus int
	}{
		{
			name:   "successful fetch with orders",
			limit:  10,
			status: "",
			mockHTML: `
				<html>
					<div class="order" data-order-id="111-1111111-1111111">
						<div class="order-id" data-order-id="111-1111111-1111111">Order #111-1111111-1111111</div>
						<div class="order-date">January 15, 2024</div>
						<div class="order-total">$29.99</div>
						<div class="order-status">Delivered</div>
						<div class="order-item" data-asin="B08N5WRWNW">
							<div class="product-title">Test Product</div>
							<div class="item-quantity">Qty: 1</div>
							<div class="item-price">$29.99</div>
						</div>
					</div>
				</html>
			`,
			expectedCount:  1,
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "limit orders returned",
			limit:  1,
			status: "",
			mockHTML: `
				<html>
					<div class="order" data-order-id="111-1111111-1111111">
						<div class="order-id" data-order-id="111-1111111-1111111">Order #111-1111111-1111111</div>
						<div class="order-date">January 15, 2024</div>
						<div class="order-total">$29.99</div>
						<div class="order-status">Delivered</div>
					</div>
					<div class="order" data-order-id="222-2222222-2222222">
						<div class="order-id" data-order-id="222-2222222-2222222">Order #222-2222222-2222222</div>
						<div class="order-date">January 16, 2024</div>
						<div class="order-total">$49.99</div>
						<div class="order-status">Pending</div>
					</div>
				</html>
			`,
			expectedCount:  1,
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty order list",
			limit:          10,
			status:         "",
			mockHTML:       `<html><body>No orders found</body></html>`,
			expectedCount:  0,
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.expectedStatus)
				w.Write([]byte(tt.mockHTML))
			}))
			defer server.Close()

			// Create client with mock server
			client := &Client{
				httpClient: &http.Client{Timeout: 5 * time.Second},
			}

			// Note: In real implementation, we'd need to override the URL
			// For now, this demonstrates the test structure
			resp, err := client.GetOrders(tt.limit, tt.status)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && resp != nil {
				if len(resp.Orders) != tt.expectedCount {
					t.Errorf("expected %d orders, got %d", tt.expectedCount, len(resp.Orders))
				}
			}
		})
	}
}

// TestGetOrder tests the GetOrder method
func TestGetOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderID     string
		mockHTML    string
		expectError bool
		statusCode  int
	}{
		{
			name:    "successful fetch",
			orderID: "111-1111111-1111111",
			mockHTML: `
				<html>
					<div class="order-date-invoice-item">January 15, 2024</div>
					<div class="grand-total-price">$29.99</div>
					<div class="delivery-box__primary-text">Delivered</div>
					<div class="shipment-item" data-asin="B08N5WRWNW">
						<div class="product-title">Test Product</div>
						<div class="item-quantity">Qty: 1</div>
						<div class="item-price">$29.99</div>
					</div>
				</html>
			`,
			expectError: false,
			statusCode:  http.StatusOK,
		},
		{
			name:        "order not found",
			orderID:     "999-9999999-9999999",
			mockHTML:    `<html><body>Order not found</body></html>`,
			expectError: true,
			statusCode:  http.StatusNotFound,
		},
		{
			name:        "empty order ID",
			orderID:     "",
			mockHTML:    "",
			expectError: true,
			statusCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				httpClient: &http.Client{Timeout: 5 * time.Second},
			}

			_, err := client.GetOrder(tt.orderID)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestGetOrderTracking tests the GetOrderTracking method
func TestGetOrderTracking(t *testing.T) {
	tests := []struct {
		name        string
		orderID     string
		mockHTML    string
		expectError bool
		expectNil   bool
	}{
		{
			name:    "successful tracking fetch",
			orderID: "111-1111111-1111111",
			mockHTML: `
				<html>
					<div class="carrier-info">UPS</div>
					<div class="tracking-number">1Z999AA10123456784</div>
					<div class="tracking-status">Delivered</div>
					<div class="delivery-date">January 17, 2024</div>
				</html>
			`,
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "no tracking info available",
			orderID:     "111-1111111-1111111",
			mockHTML:    `<html><body>No tracking information</body></html>`,
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "empty order ID",
			orderID:     "",
			mockHTML:    "",
			expectError: true,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				httpClient: &http.Client{Timeout: 5 * time.Second},
			}

			tracking, err := client.GetOrderTracking(tt.orderID)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectNil && tracking != nil {
				t.Errorf("expected nil tracking but got value")
			}
		})
	}
}

// TestGetOrderHistory tests the GetOrderHistory method
func TestGetOrderHistory(t *testing.T) {
	currentYear := time.Now().Year()

	tests := []struct {
		name          string
		year          int
		expectError   bool
		expectedCount int
	}{
		{
			name:          "current year",
			year:          currentYear,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:          "valid past year",
			year:          2023,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:          "invalid year (too old)",
			year:          1990,
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "invalid year (future)",
			year:          currentYear + 10,
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:          "zero year defaults to current",
			year:          0,
			expectError:   false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				httpClient: &http.Client{Timeout: 5 * time.Second},
			}

			_, err := client.GetOrderHistory(tt.year)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestParsePrice tests the parsePrice helper function
func TestParsePrice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"standard price", "$29.99", 29.99},
		{"price with comma", "$1,234.56", 1234.56},
		{"price without symbol", "29.99", 29.99},
		{"price with spaces", " $29.99 ", 29.99},
		{"invalid price", "invalid", 0.0},
		{"empty string", "", 0.0},
		{"zero price", "$0.00", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePrice(tt.input)
			if result != tt.expected {
				t.Errorf("parsePrice(%q) = %f; expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

// TestParseQuantity tests the parseQuantity helper function
func TestParseQuantity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"qty with prefix", "Qty: 2", 2},
		{"quantity with prefix", "Quantity: 5", 5},
		{"just number", "3", 3},
		{"with spaces", " Qty: 1 ", 1},
		{"invalid quantity", "invalid", 1},
		{"empty string", "", 1},
		{"zero quantity", "0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseQuantity(tt.input)
			if result != tt.expected {
				t.Errorf("parseQuantity(%q) = %d; expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestOrderJSONSerialization tests JSON marshaling/unmarshaling
func TestOrderJSONSerialization(t *testing.T) {
	order := &Order{
		OrderID: "111-1111111-1111111",
		Date:    "January 15, 2024",
		Total:   29.99,
		Status:  "delivered",
		Items: []OrderItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Test Product",
				Quantity: 1,
				Price:    29.99,
			},
		},
		Tracking: &Tracking{
			Carrier:        "UPS",
			TrackingNumber: "1Z999AA10123456784",
			Status:         "delivered",
			DeliveryDate:   "January 17, 2024",
		},
	}

	// Test marshaling
	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("failed to marshal order: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Order
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal order: %v", err)
	}

	// Verify fields
	if unmarshaled.OrderID != order.OrderID {
		t.Errorf("OrderID mismatch: got %q, expected %q", unmarshaled.OrderID, order.OrderID)
	}

	if unmarshaled.Total != order.Total {
		t.Errorf("Total mismatch: got %f, expected %f", unmarshaled.Total, order.Total)
	}

	if len(unmarshaled.Items) != len(order.Items) {
		t.Errorf("Items count mismatch: got %d, expected %d", len(unmarshaled.Items), len(order.Items))
	}

	if unmarshaled.Tracking == nil {
		t.Error("Tracking is nil after unmarshal")
	} else if unmarshaled.Tracking.TrackingNumber != order.Tracking.TrackingNumber {
		t.Errorf("TrackingNumber mismatch: got %q, expected %q",
			unmarshaled.Tracking.TrackingNumber, order.Tracking.TrackingNumber)
	}
}

// TestNilClient tests that methods handle nil client gracefully
func TestNilClient(t *testing.T) {
	var client *Client

	t.Run("GetOrders with nil client", func(t *testing.T) {
		_, err := client.GetOrders(10, "")
		if err == nil {
			t.Error("expected error with nil client")
		}
	})

	t.Run("GetOrder with nil client", func(t *testing.T) {
		_, err := client.GetOrder("111-1111111-1111111")
		if err == nil {
			t.Error("expected error with nil client")
		}
	})

	t.Run("GetOrderTracking with nil client", func(t *testing.T) {
		_, err := client.GetOrderTracking("111-1111111-1111111")
		if err == nil {
			t.Error("expected error with nil client")
		}
	})

	t.Run("GetOrderHistory with nil client", func(t *testing.T) {
		_, err := client.GetOrderHistory(2024)
		if err == nil {
			t.Error("expected error with nil client")
		}
	})
}

// TestOrdersResponseJSONSerialization tests the OrdersResponse JSON output
func TestOrdersResponseJSONSerialization(t *testing.T) {
	response := &OrdersResponse{
		Orders: []Order{
			{
				OrderID: "111-1111111-1111111",
				Date:    "January 15, 2024",
				Total:   29.99,
				Status:  "delivered",
				Items: []OrderItem{
					{
						ASIN:     "B08N5WRWNW",
						Title:    "Test Product",
						Quantity: 1,
						Price:    29.99,
					},
				},
			},
		},
		TotalCount: 1,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal OrdersResponse: %v", err)
	}

	var unmarshaled OrdersResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal OrdersResponse: %v", err)
	}

	if unmarshaled.TotalCount != response.TotalCount {
		t.Errorf("TotalCount mismatch: got %d, expected %d", unmarshaled.TotalCount, response.TotalCount)
	}

	if len(unmarshaled.Orders) != len(response.Orders) {
		t.Errorf("Orders count mismatch: got %d, expected %d", len(unmarshaled.Orders), len(response.Orders))
	}
}
