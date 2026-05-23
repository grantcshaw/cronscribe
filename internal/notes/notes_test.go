package notes_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/cronscribe/internal/notes"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "notes.json")
}

func TestSetAndGet(t *testing.T) {
	s, err := notes.New(tempPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Set("0 * * * *", "runs every hour"); err != nil {
		t.Fatal(err)
	}
	n, ok := s.Get("0 * * * *")
	if !ok {
		t.Fatal("expected note to exist")
	}
	if n.Text != "runs every hour" {
		t.Errorf("got %q, want %q", n.Text, "runs every hour")
	}
}

func TestGetMissing(t *testing.T) {
	s, _ := notes.New(tempPath(t))
	_, ok := s.Get("0 0 * * *")
	if ok {
		t.Error("expected missing note")
	}
}

func TestSetEmptyExpression(t *testing.T) {
	s, _ := notes.New(tempPath(t))
	if err := s.Set("", "some text"); err == nil {
		t.Error("expected error for empty expression")
	}
}

func TestDelete(t *testing.T) {
	s, _ := notes.New(tempPath(t))
	_ = s.Set("*/5 * * * *", "every 5 min")
	if err := s.Delete("*/5 * * * *"); err != nil {
		t.Fatal(err)
	}
	_, ok := s.Get("*/5 * * * *")
	if ok {
		t.Error("expected note to be deleted")
	}
}

func TestDeleteMissing(t *testing.T) {
	s, _ := notes.New(tempPath(t))
	if err := s.Delete("0 0 * * *"); err == nil {
		t.Error("expected error deleting non-existent note")
	}
}

func TestPersistence(t *testing.T) {
	p := tempPath(t)
	s1, _ := notes.New(p)
	_ = s1.Set("0 9 * * 1", "monday morning")

	s2, err := notes.New(p)
	if err != nil {
		t.Fatal(err)
	}
	n, ok := s2.Get("0 9 * * 1")
	if !ok || n.Text != "monday morning" {
		t.Errorf("persistence failed: got %+v, ok=%v", n, ok)
	}
}

func TestAll(t *testing.T) {
	s, _ := notes.New(tempPath(t))
	_ = s.Set("0 * * * *", "hourly")
	_ = s.Set("0 0 * * *", "daily")
	if len(s.All()) != 2 {
		t.Errorf("expected 2 notes, got %d", len(s.All()))
	}
}

func TestLoadMissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	s, err := notes.New(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.All()) != 0 {
		t.Error("expected empty store for missing file")
	}
	_ = os.Remove(p)
}
