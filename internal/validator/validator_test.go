package validator_test

import (
	"testing"

	"github.com/cronscribe/cronscribe/internal/validator"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{"wildcard all", "* * * * *", false},
		{"every 5 minutes", "*/5 * * * *", false},
		{"specific time", "30 9 * * *", false},
		{"range minutes", "0-30 * * * *", false},
		{"list hours", "0 8,12,18 * * *", false},
		{"weekday", "0 9 * * 1-5", false},
		{"step with range", "0-59/10 * * * *", false},
		{"invalid minute high", "60 * * * *", true},
		{"invalid hour high", "* 24 * * *", true},
		{"invalid dom low", "* * 0 * *", true},
		{"invalid month high", "* * * 13 *", true},
		{"invalid dow high", "* * * * 8", true},
		{"too few fields", "* * * *", true},
		{"too many fields", "* * * * * *", true},
		{"bad range", "10-5 * * * *", true},
		{"non-numeric", "abc * * * *", true},
		{"invalid step", "*/0 * * * *", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%q) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
			}
		})
	}
}

func TestValidationErrorMessage(t *testing.T) {
	err := validator.Validate("60 * * * *")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}
