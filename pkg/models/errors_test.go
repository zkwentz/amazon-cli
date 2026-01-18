package models

import (
	"testing"
)

func TestNewCLIError(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		message string
		details map[string]interface{}
	}{
		{
			name:    "error with details",
			code:    ErrCodeInvalidInput,
			message: "test error message",
			details: map[string]interface{}{"field": "order_id"},
		},
		{
			name:    "error without details",
			code:    ErrCodeNotFound,
			message: "resource not found",
			details: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCLIError(tt.code, tt.message, tt.details)

			if err.Code != tt.code {
				t.Errorf("NewCLIError() code = %v, want %v", err.Code, tt.code)
			}
			if err.Message != tt.message {
				t.Errorf("NewCLIError() message = %v, want %v", err.Message, tt.message)
			}
			if err.Details == nil {
				t.Errorf("NewCLIError() details should never be nil")
			}
		})
	}
}

func TestCLIError_Error(t *testing.T) {
	err := NewCLIError(ErrCodeAuthExpired, "token expired", nil)
	expected := "AUTH_EXPIRED: token expired"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}
