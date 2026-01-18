package models

import "fmt"

// Error codes as defined in the PRD
const (
	ErrAuthRequired   = "AUTH_REQUIRED"
	ErrAuthExpired    = "AUTH_EXPIRED"
	ErrNotFound       = "NOT_FOUND"
	ErrRateLimited    = "RATE_LIMITED"
	ErrInvalidInput   = "INVALID_INPUT"
	ErrPurchaseFailed = "PURCHASE_FAILED"
	ErrNetworkError   = "NETWORK_ERROR"
	ErrAmazonError    = "AMAZON_ERROR"
)

// CLIError represents a structured error with code, message and details
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewCLIError creates a new CLIError
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

// ErrorResponse wraps the error in the standard response format
type ErrorResponse struct {
	Error *CLIError `json:"error"`
}
