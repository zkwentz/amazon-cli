package models

import "fmt"

// Error codes as defined in the PRD
const (
	// ErrAuthRequired indicates that authentication is required
	ErrAuthRequired = "AUTH_REQUIRED"
	// ErrAuthExpired indicates that the authentication token has expired
	ErrAuthExpired = "AUTH_EXPIRED"
	// ErrNotFound indicates that the requested resource was not found
	ErrNotFound = "NOT_FOUND"
	// ErrRateLimited indicates that too many requests have been made
	ErrRateLimited = "RATE_LIMITED"
	// ErrInvalidInput indicates that the command input was invalid
	ErrInvalidInput = "INVALID_INPUT"
	// ErrPurchaseFailed indicates that a purchase could not be completed
	ErrPurchaseFailed = "PURCHASE_FAILED"
	// ErrNetworkError indicates a network connectivity issue
	ErrNetworkError = "NETWORK_ERROR"
	// ErrAmazonError indicates that Amazon returned an error
	ErrAmazonError = "AMAZON_ERROR"
)

// CLIError represents a structured error with code, message, and details
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	// Internal fields for verbose output
	cause      error
	stackTrace string
}

// Error implements the error interface
func (e *CLIError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error for error unwrapping
func (e *CLIError) Unwrap() error {
	return e.cause
}

// Cause returns the underlying error
func (e *CLIError) Cause() error {
	return e.cause
}

// StackTrace returns the stack trace if available
func (e *CLIError) StackTrace() string {
	return e.stackTrace
}

// NewCLIError creates a new CLIError with the specified code and message
func NewCLIError(code, message string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// NewCLIErrorWithCause creates a new CLIError with the specified code, message, and underlying cause
func NewCLIErrorWithCause(code, message string, cause error) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
		cause:   cause,
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

// WithStackTrace adds a stack trace to the error
func (e *CLIError) WithStackTrace(trace string) *CLIError {
	e.stackTrace = trace
	return e
}

// ErrorResponse represents the JSON structure for error output
type ErrorResponse struct {
	Error *ErrorDetail `json:"error"`
}

// ErrorDetail represents the error detail in the JSON response
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	// Debug info only included when verbose flag is set
	Debug *DebugInfo `json:"debug,omitempty"`
}

// DebugInfo contains debug information for verbose error output
type DebugInfo struct {
	Cause      string `json:"cause,omitempty"`
	StackTrace string `json:"stack_trace,omitempty"`
}
