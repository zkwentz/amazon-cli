package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/michaelshimeles/amazon-cli/pkg/models"
	"github.com/spf13/cobra"
)

var (
	buyConfirm   bool
	buyQuantity  int
	buyAddressID string
	buyPaymentID string
)

// buyCmd represents the buy command
var buyCmd = &cobra.Command{
	Use:   "buy <asin>",
	Short: "Quick purchase - buy a product immediately",
	Long: `Quick purchase a product by ASIN.

This command combines adding to cart and checkout into a single operation.
REQUIRES --confirm flag to actually complete the purchase.

Without --confirm, shows what would be purchased (dry run).
With --confirm, adds the item to cart and completes checkout.

Examples:
  # Preview purchase (dry run)
  amazon-cli buy B08N5WRWNW --quantity 2

  # Complete purchase
  amazon-cli buy B08N5WRWNW --quantity 2 --confirm --address-id addr_123 --payment-id pay_456

  # Use default address and payment
  amazon-cli buy B08N5WRWNW --confirm`,
	Args: cobra.ExactArgs(1),
	RunE: runBuy,
}

func init() {
	rootCmd.AddCommand(buyCmd)

	buyCmd.Flags().BoolVar(&buyConfirm, "confirm", false, "Confirm the purchase (REQUIRED to execute)")
	buyCmd.Flags().IntVar(&buyQuantity, "quantity", 1, "Quantity to purchase")
	buyCmd.Flags().StringVar(&buyAddressID, "address-id", "", "Shipping address ID (uses default if not specified)")
	buyCmd.Flags().StringVar(&buyPaymentID, "payment-id", "", "Payment method ID (uses default if not specified)")
}

func runBuy(cmd *cobra.Command, args []string) error {
	asin := args[0]

	// Validate ASIN format (basic validation)
	if len(asin) != 10 {
		return fmt.Errorf("invalid ASIN format: must be 10 characters")
	}

	// Validate quantity
	if buyQuantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	// Create Amazon client
	client := amazon.NewClient()

	// Without --confirm, show preview (dry run)
	if !buyConfirm {
		product, err := client.GetProduct(asin)
		if err != nil {
			return fmt.Errorf("failed to get product details: %w", err)
		}

		total := product.Price * float64(buyQuantity)
		preview := models.BuyPreview{
			DryRun:   true,
			Product:  product,
			Quantity: buyQuantity,
			Total:    total,
			Message:  "Add --confirm to execute this purchase",
		}

		output, err := json.MarshalIndent(preview, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	}

	// With --confirm, complete the purchase
	// Use default address/payment if not specified
	addressID := buyAddressID
	if addressID == "" {
		addressID = "default"
	}

	paymentID := buyPaymentID
	if paymentID == "" {
		paymentID = "default"
	}

	// Perform quick buy
	confirmation, err := client.QuickBuy(asin, buyQuantity, addressID, paymentID)
	if err != nil {
		return fmt.Errorf("purchase failed: %w", err)
	}

	// Output confirmation
	output, err := json.MarshalIndent(confirmation, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Println(string(output))
	return nil
}
