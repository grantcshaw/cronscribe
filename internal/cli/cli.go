package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/cronscribe/internal/formatter"
	"github.com/user/cronscribe/internal/parser"
	"github.com/user/cronscribe/internal/preview"
	"github.com/user/cronscribe/internal/validator"
)

// Config holds CLI configuration options.
type Config struct {
	PreviewCount int
	Verbose      bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		PreviewCount: 5,
		Verbose:      false,
	}
}

// Run executes the CLI logic for a given human-readable schedule input.
func Run(input string, cfg Config) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return fmt.Errorf("input schedule cannot be empty")
	}

	cronExpr, err := parser.Parse(input)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	if err := validator.Validate(cronExpr); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Cron Expression: %s\n", cronExpr)

	if cfg.Verbose {
		description, err := formatter.Describe(cronExpr)
		if err == nil {
			fmt.Fprintf(os.Stdout, "Description:     %s\n", description)
		}
	}

	runs, err := preview.NextRuns(cronExpr, cfg.PreviewCount)
	if err != nil {
		return fmt.Errorf("preview error: %w", err)
	}

	fmt.Fprintln(os.Stdout, "Next runs:")
	for _, line := range preview.FormatRuns(runs) {
		fmt.Fprintf(os.Stdout, "  %s\n", line)
	}

	return nil
}
