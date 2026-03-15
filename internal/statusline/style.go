package statusline

import (
	"fmt"
	"strconv"
	"strings"
)

var colorLevel = 1

// SetColorLevel controls ANSI output. 0 disables colors (plain mode).
func SetColorLevel(level int) {
	colorLevel = level
}

var namedColors = map[string]string{
	"black":         "30",
	"red":           "31",
	"green":         "32",
	"yellow":        "33",
	"blue":          "34",
	"magenta":       "35",
	"cyan":          "36",
	"white":         "37",
	"blackBright":   "90",
	"redBright":     "91",
	"greenBright":   "92",
	"yellowBright":  "93",
	"blueBright":    "94",
	"magentaBright": "95",
	"cyanBright":    "96",
	"whiteBright":   "97",
}

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiItalic = "\x1b[3m"
)

// ApplyStyle wraps a value with ANSI escape codes and prefix/suffix.
func ApplyStyle(value string, style *StyleAttrs) string {
	if style == nil {
		return value
	}

	styled := value

	if colorLevel > 0 {
		var mods strings.Builder
		if style.Bold {
			mods.WriteString(ansiBold)
		}
		if style.Italic {
			mods.WriteString(ansiItalic)
		}

		colorCode := resolveColor(style.Color)
		if colorCode != "" || mods.Len() > 0 {
			styled = ansiReset + mods.String() + colorCode + value + ansiReset
		}
	}

	if style.Prefix != "" {
		styled = style.Prefix + styled
	}
	if style.Suffix != "" {
		styled = styled + style.Suffix
	}

	return styled
}

func resolveColor(color string) string {
	if color == "" {
		return ""
	}

	if code, ok := namedColors[color]; ok {
		return fmt.Sprintf("\x1b[%sm", code)
	}

	if strings.HasPrefix(color, "#") && len(color) == 7 {
		r, _ := strconv.ParseInt(color[1:3], 16, 64)
		g, _ := strconv.ParseInt(color[3:5], 16, 64)
		b, _ := strconv.ParseInt(color[5:7], 16, 64)
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	}

	if n, err := strconv.Atoi(color); err == nil && n >= 0 && n <= 255 {
		return fmt.Sprintf("\x1b[38;5;%dm", n)
	}

	return ""
}
