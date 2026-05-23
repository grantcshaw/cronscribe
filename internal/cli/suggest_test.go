package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunSuggestEmpty(t *testing.T) {
	var buf bytes.Buffer
	err := RunSuggest("", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Common schedule suggestions") {
		t.Errorf("expected header for empty input, got: %s", out)
	}
}

func TestRunSuggestKnownInput(t *testing.T) {
	var buf bytes.Buffer
	err := RunSuggest("every hour", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "0 * * * *") {
		t.Errorf("expected cron expression in output, got: %s", out)
	}
}

func TestRunSuggestUnknownInput(t *testing.T) {
	var buf bytes.Buffer
	err := RunSuggest("zzzunknown", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No suggestions found") {
		t.Errorf("expected no-match message, got: %s", out)
	}
}

func TestRunSuggestWhitespace(t *testing.T) {
	var buf bytes.Buffer
	err := RunSuggest("   ", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	// Whitespace-only should behave like empty input
	if !strings.Contains(out, "Common schedule suggestions") {
		t.Errorf("expected default suggestions for whitespace input, got: %s", out)
	}
}
