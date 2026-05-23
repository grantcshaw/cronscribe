package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/cronscribe/internal/suggester"
)

// RunSuggest prints cron expression suggestions for the given input to w.
// If input is empty it prints a default set of common schedules.
func RunSuggest(input string, w io.Writer) error {
	trimmed := strings.TrimSpace(input)
	suggestions := suggester.Suggest(trimmed)

	if len(suggestions) == 0 {
		fmt.Fprintf(w, "No suggestions found for %q.\n", trimmed)
		fmt.Fprintln(w, "Try phrases like: 'every hour', 'daily', 'every 15 minutes'")
		return nil
	}

	if trimmed == "" {
		fmt.Fprintln(w, "Common schedule suggestions:")
	} else {
		fmt.Fprintf(w, "Suggestions for %q:\n", trimmed)
	}

	for _, s := range suggestions {
		fmt.Fprintln(w, " ", suggester.Format(s))
	}
	return nil
}
