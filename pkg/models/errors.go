package models

import "fmt"

// Error code constants matching PRD specification
const (
	// AuthRequired indicates user is not logged in
	AuthRequired = "AUTH_REQUIRED"
	// AuthExpired indicates authentication token has expired
	AuthExpired = "AUTH_EXPIRED"
	// NotFound indicates requested resource was not found
	NotFound = "NOT_FOUND"
	// RateLimited indicates too many requests have been made
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

// CLIError represents a structured error response for the CLI
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("%s: %s (details: %v)", e.Code, e.Message, e.Details)
	}
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

// WithDetails adds details to an existing CLIError
func (e *CLIError) WithDetails(key string, value interface{}) *CLIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}
