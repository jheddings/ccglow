package provider

import (
	"time"

	"github.com/jheddings/ccglow/internal/types"
)

type limitsProvider struct{}

func (p *limitsProvider) Name() string { return "limits" }

func (p *limitsProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	limits := map[string]any{
		"session": emptyLimitWindow(),
		"weekly":  emptyLimitWindow(),
	}

	result := &types.ProviderResult{
		Values: map[string]any{"limits": limits},
		Formats: map[string]string{
			"limits.session.percent": "%.0f%%",
			"limits.weekly.percent":  "%.0f%%",
		},
	}

	if session.RateLimits == nil {
		return result, nil
	}

	now := time.Now()
	applyLimitWindow(limits["session"].(map[string]any), session.RateLimits.FiveHour, now)
	applyLimitWindow(limits["weekly"].(map[string]any), session.RateLimits.SevenDay, now)

	return result, nil
}

func emptyLimitWindow() map[string]any {
	return map[string]any{
		"percent":  0.0,
		"reset":    "",
		"reset_at": 0,
	}
}

// applyLimitWindow fills a window map from stdin data. A nil window leaves the
// zero values in place, since each window is independently absent.
func applyLimitWindow(dst map[string]any, window *types.RateLimitWindow, now time.Time) {
	if window == nil {
		return
	}
	dst["percent"] = window.UsedPercentage
	dst["reset_at"] = int(window.ResetsAt)
	dst["reset"] = formatResetIn(window.ResetsAt, now)
}

// formatResetIn renders the time remaining until an epoch-seconds reset as a
// human-readable duration. Returns empty for an absent reset (0) or one that
// has already elapsed, since statusline updates are event-driven and the
// payload can be stale.
func formatResetIn(resetsAt int64, now time.Time) string {
	if resetsAt <= 0 {
		return ""
	}
	remaining := time.Unix(resetsAt, 0).Sub(now)
	if remaining <= 0 {
		return ""
	}
	return FormatDuration(float64(remaining.Milliseconds()))
}
