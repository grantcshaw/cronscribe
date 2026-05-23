package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// Entry represents a scheduled cron entry with metadata.
type Entry struct {
	ID          cron.EntryID
	Expression  string
	Description string
	Next        time.Time
	Prev        time.Time
}

// Scheduler wraps a cron runner with named entries.
type Scheduler struct {
	c       *cron.Cron
	entries map[string]cron.EntryID
}

// New creates a new Scheduler using local time.
func New() *Scheduler {
	return &Scheduler{
		c:       cron.New(),
		entries: make(map[string]cron.EntryID),
	}
}

// Add registers a cron expression with a name and a job function.
// Returns an error if the expression is invalid or the name is already registered.
func (s *Scheduler) Add(name, expr string, job func()) error {
	if name == "" {
		return fmt.Errorf("scheduler: name must not be empty")
	}
	if _, exists := s.entries[name]; exists {
		return fmt.Errorf("scheduler: entry %q already exists", name)
	}
	id, err := s.c.AddFunc(expr, job)
	if err != nil {
		return fmt.Errorf("scheduler: invalid expression %q: %w", expr, err)
	}
	s.entries[name] = id
	return nil
}

// Remove unregisters a named entry.
func (s *Scheduler) Remove(name string) error {
	id, exists := s.entries[name]
	if !exists {
		return fmt.Errorf("scheduler: entry %q not found", name)
	}
	s.c.Remove(id)
	delete(s.entries, name)
	return nil
}

// Get returns the Entry metadata for a named schedule.
func (s *Scheduler) Get(name string) (Entry, error) {
	id, exists := s.entries[name]
	if !exists {
		return Entry{}, fmt.Errorf("scheduler: entry %q not found", name)
	}
	e := s.c.Entry(id)
	return Entry{
		ID:         e.ID,
		Expression: name,
		Next:       e.Next,
		Prev:       e.Prev,
	}, nil
}

// Start begins the scheduler in the background.
func (s *Scheduler) Start() { s.c.Start() }

// Stop halts the scheduler gracefully.
func (s *Scheduler) Stop() { s.c.Stop() }

// Names returns all registered entry names.
func (s *Scheduler) Names() []string {
	names := make([]string, 0, len(s.entries))
	for n := range s.entries {
		names = append(names, n)
	}
	return names
}
