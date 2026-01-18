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
	// TODO: Implement actual Amazon Subscribe & Save dashboard scraping/API call
	// For now, return mock data
	subscriptions := []models.Subscription{
		{
			SubscriptionID:  "S01-1234567-8901234",
			ASIN:            "B00EXAMPLE01",
			Title:           "Coffee Pods 100 Count",
			Price:           45.99,
			DiscountPercent: 15,
			FrequencyWeeks:  4,
			NextDelivery:    "2026-02-15",
			Status:          "active",
			Quantity:        1,
		},
		{
			SubscriptionID:  "S01-2345678-9012345",
			ASIN:            "B00EXAMPLE02",
			Title:           "Paper Towels 12 Pack",
			Price:           32.50,
			DiscountPercent: 10,
			FrequencyWeeks:  8,
			NextDelivery:    "2026-03-01",
			Status:          "active",
			Quantity:        2,
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
		return nil, fmt.Errorf("subscription ID cannot be empty")
	}

	// TODO: Implement actual Amazon subscription details retrieval
	// For now, return mock data based on the subscriptionID
	subscription := &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE01",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "2026-02-15",
		Status:          "active",
		Quantity:        1,
	}

	return subscription, nil
}

// SkipDelivery skips the next delivery for a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) SkipDelivery(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscription ID cannot be empty")
	}

	// TODO: Implement actual Amazon skip delivery API call
	// Get the subscription first
	subscription, err := c.GetSubscription(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Parse the next delivery date and add the frequency to get the new date
	nextDate, err := time.Parse("2006-01-02", subscription.NextDelivery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse next delivery date: %w", err)
	}

	// Add the frequency in weeks to skip to the next period
	newNextDate := nextDate.AddDate(0, 0, subscription.FrequencyWeeks*7)
	subscription.NextDelivery = newNextDate.Format("2006-01-02")

	return subscription, nil
}

// UpdateFrequency changes the delivery frequency for a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) UpdateFrequency(subscriptionID string, weeks int) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscription ID cannot be empty")
	}

	if weeks <= 0 || weeks > 26 {
		return nil, fmt.Errorf("frequency must be between 1 and 26 weeks")
	}

	// TODO: Implement actual Amazon frequency update API call
	subscription, err := c.GetSubscription(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Update the frequency
	subscription.FrequencyWeeks = weeks

	return subscription, nil
}

// CancelSubscription cancels a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) CancelSubscription(subscriptionID string) (*models.Subscription, error) {
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscription ID cannot be empty")
	}

	// TODO: Implement actual Amazon cancel subscription API call
	subscription, err := c.GetSubscription(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Update the status to cancelled
	subscription.Status = "cancelled"

	return subscription, nil
}

// GetUpcomingDeliveries retrieves all upcoming deliveries across all subscriptions
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) GetUpcomingDeliveries() ([]models.UpcomingDelivery, error) {
	// TODO: Implement actual Amazon upcoming deliveries retrieval
	// Get all subscriptions
	response, err := c.GetSubscriptions()
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	// Convert to upcoming deliveries
	deliveries := make([]models.UpcomingDelivery, 0, len(response.Subscriptions))
	for _, sub := range response.Subscriptions {
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

	// Sort by delivery date
	sort.Slice(deliveries, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", deliveries[i].DeliveryDate)
		dateJ, _ := time.Parse("2006-01-02", deliveries[j].DeliveryDate)
		return dateI.Before(dateJ)
	})

	return deliveries, nil
}
