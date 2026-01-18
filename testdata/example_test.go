package testdata_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/zkwentz/amazon-cli/testdata"
)

// Example struct matching the order response schema
type Order struct {
	OrderID string  `json:"order_id"`
	Date    string  `json:"date"`
	Total   float64 `json:"total"`
	Status  string  `json:"status"`
}

type OrdersResponse struct {
	Orders     []Order `json:"orders"`
	TotalCount int     `json:"total_count"`
}

// TestLoadMockFile demonstrates how to load a mock file
func TestLoadMockFile(t *testing.T) {
	data := testdata.LoadMockFile(t, testdata.MockOrdersList)

	if len(data) == 0 {
		t.Error("Expected mock data, got empty")
	}

	// Verify it's valid JSON
	var response OrdersResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock data: %v", err)
	}

	if len(response.Orders) == 0 {
		t.Error("Expected orders in mock data")
	}
}

// TestUnmarshalMockFile demonstrates direct unmarshaling
func TestUnmarshalMockFile(t *testing.T) {
	var response OrdersResponse
	testdata.UnmarshalMockFile(t, testdata.MockOrdersList, &response)

	if response.TotalCount != 3 {
		t.Errorf("Expected 3 orders, got %d", response.TotalCount)
	}

	if len(response.Orders) != 3 {
		t.Errorf("Expected 3 orders in array, got %d", len(response.Orders))
	}

	// Verify first order
	firstOrder := response.Orders[0]
	if firstOrder.OrderID != "123-4567890-1234567" {
		t.Errorf("Expected order ID '123-4567890-1234567', got '%s'", firstOrder.OrderID)
	}

	if firstOrder.Status != "delivered" {
		t.Errorf("Expected status 'delivered', got '%s'", firstOrder.Status)
	}
}

// TestMockServer demonstrates using the mock HTTP server
func TestMockServer(t *testing.T) {
	server := testdata.NewMockServer()
	defer server.Close()

	// Register a handler that serves mock data
	server.OnWithMockFile(t, "GET", "/orders", testdata.MockOrdersList)

	// Make a request to the mock server
	resp, err := http.Get(server.URL + "/orders")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response OrdersResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Orders) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(response.Orders))
	}
}

// TestMockServerWithError demonstrates testing error responses
func TestMockServerWithError(t *testing.T) {
	server := testdata.NewMockServer()
	defer server.Close()

	// Register an error handler
	server.OnWithError(t, "GET", "/orders", testdata.MockErrorAuthRequired, 401)

	// Make a request
	resp, err := http.Get(server.URL + "/orders")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 401 {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	// Parse error response
	var errorResp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errorResp.Error.Code != "AUTH_REQUIRED" {
		t.Errorf("Expected error code 'AUTH_REQUIRED', got '%s'", errorResp.Error.Code)
	}
}

// TestMultipleMockResponses demonstrates testing different endpoints
func TestMultipleMockResponses(t *testing.T) {
	server := testdata.NewMockServer()
	defer server.Close()

	// Register multiple endpoints
	server.OnWithMockFile(t, "GET", "/orders", testdata.MockOrdersList)
	server.OnWithMockFile(t, "GET", "/cart", testdata.MockCartGet)
	server.OnWithMockFile(t, "GET", "/subscriptions", testdata.MockSubscriptionsList)

	// Test orders endpoint
	resp, err := http.Get(server.URL + "/orders")
	if err != nil {
		t.Fatalf("Failed to get orders: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Orders: Expected status 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Test cart endpoint
	resp, err = http.Get(server.URL + "/cart")
	if err != nil {
		t.Fatalf("Failed to get cart: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Cart: Expected status 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Test subscriptions endpoint
	resp, err = http.Get(server.URL + "/subscriptions")
	if err != nil {
		t.Fatalf("Failed to get subscriptions: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Subscriptions: Expected status 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
