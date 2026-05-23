package suggester

import (
	"strings"
	"testing"
)

func TestSuggestEmpty(t *testing.T) {
	results := Suggest("")
	if len(results) == 0 {
		t.Fatal("expected default suggestions for empty input, got none")
	}
}

func TestSuggestExactMatch(t *testing.T) {
	results := Suggest("every hour")
	if len(results) == 0 {
		t.Fatal("expected at least one suggestion for 'every hour'")
	}
	found := false
	for _, r := range results {
		if r.Cron == "0 * * * *" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected cron '0 * * * *' in suggestions, got %+v", results)
	}
}

func TestSuggestPartialMatch(t *testing.T) {
	results := Suggest("daily")
	if len(results) == 0 {
		t.Fatal("expected suggestions for 'daily'")
	}
}

func TestSuggestNoMatch(t *testing.T) {
	results := Suggest("xyzzy unknown schedule")
	if len(results) != 0 {
		t.Errorf("expected no suggestions for unknown input, got %d", len(results))
	}
}

func TestSuggestCaseInsensitive(t *testing.T) {
	results := Suggest("EVERY MINUTE")
	if len(results) == 0 {
		t.Fatal("expected suggestions for uppercase input")
	}
}

func TestFormat(t *testing.T) {
	s := Suggestion{
		Input:       "every hour",
		Cron:        "0 * * * *",
		Description: "Runs at the start of every hour",
	}
	out := Format(s)
	if !strings.Contains(out, "0 * * * *") {
		t.Errorf("expected cron in formatted output, got: %s", out)
	}
	if !strings.Contains(out, "every hour") {
		t.Errorf("expected input in formatted output, got: %s", out)
	}
}
