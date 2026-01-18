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
			name:  "valid ASIN with all caps and numbers",
			asin:  "B08N5WRWNW",
			valid: true,
		},
		{
			name:  "valid ASIN all numbers",
			asin:  "1234567890",
			valid: true,
		},
		{
			name:  "valid ASIN all letters",
			asin:  "ABCDEFGHIJ",
			valid: true,
		},
		{
			name:  "invalid ASIN too short",
			asin:  "B08N5WRW",
			valid: false,
		},
		{
			name:  "invalid ASIN too long",
			asin:  "B08N5WRWNWX",
			valid: false,
		},
		{
			name:  "invalid ASIN with lowercase",
			asin:  "b08n5wrwnw",
			valid: false,
		},
		{
			name:  "invalid ASIN with special characters",
			asin:  "B08N5WRW-W",
			valid: false,
		},
		{
			name:  "invalid empty ASIN",
			asin:  "",
			valid: false,
		},
		{
			name:  "invalid ASIN with spaces",
			asin:  "B08N5 WRWN",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidASIN(tt.asin)
			if got != tt.valid {
				t.Errorf("isValidASIN(%q) = %v, want %v", tt.asin, got, tt.valid)
			}
		})
	}
}

func TestCartAddCommand(t *testing.T) {
	// Test that the cart command exists and has proper structure
	if cartCmd == nil {
		t.Fatal("cartCmd should not be nil")
	}

	if cartCmd.Use != "cart" {
		t.Errorf("cartCmd.Use = %q, want %q", cartCmd.Use, "cart")
	}

	// Test that cart add subcommand exists
	if cartAddCmd == nil {
		t.Fatal("cartAddCmd should not be nil")
	}

	if cartAddCmd.Use != "add <asin>" {
		t.Errorf("cartAddCmd.Use = %q, want %q", cartAddCmd.Use, "add <asin>")
	}

	// Check that quantity flag exists
	flag := cartAddCmd.Flags().Lookup("quantity")
	if flag == nil {
		t.Fatal("quantity flag should exist")
	}

	if flag.Shorthand != "n" {
		t.Errorf("quantity flag shorthand = %q, want %q", flag.Shorthand, "n")
	}

	if flag.DefValue != "1" {
		t.Errorf("quantity flag default = %q, want %q", flag.DefValue, "1")
	}
}
