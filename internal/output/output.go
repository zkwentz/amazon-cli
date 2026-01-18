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

// NewPrinter creates a new Printer
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

// PrintError outputs an error in the configured format
func (p *Printer) PrintError(err error) error {
	if cliErr, ok := err.(*models.CLIError); ok {
		return p.printJSON(models.ErrorResponse{Error: cliErr})
	}

	// Wrap generic errors as CLIError
	cliErr := models.NewCLIError(
		models.ErrorCodeAmazonError,
		err.Error(),
		nil,
	)
	return p.printJSON(models.ErrorResponse{Error: cliErr})
}

func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (p *Printer) printTable(data interface{}) error {
	// Table formatting can be implemented later with tablewriter
	// For now, fall back to JSON
	return p.printJSON(data)
}

func (p *Printer) printRaw(data interface{}) error {
	fmt.Println(data)
	return nil
}
