package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func effortValues(result *types.ProviderResult) map[string]any {
	return result.Values["effort"].(map[string]any)
}

func TestEffortProvider(t *testing.T) {
	p := &effortProvider{}
	sess := &types.SessionData{
		CWD:      "/tmp",
		Effort:   &types.EffortInfo{Level: "high"},
		Thinking: &types.ThinkingInfo{Enabled: true},
		FastMode: true,
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	eff := effortValues(result)
	if eff["level"] != "high" {
		t.Errorf("expected level high, got %v", eff["level"])
	}
	if eff["thinking"] != true {
		t.Errorf("expected thinking true, got %v", eff["thinking"])
	}
	if eff["fast"] != true {
		t.Errorf("expected fast true, got %v", eff["fast"])
	}
}

// The effort object is absent when the model does not support the parameter,
// and thinking/fast_mode are absent booleans. All must fail silent.
func TestEffortProvider_NoData(t *testing.T) {
	p := &effortProvider{}
	result, err := p.Resolve(&types.SessionData{CWD: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}

	eff := effortValues(result)
	if eff["level"] != "" {
		t.Errorf("expected empty level, got %v", eff["level"])
	}
	if eff["thinking"] != false {
		t.Errorf("expected thinking false, got %v", eff["thinking"])
	}
	if eff["fast"] != false {
		t.Errorf("expected fast false, got %v", eff["fast"])
	}
}

// Ultracode is not a distinct level and reports as xhigh; the provider passes
// the level through verbatim rather than interpreting it.
func TestEffortProvider_LevelsPassThrough(t *testing.T) {
	p := &effortProvider{}
	for _, level := range []string{"low", "medium", "high", "xhigh", "max"} {
		result, err := p.Resolve(&types.SessionData{
			CWD:    "/tmp",
			Effort: &types.EffortInfo{Level: level},
		})
		if err != nil {
			t.Fatal(err)
		}
		if got := effortValues(result)["level"]; got != level {
			t.Errorf("level %q = %v, want %q", level, got, level)
		}
	}
}

func TestEffortProvider_Registered(t *testing.T) {
	registry := NewRegistry()
	RegisterBuiltin(registry)
	if _, ok := registry.All()["effort"]; !ok {
		t.Error("effort provider not registered")
	}
}
