package cli_test

import (
	"testing"

	"github.com/user/cronscribe/internal/cli"
)

func TestDefaultConfig(t *testing.T) {
	cfg := cli.DefaultConfig()
	if cfg.PreviewCount != 5 {
		t.Errorf("expected PreviewCount=5, got %d", cfg.PreviewCount)
	}
	if cfg.Verbose {
		t.Error("expected Verbose=false by default")
	}
}

func TestRunEmptyInput(t *testing.T) {
	cfg := cli.DefaultConfig()
	err := cli.Run("", cfg)
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestRunWhitespaceInput(t *testing.T) {
	cfg := cli.DefaultConfig()
	err := cli.Run("   ", cfg)
	if err == nil {
		t.Fatal("expected error for whitespace-only input, got nil")
	}
}

func TestRunValidSchedule(t *testing.T) {
	cfg := cli.DefaultConfig()
	cfg.PreviewCount = 3

	tests := []struct {
		input   string
		wantErr bool
	}{
		{"every minute", false},
		{"every hour", false},
		{"every day at midnight", false},
		{"gibberish input xyz", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := cli.Run(tt.input, cfg)
			if tt.wantErr && err == nil {
				t.Errorf("expected error for input %q, got nil", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
			}
		})
	}
}

func TestRunVerbose(t *testing.T) {
	cfg := cli.DefaultConfig()
	cfg.Verbose = true
	cfg.PreviewCount = 2

	err := cli.Run("every hour", cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
