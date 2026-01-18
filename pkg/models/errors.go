package models

import "fmt"

// Exit codes as defined in PRD
const (
	ExitSuccess         = 0 // Success
	ExitGeneralError    = 1 // General error
	ExitInvalidArgs     = 2 // Invalid arguments
	ExitAuthError       = 3 // Authentication error
	ExitNetworkError    = 4 // Network error
	ExitRateLimited     = 5 // Rate limited
	ExitNotFound        = 6 // Not found
)

// Error codes for CLI operations
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

// CLIError represents a structured error with an error code
type CLIError struct {
	Code     string                 `json:"code"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	ExitCode int                    `json:"-"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewCLIError creates a new CLI error with the given code and message
func NewCLIError(code, message string, exitCode int) *CLIError {
	return &CLIError{
		Code:     code,
		Message:  message,
		Details:  make(map[string]interface{}),
		ExitCode: exitCode,
	}
}

// WithDetails adds details to the error
func (e *CLIError) WithDetails(key string, value interface{}) *CLIError {
	e.Details[key] = value
	return e
}

// GetExitCode extracts the exit code from an error, defaulting to ExitGeneralError
func GetExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	if cliErr, ok := err.(*CLIError); ok {
		return cliErr.ExitCode
	}

	return ExitGeneralError
}

// Common error constructors
func NewAuthRequiredError(message string) *CLIError {
	if message == "" {
		message = "Authentication required. Run 'amazon-cli auth login' to authenticate."
	}
	return NewCLIError(ErrCodeAuthRequired, message, ExitAuthError)
}

func NewAuthExpiredError(message string) *CLIError {
	if message == "" {
		message = "Authentication token has expired. Run 'amazon-cli auth login' to re-authenticate."
	}
	return NewCLIError(ErrCodeAuthExpired, message, ExitAuthError)
}

func NewNotFoundError(message string) *CLIError {
	if message == "" {
		message = "Resource not found."
	}
	return NewCLIError(ErrCodeNotFound, message, ExitNotFound)
}

func NewRateLimitedError(message string) *CLIError {
	if message == "" {
		message = "Rate limited. Please try again later."
	}
	return NewCLIError(ErrCodeRateLimited, message, ExitRateLimited)
}

func NewInvalidInputError(message string) *CLIError {
	if message == "" {
		message = "Invalid input provided."
	}
	return NewCLIError(ErrCodeInvalidInput, message, ExitInvalidArgs)
}

func NewNetworkError(message string) *CLIError {
	if message == "" {
		message = "Network error occurred."
	}
	return NewCLIError(ErrCodeNetworkError, message, ExitNetworkError)
}

func NewAmazonError(message string) *CLIError {
	if message == "" {
		message = "Amazon returned an error."
	}
	return NewCLIError(ErrCodeAmazonError, message, ExitGeneralError)
}

func NewPurchaseFailedError(message string) *CLIError {
	if message == "" {
		message = "Purchase could not be completed."
	}
	return NewCLIError(ErrCodePurchaseFailed, message, ExitGeneralError)
}
