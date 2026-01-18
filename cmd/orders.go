package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/internal/output"
)

var (
	ordersLimit  int
	ordersStatus string
)

// ordersCmd represents the orders command
var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage Amazon orders",
	Long:  `Commands for viewing and managing Amazon orders including list, get, track, and history.`,
}

// ordersListCmd represents the orders list command
var ordersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent orders",
	Long: `List recent Amazon orders with optional filtering by status and limit.

Examples:
  amazon-cli orders list
  amazon-cli orders list --limit 5
  amazon-cli orders list --status delivered
  amazon-cli orders list --limit 10 --status pending`,
	RunE: runOrdersList,
}

func init() {
	rootCmd.AddCommand(ordersCmd)
	ordersCmd.AddCommand(ordersListCmd)

	// Add flags for orders list command
	ordersListCmd.Flags().IntVar(&ordersLimit, "limit", 10, "Maximum number of orders to return")
	ordersListCmd.Flags().StringVar(&ordersStatus, "status", "", "Filter by status: pending, delivered, returned")
}

func runOrdersList(cmd *cobra.Command, args []string) error {
	// Create output printer
	printer := output.NewPrinter(outputFormat, quiet)

	// Create Amazon client
	client := amazon.NewClient()

	// Get orders
	ordersResponse, err := client.GetOrders(ordersLimit, ordersStatus)
	if err != nil {
		printer.PrintError(err)
		os.Exit(1)
	}

	// Print response
	return printer.Print(ordersResponse)
}
