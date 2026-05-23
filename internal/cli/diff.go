package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/cronscribe/internal/diff"
)

// DiffConfig holds configuration for the diff subcommand.
type DiffConfig struct {
	Out     io.Writer
	Verbose bool
}

// DefaultDiffConfig returns a DiffConfig writing to stdout.
func DefaultDiffConfig() DiffConfig {
	return DiffConfig{Out: os.Stdout}
}

// RunDiff compares two cron expressions provided as arguments and prints
// a field-level diff to the configured output.
func RunDiff(cfg DiffConfig, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("diff requires two cron expressions as arguments")
	}

	oldExpr := strings.TrimSpace(args[0])
	newExpr := strings.TrimSpace(args[1])

	if oldExpr == "" || newExpr == "" {
		return fmt.Errorf("expressions must not be empty")
	}

	result, err := diff.Compare(oldExpr, newExpr)
	if err != nil {
		return fmt.Errorf("comparison failed: %w", err)
	}

	if cfg.Verbose {
		fmt.Fprintf(cfg.Out, "Old: %s\n", result.Old)
		fmt.Fprintf(cfg.Out, "New: %s\n", result.New)
		fmt.Fprintln(cfg.Out)
	}

	summary := diff.Summary(result)
	fmt.Fprintln(cfg.Out, summary)

	if cfg.Verbose && result.Changed {
		fmt.Fprintln(cfg.Out)
		fmt.Fprintln(cfg.Out, "All fields:")
		for _, f := range result.Fields {
			marker := " "
			if f.Changed {
				marker = "*"
			}
			fmt.Fprintf(cfg.Out, "  %s %-14s %s → %s\n", marker, f.Name+":", f.OldVal, f.NewVal)
		}
	}

	return nil
}
