package amazon

import (
	"encoding/json"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestParseOrderFromJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *models.Order
		wantErr bool
	}{
		{
			name: "valid complete order",
			input: `{
				"order_id": "123-4567890-1234567",
				"date": "2024-01-15",
				"total": 29.99,
				"status": "delivered",
				"items": [
					{
						"asin": "B08N5WRWNW",
						"title": "Product Name",
						"quantity": 1,
						"price": 29.99
					}
				],
				"tracking": {
					"carrier": "UPS",
					"tracking_number": "1Z999AA10123456784",
					"status": "delivered",
					"delivery_date": "2024-01-17"
				}
			}`,
			want: &models.Order{
				OrderID: "123-4567890-1234567",
				Date:    "2024-01-15",
				Total:   29.99,
				Status:  "delivered",
				Items: []models.OrderItem{
					{
						ASIN:     "B08N5WRWNW",
						Title:    "Product Name",
						Quantity: 1,
						Price:    29.99,
					},
				},
				Tracking: &models.Tracking{
					Carrier:        "UPS",
					TrackingNumber: "1Z999AA10123456784",
					Status:         "delivered",
					DeliveryDate:   "2024-01-17",
				},
			},
			wantErr: false,
		},
		{
			name: "order without tracking",
			input: `{
				"order_id": "111-2222222-3333333",
				"date": "2024-02-01",
				"total": 49.99,
				"status": "pending",
				"items": [
					{
						"asin": "B07XYZ1234",
						"title": "Another Product",
						"quantity": 2,
						"price": 24.995
					}
				]
			}`,
			want: &models.Order{
				OrderID: "111-2222222-3333333",
				Date:    "2024-02-01",
				Total:   49.99,
				Status:  "pending",
				Items: []models.OrderItem{
					{
						ASIN:     "B07XYZ1234",
						Title:    "Another Product",
						Quantity: 2,
						Price:    24.995,
					},
				},
				Tracking: nil,
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"order_id": "invalid json`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty JSON",
			input:   ``,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderFromJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.OrderID != tt.want.OrderID {
					t.Errorf("OrderID = %v, want %v", got.OrderID, tt.want.OrderID)
				}
				if got.Date != tt.want.Date {
					t.Errorf("Date = %v, want %v", got.Date, tt.want.Date)
				}
				if got.Total != tt.want.Total {
					t.Errorf("Total = %v, want %v", got.Total, tt.want.Total)
				}
				if got.Status != tt.want.Status {
					t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
				}
				if len(got.Items) != len(tt.want.Items) {
					t.Errorf("Items length = %v, want %v", len(got.Items), len(tt.want.Items))
				}
				if (got.Tracking == nil) != (tt.want.Tracking == nil) {
					t.Errorf("Tracking presence mismatch")
				}
			}
		})
	}
}

func TestParseOrdersResponse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *models.OrdersResponse
		wantErr bool
	}{
		{
			name: "valid orders response",
			input: `{
				"orders": [
					{
						"order_id": "123-4567890-1234567",
						"date": "2024-01-15",
						"total": 29.99,
						"status": "delivered",
						"items": []
					}
				],
				"total_count": 1
			}`,
			want: &models.OrdersResponse{
				Orders: []models.Order{
					{
						OrderID: "123-4567890-1234567",
						Date:    "2024-01-15",
						Total:   29.99,
						Status:  "delivered",
						Items:   []models.OrderItem{},
					},
				},
				TotalCount: 1,
			},
			wantErr: false,
		},
		{
			name: "empty orders list",
			input: `{
				"orders": [],
				"total_count": 0
			}`,
			want: &models.OrdersResponse{
				Orders:     []models.Order{},
				TotalCount: 0,
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"orders": [}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrdersResponse([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrdersResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.TotalCount != tt.want.TotalCount {
					t.Errorf("TotalCount = %v, want %v", got.TotalCount, tt.want.TotalCount)
				}
				if len(got.Orders) != len(tt.want.Orders) {
					t.Errorf("Orders length = %v, want %v", len(got.Orders), len(tt.want.Orders))
				}
			}
		})
	}
}

func TestParseOrderFromHTML(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    *models.Order
		wantErr bool
	}{
		{
			name: "valid HTML with all fields",
			html: `Order ID: 123-4567890-1234567
Date: 2024-01-15
Total: $29.99
Status: delivered`,
			want: &models.Order{
				OrderID: "123-4567890-1234567",
				Date:    "2024-01-15",
				Total:   29.99,
				Status:  "delivered",
				Items:   []models.OrderItem{},
			},
			wantErr: false,
		},
		{
			name: "HTML with comma in total",
			html: `Order ID: 111-2222222-3333333
Date: 2024-02-01
Total: $1,234.56
Status: pending`,
			want: &models.Order{
				OrderID: "111-2222222-3333333",
				Date:    "2024-02-01",
				Total:   1234.56,
				Status:  "pending",
				Items:   []models.OrderItem{},
			},
			wantErr: false,
		},
		{
			name: "HTML missing order ID",
			html: `Date: 2024-01-15
Total: $29.99
Status: delivered`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty HTML",
			html:    "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderFromHTML(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderFromHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.OrderID != tt.want.OrderID {
					t.Errorf("OrderID = %v, want %v", got.OrderID, tt.want.OrderID)
				}
				if got.Date != tt.want.Date {
					t.Errorf("Date = %v, want %v", got.Date, tt.want.Date)
				}
				if got.Total != tt.want.Total {
					t.Errorf("Total = %v, want %v", got.Total, tt.want.Total)
				}
				if got.Status != tt.want.Status {
					t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
				}
			}
		})
	}
}

func TestValidateOrder(t *testing.T) {
	tests := []struct {
		name    string
		order   *models.Order
		wantErr bool
	}{
		{
			name: "valid order",
			order: &models.Order{
				OrderID: "123-4567890-1234567",
				Date:    "2024-01-15",
				Total:   29.99,
				Status:  "delivered",
				Items:   []models.OrderItem{},
			},
			wantErr: false,
		},
		{
			name:    "nil order",
			order:   nil,
			wantErr: true,
		},
		{
			name: "missing order ID",
			order: &models.Order{
				Date:   "2024-01-15",
				Total:  29.99,
				Status: "delivered",
				Items:  []models.OrderItem{},
			},
			wantErr: true,
		},
		{
			name: "invalid order ID format",
			order: &models.Order{
				OrderID: "invalid-id",
				Date:    "2024-01-15",
				Total:   29.99,
				Status:  "delivered",
				Items:   []models.OrderItem{},
			},
			wantErr: true,
		},
		{
			name: "missing date",
			order: &models.Order{
				OrderID: "123-4567890-1234567",
				Total:   29.99,
				Status:  "delivered",
				Items:   []models.OrderItem{},
			},
			wantErr: true,
		},
		{
			name: "negative total",
			order: &models.Order{
				OrderID: "123-4567890-1234567",
				Date:    "2024-01-15",
				Total:   -10.00,
				Status:  "delivered",
				Items:   []models.OrderItem{},
			},
			wantErr: true,
		},
		{
			name: "missing status",
			order: &models.Order{
				OrderID: "123-4567890-1234567",
				Date:    "2024-01-15",
				Total:   29.99,
				Items:   []models.OrderItem{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrder(tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidOrderID(t *testing.T) {
	tests := []struct {
		name    string
		orderID string
		want    bool
	}{
		{
			name:    "valid order ID",
			orderID: "123-4567890-1234567",
			want:    true,
		},
		{
			name:    "another valid order ID",
			orderID: "111-2222222-3333333",
			want:    true,
		},
		{
			name:    "too short",
			orderID: "123-456-789",
			want:    false,
		},
		{
			name:    "missing dashes",
			orderID: "12345678901234567",
			want:    false,
		},
		{
			name:    "wrong format - first part too long",
			orderID: "1234-567890-1234567",
			want:    false,
		},
		{
			name:    "wrong format - middle part too short",
			orderID: "123-456789-1234567",
			want:    false,
		},
		{
			name:    "empty string",
			orderID: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidOrderID(tt.orderID)
			if got != tt.want {
				t.Errorf("IsValidOrderID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterOrdersByStatus(t *testing.T) {
	orders := []models.Order{
		{OrderID: "123-4567890-1234567", Status: "delivered"},
		{OrderID: "111-2222222-3333333", Status: "pending"},
		{OrderID: "444-5555555-6666666", Status: "delivered"},
		{OrderID: "777-8888888-9999999", Status: "returned"},
	}

	tests := []struct {
		name   string
		orders []models.Order
		status string
		want   int
	}{
		{
			name:   "filter delivered",
			orders: orders,
			status: "delivered",
			want:   2,
		},
		{
			name:   "filter pending",
			orders: orders,
			status: "pending",
			want:   1,
		},
		{
			name:   "filter returned",
			orders: orders,
			status: "returned",
			want:   1,
		},
		{
			name:   "filter non-existent status",
			orders: orders,
			status: "cancelled",
			want:   0,
		},
		{
			name:   "empty status returns all",
			orders: orders,
			status: "",
			want:   4,
		},
		{
			name:   "case insensitive",
			orders: orders,
			status: "DELIVERED",
			want:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterOrdersByStatus(tt.orders, tt.status)
			if len(got) != tt.want {
				t.Errorf("FilterOrdersByStatus() returned %v orders, want %v", len(got), tt.want)
			}
		})
	}
}

func TestCalculateOrderTotal(t *testing.T) {
	tests := []struct {
		name  string
		items []models.OrderItem
		want  float64
	}{
		{
			name: "single item",
			items: []models.OrderItem{
				{ASIN: "B08N5WRWNW", Title: "Product 1", Quantity: 1, Price: 29.99},
			},
			want: 29.99,
		},
		{
			name: "multiple items same quantity",
			items: []models.OrderItem{
				{ASIN: "B08N5WRWNW", Title: "Product 1", Quantity: 1, Price: 29.99},
				{ASIN: "B07XYZ1234", Title: "Product 2", Quantity: 1, Price: 19.99},
			},
			want: 49.98,
		},
		{
			name: "items with different quantities",
			items: []models.OrderItem{
				{ASIN: "B08N5WRWNW", Title: "Product 1", Quantity: 2, Price: 10.00},
				{ASIN: "B07XYZ1234", Title: "Product 2", Quantity: 3, Price: 5.00},
			},
			want: 35.00,
		},
		{
			name:  "empty items",
			items: []models.OrderItem{},
			want:  0.0,
		},
		{
			name: "zero price items",
			items: []models.OrderItem{
				{ASIN: "B08N5WRWNW", Title: "Free Item", Quantity: 1, Price: 0.00},
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateOrderTotal(tt.items)
			if got != tt.want {
				t.Errorf("CalculateOrderTotal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderJSONRoundTrip(t *testing.T) {
	original := &models.Order{
		OrderID: "123-4567890-1234567",
		Date:    "2024-01-15",
		Total:   29.99,
		Status:  "delivered",
		Items: []models.OrderItem{
			{
				ASIN:     "B08N5WRWNW",
				Title:    "Product Name",
				Quantity: 1,
				Price:    29.99,
			},
		},
		Tracking: &models.Tracking{
			Carrier:        "UPS",
			TrackingNumber: "1Z999AA10123456784",
			Status:         "delivered",
			DeliveryDate:   "2024-01-17",
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal order: %v", err)
	}

	// Parse back from JSON
	parsed, err := ParseOrderFromJSON(data)
	if err != nil {
		t.Fatalf("Failed to parse order: %v", err)
	}

	// Verify all fields match
	if parsed.OrderID != original.OrderID {
		t.Errorf("OrderID = %v, want %v", parsed.OrderID, original.OrderID)
	}
	if parsed.Date != original.Date {
		t.Errorf("Date = %v, want %v", parsed.Date, original.Date)
	}
	if parsed.Total != original.Total {
		t.Errorf("Total = %v, want %v", parsed.Total, original.Total)
	}
	if parsed.Status != original.Status {
		t.Errorf("Status = %v, want %v", parsed.Status, original.Status)
	}
	if len(parsed.Items) != len(original.Items) {
		t.Errorf("Items length = %v, want %v", len(parsed.Items), len(original.Items))
	}
	if parsed.Tracking.Carrier != original.Tracking.Carrier {
		t.Errorf("Tracking.Carrier = %v, want %v", parsed.Tracking.Carrier, original.Tracking.Carrier)
	}
}
