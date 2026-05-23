package template

import (
	"fmt"
	"strings"
)

// Template represents a named cron schedule template.
type Template struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Expression  string `json:"expression"`
	Tags        []string `json:"tags"`
}

// Store holds a collection of templates.
type Store struct {
	templates map[string]Template
}

// New creates a new Store pre-loaded with built-in templates.
func New() *Store {
	s := &Store{templates: make(map[string]Template)}
	for _, t := range builtins {
		s.templates[strings.ToLower(t.Name)] = t
	}
	return s
}

// Add inserts or replaces a template by name.
func (s *Store) Add(t Template) error {
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("template name must not be empty")
	}
	if strings.TrimSpace(t.Expression) == "" {
		return fmt.Errorf("template expression must not be empty")
	}
	s.templates[strings.ToLower(t.Name)] = t
	return nil
}

// Get retrieves a template by name (case-insensitive).
func (s *Store) Get(name string) (Template, bool) {
	t, ok := s.templates[strings.ToLower(name)]
	return t, ok
}

// List returns all templates sorted by name.
func (s *Store) List() []Template {
	out := make([]Template, 0, len(s.templates))
	for _, t := range s.templates {
		out = append(out, t)
	}
	return out
}

// Search returns templates whose name, description, or tags contain the query.
func (s *Store) Search(query string) []Template {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return s.List()
	}
	var results []Template
	for _, t := range s.templates {
		if strings.Contains(strings.ToLower(t.Name), q) ||
			strings.Contains(strings.ToLower(t.Description), q) ||
			tagsContain(t.Tags, q) {
			results = append(results, t)
		}
	}
	return results
}

func tagsContain(tags []string, q string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), q) {
			return true
		}
	}
	return false
}
