package models

import "errors"

// Sentinel errors for common error conditions
// These can be used with errors.Is() to check for specific error types
var (
	// ErrEmptyCart is returned when attempting to checkout with an empty cart
	ErrEmptyCart = errors.New("cart is empty")

	// ErrInvalidASIN is returned when an ASIN is empty or invalid
	ErrInvalidASIN = errors.New("invalid ASIN")

	// ErrInvalidQuantity is returned when quantity is zero or negative
	ErrInvalidQuantity = errors.New("invalid quantity")

	// ErrAddressNotFound is returned when the specified address does not exist
	ErrAddressNotFound = errors.New("address not found")

	// ErrPaymentMethodNotFound is returned when the specified payment method does not exist
	ErrPaymentMethodNotFound = errors.New("payment method not found")

	// ErrCheckoutFailed is returned when the checkout process fails
	ErrCheckoutFailed = errors.New("checkout failed")
)

// CLIError represents a structured error with a code for CLI output
type CLIError struct {
	Code    string
	Message string
	Err     error // Wrapped underlying error
}

// Error implements the error interface
func (e *CLIError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the wrapped error for use with errors.Is and errors.As
func (e *CLIError) Unwrap() error {
	return e.Err
}

// Error code constants for CLI output
const (
	ErrCodeAuthRequired  = "AUTH_REQUIRED"
	ErrCodeAuthExpired   = "AUTH_EXPIRED"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeRateLimited   = "RATE_LIMITED"
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodePurchaseFail  = "PURCHASE_FAILED"
	ErrCodeNetworkError  = "NETWORK_ERROR"
	ErrCodeAmazonError   = "AMAZON_ERROR"
)

// NewCLIError creates a new CLIError with the given code, message, and optional wrapped error
func NewCLIError(code, message string, err error) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
