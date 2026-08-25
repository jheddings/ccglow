package provider

import (
	"testing"
	"time"

	"github.com/jheddings/ccglow/internal/types"
)

func limitsValues(result *types.ProviderResult) map[string]any {
	return result.Values["limits"].(map[string]any)
}

func limitsWindow(result *types.ProviderResult, name string) map[string]any {
	return limitsValues(result)[name].(map[string]any)
}

func TestLimitsProvider(t *testing.T) {
	// The provider reads time.Now() itself, and epoch seconds truncate, so
	// resets are offset an extra 30s to keep the floored minute stable
	// regardless of where in the second this test runs. formatResetIn is
	// tested exactly, with an injected clock, in TestFormatResetIn.
	now := time.Now()
	fiveHourReset := now.Add(2*time.Hour + 14*time.Minute + 30*time.Second).Unix()
	sevenDayReset := now.Add(30*time.Hour + 30*time.Second).Unix()

	p := &limitsProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		RateLimits: &types.RateLimits{
			FiveHour: &types.RateLimitWindow{UsedPercentage: 23.5, ResetsAt: fiveHourReset},
			SevenDay: &types.RateLimitWindow{UsedPercentage: 41.2, ResetsAt: sevenDayReset},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	session := limitsWindow(result, "session")
	if session["percent"] != 23.5 {
		t.Errorf("expected session percent 23.5, got %v", session["percent"])
	}
	if session["reset_at"] != int(fiveHourReset) {
		t.Errorf("expected session reset_at %d, got %v", fiveHourReset, session["reset_at"])
	}
	if session["reset"] != "2h 14m" {
		t.Errorf("expected session reset 2h 14m, got %v", session["reset"])
	}

	weekly := limitsWindow(result, "weekly")
	if weekly["percent"] != 41.2 {
		t.Errorf("expected weekly percent 41.2, got %v", weekly["percent"])
	}
	if weekly["reset"] != "30h 0m" {
		t.Errorf("expected weekly reset 30h 0m, got %v", weekly["reset"])
	}
}

// rate_limits is absent for API-key users and until the first API response.
func TestLimitsProvider_NoData(t *testing.T) {
	p := &limitsProvider{}
	result, err := p.Resolve(&types.SessionData{CWD: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"session", "weekly"} {
		w := limitsWindow(result, name)
		if w["percent"] != 0.0 {
			t.Errorf("%s: expected percent 0, got %v", name, w["percent"])
		}
		if w["reset"] != "" {
			t.Errorf("%s: expected empty reset, got %v", name, w["reset"])
		}
		if w["reset_at"] != 0 {
			t.Errorf("%s: expected reset_at 0, got %v", name, w["reset_at"])
		}
	}
}

// Each window is documented as independently absent.
func TestLimitsProvider_PartialWindows(t *testing.T) {
	p := &limitsProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		RateLimits: &types.RateLimits{
			FiveHour: &types.RateLimitWindow{UsedPercentage: 60},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if got := limitsWindow(result, "session")["percent"]; got != 60.0 {
		t.Errorf("expected session percent 60, got %v", got)
	}
	if got := limitsWindow(result, "weekly")["percent"]; got != 0.0 {
		t.Errorf("expected weekly percent 0, got %v", got)
	}
}

func TestLimitsProvider_DefaultFormats(t *testing.T) {
	p := &limitsProvider{}
	result, err := p.Resolve(&types.SessionData{CWD: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"limits.session.percent", "limits.weekly.percent"} {
		if result.Formats[key] != "%.0f%%" {
			t.Errorf("expected %s format %%.0f%%%%, got %q", key, result.Formats[key])
		}
	}
}

func TestLimitsProvider_Registered(t *testing.T) {
	registry := NewRegistry()
	RegisterBuiltin(registry)
	if _, ok := registry.All()["limits"]; !ok {
		t.Error("limits provider not registered")
	}
}

func TestFormatResetIn(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	tests := []struct {
		name     string
		resetsAt int64
		expected string
	}{
		{"absent", 0, ""},
		{"two hours out", now.Add(2 * time.Hour).Unix(), "2h 0m"},
		{"forty five minutes out", now.Add(45 * time.Minute).Unix(), "45m"},
		{"already elapsed", now.Add(-time.Hour).Unix(), ""},
		{"exactly now", now.Unix(), ""},
	}

	for _, tt := range tests {
		if got := formatResetIn(tt.resetsAt, now); got != tt.expected {
			t.Errorf("%s: formatResetIn(%d) = %q, want %q", tt.name, tt.resetsAt, got, tt.expected)
		}
	}
}
