package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	searchCategory  string
	searchMinPrice  float64
	searchMaxPrice  float64
	searchPrimeOnly bool
	searchPage      int
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for products on Amazon",
	Long: `Search for products on Amazon with optional filtering.

Examples:
  # Basic search
  amazon-cli search "wireless headphones"

  # Search with price range
  amazon-cli search "laptop" --min-price 500 --max-price 1000

  # Search for Prime-eligible items only
  amazon-cli search "books" --prime-only

  # Search with category
  amazon-cli search "coffee" --category electronics

  # Search with pagination
  amazon-cli search "keyboards" --page 2`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Add flags
	searchCmd.Flags().StringVar(&searchCategory, "category", "", "Product category to filter by")
	searchCmd.Flags().Float64Var(&searchMinPrice, "min-price", 0, "Minimum price filter")
	searchCmd.Flags().Float64Var(&searchMaxPrice, "max-price", 0, "Maximum price filter")
	searchCmd.Flags().BoolVar(&searchPrimeOnly, "prime-only", false, "Show only Prime-eligible items")
	searchCmd.Flags().IntVar(&searchPage, "page", 1, "Page number for results")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Validate inputs
	if searchMinPrice < 0 {
		return models.NewCLIError(
			models.ErrorCodeInvalidInput,
			"min-price must be greater than or equal to 0",
			nil,
		)
	}

	if searchMaxPrice < 0 {
		return models.NewCLIError(
			models.ErrorCodeInvalidInput,
			"max-price must be greater than or equal to 0",
			nil,
		)
	}

	if searchMinPrice > 0 && searchMaxPrice > 0 && searchMinPrice > searchMaxPrice {
		return models.NewCLIError(
			models.ErrorCodeInvalidInput,
			"min-price must be less than or equal to max-price",
			nil,
		)
	}

	if searchPage < 1 {
		return models.NewCLIError(
			models.ErrorCodeInvalidInput,
			"page must be greater than or equal to 1",
			nil,
		)
	}

	// Create search options
	opts := models.SearchOptions{
		Category:  searchCategory,
		MinPrice:  searchMinPrice,
		MaxPrice:  searchMaxPrice,
		PrimeOnly: searchPrimeOnly,
		Page:      searchPage,
	}

	// Create Amazon client and perform search
	client := amazon.NewClient()
	response, err := client.Search(query, opts)
	if err != nil {
		return err
	}

	// Print results
	printer := GetPrinter()
	return printer.Print(response)
}
