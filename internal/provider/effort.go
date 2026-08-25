package provider

import "github.com/jheddings/ccglow/internal/types"

type effortProvider struct{}

func (p *effortProvider) Name() string { return "effort" }

func (p *effortProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	effort := map[string]any{
		"level":    "",
		"thinking": false,
		"fast":     false,
	}

	if session.Effort != nil && session.Effort.Level != "" {
		effort["level"] = session.Effort.Level
	}
	if session.Thinking != nil {
		effort["thinking"] = session.Thinking.Enabled
	}
	effort["fast"] = session.FastMode

	return &types.ProviderResult{
		Values: map[string]any{"effort": effort},
	}, nil
}
