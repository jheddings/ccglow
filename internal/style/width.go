package style

import (
	"regexp"

	"github.com/mattn/go-runewidth"
)

// ansiSGR matches Select Graphic Rendition escape sequences (\x1b[...m).
// These are the only escapes ccglow emits, so a tighter regex is sufficient.
var ansiSGR = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// wideOverrides patches runes that runewidth's default rules report as 1
// but which modern terminals (iTerm2, kitty, wezterm, Alacritty) actually
// render at 2 cells — typically narrow-classified emoji like ✏ U+270F.
var wideOverrides = map[rune]bool{
	'\u270F': true, // ✏ pencil
	'\u270D': true, // ✍ writing hand
	'\u2702': true, // ✂ scissors
	'\u2708': true, // ✈ airplane
	'\u2709': true, // ✉ envelope
	'\u270A': true, // ✊ raised fist
	'\u270B': true, // ✋ raised hand
	'\u270C': true, // ✌ victory hand
	'\u2712': true, // ✒ black nib
	'\u2714': true, // ✔ check mark
	'\u2716': true, // ✖ multiplication x
	'\u2733': true, // ✳ eight spoked asterisk
	'\u2734': true, // ✴ eight pointed star
	'\u2744': true, // ❄ snowflake
	'\u2747': true, // ❇ sparkle
	'\u2757': true, // ❗ heavy exclamation
	'\u2764': true, // ❤ heart
}

// VisibleWidth returns the rendered cell width of s after stripping ANSI
// SGR escape sequences. Uses runewidth's default condition with overrides
// for narrow-classified emoji that terminals render at 2 cells.
func VisibleWidth(s string) int {
	if s == "" {
		return 0
	}
	stripped := ansiSGR.ReplaceAllString(s, "")
	width := 0
	for _, r := range stripped {
		if wideOverrides[r] {
			width += 2
			continue
		}
		width += runewidth.RuneWidth(r)
	}
	return width
}
