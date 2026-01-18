package cmd

import (
	"testing"
)

func TestIsValidASIN(t *testing.T) {
	tests := []struct {
		name  string
		asin  string
		valid bool
	}{
		{
			name:  "valid ASIN",
			asin:  "B08N5WRWNW",
			valid: true,
		},
		{
			name:  "valid ASIN with numbers",
			asin:  "B012345678",
			valid: true,
		},
		{
			name:  "invalid - too short",
			asin:  "B08N5WRW",
			valid: false,
		},
		{
			name:  "invalid - too long",
			asin:  "B08N5WRWNW1",
			valid: false,
		},
		{
			name:  "invalid - lowercase",
			asin:  "b08n5wrwnw",
			valid: false,
		},
		{
			name:  "invalid - special characters",
			asin:  "B08N5WRW-W",
			valid: false,
		},
		{
			name:  "invalid - empty",
			asin:  "",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidASIN(tt.asin)
			if result != tt.valid {
				t.Errorf("isValidASIN(%q) = %v, want %v", tt.asin, result, tt.valid)
			}
		})
	}
}

func TestCartCommands(t *testing.T) {
	// Test that all cart commands are properly initialized
	if cartCmd == nil {
		t.Error("cartCmd is nil")
	}
	if cartAddCmd == nil {
		t.Error("cartAddCmd is nil")
	}
	if cartListCmd == nil {
		t.Error("cartListCmd is nil")
	}
	if cartRemoveCmd == nil {
		t.Error("cartRemoveCmd is nil")
	}
	if cartClearCmd == nil {
		t.Error("cartClearCmd is nil")
	}
	if cartCheckoutCmd == nil {
		t.Error("cartCheckoutCmd is nil")
	}

	// Test that cart command has correct number of subcommands
	if len(cartCmd.Commands()) != 5 {
		t.Errorf("cartCmd has %d subcommands, want 5", len(cartCmd.Commands()))
	}
}

func TestCartAddCmdFlags(t *testing.T) {
	// Test that quantity flag exists
	flag := cartAddCmd.Flags().Lookup("quantity")
	if flag == nil {
		t.Error("quantity flag not found")
	}
	if flag.DefValue != "1" {
		t.Errorf("quantity default value = %q, want %q", flag.DefValue, "1")
	}
}

func TestCartClearCmdFlags(t *testing.T) {
	// Test that confirm flag exists
	flag := cartClearCmd.Flags().Lookup("confirm")
	if flag == nil {
		t.Error("confirm flag not found")
	}
	if flag.DefValue != "false" {
		t.Errorf("confirm default value = %q, want %q", flag.DefValue, "false")
	}
}

func TestCartCheckoutCmdFlags(t *testing.T) {
	// Test that all checkout flags exist
	confirmFlag := cartCheckoutCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Error("confirm flag not found")
	}

	addressFlag := cartCheckoutCmd.Flags().Lookup("address-id")
	if addressFlag == nil {
		t.Error("address-id flag not found")
	}

	paymentFlag := cartCheckoutCmd.Flags().Lookup("payment-id")
	if paymentFlag == nil {
		t.Error("payment-id flag not found")
	}
}
