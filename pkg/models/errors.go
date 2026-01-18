package models

import "fmt"

// Error codes matching PRD section "Error Handling"
const (
	ErrCodeAuthRequired  = "AUTH_REQUIRED"
	ErrCodeAuthExpired   = "AUTH_EXPIRED"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeRateLimited   = "RATE_LIMITED"
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodePurchaseFailed = "PURCHASE_FAILED"
	ErrCodeNetworkError  = "NETWORK_ERROR"
	ErrCodeAmazonError   = "AMAZON_ERROR"
	ErrCodeGeneralError  = "GENERAL_ERROR"
)

// Exit codes matching PRD section "Exit Codes"
const (
	ExitSuccess          = 0
	ExitGeneralError     = 1
	ExitInvalidArguments = 2
	ExitAuthError        = 3
	ExitNetworkError     = 4
	ExitRateLimited      = 5
	ExitNotFound         = 6
)

// CLIError represents a structured error with code, message, and optional details
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
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

// WithDetails adds details to the error
func (e *CLIError) WithDetails(key string, value interface{}) *CLIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// GetExitCode maps error codes to exit codes according to PRD
func (e *CLIError) GetExitCode() int {
	switch e.Code {
	case ErrCodeAuthRequired, ErrCodeAuthExpired:
		return ExitAuthError
	case ErrCodeNetworkError:
		return ExitNetworkError
	case ErrCodeRateLimited:
		return ExitRateLimited
	case ErrCodeNotFound:
		return ExitNotFound
	case ErrCodeInvalidInput:
		return ExitInvalidArguments
	case ErrCodePurchaseFailed, ErrCodeAmazonError, ErrCodeGeneralError:
		return ExitGeneralError
	default:
		return ExitGeneralError
	}
}

// GetExitCodeFromError determines the appropriate exit code from any error
// If the error is a CLIError, it uses the specific exit code mapping
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
