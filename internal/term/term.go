// Package term detects the size of the terminal ccglow is rendering into.
package term

import (
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

// DefaultWidth and DefaultHeight are used when no real size can be detected.
const (
	DefaultWidth  = 80
	DefaultHeight = 24
)

// Width detects the terminal width in columns. Resolution order:
//
//  1. $CCGLOW_WIDTH (full override; useful when a host TUI controls layout)
//  2. $COLUMNS
//  3. TIOCGWINSZ ioctl on stdout, stderr, then stdin
//  4. DefaultWidth (80)
//
// The detected width is then reduced by $CCGLOW_WIDTH_OFFSET (default 0) to
// account for host chrome (e.g. Claude Code renders the statusline inside a
// bordered, padded box that consumes a few cells of usable width).
//
// Claude Code sets COLUMNS and LINES to the current terminal dimensions before
// running the statusline command, so step 2 is the usual result there. The
// ioctl covers running ccglow directly in a terminal, where stdout is a tty.
// Nothing further is attempted: the statusline subprocess has its output
// captured rather than connected to the terminal, which means `tput cols` and
// other tty queries cannot read the real size from inside it and would report
// a misleading number rather than no number.
func Width() int {
	w := detectWidth()
	if off := envInt("CCGLOW_WIDTH_OFFSET"); off > 0 && w > off {
		w -= off
	}
	return w
}

// Height detects the terminal height in rows, mirroring Width's resolution
// order via $CCGLOW_HEIGHT and $LINES. The width offset does not apply.
func Height() int {
	if n := envInt("CCGLOW_HEIGHT"); n > 0 {
		return n
	}
	if n := envInt("LINES"); n > 0 {
		return n
	}
	if _, rows := winsize(); rows > 0 {
		return rows
	}
	return DefaultHeight
}

func detectWidth() int {
	if n := envInt("CCGLOW_WIDTH"); n > 0 {
		return n
	}
	if n := envInt("COLUMNS"); n > 0 {
		return n
	}
	if cols, _ := winsize(); cols > 0 {
		return cols
	}
	return DefaultWidth
}

// winsize queries the terminal size on the standard file descriptors,
// returning zeroes when none of them is a terminal.
func winsize() (cols, rows int) {
	for _, fd := range []uintptr{os.Stdout.Fd(), os.Stderr.Fd(), os.Stdin.Fd()} {
		if ws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ); err == nil && ws.Col > 0 {
			return int(ws.Col), int(ws.Row)
		}
	}
	return 0, 0
}

// envInt reads a positive integer from an environment variable, returning 0
// when unset, unparseable, or non-positive.
func envInt(name string) int {
	v := os.Getenv(name)
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return 0
	}
	return n
}
