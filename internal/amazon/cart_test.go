package amazon

import (
	"fmt"
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// mockCartOps implements cartOperations interface for testing
type mockCartOps struct {
	cart           *models.Cart
	cartErr        error
	removeErr      error
	removedASINs   []string
	failOnASIN     string
	rateLimit      bool
	unauthorized   bool
}

func (m *mockCartOps) getCartInternal() (*models.Cart, error) {
	if m.unauthorized {
		return nil, models.NewCLIError(models.AuthRequired, "Authentication required", nil)
	}
	if m.cartErr != nil {
		return nil, m.cartErr
	}
	return m.cart, nil
}

func (m *mockCartOps) removeItemInternal(asin string) error {
	m.removedASINs = append(m.removedASINs, asin)

	if m.unauthorized {
		return models.NewCLIError(models.AuthRequired, "Authentication required", nil)
	}

	if m.rateLimit {
		return models.NewCLIError(models.RateLimited, "Rate limited by Amazon", nil)
	}

	if m.failOnASIN != "" && asin == m.failOnASIN {
		return fmt.Errorf("failed to remove item")
	}

	if m.removeErr != nil {
		return m.removeErr
	}

	return nil
}

// TestClearCart_EmptyCart tests clearing an already empty cart
func TestClearCart_EmptyCart(t *testing.T) {
	client := NewClient()

	mock := &mockCartOps{
		cart: &models.Cart{
			Items:     []models.CartItem{},
			ItemCount: 0,
		},
	}

	err := client.clearCartWithOps(mock)
	if err != nil {
		t.Errorf("ClearCart() on empty cart should not return error, got: %v", err)
	}

	if len(mock.removedASINs) != 0 {
		t.Errorf("Expected 0 remove calls for empty cart, got %d", len(mock.removedASINs))
	}
}

// TestClearCart_WithItems tests clearing a cart with items
func TestClearCart_WithItems(t *testing.T) {
	client := NewClient()
	expectedASINs := []string{"B001", "B002", "B003"}

	mock := &mockCartOps{
		cart: &models.Cart{
			Items: []models.CartItem{
				{ASIN: "B001", Title: "Item 1", Price: 10.00, Quantity: 1, Subtotal: 10.00},
				{ASIN: "B002", Title: "Item 2", Price: 20.00, Quantity: 2, Subtotal: 40.00},
				{ASIN: "B003", Title: "Item 3", Price: 15.00, Quantity: 1, Subtotal: 15.00},
			},
			Subtotal:     65.00,
			EstimatedTax: 5.20,
			Total:        70.20,
			ItemCount:    3,
		},
	}

	err := client.clearCartWithOps(mock)
	if err != nil {
		t.Errorf("ClearCart() should not return error, got: %v", err)
	}

	if len(mock.removedASINs) != len(expectedASINs) {
		t.Errorf("Expected %d remove calls, got %d", len(expectedASINs), len(mock.removedASINs))
	}

	// Verify all expected ASINs were removed
	for i, expectedASIN := range expectedASINs {
		if mock.removedASINs[i] != expectedASIN {
			t.Errorf("Expected ASIN %s at position %d, got %s", expectedASIN, i, mock.removedASINs[i])
		}
	}
}

// TestClearCart_NetworkError tests handling of network errors
func TestClearCart_NetworkError(t *testing.T) {
	client := NewClient()

	mock := &mockCartOps{
		cartErr: fmt.Errorf("network error: connection refused"),
	}

	err := client.clearCartWithOps(mock)
	if err == nil {
		t.Error("ClearCart() should return error on network failure")
	}

	// Check if it's a CLIError with NetworkError code
	if cliErr, ok := err.(*models.CLIError); ok {
		if cliErr.Code != models.NetworkError {
			t.Errorf("Expected error code %s, got %s", models.NetworkError, cliErr.Code)
		}
	} else {
		t.Errorf("Expected CLIError, got %T", err)
	}
}

// TestClearCart_UnauthorizedError tests handling of auth errors
func TestClearCart_UnauthorizedError(t *testing.T) {
	client := NewClient()

	mock := &mockCartOps{
		unauthorized: true,
	}

	err := client.clearCartWithOps(mock)
	if err == nil {
		t.Error("ClearCart() should return error on unauthorized response")
	}

	// Check if it's a CLIError with NetworkError code (auth error wrapped as network error)
	if cliErr, ok := err.(*models.CLIError); ok {
		if cliErr.Code != models.NetworkError {
			t.Errorf("Expected error code %s (auth error wrapped), got %s", models.NetworkError, cliErr.Code)
		}
	} else {
		t.Errorf("Expected CLIError, got %T", err)
	}
}

// TestClearCart_PartialFailure tests handling when some items fail to remove
func TestClearCart_PartialFailure(t *testing.T) {
	client := NewClient()

	mock := &mockCartOps{
		cart: &models.Cart{
			Items: []models.CartItem{
				{ASIN: "B001", Title: "Item 1", Price: 10.00},
				{ASIN: "B002", Title: "Item 2", Price: 20.00},
				{ASIN: "B003", Title: "Item 3", Price: 15.00},
			},
			ItemCount: 3,
		},
		failOnASIN: "B002", // Fail when trying to remove B002
	}

	err := client.clearCartWithOps(mock)
	if err == nil {
		t.Error("ClearCart() should return error when item removal fails")
	}

	// Check if it's a CLIError with AmazonError code
	if cliErr, ok := err.(*models.CLIError); ok {
		if cliErr.Code != models.AmazonError {
			t.Errorf("Expected error code %s, got %s", models.AmazonError, cliErr.Code)
		}
		// Check if the ASIN is included in details
		if asin, ok := cliErr.Details["asin"]; !ok || asin != "B002" {
			t.Errorf("Expected ASIN B002 in error details, got %v", cliErr.Details)
		}
	} else {
		t.Errorf("Expected CLIError, got %T", err)
	}

	// Verify that B001 was attempted but B003 was not (stopped after B002 failed)
	if len(mock.removedASINs) != 2 {
		t.Errorf("Expected 2 remove attempts, got %d", len(mock.removedASINs))
	}
}

// TestClearCart_RateLimited tests handling of rate limit errors
func TestClearCart_RateLimited(t *testing.T) {
	client := NewClient()

	mock := &mockCartOps{
		cart: &models.Cart{
			Items: []models.CartItem{
				{ASIN: "B001", Title: "Item 1", Price: 10.00},
			},
			ItemCount: 1,
		},
		rateLimit: true,
	}

	err := client.clearCartWithOps(mock)
	if err == nil {
		t.Error("ClearCart() should return error when rate limited")
	}

	// Check if it's a CLIError with AmazonError code (rate limit wrapped)
	if cliErr, ok := err.(*models.CLIError); ok {
		if cliErr.Code != models.AmazonError {
			t.Errorf("Expected error code %s (rate limit wrapped), got %s", models.AmazonError, cliErr.Code)
		}
	} else {
		t.Errorf("Expected CLIError, got %T", err)
	}
}
