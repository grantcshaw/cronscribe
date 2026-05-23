package template

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	if len(s.List()) == 0 {
		t.Fatal("expected built-in templates to be loaded")
	}
}

func TestGetBuiltin(t *testing.T) {
	s := New()
	tmpl, ok := s.Get("daily")
	if !ok {
		t.Fatal("expected to find 'daily' template")
	}
	if tmpl.Expression != "0 0 * * *" {
		t.Errorf("unexpected expression: %s", tmpl.Expression)
	}
}

func TestGetCaseInsensitive(t *testing.T) {
	s := New()
	_, ok := s.Get("HOURLY")
	if !ok {
		t.Error("expected case-insensitive lookup to succeed")
	}
}

func TestGetMissing(t *testing.T) {
	s := New()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected missing template to return false")
	}
}

func TestAdd(t *testing.T) {
	s := New()
	err := s.Add(Template{Name: "custom", Expression: "*/30 * * * *", Description: "every 30 min"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tmpl, ok := s.Get("custom")
	if !ok || tmpl.Expression != "*/30 * * * *" {
		t.Error("expected custom template to be stored")
	}
}

func TestAddEmptyName(t *testing.T) {
	s := New()
	err := s.Add(Template{Name: "", Expression: "* * * * *"})
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestAddEmptyExpression(t *testing.T) {
	s := New()
	err := s.Add(Template{Name: "bad", Expression: ""})
	if err == nil {
		t.Error("expected error for empty expression")
	}
}

func TestSearch(t *testing.T) {
	s := New()
	results := s.Search("daily")
	if len(results) == 0 {
		t.Error("expected search results for 'daily'")
	}
}

func TestSearchNoQuery(t *testing.T) {
	s := New()
	results := s.Search("")
	if len(results) != len(s.List()) {
		t.Error("empty query should return all templates")
	}
}

func TestSearchByTag(t *testing.T) {
	s := New()
	results := s.Search("business")
	if len(results) == 0 {
		t.Error("expected to find templates tagged 'business'")
	}
}
