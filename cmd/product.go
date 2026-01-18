package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
)

var (
	reviewLimit int
)

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Product-related commands",
	Long:  "Commands for retrieving product information and reviews from Amazon.",
}

var reviewsCmd = &cobra.Command{
	Use:   "reviews <asin>",
	Short: "Get product reviews",
	Long: `Fetch and display product reviews for a given ASIN (Amazon Standard Identification Number).

The ASIN is a 10-character alphanumeric identifier found on Amazon product pages.

Example:
  amazon-cli product reviews B08N5WRWNW --limit 5`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]

		// Validate ASIN format (10 alphanumeric characters)
		if !isValidASIN(asin) {
			return fmt.Errorf("invalid ASIN format: %s (must be 10 alphanumeric characters)", asin)
		}

		// Create Amazon client
		client := amazon.NewClient()

		// Fetch reviews
		reviews, err := client.GetProductReviews(asin, reviewLimit)
		if err != nil {
			return fmt.Errorf("failed to fetch reviews: %w", err)
		}

		// Output as JSON
		output, err := json.MarshalIndent(reviews, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	// Add flags
	reviewsCmd.Flags().IntVarP(&reviewLimit, "limit", "l", 10, "Maximum number of reviews to fetch")

	// Add subcommands
	productCmd.AddCommand(reviewsCmd)
	rootCmd.AddCommand(productCmd)
}

// isValidASIN validates that an ASIN is 10 alphanumeric characters
func isValidASIN(asin string) bool {
	match, _ := regexp.MatchString(`^[A-Z0-9]{10}$`, asin)
	return match
}

// printJSON outputs data as formatted JSON
func printJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// printError outputs an error as JSON
func printError(code, message string) {
	errorOutput := map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}
	json.NewEncoder(os.Stderr).Encode(errorOutput)
}
