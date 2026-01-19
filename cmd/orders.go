package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	ordersLimit  int
	ordersStatus string
	ordersYear   int
)

// ordersCmd represents the orders command
var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage orders",
	Long:  `List orders, get order details, and track shipments.`,
}

// ordersListCmd represents the orders list command
var ordersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent orders",
	Long:  `Display a list of your recent Amazon orders with status and tracking info.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		orders, err := c.GetOrders(ordersLimit, ordersStatus)
		if err != nil {
			output.Error(models.ErrAmazonError, err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		output.JSON(orders)
	},
}

// ordersGetCmd represents the orders get command
var ordersGetCmd = &cobra.Command{
	Use:   "get <order-id>",
	Short: "Get order details",
	Long:  `Display detailed information about a specific order.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		orderID := args[0]
		c := getClient()

		order, err := c.GetOrder(orderID)
		if err != nil {
			output.Error(models.ErrNotFound, err.Error(), nil)
			os.Exit(models.ExitNotFound)
		}

		output.JSON(order)
	},
}

// ordersTrackCmd represents the orders track command
var ordersTrackCmd = &cobra.Command{
	Use:   "track <order-id>",
	Short: "Track order shipment",
	Long:  `Display tracking information for an order's shipment.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		orderID := args[0]
		c := getClient()

		tracking, err := c.GetOrderTracking(orderID)
		if err != nil {
			output.Error(models.ErrNotFound, err.Error(), nil)
			os.Exit(models.ExitNotFound)
		}

		output.JSON(tracking)
	},
}

// ordersHistoryCmd represents the orders history command
var ordersHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Get order history",
	Long:  `Display order history for a specific year.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		year := ordersYear
		if year == 0 {
			year = time.Now().Year()
		}

		orders, err := c.GetOrderHistory(year)
		if err != nil {
			output.Error(models.ErrAmazonError, err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		output.JSON(orders)
	},
}

func init() {
	rootCmd.AddCommand(ordersCmd)

	// Add subcommands
	ordersCmd.AddCommand(ordersListCmd)
	ordersCmd.AddCommand(ordersGetCmd)
	ordersCmd.AddCommand(ordersTrackCmd)
	ordersCmd.AddCommand(ordersHistoryCmd)

	// Flags for orders list
	ordersListCmd.Flags().IntVar(&ordersLimit, "limit", 10, "Number of orders to return")
	ordersListCmd.Flags().StringVar(&ordersStatus, "status", "", "Filter by status: pending, delivered, returned")

	// Flags for orders history
	ordersHistoryCmd.Flags().IntVar(&ordersYear, "year", 0, "Year to fetch orders from (default: current year)")
}
