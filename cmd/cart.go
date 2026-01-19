package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	cartQuantity  int
	cartConfirm   bool
	cartAddressID string
	cartPaymentID string
)

// Shared client instance for cart operations
var client *amazon.Client

func getClient() *amazon.Client {
	if client == nil {
		client = amazon.NewClient()
	}
	return client
}

// cartCmd represents the cart command
var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage shopping cart",
	Long:  `Add, remove, view, and checkout items in your Amazon shopping cart.`,
}

// cartAddCmd represents the cart add command
var cartAddCmd = &cobra.Command{
	Use:   "add <asin>",
	Short: "Add item to cart",
	Long:  `Add an item to your shopping cart by ASIN (Amazon Standard Identification Number).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		asin := args[0]
		c := getClient()

		cart, err := c.AddToCart(asin, cartQuantity)
		if err != nil {
			_ = output.Error(models.ErrInvalidInput, err.Error(), nil)
			os.Exit(models.ExitInvalidArgs)
		}

		_ = output.JSON(cart)
	},
}

// cartListCmd represents the cart list command
var cartListCmd = &cobra.Command{
	Use:   "list",
	Short: "View cart contents",
	Long:  `Display all items currently in your shopping cart with prices and totals.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		cart, err := c.GetCart()
		if err != nil {
			_ = output.Error(models.ErrAmazonError, err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		_ = output.JSON(cart)
	},
}

// cartRemoveCmd represents the cart remove command
var cartRemoveCmd = &cobra.Command{
	Use:   "remove <asin>",
	Short: "Remove item from cart",
	Long:  `Remove an item from your shopping cart by ASIN.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		asin := args[0]
		c := getClient()

		cart, err := c.RemoveFromCart(asin)
		if err != nil {
			_ = output.Error(models.ErrInvalidInput, err.Error(), nil)
			os.Exit(models.ExitInvalidArgs)
		}

		_ = output.JSON(cart)
	},
}

// cartClearCmd represents the cart clear command
var cartClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all items from cart",
	Long:  `Remove all items from your shopping cart. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		if !cartConfirm {
			// Dry run - show what would be cleared
			cart, _ := c.GetCart()
			_ = output.JSON(map[string]interface{}{
				"dry_run":       true,
				"would_clear":   cart.ItemCount,
				"current_total": cart.Total,
				"message":       "Add --confirm to execute",
			})
			return
		}

		cart, _ := c.GetCart()
		itemCount := cart.ItemCount

		err := c.ClearCart()
		if err != nil {
			_ = output.Error(models.ErrAmazonError, err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		_ = output.JSON(map[string]interface{}{
			"status":        "cleared",
			"items_removed": itemCount,
		})
	},
}

// cartCheckoutCmd represents the cart checkout command
var cartCheckoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checkout cart",
	Long: `Complete purchase of items in cart.
Requires --confirm flag to execute the purchase.
Without --confirm, shows a preview of the order.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		// Check confirm flag BEFORE any checkout logic
		// Default is preview mode - if --confirm is not set, show preview
		if !cartConfirm {
			// Get address and payment IDs for preview
			addressID := cartAddressID
			paymentID := cartPaymentID

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

			// Preview checkout
			preview, err := c.PreviewCheckout(addressID, paymentID)
			if err != nil {
				_ = output.Error(models.ErrInvalidInput, err.Error(), nil)
				os.Exit(models.ExitInvalidArgs)
			}

			_ = output.JSON(map[string]interface{}{
				"dry_run":        true,
				"cart":           preview.Cart,
				"address":        preview.Address,
				"payment_method": preview.PaymentMethod,
				"message":        "Add --confirm to complete purchase",
			})
			return
		}

		// Execute checkout only when --confirm flag is set
		// Get address and payment IDs for actual checkout
		addressID := cartAddressID
		paymentID := cartPaymentID

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

		confirmation, err := c.CompleteCheckout(addressID, paymentID)
		if err != nil {
			_ = output.Error(models.ErrPurchaseFailed, err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		_ = output.JSON(confirmation)
	},
}

func init() {
	rootCmd.AddCommand(cartCmd)

	// Add subcommands
	cartCmd.AddCommand(cartAddCmd)
	cartCmd.AddCommand(cartListCmd)
	cartCmd.AddCommand(cartRemoveCmd)
	cartCmd.AddCommand(cartClearCmd)
	cartCmd.AddCommand(cartCheckoutCmd)

	// Flags for cart add
	cartAddCmd.Flags().IntVarP(&cartQuantity, "quantity", "n", 1, "Quantity to add")

	// Flags for cart clear
	cartClearCmd.Flags().BoolVar(&cartConfirm, "confirm", false, "Confirm the operation")

	// Flags for cart checkout
	cartCheckoutCmd.Flags().BoolVar(&cartConfirm, "confirm", false, "Confirm the purchase")
	cartCheckoutCmd.Flags().StringVar(&cartAddressID, "address-id", "", "Shipping address ID")
	cartCheckoutCmd.Flags().StringVar(&cartPaymentID, "payment-id", "", "Payment method ID")
}
