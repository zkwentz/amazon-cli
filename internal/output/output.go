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
	format OutputFormat
	quiet  bool
}

// NewPrinter creates a new output printer
func NewPrinter(format string, quiet bool) *Printer {
	return &Printer{
		format: OutputFormat(format),
		quiet:  quiet,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	switch p.format {
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

// printJSON outputs data as formatted JSON
func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printTable outputs data in a table format (simplified for now)
func (p *Printer) printTable(data interface{}) error {
	// For now, just use JSON format
	// TODO: Implement table formatting using tablewriter
	return p.printJSON(data)
}

// printRaw outputs raw string representation
func (p *Printer) printRaw(data interface{}) error {
	fmt.Println(data)
	return nil
}

// PrintError outputs an error in the configured format
func (p *Printer) PrintError(err error) error {
	if cliErr, ok := err.(*models.CLIError); ok {
		return p.printJSON(map[string]interface{}{
			"error": cliErr,
		})
	}

	// Wrap non-CLIError errors
	wrappedErr := models.NewCLIError(
		models.ErrAmazonError,
		err.Error(),
		nil,
	)

	return p.printJSON(map[string]interface{}{
		"error": wrappedErr,
	})
}
