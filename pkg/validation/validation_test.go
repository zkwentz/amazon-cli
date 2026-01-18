package validation

import (
	"testing"
)

func TestValidateOrderID(t *testing.T) {
	tests := []struct {
		name        string
		orderID     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid order ID",
			orderID:     "123-4567890-1234567",
			expectError: false,
		},
		{
			name:        "valid order ID with zeros",
			orderID:     "000-0000000-0000000",
			expectError: false,
		},
		{
			name:        "valid order ID with all nines",
			orderID:     "999-9999999-9999999",
			expectError: false,
		},
		{
			name:        "empty order ID",
			orderID:     "",
			expectError: true,
			errorMsg:    "order ID cannot be empty",
		},
		{
			name:        "too short - missing last segment",
			orderID:     "123-4567890",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "first segment too short",
			orderID:     "12-4567890-1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "first segment too long",
			orderID:     "1234-4567890-1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "second segment too short",
			orderID:     "123-456789-1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "second segment too long",
			orderID:     "123-45678901-1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "third segment too short",
			orderID:     "123-4567890-123456",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "third segment too long",
			orderID:     "123-4567890-12345678",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "missing hyphens",
			orderID:     "12345678901234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "contains letters",
			orderID:     "ABC-4567890-1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "contains special characters",
			orderID:     "123-456789@-1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "spaces instead of hyphens",
			orderID:     "123 4567890 1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "trailing hyphen",
			orderID:     "123-4567890-1234567-",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "leading hyphen",
			orderID:     "-123-4567890-1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "double hyphen",
			orderID:     "123--4567890-1234567",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "contains spaces",
			orderID:     "123-4567890-1234567 ",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
		{
			name:        "lowercase letters",
			orderID:     "abc-defghij-klmnopq",
			expectError: true,
			errorMsg:    "invalid order ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrderID(tt.orderID)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil for order ID: %s", tt.orderID)
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					// Check if error message contains the expected substring
					if len(tt.errorMsg) > 0 && !contains(err.Error(), tt.errorMsg) {
						t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v for order ID: %s", err, tt.orderID)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
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
