package validation

import (
	"testing"
)

func TestValidateASIN(t *testing.T) {
	tests := []struct {
		name    string
		asin    string
		wantErr bool
	}{
		{
			name:    "valid ASIN with letters and numbers",
			asin:    "B08N5WRWN1",
			wantErr: false,
		},
		{
			name:    "valid ASIN all uppercase",
			asin:    "ABCDEFGHIJ",
			wantErr: false,
		},
		{
			name:    "valid ASIN all digits",
			asin:    "1234567890",
			wantErr: false,
		},
		{
			name:    "valid ASIN mixed case",
			asin:    "AbCd123456",
			wantErr: false,
		},
		{
			name:    "invalid ASIN too short",
			asin:    "B08N5WRWN",
			wantErr: true,
		},
		{
			name:    "invalid ASIN too long",
			asin:    "B08N5WRWN12",
			wantErr: true,
		},
		{
			name:    "invalid ASIN empty string",
			asin:    "",
			wantErr: true,
		},
		{
			name:    "invalid ASIN with special characters",
			asin:    "B08N5-RWWN",
			wantErr: true,
		},
		{
			name:    "invalid ASIN with spaces",
			asin:    "B08N5 WRWN",
			wantErr: true,
		},
		{
			name:    "invalid ASIN with underscore",
			asin:    "B08N5_WRWN",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateASIN(tt.asin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateASIN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOrderID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid order ID",
			id:      "123-4567890-1234567",
			wantErr: false,
		},
		{
			name:    "valid order ID all zeros",
			id:      "000-0000000-0000000",
			wantErr: false,
		},
		{
			name:    "valid order ID all nines",
			id:      "999-9999999-9999999",
			wantErr: false,
		},
		{
			name:    "invalid order ID too few digits in first segment",
			id:      "12-4567890-1234567",
			wantErr: true,
		},
		{
			name:    "invalid order ID too many digits in first segment",
			id:      "1234-4567890-1234567",
			wantErr: true,
		},
		{
			name:    "invalid order ID too few digits in second segment",
			id:      "123-456789-1234567",
			wantErr: true,
		},
		{
			name:    "invalid order ID too many digits in second segment",
			id:      "123-45678901-1234567",
			wantErr: true,
		},
		{
			name:    "invalid order ID too few digits in third segment",
			id:      "123-4567890-123456",
			wantErr: true,
		},
		{
			name:    "invalid order ID too many digits in third segment",
			id:      "123-4567890-12345678",
			wantErr: true,
		},
		{
			name:    "invalid order ID missing dashes",
			id:      "12345678901234567",
			wantErr: true,
		},
		{
			name:    "invalid order ID with letters",
			id:      "ABC-4567890-1234567",
			wantErr: true,
		},
		{
			name:    "invalid order ID empty string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "invalid order ID with spaces",
			id:      "123-4567890-123456 ",
			wantErr: true,
		},
		{
			name:    "invalid order ID wrong delimiter",
			id:      "123_4567890_1234567",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrderID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrderID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		name    string
		qty     int
		wantErr bool
	}{
		{
			name:    "valid quantity minimum",
			qty:     1,
			wantErr: false,
		},
		{
			name:    "valid quantity maximum",
			qty:     999,
			wantErr: false,
		},
		{
			name:    "valid quantity middle range",
			qty:     50,
			wantErr: false,
		},
		{
			name:    "valid quantity 100",
			qty:     100,
			wantErr: false,
		},
		{
			name:    "invalid quantity zero",
			qty:     0,
			wantErr: true,
		},
		{
			name:    "invalid quantity negative",
			qty:     -1,
			wantErr: true,
		},
		{
			name:    "invalid quantity negative large",
			qty:     -100,
			wantErr: true,
		},
		{
			name:    "invalid quantity too large",
			qty:     1000,
			wantErr: true,
		},
		{
			name:    "invalid quantity way too large",
			qty:     10000,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuantity(tt.qty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuantity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePriceRange(t *testing.T) {
	tests := []struct {
		name    string
		min     float64
		max     float64
		wantErr bool
	}{
		{
			name:    "valid price range",
			min:     0.0,
			max:     100.0,
			wantErr: false,
		},
		{
			name:    "valid price range decimal values",
			min:     9.99,
			max:     99.99,
			wantErr: false,
		},
		{
			name:    "valid price range minimum zero",
			min:     0.0,
			max:     0.01,
			wantErr: false,
		},
		{
			name:    "valid price range large values",
			min:     100.0,
			max:     1000.0,
			wantErr: false,
		},
		{
			name:    "valid price range small difference",
			min:     10.0,
			max:     10.01,
			wantErr: false,
		},
		{
			name:    "invalid price range negative minimum",
			min:     -10.0,
			max:     100.0,
			wantErr: true,
		},
		{
			name:    "invalid price range max equals min",
			min:     50.0,
			max:     50.0,
			wantErr: true,
		},
		{
			name:    "invalid price range max less than min",
			min:     100.0,
			max:     50.0,
			wantErr: true,
		},
		{
			name:    "invalid price range both zero",
			min:     0.0,
			max:     0.0,
			wantErr: true,
		},
		{
			name:    "invalid price range negative max",
			min:     0.0,
			max:     -10.0,
			wantErr: true,
		},
		{
			name:    "invalid price range both negative",
			min:     -50.0,
			max:     -10.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePriceRange(tt.min, tt.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePriceRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
