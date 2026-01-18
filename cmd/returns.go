package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/internal/output"
)

// returnsCmd represents the returns parent command
var returnsCmd = &cobra.Command{
	Use:   "returns",
	Short: "Manage Amazon returns",
	Long:  `Commands for managing Amazon returns including listing returnable items, creating returns, and checking return status.`,
}

// returnsListCmd represents the returns list command
var returnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List returnable items",
	Long:  `List all items eligible for return from your recent Amazon orders.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create printer for output
		printer := output.NewPrinter(outputFormat, quiet)

		// Create Amazon client
		client := amazon.NewClient()

		// Get returnable items
		response, err := client.GetReturnableItems()
		if err != nil {
			printer.PrintError(err)
			output.Exit(err)
			return err
		}

		// Print response
		if err := printer.Print(response); err != nil {
			printer.PrintError(err)
			output.Exit(err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(returnsCmd)
	returnsCmd.AddCommand(returnsListCmd)
}
