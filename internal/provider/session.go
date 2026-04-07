package provider

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jheddings/ccglow/internal/types"
)

type sessionProvider struct{}

func (p *sessionProvider) Name() string { return "session" }

func (p *sessionProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	sess := map[string]any{
		"duration": map[string]any{
			"total":     "",
			"api":       "",
			"total_min": 0,
			"api_min":   0,
		},
		"lines-added":   0,
		"lines-removed": 0,
		"id":            "",
		"name":          readSessionName(session.TranscriptPath),
	}

	result := &types.ProviderResult{
		Values: map[string]any{"session": sess},
	}

	if session.SessionID != "" {
		sess["id"] = session.SessionID
	}

	if session.Cost == nil {
		return result, nil
	}

	dur := sess["duration"].(map[string]any)
	dur["total"] = FormatDuration(session.Cost.TotalDurationMS)
	dur["api"] = FormatDuration(session.Cost.TotalAPIDurationMS)
	dur["total_min"] = int(session.Cost.TotalDurationMS / 60_000)
	dur["api_min"] = int(session.Cost.TotalAPIDurationMS / 60_000)

	if session.Cost.TotalLinesAdded > 0 {
		sess["lines-added"] = session.Cost.TotalLinesAdded
	}
	if session.Cost.TotalLinesRemoved > 0 {
		sess["lines-removed"] = session.Cost.TotalLinesRemoved
	}

	return result, nil
}

// readSessionName scans a Claude Code transcript JSONL file for the most
// recent custom-title entry (set via /rename) and returns its title. Returns
// empty string when the file is missing, unreadable, or has no title.
func readSessionName(path string) string {
	if path == "" {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	type titleEntry struct {
		Type        string `json:"type"`
		CustomTitle string `json:"customTitle"`
	}

	var latest string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry titleEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}
		if entry.Type == "custom-title" && entry.CustomTitle != "" {
			latest = entry.CustomTitle
		}
	}
	return latest
}

// FormatDuration formats milliseconds into a human-readable duration.
func FormatDuration(ms float64) string {
	totalMinutes := int(ms / 60_000)
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
