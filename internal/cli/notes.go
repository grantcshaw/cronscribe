package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/cronscribe/internal/notes"
)

// NotesConfig holds configuration for the notes sub-command.
type NotesConfig struct {
	StorePath  string
	Action     string // "set", "get", "delete", "list"
	Expression string
	Text       string
	Output     io.Writer
}

// DefaultNotesConfig returns a NotesConfig with sensible defaults.
func DefaultNotesConfig() NotesConfig {
	return NotesConfig{
		StorePath: ".cronscribe_notes.json",
		Action:    "list",
		Output:    io.Discard,
	}
}

// RunNotes executes the notes command based on the provided config.
func RunNotes(cfg NotesConfig) error {
	store, err := notes.New(cfg.StorePath)
	if err != nil {
		return fmt.Errorf("opening notes store: %w", err)
	}

	switch strings.ToLower(cfg.Action) {
	case "set":
		if err := store.Set(cfg.Expression, cfg.Text); err != nil {
			return err
		}
		fmt.Fprintf(cfg.Output, "Note saved for %q\n", cfg.Expression)

	case "get":
		n, ok := store.Get(cfg.Expression)
		if !ok {
			return fmt.Errorf("no note found for %q", cfg.Expression)
		}
		fmt.Fprintf(cfg.Output, "%s\n  %s\n", n.Expression, n.Text)

	case "delete":
		if err := store.Delete(cfg.Expression); err != nil {
			return err
		}
		fmt.Fprintf(cfg.Output, "Note deleted for %q\n", cfg.Expression)

	case "list":
		all := store.All()
		if len(all) == 0 {
			fmt.Fprintln(cfg.Output, "No notes stored.")
			return nil
		}
		sort.Slice(all, func(i, j int) bool {
			return all[i].Expression < all[j].Expression
		})
		for _, n := range all {
			fmt.Fprintf(cfg.Output, "%-20s  %s\n", n.Expression, n.Text)
		}

	default:
		return fmt.Errorf("unknown action %q; use set, get, delete, or list", cfg.Action)
	}
	return nil
}
