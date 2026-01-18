package models

import "fmt"

const (
	ErrorCodeAuthRequired  = "AUTH_REQUIRED"
	ErrorCodeAuthExpired   = "AUTH_EXPIRED"
	ErrorCodeNotFound      = "NOT_FOUND"
	ErrorCodeRateLimited   = "RATE_LIMITED"
	ErrorCodeInvalidInput  = "INVALID_INPUT"
	ErrorCodePurchaseFailed = "PURCHASE_FAILED"
	ErrorCodeNetworkError  = "NETWORK_ERROR"
	ErrorCodeAmazonError   = "AMAZON_ERROR"
)

type CLIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *CLIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewCLIError(code, message string, details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
