package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	cartConfirm  bool
	cartQuantity int
	addressID    string
	paymentID    string
)

// cartCmd represents the cart command
var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage your Amazon shopping cart",
	Long: `Manage your Amazon shopping cart including adding items, viewing cart contents,
removing items, clearing the cart, and checking out.`,
}

// cartCheckoutCmd represents the cart checkout command
var cartCheckoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checkout your cart",
	Long: `Complete the checkout process for items in your cart.

IMPORTANT: This command REQUIRES the --confirm flag to execute the actual purchase.
Without --confirm, it will show a preview of what would be purchased (dry run).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()

		// If no address or payment specified, use defaults
		if addressID == "" {
			addressID = "addr_default"
		}
		if paymentID == "" {
			paymentID = "pay_default"
		}

		// Without --confirm, show preview (dry run)
		if !cartConfirm {
			preview, err := client.PreviewCheckout(addressID, paymentID)
			if err != nil {
				return fmt.Errorf("failed to preview checkout: %w", err)
			}

			// Output dry run response
			response := map[string]interface{}{
				"dry_run": true,
				"preview": preview,
				"message": "Add --confirm to execute the purchase",
			}

			output, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}

			fmt.Fprintln(os.Stdout, string(output))
			return nil
		}

		// With --confirm, complete the checkout
		confirmation, err := client.CompleteCheckout(addressID, paymentID)
		if err != nil {
			return fmt.Errorf("failed to complete checkout: %w", err)
		}

		// Output order confirmation
		output, err := json.MarshalIndent(confirmation, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Fprintln(os.Stdout, string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cartCmd)
	cartCmd.AddCommand(cartCheckoutCmd)

	// Flags for checkout command
	cartCheckoutCmd.Flags().BoolVar(&cartConfirm, "confirm", false, "Confirm the purchase (required to execute)")
	cartCheckoutCmd.Flags().StringVar(&addressID, "address-id", "", "Address ID to use for shipping")
	cartCheckoutCmd.Flags().StringVar(&paymentID, "payment-id", "", "Payment method ID to use")
}
