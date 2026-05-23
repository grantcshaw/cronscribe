package alias

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// Alias maps a short name to a cron expression.
type Alias struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
}

// Store holds a collection of aliases backed by a JSON file.
type Store struct {
	path    string
	aliases map[string]string
}

// New creates a new Store backed by the given file path.
func New(path string) (*Store, error) {
	s := &Store{path: path, aliases: make(map[string]string)}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Set adds or updates an alias.
func (s *Store) Set(name, cron string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return errors.New("alias name must not be empty")
	}
	if cron == "" {
		return errors.New("cron expression must not be empty")
	}
	s.aliases[name] = cron
	return s.save()
}

// Get retrieves the cron expression for a given alias name.
func (s *Store) Get(name string) (string, bool) {
	name = strings.TrimSpace(strings.ToLower(name))
	cron, ok := s.aliases[name]
	return cron, ok
}

// Delete removes an alias by name. Returns false if it did not exist.
func (s *Store) Delete(name string) (bool, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	_, ok := s.aliases[name]
	if !ok {
		return false, nil
	}
	delete(s.aliases, name)
	return true, s.save()
}

// List returns all stored aliases sorted by name.
func (s *Store) List() []Alias {
	result := make([]Alias, 0, len(s.aliases))
	for name, cron := range s.aliases {
		result = append(result, Alias{Name: name, Cron: cron})
	}
	return result
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.aliases)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.aliases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
