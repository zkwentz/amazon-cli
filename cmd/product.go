package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
)

// productCmd represents the product command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Product information commands",
	Long:  `Commands for retrieving product information from Amazon.`,
}

// productGetCmd represents the product get command
var productGetCmd = &cobra.Command{
	Use:   "get <asin>",
	Short: "Get product details by ASIN",
	Long: `Retrieve detailed information about a product using its ASIN (Amazon Standard Identification Number).

ASIN is a 10-character alphanumeric unique identifier assigned by Amazon.
You can find it in the product URL or in the product details section.

Example:
  amazon-cli product get B08N5WRWNW`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]

		// Create Amazon client
		client, err := amazon.NewClient(rootConfig)
		if err != nil {
			printer.PrintError(err)
			return err
		}

		// Get product
		product, err := client.GetProduct(asin)
		if err != nil {
			printer.PrintError(err)
			return err
		}

		// Print result
		return printer.Print(product)
	},
}

func init() {
	rootCmd.AddCommand(productCmd)
	productCmd.AddCommand(productGetCmd)
}
