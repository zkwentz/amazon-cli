package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Product management commands",
	Long:  `Manage Amazon products: get product details, read reviews, and more.`,
}

var productGetCmd = &cobra.Command{
	Use:   "get <asin>",
	Short: "Get product details by ASIN",
	Long: `Get detailed information about a product by its ASIN (Amazon Standard Identification Number).

Returns product details including title, price, ratings, reviews, Prime eligibility, stock status,
delivery estimates, description, features, and images.

Example:
  amazon-cli product get B08N5WRWNW`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]

		// Validate ASIN format (10 alphanumeric characters)
		if err := validateASIN(asin); err != nil {
			return err
		}

		// TODO: Implement actual API call
		// For now, return a mock response
		product := map[string]interface{}{
			"asin":              asin,
			"title":             "Sample Product",
			"price":             29.99,
			"original_price":    39.99,
			"rating":            4.5,
			"review_count":      1234,
			"prime":             true,
			"in_stock":          true,
			"delivery_estimate": "Tomorrow",
			"description":       "This is a sample product description.",
			"features": []string{
				"Feature 1",
				"Feature 2",
				"Feature 3",
			},
			"images": []string{
				"https://example.com/image1.jpg",
				"https://example.com/image2.jpg",
			},
		}

		output, err := json.MarshalIndent(product, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var productReviewsCmd = &cobra.Command{
	Use:   "reviews <asin>",
	Short: "Get product reviews by ASIN",
	Long: `Get customer reviews for a product by its ASIN.

Returns a list of reviews including rating, title, body, author, date, and verified purchase status.

Example:
  amazon-cli product reviews B08N5WRWNW --limit 5`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		// Validate ASIN format
		if err := validateASIN(asin); err != nil {
			return err
		}

		// Validate limit
		if limit <= 0 {
			return fmt.Errorf("limit must be a positive integer")
		}

		// TODO: Implement actual API call
		// For now, return a mock response
		reviews := map[string]interface{}{
			"asin":           asin,
			"average_rating": 4.5,
			"total_reviews":  1234,
			"reviews": []map[string]interface{}{
				{
					"rating":   5,
					"title":    "Great product!",
					"body":     "This product exceeded my expectations. Highly recommended!",
					"author":   "John Doe",
					"date":     "2024-01-15",
					"verified": true,
				},
				{
					"rating":   4,
					"title":    "Good value",
					"body":     "Works well for the price. Minor issues but overall satisfied.",
					"author":   "Jane Smith",
					"date":     "2024-01-14",
					"verified": true,
				},
			},
		}

		// Limit reviews to requested count
		if reviewsList, ok := reviews["reviews"].([]map[string]interface{}); ok {
			if len(reviewsList) > limit {
				reviews["reviews"] = reviewsList[:limit]
			}
		}

		output, err := json.MarshalIndent(reviews, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// validateASIN validates that the ASIN is in the correct format
// ASINs are 10 alphanumeric characters
func validateASIN(asin string) error {
	if len(asin) != 10 {
		return fmt.Errorf("invalid ASIN: must be exactly 10 characters (got %d)", len(asin))
	}

	matched, err := regexp.MatchString("^[A-Z0-9]{10}$", asin)
	if err != nil {
		return fmt.Errorf("error validating ASIN: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid ASIN format: must contain only uppercase letters and numbers")
	}

	return nil
}

func init() {
	// Add flags
	productReviewsCmd.Flags().IntP("limit", "l", 10, "Maximum number of reviews to retrieve")

	// Add subcommands
	productCmd.AddCommand(productGetCmd)
	productCmd.AddCommand(productReviewsCmd)

	// Register with root command (assumes rootCmd exists)
	// Note: This will be properly integrated when root.go is implemented
	if rootCmd != nil {
		rootCmd.AddCommand(productCmd)
	}
}

// rootCmd reference - this would typically be imported from root.go
// For now, we check if it exists to prevent initialization errors
var rootCmd *cobra.Command

// SetRootCmd allows the root command to be set from outside this package
func SetRootCmd(cmd *cobra.Command) {
	rootCmd = cmd
	if rootCmd != nil && productCmd != nil {
		rootCmd.AddCommand(productCmd)
	}
}

// GetProductCmd returns the product command for testing purposes
func GetProductCmd() *cobra.Command {
	return productCmd
}
