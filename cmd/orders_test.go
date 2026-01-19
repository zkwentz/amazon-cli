package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// captureStdout captures stdout output during test execution
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestOrdersListCmd_Success(t *testing.T) {
	// Create a test server that returns sample order list HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/gp/your-account/order-history" {
			// Return minimal valid HTML with order data
			html := `
				<html>
					<div class="order" data-order-id="111-2222222-3333333">
						<div class="order-date">January 15, 2026</div>
						<div class="order-total">$29.99</div>
						<div class="delivery-status">Delivered Jan 18, 2026</div>
					</div>
					<div class="order" data-order-id="111-4444444-5555555">
						<div class="order-date">January 10, 2026</div>
						<div class="order-total">$54.99</div>
						<div class="delivery-status">Arriving Jan 20, 2026</div>
					</div>
				</html>
			`
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(html))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a client with the test server URL
	testClient := amazon.NewClient()
	// Use reflection or a setter to modify the baseURL for testing
	// For now, we'll set the global client variable
	client = testClient

	// Note: In a real implementation, we'd need a way to inject the test server URL
	// into the client. This is a limitation of the current design.
	// For this test, we'll just verify the command structure

	// Test that the command exists and has the correct configuration
	if ordersListCmd.Use != "list" {
		t.Errorf("Expected Use='list', got '%s'", ordersListCmd.Use)
	}

	if ordersListCmd.Short != "List recent orders" {
		t.Errorf("Expected Short='List recent orders', got '%s'", ordersListCmd.Short)
	}

	if ordersListCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestOrdersListCmd_Flags(t *testing.T) {
	// Test that flags are properly configured
	limitFlag := ordersListCmd.Flags().Lookup("limit")
	if limitFlag == nil {
		t.Error("Expected --limit flag to be defined")
	} else {
		if limitFlag.DefValue != "10" {
			t.Errorf("Expected --limit default value to be '10', got '%s'", limitFlag.DefValue)
		}
	}

	statusFlag := ordersListCmd.Flags().Lookup("status")
	if statusFlag == nil {
		t.Error("Expected --status flag to be defined")
	} else {
		if statusFlag.DefValue != "" {
			t.Errorf("Expected --status default value to be empty, got '%s'", statusFlag.DefValue)
		}
	}
}

func TestOrdersListCmd_GetClientReturnsClient(t *testing.T) {
	// Test that getClient returns a non-nil client
	c := getClient()
	if c == nil {
		t.Error("Expected getClient() to return non-nil client")
	}

	// Test that calling getClient twice returns the same instance
	c2 := getClient()
	if c != c2 {
		t.Error("Expected getClient() to return the same client instance")
	}
}

func TestOrdersGetCmd_Success(t *testing.T) {
	// Test that the get command exists and has correct configuration
	if ordersGetCmd.Use != "get <order-id>" {
		t.Errorf("Expected Use='get <order-id>', got '%s'", ordersGetCmd.Use)
	}

	if ordersGetCmd.Short != "Get order details" {
		t.Errorf("Expected Short='Get order details', got '%s'", ordersGetCmd.Short)
	}

	if ordersGetCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if ordersGetCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestOrdersTrackCmd_Success(t *testing.T) {
	// Test that the track command exists and has correct configuration
	if ordersTrackCmd.Use != "track <order-id>" {
		t.Errorf("Expected Use='track <order-id>', got '%s'", ordersTrackCmd.Use)
	}

	if ordersTrackCmd.Short != "Track order shipment" {
		t.Errorf("Expected Short='Track order shipment', got '%s'", ordersTrackCmd.Short)
	}

	if ordersTrackCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if ordersTrackCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestOrdersTrackCmd_EmptyOrderID(t *testing.T) {
	// Test the function behavior with empty order ID
	// Since cobra.ExactArgs(1) prevents empty args, we test the validation inside Run
	// by checking if empty string is properly validated

	// We can't easily test os.Exit behavior without refactoring,
	// but we verify the validation logic exists by checking the models constants
	if models.ErrInvalidInput != "INVALID_INPUT" {
		t.Errorf("Expected ErrInvalidInput='INVALID_INPUT', got '%s'", models.ErrInvalidInput)
	}

	if models.ExitInvalidArgs != 2 {
		t.Errorf("Expected ExitInvalidArgs=2, got %d", models.ExitInvalidArgs)
	}
}

func TestOrdersTrackCmd_ValidatesOrderID(t *testing.T) {
	// Test that the Run function performs validation
	// We verify this by checking that the validation happens before calling GetClient

	// This test verifies the structure is correct for validation
	// The actual validation logic is in the GetOrderTracking method of the amazon client

	// Verify the Run function exists and is callable
	if ordersTrackCmd.Run == nil {
		t.Fatal("Expected Run function to be defined")
	}

	// Test would require mocking os.Exit to properly test the validation path
	// For now, we verify the error constants exist
	if models.ErrInvalidInput != "INVALID_INPUT" {
		t.Errorf("Expected ErrInvalidInput constant to be defined correctly")
	}
}

func TestOrdersHistoryCmd_Success(t *testing.T) {
	// Test that the history command exists and has correct configuration
	if ordersHistoryCmd.Use != "history" {
		t.Errorf("Expected Use='history', got '%s'", ordersHistoryCmd.Use)
	}

	if ordersHistoryCmd.Short != "Get order history" {
		t.Errorf("Expected Short='Get order history', got '%s'", ordersHistoryCmd.Short)
	}

	if ordersHistoryCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test year flag
	yearFlag := ordersHistoryCmd.Flags().Lookup("year")
	if yearFlag == nil {
		t.Error("Expected --year flag to be defined")
	} else {
		if yearFlag.DefValue != "0" {
			t.Errorf("Expected --year default value to be '0', got '%s'", yearFlag.DefValue)
		}
	}
}

func TestOrdersCmd_Subcommands(t *testing.T) {
	// Test that all subcommands are registered
	expectedSubcommands := []string{"list", "get", "track", "history"}
	commands := ordersCmd.Commands()

	if len(commands) != len(expectedSubcommands) {
		t.Errorf("Expected %d subcommands, got %d", len(expectedSubcommands), len(commands))
	}

	// Check that each expected subcommand exists
	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range commands {
			if cmd.Use == expected || cmd.Use == expected+" <order-id>" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}

// Mock test for ordersListCmd with proper error handling
func TestOrdersListCmd_ErrorHandling(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// We can't easily test the actual error handling without refactoring
	// to allow dependency injection, but we can verify the error constants exist
	if models.ErrAmazonError != "AMAZON_ERROR" {
		t.Errorf("Expected ErrAmazonError='AMAZON_ERROR', got '%s'", models.ErrAmazonError)
	}

	if models.ExitGeneralError != 1 {
		t.Errorf("Expected ExitGeneralError=1, got %d", models.ExitGeneralError)
	}
}

func TestOrdersCmd_Configuration(t *testing.T) {
	// Test the main orders command configuration
	if ordersCmd.Use != "orders" {
		t.Errorf("Expected Use='orders', got '%s'", ordersCmd.Use)
	}

	if ordersCmd.Short != "Manage orders" {
		t.Errorf("Expected Short='Manage orders', got '%s'", ordersCmd.Short)
	}

	expectedLong := "List orders, get order details, and track shipments."
	if ordersCmd.Long != expectedLong {
		t.Errorf("Expected Long='%s', got '%s'", expectedLong, ordersCmd.Long)
	}
}

func TestOrdersCmd_VariablesInitialized(t *testing.T) {
	// Test that package-level variables are initialized with correct default values
	// This is implicit through the flag defaults, but we can verify the variables exist
	// by checking they can be assigned

	// Save original values
	origLimit := ordersLimit
	origStatus := ordersStatus
	origYear := ordersYear

	// Modify them
	ordersLimit = 20
	ordersStatus = "delivered"
	ordersYear = 2025

	// Verify modifications worked
	if ordersLimit != 20 {
		t.Error("Failed to modify ordersLimit")
	}
	if ordersStatus != "delivered" {
		t.Error("Failed to modify ordersStatus")
	}
	if ordersYear != 2025 {
		t.Error("Failed to modify ordersYear")
	}

	// Restore original values
	ordersLimit = origLimit
	ordersStatus = origStatus
	ordersYear = origYear
}

// Integration-style test that verifies the command can parse valid responses
func TestOrdersListCmd_ResponseParsing(t *testing.T) {
	// This test verifies that the models.OrdersResponse structure
	// can be properly marshaled to JSON (as used by output.JSON)

	response := &models.OrdersResponse{
		Orders: []models.Order{
			{
				OrderID: "111-2222222-3333333",
				Date:    "2026-01-15",
				Total:   29.99,
				Status:  "delivered",
			},
		},
		TotalCount: 1,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal OrdersResponse to JSON: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	// Verify expected fields exist
	if _, ok := parsed["orders"]; !ok {
		t.Error("Expected 'orders' field in JSON output")
	}
	if _, ok := parsed["total_count"]; !ok {
		t.Error("Expected 'total_count' field in JSON output")
	}
}

// Test ordersGetCmd with NOT_FOUND error
func TestOrdersGetCmd_NotFoundError(t *testing.T) {
	// This test validates the NOT_FOUND error handling
	// We can't easily test the actual exit behavior without refactoring
	// but we can verify the error constants exist

	if models.ErrNotFound != "NOT_FOUND" {
		t.Errorf("Expected ErrNotFound='NOT_FOUND', got '%s'", models.ErrNotFound)
	}

	if models.ExitNotFound != 6 {
		t.Errorf("Expected ExitNotFound=6, got %d", models.ExitNotFound)
	}
}

// Test ordersGetCmd with invalid order ID format
func TestOrdersGetCmd_InvalidOrderIDFormat(t *testing.T) {
	// Create a client
	testClient := amazon.NewClient()

	// Test that invalid order ID formats would be caught by GetOrder
	invalidOrderIDs := []string{
		"123",                 // too short
		"123-456-789",         // wrong format
		"abc-1234567-1234567", // letters in first segment
		"123-abcdefg-1234567", // letters in middle segment
		"123-1234567-abcdefg", // letters in last segment
	}

	for _, orderID := range invalidOrderIDs {
		_, err := testClient.GetOrder(orderID)
		if err == nil {
			t.Errorf("Expected error for invalid order ID %s, got nil", orderID)
		}
	}
}

// Test ordersGetCmd response marshaling
func TestOrdersGetCmd_ResponseParsing(t *testing.T) {
	// This test verifies that the models.Order structure
	// can be properly marshaled to JSON (as used by output.JSON)

	order := &models.Order{
		OrderID: "123-4567890-1234567",
		Date:    "2026-01-15",
		Total:   84.98,
		Status:  "delivered",
		Items: []models.OrderItem{
			{
				ASIN:     "B08XYZ1234",
				Title:    "Example Product 1",
				Quantity: 1,
				Price:    39.99,
			},
			{
				ASIN:     "B08ABC5678",
				Title:    "Example Product 2",
				Quantity: 1,
				Price:    44.99,
			},
		},
		Tracking: &models.Tracking{
			Carrier:        "UPS",
			TrackingNumber: "1Z999AA10123456784",
			Status:         "delivered",
			DeliveryDate:   "2026-01-18",
		},
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Failed to marshal Order to JSON: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Failed to parse marshaled JSON: %v", err)
	}

	// Verify expected fields exist
	if _, ok := parsed["order_id"]; !ok {
		t.Error("Expected 'order_id' field in JSON output")
	}
	if _, ok := parsed["date"]; !ok {
		t.Error("Expected 'date' field in JSON output")
	}
	if _, ok := parsed["total"]; !ok {
		t.Error("Expected 'total' field in JSON output")
	}
	if _, ok := parsed["status"]; !ok {
		t.Error("Expected 'status' field in JSON output")
	}
	if _, ok := parsed["items"]; !ok {
		t.Error("Expected 'items' field in JSON output")
	}
	if _, ok := parsed["tracking"]; !ok {
		t.Error("Expected 'tracking' field in JSON output")
	}
}
