package suggester

import (
	"fmt"
	"strings"
)

// Suggestion represents a cron expression suggestion with a description.
type Suggestion struct {
	Input       string
	Cron        string
	Description string
}

// commonPatterns maps keywords to cron expressions.
var commonPatterns = []struct {
	keywords []string
	cron     string
	desc     string
}{
	{[]string{"every minute"}, "* * * * *", "Runs every minute"},
	{[]string{"every hour", "hourly"}, "0 * * * *", "Runs at the start of every hour"},
	{[]string{"every day", "daily"}, "0 0 * * *", "Runs once a day at midnight"},
	{[]string{"every week", "weekly"}, "0 0 * * 0", "Runs once a week on Sunday at midnight"},
	{[]string{"every month", "monthly"}, "0 0 1 * *", "Runs on the first day of every month"},
	{[]string{"every weekday", "weekdays"}, "0 0 * * 1-5", "Runs every weekday (Mon-Fri) at midnight"},
	{[]string{"every weekend"}, "0 0 * * 6,0", "Runs every Saturday and Sunday at midnight"},
	{[]string{"midnight"}, "0 0 * * *", "Runs at midnight every day"},
	{[]string{"noon"}, "0 12 * * *", "Runs at noon every day"},
	{[]string{"every 5 minutes"}, "*/5 * * * *", "Runs every 5 minutes"},
	{[]string{"every 15 minutes"}, "*/15 * * * *", "Runs every 15 minutes"},
	{[]string{"every 30 minutes"}, "*/30 * * * *", "Runs every 30 minutes"},
}

// Suggest returns cron expression suggestions based on a partial or fuzzy input.
func Suggest(input string) []Suggestion {
	normalized := strings.ToLower(strings.TrimSpace(input))
	if normalized == "" {
		return defaultSuggestions()
	}

	var results []Suggestion
	for _, p := range commonPatterns {
		for _, kw := range p.keywords {
			if strings.Contains(kw, normalized) || strings.Contains(normalized, kw) {
				results = append(results, Suggestion{
					Input:       kw,
					Cron:        p.cron,
					Description: p.desc,
				})
				break
			}
		}
	}
	return results
}

// defaultSuggestions returns a small set of common suggestions.
func defaultSuggestions() []Suggestion {
	return []Suggestion{
		{Input: "every hour", Cron: "0 * * * *", Description: "Runs at the start of every hour"},
		{Input: "every day", Cron: "0 0 * * *", Description: "Runs once a day at midnight"},
		{Input: "every week", Cron: "0 0 * * 0", Description: "Runs once a week on Sunday"},
	}
}

// Format renders a suggestion as a human-readable string.
func Format(s Suggestion) string {
	return fmt.Sprintf("%-20s => %-15s # %s", s.Input, s.Cron, s.Description)
}
