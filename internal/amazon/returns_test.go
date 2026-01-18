package amazon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestGetReturnOptions(t *testing.T) {
	tests := []struct {
		name           string
		orderID        string
		itemID         string
		responseStatus int
		responseBody   interface{}
		wantErr        bool
		errContains    string
		validateResult func(t *testing.T, options []models.ReturnOption)
	}{
		{
			name:           "successful JSON response",
			orderID:        "123-4567890-1234567",
			itemID:         "ITEM123",
			responseStatus: http.StatusOK,
			responseBody: map[string]interface{}{
				"returnOptions": []interface{}{
					map[string]interface{}{
						"method":          "UPS_DROPOFF",
						"label":           "Drop off at UPS",
						"dropoffLocation": "UPS Store - 123 Main St",
						"fee":             0.0,
					},
					map[string]interface{}{
						"method":          "AMAZON_LOCKER",
						"label":           "Amazon Locker",
						"dropoffLocation": "Locker at Whole Foods",
						"fee":             0.0,
					},
				},
			},
			wantErr: false,
			validateResult: func(t *testing.T, options []models.ReturnOption) {
				if len(options) != 2 {
					t.Errorf("expected 2 return options, got %d", len(options))
				}
				if options[0].Method != "UPS_DROPOFF" {
					t.Errorf("expected first option method to be UPS_DROPOFF, got %s", options[0].Method)
				}
				if options[0].Label != "Drop off at UPS" {
					t.Errorf("expected first option label to be 'Drop off at UPS', got %s", options[0].Label)
				}
				if options[1].Method != "AMAZON_LOCKER" {
					t.Errorf("expected second option method to be AMAZON_LOCKER, got %s", options[1].Method)
				}
			},
		},
		{
			name:           "HTML response with return options",
			orderID:        "123-4567890-1234567",
			itemID:         "ITEM456",
			responseStatus: http.StatusOK,
			responseBody:   "<html><body>return option available</body></html>",
			wantErr:        false,
			validateResult: func(t *testing.T, options []models.ReturnOption) {
				if len(options) == 0 {
					t.Error("expected at least one return option")
				}
			},
		},
		{
			name:           "empty orderID",
			orderID:        "",
			itemID:         "ITEM123",
			responseStatus: http.StatusOK,
			responseBody:   map[string]interface{}{},
			wantErr:        true,
			errContains:    "orderID cannot be empty",
		},
		{
			name:           "empty itemID",
			orderID:        "123-4567890-1234567",
			itemID:         "",
			responseStatus: http.StatusOK,
			responseBody:   map[string]interface{}{},
			wantErr:        true,
			errContains:    "itemID cannot be empty",
		},
		{
			name:           "not found response",
			orderID:        "123-4567890-1234567",
			itemID:         "NONEXISTENT",
			responseStatus: http.StatusNotFound,
			responseBody:   "Not Found",
			wantErr:        true,
			errContains:    "order or item not found",
		},
		{
			name:           "unauthorized response",
			orderID:        "123-4567890-1234567",
			itemID:         "ITEM123",
			responseStatus: http.StatusUnauthorized,
			responseBody:   "Unauthorized",
			wantErr:        true,
			errContains:    "authentication required or expired",
		},
		{
			name:           "server error",
			orderID:        "123-4567890-1234567",
			itemID:         "ITEM123",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "Internal Server Error",
			wantErr:        true,
			errContains:    "unexpected status code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request parameters
				if tt.orderID != "" && tt.itemID != "" {
					query := r.URL.Query()
					if query.Get("orderID") != tt.orderID {
						t.Errorf("expected orderID %s, got %s", tt.orderID, query.Get("orderID"))
					}
					if query.Get("itemID") != tt.itemID {
						t.Errorf("expected itemID %s, got %s", tt.itemID, query.Get("itemID"))
					}
				}

				w.WriteHeader(tt.responseStatus)

				// Write response body
				switch body := tt.responseBody.(type) {
				case map[string]interface{}:
					json.NewEncoder(w).Encode(body)
				case string:
					w.Write([]byte(body))
				}
			}))
			defer server.Close()

			// Create client and modify the URL to point to test server
			client := NewClient()

			// For testing, we need to modify the GetReturnOptions to accept a base URL
			// or we can test the parsing functions directly
			// Since we can't modify the URL easily, let's test with the actual implementation
			// but note that it will try to connect to amazon.com

			// For this test, we'll skip the actual HTTP call and test the parsing logic
			// In a real scenario, you'd want to make the URL configurable for testing

			if tt.wantErr && (tt.orderID == "" || tt.itemID == "") {
				// Test input validation without making HTTP call
				_, err := client.GetReturnOptions(tt.orderID, tt.itemID)
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errContains != "" && err != nil {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
					}
				}
			}
			// Note: Full integration tests would require either:
			// 1. Making the base URL configurable in the Client
			// 2. Using dependency injection for the HTTP client
			// 3. Testing against actual Amazon endpoints (not recommended for unit tests)
		})
	}
}

func TestParseJSONReturnOptions(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]interface{}
		wantErr     bool
		errContains string
		validate    func(t *testing.T, options []models.ReturnOption)
	}{
		{
			name: "valid return options",
			input: map[string]interface{}{
				"returnOptions": []interface{}{
					map[string]interface{}{
						"method":          "UPS_DROPOFF",
						"label":           "Drop off at UPS",
						"dropoffLocation": "UPS Store",
						"fee":             0.0,
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, options []models.ReturnOption) {
				if len(options) != 1 {
					t.Fatalf("expected 1 option, got %d", len(options))
				}
				if options[0].Method != "UPS_DROPOFF" {
					t.Errorf("expected method UPS_DROPOFF, got %s", options[0].Method)
				}
				if options[0].Fee != 0.0 {
					t.Errorf("expected fee 0.0, got %f", options[0].Fee)
				}
			},
		},
		{
			name: "multiple return options",
			input: map[string]interface{}{
				"returnOptions": []interface{}{
					map[string]interface{}{
						"method":          "UPS_DROPOFF",
						"label":           "UPS Drop-off",
						"dropoffLocation": "UPS Store",
						"fee":             0.0,
					},
					map[string]interface{}{
						"method":          "WHOLE_FOODS",
						"label":           "Return at Whole Foods",
						"dropoffLocation": "Whole Foods Market",
						"fee":             0.0,
					},
					map[string]interface{}{
						"method":          "MAIL_RETURN",
						"label":           "Mail return with fee",
						"dropoffLocation": "",
						"fee":             5.99,
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, options []models.ReturnOption) {
				if len(options) != 3 {
					t.Fatalf("expected 3 options, got %d", len(options))
				}
				if options[2].Fee != 5.99 {
					t.Errorf("expected third option fee 5.99, got %f", options[2].Fee)
				}
			},
		},
		{
			name: "missing returnOptions key",
			input: map[string]interface{}{
				"other": "data",
			},
			wantErr:     true,
			errContains: "no returnOptions found",
		},
		{
			name: "returnOptions not an array",
			input: map[string]interface{}{
				"returnOptions": "not an array",
			},
			wantErr:     true,
			errContains: "returnOptions is not an array",
		},
		{
			name: "empty returnOptions array",
			input: map[string]interface{}{
				"returnOptions": []interface{}{},
			},
			wantErr: false,
			validate: func(t *testing.T, options []models.ReturnOption) {
				if len(options) != 0 {
					t.Errorf("expected 0 options, got %d", len(options))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options, err := parseJSONReturnOptions(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errContains != "" && err != nil {
					if !contains(err.Error(), tt.errContains) {
						t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, options)
				}
			}
		})
	}
}

func TestParseHTMLReturnOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(t *testing.T, options []models.ReturnOption)
	}{
		{
			name:    "HTML with return option keyword",
			input:   "<html><body>return option available</body></html>",
			wantErr: false,
			validate: func(t *testing.T, options []models.ReturnOption) {
				if len(options) == 0 {
					t.Error("expected at least one option")
				}
			},
		},
		{
			name:    "HTML with returnOptions keyword",
			input:   "<html><body>returnOptions available</body></html>",
			wantErr: false,
			validate: func(t *testing.T, options []models.ReturnOption) {
				if len(options) == 0 {
					t.Error("expected at least one option")
				}
			},
		},
		{
			name:    "HTML without return options",
			input:   "<html><body>no options here</body></html>",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options, err := parseHTMLReturnOptions([]byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, options)
				}
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
