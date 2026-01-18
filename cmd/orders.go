package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// ordersCmd represents the orders command
var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage Amazon orders",
	Long:  `View and track your Amazon orders, including order history and tracking information.`,
}

// trackCmd represents the track subcommand
var trackCmd = &cobra.Command{
	Use:   "track <order-id>",
	Short: "Track shipment for an order",
	Long:  `Retrieve tracking information for a specific Amazon order, including carrier, tracking number, status, and delivery date.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runTrackCommand,
}

func init() {
	rootCmd.AddCommand(ordersCmd)
	ordersCmd.AddCommand(trackCmd)
}

// runTrackCommand executes the track command
func runTrackCommand(cmd *cobra.Command, args []string) error {
	orderID := args[0]
	printer := getPrinter()

	// Load configuration
	cfg, err := getConfig()
	if err != nil {
		return printer.PrintError(models.NewCLIError(
			models.ErrAmazonError,
			fmt.Sprintf("failed to load configuration: %v", err),
			nil,
		))
	}

	// Create Amazon client
	client, err := amazon.NewClient(cfg)
	if err != nil {
		return printer.PrintError(models.NewCLIError(
			models.ErrAmazonError,
			fmt.Sprintf("failed to create client: %v", err),
			nil,
		))
	}

	// Get tracking information
	tracking, err := client.GetOrderTracking(orderID)
	if err != nil {
		return printer.PrintError(err)
	}

	// Output tracking information
	return printer.Print(tracking)
}
