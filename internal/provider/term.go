package provider

import (
	"github.com/jheddings/ccglow/internal/term"
	"github.com/jheddings/ccglow/internal/types"
)

type termProvider struct{}

func (p *termProvider) Name() string { return "term" }

func (p *termProvider) Resolve(_ *types.SessionData) (*types.ProviderResult, error) {
	return &types.ProviderResult{
		Values: map[string]any{
			"term": map[string]any{
				"columns": term.Width(),
				"lines":   term.Height(),
			},
		},
	}, nil
}
