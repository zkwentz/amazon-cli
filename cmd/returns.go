package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var returnsCmd = &cobra.Command{
	Use:   "returns",
	Short: "Manage product returns",
	Long:  `Commands for listing returnable items, initiating returns, and tracking return status.`,
}

var returnsLabelCmd = &cobra.Command{
	Use:   "label <return-id>",
	Short: "Get return label for a return",
	Long:  `Retrieves the return shipping label URL and instructions for a specific return ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		returnID := args[0]

		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			return err
		}

		client := amazon.NewClient(cfg)
		printer := output.NewPrinter(outputFormat, quiet)

		label, err := client.GetReturnLabel(returnID)
		if err != nil {
			printer.PrintError(err)
			if cliErr, ok := err.(*models.CLIError); ok {
				switch cliErr.Code {
				case models.ErrorCodeAuthRequired, models.ErrorCodeAuthExpired:
					return cmd.Help()
				case models.ErrorCodeNotFound:
					cmd.SilenceUsage = true
					return err
				}
			}
			cmd.SilenceUsage = true
			return err
		}

		return printer.Print(label)
	},
}

var returnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List returnable items",
	Long:  `Lists all items eligible for return.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			return err
		}

		client := amazon.NewClient(cfg)
		printer := output.NewPrinter(outputFormat, quiet)

		items, err := client.GetReturnableItems()
		if err != nil {
			printer.PrintError(err)
			cmd.SilenceUsage = true
			return err
		}

		return printer.Print(map[string]interface{}{
			"returnable_items": items,
		})
	},
}

var returnsOptionsCmd = &cobra.Command{
	Use:   "options <order-id> <item-id>",
	Short: "Get return options for an item",
	Long:  `Retrieves available return methods and options for a specific order item.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		orderID := args[0]
		itemID := args[1]

		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			return err
		}

		client := amazon.NewClient(cfg)
		printer := output.NewPrinter(outputFormat, quiet)

		options, err := client.GetReturnOptions(orderID, itemID)
		if err != nil {
			printer.PrintError(err)
			cmd.SilenceUsage = true
			return err
		}

		return printer.Print(map[string]interface{}{
			"return_options": options,
		})
	},
}

var returnsCreateCmd = &cobra.Command{
	Use:   "create <order-id> <item-id>",
	Short: "Initiate a return",
	Long:  `Creates a new return for a specific order item. Requires --reason and --confirm flags.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		orderID := args[0]
		itemID := args[1]

		reason, _ := cmd.Flags().GetString("reason")
		confirm, _ := cmd.Flags().GetBool("confirm")

		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			return err
		}

		client := amazon.NewClient(cfg)
		printer := output.NewPrinter(outputFormat, quiet)

		if !confirm {
			return printer.Print(map[string]interface{}{
				"dry_run": true,
				"would_return": map[string]string{
					"order_id": orderID,
					"item_id":  itemID,
					"reason":   reason,
				},
				"message": "Add --confirm to execute",
			})
		}

		returnObj, err := client.CreateReturn(orderID, itemID, reason)
		if err != nil {
			printer.PrintError(err)
			cmd.SilenceUsage = true
			return err
		}

		return printer.Print(returnObj)
	},
}

var returnsStatusCmd = &cobra.Command{
	Use:   "status <return-id>",
	Short: "Check return status",
	Long:  `Retrieves the current status of a return.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		returnID := args[0]

		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			return err
		}

		client := amazon.NewClient(cfg)
		printer := output.NewPrinter(outputFormat, quiet)

		status, err := client.GetReturnStatus(returnID)
		if err != nil {
			printer.PrintError(err)
			cmd.SilenceUsage = true
			return err
		}

		return printer.Print(status)
	},
}

func init() {
	rootCmd.AddCommand(returnsCmd)
	returnsCmd.AddCommand(returnsLabelCmd)
	returnsCmd.AddCommand(returnsListCmd)
	returnsCmd.AddCommand(returnsOptionsCmd)
	returnsCmd.AddCommand(returnsCreateCmd)
	returnsCmd.AddCommand(returnsStatusCmd)

	returnsCreateCmd.Flags().String("reason", "", "reason code for return (required)")
	returnsCreateCmd.Flags().Bool("confirm", false, "confirm the return action")
	returnsCreateCmd.MarkFlagRequired("reason")
}
