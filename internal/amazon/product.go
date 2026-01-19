package amazon

import (
	"fmt"
	"time"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// GetProduct retrieves detailed product information
func (c *Client) GetProduct(asin string) (*models.Product, error) {
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}

	// TODO: Implement actual Amazon API call
	// For now, return mock data

	originalPrice := 349.99

	return &models.Product{
		ASIN:             asin,
		Title:            "Sony WH-1000XM4 Wireless Premium Noise Canceling Overhead Headphones",
		Price:            278.00,
		OriginalPrice:    &originalPrice,
		Rating:           4.7,
		ReviewCount:      52431,
		Prime:            true,
		InStock:          true,
		DeliveryEstimate: "Tomorrow",
		Description:      "Industry-leading noise canceling with Dual Noise Sensor technology. Next-level music with Edge-AI, co-developed with Sony Music Studios Tokyo. Up to 30-hour battery life with quick charging (10 min charge for 5 hours of playback).",
		Features: []string{
			"Industry-leading noise cancellation",
			"30-hour battery life",
			"Touch sensor controls",
			"Speak-to-chat technology",
			"Wearing detection",
			"Multipoint connection",
		},
		Images: []string{
			"https://images-na.ssl-images-amazon.com/images/I/71o8Q5XJS5L._AC_SL1500_.jpg",
			"https://images-na.ssl-images-amazon.com/images/I/81WpXBD4uWL._AC_SL1500_.jpg",
		},
	}, nil
}

// GetProductReviews retrieves reviews for a product
func (c *Client) GetProductReviews(asin string, limit int) (*models.ReviewsResponse, error) {
	if asin == "" {
		return nil, fmt.Errorf("ASIN cannot be empty")
	}

	if limit <= 0 {
		limit = 10
	}

	// TODO: Implement actual Amazon API call
	// For now, return mock data

	reviews := []models.Review{
		{
			Rating:   5,
			Title:    "Best headphones I've ever owned",
			Body:     "The noise canceling is incredible. I use these daily for work calls and music. Battery life is exactly as advertised.",
			Author:   "John D.",
			Date:     time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			Verified: true,
		},
		{
			Rating:   4,
			Title:    "Great but pricey",
			Body:     "Sound quality is excellent and the ANC is top-notch. Only complaint is the price, but you get what you pay for.",
			Author:   "Sarah M.",
			Date:     time.Now().AddDate(0, 0, -25).Format("2006-01-02"),
			Verified: true,
		},
		{
			Rating:   5,
			Title:    "Perfect for travel",
			Body:     "Used these on a 12-hour flight and they were amazing. The noise canceling blocked out all the engine noise.",
			Author:   "Mike R.",
			Date:     time.Now().AddDate(0, -1, -5).Format("2006-01-02"),
			Verified: true,
		},
	}

	if len(reviews) > limit {
		reviews = reviews[:limit]
	}

	return &models.ReviewsResponse{
		ASIN:          asin,
		AverageRating: 4.7,
		TotalReviews:  52431,
		Reviews:       reviews,
	}, nil
}
