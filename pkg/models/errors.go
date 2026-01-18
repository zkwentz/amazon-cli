package models

import (
	"encoding/json"
	"fmt"
)

// Error codes as defined in the PRD
const (
	ErrCodeAuthRequired   = "AUTH_REQUIRED"
	ErrCodeAuthExpired    = "AUTH_EXPIRED"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeRateLimited    = "RATE_LIMITED"
	ErrCodeInvalidInput   = "INVALID_INPUT"
	ErrCodePurchaseFailed = "PURCHASE_FAILED"
	ErrCodeNetworkError   = "NETWORK_ERROR"
	ErrCodeAmazonError    = "AMAZON_ERROR"
	ErrCodeInternalError  = "INTERNAL_ERROR"
)

// CLIError represents a structured error with code, message, and optional details
type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ErrorResponse wraps the error for JSON output
type ErrorResponse struct {
	Error *CLIError `json:"error"`
}

// NewCLIError creates a new CLIError with the specified code and message
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

// ToJSON converts the error to a JSON string for output
func (e *CLIError) ToJSON() string {
	resp := ErrorResponse{Error: e}
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		// Fallback to simple JSON if marshaling fails
		return fmt.Sprintf(`{"error":{"code":"%s","message":"%s"}}`, e.Code, e.Message)
	}
	return string(data)
}

// WithDetail adds a detail field to the error
func (e *CLIError) WithDetail(key string, value interface{}) *CLIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}
