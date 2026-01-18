package models

import (
	"encoding/json"
	"testing"
)

func TestProductJSONMarshaling(t *testing.T) {
	originalPrice := 349.99
	product := Product{
		ASIN:             "B08N5WRWNW",
		Title:            "Sony WH-1000XM4 Wireless Headphones",
		Price:            278.00,
		OriginalPrice:    &originalPrice,
		Rating:           4.7,
		ReviewCount:      52431,
		Prime:            true,
		InStock:          true,
		DeliveryEstimate: "Tomorrow",
		Description:      "Industry-leading noise canceling headphones",
		Features:         []string{"Noise Canceling", "Bluetooth", "30hr Battery"},
		Images:           []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
	}

	// Test marshaling
	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal product: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Product
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal product: %v", err)
	}

	// Verify key fields
	if unmarshaled.ASIN != product.ASIN {
		t.Errorf("ASIN mismatch: got %s, want %s", unmarshaled.ASIN, product.ASIN)
	}
	if unmarshaled.Title != product.Title {
		t.Errorf("Title mismatch: got %s, want %s", unmarshaled.Title, product.Title)
	}
	if unmarshaled.Price != product.Price {
		t.Errorf("Price mismatch: got %f, want %f", unmarshaled.Price, product.Price)
	}
	if unmarshaled.OriginalPrice == nil || *unmarshaled.OriginalPrice != *product.OriginalPrice {
		t.Errorf("OriginalPrice mismatch")
	}
	if unmarshaled.Prime != product.Prime {
		t.Errorf("Prime mismatch: got %v, want %v", unmarshaled.Prime, product.Prime)
	}
}

func TestProductWithoutOptionalFields(t *testing.T) {
	product := Product{
		ASIN:             "B07EXAMPLE",
		Title:            "Basic Product",
		Price:            19.99,
		Rating:           4.0,
		ReviewCount:      100,
		Prime:            false,
		InStock:          true,
		DeliveryEstimate: "2-3 days",
	}

	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal product: %v", err)
	}

	// Verify optional fields are omitted from JSON
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if _, exists := jsonMap["original_price"]; exists {
		t.Error("original_price should be omitted when nil")
	}
	if _, exists := jsonMap["description"]; exists {
		t.Error("description should be omitted when empty")
	}
	if _, exists := jsonMap["features"]; exists {
		t.Error("features should be omitted when nil")
	}
	if _, exists := jsonMap["images"]; exists {
		t.Error("images should be omitted when nil")
	}
}

func TestSearchResponseJSONMarshaling(t *testing.T) {
	response := SearchResponse{
		Query: "wireless headphones",
		Results: []Product{
			{
				ASIN:             "B08N5WRWNW",
				Title:            "Sony WH-1000XM4",
				Price:            278.00,
				Rating:           4.7,
				ReviewCount:      52431,
				Prime:            true,
				InStock:          true,
				DeliveryEstimate: "Tomorrow",
			},
		},
		TotalResults: 1000,
		Page:         1,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal SearchResponse: %v", err)
	}

	var unmarshaled SearchResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SearchResponse: %v", err)
	}

	if unmarshaled.Query != response.Query {
		t.Errorf("Query mismatch: got %s, want %s", unmarshaled.Query, response.Query)
	}
	if unmarshaled.TotalResults != response.TotalResults {
		t.Errorf("TotalResults mismatch: got %d, want %d", unmarshaled.TotalResults, response.TotalResults)
	}
	if len(unmarshaled.Results) != len(response.Results) {
		t.Errorf("Results length mismatch: got %d, want %d", len(unmarshaled.Results), len(response.Results))
	}
}

func TestReviewJSONMarshaling(t *testing.T) {
	review := Review{
		Rating:   5.0,
		Title:    "Excellent headphones!",
		Body:     "These are the best headphones I've ever owned.",
		Author:   "John Doe",
		Date:     "2024-01-15",
		Verified: true,
	}

	data, err := json.Marshal(review)
	if err != nil {
		t.Fatalf("Failed to marshal Review: %v", err)
	}

	var unmarshaled Review
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Review: %v", err)
	}

	if unmarshaled.Rating != review.Rating {
		t.Errorf("Rating mismatch: got %f, want %f", unmarshaled.Rating, review.Rating)
	}
	if unmarshaled.Title != review.Title {
		t.Errorf("Title mismatch: got %s, want %s", unmarshaled.Title, review.Title)
	}
	if unmarshaled.Verified != review.Verified {
		t.Errorf("Verified mismatch: got %v, want %v", unmarshaled.Verified, review.Verified)
	}
}

func TestReviewsResponseJSONMarshaling(t *testing.T) {
	response := ReviewsResponse{
		ASIN: "B08N5WRWNW",
		Reviews: []Review{
			{
				Rating:   5.0,
				Title:    "Great!",
				Body:     "Love it",
				Author:   "User1",
				Date:     "2024-01-15",
				Verified: true,
			},
			{
				Rating:   4.0,
				Title:    "Good",
				Body:     "Nice product",
				Author:   "User2",
				Date:     "2024-01-14",
				Verified: false,
			},
		},
		AverageRating: 4.5,
		TotalReviews:  1000,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal ReviewsResponse: %v", err)
	}

	var unmarshaled ReviewsResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ReviewsResponse: %v", err)
	}

	if unmarshaled.ASIN != response.ASIN {
		t.Errorf("ASIN mismatch: got %s, want %s", unmarshaled.ASIN, response.ASIN)
	}
	if unmarshaled.AverageRating != response.AverageRating {
		t.Errorf("AverageRating mismatch: got %f, want %f", unmarshaled.AverageRating, response.AverageRating)
	}
	if unmarshaled.TotalReviews != response.TotalReviews {
		t.Errorf("TotalReviews mismatch: got %d, want %d", unmarshaled.TotalReviews, response.TotalReviews)
	}
	if len(unmarshaled.Reviews) != len(response.Reviews) {
		t.Errorf("Reviews length mismatch: got %d, want %d", len(unmarshaled.Reviews), len(response.Reviews))
	}
}

func TestEmptySearchResponse(t *testing.T) {
	response := SearchResponse{
		Query:        "nonexistent product xyz123",
		Results:      []Product{},
		TotalResults: 0,
		Page:         1,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal empty SearchResponse: %v", err)
	}

	var unmarshaled SearchResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty SearchResponse: %v", err)
	}

	if len(unmarshaled.Results) != 0 {
		t.Errorf("Expected empty results, got %d items", len(unmarshaled.Results))
	}
	if unmarshaled.TotalResults != 0 {
		t.Errorf("Expected 0 total results, got %d", unmarshaled.TotalResults)
	}
}
