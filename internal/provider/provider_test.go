package provider

import (
	"os"
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{500, "500"},
		{1000, "1K"},
		{24500, "24K"},
		{1000000, "1M"},
		{1500000, "1.5M"},
		{2000000, "2M"},
	}

	for _, tt := range tests {
		result := FormatTokens(tt.input)
		if result != tt.expected {
			t.Errorf("FormatTokens(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		ms       float64
		expected string
	}{
		{0, "0m"},
		{30000, "0m"},
		{60000, "1m"},
		{300000, "5m"},
		{3600000, "1h 0m"},
		{5400000, "1h 30m"},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.ms)
		if result != tt.expected {
			t.Errorf("FormatDuration(%f) = %q, want %q", tt.ms, result, tt.expected)
		}
	}
}

func TestContextProviderRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage:      36,
			RemainingPercentage: 64,
			ContextWindowSize:   1000000,
			CurrentUsage:        &types.CurrentUsage{InputTokens: 100},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Remaining == nil || *data.Remaining != 64 {
		t.Errorf("expected remaining 64, got %v", data.Remaining)
	}
}

func TestContextProviderNoRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Remaining != nil {
		t.Errorf("expected nil remaining, got %v", data.Remaining)
	}
}

func TestContextProviderZeroRemaining(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD:           "/tmp",
		ContextWindow: &types.ContextWindow{},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Remaining != nil {
		t.Errorf("expected nil remaining for zero value with no usage, got %v", data.Remaining)
	}
}

func TestContextProvider(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage:    36,
			ContextWindowSize: 1000000,
			CurrentUsage: &types.CurrentUsage{
				InputTokens:              100,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     300,
			},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Tokens != "600" {
		t.Errorf("expected 600 tokens, got %s", data.Tokens)
	}
	if data.Size != "1M" {
		t.Errorf("expected 1M size, got %s", data.Size)
	}
	if *data.Percent != 36 {
		t.Errorf("expected 36%%, got %d", *data.Percent)
	}
}

func TestContextProviderWithTotalTokens(t *testing.T) {
	p := &contextProvider{}
	inputTokens := 50000
	outputTokens := 8000
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage:    36,
			ContextWindowSize: 1000000,
			TotalInputTokens:  &inputTokens,
			TotalOutputTokens: &outputTokens,
			CurrentUsage: &types.CurrentUsage{
				InputTokens:              100,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     300,
			},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Input != "50K" {
		t.Errorf("expected Input 50K, got %s", data.Input)
	}
	if data.Output != "8K" {
		t.Errorf("expected Output 8K, got %s", data.Output)
	}
}

func TestContextProviderInputFallback(t *testing.T) {
	p := &contextProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		ContextWindow: &types.ContextWindow{
			UsedPercentage: 10,
			CurrentUsage: &types.CurrentUsage{
				InputTokens:              100,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     300,
			},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ContextData)
	if data.Input != "600" {
		t.Errorf("expected Input 600, got %s", data.Input)
	}
	if data.Output != "" {
		t.Errorf("expected empty Output, got %s", data.Output)
	}
}

func TestCostProvider(t *testing.T) {
	p := &costProvider{}
	sess := &types.SessionData{
		CWD:  "/tmp",
		Cost: &types.CostInfo{TotalCostUSD: 12.5},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*CostData)
	if *data.USD != "$12.50" {
		t.Errorf("expected $12.50, got %s", *data.USD)
	}
}

func TestModelProvider(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{
		CWD:   "/tmp",
		Model: &types.ModelInfo{DisplayName: "Opus 4.6"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ModelData)
	if *data.Name != "Opus 4.6" {
		t.Errorf("expected Opus 4.6, got %s", *data.Name)
	}
}

func TestModelProviderID(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{
		CWD:   "/tmp",
		Model: &types.ModelInfo{ID: "claude-opus-4-6[1m]", DisplayName: "Opus 4.6"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ModelData)
	if data.ID == nil || *data.ID != "claude-opus-4-6[1m]" {
		t.Errorf("expected claude-opus-4-6[1m], got %v", data.ID)
	}
}

func TestModelProviderNoModel(t *testing.T) {
	p := &modelProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ModelData)
	if data.ID != nil {
		t.Errorf("expected nil ID, got %v", data.ID)
	}
	if data.Name != nil {
		t.Errorf("expected nil Name, got %v", data.Name)
	}
}

func TestSessionProvider(t *testing.T) {
	p := &sessionProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		Cost: &types.CostInfo{
			TotalDurationMS:   5400000,
			TotalLinesAdded:   100,
			TotalLinesRemoved: 50,
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*SessionData)
	if *data.Duration != "1h 30m" {
		t.Errorf("expected 1h 30m, got %s", *data.Duration)
	}
	if *data.LinesAdded != 100 {
		t.Errorf("expected 100 lines added, got %d", *data.LinesAdded)
	}
	if *data.LinesRemoved != 50 {
		t.Errorf("expected 50 lines removed, got %d", *data.LinesRemoved)
	}
}

func TestClaudeProvider(t *testing.T) {
	p := &claudeProvider{}
	sess := &types.SessionData{
		CWD:     "/tmp",
		Version: "2.1.75",
		OutputStyle: &types.OutputStyleInfo{Name: "concise"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ClaudeData)
	if data.Version == nil || *data.Version != "2.1.75" {
		t.Errorf("expected version 2.1.75, got %v", data.Version)
	}
	if data.Style == nil || *data.Style != "concise" {
		t.Errorf("expected style concise, got %v", data.Style)
	}
}

func TestClaudeProviderEmpty(t *testing.T) {
	p := &claudeProvider{}
	sess := &types.SessionData{CWD: "/tmp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*ClaudeData)
	if data.Version != nil {
		t.Errorf("expected nil version, got %v", data.Version)
	}
	if data.Style != nil {
		t.Errorf("expected nil style, got %v", data.Style)
	}
}

func TestPwdProvider(t *testing.T) {
	p := &pwdProvider{}
	sess := &types.SessionData{CWD: "/home/user/projects/myapp"}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*PwdData)
	if data.Name != "myapp" {
		t.Errorf("expected myapp, got %s", data.Name)
	}
	if data.Path != "/home/user/projects/" {
		t.Errorf("expected /home/user/projects/, got %s", data.Path)
	}
}

func TestSmartPrefix(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		cwd      string
		expected string
	}{
		// Root and top-level
		{"/", ""},
		{"/tmp", ""},
		{"/usr", ""},

		// Absolute paths (not under home)
		{"/usr/local", "/usr/"},
		{"/usr/local/bin", "/usr/local/"},
		{"/var/log/syslog", "/var/log/"},

		// Home directory itself
		{home, ""},

		// First level under home (the bug case — was producing "~//")
		{home + "/Projects", "~/"},

		// Two levels under home
		{home + "/Projects/myapp", "~/Projects/"},

		// Three levels under home
		{home + "/Projects/myapp/src", "~/Projects/myapp/"},

		// Four levels under home (abbreviation kicks in)
		{home + "/Projects/myapp/src/pkg", "~/P/m/…/"},
	}

	for _, tt := range tests {
		result := smartPrefix(tt.cwd)
		if result != tt.expected {
			t.Errorf("smartPrefix(%q) = %q, want %q", tt.cwd, result, tt.expected)
		}
	}
}
