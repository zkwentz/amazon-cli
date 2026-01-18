package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
	FormatRaw   OutputFormat = "raw"
)

// Printer handles output formatting
type Printer struct {
	Format OutputFormat
	Quiet  bool
}

// NewPrinter creates a new Printer with the specified format
func NewPrinter(format string, quiet bool) *Printer {
	return &Printer{
		Format: OutputFormat(format),
		Quiet:  quiet,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	switch p.Format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatTable:
		return p.printTable(data)
	case FormatRaw:
		return p.printRaw(data)
	default:
		return p.printJSON(data)
	}
}

// PrintError outputs an error in JSON format
func (p *Printer) PrintError(err error) error {
	// Check if it's already a CLIError
	if cliErr, ok := err.(*models.CLIError); ok {
		return p.printJSON(&models.ErrorResponse{Error: cliErr})
	}

	// Wrap generic errors as AMAZON_ERROR
	cliErr := models.NewCLIError(
		models.ErrCodeAmazonError,
		err.Error(),
		nil,
	)

	return p.printJSON(&models.ErrorResponse{Error: cliErr})
}

// printJSON outputs data as formatted JSON
func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printTable outputs data in table format (basic implementation)
func (p *Printer) printTable(data interface{}) error {
	// For now, just print JSON as table format would require more complex logic
	// This can be enhanced later with a table library
	return p.printJSON(data)
}

// printRaw outputs data in raw string format
func (p *Printer) printRaw(data interface{}) error {
	fmt.Println(data)
	return nil
}
