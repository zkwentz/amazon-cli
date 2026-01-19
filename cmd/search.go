package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	searchCategory string
	searchMinPrice float64
	searchMaxPrice float64
	searchPrimeOnly bool
	searchPage     int
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for products",
	Long:  `Search for products on Amazon by keyword with optional filters.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		c := getClient()

		opts := models.SearchOptions{
			Category:  searchCategory,
			MinPrice:  searchMinPrice,
			MaxPrice:  searchMaxPrice,
			PrimeOnly: searchPrimeOnly,
			Page:      searchPage,
		}

		results, err := c.Search(query, opts)
		if err != nil {
			output.Error(models.ErrAmazonError, err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		output.JSON(results)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchCategory, "category", "", "Product category to search in")
	searchCmd.Flags().Float64Var(&searchMinPrice, "min-price", 0, "Minimum price filter")
	searchCmd.Flags().Float64Var(&searchMaxPrice, "max-price", 0, "Maximum price filter")
	searchCmd.Flags().BoolVar(&searchPrimeOnly, "prime-only", false, "Only show Prime eligible items")
	searchCmd.Flags().IntVar(&searchPage, "page", 1, "Results page number")
}
