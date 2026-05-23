package preview

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// NextRuns returns the next n scheduled run times for a given cron expression,
// starting from the provided reference time.
func NextRuns(expression string, from time.Time, count int) ([]time.Time, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than zero")
	}

	sched, err := cron.ParseStandard(expression)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression %q: %w", expression, err)
	}

	runs := make([]time.Time, 0, count)
	next := from
	for i := 0; i < count; i++ {
		next = sched.Next(next)
		runs = append(runs, next)
	}
	return runs, nil
}

// FormatRuns returns a human-readable list of upcoming run times.
func FormatRuns(expression string, from time.Time, count int) ([]string, error) {
	runs, err := NextRuns(expression, from, count)
	if err != nil {
		return nil, err
	}

	formatted := make([]string, len(runs))
	for i, t := range runs {
		formatted[i] = t.Format("Mon, 02 Jan 2006 15:04:05 MST")
	}
	return formatted, nil
}
