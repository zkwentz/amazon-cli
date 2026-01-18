package models

import "fmt"

// Error codes matching PRD specification
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

// CLIError represents a structured error with code and details
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
func NewCLIError(code, message string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *CLIError) WithDetails(key string, value interface{}) *CLIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}
