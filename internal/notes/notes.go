package notes

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Note holds a user-defined annotation attached to a cron expression.
type Note struct {
	Expression string    `json:"expression"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}

// Store manages notes persisted to a JSON file.
type Store struct {
	path  string
	notes map[string]Note // keyed by expression
}

// New creates a Store backed by the given file path.
func New(path string) (*Store, error) {
	s := &Store{path: path, notes: make(map[string]Note)}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Set adds or updates the note for the given cron expression.
func (s *Store) Set(expression, text string) error {
	if expression == "" {
		return errors.New("expression must not be empty")
	}
	s.notes[expression] = Note{
		Expression: expression,
		Text:       text,
		CreatedAt:  time.Now().UTC(),
	}
	return s.save()
}

// Get retrieves the note for the given expression.
func (s *Store) Get(expression string) (Note, bool) {
	n, ok := s.notes[expression]
	return n, ok
}

// Delete removes the note for the given expression.
func (s *Store) Delete(expression string) error {
	if _, ok := s.notes[expression]; !ok {
		return errors.New("note not found")
	}
	delete(s.notes, expression)
	return s.save()
}

// All returns every stored note.
func (s *Store) All() []Note {
	out := make([]Note, 0, len(s.notes))
	for _, n := range s.notes {
		out = append(out, n)
	}
	return out
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.notes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.notes)
}
