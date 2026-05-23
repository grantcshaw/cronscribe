package exporter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/cronscribe/internal/exporter"
)

func sampleEntries() []exporter.Entry {
	return []exporter.Entry{
		{
			Input:       "every day at noon",
			Expression:  "0 12 * * *",
			Description: "At 12:00 PM every day",
			NextRuns:    []string{"2024-01-15 12:00", "2024-01-16 12:00"},
			ExportedAt:  time.Date(2024, 1, 14, 9, 0, 0, 0, time.UTC),
		},
		{
			Input:      "every monday at 9am",
			Expression: "0 9 * * 1",
			ExportedAt: time.Date(2024, 1, 14, 9, 0, 0, 0, time.UTC),
		},
	}
}

func TestExportJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.json")
	entries := sampleEntries()

	if err := exporter.Export(entries, tmp, exporter.FormatJSON); err != nil {
		t.Fatalf("Export JSON failed: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var got []exporter.Entry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(got) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(got))
	}
	if got[0].Expression != "0 12 * * *" {
		t.Errorf("unexpected expression: %s", got[0].Expression)
	}
}

func TestExportCSV(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.csv")
	entries := sampleEntries()

	if err := exporter.Export(entries, tmp, exporter.FormatCSV); err != nil {
		t.Fatalf("Export CSV failed: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 { // header + 2 entries
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "input,expression") {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if !strings.Contains(lines[1], "0 12 * * *") {
		t.Errorf("expected expression in CSV row: %s", lines[1])
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.txt")
	err := exporter.Export(sampleEntries(), tmp, exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportInvalidPath(t *testing.T) {
	err := exporter.Export(sampleEntries(), "/nonexistent/dir/out.json", exporter.FormatJSON)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
