package models

// Subscription represents an Amazon Subscribe & Save subscription
type Subscription struct {
	SubscriptionID  string  `json:"subscription_id"`
	ASIN            string  `json:"asin"`
	Title           string  `json:"title"`
	Price           float64 `json:"price"`
	DiscountPercent int     `json:"discount_percent"`
	FrequencyWeeks  int     `json:"frequency_weeks"`
	NextDelivery    string  `json:"next_delivery"`
	Status          string  `json:"status"`
	Quantity        int     `json:"quantity"`
}

// SubscriptionsResponse represents the response containing all subscriptions
type SubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

// UpcomingDelivery represents a scheduled delivery from a subscription
type UpcomingDelivery struct {
	SubscriptionID string `json:"subscription_id"`
	ASIN           string `json:"asin"`
	Title          string `json:"title"`
	DeliveryDate   string `json:"delivery_date"`
	Quantity       int    `json:"quantity"`
}
