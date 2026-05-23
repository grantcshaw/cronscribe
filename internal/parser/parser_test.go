package parser

import (
	"testing"
)

type testCase struct {
	input    string
	expected string
	wantErr  bool
}

func TestParse(t *testing.T) {
	cases := []testCase{
		{"every minute", "* * * * *", false},
		{"every hour", "0 * * * *", false},
		{"every day", "0 0 * * *", false},
		{"daily", "0 0 * * *", false},
		{"every week", "0 0 * * 0", false},
		{"weekly", "0 0 * * 0", false},
		{"every month", "0 0 1 * *", false},
		{"monthly", "0 0 1 * *", false},
		{"every year", "0 0 1 1 *", false},
		{"annually", "0 0 1 1 *", false},
		{"every 15 minutes", "*/15 * * * *", false},
		{"every 2 hours", "0 */2 * * *", false},
		{"every monday", "0 0 * * 1", false},
		{"every friday at 09:30", "30 9 * * 5", false},
		{"at 08:00", "0 8 * * *", false},
		{"at 23:59", "59 23 * * *", false},
		{"fly to the moon", "", true},
		{"every 0 minutes", "*/0 * * * *", false},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			expr, err := Parse(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tc.input, err)
				return
			}
			if got := expr.String(); got != tc.expected {
				t.Errorf("Parse(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}
