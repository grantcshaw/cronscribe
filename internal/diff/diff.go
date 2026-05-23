package diff

import (
	"fmt"
	"strings"
)

// Field represents a single cron field with its name and value.
type Field struct {
	Name  string
	Value string
}

// Result holds the comparison between two cron expressions.
type Result struct {
	Old    string
	New    string
	Fields []FieldDiff
	Changed bool
}

// FieldDiff represents the difference in a single cron field.
type FieldDiff struct {
	Name    string
	OldVal  string
	NewVal  string
	Changed bool
}

var fieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// Compare takes two cron expression strings and returns a Result
// describing the field-level differences between them.
func Compare(oldExpr, newExpr string) (*Result, error) {
	oldParts, err := split(oldExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid old expression: %w", err)
	}
	newParts, err := split(newExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid new expression: %w", err)
	}

	result := &Result{
		Old:    oldExpr,
		New:    newExpr,
		Fields: make([]FieldDiff, 5),
	}

	for i := 0; i < 5; i++ {
		changed := oldParts[i] != newParts[i]
		result.Fields[i] = FieldDiff{
			Name:    fieldNames[i],
			OldVal:  oldParts[i],
			NewVal:  newParts[i],
			Changed: changed,
		}
		if changed {
			result.Changed = true
		}
	}

	return result, nil
}

// Summary returns a human-readable summary of the diff result.
func Summary(r *Result) string {
	if !r.Changed {
		return "No changes between expressions."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Diff: %q → %q\n", r.Old, r.New))
	for _, f := range r.Fields {
		if f.Changed {
			sb.WriteString(fmt.Sprintf("  %-14s %s → %s\n", f.Name+":", f.OldVal, f.NewVal))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func split(expr string) ([]string, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}
	return parts, nil
}
