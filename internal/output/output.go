package output

import (
	"encoding/json"
	"fmt"
	"io"
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
	writer io.Writer
}

// NewPrinter creates a new Printer with the specified format and quiet mode
func NewPrinter(format string, quiet bool) *Printer {
	return &Printer{
		format: OutputFormat(format),
		quiet:  quiet,
		writer: os.Stdout,
	}
}

// NewPrinterWithWriter creates a new Printer with a custom writer (useful for testing)
func NewPrinterWithWriter(format string, quiet bool, writer io.Writer) *Printer {
	return &Printer{
		format: OutputFormat(format),
		quiet:  quiet,
		writer: writer,
	}
}

// Print outputs the data in the configured format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		return nil
	}

	switch p.format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatTable:
		// Table format will be implemented later with tablewriter
		return fmt.Errorf("table format not yet implemented")
	case FormatRaw:
		// Raw format prints the string representation
		_, err := fmt.Fprintf(p.writer, "%v\n", data)
		return err
	default:
		return fmt.Errorf("unknown output format: %s", p.format)
	}
}

// printJSON marshals data to JSON and writes to output
func (p *Printer) printJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintf(p.writer, "%s\n", jsonData)
	return err
}

// PrintError outputs an error in the standard error format
func (p *Printer) PrintError(err error) error {
	if p.quiet {
		return nil
	}

	var cliErr *models.CLIError

	// Check if it's already a CLIError
	if e, ok := err.(*models.CLIError); ok {
		cliErr = e
	} else {
		// Wrap generic errors as general errors
		cliErr = models.NewCLIError("GENERAL_ERROR", err.Error(), nil)
	}

	response := models.ErrorResponse{
		Error: cliErr,
	}

	return p.printJSON(response)
}
