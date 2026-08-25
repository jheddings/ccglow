package session

import "testing"

func TestParse_Valid(t *testing.T) {
	s := Parse(`{"cwd": "/home/user/project"}`)
	if s == nil {
		t.Fatal("expected non-nil session")
	}
	if s.CWD != "/home/user/project" {
		t.Errorf("expected cwd /home/user/project, got %s", s.CWD)
	}
}

func TestParse_Empty(t *testing.T) {
	if Parse("") != nil {
		t.Error("expected nil for empty input")
	}
	if Parse("   ") != nil {
		t.Error("expected nil for whitespace input")
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	if Parse("not json") != nil {
		t.Error("expected nil for invalid JSON")
	}
}

func TestParse_MissingCWD(t *testing.T) {
	if Parse(`{"model": {}}`) != nil {
		t.Error("expected nil when cwd is missing")
	}
}

func TestParse_FullData(t *testing.T) {
	input := `{
		"cwd": "/tmp",
		"model": {"id": "claude-opus-4-6", "display_name": "Opus 4.6"},
		"cost": {"total_cost_usd": 1.5, "total_duration_ms": 120000},
		"context_window": {
			"used_percentage": 42,
			"context_window_size": 1000000,
			"current_usage": {"input_tokens": 100, "cache_creation_input_tokens": 200, "cache_read_input_tokens": 300}
		}
	}`
	s := Parse(input)
	if s == nil {
		t.Fatal("expected non-nil session")
	}
	if s.Model.DisplayName != "Opus 4.6" {
		t.Errorf("expected model name Opus 4.6, got %s", s.Model.DisplayName)
	}
	if s.ContextWindow.UsedPercentage != 42 {
		t.Errorf("expected 42%% usage, got %d", s.ContextWindow.UsedPercentage)
	}
}

func TestParse_EffortFields(t *testing.T) {
	input := `{
		"cwd": "/tmp",
		"fast_mode": true,
		"effort": {"level": "xhigh"},
		"thinking": {"enabled": true}
	}`
	s := Parse(input)
	if s == nil {
		t.Fatal("expected non-nil session")
	}
	if s.Effort == nil || s.Effort.Level != "xhigh" {
		t.Errorf("expected effort level xhigh, got %+v", s.Effort)
	}
	if s.Thinking == nil || !s.Thinking.Enabled {
		t.Errorf("expected thinking enabled, got %+v", s.Thinking)
	}
	if !s.FastMode {
		t.Error("expected fast mode true")
	}
}

// effort is absent when the model does not support the parameter.
func TestParse_EffortAbsent(t *testing.T) {
	s := Parse(`{"cwd": "/tmp"}`)
	if s == nil {
		t.Fatal("expected non-nil session")
	}
	if s.Effort != nil {
		t.Errorf("expected nil effort, got %+v", s.Effort)
	}
	if s.Thinking != nil {
		t.Errorf("expected nil thinking, got %+v", s.Thinking)
	}
	if s.FastMode {
		t.Error("expected fast mode false")
	}
}

// used_percentage is 0-100 (not 0-1) and resets_at is epoch seconds (not ms).
func TestParse_RateLimits(t *testing.T) {
	input := `{
		"cwd": "/tmp",
		"rate_limits": {
			"five_hour": {"used_percentage": 23.5, "resets_at": 1738425600},
			"seven_day": {"used_percentage": 41.2, "resets_at": 1738857600}
		}
	}`
	s := Parse(input)
	if s == nil {
		t.Fatal("expected non-nil session")
	}
	if s.RateLimits == nil {
		t.Fatal("expected non-nil rate limits")
	}
	if s.RateLimits.FiveHour.UsedPercentage != 23.5 {
		t.Errorf("expected 23.5, got %v", s.RateLimits.FiveHour.UsedPercentage)
	}
	if s.RateLimits.FiveHour.ResetsAt != 1738425600 {
		t.Errorf("expected 1738425600, got %d", s.RateLimits.FiveHour.ResetsAt)
	}
	if s.RateLimits.SevenDay.UsedPercentage != 41.2 {
		t.Errorf("expected 41.2, got %v", s.RateLimits.SevenDay.UsedPercentage)
	}
}

// Each window may be independently absent.
func TestParse_RateLimitsPartial(t *testing.T) {
	s := Parse(`{"cwd": "/tmp", "rate_limits": {"five_hour": {"used_percentage": 10}}}`)
	if s == nil {
		t.Fatal("expected non-nil session")
	}
	if s.RateLimits.FiveHour == nil {
		t.Fatal("expected non-nil five hour window")
	}
	if s.RateLimits.SevenDay != nil {
		t.Errorf("expected nil seven day window, got %+v", s.RateLimits.SevenDay)
	}
}
