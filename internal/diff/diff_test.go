package diff

import (
	"strings"
	"testing"
)

func TestCompareNoChange(t *testing.T) {
	expr := "0 9 * * 1"
	r, err := Compare(expr, expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Changed {
		t.Error("expected no change")
	}
	for _, f := range r.Fields {
		if f.Changed {
			t.Errorf("field %q should not be changed", f.Name)
		}
	}
}

func TestCompareWithChange(t *testing.T) {
	r, err := Compare("0 9 * * 1", "30 10 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Changed {
		t.Fatal("expected changes")
	}
	if r.Fields[0].OldVal != "0" || r.Fields[0].NewVal != "30" {
		t.Errorf("minute diff wrong: got %v", r.Fields[0])
	}
	if r.Fields[1].OldVal != "9" || r.Fields[1].NewVal != "10" {
		t.Errorf("hour diff wrong: got %v", r.Fields[1])
	}
	if r.Fields[4].Changed {
		t.Error("day-of-week should not have changed")
	}
}

func TestCompareInvalidOld(t *testing.T) {
	_, err := Compare("bad", "0 9 * * 1")
	if err == nil {
		t.Error("expected error for invalid old expression")
	}
}

func TestCompareInvalidNew(t *testing.T) {
	_, err := Compare("0 9 * * 1", "bad expr here extra")
	if err == nil {
		t.Error("expected error for invalid new expression")
	}
}

func TestSummaryNoChange(t *testing.T) {
	r, _ := Compare("* * * * *", "* * * * *")
	got := Summary(r)
	if got != "No changes between expressions." {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummaryWithChange(t *testing.T) {
	r, _ := Compare("0 9 * * *", "0 12 * * *")
	got := Summary(r)
	if !strings.Contains(got, "hour") {
		t.Errorf("expected 'hour' in summary, got: %q", got)
	}
	if !strings.Contains(got, "9") || !strings.Contains(got, "12") {
		t.Errorf("expected old/new values in summary, got: %q", got)
	}
}
