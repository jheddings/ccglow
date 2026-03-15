package statusline

import "testing"

func TestParseSession_Valid(t *testing.T) {
	input := `{"cwd": "/home/user/project"}`
	session := ParseSession(input)
	if session == nil {
		t.Fatal("expected non-nil session")
	}
	if session.CWD != "/home/user/project" {
		t.Errorf("expected cwd /home/user/project, got %s", session.CWD)
	}
}

func TestParseSession_Empty(t *testing.T) {
	if ParseSession("") != nil {
		t.Error("expected nil for empty input")
	}
	if ParseSession("   ") != nil {
		t.Error("expected nil for whitespace input")
	}
}

func TestParseSession_InvalidJSON(t *testing.T) {
	if ParseSession("not json") != nil {
		t.Error("expected nil for invalid JSON")
	}
}

func TestParseSession_MissingCWD(t *testing.T) {
	if ParseSession(`{"model": {}}`) != nil {
		t.Error("expected nil when cwd is missing")
	}
}

func TestParseSession_FullData(t *testing.T) {
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
	session := ParseSession(input)
	if session == nil {
		t.Fatal("expected non-nil session")
	}
	if session.Model.DisplayName != "Opus 4.6" {
		t.Errorf("expected model name Opus 4.6, got %s", session.Model.DisplayName)
	}
	if session.ContextWindow.UsedPercentage != 42 {
		t.Errorf("expected 42%% usage, got %d", session.ContextWindow.UsedPercentage)
	}
}
