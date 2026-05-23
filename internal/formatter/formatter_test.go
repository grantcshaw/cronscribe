package formatter

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantMin string
		wantHour string
	}{
		{"valid five fields", "*/5 * * * *", false, "*/5", "*"},
		{"valid daily", "0 9 * * *", false, "0", "9"},
		{"too few fields", "* * *", true, "", ""},
		{"too many fields", "* * * * * *", true, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got.Minute != tt.wantMin {
					t.Errorf("Minute = %q, want %q", got.Minute, tt.wantMin)
				}
				if got.Hour != tt.wantHour {
					t.Errorf("Hour = %q, want %q", got.Hour, tt.wantHour)
				}
			}
		})
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    string
	}{
		{
			name:  "every minute",
			input: "* * * * *",
			want:  "every minute",
		},
		{
			name:  "every 5 minutes",
			input: "*/5 * * * *",
			want:  "every 5 minutes",
		},
		{
			name:  "at specific hour",
			input: "0 9 * * *",
			want:  "at minute 0, past hour 9",
		},
		{
			name:  "on monday",
			input: "0 8 * * 1",
			want:  "at minute 0, past hour 8, on Monday",
		},
		{
			name:    "invalid expression",
			input:   "bad",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Describe(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Describe() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Describe() = %q, want %q", got, tt.want)
			}
		})
	}
}
