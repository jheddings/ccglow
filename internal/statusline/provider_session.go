package statusline

import "fmt"

// SessionProviderData holds resolved session timing and line-change data.
type SessionProviderData struct {
	Duration     *string
	LinesAdded   *int
	LinesRemoved *int
}

type sessionProvider struct{}

func (p *sessionProvider) Name() string { return "session" }

func (p *sessionProvider) Resolve(session *SessionData) (any, error) {
	data := &SessionProviderData{}
	if session.Cost == nil {
		return data, nil
	}

	dur := FormatDuration(session.Cost.TotalDurationMS)
	data.Duration = &dur

	if session.Cost.TotalLinesAdded > 0 {
		n := session.Cost.TotalLinesAdded
		data.LinesAdded = &n
	}
	if session.Cost.TotalLinesRemoved > 0 {
		n := session.Cost.TotalLinesRemoved
		data.LinesRemoved = &n
	}

	return data, nil
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
