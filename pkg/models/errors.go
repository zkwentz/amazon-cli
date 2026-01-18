package models

import "fmt"

// Error codes as defined in PRD
const (
	AuthRequired    = "AUTH_REQUIRED"
	AuthExpired     = "AUTH_EXPIRED"
	NotFound        = "NOT_FOUND"
	RateLimited     = "RATE_LIMITED"
	InvalidInput    = "INVALID_INPUT"
	PurchaseFailed  = "PURCHASE_FAILED"
	NetworkError    = "NETWORK_ERROR"
	AmazonError     = "AMAZON_ERROR"
)

// CLIError represents a structured error for the CLI
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
