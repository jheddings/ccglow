package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func agentValues(result *types.ProviderResult) map[string]any {
	return result.Values["agent"].(map[string]any)
}

func TestAgentProvider(t *testing.T) {
	p := &agentProvider{}
	sess := &types.SessionData{
		CWD:   "/tmp",
		Agent: &types.AgentInfo{Name: "security-reviewer"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	agent := agentValues(result)
	if agent["name"] != "security-reviewer" {
		t.Errorf("expected name security-reviewer, got %v", agent["name"])
	}
	if agent["active"] != true {
		t.Errorf("expected active true, got %v", agent["active"])
	}
}

// agent is absent unless the session runs with --agent or agent settings.
func TestAgentProvider_NoData(t *testing.T) {
	p := &agentProvider{}
	result, err := p.Resolve(&types.SessionData{CWD: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}

	agent := agentValues(result)
	if agent["name"] != "" {
		t.Errorf("expected empty name, got %v", agent["name"])
	}
	if agent["active"] != false {
		t.Errorf("expected active false, got %v", agent["active"])
	}
}

// An agent object carrying an empty name is not an active agent session.
func TestAgentProvider_EmptyName(t *testing.T) {
	p := &agentProvider{}
	result, err := p.Resolve(&types.SessionData{
		CWD:   "/tmp",
		Agent: &types.AgentInfo{},
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := agentValues(result)["active"]; got != false {
		t.Errorf("expected active false, got %v", got)
	}
}

func TestAgentProvider_Registered(t *testing.T) {
	registry := NewRegistry()
	RegisterBuiltin(registry)
	if _, ok := registry.All()["agent"]; !ok {
		t.Error("agent provider not registered")
	}
}
