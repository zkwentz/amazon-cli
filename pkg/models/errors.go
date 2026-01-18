package models

import "fmt"

// Error codes as defined in the PRD
const (
	ErrorCodeAuthRequired   = "AUTH_REQUIRED"
	ErrorCodeAuthExpired    = "AUTH_EXPIRED"
	ErrorCodeNotFound       = "NOT_FOUND"
	ErrorCodeRateLimited    = "RATE_LIMITED"
	ErrorCodeInvalidInput   = "INVALID_INPUT"
	ErrorCodePurchaseFailed = "PURCHASE_FAILED"
	ErrorCodeNetworkError   = "NETWORK_ERROR"
	ErrorCodeAmazonError    = "AMAZON_ERROR"
)

// CLIError represents an error in the CLI
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewCLIError creates a new CLI error
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

// ErrorResponse represents the JSON error response
type ErrorResponse struct {
	Error *CLIError `json:"error"`
}
