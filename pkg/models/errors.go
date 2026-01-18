package models

import "fmt"

// Error codes as specified in PRD
const (
	ErrCodeAuthRequired    = "AUTH_REQUIRED"
	ErrCodeAuthExpired     = "AUTH_EXPIRED"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeRateLimited     = "RATE_LIMITED"
	ErrCodeInvalidInput    = "INVALID_INPUT"
	ErrCodePurchaseFailed  = "PURCHASE_FAILED"
	ErrCodeNetworkError    = "NETWORK_ERROR"
	ErrCodeAmazonError     = "AMAZON_ERROR"
)

// CLIError represents a structured error with code, message, and details
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewCLIError creates a new CLIError with the given code and message
func NewCLIError(code, message string, details map[string]interface{}) *CLIError {
	if details == nil {
		details = make(map[string]interface{})
	}
	return &CLIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error constructors
func NewAuthExpiredError() *CLIError {
	return NewCLIError(
		ErrCodeAuthExpired,
		"Authentication token has expired. Run 'amazon-cli auth login' to re-authenticate.",
		nil,
	)
}

func NewNotFoundError(resource string) *CLIError {
	return NewCLIError(
		ErrCodeNotFound,
		fmt.Sprintf("Resource not found: %s", resource),
		map[string]interface{}{"resource": resource},
	)
}

func NewNetworkError(err error) *CLIError {
	return NewCLIError(
		ErrCodeNetworkError,
		fmt.Sprintf("Network connectivity issue: %v", err),
		map[string]interface{}{"original_error": err.Error()},
	)
}

func NewInvalidInputError(field, reason string) *CLIError {
	return NewCLIError(
		ErrCodeInvalidInput,
		fmt.Sprintf("Invalid input for %s: %s", field, reason),
		map[string]interface{}{"field": field, "reason": reason},
	)
}
