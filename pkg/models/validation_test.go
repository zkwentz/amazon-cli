package models

import (
	"testing"
)

func TestValidateSubscriptionID(t *testing.T) {
	tests := []struct {
		name          string
		subscriptionID string
		wantErr       bool
		errContains   string
	}{
		{
			name:          "valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:       false,
		},
		{
			name:          "valid subscription ID with different digits",
			subscriptionID: "S99-9999999-9999999",
			wantErr:       false,
		},
		{
			name:          "valid subscription ID with zeros",
			subscriptionID: "S00-0000000-0000000",
			wantErr:       false,
		},
		{
			name:          "empty subscription ID",
			subscriptionID: "",
			wantErr:       true,
			errContains:   "cannot be empty",
		},
		{
			name:          "missing S prefix",
			subscriptionID: "01-1234567-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "lowercase s prefix",
			subscriptionID: "s01-1234567-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "too few digits after S",
			subscriptionID: "S1-1234567-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "too many digits after S",
			subscriptionID: "S001-1234567-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "too few digits in first segment",
			subscriptionID: "S01-123456-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "too many digits in first segment",
			subscriptionID: "S01-12345678-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "too few digits in second segment",
			subscriptionID: "S01-1234567-890123",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "too many digits in second segment",
			subscriptionID: "S01-1234567-89012345",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "missing first hyphen",
			subscriptionID: "S011234567-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "missing second hyphen",
			subscriptionID: "S01-12345678901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "missing both hyphens",
			subscriptionID: "S0112345678901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "contains letters in digit sections",
			subscriptionID: "S01-123456A-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "contains special characters",
			subscriptionID: "S01-1234567-890123!",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "extra characters at end",
			subscriptionID: "S01-1234567-8901234X",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "extra characters at start",
			subscriptionID: "XS01-1234567-8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "spaces in ID",
			subscriptionID: "S01 1234567 8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
		{
			name:          "underscores instead of hyphens",
			subscriptionID: "S01_1234567_8901234",
			wantErr:       true,
			errContains:   "invalid subscription ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSubscriptionID(tt.subscriptionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateSubscriptionID() expected error but got nil")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateSubscriptionID() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateSubscriptionID() unexpected error = %v", err)
				}
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
