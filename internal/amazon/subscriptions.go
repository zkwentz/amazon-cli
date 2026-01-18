package amazon

import (
	"fmt"
	"sort"
	"time"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// GetSubscriptions retrieves all active and paused subscriptions
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) GetSubscriptions() (*models.SubscriptionsResponse, error) {
	// TODO: Implement actual Amazon Subscribe & Save dashboard API call
	// For now, return mock subscriptions for testing/development
	subscriptions := []models.Subscription{
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
	}

	return &models.SubscriptionsResponse{
		Subscriptions: subscriptions,
	}, nil
}

// GetSubscription retrieves details for a specific subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) GetSubscription(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	// TODO: Implement actual Amazon subscription details API call
	// For now, return a mock subscription
	subscription := &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "2024-02-01",
		Status:          "active",
		Quantity:        1,
	}

	return subscription, nil
}

// SkipDelivery skips the next delivery for a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) SkipDelivery(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	// TODO: Implement actual Amazon skip delivery API call
	// For now, get the subscription and calculate next delivery date
	subscription, err := c.GetSubscription(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Parse the current next delivery date
	currentDate, err := time.Parse("2006-01-02", subscription.NextDelivery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse next delivery date: %w", err)
	}

	// Calculate new next delivery date (current + frequency)
	newDate := currentDate.AddDate(0, 0, subscription.FrequencyWeeks*7)
	subscription.NextDelivery = newDate.Format("2006-01-02")

	return subscription, nil
}

// UpdateFrequency changes the delivery frequency for a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) UpdateFrequency(subscriptionID string, weeks int) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	// Validate frequency is reasonable (1-26 weeks)
	if weeks < 1 || weeks > 26 {
		return nil, fmt.Errorf("frequency must be between 1 and 26 weeks, got %d", weeks)
	}

	// TODO: Implement actual Amazon frequency change API call
	// For now, get the subscription and update frequency
	subscription, err := c.GetSubscription(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	subscription.FrequencyWeeks = weeks

	return subscription, nil
}

// CancelSubscription cancels a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) CancelSubscription(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscriptionID cannot be empty")
	}

	// TODO: Implement actual Amazon subscription cancellation API call
	// For now, get the subscription and mark as cancelled
	subscription, err := c.GetSubscription(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	subscription.Status = "cancelled"

	return subscription, nil
}

// GetUpcomingDeliveries retrieves upcoming deliveries across all subscriptions
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) GetUpcomingDeliveries() ([]models.UpcomingDelivery, error) {
	// TODO: Implement actual Amazon upcoming deliveries API call
	// For now, get all subscriptions and extract delivery info
	subscriptionsResp, err := c.GetSubscriptions()
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	var deliveries []models.UpcomingDelivery
	for _, sub := range subscriptionsResp.Subscriptions {
		// Only include active subscriptions
		if sub.Status != "active" {
			continue
		}

		delivery := models.UpcomingDelivery{
			SubscriptionID: sub.SubscriptionID,
			ASIN:           sub.ASIN,
			Title:          sub.Title,
			DeliveryDate:   sub.NextDelivery,
			Quantity:       sub.Quantity,
		}
		deliveries = append(deliveries, delivery)
	}

	// Sort deliveries by date (earliest first)
	sort.Slice(deliveries, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", deliveries[i].DeliveryDate)
		dateJ, errJ := time.Parse("2006-01-02", deliveries[j].DeliveryDate)
		if errI != nil || errJ != nil {
			return false
		}
		return dateI.Before(dateJ)
	})

	return deliveries, nil
}
