package style

import "testing"

func TestVisibleWidth(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"empty", "", 0},
		{"plain ascii", "hello", 5},
		{"sgr reset", "\x1b[0m", 0},
		{"sgr wrapped", "\x1b[31mhello\x1b[0m", 5},
		{"multiple sgr", "\x1b[1m\x1b[38;5;240mfoo\x1b[0m", 3},
		{"unicode block", "███░░░", 6},
		{"styled with prefix", "\x1b[0m\x1b[1m » main\x1b[0m", 7},
		{"emoji 2-wide", "✏", 2},
		{"clock 2-wide", "⏰", 2},
		{"cjk 2-wide", "你好", 4},
		{"mixed emoji ascii", "x✏y", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VisibleWidth(tt.in)
			if got != tt.want {
				t.Errorf("VisibleWidth(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
