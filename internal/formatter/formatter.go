package formatter

import (
	"fmt"
	"strings"
)

// CronExpression holds the parsed components of a cron expression.
type CronExpression struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

// Parse splits a raw cron string into its components.
func Parse(expr string) (*CronExpression, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(parts))
	}
	return &CronExpression{
		Minute:     parts[0],
		Hour:       parts[1],
		DayOfMonth: parts[2],
		Month:      parts[3],
		DayOfWeek:  parts[4],
	}, nil
}

// String reassembles the CronExpression into a standard cron string.
func (c *CronExpression) String() string {
	return fmt.Sprintf("%s %s %s %s %s",
		c.Minute, c.Hour, c.DayOfMonth, c.Month, c.DayOfWeek)
}

// Describe returns a human-readable summary of the cron expression.
func Describe(expr string) (string, error) {
	c, err := Parse(expr)
	if err != nil {
		return "", err
	}

	var parts []string

	// Minute
	if c.Minute == "*" {
		parts = append(parts, "every minute")
	} else if strings.HasPrefix(c.Minute, "*/") {
		parts = append(parts, fmt.Sprintf("every %s minutes", c.Minute[2:]))
	} else {
		parts = append(parts, fmt.Sprintf("at minute %s", c.Minute))
	}

	// Hour
	if c.Hour != "*" {
		if strings.HasPrefix(c.Hour, "*/") {
			parts = append(parts, fmt.Sprintf("every %s hours", c.Hour[2:]))
		} else {
			parts = append(parts, fmt.Sprintf("past hour %s", c.Hour))
		}
	}

	// Day of week
	if c.DayOfWeek != "*" {
		dayNames := map[string]string{
			"0": "Sunday", "1": "Monday", "2": "Tuesday",
			"3": "Wednesday", "4": "Thursday", "5": "Friday", "6": "Saturday",
		}
		if name, ok := dayNames[c.DayOfWeek]; ok {
			parts = append(parts, fmt.Sprintf("on %s", name))
		} else {
			parts = append(parts, fmt.Sprintf("on day-of-week %s", c.DayOfWeek))
		}
	}

	// Day of month
	if c.DayOfMonth != "*" {
		parts = append(parts, fmt.Sprintf("on day %s of the month", c.DayOfMonth))
	}

	return strings.Join(parts, ", "), nil
}
