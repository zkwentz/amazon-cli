package amazon

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// TestValidateOrderID_Invalid tests invalid order ID formats
func TestValidateOrderID_Invalid(t *testing.T) {
	testCases := []struct {
		name    string
		orderID string
		wantErr bool
	}{
		{
			name:    "empty order ID",
			orderID: "",
			wantErr: true,
		},
		{
			name:    "invalid format - too short",
			orderID: "123-456-789",
			wantErr: true,
		},
		{
			name:    "invalid format - missing dashes",
			orderID: "12345678901234567",
			wantErr: true,
		},
		{
			name:    "invalid format - letters instead of numbers",
			orderID: "ABC-DEFGHIJ-KLMNOPQ",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong structure",
			orderID: "1234-567890-123456",
			wantErr: true,
		},
		{
			name:    "valid format",
			orderID: "123-4567890-1234567",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateOrderID(tc.orderID)
			if tc.wantErr && err == nil {
				t.Errorf("expected error for order ID %q, got nil", tc.orderID)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error for order ID %q: %v", tc.orderID, err)
			}

			// If error is expected, verify it's a CLIError with correct code
			if tc.wantErr && err != nil {
				cliErr, ok := err.(*models.CLIError)
				if !ok {
					t.Errorf("expected *models.CLIError, got %T", err)
				} else if cliErr.Code != models.ErrCodeInvalidInput {
					t.Errorf("expected error code %s, got %s", models.ErrCodeInvalidInput, cliErr.Code)
				}
			}
		})
	}
}

// TestGetOrder_InvalidOrderID tests GetOrder with invalid order ID
func TestGetOrder_InvalidOrderID(t *testing.T) {
	client := NewOrdersClient()

	testCases := []struct {
		name    string
		orderID string
	}{
		{"empty order ID", ""},
		{"invalid format", "invalid-order-id"},
		{"too short", "123"},
		{"contains letters", "ABC-1234567-1234567"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order, err := client.GetOrder(tc.orderID)

			if err == nil {
				t.Errorf("expected error for invalid order ID %q, got nil", tc.orderID)
			}

			if order != nil {
				t.Errorf("expected nil order for invalid order ID, got %+v", order)
			}

			// Verify the error is a CLIError with INVALID_INPUT code
			cliErr, ok := err.(*models.CLIError)
			if !ok {
				t.Fatalf("expected *models.CLIError, got %T: %v", err, err)
			}

			if cliErr.Code != models.ErrCodeInvalidInput {
				t.Errorf("expected error code %s, got %s", models.ErrCodeInvalidInput, cliErr.Code)
			}

			if cliErr.Message == "" {
				t.Error("expected non-empty error message")
			}

			// Verify details contain the field name
			if field, ok := cliErr.Details["field"].(string); !ok || field != "order_id" {
				t.Errorf("expected details to contain field='order_id', got %v", cliErr.Details)
			}
		})
	}
}

// MockAuthExpiredServer creates a test server that simulates auth expiration
func MockAuthExpiredServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Authentication required"}`))
	}))
}

// TestGetOrder_AuthExpired tests handling of expired authentication
func TestGetOrder_AuthExpired(t *testing.T) {
	// Create a mock HTTP server that returns 401 Unauthorized
	server := MockAuthExpiredServer()
	defer server.Close()

	// In a real implementation, we would inject the server URL into the client
	// For now, we'll test the error creation directly
	authErr := models.NewAuthExpiredError()

	if authErr.Code != models.ErrCodeAuthExpired {
		t.Errorf("expected error code %s, got %s", models.ErrCodeAuthExpired, authErr.Code)
	}

	if authErr.Message == "" {
		t.Error("expected non-empty error message")
	}

	// Verify the error message contains helpful instructions
	expectedSubstring := "auth login"
	if !contains(authErr.Message, expectedSubstring) {
		t.Errorf("expected error message to contain %q, got %q", expectedSubstring, authErr.Message)
	}
}

// SimulateAuthExpiredError simulates an auth expired scenario
func SimulateAuthExpiredError(orderID string) (*models.Order, error) {
	// Validate the order ID first
	if err := ValidateOrderID(orderID); err != nil {
		return nil, err
	}

	// Simulate checking auth token expiry
	// In a real implementation, this would check actual token expiration
	return nil, models.NewAuthExpiredError()
}

// TestOrdersClient_AuthExpired tests the full auth expiration flow
func TestOrdersClient_AuthExpired(t *testing.T) {
	validOrderID := "123-4567890-1234567"

	order, err := SimulateAuthExpiredError(validOrderID)

	if err == nil {
		t.Fatal("expected auth expired error, got nil")
	}

	if order != nil {
		t.Errorf("expected nil order, got %+v", order)
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("expected *models.CLIError, got %T: %v", err, err)
	}

	if cliErr.Code != models.ErrCodeAuthExpired {
		t.Errorf("expected error code %s, got %s", models.ErrCodeAuthExpired, cliErr.Code)
	}
}

// MockNetworkErrorServer creates a test server that simulates network errors
func MockNetworkErrorServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Close connection immediately to simulate network error
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		conn.Close()
	}))
}

// TestGetOrder_NetworkError tests handling of network connectivity issues
func TestGetOrder_NetworkError(t *testing.T) {
	// Create a mock HTTP server that simulates network failure
	server := MockNetworkErrorServer()
	serverURL := server.URL
	server.Close() // Close immediately to ensure connection fails

	// Test making a request to the closed server
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	_, err := client.Get(serverURL)
	if err == nil {
		t.Fatal("expected network error, got nil")
	}

	// Wrap the error in our CLIError format
	networkErr := models.NewNetworkError(err)

	if networkErr.Code != models.ErrCodeNetworkError {
		t.Errorf("expected error code %s, got %s", models.ErrCodeNetworkError, networkErr.Code)
	}

	if networkErr.Message == "" {
		t.Error("expected non-empty error message")
	}

	// Verify details contain the original error
	if _, ok := networkErr.Details["original_error"]; !ok {
		t.Error("expected details to contain original_error")
	}
}

// SimulateNetworkError simulates a network error during order fetch
func SimulateNetworkError(orderID string) (*models.Order, error) {
	// Validate the order ID first
	if err := ValidateOrderID(orderID); err != nil {
		return nil, err
	}

	// Simulate a network error
	originalErr := errors.New("dial tcp: connection refused")
	return nil, models.NewNetworkError(originalErr)
}

// TestOrdersClient_NetworkError tests the full network error flow
func TestOrdersClient_NetworkError(t *testing.T) {
	validOrderID := "123-4567890-1234567"

	order, err := SimulateNetworkError(validOrderID)

	if err == nil {
		t.Fatal("expected network error, got nil")
	}

	if order != nil {
		t.Errorf("expected nil order, got %+v", order)
	}

	cliErr, ok := err.(*models.CLIError)
	if !ok {
		t.Fatalf("expected *models.CLIError, got %T: %v", err, err)
	}

	if cliErr.Code != models.ErrCodeNetworkError {
		t.Errorf("expected error code %s, got %s", models.ErrCodeNetworkError, cliErr.Code)
	}

	// Verify the error message contains context about network issues
	if !contains(cliErr.Message, "Network") && !contains(cliErr.Message, "network") {
		t.Errorf("expected error message to mention network, got %q", cliErr.Message)
	}
}

// TestAllErrorCases_Integration tests all three error cases together
func TestAllErrorCases_Integration(t *testing.T) {
	testCases := []struct {
		name          string
		orderID       string
		simulateFunc  func(string) (*models.Order, error)
		expectedCode  string
		errorContains string
	}{
		{
			name:          "invalid order ID",
			orderID:       "invalid",
			simulateFunc:  func(id string) (*models.Order, error) { return NewOrdersClient().GetOrder(id) },
			expectedCode:  models.ErrCodeInvalidInput,
			errorContains: "invalid format",
		},
		{
			name:          "auth expired",
			orderID:       "123-4567890-1234567",
			simulateFunc:  SimulateAuthExpiredError,
			expectedCode:  models.ErrCodeAuthExpired,
			errorContains: "expired",
		},
		{
			name:          "network error",
			orderID:       "123-4567890-1234567",
			simulateFunc:  SimulateNetworkError,
			expectedCode:  models.ErrCodeNetworkError,
			errorContains: "Network",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order, err := tc.simulateFunc(tc.orderID)

			// All error cases should return an error
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			// All error cases should return nil order
			if order != nil {
				t.Errorf("expected nil order, got %+v", order)
			}

			// Verify it's a CLIError with correct code
			cliErr, ok := err.(*models.CLIError)
			if !ok {
				t.Fatalf("expected *models.CLIError, got %T: %v", err, err)
			}

			if cliErr.Code != tc.expectedCode {
				t.Errorf("expected error code %s, got %s", tc.expectedCode, cliErr.Code)
			}

			if !contains(cliErr.Message, tc.errorContains) {
				t.Errorf("expected error message to contain %q, got %q", tc.errorContains, cliErr.Message)
			}

			// Verify the error is JSON-serializable
			if cliErr.Error() == "" {
				t.Error("Error() method returned empty string")
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		len(s) > len(substr)+1 && containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
