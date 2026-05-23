package validator

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidationError represents a cron expression validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Message)
}

// Validate checks whether a cron expression string is valid.
// It expects the standard 5-field format: minute hour dom month dow
func Validate(expr string) error {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return &ValidationError{
			Field:   "expression",
			Message: fmt.Sprintf("expected 5 fields, got %d", len(fields)),
		}
	}

	limits := []struct {
		name string
		min  int
		max  int
	}{
		{"minute", 0, 59},
		{"hour", 0, 23},
		{"day-of-month", 1, 31},
		{"month", 1, 12},
		{"day-of-week", 0, 7},
	}

	for i, field := range fields {
		if err := validateField(field, limits[i].name, limits[i].min, limits[i].max); err != nil {
			return err
		}
	}
	return nil
}

func validateField(field, name string, min, max int) error {
	if field == "*" {
		return nil
	}

	// Handle step values like */5 or 1-5/2
	parts := strings.SplitN(field, "/", 2)
	if len(parts) == 2 {
		step, err := strconv.Atoi(parts[1])
		if err != nil || step < 1 {
			return &ValidationError{Field: name, Message: fmt.Sprintf("invalid step value %q", parts[1])}
		}
		field = parts[0]
		if field == "*" {
			return nil
		}
	}

	// Handle ranges like 1-5
	if strings.Contains(field, "-") {
		rangeParts := strings.SplitN(field, "-", 2)
		lo, err1 := strconv.Atoi(rangeParts[0])
		hi, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil || lo > hi {
			return &ValidationError{Field: name, Message: fmt.Sprintf("invalid range %q", field)}
		}
		if lo < min || hi > max {
			return &ValidationError{Field: name, Message: fmt.Sprintf("range %d-%d out of bounds [%d,%d]", lo, hi, min, max)}
		}
		return nil
	}

	// Handle lists like 1,2,3
	for _, part := range strings.Split(field, ",") {
		v, err := strconv.Atoi(part)
		if err != nil {
			return &ValidationError{Field: name, Message: fmt.Sprintf("non-numeric value %q", part)}
		}
		if v < min || v > max {
			return &ValidationError{Field: name, Message: fmt.Sprintf("value %d out of bounds [%d,%d]", v, min, max)}
		}
	}
	return nil
}
