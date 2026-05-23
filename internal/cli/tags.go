package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/cronscribe/internal/tags"
)

// TagsConfig holds configuration for the tags subcommand.
type TagsConfig struct {
	StorePath string
	Out       io.Writer
}

// DefaultTagsConfig returns a TagsConfig with sensible defaults.
func DefaultTagsConfig() TagsConfig {
	home, _ := os.UserHomeDir()
	return TagsConfig{
		StorePath: home + "/.cronscribe_tags.json",
		Out:       os.Stdout,
	}
}

// RunTags executes a tags operation: add, remove, get, or list.
// action is one of "add", "remove", "get", "list".
func RunTags(cfg TagsConfig, action, tag, expr string) error {
	store := tags.New(cfg.StorePath)
	if err := store.Load(); err != nil {
		return fmt.Errorf("loading tags: %w", err)
	}

	switch strings.ToLower(action) {
	case "add":
		if err := store.Add(tag, expr); err != nil {
			return err
		}
		fmt.Fprintf(cfg.Out, "Tagged %q as %q\n", expr, tag)

	case "remove":
		if err := store.Remove(tag, expr); err != nil {
			return err
		}
		fmt.Fprintf(cfg.Out, "Removed %q from tag %q\n", expr, tag)

	case "get":
		exprs, ok := store.Get(tag)
		if !ok {
			fmt.Fprintf(cfg.Out, "No expressions found for tag %q\n", tag)
			return nil
		}
		for _, e := range exprs {
			fmt.Fprintln(cfg.Out, e)
		}

	case "list":
		list := store.List()
		if len(list) == 0 {
			fmt.Fprintln(cfg.Out, "No tags defined.")
			return nil
		}
		for _, t := range list {
			exprs, _ := store.Get(t)
			fmt.Fprintf(cfg.Out, "%s (%d)\n", t, len(exprs))
		}

	default:
		return fmt.Errorf("unknown action %q: use add, remove, get, or list", action)
	}
	return nil
}
