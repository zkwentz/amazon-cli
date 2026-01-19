package amazon

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// validReturnReasons contains the list of allowed return reasons
var validReturnReasons = map[string]bool{
	"defective":        true,
	"wrong_item":       true,
	"not_as_described": true,
	"no_longer_needed": true,
	"better_price":     true,
	"other":            true,
}

// CreateReturn creates a return request for an order item
func (c *Client) CreateReturn(orderID, itemID, reason string) (*models.Return, error) {
	// Validate orderID is not empty
	if orderID == "" {
		return nil, fmt.Errorf("order ID cannot be empty")
	}

	// Validate itemID is not empty
	if itemID == "" {
		return nil, fmt.Errorf("item ID cannot be empty")
	}

	// Validate reason is not empty
	if reason == "" {
		return nil, fmt.Errorf("reason cannot be empty")
	}

	// Validate reason is in the allowed list
	if !validReturnReasons[reason] {
		return nil, fmt.Errorf("invalid return reason: %s (allowed: defective, wrong_item, not_as_described, no_longer_needed, better_price, other)", reason)
	}

	// Generate a unique return ID
	returnID := fmt.Sprintf("RET-%s", uuid.New().String())

	// Create the return object with mock data
	ret := &models.Return{
		ReturnID:  returnID,
		OrderID:   orderID,
		ItemID:    itemID,
		Reason:    reason,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	return ret, nil
}
