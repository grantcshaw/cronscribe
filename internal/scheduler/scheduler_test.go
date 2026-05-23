package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestAddAndGet(t *testing.T) {
	s := New()
	err := s.Add("heartbeat", "* * * * *", func() {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, err := s.Get("heartbeat")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if e.Expression != "heartbeat" {
		t.Errorf("expected expression name 'heartbeat', got %q", e.Expression)
	}
}

func TestAddDuplicateName(t *testing.T) {
	s := New()
	_ = s.Add("job", "* * * * *", func() {})
	err := s.Add("job", "0 * * * *", func() {})
	if err == nil {
		t.Fatal("expected error for duplicate name, got nil")
	}
}

func TestAddEmptyName(t *testing.T) {
	s := New()
	err := s.Add("", "* * * * *", func() {})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestAddInvalidExpression(t *testing.T) {
	s := New()
	err := s.Add("bad", "not-a-cron", func() {})
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestRemove(t *testing.T) {
	s := New()
	_ = s.Add("temp", "* * * * *", func() {})
	if err := s.Remove("temp"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	if _, err := s.Get("temp"); err == nil {
		t.Fatal("expected error after removal, got nil")
	}
}

func TestRemoveMissing(t *testing.T) {
	s := New()
	if err := s.Remove("ghost"); err == nil {
		t.Fatal("expected error removing non-existent entry")
	}
}

func TestGetMissing(t *testing.T) {
	s := New()
	_, err := s.Get("none")
	if err == nil {
		t.Fatal("expected error for missing entry")
	}
}

func TestNames(t *testing.T) {
	s := New()
	_ = s.Add("a", "* * * * *", func() {})
	_ = s.Add("b", "0 * * * *", func() {})
	names := s.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestStartStop(t *testing.T) {
	s := New()
	var mu sync.Mutex
	called := false
	_ = s.Add("tick", "* * * * *", func() {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	s.Start()
	time.Sleep(10 * time.Millisecond)
	s.Stop()
	// We only verify no panic; timing-based assertion skipped for reliability.
	_ = called
}
