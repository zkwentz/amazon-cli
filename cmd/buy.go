package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	buyQuantity  int
	buyConfirm   bool
	buyAddressID string
	buyPaymentID string
)

// buyCmd represents the buy command
var buyCmd = &cobra.Command{
	Use:   "buy <asin>",
	Short: "Quick purchase an item",
	Long: `Quickly purchase an item by ASIN without adding to cart first.
Requires --confirm flag to execute the purchase.
Without --confirm, shows a preview of what would be purchased.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		asin := args[0]
		c := getClient()

		// Get product details
		product, err := c.GetProduct(asin)
		if err != nil {
			output.Error(models.ErrNotFound, "Product not found: "+err.Error(), nil)
			os.Exit(models.ExitNotFound)
		}

		// Calculate total (price * quantity + estimated tax)
		subtotal := product.Price * float64(buyQuantity)
		tax := subtotal * 0.08
		total := subtotal + tax

		if !buyConfirm {
			// Preview purchase
			output.JSON(map[string]interface{}{
				"dry_run": true,
				"product": map[string]interface{}{
					"asin":  product.ASIN,
					"title": product.Title,
					"price": product.Price,
				},
				"quantity":      buyQuantity,
				"subtotal":      subtotal,
				"estimated_tax": tax,
				"total":         total,
				"message":       "Add --confirm to complete purchase",
			})
			return
		}

		// Get address and payment IDs, use defaults if not provided
		addressID := buyAddressID
		paymentID := buyPaymentID

		if addressID == "" {
			addresses, _ := c.GetAddresses()
			for _, addr := range addresses {
				if addr.Default {
					addressID = addr.ID
					break
				}
			}
			if addressID == "" && len(addresses) > 0 {
				addressID = addresses[0].ID
			}
		}

		if paymentID == "" {
			payments, _ := c.GetPaymentMethods()
			for _, pm := range payments {
				if pm.Default {
					paymentID = pm.ID
					break
				}
			}
			if paymentID == "" && len(payments) > 0 {
				paymentID = payments[0].ID
			}
		}

		// Add to cart and checkout
		_, err = c.AddToCart(asin, buyQuantity)
		if err != nil {
			output.Error(models.ErrInvalidInput, "Failed to add to cart: "+err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		confirmation, err := c.CompleteCheckout(addressID, paymentID)
		if err != nil {
			output.Error(models.ErrPurchaseFailed, "Checkout failed: "+err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		output.JSON(confirmation)
	},
}

func init() {
	rootCmd.AddCommand(buyCmd)

	buyCmd.Flags().IntVarP(&buyQuantity, "quantity", "n", 1, "Quantity to purchase")
	buyCmd.Flags().BoolVar(&buyConfirm, "confirm", false, "Confirm the purchase")
	buyCmd.Flags().StringVar(&buyAddressID, "address-id", "", "Shipping address ID")
	buyCmd.Flags().StringVar(&buyPaymentID, "payment-id", "", "Payment method ID")
}
