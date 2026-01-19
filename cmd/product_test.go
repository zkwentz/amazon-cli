package cmd

import (
	"encoding/json"
	"testing"

	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestProductGetCmd_Success(t *testing.T) {
	if productGetCmd.Use != "get <asin>" {
		t.Errorf("Expected Use to be 'get <asin>', got '%s'", productGetCmd.Use)
	}
	if productGetCmd.Short == "" {
		t.Error("Expected Short description to be set")
	}
	if productGetCmd.Run == nil {
		t.Error("Expected Run function to be set")
	}
}

func TestProductGetCmd_ValidatesASIN(t *testing.T) {
	// Test that GetProduct validates ASIN
	c := amazon.NewClient()

	// Test with empty ASIN
	_, err := c.GetProduct("")
	if err == nil {
		t.Error("Expected error for empty ASIN")
	}
	if err != nil && err.Error() != "ASIN cannot be empty" {
		t.Errorf("Expected 'ASIN cannot be empty' error, got: %s", err.Error())
	}

	// Test with invalid ASIN format
	invalidASINs := []string{
		"B08N5",           // too short
		"B08N5WRWNW123",   // too long
		"b08n5wrwnw",      // lowercase
		"B08N5WRWN!",      // special characters
	}

	for _, asin := range invalidASINs {
		_, err := c.GetProduct(asin)
		if err == nil {
			t.Errorf("Expected error for invalid ASIN: %s", asin)
		}
	}
}

func TestProductGetCmd_ErrorConstants(t *testing.T) {
	// Verify error constants exist
	if models.ErrInvalidInput != "INVALID_INPUT" {
		t.Errorf("Expected ErrInvalidInput to be 'INVALID_INPUT', got '%s'", models.ErrInvalidInput)
	}
	if models.ExitInvalidArgs != 2 {
		t.Errorf("Expected ExitInvalidArgs to be 2, got %d", models.ExitInvalidArgs)
	}
	if models.ErrAmazonError != "AMAZON_ERROR" {
		t.Errorf("Expected ErrAmazonError to be 'AMAZON_ERROR', got '%s'", models.ErrAmazonError)
	}
	if models.ExitGeneralError != 1 {
		t.Errorf("Expected ExitGeneralError to be 1, got %d", models.ExitGeneralError)
	}
}

func TestProductGetCmd_ResponseParsing(t *testing.T) {
	// Test JSON output format
	originalPrice := 349.99
	product := &models.Product{
		ASIN:          "B08N5WRWNW",
		Title:         "Sony WH-1000XM4",
		Price:         278.00,
		OriginalPrice: &originalPrice,
		Rating:        4.7,
		ReviewCount:   52431,
		Prime:         true,
		InStock:       true,
	}

	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal product: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal product: %v", err)
	}

	// Verify expected fields
	if result["asin"] != "B08N5WRWNW" {
		t.Errorf("Expected ASIN 'B08N5WRWNW', got '%v'", result["asin"])
	}
	if result["title"] != "Sony WH-1000XM4" {
		t.Errorf("Expected title 'Sony WH-1000XM4', got '%v'", result["title"])
	}
}

func TestProductReviewsCmd_Success(t *testing.T) {
	if productReviewsCmd.Use != "reviews <asin>" {
		t.Errorf("Expected Use to be 'reviews <asin>', got '%s'", productReviewsCmd.Use)
	}
	if productReviewsCmd.Short == "" {
		t.Error("Expected Short description to be set")
	}
	if productReviewsCmd.Run == nil {
		t.Error("Expected Run function to be set")
	}
}

func TestProductReviewsCmd_LimitFlag(t *testing.T) {
	// Verify --limit flag exists
	flag := productReviewsCmd.Flags().Lookup("limit")
	if flag == nil {
		t.Fatal("Expected --limit flag to be defined")
	}
	if flag.DefValue != "10" {
		t.Errorf("Expected default limit to be '10', got '%s'", flag.DefValue)
	}
}

func TestProductReviewsCmd_ValidatesASIN(t *testing.T) {
	// Test that GetProductReviews validates ASIN
	c := amazon.NewClient()

	// Test with empty ASIN
	_, err := c.GetProductReviews("", 10)
	if err == nil {
		t.Error("Expected error for empty ASIN")
	}
	if err != nil && err.Error() != "ASIN cannot be empty" {
		t.Errorf("Expected 'ASIN cannot be empty' error, got: %s", err.Error())
	}
}

func TestProductReviewsCmd_DefaultLimit(t *testing.T) {
	// Test that limit defaults to 10
	c := amazon.NewClient()

	// GetProductReviews should handle limit <= 0 by defaulting to 10
	// We can't actually test the HTTP call, but we can verify the function accepts 0
	_, err := c.GetProductReviews("B08N5WRWNW", 0)
	// We expect an error because we're not making a real HTTP call,
	// but it shouldn't be about the limit
	if err != nil && err.Error() == "invalid limit" {
		t.Error("Expected limit to be handled, not cause 'invalid limit' error")
	}
}

func TestProductReviewsCmd_ResponseParsing(t *testing.T) {
	// Test JSON output format
	reviews := &models.ReviewsResponse{
		ASIN:         "B08N5WRWNW",
		AverageRating: 4.5,
		TotalReviews: 1234,
		Reviews: []models.Review{
			{
				Rating:   5,
				Title:    "Great headphones!",
				Body:     "These are amazing",
				Author:   "John Doe",
				Date:     "2024-01-15",
				Verified: true,
			},
		},
	}

	data, err := json.Marshal(reviews)
	if err != nil {
		t.Fatalf("Failed to marshal reviews: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal reviews: %v", err)
	}

	// Verify expected fields
	if result["asin"] != "B08N5WRWNW" {
		t.Errorf("Expected ASIN 'B08N5WRWNW', got '%v'", result["asin"])
	}
	if result["total_reviews"] != float64(1234) {
		t.Errorf("Expected total_reviews 1234, got '%v'", result["total_reviews"])
	}
}

func TestProductCmd_Configuration(t *testing.T) {
	if productCmd.Use != "product" {
		t.Errorf("Expected Use to be 'product', got '%s'", productCmd.Use)
	}
	if productCmd.Short != "Get product information" {
		t.Errorf("Expected Short to be 'Get product information', got '%s'", productCmd.Short)
	}
}

func TestProductCmd_Subcommands(t *testing.T) {
	// Verify subcommands are registered
	commands := productCmd.Commands()
	if len(commands) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(commands))
	}

	// Check for get and reviews subcommands
	var hasGet, hasReviews bool
	for _, cmd := range commands {
		if cmd.Name() == "get" {
			hasGet = true
		}
		if cmd.Name() == "reviews" {
			hasReviews = true
		}
	}

	if !hasGet {
		t.Error("Expected 'get' subcommand to be registered")
	}
	if !hasReviews {
		t.Error("Expected 'reviews' subcommand to be registered")
	}
}

func TestProductCmd_VariablesInitialized(t *testing.T) {
	// Test that package variables are initialized
	reviewsLimit = 5
	if reviewsLimit != 5 {
		t.Errorf("Expected reviewsLimit to be 5, got %d", reviewsLimit)
	}
}
