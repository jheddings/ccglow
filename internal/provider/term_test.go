package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func termValues(result *types.ProviderResult) map[string]any {
	return result.Values["term"].(map[string]any)
}

func TestTermProvider(t *testing.T) {
	t.Setenv("COLUMNS", "160")
	t.Setenv("LINES", "48")

	p := &termProvider{}
	result, err := p.Resolve(&types.SessionData{CWD: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}

	term := termValues(result)
	if term["columns"] != 160 {
		t.Errorf("expected columns 160, got %v", term["columns"])
	}
	if term["lines"] != 48 {
		t.Errorf("expected lines 48, got %v", term["lines"])
	}
}

// Both values must stay positive so `when` comparisons behave predictably
// even when the host supplies nothing.
func TestTermProvider_Fallback(t *testing.T) {
	p := &termProvider{}
	result, err := p.Resolve(&types.SessionData{CWD: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}

	term := termValues(result)
	if cols, ok := term["columns"].(int); !ok || cols <= 0 {
		t.Errorf("expected positive columns, got %v", term["columns"])
	}
	if lines, ok := term["lines"].(int); !ok || lines <= 0 {
		t.Errorf("expected positive lines, got %v", term["lines"])
	}
}

func TestTermProvider_Registered(t *testing.T) {
	registry := NewRegistry()
	RegisterBuiltin(registry)
	if _, ok := registry.All()["term"]; !ok {
		t.Error("term provider not registered")
	}
}
