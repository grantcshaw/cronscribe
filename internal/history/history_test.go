package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestAddAndLast(t *testing.T) {
	h := New(tempPath(t))
	h.Add("every minute", "* * * * *")
	h.Add("every hour", "0 * * * *")

	if got := len(h.Entries); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}

	last := h.Last(1)
	if len(last) != 1 || last[0].Input != "every hour" {
		t.Errorf("unexpected last entry: %+v", last)
	}
}

func TestAddDeduplication(t *testing.T) {
	h := New(tempPath(t))
	h.Add("every minute", "* * * * *")
	h.Add("every minute", "* * * * *")

	if got := len(h.Entries); got != 1 {
		t.Errorf("expected 1 entry after dedup, got %d", got)
	}
}

func TestSaveAndLoad(t *testing.T) {
	p := tempPath(t)
	h := New(p)
	h.Add("every minute", "* * * * *")
	h.Add("every hour", "0 * * * *")

	if err := h.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	h2 := New(p)
	if err := h2.Load(); err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(h2.Entries) != 2 {
		t.Errorf("expected 2 entries after reload, got %d", len(h2.Entries))
	}
	if h2.Entries[0].Input != "every minute" {
		t.Errorf("unexpected first entry: %+v", h2.Entries[0])
	}
}

func TestLoadMissingFile(t *testing.T) {
	h := New(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err := h.Load(); err != nil {
		t.Errorf("Load() on missing file should not error, got: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty entries for missing file, got %d", len(h.Entries))
	}
}

func TestLastEmpty(t *testing.T) {
	h := New(tempPath(t))
	if got := h.Last(5); got != nil {
		t.Errorf("expected nil for empty history, got %v", got)
	}
}

func TestSaveCreatesDirectories(t *testing.T) {
	base := t.TempDir()
	p := filepath.Join(base, "a", "b", "history.json")
	h := New(p)
	h.Add("every minute", "* * * * *")
	if err := h.Save(); err != nil {
		t.Fatalf("Save() should create dirs, got error: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Errorf("expected file to exist at %s: %v", p, err)
	}
}
