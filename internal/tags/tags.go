package tags

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

// Store manages tags associated with cron expressions.
type Store struct {
	path string
	data map[string][]string // tag -> list of cron expressions
}

// New creates a new Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path, data: make(map[string][]string)}
}

// Add associates a cron expression with a tag.
func (s *Store) Add(tag, expr string) error {
	if tag == "" {
		return errors.New("tag must not be empty")
	}
	if expr == "" {
		return errors.New("expression must not be empty")
	}
	for _, e := range s.data[tag] {
		if e == expr {
			return nil
		}
	}
	s.data[tag] = append(s.data[tag], expr)
	return s.save()
}

// Remove removes a cron expression from a tag.
func (s *Store) Remove(tag, expr string) error {
	list, ok := s.data[tag]
	if !ok {
		return errors.New("tag not found")
	}
	newList := list[:0]
	for _, e := range list {
		if e != expr {
			newList = append(newList, e)
		}
	}
	if len(newList) == 0 {
		delete(s.data, tag)
	} else {
		s.data[tag] = newList
	}
	return s.save()
}

// Get returns all cron expressions associated with a tag.
func (s *Store) Get(tag string) ([]string, bool) {
	exprs, ok := s.data[tag]
	return exprs, ok
}

// List returns all known tags in sorted order.
func (s *Store) List() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Load reads the store from disk.
func (s *Store) Load() error {
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *Store) save() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0644)
}
