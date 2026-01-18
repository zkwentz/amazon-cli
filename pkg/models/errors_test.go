package models

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorWrapping(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantIs         error
		wantMessage    string
		wantUnwrapping bool
	}{
		{
			name:           "wrapped ErrInvalidASIN",
			err:            fmt.Errorf("add to cart: %w", ErrInvalidASIN),
			wantIs:         ErrInvalidASIN,
			wantMessage:    "add to cart: invalid ASIN",
			wantUnwrapping: true,
		},
		{
			name:           "wrapped ErrInvalidQuantity",
			err:            fmt.Errorf("add to cart: %w (got %d)", ErrInvalidQuantity, -1),
			wantIs:         ErrInvalidQuantity,
			wantMessage:    "add to cart: invalid quantity (got -1)",
			wantUnwrapping: true,
		},
		{
			name:           "wrapped ErrEmptyCart",
			err:            fmt.Errorf("complete checkout: %w", ErrEmptyCart),
			wantIs:         ErrEmptyCart,
			wantMessage:    "complete checkout: cart is empty",
			wantUnwrapping: true,
		},
		{
			name:           "wrapped ErrAddressNotFound",
			err:            fmt.Errorf("complete checkout: %w: %s", ErrAddressNotFound, "addr123"),
			wantIs:         ErrAddressNotFound,
			wantMessage:    "complete checkout: address not found: addr123",
			wantUnwrapping: true,
		},
		{
			name:           "wrapped ErrPaymentMethodNotFound",
			err:            fmt.Errorf("complete checkout: %w: %s", ErrPaymentMethodNotFound, "pay456"),
			wantIs:         ErrPaymentMethodNotFound,
			wantMessage:    "complete checkout: payment method not found: pay456",
			wantUnwrapping: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test error message
			if tt.err.Error() != tt.wantMessage {
				t.Errorf("error message = %q, want %q", tt.err.Error(), tt.wantMessage)
			}

			// Test error unwrapping with errors.Is
			if tt.wantUnwrapping {
				if !errors.Is(tt.err, tt.wantIs) {
					t.Errorf("errors.Is() = false, want true for %v", tt.wantIs)
				}
			}
		})
	}
}

func TestCLIError(t *testing.T) {
	tests := []struct {
		name        string
		cliErr      *CLIError
		wantMessage string
		wantCode    string
		wantUnwrap  error
	}{
		{
			name: "CLIError without wrapped error",
			cliErr: &CLIError{
				Code:    ErrCodeInvalidInput,
				Message: "invalid ASIN provided",
				Err:     nil,
			},
			wantMessage: "invalid ASIN provided",
			wantCode:    ErrCodeInvalidInput,
			wantUnwrap:  nil,
		},
		{
			name: "CLIError with wrapped sentinel error",
			cliErr: &CLIError{
				Code:    ErrCodeInvalidInput,
				Message: "validation failed",
				Err:     ErrInvalidASIN,
			},
			wantMessage: "validation failed: invalid ASIN",
			wantCode:    ErrCodeInvalidInput,
			wantUnwrap:  ErrInvalidASIN,
		},
		{
			name: "CLIError with wrapped fmt error",
			cliErr: &CLIError{
				Code:    ErrCodePurchaseFail,
				Message: "checkout failed",
				Err:     fmt.Errorf("add to cart: %w", ErrEmptyCart),
			},
			wantMessage: "checkout failed: add to cart: cart is empty",
			wantCode:    ErrCodePurchaseFail,
			wantUnwrap:  fmt.Errorf("add to cart: %w", ErrEmptyCart),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Error() method
			if tt.cliErr.Error() != tt.wantMessage {
				t.Errorf("Error() = %q, want %q", tt.cliErr.Error(), tt.wantMessage)
			}

			// Test Code field
			if tt.cliErr.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", tt.cliErr.Code, tt.wantCode)
			}

			// Test Unwrap() method
			unwrapped := tt.cliErr.Unwrap()
			if tt.wantUnwrap == nil {
				if unwrapped != nil {
					t.Errorf("Unwrap() = %v, want nil", unwrapped)
				}
			} else {
				if unwrapped == nil {
					t.Errorf("Unwrap() = nil, want non-nil")
				} else if unwrapped.Error() != tt.wantUnwrap.Error() {
					t.Errorf("Unwrap().Error() = %q, want %q", unwrapped.Error(), tt.wantUnwrap.Error())
				}
			}
		})
	}
}

func TestNewCLIError(t *testing.T) {
	baseErr := fmt.Errorf("network timeout")
	cliErr := NewCLIError(ErrCodeNetworkError, "failed to connect", baseErr)

	if cliErr.Code != ErrCodeNetworkError {
		t.Errorf("Code = %q, want %q", cliErr.Code, ErrCodeNetworkError)
	}

	if cliErr.Message != "failed to connect" {
		t.Errorf("Message = %q, want %q", cliErr.Message, "failed to connect")
	}

	if cliErr.Err != baseErr {
		t.Errorf("Err = %v, want %v", cliErr.Err, baseErr)
	}

	expectedMsg := "failed to connect: network timeout"
	if cliErr.Error() != expectedMsg {
		t.Errorf("Error() = %q, want %q", cliErr.Error(), expectedMsg)
	}
}

func TestErrorChaining(t *testing.T) {
	// Create a chain of wrapped errors
	baseErr := ErrInvalidASIN
	wrappedOnce := fmt.Errorf("validation failed: %w", baseErr)
	wrappedTwice := fmt.Errorf("add to cart: %w", wrappedOnce)
	wrappedThrice := fmt.Errorf("operation failed: %w", wrappedTwice)

	// Test that errors.Is can find the base error through multiple wraps
	if !errors.Is(wrappedThrice, ErrInvalidASIN) {
		t.Error("errors.Is() should find ErrInvalidASIN in chain")
	}

	// Test error message contains all context
	expectedMsg := "operation failed: add to cart: validation failed: invalid ASIN"
	if wrappedThrice.Error() != expectedMsg {
		t.Errorf("Error() = %q, want %q", wrappedThrice.Error(), expectedMsg)
	}

	// Test unwrapping step by step
	unwrapped1 := errors.Unwrap(wrappedThrice)
	if unwrapped1 == nil || unwrapped1.Error() != wrappedTwice.Error() {
		t.Error("First unwrap should return wrappedTwice")
	}

	unwrapped2 := errors.Unwrap(unwrapped1)
	if unwrapped2 == nil || unwrapped2.Error() != wrappedOnce.Error() {
		t.Error("Second unwrap should return wrappedOnce")
	}

	unwrapped3 := errors.Unwrap(unwrapped2)
	if unwrapped3 == nil || unwrapped3.Error() != baseErr.Error() {
		t.Error("Third unwrap should return baseErr")
	}
}

func TestCLIErrorWithErrorsIs(t *testing.T) {
	// Create a CLIError wrapping a sentinel error
	baseErr := ErrEmptyCart
	cliErr := NewCLIError(ErrCodePurchaseFail, "checkout failed", baseErr)

	// Test that errors.Is works with CLIError
	if !errors.Is(cliErr, ErrEmptyCart) {
		t.Error("errors.Is() should find ErrEmptyCart in CLIError")
	}

	// Test with multi-level wrapping
	wrappedCLI := fmt.Errorf("operation failed: %w", cliErr)
	if !errors.Is(wrappedCLI, ErrEmptyCart) {
		t.Error("errors.Is() should find ErrEmptyCart through fmt.Errorf and CLIError")
	}
}
