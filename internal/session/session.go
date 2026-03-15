package session

import (
	"encoding/json"
	"strings"

	"github.com/jheddings/ccglow/internal/types"
)

// Parse parses session JSON from stdin. Returns nil on invalid input.
func Parse(input string) *types.SessionData {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	var session types.SessionData
	if err := json.Unmarshal([]byte(input), &session); err != nil {
		return nil
	}

	if session.CWD == "" {
		return nil
	}

	return &session
}
