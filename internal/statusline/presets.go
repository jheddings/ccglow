package statusline

// GetPreset returns the segment tree for a named preset, or nil.
func GetPreset(name string) []SegmentNode {
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

// ListPresets returns all available preset names.
func ListPresets() []string {
	return []string{"default", "minimal", "full"}
}

func defaultPreset() []SegmentNode {
	return []SegmentNode{
		{Type: "pwd.smart", Provider: "pwd", Style: &StyleAttrs{Color: "31"}},
		{Type: "pwd.name", Provider: "pwd", Style: &StyleAttrs{Color: "39", Bold: true}},
		{
			Type:  "group",
			Style: &StyleAttrs{Prefix: " | ", Color: "240"},
			Children: []SegmentNode{
				{Type: "git.branch", Provider: "git", Style: &StyleAttrs{Color: "whiteBright", Bold: true, Prefix: "\ue0a0 "}},
				{Type: "git.insertions", Provider: "git", Style: &StyleAttrs{Color: "green", Prefix: " \u00b7 +"}},
				{Type: "git.deletions", Provider: "git", Style: &StyleAttrs{Color: "red", Prefix: " -"}},
			},
		},
		{
			Type:  "group",
			Style: &StyleAttrs{Prefix: " | "},
			Children: []SegmentNode{
				{Type: "context.tokens", Provider: "context", Style: &StyleAttrs{Color: "white", Bold: true}},
				{Type: "context.percent", Provider: "context", Style: &StyleAttrs{Color: "white", Prefix: " (", Suffix: ")"}},
			},
		},
		{Type: "session.duration", Provider: "session", Style: &StyleAttrs{Color: "magenta", Prefix: " \u00b7 "}},
	}
}

func minimalPreset() []SegmentNode {
	return []SegmentNode{
		{Type: "pwd.name", Provider: "pwd", Style: &StyleAttrs{Color: "39"}},
		{Type: "git.branch", Provider: "git", Style: &StyleAttrs{Color: "whiteBright", Bold: true, Prefix: " | "}},
		{
			Type:  "group",
			Style: &StyleAttrs{Prefix: " | "},
			Children: []SegmentNode{
				{Type: "context.tokens", Provider: "context", Style: &StyleAttrs{Color: "white"}},
				{Type: "context.size", Provider: "context", Style: &StyleAttrs{Color: "white", Prefix: "/"}},
			},
		},
	}
}

func fullPreset() []SegmentNode {
	return []SegmentNode{
		{Type: "pwd.smart", Provider: "pwd", Style: &StyleAttrs{Color: "31"}},
		{Type: "pwd.name", Provider: "pwd", Style: &StyleAttrs{Color: "39", Bold: true}},
		{
			Type:  "group",
			Style: &StyleAttrs{Prefix: " | ", Color: "240"},
			Children: []SegmentNode{
				{Type: "git.branch", Provider: "git", Style: &StyleAttrs{Color: "whiteBright", Bold: true, Prefix: "\ue0a0 "}},
				{Type: "git.insertions", Provider: "git", Style: &StyleAttrs{Color: "green", Prefix: " \u00b7 +"}},
				{Type: "git.deletions", Provider: "git", Style: &StyleAttrs{Color: "red", Prefix: " -"}},
			},
		},
		{Type: "model.name", Provider: "model", Style: &StyleAttrs{Prefix: " | "}},
		{
			Type:  "group",
			Style: &StyleAttrs{Prefix: " \u00b7 "},
			Children: []SegmentNode{
				{Type: "context.tokens", Provider: "context", Style: &StyleAttrs{Color: "white", Bold: true}},
				{Type: "context.size", Provider: "context", Style: &StyleAttrs{Color: "white", Prefix: "/"}},
				{Type: "context.percent", Provider: "context", Style: &StyleAttrs{Color: "white", Prefix: " (", Suffix: ")"}},
			},
		},
		{Type: "cost.usd", Provider: "cost", Style: &StyleAttrs{Color: "yellow", Bold: true, Prefix: " \u00b7 "}},
		{Type: "session.duration", Provider: "session", Style: &StyleAttrs{Color: "magenta", Prefix: " \u00b7 "}},
		{
			Type:  "group",
			Style: &StyleAttrs{Prefix: " \u00b7 "},
			Children: []SegmentNode{
				{Type: "session.lines-added", Provider: "session", Style: &StyleAttrs{Color: "green", Prefix: "+"}},
				{Type: "session.lines-removed", Provider: "session", Style: &StyleAttrs{Color: "red", Prefix: " -"}},
			},
		},
	}
}
