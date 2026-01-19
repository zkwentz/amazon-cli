package models

import "time"

// Subscription represents an Amazon Subscribe & Save subscription
type Subscription struct {
	ID             string    `json:"id"`
	ASIN           string    `json:"asin"`
	Title          string    `json:"title"`
	Price          float64   `json:"price"`
	Discount       float64   `json:"discount"`
	FrequencyWeeks int       `json:"frequency_weeks"`
	NextDelivery   time.Time `json:"next_delivery"`
	Status         string    `json:"status"` // "active", "cancelled", "paused"
	Quantity       int       `json:"quantity"`
}

// SubscriptionList represents a list of subscriptions
type SubscriptionList struct {
	Subscriptions []Subscription `json:"subscriptions"`
	TotalCount    int            `json:"total_count"`
}
