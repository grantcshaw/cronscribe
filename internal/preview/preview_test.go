package preview

import (
	"testing"
	"time"
)

func TestNextRuns(t *testing.T) {
	// Use a fixed reference time: 2024-01-15 12:00:00 UTC (Monday)
	ref := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		expr       string
		count      int
		wantCount  int
		wantErr    bool
	}{
		{
			name:      "every minute",
			expr:      "* * * * *",
			count:     3,
			wantCount: 3,
		},
		{
			name:      "daily at noon",
			expr:      "0 12 * * *",
			count:     5,
			wantCount: 5,
		},
		{
			name:    "invalid expression",
			expr:    "invalid",
			count:   3,
			wantErr: true,
		},
		{
			name:    "zero count",
			expr:    "* * * * *",
			count:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runs, err := NextRuns(tt.expr, ref, tt.count)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(runs) != tt.wantCount {
				t.Errorf("got %d runs, want %d", len(runs), tt.wantCount)
			}
			for i := 1; i < len(runs); i++ {
				if !runs[i].After(runs[i-1]) {
					t.Errorf("run[%d] (%v) is not after run[%d] (%v)", i, runs[i], i-1, runs[i-1])
				}
			}
		})
	}
}

func TestFormatRuns(t *testing.T) {
	ref := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	formatted, err := FormatRuns("0 9 * * 1", ref, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(formatted) != 2 {
		t.Errorf("expected 2 formatted strings, got %d", len(formatted))
	}
	for _, s := range formatted {
		if s == "" {
			t.Error("formatted run time should not be empty")
		}
	}
}
