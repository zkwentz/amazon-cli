package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var returnsCmd = &cobra.Command{
	Use:   "returns",
	Short: "Manage Amazon returns",
	Long:  `List returnable items, get return options, create returns, and track return status.`,
}

var returnsOptionsCmd = &cobra.Command{
	Use:   "options <order-id> <item-id>",
	Short: "Get return options for an item",
	Long: `Get available return options for a specific item from an order.

This command displays all available return methods including:
- UPS drop-off locations
- Amazon Locker locations
- USPS drop-off
- Other carrier options

Each option includes details about drop-off locations, fees, and instructions.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		orderID := args[0]
		itemID := args[1]

		printer := output.NewPrinter(outputFormat, quiet)
		client := amazon.NewClient()

		options, err := client.GetReturnOptions(orderID, itemID)
		if err != nil {
			printer.PrintError(err)
			if cliErr, ok := err.(*models.CLIError); ok {
				os.Exit(getExitCode(cliErr.Code))
			}
			os.Exit(1)
		}

		if err := printer.Print(options); err != nil {
			printer.PrintError(err)
			os.Exit(1)
		}
	},
}

var returnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List returnable items",
	Long:  `List all items that are eligible for return.`,
	Run: func(cmd *cobra.Command, args []string) {
		printer := output.NewPrinter(outputFormat, quiet)
		client := amazon.NewClient()

		items, err := client.GetReturnableItems()
		if err != nil {
			printer.PrintError(err)
			os.Exit(1)
		}

		if err := printer.Print(items); err != nil {
			printer.PrintError(err)
			os.Exit(1)
		}
	},
}

func getExitCode(errorCode string) int {
	switch errorCode {
	case models.ErrCodeAuthRequired, models.ErrCodeAuthExpired:
		return 3
	case models.ErrCodeNetworkError:
		return 4
	case models.ErrCodeRateLimited:
		return 5
	case models.ErrCodeNotFound:
		return 6
	case models.ErrCodeInvalidInput:
		return 2
	default:
		return 1
	}
}

func init() {
	rootCmd.AddCommand(returnsCmd)
	returnsCmd.AddCommand(returnsOptionsCmd)
	returnsCmd.AddCommand(returnsListCmd)
}
