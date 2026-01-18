package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

// OutputFormat represents output format type
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

// NewPrinter creates a new printer
func NewPrinter(format string, quiet bool) *Printer {
	f := FormatJSON
	if format == "table" {
		f = FormatTable
	} else if format == "raw" {
		f = FormatRaw
	}

	return &Printer{
		format: f,
		quiet:  quiet,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		return nil
	}

	switch p.format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatTable:
		// For now, fall back to JSON
		// TODO: implement table formatting with tablewriter
		return p.printJSON(data)
	case FormatRaw:
		fmt.Println(data)
		return nil
	default:
		return p.printJSON(data)
	}
}

func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintError outputs an error in JSON format
func (p *Printer) PrintError(err error) error {
	errorResponse := map[string]interface{}{
		"error": nil,
	}

	if cliErr, ok := err.(*models.CLIError); ok {
		errorResponse["error"] = cliErr
	} else {
		errorResponse["error"] = models.NewCLIError(models.ErrCodeAmazonError, err.Error())
	}

	encoder := json.NewEncoder(os.Stderr)
	encoder.SetIndent("", "  ")
	return encoder.Encode(errorResponse)
}
