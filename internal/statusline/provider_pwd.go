package statusline

import (
	"os"
	"path/filepath"
	"strings"
)

// PwdData holds resolved working directory information.
type PwdData struct {
	Name  string
	Path  string
	Smart string
}

type pwdProvider struct{}

func (p *pwdProvider) Name() string { return "pwd" }

func (p *pwdProvider) Resolve(session *SessionData) (any, error) {
	cwd := session.CWD
	name := filepath.Base(cwd)
	dir := filepath.Dir(cwd)
	if dir != "/" {
		dir += "/"
	}

	return &PwdData{
		Name:  name,
		Path:  dir,
		Smart: smartPrefix(cwd),
	}, nil
}

func smartPrefix(cwd string) string {
	if cwd == "/" {
		return ""
	}

	home, _ := os.UserHomeDir()

	display := cwd
	if home != "" && strings.HasPrefix(cwd, home) {
		display = "~" + cwd[len(home):]
	}

	dir := filepath.Dir(display)
	if dir == "." || dir == "/" {
		return ""
	}

	parts := strings.Split(dir, "/")
	var segments []string
	for _, p := range parts {
		if p != "" {
			segments = append(segments, p)
		}
	}

	root := ""
	if strings.HasPrefix(dir, "~") {
		root = "~/"
		if len(segments) > 0 && segments[0] == "~" {
			segments = segments[1:]
		}
	} else if strings.HasPrefix(dir, "/") {
		root = "/"
	}

	if len(segments) <= 2 {
		return root + strings.Join(segments, "/") + "/"
	}

	// Abbreviate: first char of leading parts, then ellipsis
	var abbrev []string
	for i := 0; i < len(segments)-1 && i < 2; i++ {
		if len(segments[i]) > 0 {
			abbrev = append(abbrev, string(segments[i][0]))
		}
	}
	abbrev = append(abbrev, "\u2026")

	return root + strings.Join(abbrev, "/") + "/"
}
