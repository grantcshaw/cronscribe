package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/user/cronscribe/internal/exporter"
	"github.com/user/cronscribe/internal/formatter"
	"github.com/user/cronscribe/internal/parser"
	"github.com/user/cronscribe/internal/preview"
)

// ExportConfig holds options for the export command.
type ExportConfig struct {
	OutputPath string
	Format     exporter.Format
	NextRunCount int
}

// DefaultExportConfig returns sensible defaults for export.
func DefaultExportConfig() ExportConfig {
	return ExportConfig{
		OutputPath:   "cronscribe_export.json",
		Format:       exporter.FormatJSON,
		NextRunCount: 3,
	}
}

// RunExport parses the given inputs and exports them to a file.
func RunExport(inputs []string, cfg ExportConfig, out io.Writer) error {
	if len(inputs) == 0 {
		return fmt.Errorf("export: no inputs provided")
	}

	entries := make([]exporter.Entry, 0, len(inputs))
	for _, input := range inputs {
		expr, err := parser.Parse(input)
		if err != nil {
			fmt.Fprintf(out, "skipping %q: %v\n", input, err)
			continue
		}

		desc := ""
		if d, err := formatter.Describe(expr); err == nil {
			desc = d
		}

		runs, _ := preview.NextRuns(expr, cfg.NextRunCount)
		formatted := preview.FormatRuns(runs)

		entries = append(entries, exporter.Entry{
			Input:       input,
			Expression:  expr,
			Description: desc,
			NextRuns:    formatted,
			ExportedAt:  time.Now().UTC(),
		})
	}

	if len(entries) == 0 {
		return fmt.Errorf("export: no valid schedules to export")
	}

	if err := exporter.Export(entries, cfg.OutputPath, cfg.Format); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	fmt.Fprintf(out, "Exported %d schedule(s) to %s\n", len(entries), cfg.OutputPath)
	return nil
}
