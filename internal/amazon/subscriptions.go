package amazon

import (
	"sort"
	"time"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// GetSubscriptions retrieves all active subscriptions
func (c *Client) GetSubscriptions() (*models.SubscriptionsResponse, error) {
	// Mock implementation - returns sample subscriptions
	subscriptions := []models.Subscription{
		{
			SubscriptionID:  "S01-1234567-8901234",
			ASIN:            "B00EXAMPLE1",
			Title:           "Coffee Pods 100 Count",
			Price:           45.99,
			DiscountPercent: 15,
			FrequencyWeeks:  4,
			NextDelivery:    "2024-02-01",
			Status:          "active",
			Quantity:        1,
		},
		{
			SubscriptionID:  "S01-2345678-9012345",
			ASIN:            "B00EXAMPLE2",
			Title:           "Paper Towels 12-Pack",
			Price:           29.99,
			DiscountPercent: 10,
			FrequencyWeeks:  8,
			NextDelivery:    "2024-03-15",
			Status:          "active",
			Quantity:        2,
		},
	}

	return &models.SubscriptionsResponse{
		Subscriptions: subscriptions,
	}, nil
}

// GetSubscription retrieves details for a specific subscription
func (c *Client) GetSubscription(subscriptionID string) (*models.Subscription, error) {
	// Mock implementation
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE1",
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
func (c *Client) SkipDelivery(subscriptionID string) (*models.Subscription, error) {
	// Mock implementation - would update next delivery date
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE1",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "2024-03-01", // Skipped one delivery
		Status:          "active",
		Quantity:        1,
	}, nil
}

// UpdateFrequency updates the delivery frequency for a subscription
func (c *Client) UpdateFrequency(subscriptionID string, weeks int) (*models.Subscription, error) {
	// Mock implementation
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE1",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  weeks,
		NextDelivery:    "2024-02-01",
		Status:          "active",
		Quantity:        1,
	}, nil
}

// CancelSubscription cancels a subscription
func (c *Client) CancelSubscription(subscriptionID string) (*models.Subscription, error) {
	// Mock implementation
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE1",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    "",
		Status:          "cancelled",
		Quantity:        1,
	}, nil
}

// GetUpcomingDeliveries retrieves all upcoming deliveries sorted by date
func (c *Client) GetUpcomingDeliveries() ([]models.UpcomingDelivery, error) {
	// Get all subscriptions
	resp, err := c.GetSubscriptions()
	if err != nil {
		return nil, err
	}

	// Convert subscriptions to upcoming deliveries
	deliveries := make([]models.UpcomingDelivery, 0, len(resp.Subscriptions))
	for _, sub := range resp.Subscriptions {
		// Only include active subscriptions with next delivery dates
		if sub.Status == "active" && sub.NextDelivery != "" {
			deliveries = append(deliveries, models.UpcomingDelivery{
				SubscriptionID: sub.SubscriptionID,
				ASIN:           sub.ASIN,
				Title:          sub.Title,
				DeliveryDate:   sub.NextDelivery,
				Quantity:       sub.Quantity,
			})
		}
	}

	// Sort deliveries by date (earliest first)
	sort.Slice(deliveries, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", deliveries[i].DeliveryDate)
		dateJ, errJ := time.Parse("2006-01-02", deliveries[j].DeliveryDate)

		// If parsing fails, maintain original order
		if errI != nil || errJ != nil {
			return false
		}

		return dateI.Before(dateJ)
	})

	return deliveries, nil
}
