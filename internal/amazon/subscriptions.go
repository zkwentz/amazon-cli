package amazon

import (
	"sort"
	"time"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// GetSubscriptions retrieves all active and paused Subscribe & Save subscriptions
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) GetSubscriptions() (*models.SubscriptionsResponse, error) {
	// TODO: Implement actual Amazon Subscribe & Save API call
	// For now, return mock subscriptions for testing
	subscriptions := []models.Subscription{
		{
			SubscriptionID:  "S01-1234567-8901234",
			ASIN:            "B00EXAMPLE1",
			Title:           "Coffee Pods 100 Count",
			Price:           45.99,
			DiscountPercent: 15,
			FrequencyWeeks:  4,
			NextDelivery:    time.Now().AddDate(0, 0, 14).Format("2006-01-02"),
			Status:          "active",
			Quantity:        1,
		},
		{
			SubscriptionID:  "S01-2345678-9012345",
			ASIN:            "B00EXAMPLE2",
			Title:           "Paper Towels 12 Pack",
			Price:           28.50,
			DiscountPercent: 10,
			FrequencyWeeks:  8,
			NextDelivery:    time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
			Status:          "active",
			Quantity:        2,
		},
		{
			SubscriptionID:  "S01-3456789-0123456",
			ASIN:            "B00EXAMPLE3",
			Title:           "Laundry Detergent 150oz",
			Price:           19.99,
			DiscountPercent: 5,
			FrequencyWeeks:  12,
			NextDelivery:    time.Now().AddDate(0, 0, 45).Format("2006-01-02"),
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
	// TODO: Implement actual Amazon Subscribe & Save API call
	// For now, return a mock subscription
	return &models.Subscription{
		SubscriptionID:  subscriptionID,
		ASIN:            "B00EXAMPLE1",
		Title:           "Coffee Pods 100 Count",
		Price:           45.99,
		DiscountPercent: 15,
		FrequencyWeeks:  4,
		NextDelivery:    time.Now().AddDate(0, 0, 14).Format("2006-01-02"),
		Status:          "active",
		Quantity:        1,
	}, nil
}

// SkipDelivery skips the next delivery for a subscription
// This is a placeholder implementation that will be expanded with actual Amazon API calls
func (c *Client) SkipDelivery(subscriptionID string) (*models.Subscription, error) {
	// TODO: Implement actual Amazon Subscribe & Save API call
	// For now, return the subscription with updated next delivery date
	subscription, err := c.GetSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Parse current next delivery date
	currentDelivery, err := time.Parse("2006-01-02", subscription.NextDelivery)
	if err != nil {
		return nil, err
	}

	// Add the frequency weeks to skip to next delivery
	daysToAdd := subscription.FrequencyWeeks * 7
	subscription.NextDelivery = currentDelivery.AddDate(0, 0, daysToAdd).Format("2006-01-02")

	return subscription, nil
}

// GetUpcomingDeliveries retrieves all upcoming deliveries across all subscriptions
// Returns a slice of UpcomingDelivery sorted by delivery date
func (c *Client) GetUpcomingDeliveries() ([]models.UpcomingDelivery, error) {
	// Step 1: Get all subscriptions
	subscriptionsResp, err := c.GetSubscriptions()
	if err != nil {
		return nil, err
	}

	// Step 2: Build upcoming deliveries from active subscriptions
	var upcomingDeliveries []models.UpcomingDelivery
	for _, sub := range subscriptionsResp.Subscriptions {
		// Only include active subscriptions
		if sub.Status == "active" {
			upcomingDeliveries = append(upcomingDeliveries, models.UpcomingDelivery{
				SubscriptionID: sub.SubscriptionID,
				ASIN:           sub.ASIN,
				Title:          sub.Title,
				DeliveryDate:   sub.NextDelivery,
				Quantity:       sub.Quantity,
			})
		}
	}

	// Step 3: Sort by delivery date (earliest first)
	sort.Slice(upcomingDeliveries, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", upcomingDeliveries[i].DeliveryDate)
		dateJ, errJ := time.Parse("2006-01-02", upcomingDeliveries[j].DeliveryDate)

		// If either date can't be parsed, put it at the end
		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}

		return dateI.Before(dateJ)
	})

	return upcomingDeliveries, nil
}
