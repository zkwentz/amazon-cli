package models

import "fmt"

// Error codes as defined in PRD
const (
	// AuthRequired indicates user is not logged in
	AuthRequired = "AUTH_REQUIRED"
	// AuthExpired indicates authentication token has expired
	AuthExpired = "AUTH_EXPIRED"
	// NotFound indicates requested resource was not found
	NotFound = "NOT_FOUND"
	// RateLimited indicates too many requests were made
	RateLimited = "RATE_LIMITED"
	// InvalidInput indicates invalid command input
	InvalidInput = "INVALID_INPUT"
	// PurchaseFailed indicates purchase could not be completed
	PurchaseFailed = "PURCHASE_FAILED"
	// NetworkError indicates network connectivity issue
	NetworkError = "NETWORK_ERROR"
	// AmazonError indicates Amazon returned an error
	AmazonError = "AMAZON_ERROR"
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

// NewCLIError creates a new CLIError with the given code and message
func NewCLIError(code, message string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// NewCLIErrorWithDetails creates a new CLIError with code, message, and details
func NewCLIErrorWithDetails(code, message string, details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ErrorResponse wraps a CLIError for JSON output
type ErrorResponse struct {
	Error *CLIError `json:"error"`
}
