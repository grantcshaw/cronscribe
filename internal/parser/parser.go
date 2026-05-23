package parser

import (
	"fmt"
	"strings"

	"github.com/cronscribe/cronscribe/internal/validator"
)

// Result holds the parsed cron expression and its human-readable description.
type Result struct {
	Expression  string
	Description string
}

// Parse converts a human-readable schedule string into a cron expression.
// It validates the resulting expression before returning.
func Parse(input string) (*Result, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	var expr string
	var desc string

	switch {
	case input == "every minute":
		expr = "* * * * *"
		desc = "Every minute"

	case input == "every hour":
		expr = "0 * * * *"
		desc = "At the start of every hour"

	case input == "every day" || input == "daily":
		expr = "0 0 * * *"
		desc = "Every day at midnight"

	case input == "every week" || input == "weekly":
		expr = "0 0 * * 0"
		desc = "Every Sunday at midnight"

	case input == "every month" || input == "monthly":
		expr = "0 0 1 * *"
		desc = "On the 1st of every month at midnight"

	case strings.HasPrefix(input, "every ") && strings.Contains(input, " minutes"):
		step, err := parseEveryNMinutes(input)
		if err != nil {
			return nil, err
		}
		expr = fmt.Sprintf("*/%d * * * *", step)
		desc = fmt.Sprintf("Every %d minutes", step)

	case strings.HasPrefix(input, "at ") && strings.Contains(input, "on "):
		hour, minute, day, err := parseAtTimeOnDay(input)
		if err != nil {
			return nil, err
		}
		expr = fmt.Sprintf("%d %d * * %d", minute, hour, day)
		desc = fmt.Sprintf("At %02d:%02d on %s", hour, minute, dayName(day))

	case strings.HasPrefix(input, "at "):
		hour, minute, err := parseAtTime(input)
		if err != nil {
			return nil, err
		}
		expr = fmt.Sprintf("%d %d * * *", minute, hour)
		desc = fmt.Sprintf("Every day at %02d:%02d", hour, minute)

	default:
		return nil, fmt.Errorf("unrecognized schedule: %q", input)
	}

	if err := validator.Validate(expr); err != nil {
		return nil, fmt.Errorf("generated invalid expression: %w", err)
	}

	return &Result{Expression: expr, Description: desc}, nil
}

func dayName(dow int) string {
	names := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	if dow >= 0 && dow < len(names) {
		return names[dow]
	}
	return "Unknown"
}
