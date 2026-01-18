package models

import "fmt"

// Error codes matching PRD specification
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

// CLIError represents a structured error with code, message, and additional details
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Error implements the error interface for CLIError
func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewCLIError creates a new CLIError with the given code and message
func NewCLIError(code, message string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds additional details to the error
func (e *CLIError) WithDetails(details map[string]interface{}) *CLIError {
	e.Details = details
	return e
}

// WithDetail adds a single detail to the error
func (e *CLIError) WithDetail(key string, value interface{}) *CLIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}
