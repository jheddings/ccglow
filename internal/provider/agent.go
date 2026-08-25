package provider

import "github.com/jheddings/ccglow/internal/types"

type agentProvider struct{}

func (p *agentProvider) Name() string { return "agent" }

func (p *agentProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	agent := map[string]any{
		"name":   "",
		"active": false,
	}

	if session.Agent != nil && session.Agent.Name != "" {
		agent["name"] = session.Agent.Name
		agent["active"] = true
	}

	return &types.ProviderResult{
		Values: map[string]any{"agent": agent},
	}, nil
}
