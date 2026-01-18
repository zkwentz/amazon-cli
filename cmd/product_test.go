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
			asin:  "B0123456XY",
			valid: true,
		},
		{
			name:  "invalid ASIN - too short",
			asin:  "B08N5WRW",
			valid: false,
		},
		{
			name:  "invalid ASIN - too long",
			asin:  "B08N5WRWNW1",
			valid: false,
		},
		{
			name:  "invalid ASIN - lowercase",
			asin:  "b08n5wrwnw",
			valid: false,
		},
		{
			name:  "invalid ASIN - special characters",
			asin:  "B08N5WRW-W",
			valid: false,
		},
		{
			name:  "empty ASIN",
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
