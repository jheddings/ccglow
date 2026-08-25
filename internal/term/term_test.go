package term

import (
	"os"
	"testing"
)

func TestWidthFromEnv(t *testing.T) {
	t.Setenv("COLUMNS", "123")
	if got := Width(); got != 123 {
		t.Errorf("Width() = %d, want 123", got)
	}
}

func TestWidthFallback(t *testing.T) {
	// Unset COLUMNS and run with stdout not a tty (test runner context).
	os.Unsetenv("COLUMNS")
	got := Width()
	if got <= 0 {
		t.Errorf("Width() = %d, want positive fallback", got)
	}
}

func TestWidthOverride(t *testing.T) {
	t.Setenv("CCGLOW_WIDTH", "55")
	t.Setenv("COLUMNS", "200") // should be ignored
	if got := Width(); got != 55 {
		t.Errorf("CCGLOW_WIDTH override = %d, want 55", got)
	}
}

func TestWidthOffset(t *testing.T) {
	t.Setenv("COLUMNS", "100")
	t.Setenv("CCGLOW_WIDTH_OFFSET", "4")
	if got := Width(); got != 96 {
		t.Errorf("offset 4 from 100 = %d, want 96", got)
	}
}

func TestWidthOffsetIgnoredWhenLarger(t *testing.T) {
	t.Setenv("COLUMNS", "10")
	t.Setenv("CCGLOW_WIDTH_OFFSET", "20")
	if got := Width(); got != 10 {
		t.Errorf("oversized offset should be ignored, got %d", got)
	}
}

func TestWidthInvalidEnv(t *testing.T) {
	t.Setenv("COLUMNS", "not-a-number")
	got := Width()
	if got <= 0 {
		t.Errorf("Width() with bogus env = %d, want positive fallback", got)
	}
}

func TestHeightFromEnv(t *testing.T) {
	t.Setenv("LINES", "42")
	if got := Height(); got != 42 {
		t.Errorf("Height() = %d, want 42", got)
	}
}

func TestHeightOverride(t *testing.T) {
	t.Setenv("CCGLOW_HEIGHT", "7")
	t.Setenv("LINES", "42") // should be ignored
	if got := Height(); got != 7 {
		t.Errorf("CCGLOW_HEIGHT override = %d, want 7", got)
	}
}

func TestHeightFallback(t *testing.T) {
	os.Unsetenv("LINES")
	got := Height()
	if got <= 0 {
		t.Errorf("Height() = %d, want positive fallback", got)
	}
}

func TestHeightInvalidEnv(t *testing.T) {
	t.Setenv("LINES", "not-a-number")
	got := Height()
	if got <= 0 {
		t.Errorf("Height() with bogus env = %d, want positive fallback", got)
	}
}

// The width offset accounts for host chrome around the statusline and must
// not bleed into the row count.
func TestHeightIgnoresWidthOffset(t *testing.T) {
	t.Setenv("LINES", "42")
	t.Setenv("CCGLOW_WIDTH_OFFSET", "4")
	if got := Height(); got != 42 {
		t.Errorf("Height() = %d, want 42", got)
	}
}
