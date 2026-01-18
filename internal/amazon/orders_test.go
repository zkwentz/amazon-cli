package amazon

import (
	"testing"
)

func TestGetOrders(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		limit       int
		status      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid request with no limit and no status filter",
			limit:   0,
			status:  "",
			wantErr: false,
		},
		{
			name:    "valid request with limit",
			limit:   10,
			status:  "",
			wantErr: false,
		},
		{
			name:    "valid request with status filter - pending",
			limit:   0,
			status:  "pending",
			wantErr: false,
		},
		{
			name:    "valid request with status filter - delivered",
			limit:   0,
			status:  "delivered",
			wantErr: false,
		},
		{
			name:    "valid request with status filter - returned",
			limit:   0,
			status:  "returned",
			wantErr: false,
		},
		{
			name:    "valid request with limit and status filter",
			limit:   5,
			status:  "delivered",
			wantErr: false,
		},
		{
			name:        "invalid limit - negative",
			limit:       -1,
			status:      "",
			wantErr:     true,
			errContains: "limit must be non-negative",
		},
		{
			name:        "invalid status",
			limit:       0,
			status:      "invalid_status",
			wantErr:     true,
			errContains: "invalid status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.GetOrders(tt.limit, tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetOrders() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("GetOrders() error = %v, expected to contain %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GetOrders() unexpected error = %v", err)
				return
			}

			if resp == nil {
				t.Errorf("GetOrders() returned nil response")
				return
			}

			// Validate response structure
			if resp.Orders == nil {
				t.Errorf("GetOrders() response.Orders is nil, expected empty slice")
			}

			if resp.TotalCount < 0 {
				t.Errorf("GetOrders() response.TotalCount = %d, expected non-negative", resp.TotalCount)
			}

			// Verify limit is respected (when there's data)
			if tt.limit > 0 && len(resp.Orders) > tt.limit {
				t.Errorf("GetOrders() returned %d orders, expected at most %d", len(resp.Orders), tt.limit)
			}

			// Verify status filter is applied (when there's data and status is set)
			if tt.status != "" {
				for i, order := range resp.Orders {
					if order.Status != tt.status {
						t.Errorf("GetOrders() order[%d].Status = %s, expected %s", i, order.Status, tt.status)
					}
				}
			}
		})
	}
}

func TestGetOrders_ResponseStructure(t *testing.T) {
	client := NewClient()

	resp, err := client.GetOrders(0, "")
	if err != nil {
		t.Fatalf("GetOrders() unexpected error = %v", err)
	}

	if resp == nil {
		t.Fatal("GetOrders() returned nil response")
	}

	// Verify response has required fields
	if resp.Orders == nil {
		t.Error("response.Orders should not be nil")
	}

	// TotalCount should match the length of Orders
	if len(resp.Orders) != resp.TotalCount {
		t.Errorf("response.TotalCount = %d, but len(Orders) = %d", resp.TotalCount, len(resp.Orders))
	}
}

func TestGetOrders_LimitBehavior(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name  string
		limit int
	}{
		{"zero limit means no limit", 0},
		{"limit of 1", 1},
		{"limit of 5", 5},
		{"limit of 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.GetOrders(tt.limit, "")
			if err != nil {
				t.Fatalf("GetOrders() unexpected error = %v", err)
			}

			if tt.limit > 0 && len(resp.Orders) > tt.limit {
				t.Errorf("GetOrders() with limit=%d returned %d orders", tt.limit, len(resp.Orders))
			}
		})
	}
}

func TestGetOrders_StatusFilter(t *testing.T) {
	client := NewClient()

	statuses := []string{"", "pending", "delivered", "returned"}

	for _, status := range statuses {
		t.Run("status="+status, func(t *testing.T) {
			resp, err := client.GetOrders(0, status)
			if err != nil {
				t.Fatalf("GetOrders() unexpected error = %v", err)
			}

			if status != "" {
				for _, order := range resp.Orders {
					if order.Status != status {
						t.Errorf("Expected all orders to have status=%s, got %s", status, order.Status)
					}
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
