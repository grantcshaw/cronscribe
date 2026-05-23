package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single history record of a parsed schedule.
type Entry struct {
	Input      string    `json:"input"`
	Cron       string    `json:"cron"`
	ParsedAt   time.Time `json:"parsed_at"`
}

// History manages a list of past schedule entries.
type History struct {
	Entries []Entry `json:"entries"`
	path    string
}

// New creates a new History backed by the given file path.
func New(path string) *History {
	return &History{path: path}
}

// Load reads history entries from disk. Returns an empty history if the file
// does not exist yet.
func (h *History) Load() error {
	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		h.Entries = nil
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, h)
}

// Save persists the current entries to disk, creating parent directories as
// needed.
func (h *History) Save() error {
	if err := os.MkdirAll(filepath.Dir(h.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}

// Add appends a new entry. Duplicate cron expressions for the same input are
// not stored again.
func (h *History) Add(input, cron string) {
	for _, e := range h.Entries {
		if e.Input == input && e.Cron == cron {
			return
		}
	}
	h.Entries = append(h.Entries, Entry{
		Input:    input,
		Cron:     cron,
		ParsedAt: time.Now().UTC(),
	})
}

// Last returns up to n most-recent entries.
func (h *History) Last(n int) []Entry {
	if n <= 0 || len(h.Entries) == 0 {
		return nil
	}
	start := len(h.Entries) - n
	if start < 0 {
		start = 0
	}
	return h.Entries[start:]
}
