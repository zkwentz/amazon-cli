package amazon

import (
	"testing"
	"time"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestValidateReturnEligibility_ItemNotReturnable(t *testing.T) {
	tests := []struct {
		name        string
		item        *models.ReturnableItem
		expectError bool
		errorCode   string
	}{
		{
			name: "item marked as non-returnable",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234567",
				ItemID:       "ITEM123",
				ASIN:         "B08N5WRWNW",
				Title:        "Non-returnable Product",
				Price:        29.99,
				PurchaseDate: "2024-01-01",
				ReturnWindow: "2024-02-01",
				Returnable:   false, // Item is not returnable
			},
			expectError: true,
			errorCode:   models.ErrCodeItemNotReturnable,
		},
		{
			name: "digital content not returnable",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234568",
				ItemID:       "ITEM124",
				ASIN:         "B08N5WRWNY",
				Title:        "Digital Download",
				Price:        9.99,
				PurchaseDate: "2024-01-01",
				ReturnWindow: time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				Returnable:   false, // Digital items typically not returnable
			},
			expectError: true,
			errorCode:   models.ErrCodeItemNotReturnable,
		},
		{
			name: "returnable item within window",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234569",
				ItemID:       "ITEM125",
				ASIN:         "B08N5WRWNZ",
				Title:        "Returnable Product",
				Price:        49.99,
				PurchaseDate: "2024-01-01",
				ReturnWindow: time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02"),
				Returnable:   true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnEligibility(tt.item)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				cliErr, ok := err.(*models.CLIError)
				if !ok {
					t.Errorf("Expected CLIError but got %T", err)
					return
				}

				if cliErr.Code != tt.errorCode {
					t.Errorf("Expected error code %s but got %s", tt.errorCode, cliErr.Code)
				}

				if cliErr.Details == nil {
					t.Errorf("Expected error details but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateReturnEligibility_ReturnWindowExpired(t *testing.T) {
	tests := []struct {
		name        string
		item        *models.ReturnableItem
		expectError bool
		errorCode   string
	}{
		{
			name: "return window expired by 1 day",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234570",
				ItemID:       "ITEM126",
				ASIN:         "B08N5WRWNA",
				Title:        "Product with Expired Return Window",
				Price:        39.99,
				PurchaseDate: "2023-12-01",
				ReturnWindow: time.Now().Add(-24 * time.Hour).Format("2006-01-02"), // Yesterday
				Returnable:   true,
			},
			expectError: true,
			errorCode:   models.ErrCodeReturnWindowExpired,
		},
		{
			name: "return window expired by 30 days",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234571",
				ItemID:       "ITEM127",
				ASIN:         "B08N5WRWNB",
				Title:        "Product with Long Expired Return Window",
				Price:        99.99,
				PurchaseDate: "2023-10-01",
				ReturnWindow: time.Now().Add(-30 * 24 * time.Hour).Format("2006-01-02"), // 30 days ago
				Returnable:   true,
			},
			expectError: true,
			errorCode:   models.ErrCodeReturnWindowExpired,
		},
		{
			name: "empty return window (expired)",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234572",
				ItemID:       "ITEM128",
				ASIN:         "B08N5WRWNC",
				Title:        "Product with Empty Return Window",
				Price:        19.99,
				PurchaseDate: "2023-11-01",
				ReturnWindow: "", // Empty return window should be considered expired
				Returnable:   true,
			},
			expectError: true,
			errorCode:   models.ErrCodeReturnWindowExpired,
		},
		{
			name: "return window still valid (7 days left)",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234573",
				ItemID:       "ITEM129",
				ASIN:         "B08N5WRWND",
				Title:        "Product with Valid Return Window",
				Price:        59.99,
				PurchaseDate: time.Now().Add(-23 * 24 * time.Hour).Format("2006-01-02"),
				ReturnWindow: time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02"), // 7 days from now
				Returnable:   true,
			},
			expectError: false,
		},
		{
			name: "return window valid (tomorrow)",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234574",
				ItemID:       "ITEM130",
				ASIN:         "B08N5WRWNE",
				Title:        "Product with Return Window Expiring Tomorrow",
				Price:        29.99,
				PurchaseDate: time.Now().Add(-30 * 24 * time.Hour).Format("2006-01-02"),
				ReturnWindow: time.Now().Add(24 * time.Hour).Format("2006-01-02"), // Tomorrow
				Returnable:   true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReturnEligibility(tt.item)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				cliErr, ok := err.(*models.CLIError)
				if !ok {
					t.Errorf("Expected CLIError but got %T", err)
					return
				}

				if cliErr.Code != tt.errorCode {
					t.Errorf("Expected error code %s but got %s", tt.errorCode, cliErr.Code)
				}

				if cliErr.Details == nil {
					t.Errorf("Expected error details but got nil")
					return
				}

				// Verify error details contain order_date and expiry_date
				if _, ok := cliErr.Details["order_date"]; !ok {
					t.Errorf("Expected error details to contain 'order_date'")
				}
				if _, ok := cliErr.Details["expiry_date"]; !ok {
					t.Errorf("Expected error details to contain 'expiry_date'")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestCreateReturn_WithErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		item        *models.ReturnableItem
		reason      string
		expectError bool
		errorCode   string
	}{
		{
			name: "successful return creation",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234575",
				ItemID:       "ITEM131",
				ASIN:         "B08N5WRWNF",
				Title:        "Valid Product",
				Price:        49.99,
				PurchaseDate: time.Now().Add(-10 * 24 * time.Hour).Format("2006-01-02"),
				ReturnWindow: time.Now().Add(20 * 24 * time.Hour).Format("2006-01-02"),
				Returnable:   true,
			},
			reason:      "defective",
			expectError: false,
		},
		{
			name: "return creation fails - item not returnable",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234576",
				ItemID:       "ITEM132",
				ASIN:         "B08N5WRWNG",
				Title:        "Non-returnable Product",
				Price:        29.99,
				PurchaseDate: time.Now().Add(-5 * 24 * time.Hour).Format("2006-01-02"),
				ReturnWindow: time.Now().Add(25 * 24 * time.Hour).Format("2006-01-02"),
				Returnable:   false, // Not returnable
			},
			reason:      "wrong_item",
			expectError: true,
			errorCode:   models.ErrCodeItemNotReturnable,
		},
		{
			name: "return creation fails - window expired",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234577",
				ItemID:       "ITEM133",
				ASIN:         "B08N5WRWNH",
				Title:        "Expired Return Window Product",
				Price:        79.99,
				PurchaseDate: time.Now().Add(-45 * 24 * time.Hour).Format("2006-01-02"),
				ReturnWindow: time.Now().Add(-5 * 24 * time.Hour).Format("2006-01-02"), // 5 days ago
				Returnable:   true,
			},
			reason:      "not_as_described",
			expectError: true,
			errorCode:   models.ErrCodeReturnWindowExpired,
		},
		{
			name: "return creation fails - invalid reason",
			item: &models.ReturnableItem{
				OrderID:      "123-4567890-1234578",
				ItemID:       "ITEM134",
				ASIN:         "B08N5WRWNI",
				Title:        "Valid Product",
				Price:        39.99,
				PurchaseDate: time.Now().Add(-5 * 24 * time.Hour).Format("2006-01-02"),
				ReturnWindow: time.Now().Add(25 * 24 * time.Hour).Format("2006-01-02"),
				Returnable:   true,
			},
			reason:      "invalid_reason", // Invalid reason code
			expectError: true,
			errorCode:   models.ErrCodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CreateReturn(tt.item, tt.reason)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				cliErr, ok := err.(*models.CLIError)
				if !ok {
					t.Errorf("Expected CLIError but got %T", err)
					return
				}

				if cliErr.Code != tt.errorCode {
					t.Errorf("Expected error code %s but got %s", tt.errorCode, cliErr.Code)
				}

				if result != nil {
					t.Errorf("Expected nil result on error but got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
					return
				}

				if result == nil {
					t.Errorf("Expected return result but got nil")
					return
				}

				if result.OrderID != tt.item.OrderID {
					t.Errorf("Expected OrderID %s but got %s", tt.item.OrderID, result.OrderID)
				}

				if result.ItemID != tt.item.ItemID {
					t.Errorf("Expected ItemID %s but got %s", tt.item.ItemID, result.ItemID)
				}

				if result.Reason != tt.reason {
					t.Errorf("Expected Reason %s but got %s", tt.reason, result.Reason)
				}

				if result.Status != "initiated" {
					t.Errorf("Expected Status 'initiated' but got %s", result.Status)
				}
			}
		})
	}
}

func TestReturnableItem_IsReturnWindowExpired(t *testing.T) {
	tests := []struct {
		name     string
		item     *models.ReturnableItem
		expected bool
	}{
		{
			name: "expired window",
			item: &models.ReturnableItem{
				ReturnWindow: time.Now().Add(-24 * time.Hour).Format("2006-01-02"),
			},
			expected: true,
		},
		{
			name: "valid window",
			item: &models.ReturnableItem{
				ReturnWindow: time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			},
			expected: false,
		},
		{
			name: "empty window",
			item: &models.ReturnableItem{
				ReturnWindow: "",
			},
			expected: true,
		},
		{
			name: "invalid date format",
			item: &models.ReturnableItem{
				ReturnWindow: "invalid-date",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsReturnWindowExpired()
			if result != tt.expected {
				t.Errorf("Expected %v but got %v", tt.expected, result)
			}
		})
	}
}

func TestReturnableItem_IsReturnable(t *testing.T) {
	tests := []struct {
		name     string
		item     *models.ReturnableItem
		expected bool
	}{
		{
			name: "returnable and within window",
			item: &models.ReturnableItem{
				Returnable:   true,
				ReturnWindow: time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02"),
			},
			expected: true,
		},
		{
			name: "returnable but window expired",
			item: &models.ReturnableItem{
				Returnable:   true,
				ReturnWindow: time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02"),
			},
			expected: false,
		},
		{
			name: "not returnable but within window",
			item: &models.ReturnableItem{
				Returnable:   false,
				ReturnWindow: time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02"),
			},
			expected: false,
		},
		{
			name: "not returnable and window expired",
			item: &models.ReturnableItem{
				Returnable:   false,
				ReturnWindow: time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsReturnable()
			if result != tt.expected {
				t.Errorf("Expected %v but got %v", tt.expected, result)
			}
		})
	}
}
