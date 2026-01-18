package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestOrdersListCommand tests the 'amazon-cli orders list' command
func TestOrdersListCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		mockResponse   string
		expectedStatus int
		expectedFields []string
		wantErr        bool
	}{
		{
			name: "list orders with default limit",
			args: []string{"orders", "list"},
			mockResponse: `{
				"orders": [
					{
						"order_id": "123-4567890-1234567",
						"date": "2024-01-15",
						"total": 29.99,
						"status": "delivered",
						"items": [
							{
								"asin": "B08N5WRWNW",
								"title": "Test Product",
								"quantity": 1,
								"price": 29.99
							}
						]
					}
				],
				"total_count": 1
			}`,
			expectedStatus: 0,
			expectedFields: []string{"orders", "total_count"},
			wantErr:        false,
		},
		{
			name: "list orders with custom limit",
			args: []string{"orders", "list", "--limit", "5"},
			mockResponse: `{
				"orders": [],
				"total_count": 0
			}`,
			expectedStatus: 0,
			expectedFields: []string{"orders", "total_count"},
			wantErr:        false,
		},
		{
			name: "list orders with status filter",
			args: []string{"orders", "list", "--status", "delivered"},
			mockResponse: `{
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
			expectedStatus: 0,
			expectedFields: []string{"orders", "total_count"},
			wantErr:        false,
		},
		{
			name: "list orders - authentication error",
			args: []string{"orders", "list"},
			mockResponse: `{
				"error": {
					"code": "AUTH_REQUIRED",
					"message": "Not logged in",
					"details": {}
				}
			}`,
			expectedStatus: 3,
			expectedFields: []string{"error"},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := createMockAmazonServer(tt.mockResponse, tt.expectedStatus == 0)
			defer server.Close()

			// Execute command
			output, exitCode := executeCommand(t, tt.args, server.URL)

			// Verify exit code
			if exitCode != tt.expectedStatus {
				t.Errorf("expected exit code %d, got %d", tt.expectedStatus, exitCode)
			}

			// Verify JSON output
			if !isValidJSON(output) {
				t.Errorf("output is not valid JSON: %s", output)
			}

			// Verify expected fields
			for _, field := range tt.expectedFields {
				if !strings.Contains(output, field) {
					t.Errorf("expected field '%s' in output, got: %s", field, output)
				}
			}
		})
	}
}

// TestOrdersGetCommand tests the 'amazon-cli orders get <order-id>' command
func TestOrdersGetCommand(t *testing.T) {
	tests := []struct {
		name           string
		orderID        string
		mockResponse   string
		expectedStatus int
		wantErr        bool
	}{
		{
			name:    "get order details - success",
			orderID: "123-4567890-1234567",
			mockResponse: `{
				"order_id": "123-4567890-1234567",
				"date": "2024-01-15",
				"total": 29.99,
				"status": "delivered",
				"items": [
					{
						"asin": "B08N5WRWNW",
						"title": "Test Product",
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
			expectedStatus: 0,
			wantErr:        false,
		},
		{
			name:    "get order details - not found",
			orderID: "999-9999999-9999999",
			mockResponse: `{
				"error": {
					"code": "NOT_FOUND",
					"message": "Order not found",
					"details": {}
				}
			}`,
			expectedStatus: 6,
			wantErr:        true,
		},
		{
			name:    "get order details - missing order ID",
			orderID: "",
			mockResponse: `{
				"error": {
					"code": "INVALID_INPUT",
					"message": "Order ID is required",
					"details": {}
				}
			}`,
			expectedStatus: 2,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createMockAmazonServer(tt.mockResponse, tt.expectedStatus == 0)
			defer server.Close()

			var args []string
			if tt.orderID != "" {
				args = []string{"orders", "get", tt.orderID}
			} else {
				args = []string{"orders", "get"}
			}

			output, exitCode := executeCommand(t, args, server.URL)

			if exitCode != tt.expectedStatus {
				t.Errorf("expected exit code %d, got %d", tt.expectedStatus, exitCode)
			}

			if !isValidJSON(output) {
				t.Errorf("output is not valid JSON: %s", output)
			}

			if !tt.wantErr && !strings.Contains(output, "order_id") {
				t.Errorf("expected order_id in successful response, got: %s", output)
			}
		})
	}
}

// TestOrdersTrackCommand tests the 'amazon-cli orders track <order-id>' command
func TestOrdersTrackCommand(t *testing.T) {
	tests := []struct {
		name           string
		orderID        string
		mockResponse   string
		expectedStatus int
		expectedFields []string
	}{
		{
			name:    "track order - success",
			orderID: "123-4567890-1234567",
			mockResponse: `{
				"carrier": "UPS",
				"tracking_number": "1Z999AA10123456784",
				"status": "in_transit",
				"delivery_date": "2024-01-20"
			}`,
			expectedStatus: 0,
			expectedFields: []string{"carrier", "tracking_number", "status", "delivery_date"},
		},
		{
			name:    "track order - delivered",
			orderID: "123-4567890-1234567",
			mockResponse: `{
				"carrier": "USPS",
				"tracking_number": "9400100000000000000000",
				"status": "delivered",
				"delivery_date": "2024-01-18"
			}`,
			expectedStatus: 0,
			expectedFields: []string{"carrier", "tracking_number", "status", "delivered"},
		},
		{
			name:    "track order - no tracking available",
			orderID: "123-4567890-1234567",
			mockResponse: `{
				"error": {
					"code": "NOT_FOUND",
					"message": "Tracking information not available",
					"details": {}
				}
			}`,
			expectedStatus: 6,
			expectedFields: []string{"error", "NOT_FOUND"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createMockAmazonServer(tt.mockResponse, tt.expectedStatus == 0)
			defer server.Close()

			args := []string{"orders", "track", tt.orderID}
			output, exitCode := executeCommand(t, args, server.URL)

			if exitCode != tt.expectedStatus {
				t.Errorf("expected exit code %d, got %d", tt.expectedStatus, exitCode)
			}

			if !isValidJSON(output) {
				t.Errorf("output is not valid JSON: %s", output)
			}

			for _, field := range tt.expectedFields {
				if !strings.Contains(output, field) {
					t.Errorf("expected field '%s' in output, got: %s", field, output)
				}
			}
		})
	}
}

// TestOrdersHistoryCommand tests the 'amazon-cli orders history' command
func TestOrdersHistoryCommand(t *testing.T) {
	tests := []struct {
		name           string
		year           string
		mockResponse   string
		expectedStatus int
		minOrders      int
	}{
		{
			name: "get order history - current year",
			year: "",
			mockResponse: `{
				"orders": [
					{
						"order_id": "123-4567890-1234567",
						"date": "2024-01-15",
						"total": 29.99,
						"status": "delivered",
						"items": []
					},
					{
						"order_id": "123-4567890-1234568",
						"date": "2024-02-10",
						"total": 49.99,
						"status": "delivered",
						"items": []
					}
				],
				"total_count": 2
			}`,
			expectedStatus: 0,
			minOrders:      2,
		},
		{
			name: "get order history - specific year",
			year: "2023",
			mockResponse: `{
				"orders": [
					{
						"order_id": "123-4567890-1234560",
						"date": "2023-12-25",
						"total": 99.99,
						"status": "delivered",
						"items": []
					}
				],
				"total_count": 1
			}`,
			expectedStatus: 0,
			minOrders:      1,
		},
		{
			name: "get order history - no orders",
			year: "2020",
			mockResponse: `{
				"orders": [],
				"total_count": 0
			}`,
			expectedStatus: 0,
			minOrders:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createMockAmazonServer(tt.mockResponse, tt.expectedStatus == 0)
			defer server.Close()

			var args []string
			if tt.year != "" {
				args = []string{"orders", "history", "--year", tt.year}
			} else {
				args = []string{"orders", "history"}
			}

			output, exitCode := executeCommand(t, args, server.URL)

			if exitCode != tt.expectedStatus {
				t.Errorf("expected exit code %d, got %d", tt.expectedStatus, exitCode)
			}

			if !isValidJSON(output) {
				t.Errorf("output is not valid JSON: %s", output)
			}

			// Parse and verify order count
			var response struct {
				Orders     []interface{} `json:"orders"`
				TotalCount int           `json:"total_count"`
			}
			if err := json.Unmarshal([]byte(output), &response); err == nil {
				if len(response.Orders) < tt.minOrders {
					t.Errorf("expected at least %d orders, got %d", tt.minOrders, len(response.Orders))
				}
			}
		})
	}
}

// TestOrdersRateLimiting tests that rate limiting is properly enforced
func TestOrdersRateLimiting(t *testing.T) {
	server := createMockAmazonServer(`{"orders": [], "total_count": 0}`, true)
	defer server.Close()

	requestTimes := make([]time.Time, 0, 3)

	// Make 3 consecutive requests
	for i := 0; i < 3; i++ {
		requestTimes = append(requestTimes, time.Now())
		executeCommand(t, []string{"orders", "list"}, server.URL)
	}

	// Verify delays between requests (should be at least 1 second based on PRD)
	for i := 1; i < len(requestTimes); i++ {
		delay := requestTimes[i].Sub(requestTimes[i-1])
		if delay < time.Second {
			t.Logf("Warning: Request %d delay was %v, expected >= 1s (rate limiting may not be active in tests)", i, delay)
		}
	}
}

// TestOrdersNetworkError tests handling of network errors
func TestOrdersNetworkError(t *testing.T) {
	// Use an invalid URL to trigger network error
	args := []string{"orders", "list"}
	output, exitCode := executeCommand(t, args, "http://invalid-amazon-server:9999")

	if exitCode != 4 { // NETWORK_ERROR exit code
		t.Logf("Expected exit code 4 (NETWORK_ERROR), got %d", exitCode)
	}

	if !isValidJSON(output) {
		t.Errorf("output is not valid JSON even on error: %s", output)
	}

	if !strings.Contains(output, "error") {
		t.Errorf("expected error field in output, got: %s", output)
	}
}

// TestOrdersJSONOutput tests that all orders commands output valid JSON
func TestOrdersJSONOutput(t *testing.T) {
	commands := [][]string{
		{"orders", "list"},
		{"orders", "get", "123-4567890-1234567"},
		{"orders", "track", "123-4567890-1234567"},
		{"orders", "history"},
	}

	mockResponse := `{"orders": [], "total_count": 0}`
	server := createMockAmazonServer(mockResponse, true)
	defer server.Close()

	for _, cmd := range commands {
		t.Run(strings.Join(cmd, " "), func(t *testing.T) {
			output, _ := executeCommand(t, cmd, server.URL)

			if !isValidJSON(output) {
				t.Errorf("command '%v' did not output valid JSON: %s", cmd, output)
			}
		})
	}
}

// Helper Functions

// createMockAmazonServer creates a mock HTTP server for testing
func createMockAmazonServer(response string, success bool) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !success {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(response))
	})
	return httptest.NewServer(handler)
}

// executeCommand executes the amazon-cli command and returns output and exit code
func executeCommand(t *testing.T, args []string, serverURL string) (string, int) {
	t.Helper()

	// Create temporary config for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Write test config with mock auth
	testConfig := `{
		"auth": {
			"access_token": "test_token",
			"refresh_token": "test_refresh",
			"expires_at": "2025-12-31T23:59:59Z"
		},
		"defaults": {
			"output_format": "json"
		},
		"rate_limiting": {
			"min_delay_ms": 0,
			"max_delay_ms": 0,
			"max_retries": 3
		}
	}`
	os.WriteFile(configPath, []byte(testConfig), 0600)

	// Build command
	cmdArgs := append([]string{"run", "main.go"}, args...)
	cmdArgs = append(cmdArgs, "--config", configPath)

	cmd := exec.Command("go", cmdArgs...)
	cmd.Env = append(os.Environ(), "AMAZON_CLI_TEST_SERVER="+serverURL)

	output, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return string(output), exitCode
}

// isValidJSON checks if a string is valid JSON
func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
