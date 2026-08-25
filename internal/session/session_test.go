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
