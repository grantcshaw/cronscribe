package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/user/cronscribe/internal/template"
)

// TemplateConfig holds configuration for the template command.
type TemplateConfig struct {
	Query  string
	Name   string
	Output io.Writer
}

// DefaultTemplateConfig returns a TemplateConfig writing to stdout.
func DefaultTemplateConfig() TemplateConfig {
	return TemplateConfig{Output: os.Stdout}
}

// RunTemplate lists or searches templates, or looks up a single template by name.
func RunTemplate(cfg TemplateConfig) error {
	store := template.New()

	if cfg.Name != "" {
		t, ok := store.Get(cfg.Name)
		if !ok {
			return fmt.Errorf("template %q not found", cfg.Name)
		}
		printTemplate(cfg.Output, t)
		return nil
	}

	results := store.Search(cfg.Query)
	if len(results) == 0 {
		fmt.Fprintln(cfg.Output, "No templates found.")
		return nil
	}

	w := tabwriter.NewWriter(cfg.Output, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tEXPRESSION\tDESCRIPTION")
	fmt.Fprintln(w, "----\t----------\t-----------")
	for _, t := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\n", t.Name, t.Expression, t.Description)
	}
	return w.Flush()
}

func printTemplate(w io.Writer, t template.Template) {
	fmt.Fprintf(w, "Name:       %s\n", t.Name)
	fmt.Fprintf(w, "Expression: %s\n", t.Expression)
	fmt.Fprintf(w, "Desc:       %s\n", t.Description)
	if len(t.Tags) > 0 {
		fmt.Fprintf(w, "Tags:       %s\n", strings.Join(t.Tags, ", "))
	}
}
