package parser

import (
	"fmt"
	"strings"
)

// CronExpression holds the parsed cron fields.
type CronExpression struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

// String returns the standard 5-field cron expression.
func (c CronExpression) String() string {
	return fmt.Sprintf("%s %s %s %s %s",
		c.Minute, c.Hour, c.DayOfMonth, c.Month, c.DayOfWeek)
}

// Parse converts a human-readable schedule string into a CronExpression.
func Parse(input string) (CronExpression, error) {
	normalized := strings.ToLower(strings.TrimSpace(input))

	switch {
	case normalized == "every minute":
		return CronExpression{"*", "*", "*", "*", "*"}, nil
	case normalized == "every hour":
		return CronExpression{"0", "*", "*", "*", "*"}, nil
	case normalized == "every day" || normalized == "daily":
		return CronExpression{"0", "0", "*", "*", "*"}, nil
	case normalized == "every week" || normalized == "weekly":
		return CronExpression{"0", "0", "*", "*", "0"}, nil
	case normalized == "every month" || normalized == "monthly":
		return CronExpression{"0", "0", "1", "*", "*"}, nil
	case normalized == "every year" || normalized == "yearly" || normalized == "annually":
		return CronExpression{"0", "0", "1", "1", "*"}, nil
	}

	if expr, ok := parseDayOfWeek(normalized); ok {
		return expr, nil
	}

	if expr, ok := parseEveryNMinutes(normalized); ok {
		return expr, nil
	}

	if expr, ok := parseAtTime(normalized); ok {
		return expr, nil
	}

	return CronExpression{}, fmt.Errorf("unrecognized schedule: %q", input)
}
