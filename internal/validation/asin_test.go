package validation

import (
	"testing"
)

func TestValidateASIN(t *testing.T) {
	tests := []struct {
		name    string
		asin    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid ASIN with letters and numbers",
			asin:    "B08N5WRWNW",
			wantErr: false,
		},
		{
			name:    "valid ASIN all letters",
			asin:    "ABCDEFGHIJ",
			wantErr: false,
		},
		{
			name:    "valid ASIN all numbers",
			asin:    "1234567890",
			wantErr: false,
		},
		{
			name:    "valid ASIN mixed",
			asin:    "B00EXAMPLE",
			wantErr: false,
		},
		{
			name:    "valid ASIN starting with number",
			asin:    "0123456789",
			wantErr: false,
		},
		{
			name:    "empty ASIN",
			asin:    "",
			wantErr: true,
			errMsg:  "ASIN cannot be empty",
		},
		{
			name:    "ASIN too short",
			asin:    "B08N5WRWN",
			wantErr: true,
			errMsg:  "ASIN must be exactly 10 characters long",
		},
		{
			name:    "ASIN too long",
			asin:    "B08N5WRWNWX",
			wantErr: true,
			errMsg:  "ASIN must be exactly 10 characters long",
		},
		{
			name:    "ASIN with lowercase letters",
			asin:    "b08n5wrwnw",
			wantErr: true,
			errMsg:  "ASIN must contain only uppercase letters and digits",
		},
		{
			name:    "ASIN with special characters",
			asin:    "B08N5WRW-W",
			wantErr: true,
			errMsg:  "ASIN must contain only uppercase letters and digits",
		},
		{
			name:    "ASIN with spaces",
			asin:    "B08N5WRW W",
			wantErr: true,
			errMsg:  "ASIN must contain only uppercase letters and digits",
		},
		{
			name:    "ASIN with underscore",
			asin:    "B08N5WRW_W",
			wantErr: true,
			errMsg:  "ASIN must contain only uppercase letters and digits",
		},
		{
			name:    "ASIN with only 5 characters",
			asin:    "B08N5",
			wantErr: true,
			errMsg:  "ASIN must be exactly 10 characters long",
		},
		{
			name:    "ASIN with 15 characters",
			asin:    "B08N5WRWNWEXTRA",
			wantErr: true,
			errMsg:  "ASIN must be exactly 10 characters long",
		},
		{
			name:    "ASIN with mixed case",
			asin:    "B08n5WrWnW",
			wantErr: true,
			errMsg:  "ASIN must contain only uppercase letters and digits",
		},
		{
			name:    "ASIN with leading space",
			asin:    " B08N5WRWNW",
			wantErr: true,
			errMsg:  "ASIN must be exactly 10 characters long",
		},
		{
			name:    "ASIN with trailing space",
			asin:    "B08N5WRWNW ",
			wantErr: true,
			errMsg:  "ASIN must be exactly 10 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateASIN(tt.asin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateASIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg && len(err.Error()) > 0 {
					// Check if error message contains the expected substring
					contains := false
					if len(tt.errMsg) > 0 {
						for i := 0; i <= len(err.Error())-len(tt.errMsg); i++ {
							if err.Error()[i:i+len(tt.errMsg)] == tt.errMsg {
								contains = true
								break
							}
						}
					}
					if !contains {
						t.Errorf("ValidateASIN() error message = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				}
			}
		})
	}
}

func TestIsValidASIN(t *testing.T) {
	tests := []struct {
		name string
		asin string
		want bool
	}{
		{
			name: "valid ASIN",
			asin: "B08N5WRWNW",
			want: true,
		},
		{
			name: "valid ASIN all letters",
			asin: "ABCDEFGHIJ",
			want: true,
		},
		{
			name: "valid ASIN all numbers",
			asin: "1234567890",
			want: true,
		},
		{
			name: "invalid ASIN - too short",
			asin: "B08N5WRW",
			want: false,
		},
		{
			name: "invalid ASIN - too long",
			asin: "B08N5WRWNWX",
			want: false,
		},
		{
			name: "invalid ASIN - lowercase",
			asin: "b08n5wrwnw",
			want: false,
		},
		{
			name: "invalid ASIN - special chars",
			asin: "B08N5WRW-W",
			want: false,
		},
		{
			name: "invalid ASIN - empty",
			asin: "",
			want: false,
		},
		{
			name: "invalid ASIN - with space",
			asin: "B08N5WRW W",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidASIN(tt.asin); got != tt.want {
				t.Errorf("IsValidASIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark tests to ensure validation is performant
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
