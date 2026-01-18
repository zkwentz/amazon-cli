package validation

import (
	"testing"
)

func TestValidateASIN(t *testing.T) {
	tests := []struct {
		name      string
		asin      string
		wantError bool
	}{
		// Valid ASINs
		{
			name:      "valid ASIN with all uppercase letters",
			asin:      "ABCDEFGHIJ",
			wantError: false,
		},
		{
			name:      "valid ASIN with all digits",
			asin:      "1234567890",
			wantError: false,
		},
		{
			name:      "valid ASIN with mixed letters and digits",
			asin:      "B08N5WRWNW",
			wantError: false,
		},
		{
			name:      "valid ASIN starting with B",
			asin:      "B00EXAMPLE",
			wantError: false,
		},
		{
			name:      "valid ASIN with consecutive letters",
			asin:      "ABC1234DEF",
			wantError: false,
		},
		// Invalid ASINs - empty or wrong length
		{
			name:      "empty ASIN",
			asin:      "",
			wantError: true,
		},
		{
			name:      "ASIN too short (9 characters)",
			asin:      "B08N5WRWN",
			wantError: true,
		},
		{
			name:      "ASIN too long (11 characters)",
			asin:      "B08N5WRWNWX",
			wantError: true,
		},
		{
			name:      "ASIN way too short (5 characters)",
			asin:      "B08N5",
			wantError: true,
		},
		{
			name:      "ASIN way too long (15 characters)",
			asin:      "B08N5WRWNWEXTRA",
			wantError: true,
		},
		// Invalid ASINs - wrong character types
		{
			name:      "ASIN with lowercase letters",
			asin:      "b08n5wrwnw",
			wantError: true,
		},
		{
			name:      "ASIN with mixed case",
			asin:      "B08n5WRWNw",
			wantError: true,
		},
		{
			name:      "ASIN with special characters (hyphen)",
			asin:      "B08-5WRWNW",
			wantError: true,
		},
		{
			name:      "ASIN with special characters (underscore)",
			asin:      "B08_5WRWNW",
			wantError: true,
		},
		{
			name:      "ASIN with space",
			asin:      "B08 5WRWNW",
			wantError: true,
		},
		{
			name:      "ASIN with special character at start",
			asin:      "@B08N5WRWN",
			wantError: true,
		},
		{
			name:      "ASIN with special character at end",
			asin:      "B08N5WRWN!",
			wantError: true,
		},
		{
			name:      "ASIN with unicode characters",
			asin:      "B08N5WRWNÉ",
			wantError: true,
		},
		{
			name:      "ASIN with only special characters",
			asin:      "!@#$%^&*()",
			wantError: true,
		},
		// Edge cases
		{
			name:      "ASIN with null byte",
			asin:      "B08N5WRW\x00W",
			wantError: true,
		},
		{
			name:      "ASIN with tab",
			asin:      "B08N5WR\tNW",
			wantError: true,
		},
		{
			name:      "ASIN with newline",
			asin:      "B08N5WR\nNW",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateASIN(tt.asin)
			if tt.wantError && err == nil {
				t.Errorf("ValidateASIN(%q) expected error, got nil", tt.asin)
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateASIN(%q) expected no error, got %v", tt.asin, err)
			}
		})
	}
}

func TestIsValidASIN(t *testing.T) {
	tests := []struct {
		name  string
		asin  string
		valid bool
	}{
		{
			name:  "valid ASIN returns true",
			asin:  "B08N5WRWNW",
			valid: true,
		},
		{
			name:  "empty ASIN returns false",
			asin:  "",
			valid: false,
		},
		{
			name:  "short ASIN returns false",
			asin:  "B08N5",
			valid: false,
		},
		{
			name:  "long ASIN returns false",
			asin:  "B08N5WRWNWX",
			valid: false,
		},
		{
			name:  "lowercase ASIN returns false",
			asin:  "b08n5wrwnw",
			valid: false,
		},
		{
			name:  "ASIN with special chars returns false",
			asin:  "B08-5WRWNW",
			valid: false,
		},
		{
			name:  "all digits ASIN returns true",
			asin:  "1234567890",
			valid: true,
		},
		{
			name:  "all letters ASIN returns true",
			asin:  "ABCDEFGHIJ",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidASIN(tt.asin)
			if got != tt.valid {
				t.Errorf("IsValidASIN(%q) = %v, want %v", tt.asin, got, tt.valid)
			}
		})
	}
}

func TestValidateASINErrorMessages(t *testing.T) {
	tests := []struct {
		name            string
		asin            string
		wantErrContains string
	}{
		{
			name:            "empty ASIN error message",
			asin:            "",
			wantErrContains: "cannot be empty",
		},
		{
			name:            "short ASIN error message",
			asin:            "ABC",
			wantErrContains: "must be exactly 10 characters long",
		},
		{
			name:            "long ASIN error message",
			asin:            "ABCDEFGHIJKLM",
			wantErrContains: "must be exactly 10 characters long",
		},
		{
			name:            "lowercase ASIN error message",
			asin:            "abcdefghij",
			wantErrContains: "must contain only uppercase letters and digits",
		},
		{
			name:            "special char ASIN error message",
			asin:            "ABC-DEFGHI",
			wantErrContains: "must contain only uppercase letters and digits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateASIN(tt.asin)
			if err == nil {
				t.Fatalf("ValidateASIN(%q) expected error, got nil", tt.asin)
			}
			if !contains(err.Error(), tt.wantErrContains) {
				t.Errorf("ValidateASIN(%q) error = %q, want error containing %q", tt.asin, err.Error(), tt.wantErrContains)
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

// Benchmark tests
func BenchmarkValidateASIN(b *testing.B) {
	asin := "B08N5WRWNW"
	for i := 0; i < b.N; i++ {
		ValidateASIN(asin)
	}
}

func BenchmarkIsValidASIN(b *testing.B) {
	asin := "B08N5WRWNW"
	for i := 0; i < b.N; i++ {
		IsValidASIN(asin)
	}
}

func BenchmarkValidateASINInvalid(b *testing.B) {
	asin := "invalid"
	for i := 0; i < b.N; i++ {
		ValidateASIN(asin)
	}
}
