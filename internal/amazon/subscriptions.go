package amazon

import (
	"fmt"
	"time"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// GetSubscriptions retrieves all active subscriptions for the user
func (c *Client) GetSubscriptions() (*models.SubscriptionList, error) {
	// TODO: Implement actual Amazon API call to get subscriptions
	// For now, return mock data
	subscriptions := []models.Subscription{
		{
			ID:             "sub001",
			ASIN:           "B08XYZ1234",
			Title:          "Coffee Pods - Subscribe & Save",
			Price:          24.99,
			Discount:       5.0,
			FrequencyWeeks: 4,
			NextDelivery:   time.Now().AddDate(0, 0, 14),
			Status:         "active",
			Quantity:       1,
		},
		{
			ID:             "sub002",
			ASIN:           "B09ABC5678",
			Title:          "Paper Towels - 12 Pack",
			Price:          29.99,
			Discount:       10.0,
			FrequencyWeeks: 8,
			NextDelivery:   time.Now().AddDate(0, 0, 21),
			Status:         "active",
			Quantity:       2,
		},
	}

	return &models.SubscriptionList{
		Subscriptions: subscriptions,
		TotalCount:    len(subscriptions),
	}, nil
}

// SkipDelivery skips the next delivery for a subscription by advancing NextDelivery by FrequencyWeeks
func (c *Client) SkipDelivery(id string) (*models.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("subscription ID cannot be empty")
	}

	// TODO: Implement actual Amazon API call to skip delivery
	// For now, simulate the operation with mock data
	subscription := &models.Subscription{
		ID:             id,
		ASIN:           "B08XYZ1234",
		Title:          "Coffee Pods - Subscribe & Save",
		Price:          24.99,
		Discount:       5.0,
		FrequencyWeeks: 4,
		NextDelivery:   time.Now().AddDate(0, 0, 14), // Current next delivery (2 weeks from now)
		Status:         "active",
		Quantity:       1,
	}

	// Advance NextDelivery by FrequencyWeeks
	subscription.NextDelivery = subscription.NextDelivery.AddDate(0, 0, subscription.FrequencyWeeks*7)

	return subscription, nil
}

// CancelSubscription cancels a subscription by setting its Status to "cancelled"
func (c *Client) CancelSubscription(id string) (*models.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("subscription ID cannot be empty")
	}

	// TODO: Implement actual Amazon API call to cancel subscription
	// For now, simulate the operation with mock data
	subscription := &models.Subscription{
		ID:             id,
		ASIN:           "B08XYZ1234",
		Title:          "Coffee Pods - Subscribe & Save",
		Price:          24.99,
		Discount:       5.0,
		FrequencyWeeks: 4,
		NextDelivery:   time.Now().AddDate(0, 0, 14),
		Status:         "active",
		Quantity:       1,
	}

	// Set status to cancelled
	subscription.Status = "cancelled"

	return subscription, nil
}

// UpdateFrequency updates the delivery frequency for a subscription
func (c *Client) UpdateFrequency(id string, intervalWeeks int) (*models.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("subscription ID cannot be empty")
	}

	if intervalWeeks < 1 || intervalWeeks > 26 {
		return nil, fmt.Errorf("interval must be between 1 and 26 weeks")
	}

	// TODO: Implement actual Amazon API call to update subscription frequency
	// For now, simulate the operation with mock data
	subscription := &models.Subscription{
		ID:             id,
		ASIN:           "B08XYZ1234",
		Title:          "Coffee Pods - Subscribe & Save",
		Price:          24.99,
		Discount:       5.0,
		FrequencyWeeks: 4,
		NextDelivery:   time.Now().AddDate(0, 0, 14),
		Status:         "active",
		Quantity:       1,
	}

	// Update the frequency
	subscription.FrequencyWeeks = intervalWeeks

	return subscription, nil
}
