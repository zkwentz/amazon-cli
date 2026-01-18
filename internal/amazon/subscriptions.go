package amazon

import (
	"fmt"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// SubscriptionClient handles Amazon Subscribe & Save operations
type SubscriptionClient struct {
	// In a real implementation, this would contain HTTP client, auth, etc.
}

// NewSubscriptionClient creates a new subscription client
func NewSubscriptionClient() *SubscriptionClient {
	return &SubscriptionClient{}
}

// ValidFrequencies defines the valid subscription frequency intervals in weeks
var ValidFrequencies = []int{1, 2, 3, 4, 5, 6, 8, 10, 12, 16, 20, 24, 26}

// IsValidFrequency checks if the given frequency in weeks is valid
func IsValidFrequency(weeks int) bool {
	for _, valid := range ValidFrequencies {
		if weeks == valid {
			return true
		}
	}
	return false
}

// UpdateFrequency changes the delivery frequency for a subscription
// It validates that the frequency is within the allowed range (1-26 weeks)
// and matches Amazon's supported frequency intervals
func (c *SubscriptionClient) UpdateFrequency(subscriptionID string, weeks int) (*models.Subscription, error) {
	// Validate subscription ID
	if subscriptionID == "" {
		return nil, &models.CLIError{
			Code:    models.ErrorCodeInvalidInput,
			Message: "subscription ID cannot be empty",
			Details: map[string]interface{}{},
		}
	}

	// Validate frequency range (1-26 weeks as per PRD)
	if weeks < 1 {
		return nil, &models.CLIError{
			Code:    models.ErrorCodeInvalidInput,
			Message: "frequency must be at least 1 week",
			Details: map[string]interface{}{
				"provided_weeks": weeks,
				"minimum_weeks":  1,
			},
		}
	}

	if weeks > 26 {
		return nil, &models.CLIError{
			Code:    models.ErrorCodeInvalidInput,
			Message: "frequency cannot exceed 26 weeks",
			Details: map[string]interface{}{
				"provided_weeks": weeks,
				"maximum_weeks":  26,
			},
		}
	}

	// Validate that the frequency is one of Amazon's supported intervals
	if !IsValidFrequency(weeks) {
		return nil, &models.CLIError{
			Code:    models.ErrorCodeInvalidInput,
			Message: fmt.Sprintf("frequency must be one of the supported intervals: %v weeks", ValidFrequencies),
			Details: map[string]interface{}{
				"provided_weeks":     weeks,
				"valid_frequencies":  ValidFrequencies,
			},
		}
	}

	// In a real implementation, this would make an API call to Amazon
	// For now, return a mock successful response
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  weeks,
		NextDelivery:    "2024-02-01",
		Status:          "active",
		Quantity:        1,
	}, nil
}

// GetSubscriptions retrieves all active subscriptions
func (c *SubscriptionClient) GetSubscriptions() (*models.SubscriptionsResponse, error) {
	// Mock implementation
	return &models.SubscriptionsResponse{
		Subscriptions: []models.Subscription{
			{
				SubscriptionID:  "S01-1234567-8901234",
				ASIN:            "B00EXAMPLE",
				Title:           "Coffee Pods 100 Count",
				Price:           45.99,
				DiscountPercent: 15,
				FrequencyWeeks:  4,
				NextDelivery:    "2024-02-01",
				Status:          "active",
				Quantity:        1,
			},
		},
	}, nil
}

// GetSubscription retrieves details for a specific subscription
func (c *SubscriptionClient) GetSubscription(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, &models.CLIError{
			Code:    models.ErrorCodeInvalidInput,
			Message: "subscription ID cannot be empty",
			Details: map[string]interface{}{},
		}
	}

	// Mock implementation
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "2024-02-01",
		Status:          "active",
		Quantity:        1,
	}, nil
}

// SkipDelivery skips the next delivery for a subscription
func (c *SubscriptionClient) SkipDelivery(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, &models.CLIError{
			Code:    models.ErrorCodeInvalidInput,
			Message: "subscription ID cannot be empty",
			Details: map[string]interface{}{},
		}
	}

	// Mock implementation - returns subscription with updated next delivery date
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "2024-03-01",
		Status:          "active",
		Quantity:        1,
	}, nil
}

// CancelSubscription cancels a subscription
func (c *SubscriptionClient) CancelSubscription(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, &models.CLIError{
			Code:    models.ErrorCodeInvalidInput,
			Message: "subscription ID cannot be empty",
			Details: map[string]interface{}{},
		}
	}

	// Mock implementation - returns subscription with cancelled status
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "",
		Status:          "cancelled",
		Quantity:        1,
	}, nil
}

// GetUpcomingDeliveries retrieves all upcoming deliveries across subscriptions
func (c *SubscriptionClient) GetUpcomingDeliveries() ([]models.UpcomingDelivery, error) {
	// Mock implementation
	return []models.UpcomingDelivery{
		{
			SubscriptionID: "S01-1234567-8901234",
			ASIN:           "B00EXAMPLE",
			Title:          "Coffee Pods 100 Count",
			DeliveryDate:   "2024-02-01",
			Quantity:       1,
		},
	}, nil
}
