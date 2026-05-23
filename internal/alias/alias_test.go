package alias

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "aliases.json")
}

func TestSetAndGet(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set("daily", "0 0 * * *"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	cron, ok := s.Get("daily")
	if !ok || cron != "0 0 * * *" {
		t.Errorf("Get: got (%q, %v), want (\"0 0 * * *\", true)", cron, ok)
	}
}

func TestGetMissing(t *testing.T) {
	s, _ := New(tempPath(t))
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected Get to return false for missing alias")
	}
}

func TestSetEmptyName(t *testing.T) {
	s, _ := New(tempPath(t))
	if err := s.Set("", "0 0 * * *"); err == nil {
		t.Error("expected error for empty alias name")
	}
}

func TestDelete(t *testing.T) {
	s, _ := New(tempPath(t))
	_ = s.Set("weekly", "0 0 * * 0")
	removed, err := s.Delete("weekly")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !removed {
		t.Error("expected Delete to return true")
	}
	_, ok := s.Get("weekly")
	if ok {
		t.Error("alias still present after deletion")
	}
}

func TestDeleteMissing(t *testing.T) {
	s, _ := New(tempPath(t))
	removed, err := s.Delete("ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed {
		t.Error("expected Delete to return false for missing alias")
	}
}

func TestPersistence(t *testing.T) {
	path := tempPath(t)
	s1, _ := New(path)
	_ = s1.Set("hourly", "0 * * * *")

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	cron, ok := s2.Get("hourly")
	if !ok || cron != "0 * * * *" {
		t.Errorf("persistence: got (%q, %v)", cron, ok)
	}
}

func TestLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no_such_file.json")
	_, err := New(path)
	if err != nil {
		t.Errorf("expected no error for missing file, got: %v", err)
	}
}

func TestList(t *testing.T) {
	s, _ := New(tempPath(t))
	_ = s.Set("a", "* * * * *")
	_ = s.Set("b", "0 0 * * *")
	list := s.List()
	if len(list) != 2 {
		t.Errorf("List: expected 2 entries, got %d", len(list))
	}
}

func TestCaseInsensitiveName(t *testing.T) {
	s, _ := New(tempPath(t))
	_ = s.Set("Daily", "0 0 * * *")
	cron, ok := s.Get("DAILY")
	if !ok || cron != "0 0 * * *" {
		t.Errorf("case insensitive: got (%q, %v)", cron, ok)
	}
}

func init() {
	_ = os.Setenv // suppress unused import lint
}
