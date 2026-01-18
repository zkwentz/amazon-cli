package models

import "testing"

func TestIsValidReturnReason(t *testing.T) {
	tests := []struct {
		name     string
		reason   string
		expected bool
	}{
		{"Valid: defective", "defective", true},
		{"Valid: wrong_item", "wrong_item", true},
		{"Valid: not_as_described", "not_as_described", true},
		{"Valid: no_longer_needed", "no_longer_needed", true},
		{"Valid: better_price", "better_price", true},
		{"Valid: other", "other", true},
		{"Invalid: empty", "", false},
		{"Invalid: unknown", "unknown_reason", false},
		{"Invalid: wrong case", "DEFECTIVE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidReturnReason(tt.reason)
			if result != tt.expected {
				t.Errorf("IsValidReturnReason(%q) = %v, want %v", tt.reason, result, tt.expected)
			}
		})
	}
}

func TestValidReturnReasons(t *testing.T) {
	expectedReasons := []string{
		"defective",
		"wrong_item",
		"not_as_described",
		"no_longer_needed",
		"better_price",
		"other",
	}

	for _, reason := range expectedReasons {
		if _, ok := ValidReturnReasons[reason]; !ok {
			t.Errorf("ValidReturnReasons missing expected reason: %s", reason)
		}
	}

	if len(ValidReturnReasons) != len(expectedReasons) {
		t.Errorf("ValidReturnReasons has %d reasons, expected %d", len(ValidReturnReasons), len(expectedReasons))
	}
}
