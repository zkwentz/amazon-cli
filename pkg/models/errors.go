package models

import "fmt"

// Error codes as defined in PRD
const (
	ErrCodeAuthRequired   = "AUTH_REQUIRED"
	ErrCodeAuthExpired    = "AUTH_EXPIRED"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeRateLimited    = "RATE_LIMITED"
	ErrCodeInvalidInput   = "INVALID_INPUT"
	ErrCodePurchaseFailed = "PURCHASE_FAILED"
	ErrCodeNetworkError   = "NETWORK_ERROR"
	ErrCodeAmazonError    = "AMAZON_ERROR"
)

// CLIError represents a structured error with a code and details
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

// ErrorResponse represents the JSON error response format
type ErrorResponse struct {
	Error *CLIError `json:"error"`
}
