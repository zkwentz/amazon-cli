package models

import "fmt"

// Error code constants matching PRD specifications
const (
	// ErrCodeAuthRequired indicates user is not logged in
	ErrCodeAuthRequired = "AUTH_REQUIRED"
	// ErrCodeAuthExpired indicates authentication token has expired
	ErrCodeAuthExpired = "AUTH_EXPIRED"
	// ErrCodeNotFound indicates requested resource was not found
	ErrCodeNotFound = "NOT_FOUND"
	// ErrCodeRateLimited indicates too many requests have been made
	ErrCodeRateLimited = "RATE_LIMITED"
	// ErrCodeInvalidInput indicates invalid command input
	ErrCodeInvalidInput = "INVALID_INPUT"
	// ErrCodePurchaseFailed indicates purchase could not be completed
	ErrCodePurchaseFailed = "PURCHASE_FAILED"
	// ErrCodeNetworkError indicates network connectivity issue
	ErrCodeNetworkError = "NETWORK_ERROR"
	// ErrCodeAmazonError indicates Amazon returned an error
	ErrCodeAmazonError = "AMAZON_ERROR"
)

// Exit code constants matching PRD specifications
const (
	// ExitSuccess indicates successful execution
	ExitSuccess = 0
	// ExitGeneralError indicates general error
	ExitGeneralError = 1
	// ExitInvalidArgs indicates invalid arguments
	ExitInvalidArgs = 2
	// ExitAuthError indicates authentication error
	ExitAuthError = 3
	// ExitNetworkError indicates network error
	ExitNetworkError = 4
	// ExitRateLimited indicates rate limited
	ExitRateLimited = 5
	// ExitNotFound indicates resource not found
	ExitNotFound = 6
)

// CLIError represents a structured error with code, message, and optional details
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
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

// NewCLIErrorWithDetails creates a new CLIError with code, message, and details
func NewCLIErrorWithDetails(code, message string, details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// GetExitCode maps a CLIError code to the appropriate exit code
func (e *CLIError) GetExitCode() int {
	switch e.Code {
	case ErrCodeAuthRequired, ErrCodeAuthExpired:
		return ExitAuthError
	case ErrCodeNotFound:
		return ExitNotFound
	case ErrCodeRateLimited:
		return ExitRateLimited
	case ErrCodeInvalidInput:
		return ExitInvalidArgs
	case ErrCodeNetworkError:
		return ExitNetworkError
	case ErrCodePurchaseFailed, ErrCodeAmazonError:
		return ExitGeneralError
	default:
		return ExitGeneralError
	}
}

// GetExitCodeFromError returns the appropriate exit code for any error
// If the error is a CLIError, it uses the mapped exit code
// Otherwise, it returns ExitGeneralError
func GetExitCodeFromError(err error) int {
	if err == nil {
		return ExitSuccess
	}

	if cliErr, ok := err.(*CLIError); ok {
		return cliErr.GetExitCode()
	}

	return ExitGeneralError
}
