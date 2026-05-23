package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunTemplateList(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultTemplateConfig()
	cfg.Output = &buf

	if err := RunTemplate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "NAME") {
		t.Error("expected header row in output")
	}
	if !strings.Contains(out, "daily") {
		t.Error("expected 'daily' template in output")
	}
}

func TestRunTemplateSearch(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultTemplateConfig()
	cfg.Output = &buf
	cfg.Query = "hourly"

	if err := RunTemplate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "hourly") {
		t.Error("expected 'hourly' in search results")
	}
}

func TestRunTemplateGetByName(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultTemplateConfig()
	cfg.Output = &buf
	cfg.Name = "weekly"

	if err := RunTemplate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "0 0 * * 0") {
		t.Error("expected weekly expression in output")
	}
}

func TestRunTemplateGetMissing(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultTemplateConfig()
	cfg.Output = &buf
	cfg.Name = "nonexistent"

	err := RunTemplate(cfg)
	if err == nil {
		t.Error("expected error for missing template")
	}
}

func TestRunTemplateNoResults(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultTemplateConfig()
	cfg.Output = &buf
	cfg.Query = "zzznomatch"

	if err := RunTemplate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No templates found") {
		t.Error("expected no-results message")
	}
}
