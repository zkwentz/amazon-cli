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
	f := FormatJSON
	if format != "" {
		f = OutputFormat(format)
	}
	return &Printer{
		format: f,
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

// PrintError outputs an error in the configured format
func (p *Printer) PrintError(err error) error {
	var cliErr *models.CLIError
	var ok bool

	if cliErr, ok = err.(*models.CLIError); !ok {
		cliErr = models.NewCLIError(models.ErrorCodeAmazonError, err.Error(), nil)
	}

	response := models.ErrorResponse{Error: cliErr}
	return p.printJSON(response)
}

func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (p *Printer) printTable(data interface{}) error {
	// Table formatting not implemented yet - fallback to JSON
	return p.printJSON(data)
}

func (p *Printer) printRaw(data interface{}) error {
	fmt.Println(data)
	return nil
}
