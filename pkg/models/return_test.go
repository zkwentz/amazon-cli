package models

import "testing"

func TestIsValidReturnReason(t *testing.T) {
	tests := []struct {
		name     string
		reason   string
		expected bool
	}{
		{
			name:     "valid reason - defective",
			reason:   "defective",
			expected: true,
		},
		{
			name:     "valid reason - wrong_item",
			reason:   "wrong_item",
			expected: true,
		},
		{
			name:     "valid reason - not_as_described",
			reason:   "not_as_described",
			expected: true,
		},
		{
			name:     "valid reason - no_longer_needed",
			reason:   "no_longer_needed",
			expected: true,
		},
		{
			name:     "valid reason - better_price",
			reason:   "better_price",
			expected: true,
		},
		{
			name:     "valid reason - other",
			reason:   "other",
			expected: true,
		},
		{
			name:     "invalid reason - empty string",
			reason:   "",
			expected: false,
		},
		{
			name:     "invalid reason - random string",
			reason:   "invalid_reason",
			expected: false,
		},
		{
			name:     "invalid reason - wrong case",
			reason:   "DEFECTIVE",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidReturnReason(tt.reason)
			if result != tt.expected {
				t.Errorf("IsValidReturnReason(%q) = %v, expected %v", tt.reason, result, tt.expected)
			}
		})
	}
}

func TestReturnableItemStructure(t *testing.T) {
	item := ReturnableItem{
		OrderID:      "123-4567890-1234567",
		ItemID:       "ITEM123",
		ASIN:         "B08N5WRWNW",
		Title:        "Test Product",
		Price:        29.99,
		PurchaseDate: "2024-01-15",
		ReturnWindow: "2024-02-15",
	}

	if item.OrderID != "123-4567890-1234567" {
		t.Errorf("Expected OrderID to be set correctly")
	}
	if item.Price != 29.99 {
		t.Errorf("Expected Price to be 29.99, got %f", item.Price)
	}
}

func TestReturnOptionStructure(t *testing.T) {
	option := ReturnOption{
		Method:          "UPS",
		Label:           "UPS Drop-off",
		DropoffLocation: "123 Main St",
		Fee:             0.0,
	}

	if option.Method != "UPS" {
		t.Errorf("Expected Method to be UPS")
	}
	if option.Fee != 0.0 {
		t.Errorf("Expected Fee to be 0.0, got %f", option.Fee)
	}
}

func TestReturnStructure(t *testing.T) {
	returnItem := Return{
		ReturnID:  "R123456789",
		OrderID:   "123-4567890-1234567",
		ItemID:    "ITEM123",
		Status:    "initiated",
		Reason:    "defective",
		CreatedAt: "2024-01-20T12:00:00Z",
	}

	if returnItem.Status != "initiated" {
		t.Errorf("Expected Status to be initiated")
	}
	if returnItem.Reason != "defective" {
		t.Errorf("Expected Reason to be defective")
	}
}

func TestReturnLabelStructure(t *testing.T) {
	label := ReturnLabel{
		URL:          "https://amazon.com/return-label/123",
		Carrier:      "UPS",
		Instructions: "Drop off at any UPS location",
	}

	if label.Carrier != "UPS" {
		t.Errorf("Expected Carrier to be UPS")
	}
	if label.URL == "" {
		t.Errorf("Expected URL to be set")
	}
}

func TestReturnReasonCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant ReturnReasonCode
		expected string
	}{
		{"defective", ReturnReasonDefective, "defective"},
		{"wrong_item", ReturnReasonWrongItem, "wrong_item"},
		{"not_as_described", ReturnReasonNotAsDescribed, "not_as_described"},
		{"no_longer_needed", ReturnReasonNoLongerNeeded, "no_longer_needed"},
		{"better_price", ReturnReasonBetterPrice, "better_price"},
		{"other", ReturnReasonOther, "other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("Expected %s to be %s, got %s", tt.name, tt.expected, string(tt.constant))
			}
		})
	}
}
