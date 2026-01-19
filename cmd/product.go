package cmd

import (
	"os"

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
		c := getClient()

		product, err := c.GetProduct(asin)
		if err != nil {
			output.Error(models.ErrNotFound, err.Error(), nil)
			os.Exit(models.ExitNotFound)
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
		c := getClient()

		reviews, err := c.GetProductReviews(asin, reviewsLimit)
		if err != nil {
			output.Error(models.ErrNotFound, err.Error(), nil)
			os.Exit(models.ExitNotFound)
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
