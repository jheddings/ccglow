package render

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

// DefaultTerminalWidth is used when no real width can be detected.
const DefaultTerminalWidth = 80

// TerminalWidth detects the terminal width in columns. Resolution order:
//
//  1. $CCGLOW_WIDTH (full override; useful when host TUI controls layout)
//  2. $COLUMNS env var (if set and valid)
//  3. TIOCGWINSZ ioctl on stdout, stderr, then stdin
//  4. Parent process's controlling tty via `ps` + `stty size`
//  5. `tput cols`
//  6. DefaultTerminalWidth (80)
//
// The detected width is then reduced by $CCGLOW_WIDTH_OFFSET (default 0)
// to account for host chrome (e.g. Claude Code renders the statusline
// inside a bordered, padded box that consumes ~4 cells of usable width).
//
// Steps 4 and 5 are necessary because Claude Code's statusline subprocess
// receives piped stdio (no controlling tty on any standard fd) and does
// not export COLUMNS.
func TerminalWidth() int {
	w := detectWidth()
	if off := offsetFromEnv(); off > 0 && w > off {
		w -= off
	}
	return w
}

func detectWidth() int {
	if v := os.Getenv("CCGLOW_WIDTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if v := os.Getenv("COLUMNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	for _, fd := range []uintptr{os.Stdout.Fd(), os.Stderr.Fd(), os.Stdin.Fd()} {
		if ws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ); err == nil && ws.Col > 0 {
			return int(ws.Col)
		}
	}
	if n := widthFromParentTTY(); n > 0 {
		return n
	}
	if n := widthFromTput(); n > 0 {
		return n
	}
	return DefaultTerminalWidth
}

func offsetFromEnv() int {
	v := os.Getenv("CCGLOW_WIDTH_OFFSET")
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

// widthFromParentTTY walks up to the parent process, reads its controlling
// tty name, and runs `stty size` against /dev/<tty>. Returns 0 on any
// failure or if the parent has no controlling tty.
func widthFromParentTTY() int {
	ppid := os.Getppid()
	out, err := exec.Command("ps", "-o", "tty=", "-p", strconv.Itoa(ppid)).Output()
	if err != nil {
		return 0
	}
	tty := strings.TrimSpace(string(out))
	if tty == "" || tty == "?" || tty == "??" {
		return 0
	}
	// stty size reads from stdin; redirect from /dev/<tty>.
	devPath := "/dev/" + tty
	f, err := os.Open(devPath)
	if err != nil {
		return 0
	}
	defer f.Close()
	cmd := exec.Command("stty", "size")
	cmd.Stdin = f
	out, err = cmd.Output()
	if err != nil {
		return 0
	}
	var rows, cols int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(out)), "%d %d", &rows, &cols); err != nil {
		return 0
	}
	return cols
}

// widthFromTput shells out to `tput cols`. Returns 0 on failure.
func widthFromTput() int {
	out, err := exec.Command("tput", "cols").Output()
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0
	}
	return n
}
