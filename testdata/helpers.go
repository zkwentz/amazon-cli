package testdata

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// MockDataDir is the base directory for mock data files
const MockDataDir = "mocks"

// LoadMockFile reads a mock JSON file from the testdata directory
func LoadMockFile(t *testing.T, relativePath string) []byte {
	t.Helper()

	// Get the directory of the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current file path")
	}
	testDataDir := filepath.Dir(filename)

	fullPath := filepath.Join(testDataDir, MockDataDir, relativePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read mock file %s: %v", fullPath, err)
	}

	return data
}

// UnmarshalMockFile reads and unmarshals a mock JSON file into the provided interface
func UnmarshalMockFile(t *testing.T, relativePath string, v interface{}) {
	t.Helper()

	data := LoadMockFile(t, relativePath)
	err := json.Unmarshal(data, v)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock file %s: %v", relativePath, err)
	}
}

// MockServer creates a test HTTP server that serves mock responses
type MockServer struct {
	*httptest.Server
	handlers map[string]http.HandlerFunc
}

// NewMockServer creates a new mock HTTP server
func NewMockServer() *MockServer {
	ms := &MockServer{
		handlers: make(map[string]http.HandlerFunc),
	}

	ms.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		if handler, ok := ms.handlers[key]; ok {
			handler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))

	return ms
}

// On registers a handler for a specific HTTP method and path
func (ms *MockServer) On(method, path string, handler http.HandlerFunc) {
	key := method + " " + path
	ms.handlers[key] = handler
}

// OnWithMockFile registers a handler that serves a mock JSON file
func (ms *MockServer) OnWithMockFile(t *testing.T, method, path, mockFilePath string) {
	ms.On(method, path, func(w http.ResponseWriter, r *http.Request) {
		data := LoadMockFile(t, mockFilePath)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}

// OnWithError registers a handler that returns an error response
func (ms *MockServer) OnWithError(t *testing.T, method, path, errorMockFile string, statusCode int) {
	ms.On(method, path, func(w http.ResponseWriter, r *http.Request) {
		data := LoadMockFile(t, errorMockFile)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(data)
	})
}

// Example usage in tests:
//
// func TestGetOrders(t *testing.T) {
//     server := testdata.NewMockServer()
//     defer server.Close()
//
//     server.OnWithMockFile(t, "GET", "/orders", "orders/list_response.json")
//
//     client := NewClient(server.URL)
//     orders, err := client.GetOrders(10, "")
//
//     if err != nil {
//         t.Fatalf("Expected no error, got %v", err)
//     }
//
//     if len(orders.Orders) != 3 {
//         t.Errorf("Expected 3 orders, got %d", len(orders.Orders))
//     }
// }

// Common mock file paths as constants for convenience
const (
	// Auth
	MockAuthLoginSuccess            = "auth/login_success.json"
	MockAuthStatusAuthenticated     = "auth/auth_status_authenticated.json"
	MockAuthStatusUnauthenticated   = "auth/auth_status_unauthenticated.json"
	MockAuthLogoutSuccess           = "auth/logout_success.json"

	// Orders
	MockOrdersList      = "orders/list_response.json"
	MockOrdersGet       = "orders/get_order_response.json"
	MockOrdersTrack     = "orders/track_order_response.json"
	MockOrdersHistory   = "orders/history_response.json"

	// Returns
	MockReturnsList         = "returns/list_returnable_items.json"
	MockReturnsOptions      = "returns/return_options.json"
	MockReturnsCreate       = "returns/create_return_response.json"
	MockReturnsLabel        = "returns/return_label.json"
	MockReturnsStatus       = "returns/return_status.json"

	// Products
	MockProductSearch   = "products/search_response.json"
	MockProductDetails  = "products/product_details.json"
	MockProductReviews  = "products/product_reviews.json"

	// Cart
	MockCartGet             = "cart/cart_response.json"
	MockCartAdd             = "cart/add_to_cart_response.json"
	MockCartCheckoutPreview = "cart/checkout_preview.json"
	MockCartOrderConfirm    = "cart/order_confirmation.json"
	MockCartAddresses       = "cart/addresses.json"
	MockCartPaymentMethods  = "cart/payment_methods.json"

	// Subscriptions
	MockSubscriptionsList           = "subscriptions/list_subscriptions.json"
	MockSubscriptionsGet            = "subscriptions/get_subscription.json"
	MockSubscriptionsUpcoming       = "subscriptions/upcoming_deliveries.json"
	MockSubscriptionsSkip           = "subscriptions/skip_delivery_response.json"
	MockSubscriptionsUpdateFreq     = "subscriptions/update_frequency_response.json"
	MockSubscriptionsCancel         = "subscriptions/cancel_subscription_response.json"

	// Errors
	MockErrorAuthRequired   = "errors/auth_required.json"
	MockErrorAuthExpired    = "errors/auth_expired.json"
	MockErrorNotFound       = "errors/not_found.json"
	MockErrorRateLimited    = "errors/rate_limited.json"
	MockErrorInvalidInput   = "errors/invalid_input.json"
	MockErrorPurchaseFailed = "errors/purchase_failed.json"
	MockErrorNetwork        = "errors/network_error.json"
	MockErrorAmazon         = "errors/amazon_error.json"
)
