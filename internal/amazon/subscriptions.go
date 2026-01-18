package amazon

import (
	"fmt"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// subscriptionStore is an in-memory store for testing/development
// In production, this would be replaced with actual Amazon API calls
var subscriptionStore = map[string]*models.Subscription{
	"S01-1234567-8901234": {
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
	"S01-9876543-2109876": {
		SubscriptionID:  "S01-9876543-2109876",
		ASIN:            "B01SAMPLE",
		Title:           "Paper Towels 12 Pack",
		Price:           32.50,
		DiscountPercent: 10,
		FrequencyWeeks:  8,
		NextDelivery:    "2024-03-15",
		Status:          "active",
		Quantity:        2,
	},
}

// GetSubscription retrieves a specific subscription by ID
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) GetSubscription(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	// TODO: Implement actual Amazon Subscribe & Save API call
	// For now, retrieve from in-memory store
	subscription, exists := subscriptionStore[subscriptionID]
	if !exists {
		return nil, fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	// Return a copy to prevent external modifications
	result := *subscription
	return &result, nil
}

// GetSubscriptions retrieves all subscriptions for the user
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) GetSubscriptions() (*models.SubscriptionsResponse, error) {
	// TODO: Implement actual Amazon Subscribe & Save dashboard API call
	subscriptions := make([]models.Subscription, 0, len(subscriptionStore))
	for _, sub := range subscriptionStore {
		subscriptions = append(subscriptions, *sub)
	}

	return &models.SubscriptionsResponse{
		Subscriptions: subscriptions,
	}, nil
}

// SkipDelivery skips the next delivery for a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) SkipDelivery(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	// TODO: Implement actual Amazon skip delivery API call
	subscription, exists := subscriptionStore[subscriptionID]
	if !exists {
		return nil, fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	// Return a copy
	result := *subscription
	return &result, nil
}

// UpdateFrequency changes the delivery frequency of a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) UpdateFrequency(subscriptionID string, weeks int) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	if weeks <= 0 || weeks > 26 {
		return nil, fmt.Errorf("frequency must be between 1 and 26 weeks")
	}

	// TODO: Implement actual Amazon update frequency API call
	subscription, exists := subscriptionStore[subscriptionID]
	if !exists {
		return nil, fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	// Return a copy
	result := *subscription
	return &result, nil
}

// CancelSubscription cancels a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) CancelSubscription(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	// TODO: Implement actual Amazon cancel subscription API call
	subscription, exists := subscriptionStore[subscriptionID]
	if !exists {
		return nil, fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	// Return a copy with cancelled status
	result := *subscription
	result.Status = "cancelled"
	return &result, nil
}

// GetUpcomingDeliveries retrieves all upcoming deliveries across all subscriptions
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) GetUpcomingDeliveries() ([]models.UpcomingDelivery, error) {
	// TODO: Implement actual Amazon upcoming deliveries API call
	deliveries := make([]models.UpcomingDelivery, 0, len(subscriptionStore))
	for _, sub := range subscriptionStore {
		if sub.Status == "active" {
			deliveries = append(deliveries, models.UpcomingDelivery{
				SubscriptionID: sub.SubscriptionID,
				ASIN:           sub.ASIN,
				Title:          sub.Title,
				DeliveryDate:   sub.NextDelivery,
				Quantity:       sub.Quantity,
			})
		}
	}

	return deliveries, nil
}
