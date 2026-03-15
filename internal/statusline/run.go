package statusline

import (
	"fmt"
	"os"
)

// Options configures a statusline render pass.
type Options struct {
	Preset string
	Config string
	Format string
}

// Run is the main pipeline: parse session, resolve providers, render tree.
func Run(opts Options, stdin string) string {
	session := ParseSession(stdin)
	if session == nil {
		return ""
	}

	if opts.Format == "plain" {
		SetColorLevel(0)
	} else {
		SetColorLevel(1)
	}
	defer SetColorLevel(1)

	segments := NewSegmentRegistry()
	RegisterBuiltinSegments(segments)

	providers := NewProviderRegistry()
	RegisterBuiltinProviders(providers)

	tree := resolveTree(opts, segments)

	providerNames := providers.CollectProviderNames(tree)
	providerData := providers.ResolveAll(providerNames, session)

	return RenderTree(tree, segments, session, providerData)
}

func resolveTree(opts Options, segments *SegmentRegistry) []SegmentNode {
	if opts.Config != "" {
		data, err := os.ReadFile(opts.Config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ccnow: failed to load config: %v\n", err)
		} else {
			tree, err := ParseConfig(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ccnow: failed to parse config: %v\n", err)
			} else if len(tree) > 0 {
				return tree
			}
		}
	}

	if tree := GetPreset(opts.Preset); tree != nil {
		return tree
	}

	return GetPreset("default")
}
