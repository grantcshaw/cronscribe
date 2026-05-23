package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Entry represents a single exported cron schedule entry.
type Entry struct {
	Input      string    `json:"input"`
	Expression string    `json:"expression"`
	Description string   `json:"description,omitempty"`
	NextRuns   []string  `json:"next_runs,omitempty"`
	ExportedAt time.Time `json:"exported_at"`
}

// Format defines the output format for export.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Export writes entries to the given file path in the specified format.
func Export(entries []Entry, path string, format Format) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("exporter: failed to create file %q: %w", path, err)
	}
	defer f.Close()

	switch format {
	case FormatJSON:
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(entries); err != nil {
			return fmt.Errorf("exporter: failed to encode JSON: %w", err)
		}
	case FormatCSV:
		_, err := fmt.Fprintln(f, "input,expression,description,next_runs,exported_at")
		if err != nil {
			return err
		}
		for _, e := range entries {
			nextRuns := strings.Join(e.NextRuns, "|")
			line := fmt.Sprintf("%s,%s,%s,%s,%s\n",
				quoteCSV(e.Input),
				quoteCSV(e.Expression),
				quoteCSV(e.Description),
				quoteCSV(nextRuns),
				e.ExportedAt.Format(time.RFC3339),
			)
			if _, err := fmt.Fprint(f, line); err != nil {
				return fmt.Errorf("exporter: failed to write CSV row: %w", err)
			}
		}
	default:
		return fmt.Errorf("exporter: unsupported format %q", format)
	}
	return nil
}

func quoteCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}
