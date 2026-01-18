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
	JSON  OutputFormat = "json"
	Table OutputFormat = "table"
	Raw   OutputFormat = "raw"
)

// Printer handles output formatting
type Printer struct {
	Format OutputFormat
	Quiet  bool
}

// NewPrinter creates a new Printer instance
func NewPrinter(format string, quiet bool) *Printer {
	return &Printer{
		Format: OutputFormat(format),
		Quiet:  quiet,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	switch p.Format {
	case JSON:
		return p.printJSON(data)
	case Table:
		return p.printTable(data)
	case Raw:
		return p.printRaw(data)
	default:
		return p.printJSON(data)
	}
}

// PrintError outputs an error in JSON format
func (p *Printer) PrintError(err error) error {
	errorOutput := struct {
		Error interface{} `json:"error"`
	}{}

	if cliErr, ok := err.(*models.CLIError); ok {
		errorOutput.Error = cliErr
	} else {
		errorOutput.Error = models.NewCLIError("GENERAL_ERROR", err.Error(), nil)
	}

	data, jsonErr := json.MarshalIndent(errorOutput, "", "  ")
	if jsonErr != nil {
		return jsonErr
	}

	fmt.Fprintln(os.Stderr, string(data))
	return nil
}

func (p *Printer) printJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (p *Printer) printTable(data interface{}) error {
	// For now, just fall back to JSON
	// TODO: Implement table formatting
	return p.printJSON(data)
}

func (p *Printer) printRaw(data interface{}) error {
	fmt.Println(data)
	return nil
}
