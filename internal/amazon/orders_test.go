package amazon

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetOrders(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name           string
		limit          int
		status         string
		expectedCount  int
		expectedStatus string
	}{
		{
			name:          "get all orders",
			limit:         0,
			status:        "",
			expectedCount: 2,
		},
		{
			name:          "get one order",
			limit:         1,
			status:        "",
			expectedCount: 1,
		},
		{
			name:           "filter by delivered status",
			limit:          0,
			status:         "delivered",
			expectedCount:  1,
			expectedStatus: "delivered",
		},
		{
			name:           "filter by pending status",
			limit:          0,
			status:         "pending",
			expectedCount:  1,
			expectedStatus: "pending",
		},
		{
			name:          "filter by non-existent status",
			limit:         0,
			status:        "cancelled",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.GetOrders(tt.limit, tt.status)
			if err != nil {
				t.Fatalf("GetOrders() error = %v", err)
			}

			if response.TotalCount != tt.expectedCount {
				t.Errorf("GetOrders() got %d orders, want %d", response.TotalCount, tt.expectedCount)
			}

			if len(response.Orders) != tt.expectedCount {
				t.Errorf("GetOrders() got %d orders in slice, want %d", len(response.Orders), tt.expectedCount)
			}

			if tt.expectedStatus != "" && len(response.Orders) > 0 {
				for _, order := range response.Orders {
					if order.Status != tt.expectedStatus {
						t.Errorf("GetOrders() order status = %s, want %s", order.Status, tt.expectedStatus)
					}
				}
			}
		})
	}
}

func TestGetOrdersStructure(t *testing.T) {
	client := NewClient()
	response, err := client.GetOrders(0, "")
	if err != nil {
		t.Fatalf("GetOrders() error = %v", err)
	}

	if len(response.Orders) == 0 {
		t.Fatal("GetOrders() returned no orders")
	}

	// Test first order structure
	order := response.Orders[0]

	if order.OrderID == "" {
		t.Error("Order OrderID is empty")
	}

	if order.Date == "" {
		t.Error("Order Date is empty")
	}

	if order.Total <= 0 {
		t.Error("Order Total should be greater than 0")
	}

	if order.Status == "" {
		t.Error("Order Status is empty")
	}

	if len(order.Items) == 0 {
		t.Error("Order has no items")
	}

	// Test first item structure
	item := order.Items[0]
	if item.ASIN == "" {
		t.Error("OrderItem ASIN is empty")
	}

	if item.Title == "" {
		t.Error("OrderItem Title is empty")
	}

	if item.Quantity <= 0 {
		t.Error("OrderItem Quantity should be greater than 0")
	}

	if item.Price <= 0 {
		t.Error("OrderItem Price should be greater than 0")
	}

	// Test tracking if present
	if order.Tracking != nil {
		if order.Tracking.Carrier == "" {
			t.Error("Tracking Carrier is empty")
		}

		if order.Tracking.TrackingNumber == "" {
			t.Error("Tracking TrackingNumber is empty")
		}

		if order.Tracking.Status == "" {
			t.Error("Tracking Status is empty")
		}
	}
}

func TestGetOrder(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		orderID     string
		wantError   bool
		errorCode   string
		wantOrderID string
	}{
		{
			name:        "get existing order - delivered",
			orderID:     "123-4567890-1234567",
			wantError:   false,
			wantOrderID: "123-4567890-1234567",
		},
		{
			name:        "get existing order - pending",
			orderID:     "123-4567890-1234568",
			wantError:   false,
			wantOrderID: "123-4567890-1234568",
		},
		{
			name:      "get non-existent order",
			orderID:   "999-9999999-9999999",
			wantError: true,
			errorCode: models.ErrorCodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := client.GetOrder(tt.orderID)

			if tt.wantError {
				if err == nil {
					t.Error("GetOrder() expected error but got nil")
				}
				if cliErr, ok := err.(*models.CLIError); ok {
					if cliErr.Code != tt.errorCode {
						t.Errorf("GetOrder() error code = %s, want %s", cliErr.Code, tt.errorCode)
					}
				} else {
					t.Error("GetOrder() error is not a CLIError")
				}
			} else {
				if err != nil {
					t.Errorf("GetOrder() unexpected error = %v", err)
				}
				if order == nil {
					t.Fatal("GetOrder() returned nil order")
				}
				if order.OrderID != tt.wantOrderID {
					t.Errorf("GetOrder() order ID = %s, want %s", order.OrderID, tt.wantOrderID)
				}
			}
		})
	}
}

func TestGetOrderStructure(t *testing.T) {
	client := NewClient()
	order, err := client.GetOrder("123-4567890-1234567")
	if err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}

	if order == nil {
		t.Fatal("GetOrder() returned nil order")
	}

	// Verify order structure
	if order.OrderID == "" {
		t.Error("Order OrderID is empty")
	}

	if order.Date == "" {
		t.Error("Order Date is empty")
	}

	if order.Total <= 0 {
		t.Error("Order Total should be greater than 0")
	}

	if order.Status == "" {
		t.Error("Order Status is empty")
	}

	if len(order.Items) == 0 {
		t.Error("Order has no items")
	}

	// Test first item structure
	item := order.Items[0]
	if item.ASIN == "" {
		t.Error("OrderItem ASIN is empty")
	}

	if item.Title == "" {
		t.Error("OrderItem Title is empty")
	}

	if item.Quantity <= 0 {
		t.Error("OrderItem Quantity should be greater than 0")
	}

	if item.Price <= 0 {
		t.Error("OrderItem Price should be greater than 0")
	}

	// Test tracking if present (should be present for delivered order)
	if order.Tracking != nil {
		if order.Tracking.Carrier == "" {
			t.Error("Tracking Carrier is empty")
		}

		if order.Tracking.TrackingNumber == "" {
			t.Error("Tracking TrackingNumber is empty")
		}

		if order.Tracking.Status == "" {
			t.Error("Tracking Status is empty")
		}

		if order.Tracking.DeliveryDate == "" {
			t.Error("Tracking DeliveryDate is empty")
		}
	}
}

func TestGetOrderTracking(t *testing.T) {
	client := NewClient()

	_, err := client.GetOrderTracking("123-4567890-1234567")
	if err == nil {
		t.Error("GetOrderTracking() should return error for unimplemented feature")
	}

	if cliErr, ok := err.(*models.CLIError); ok {
		if cliErr.Code != models.ErrorCodeNotFound {
			t.Errorf("GetOrderTracking() error code = %s, want %s", cliErr.Code, models.ErrorCodeNotFound)
		}
	}
}

func TestGetOrderHistory(t *testing.T) {
	client := NewClient()

	_, err := client.GetOrderHistory(2024)
	if err == nil {
		t.Error("GetOrderHistory() should return error for unimplemented feature")
	}

	if cliErr, ok := err.(*models.CLIError); ok {
		if cliErr.Code != models.ErrorCodeNotFound {
			t.Errorf("GetOrderHistory() error code = %s, want %s", cliErr.Code, models.ErrorCodeNotFound)
		}
	}
}
