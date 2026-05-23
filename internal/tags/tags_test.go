package tags

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "tags.json")
}

func TestAddAndGet(t *testing.T) {
	s := New(tempPath(t))
	if err := s.Add("daily", "0 9 * * *"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	exprs, ok := s.Get("daily")
	if !ok || len(exprs) != 1 || exprs[0] != "0 9 * * *" {
		t.Errorf("unexpected Get result: %v %v", ok, exprs)
	}
}

func TestAddDeduplication(t *testing.T) {
	s := New(tempPath(t))
	s.Add("daily", "0 9 * * *")
	s.Add("daily", "0 9 * * *")
	exprs, _ := s.Get("daily")
	if len(exprs) != 1 {
		t.Errorf("expected 1 entry, got %d", len(exprs))
	}
}

func TestAddEmptyTag(t *testing.T) {
	s := New(tempPath(t))
	if err := s.Add("", "0 9 * * *"); err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestRemove(t *testing.T) {
	s := New(tempPath(t))
	s.Add("weekly", "0 9 * * 1")
	s.Add("weekly", "0 10 * * 1")
	if err := s.Remove("weekly", "0 9 * * 1"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	exprs, _ := s.Get("weekly")
	if len(exprs) != 1 || exprs[0] != "0 10 * * 1" {
		t.Errorf("unexpected exprs after remove: %v", exprs)
	}
}

func TestRemoveLastDeletesTag(t *testing.T) {
	s := New(tempPath(t))
	s.Add("solo", "* * * * *")
	s.Remove("solo", "* * * * *")
	if _, ok := s.Get("solo"); ok {
		t.Error("expected tag to be deleted")
	}
}

func TestList(t *testing.T) {
	s := New(tempPath(t))
	s.Add("zebra", "* * * * *")
	s.Add("alpha", "* * * * *")
	list := s.List()
	if len(list) != 2 || list[0] != "alpha" || list[1] != "zebra" {
		t.Errorf("unexpected list order: %v", list)
	}
}

func TestSaveAndLoad(t *testing.T) {
	p := tempPath(t)
	s1 := New(p)
	s1.Add("nightly", "0 0 * * *")

	s2 := New(p)
	if err := s2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	exprs, ok := s2.Get("nightly")
	if !ok || len(exprs) != 1 {
		t.Errorf("expected loaded data, got %v %v", ok, exprs)
	}
}

func TestLoadMissingFile(t *testing.T) {
	s := New(filepath.Join(t.TempDir(), "missing.json"))
	if err := s.Load(); err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestRemoveMissingTag(t *testing.T) {
	s := New(tempPath(t))
	if err := s.Remove("ghost", "* * * * *"); err == nil {
		t.Error("expected error removing from missing tag")
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
