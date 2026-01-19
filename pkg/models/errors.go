package models

import "fmt"

// Error codes as defined in PRD
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

// Exit codes as defined in PRD
const (
	ExitSuccess        = 0
	ExitGeneralError   = 1
	ExitInvalidArgs    = 2
	ExitAuthError      = 3
	ExitNetworkError   = 4
	ExitRateLimited    = 5
	ExitNotFound       = 6
)

// CLIError represents a structured error for CLI output
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewCLIError creates a new CLIError
func NewCLIError(code, message string) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: map[string]interface{}{},
	}
}

// WithDetails adds details to a CLIError
func (e *CLIError) WithDetails(details map[string]interface{}) *CLIError {
	e.Details = details
	return e
}

// ExitCodeForError returns the appropriate exit code for an error code
func ExitCodeForError(code string) int {
	switch code {
	case ErrAuthRequired, ErrAuthExpired:
		return ExitAuthError
	case ErrNetworkError:
		return ExitNetworkError
	case ErrRateLimited:
		return ExitRateLimited
	case ErrNotFound:
		return ExitNotFound
	case ErrInvalidInput:
		return ExitInvalidArgs
	default:
		return ExitGeneralError
	}
}
