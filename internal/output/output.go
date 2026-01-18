package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

type OutputFormat string

const (
	FormatJSON  OutputFormat = "json"
	FormatTable OutputFormat = "table"
	FormatRaw   OutputFormat = "raw"
)

type Printer struct {
	format OutputFormat
	quiet  bool
}

func NewPrinter(format string, quiet bool) *Printer {
	if format == "" {
		format = "json"
	}
	return &Printer{
		format: OutputFormat(format),
		quiet:  quiet,
	}
}

func (p *Printer) Print(data interface{}) error {
	switch p.format {
	case FormatJSON:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	case FormatTable, FormatRaw:
		return fmt.Errorf("table and raw formats not yet implemented")
	default:
		return fmt.Errorf("unknown output format: %s", p.format)
	}
}

func (p *Printer) PrintError(err error) error {
	cliErr, ok := err.(*models.CLIError)
	if !ok {
		cliErr = &models.CLIError{
			Code:    models.ErrorCodeAmazonError,
			Message: err.Error(),
		}
	}

	errorOutput := map[string]interface{}{
		"error": cliErr,
	}

	encoder := json.NewEncoder(os.Stderr)
	encoder.SetIndent("", "  ")
	return encoder.Encode(errorOutput)
}
