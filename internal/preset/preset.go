package preset

import "github.com/jheddings/ccnow/internal/types"

// Get returns the segment tree for a named preset, or nil.
func Get(name string) []types.SegmentNode {
	switch name {
	case "default":
		return defaultPreset()
	case "minimal":
		return minimalPreset()
	case "full":
		return fullPreset()
	default:
		return nil
	}
}

// List returns all available preset names.
func List() []string {
	return []string{"default", "minimal", "full"}
}

func defaultPreset() []types.SegmentNode {
	return []types.SegmentNode{
		{Type: "pwd.smart", Provider: "pwd", Style: &types.StyleAttrs{Color: "31"}},
		{Type: "pwd.name", Provider: "pwd", Style: &types.StyleAttrs{Color: "39", Bold: true}},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " | ", Color: "240"},
			Children: []types.SegmentNode{
				{Type: "git.branch", Provider: "git", Style: &types.StyleAttrs{Color: "whiteBright", Bold: true, Prefix: "\ue0a0 "}},
				{Type: "git.insertions", Provider: "git", Style: &types.StyleAttrs{Color: "green", Prefix: " \u00b7 +"}},
				{Type: "git.deletions", Provider: "git", Style: &types.StyleAttrs{Color: "red", Prefix: " -"}},
			},
		},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Type: "context.tokens", Provider: "context", Style: &types.StyleAttrs{Color: "white", Bold: true}},
				{Type: "context.percent", Provider: "context", Style: &types.StyleAttrs{Color: "white", Prefix: " (", Suffix: ")"}},
			},
		},
		{Type: "session.duration", Provider: "session", Style: &types.StyleAttrs{Color: "magenta", Prefix: " \u00b7 "}},
	}
}

func minimalPreset() []types.SegmentNode {
	return []types.SegmentNode{
		{Type: "pwd.name", Provider: "pwd", Style: &types.StyleAttrs{Color: "39"}},
		{Type: "git.branch", Provider: "git", Style: &types.StyleAttrs{Color: "whiteBright", Bold: true, Prefix: " | "}},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Type: "context.tokens", Provider: "context", Style: &types.StyleAttrs{Color: "white"}},
				{Type: "context.size", Provider: "context", Style: &types.StyleAttrs{Color: "white", Prefix: "/"}},
			},
		},
	}
}

func fullPreset() []types.SegmentNode {
	return []types.SegmentNode{
		{Type: "pwd.smart", Provider: "pwd", Style: &types.StyleAttrs{Color: "31"}},
		{Type: "pwd.name", Provider: "pwd", Style: &types.StyleAttrs{Color: "39", Bold: true}},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " | ", Color: "240"},
			Children: []types.SegmentNode{
				{Type: "git.branch", Provider: "git", Style: &types.StyleAttrs{Color: "whiteBright", Bold: true, Prefix: "\ue0a0 "}},
				{Type: "git.insertions", Provider: "git", Style: &types.StyleAttrs{Color: "green", Prefix: " \u00b7 +"}},
				{Type: "git.deletions", Provider: "git", Style: &types.StyleAttrs{Color: "red", Prefix: " -"}},
			},
		},
		{Type: "model.name", Provider: "model", Style: &types.StyleAttrs{Prefix: " | "}},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " \u00b7 "},
			Children: []types.SegmentNode{
				{Type: "context.tokens", Provider: "context", Style: &types.StyleAttrs{Color: "white", Bold: true}},
				{Type: "context.size", Provider: "context", Style: &types.StyleAttrs{Color: "white", Prefix: "/"}},
				{Type: "context.percent", Provider: "context", Style: &types.StyleAttrs{Color: "white", Prefix: " (", Suffix: ")"}},
			},
		},
		{Type: "cost.usd", Provider: "cost", Style: &types.StyleAttrs{Color: "yellow", Bold: true, Prefix: " \u00b7 "}},
		{Type: "session.duration", Provider: "session", Style: &types.StyleAttrs{Color: "magenta", Prefix: " \u00b7 "}},
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " \u00b7 "},
			Children: []types.SegmentNode{
				{Type: "session.lines-added", Provider: "session", Style: &types.StyleAttrs{Color: "green", Prefix: "+"}},
				{Type: "session.lines-removed", Provider: "session", Style: &types.StyleAttrs{Color: "red", Prefix: " -"}},
			},
		},
	}
}
