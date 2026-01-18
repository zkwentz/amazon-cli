package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	historyYear int
)

// ordersCmd represents the orders command
var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage Amazon orders",
	Long: `Commands for viewing and managing your Amazon orders.

You can list recent orders, get order details, track shipments, and view order history.`,
}

// ordersHistoryCmd represents the orders history command
var ordersHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Get order history for a specific year",
	Long: `Fetches all orders from a specific year.

By default, fetches orders from the current year. Use --year flag to specify a different year.

Example:
  amazon-cli orders history
  amazon-cli orders history --year 2023`,
	RunE: runOrdersHistory,
}

func init() {
	rootCmd.AddCommand(ordersCmd)
	ordersCmd.AddCommand(ordersHistoryCmd)

	// Flags for orders history command
	currentYear := time.Now().Year()
	ordersHistoryCmd.Flags().IntVar(&historyYear, "year", currentYear, "Year to fetch order history for")
}

func runOrdersHistory(cmd *cobra.Command, args []string) error {
	LogVerbose("Fetching order history for year %d", historyYear)

	// Create Amazon client
	client, err := amazon.NewClient()
	if err != nil {
		return models.NewCLIError(
			models.ErrorCodeNetworkError,
			"Failed to create Amazon client: "+err.Error(),
			nil,
		)
	}

	// Fetch order history
	response, err := client.GetOrderHistory(historyYear)
	if err != nil {
		return err
	}

	LogVerbose("Found %d orders for year %d", response.TotalCount, historyYear)

	// Print the response as JSON
	return PrintJSON(response)
}
