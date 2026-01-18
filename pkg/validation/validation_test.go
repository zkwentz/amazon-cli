package validation

import (
	"strings"
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
			name:    "valid ASIN",
			asin:    "B08N5WRWNW",
			wantErr: false,
		},
		{
			name:    "valid ASIN all numbers",
			asin:    "1234567890",
			wantErr: false,
		},
		{
			name:    "valid ASIN all letters",
			asin:    "ABCDEFGHIJ",
			wantErr: false,
		},
		{
			name:    "empty ASIN",
			asin:    "",
			wantErr: true,
			errMsg:  "ASIN cannot be empty",
		},
		{
			name:    "too short ASIN",
			asin:    "B08N5WRW",
			wantErr: true,
			errMsg:  "invalid ASIN format",
		},
		{
			name:    "too long ASIN",
			asin:    "B08N5WRWNW1",
			wantErr: true,
			errMsg:  "invalid ASIN format",
		},
		{
			name:    "lowercase ASIN",
			asin:    "b08n5wrwnw",
			wantErr: true,
			errMsg:  "invalid ASIN format",
		},
		{
			name:    "ASIN with special characters",
			asin:    "B08N5WRW-W",
			wantErr: true,
			errMsg:  "invalid ASIN format",
		},
		{
			name:    "ASIN with spaces",
			asin:    "B08N5 RWNW",
			wantErr: true,
			errMsg:  "invalid ASIN format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateASIN(tt.asin)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateASIN() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateASIN() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateASIN() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid quantity 1",
			quantity: 1,
			wantErr:  false,
		},
		{
			name:     "valid quantity 10",
			quantity: 10,
			wantErr:  false,
		},
		{
			name:     "valid quantity max",
			quantity: 999,
			wantErr:  false,
		},
		{
			name:     "zero quantity",
			quantity: 0,
			wantErr:  true,
			errMsg:   "quantity must be at least",
		},
		{
			name:     "negative quantity",
			quantity: -1,
			wantErr:  true,
			errMsg:   "quantity must be at least",
		},
		{
			name:     "quantity exceeds max",
			quantity: 1000,
			wantErr:  true,
			errMsg:   "quantity cannot exceed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuantity(tt.quantity)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateQuantity() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateQuantity() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateQuantity() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateOrderID(t *testing.T) {
	tests := []struct {
		name    string
		orderID string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid order ID",
			orderID: "123-4567890-1234567",
			wantErr: false,
		},
		{
			name:    "valid order ID all zeros",
			orderID: "000-0000000-0000000",
			wantErr: false,
		},
		{
			name:    "empty order ID",
			orderID: "",
			wantErr: true,
			errMsg:  "order ID cannot be empty",
		},
		{
			name:    "invalid format - missing dashes",
			orderID: "12345678901234567",
			wantErr: true,
			errMsg:  "invalid order ID format",
		},
		{
			name:    "invalid format - wrong segment length",
			orderID: "12-4567890-1234567",
			wantErr: true,
			errMsg:  "invalid order ID format",
		},
		{
			name:    "invalid format - contains letters",
			orderID: "ABC-4567890-1234567",
			wantErr: true,
			errMsg:  "invalid order ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrderID(tt.orderID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateOrderID() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateOrderID() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateOrderID() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateSubscriptionID(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "valid subscription ID",
			subscriptionID: "S01-1234567-8901234",
			wantErr:        false,
		},
		{
			name:           "valid subscription ID S99",
			subscriptionID: "S99-0000000-0000000",
			wantErr:        false,
		},
		{
			name:           "empty subscription ID",
			subscriptionID: "",
			wantErr:        true,
			errMsg:         "subscription ID cannot be empty",
		},
		{
			name:           "invalid format - missing S prefix",
			subscriptionID: "001-1234567-8901234",
			wantErr:        true,
			errMsg:         "invalid subscription ID format",
		},
		{
			name:           "invalid format - wrong prefix format",
			subscriptionID: "S1-1234567-8901234",
			wantErr:        true,
			errMsg:         "invalid subscription ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSubscriptionID(tt.subscriptionID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateSubscriptionID() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateSubscriptionID() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateSubscriptionID() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateAddressID(t *testing.T) {
	tests := []struct {
		name      string
		addressID string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid address ID",
			addressID: "addr123",
			wantErr:   false,
		},
		{
			name:      "empty address ID",
			addressID: "",
			wantErr:   true,
			errMsg:    "address ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddressID(tt.addressID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAddressID() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateAddressID() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateAddressID() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidatePaymentID(t *testing.T) {
	tests := []struct {
		name      string
		paymentID string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid payment ID",
			paymentID: "pay123",
			wantErr:   false,
		},
		{
			name:      "empty payment ID",
			paymentID: "",
			wantErr:   true,
			errMsg:    "payment ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePaymentID(tt.paymentID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePaymentID() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidatePaymentID() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePaymentID() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidatePriceRange(t *testing.T) {
	tests := []struct {
		name     string
		minPrice float64
		maxPrice float64
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid price range",
			minPrice: 10.00,
			maxPrice: 50.00,
			wantErr:  false,
		},
		{
			name:     "both zero",
			minPrice: 0,
			maxPrice: 0,
			wantErr:  false,
		},
		{
			name:     "min negative",
			minPrice: -10.00,
			maxPrice: 50.00,
			wantErr:  true,
			errMsg:   "minimum price cannot be negative",
		},
		{
			name:     "max negative",
			minPrice: 10.00,
			maxPrice: -50.00,
			wantErr:  true,
			errMsg:   "maximum price cannot be negative",
		},
		{
			name:     "min greater than max",
			minPrice: 100.00,
			maxPrice: 50.00,
			wantErr:  true,
			errMsg:   "minimum price cannot be greater than maximum price",
		},
		{
			name:     "min exceeds max allowed",
			minPrice: 9999999.99,
			maxPrice: 0,
			wantErr:  true,
			errMsg:   "minimum price exceeds maximum allowed value",
		},
		{
			name:     "max exceeds max allowed",
			minPrice: 0,
			maxPrice: 9999999.99,
			wantErr:  true,
			errMsg:   "maximum price exceeds maximum allowed value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePriceRange(tt.minPrice, tt.maxPrice)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePriceRange() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidatePriceRange() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePriceRange() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateReturnReason(t *testing.T) {
	tests := []struct {
		name    string
		reason  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid reason defective",
			reason:  "defective",
			wantErr: false,
		},
		{
			name:    "valid reason wrong_item",
			reason:  "wrong_item",
			wantErr: false,
		},
		{
			name:    "valid reason not_as_described",
			reason:  "not_as_described",
			wantErr: false,
		},
		{
			name:    "valid reason no_longer_needed",
			reason:  "no_longer_needed",
			wantErr: false,
		},
		{
			name:    "valid reason better_price",
			reason:  "better_price",
			wantErr: false,
		},
		{
			name:    "valid reason other",
			reason:  "other",
			wantErr: false,
		},
		{
			name:    "empty reason",
			reason:  "",
			wantErr: true,
			errMsg:  "return reason cannot be empty",
		},
		{
			name:    "invalid reason",
			reason:  "just_because",
			wantErr: true,
			errMsg:  "invalid return reason",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnReason(tt.reason)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateReturnReason() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateReturnReason() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateReturnReason() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateFrequencyWeeks(t *testing.T) {
	tests := []struct {
		name    string
		weeks   int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid 1 week",
			weeks:   1,
			wantErr: false,
		},
		{
			name:    "valid 4 weeks",
			weeks:   4,
			wantErr: false,
		},
		{
			name:    "valid 26 weeks",
			weeks:   26,
			wantErr: false,
		},
		{
			name:    "zero weeks",
			weeks:   0,
			wantErr: true,
			errMsg:  "frequency must be at least 1 week",
		},
		{
			name:    "negative weeks",
			weeks:   -1,
			wantErr: true,
			errMsg:  "frequency must be at least 1 week",
		},
		{
			name:    "exceeds max weeks",
			weeks:   27,
			wantErr: true,
			errMsg:  "frequency cannot exceed 26 weeks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFrequencyWeeks(tt.weeks)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateFrequencyWeeks() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateFrequencyWeeks() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateFrequencyWeeks() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateItemID(t *testing.T) {
	tests := []struct {
		name    string
		itemID  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid item ID",
			itemID:  "ITEM123",
			wantErr: false,
		},
		{
			name:    "empty item ID",
			itemID:  "",
			wantErr: true,
			errMsg:  "item ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateItemID(tt.itemID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateItemID() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateItemID() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateItemID() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateReturnID(t *testing.T) {
	tests := []struct {
		name     string
		returnID string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid return ID",
			returnID: "RET123",
			wantErr:  false,
		},
		{
			name:     "empty return ID",
			returnID: "",
			wantErr:  true,
			errMsg:   "return ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnID(tt.returnID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateReturnID() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateReturnID() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateReturnID() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateSearchQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid search query",
			query:   "wireless headphones",
			wantErr: false,
		},
		{
			name:    "empty search query",
			query:   "",
			wantErr: true,
			errMsg:  "search query cannot be empty",
		},
		{
			name:    "query too long",
			query:   strings.Repeat("a", 501),
			wantErr: true,
			errMsg:  "search query too long",
		},
		{
			name:    "query at max length",
			query:   strings.Repeat("a", 500),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSearchQuery(tt.query)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateSearchQuery() expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateSearchQuery() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateSearchQuery() unexpected error: %v", err)
				}
			}
		})
	}
}
