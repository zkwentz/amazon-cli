package models

import "fmt"

// Error codes as defined in the PRD
const (
	ErrCodeAuthRequired  = "AUTH_REQUIRED"
	ErrCodeAuthExpired   = "AUTH_EXPIRED"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeRateLimited   = "RATE_LIMITED"
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodePurchaseFailed = "PURCHASE_FAILED"
	ErrCodeNetworkError  = "NETWORK_ERROR"
	ErrCodeAmazonError   = "AMAZON_ERROR"
)

// CLIError represents a structured CLI error
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewCLIError creates a new CLI error
func NewCLIError(code, message string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to a CLI error
func (e *CLIError) WithDetails(key string, value interface{}) *CLIError {
	e.Details[key] = value
	return e
}
