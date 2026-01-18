package validation

import (
	"testing"
)

func TestValidatePriceRange(t *testing.T) {
	tests := []struct {
		name        string
		minPrice    float64
		maxPrice    float64
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid price range",
			minPrice:    10.0,
			maxPrice:    100.0,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "valid price range with small difference",
			minPrice:    9.99,
			maxPrice:    10.00,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "valid price range with large values",
			minPrice:    500.0,
			maxPrice:    1000.0,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "negative min price should fail",
			minPrice:    -10.0,
			maxPrice:    100.0,
			wantErr:     true,
			errContains: "min price must be positive",
		},
		{
			name:        "negative max price should fail",
			minPrice:    10.0,
			maxPrice:    -100.0,
			wantErr:     true,
			errContains: "max price must be positive",
		},
		{
			name:        "both prices negative should fail",
			minPrice:    -50.0,
			maxPrice:    -10.0,
			wantErr:     true,
			errContains: "min price must be positive",
		},
		{
			name:        "min price equal to max price should fail",
			minPrice:    50.0,
			maxPrice:    50.0,
			wantErr:     true,
			errContains: "min price (50.00) must be less than max price (50.00)",
		},
		{
			name:        "min price greater than max price should fail",
			minPrice:    100.0,
			maxPrice:    50.0,
			wantErr:     true,
			errContains: "min price (100.00) must be less than max price (50.00)",
		},
		{
			name:        "zero min price with positive max should succeed",
			minPrice:    0.0,
			maxPrice:    100.0,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "very small positive values should succeed",
			minPrice:    0.01,
			maxPrice:    0.02,
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "zero min and zero max should fail",
			minPrice:    0.0,
			maxPrice:    0.0,
			wantErr:     true,
			errContains: "min price (0.00) must be less than max price (0.00)",
		},
		{
			name:        "decimal prices should work",
			minPrice:    19.99,
			maxPrice:    99.99,
			wantErr:     false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePriceRange(tt.minPrice, tt.maxPrice)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidatePriceRange() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidatePriceRange() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidatePriceRange() unexpected error: %v", err)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
