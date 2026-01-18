package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
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
	Short: "Quick buy - add to cart and checkout in one step",
	Long: `Quick purchase a product by ASIN. This command combines adding to cart
and checkout into a single operation.

SAFETY: This command REQUIRES --confirm flag to execute. Without it, the
command will display what would be purchased but not execute the purchase.

Examples:
  # Preview purchase (dry run)
  amazon-cli buy B08N5WRWNW

  # Actually purchase with confirmation
  amazon-cli buy B08N5WRWNW --confirm

  # Purchase specific quantity
  amazon-cli buy B08N5WRWNW --quantity 2 --confirm

  # Purchase with specific address and payment method
  amazon-cli buy B08N5WRWNW --confirm --address-id addr_123 --payment-id pay_456`,
	Args: cobra.ExactArgs(1),
	RunE: runBuy,
}

func init() {
	rootCmd.AddCommand(buyCmd)

	buyCmd.Flags().BoolVar(&buyConfirm, "confirm", false, "Confirm the purchase (required to execute)")
	buyCmd.Flags().IntVar(&buyQuantity, "quantity", 1, "Quantity to purchase")
	buyCmd.Flags().StringVar(&buyAddressID, "address-id", "", "Address ID to use for shipping")
	buyCmd.Flags().StringVar(&buyPaymentID, "payment-id", "", "Payment method ID to use")
}

func runBuy(cmd *cobra.Command, args []string) error {
	asin := args[0]

	// Validate ASIN format (should be 10 alphanumeric characters)
	if len(asin) != 10 {
		return fmt.Errorf("invalid ASIN format: %s (should be 10 characters)", asin)
	}

	// Validate quantity
	if buyQuantity <= 0 {
		return fmt.Errorf("quantity must be positive, got: %d", buyQuantity)
	}

	// Create Amazon client
	client := amazon.NewClient()

	// If not confirmed, show dry run
	if !buyConfirm {
		dryRun := map[string]interface{}{
			"dry_run": true,
			"product": map[string]interface{}{
				"asin":  asin,
				"title": "Product Preview", // In real implementation, fetch actual product details
				"price": 29.99,
			},
			"quantity":    buyQuantity,
			"total":       29.99 * float64(buyQuantity),
			"message":     "Add --confirm to execute this purchase",
			"address_id":  buyAddressID,
			"payment_id":  buyPaymentID,
		}

		output, err := json.MarshalIndent(dryRun, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Fprintln(os.Stdout, string(output))
		return nil
	}

	// Confirmed purchase - execute the buy flow

	// Step 1: Add to cart
	if !quiet {
		fmt.Fprintln(os.Stderr, "Adding product to cart...")
	}

	cart, err := client.AddToCart(asin, buyQuantity)
	if err != nil {
		return fmt.Errorf("failed to add to cart: %w", err)
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "Added %d item(s) to cart. Cart total: $%.2f\n", buyQuantity, cart.Total)
		fmt.Fprintln(os.Stderr, "Proceeding to checkout...")
	}

	// Step 2: Complete checkout
	// Use default address/payment if not specified
	addressID := buyAddressID
	if addressID == "" {
		addressID = "addr_default"
	}

	paymentID := buyPaymentID
	if paymentID == "" {
		paymentID = "pay_default"
	}

	confirmation, err := client.CompleteCheckout(addressID, paymentID)
	if err != nil {
		return fmt.Errorf("failed to complete checkout: %w", err)
	}

	// Output the order confirmation as JSON
	output, err := json.MarshalIndent(confirmation, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	if !quiet {
		fmt.Fprintln(os.Stderr, "\nOrder placed successfully!")
	}

	fmt.Fprintln(os.Stdout, string(output))
	return nil
}
