package models

import "fmt"

// Error codes
const (
	// Authentication errors
	ErrCodeAuthRequired = "AUTH_REQUIRED"
	ErrCodeAuthExpired  = "AUTH_EXPIRED"

	// Resource errors
	ErrCodeNotFound = "NOT_FOUND"

	// Network and rate limiting errors
	ErrCodeRateLimited  = "RATE_LIMITED"
	ErrCodeNetworkError = "NETWORK_ERROR"
	ErrCodeAmazonError  = "AMAZON_ERROR"

	// Input validation errors
	ErrCodeInvalidInput = "INVALID_INPUT"

	// Purchase errors
	ErrCodePurchaseFailed = "PURCHASE_FAILED"

	// HTML response errors (CAPTCHA, login redirects)
	ErrCodeCaptchaRequired = "CAPTCHA_REQUIRED"
	ErrCodeLoginRequired   = "LOGIN_REQUIRED"
	ErrCodeHTMLResponse    = "HTML_RESPONSE"
)

// CLIError represents a structured error for the CLI
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

// NewCLIError creates a new CLIError
func NewCLIError(code, message string, details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewAuthRequiredError creates an error for when authentication is required
func NewAuthRequiredError() *CLIError {
	return &CLIError{
		Code:    ErrCodeAuthRequired,
		Message: "Authentication required. Run 'amazon-cli auth login' to authenticate.",
	}
}

// NewAuthExpiredError creates an error for when authentication has expired
func NewAuthExpiredError() *CLIError {
	return &CLIError{
		Code:    ErrCodeAuthExpired,
		Message: "Authentication token has expired. Run 'amazon-cli auth login' to re-authenticate.",
	}
}

// NewCaptchaRequiredError creates an error for when Amazon returns a CAPTCHA page
func NewCaptchaRequiredError(details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    ErrCodeCaptchaRequired,
		Message: "Amazon requires CAPTCHA verification. Please log in through the web browser and try again later.",
		Details: details,
	}
}

// NewLoginRequiredError creates an error for when Amazon redirects to login
func NewLoginRequiredError(details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    ErrCodeLoginRequired,
		Message: "Amazon requires re-authentication. Run 'amazon-cli auth login' to authenticate.",
		Details: details,
	}
}

// NewHTMLResponseError creates an error for when an unexpected HTML response is received
func NewHTMLResponseError(details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    ErrCodeHTMLResponse,
		Message: "Received unexpected HTML response instead of expected data. This may indicate a CAPTCHA, login redirect, or other Amazon security measure.",
		Details: details,
	}
}

// NewRateLimitedError creates an error for rate limiting
func NewRateLimitedError() *CLIError {
	return &CLIError{
		Code:    ErrCodeRateLimited,
		Message: "Rate limited by Amazon. Please wait a few moments and try again.",
	}
}

// NewNetworkError creates an error for network issues
func NewNetworkError(err error) *CLIError {
	return &CLIError{
		Code:    ErrCodeNetworkError,
		Message: fmt.Sprintf("Network error: %v", err),
	}
}

// NewAmazonError creates an error for Amazon-specific issues
func NewAmazonError(message string, details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    ErrCodeAmazonError,
		Message: message,
		Details: details,
	}
}

// NewNotFoundError creates an error for resource not found
func NewNotFoundError(resource string) *CLIError {
	return &CLIError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// NewInvalidInputError creates an error for invalid input
func NewInvalidInputError(message string) *CLIError {
	return &CLIError{
		Code:    ErrCodeInvalidInput,
		Message: message,
	}
}

// NewPurchaseFailedError creates an error for failed purchases
func NewPurchaseFailedError(message string, details map[string]interface{}) *CLIError {
	return &CLIError{
		Code:    ErrCodePurchaseFailed,
		Message: message,
		Details: details,
	}
}
