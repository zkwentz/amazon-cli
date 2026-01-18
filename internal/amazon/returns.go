package amazon

import (
	"fmt"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// ValidateReturnEligibility checks if an item can be returned
func ValidateReturnEligibility(item *models.ReturnableItem) error {
	// Check if item is marked as non-returnable
	if !item.Returnable {
		return models.NewItemNotReturnableError("Item is marked as non-returnable by Amazon")
	}

	// Check if return window has expired
	if item.IsReturnWindowExpired() {
		return models.NewReturnWindowExpiredError(item.PurchaseDate, item.ReturnWindow)
	}

	return nil
}

// CreateReturn initiates a return for an item
func CreateReturn(item *models.ReturnableItem, reason string) (*models.Return, error) {
	// Validate return eligibility
	if err := ValidateReturnEligibility(item); err != nil {
		return nil, err
	}

	// Validate reason code
	validReasons := map[string]bool{
		"defective":        true,
		"wrong_item":       true,
		"not_as_described": true,
		"no_longer_needed": true,
		"better_price":     true,
		"other":            true,
	}

	if !validReasons[reason] {
		return nil, models.NewCLIError(
			models.ErrCodeInvalidInput,
			fmt.Sprintf("Invalid return reason: %s", reason),
			map[string]interface{}{
				"valid_reasons": []string{
					"defective",
					"wrong_item",
					"not_as_described",
					"no_longer_needed",
					"better_price",
					"other",
				},
			},
		)
	}

	// In a real implementation, this would make API calls to Amazon
	// For now, we return a mock successful return
	return &models.Return{
		ReturnID:  fmt.Sprintf("RET-%s-%s", item.OrderID, item.ItemID),
		OrderID:   item.OrderID,
		ItemID:    item.ItemID,
		Status:    "initiated",
		Reason:    reason,
		CreatedAt: "2024-01-18T00:00:00Z",
	}, nil
}
