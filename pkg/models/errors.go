package models

import "fmt"

// Error codes matching PRD specification
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

// CLIError represents a structured error for the CLI
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
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
