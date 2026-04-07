package render

import (
	"os"
	"testing"
)

func TestTerminalWidthFromEnv(t *testing.T) {
	t.Setenv("COLUMNS", "123")
	if got := TerminalWidth(); got != 123 {
		t.Errorf("TerminalWidth() = %d, want 123", got)
	}
}

func TestTerminalWidthFallback(t *testing.T) {
	// Unset COLUMNS and run with stdout not a tty (test runner context).
	os.Unsetenv("COLUMNS")
	got := TerminalWidth()
	if got <= 0 {
		t.Errorf("TerminalWidth() = %d, want positive fallback", got)
	}
}

func TestTerminalWidthOverride(t *testing.T) {
	t.Setenv("CCGLOW_WIDTH", "55")
	t.Setenv("COLUMNS", "200") // should be ignored
	if got := TerminalWidth(); got != 55 {
		t.Errorf("CCGLOW_WIDTH override = %d, want 55", got)
	}
}

func TestTerminalWidthOffset(t *testing.T) {
	t.Setenv("COLUMNS", "100")
	t.Setenv("CCGLOW_WIDTH_OFFSET", "4")
	if got := TerminalWidth(); got != 96 {
		t.Errorf("offset 4 from 100 = %d, want 96", got)
	}
}

func TestTerminalWidthOffsetIgnoredWhenLarger(t *testing.T) {
	t.Setenv("COLUMNS", "10")
	t.Setenv("CCGLOW_WIDTH_OFFSET", "20")
	if got := TerminalWidth(); got != 10 {
		t.Errorf("oversized offset should be ignored, got %d", got)
	}
}

func TestTerminalWidthInvalidEnv(t *testing.T) {
	t.Setenv("COLUMNS", "not-a-number")
	got := TerminalWidth()
	if got <= 0 {
		t.Errorf("TerminalWidth() with bogus env = %d, want positive fallback", got)
	}
}
