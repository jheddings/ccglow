package statusline

import (
	"encoding/json"
	"strings"
)

// ParseSession parses session JSON from stdin. Returns nil on invalid input.
func ParseSession(input string) *SessionData {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	var session SessionData
	if err := json.Unmarshal([]byte(input), &session); err != nil {
		return nil
	}

	if session.CWD == "" {
		return nil
	}

	return &session
}
