package models

import (
	"strings"
	"testing"
)

func TestValidateReturnReason(t *testing.T) {
	tests := []struct {
		name      string
		reason    string
		wantError bool
	}{
		{
			name:      "valid defective",
			reason:    "defective",
			wantError: false,
		},
		{
			name:      "valid wrong_item",
			reason:    "wrong_item",
			wantError: false,
		},
		{
			name:      "valid not_as_described",
			reason:    "not_as_described",
			wantError: false,
		},
		{
			name:      "valid no_longer_needed",
			reason:    "no_longer_needed",
			wantError: false,
		},
		{
			name:      "valid better_price",
			reason:    "better_price",
			wantError: false,
		},
		{
			name:      "valid other",
			reason:    "other",
			wantError: false,
		},
		{
			name:      "invalid reason - empty",
			reason:    "",
			wantError: true,
		},
		{
			name:      "invalid reason - unknown",
			reason:    "invalid_reason",
			wantError: true,
		},
		{
			name:      "invalid reason - typo",
			reason:    "deffective",
			wantError: true,
		},
		{
			name:      "invalid reason - spaces only",
			reason:    "   ",
			wantError: true,
		},
		{
			name:      "invalid reason - partial match",
			reason:    "wrong",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnReason(tt.reason)
			if tt.wantError && err == nil {
				t.Errorf("ValidateReturnReason(%q) expected error but got nil", tt.reason)
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateReturnReason(%q) unexpected error: %v", tt.reason, err)
			}
		})
	}
}

func TestValidateReturnReason_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name   string
		reason string
	}{
		{"uppercase", "DEFECTIVE"},
		{"mixed case", "Defective"},
		{"uppercase with underscores", "WRONG_ITEM"},
		{"mixed case with underscores", "Wrong_Item"},
		{"uppercase other", "OTHER"},
		{"mixed case better price", "Better_Price"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnReason(tt.reason)
			if err != nil {
				t.Errorf("ValidateReturnReason(%q) should be case-insensitive, got error: %v", tt.reason, err)
			}
		})
	}
}

func TestValidateReturnReason_WithWhitespace(t *testing.T) {
	tests := []struct {
		name   string
		reason string
	}{
		{"leading spaces", "  defective"},
		{"trailing spaces", "defective  "},
		{"both sides", "  defective  "},
		{"leading tab", "\tdefective"},
		{"trailing newline", "defective\n"},
		{"mixed whitespace", " \tdefective \n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnReason(tt.reason)
			if err != nil {
				t.Errorf("ValidateReturnReason(%q) should trim whitespace, got error: %v", tt.reason, err)
			}
		})
	}
}

func TestValidateReturnReason_ErrorMessage(t *testing.T) {
	err := ValidateReturnReason("invalid_code")
	if err == nil {
		t.Fatal("expected error for invalid reason code")
	}

	errMsg := err.Error()
	expectedParts := []string{
		"invalid return reason",
		"invalid_code",
		"defective",
		"wrong_item",
		"not_as_described",
		"no_longer_needed",
		"better_price",
		"other",
	}

	for _, part := range expectedParts {
		if !strings.Contains(errMsg, part) {
			t.Errorf("error message should contain %q, got: %s", part, errMsg)
		}
	}
}

func TestIsValidReturnReason(t *testing.T) {
	tests := []struct {
		reason string
		want   bool
	}{
		{"defective", true},
		{"wrong_item", true},
		{"not_as_described", true},
		{"no_longer_needed", true},
		{"better_price", true},
		{"other", true},
		{"DEFECTIVE", true},
		{"  defective  ", true},
		{"invalid", false},
		{"", false},
		{"wrong", false},
	}

	for _, tt := range tests {
		t.Run(tt.reason, func(t *testing.T) {
			got := IsValidReturnReason(tt.reason)
			if got != tt.want {
				t.Errorf("IsValidReturnReason(%q) = %v, want %v", tt.reason, got, tt.want)
			}
		})
	}
}

func TestGetReturnReasonDescription(t *testing.T) {
	tests := []struct {
		reason      ReturnReasonCode
		wantContain string
	}{
		{ReasonDefective, "defective"},
		{ReasonWrongItem, "wrong item"},
		{ReasonNotAsDescribed, "not as described"},
		{ReasonNoLongerNeeded, "no longer needed"},
		{ReasonBetterPrice, "better price"},
		{ReasonOther, "other"},
	}

	for _, tt := range tests {
		t.Run(string(tt.reason), func(t *testing.T) {
			desc := GetReturnReasonDescription(tt.reason)
			if desc == "" {
				t.Errorf("GetReturnReasonDescription(%q) returned empty string", tt.reason)
			}
			lowerDesc := strings.ToLower(desc)
			if !strings.Contains(lowerDesc, tt.wantContain) {
				t.Errorf("GetReturnReasonDescription(%q) = %q, should contain %q", tt.reason, desc, tt.wantContain)
			}
		})
	}
}

func TestGetReturnReasonDescription_Unknown(t *testing.T) {
	desc := GetReturnReasonDescription("unknown_code")
	if desc == "" {
		t.Error("GetReturnReasonDescription should return non-empty string for unknown reason")
	}
	if !strings.Contains(strings.ToLower(desc), "unknown") {
		t.Errorf("GetReturnReasonDescription for unknown code should contain 'unknown', got: %s", desc)
	}
}

func TestNormalizeReturnReason(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      ReturnReasonCode
		wantError bool
	}{
		{
			name:      "valid lowercase",
			input:     "defective",
			want:      ReasonDefective,
			wantError: false,
		},
		{
			name:      "valid uppercase",
			input:     "DEFECTIVE",
			want:      ReasonDefective,
			wantError: false,
		},
		{
			name:      "valid mixed case",
			input:     "Defective",
			want:      ReasonDefective,
			wantError: false,
		},
		{
			name:      "with leading spaces",
			input:     "  wrong_item",
			want:      ReasonWrongItem,
			wantError: false,
		},
		{
			name:      "with trailing spaces",
			input:     "other  ",
			want:      ReasonOther,
			wantError: false,
		},
		{
			name:      "with whitespace both sides",
			input:     "  not_as_described  ",
			want:      ReasonNotAsDescribed,
			wantError: false,
		},
		{
			name:      "invalid reason",
			input:     "invalid",
			want:      "",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeReturnReason(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("NormalizeReturnReason(%q) expected error but got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("NormalizeReturnReason(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeReturnReason(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestAllReturnReasons(t *testing.T) {
	// Verify all expected reasons are present
	expectedReasons := []ReturnReasonCode{
		ReasonDefective,
		ReasonWrongItem,
		ReasonNotAsDescribed,
		ReasonNoLongerNeeded,
		ReasonBetterPrice,
		ReasonOther,
	}

	if len(AllReturnReasons) != len(expectedReasons) {
		t.Errorf("AllReturnReasons length = %d, want %d", len(AllReturnReasons), len(expectedReasons))
	}

	// Verify each expected reason exists in AllReturnReasons
	reasonMap := make(map[ReturnReasonCode]bool)
	for _, reason := range AllReturnReasons {
		reasonMap[reason] = true
	}

	for _, expected := range expectedReasons {
		if !reasonMap[expected] {
			t.Errorf("AllReturnReasons missing expected reason: %s", expected)
		}
	}
}

func TestReturnReasonCode_Constants(t *testing.T) {
	tests := []struct {
		constant ReturnReasonCode
		value    string
	}{
		{ReasonDefective, "defective"},
		{ReasonWrongItem, "wrong_item"},
		{ReasonNotAsDescribed, "not_as_described"},
		{ReasonNoLongerNeeded, "no_longer_needed"},
		{ReasonBetterPrice, "better_price"},
		{ReasonOther, "other"},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if string(tt.constant) != tt.value {
				t.Errorf("constant value mismatch: got %q, want %q", string(tt.constant), tt.value)
			}
		})
	}
}

func TestReturnableItem_Structure(t *testing.T) {
	// Test that ReturnableItem struct can be created and has correct fields
	item := ReturnableItem{
		OrderID:      "123-4567890-1234567",
		ItemID:       "item123",
		ASIN:         "B08N5WRWNW",
		Title:        "Test Product",
		Price:        29.99,
		PurchaseDate: "2024-01-15",
		ReturnWindow: "30 days",
	}

	if item.OrderID != "123-4567890-1234567" {
		t.Errorf("OrderID mismatch")
	}
	if item.ASIN != "B08N5WRWNW" {
		t.Errorf("ASIN mismatch")
	}
	if item.Price != 29.99 {
		t.Errorf("Price mismatch")
	}
}

func TestReturn_Structure(t *testing.T) {
	// Test that Return struct can be created with valid reason code
	ret := Return{
		ReturnID:  "RET123",
		OrderID:   "123-4567890-1234567",
		ItemID:    "item123",
		Status:    "initiated",
		Reason:    ReasonDefective,
		CreatedAt: "2024-01-20T12:00:00Z",
	}

	if ret.Reason != ReasonDefective {
		t.Errorf("Reason mismatch: got %q, want %q", ret.Reason, ReasonDefective)
	}
	if string(ret.Reason) != "defective" {
		t.Errorf("Reason string value mismatch")
	}
}

func TestReturnOption_Structure(t *testing.T) {
	option := ReturnOption{
		Method:          "UPS",
		Label:           "UPS Drop-off",
		DropoffLocation: "123 Main St",
		Fee:             0.0,
	}

	if option.Method != "UPS" {
		t.Errorf("Method mismatch")
	}
	if option.Fee != 0.0 {
		t.Errorf("Fee should be 0.0")
	}
}

func TestReturnLabel_Structure(t *testing.T) {
	label := ReturnLabel{
		URL:          "https://example.com/label.pdf",
		Carrier:      "UPS",
		Instructions: "Drop off at any UPS location",
	}

	if !strings.HasPrefix(label.URL, "https://") {
		t.Errorf("URL should start with https://")
	}
}

// Benchmark tests
func BenchmarkValidateReturnReason_Valid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ValidateReturnReason("defective")
	}
}

func BenchmarkValidateReturnReason_Invalid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ValidateReturnReason("invalid_reason")
	}
}

func BenchmarkIsValidReturnReason(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsValidReturnReason("defective")
	}
}

func BenchmarkNormalizeReturnReason(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NormalizeReturnReason("  DEFECTIVE  ")
	}
}
