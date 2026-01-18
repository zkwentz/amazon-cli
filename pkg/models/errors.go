package models

import "fmt"

// Error codes
const (
	ErrCodeAuthRequired     = "AUTH_REQUIRED"
	ErrCodeAuthExpired      = "AUTH_EXPIRED"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeRateLimited      = "RATE_LIMITED"
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodePurchaseFailed   = "PURCHASE_FAILED"
	ErrCodeNetworkError     = "NETWORK_ERROR"
	ErrCodeAmazonError      = "AMAZON_ERROR"
	ErrCodeItemNotReturnable = "ITEM_NOT_RETURNABLE"
	ErrCodeReturnWindowExpired = "RETURN_WINDOW_EXPIRED"
)

// CLIError represents a structured error with code, message, and details
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewCLIError creates a new CLIError
func NewCLIError(code, message string, details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewItemNotReturnableError creates an error for items that cannot be returned
func NewItemNotReturnableError(reason string) *CLIError {
	return &CLIError{
		Code:    ErrCodeItemNotReturnable,
		Message: "This item is not eligible for return",
		Details: map[string]interface{}{
			"reason": reason,
		},
	}
}

// NewReturnWindowExpiredError creates an error for expired return windows
func NewReturnWindowExpiredError(orderDate, expiryDate string) *CLIError {
	return &CLIError{
		Code:    ErrCodeReturnWindowExpired,
		Message: "The return window for this item has expired",
		Details: map[string]interface{}{
			"order_date":  orderDate,
			"expiry_date": expiryDate,
		},
	}
}
