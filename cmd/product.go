package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	reviewsLimit int
)

// productCmd represents the product command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Get product information",
	Long:  `Retrieve detailed product information and reviews.`,
}

// productGetCmd represents the product get command
var productGetCmd = &cobra.Command{
	Use:   "get <asin>",
	Short: "Get product details",
	Long:  `Display detailed information about a product by its ASIN.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		asin := args[0]

		// Validate ASIN argument
		if asin == "" {
			output.Error(models.ErrInvalidInput, "ASIN cannot be empty", nil)
			os.Exit(models.ExitInvalidArgs)
		}

		c := getClient()

		product, err := c.GetProduct(asin)
		if err != nil {
			// Check if error is validation related
			errMsg := err.Error()
			if strings.Contains(errMsg, "invalid ASIN format") || strings.Contains(errMsg, "ASIN cannot be empty") {
				output.Error(models.ErrInvalidInput, errMsg, nil)
				os.Exit(models.ExitInvalidArgs)
			}
			// Otherwise treat as general error
			output.Error(models.ErrAmazonError, errMsg, nil)
			os.Exit(models.ExitGeneralError)
		}

		output.JSON(product)
	},
}

// productReviewsCmd represents the product reviews command
var productReviewsCmd = &cobra.Command{
	Use:   "reviews <asin>",
	Short: "Get product reviews",
	Long:  `Display reviews for a product by its ASIN.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		asin := args[0]

		// Validate ASIN argument
		if asin == "" {
			output.Error(models.ErrInvalidInput, "ASIN cannot be empty", nil)
			os.Exit(models.ExitInvalidArgs)
		}

		c := getClient()

		reviews, err := c.GetProductReviews(asin, reviewsLimit)
		if err != nil {
			// Check if error is validation related
			errMsg := err.Error()
			if strings.Contains(errMsg, "ASIN cannot be empty") {
				output.Error(models.ErrInvalidInput, errMsg, nil)
				os.Exit(models.ExitInvalidArgs)
			}
			// Otherwise treat as general error
			output.Error(models.ErrAmazonError, errMsg, nil)
			os.Exit(models.ExitGeneralError)
		}

		output.JSON(reviews)
	},
}

func init() {
	rootCmd.AddCommand(productCmd)

	productCmd.AddCommand(productGetCmd)
	productCmd.AddCommand(productReviewsCmd)

	productReviewsCmd.Flags().IntVar(&reviewsLimit, "limit", 10, "Number of reviews to return")
}
