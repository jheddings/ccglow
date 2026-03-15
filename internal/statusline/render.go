package statusline

import "strings"

func isEnabled(node *SegmentNode, session *SessionData) bool {
	if node.EnabledFn != nil {
		defer func() { recover() }()
		return node.EnabledFn(session)
	}
	if node.Enabled != nil {
		return *node.Enabled
	}
	return true
}

func renderNode(
	node *SegmentNode,
	segments *SegmentRegistry,
	session *SessionData,
	providerData map[string]any,
) *string {
	if !isEnabled(node, session) {
		return nil
	}

	// Composite node: render children, collapse if all nil
	if len(node.Children) > 0 {
		var parts []string
		for i := range node.Children {
			rendered := renderNode(&node.Children[i], segments, session, providerData)
			if rendered != nil {
				parts = append(parts, *rendered)
			}
		}
		if len(parts) == 0 {
			return nil
		}
		joined := strings.Join(parts, "")
		styled := ApplyStyle(joined, node.Style)
		return &styled
	}

	// Atomic node: look up segment and render
	seg := segments.Get(node.Type)
	if seg == nil {
		return nil
	}

	ctx := &SegmentContext{
		Session: session,
		Props:   node.Props,
	}
	if node.Provider != "" {
		if data, ok := providerData[node.Provider]; ok {
			ctx.Provider = data
		}
	}

	value := seg.Render(ctx)
	if value == nil {
		return nil
	}

	styled := ApplyStyle(*value, node.Style)
	return &styled
}

// RenderTree performs a depth-first traversal of the segment tree,
// resolving each node against the registries and provider data.
func RenderTree(
	tree []SegmentNode,
	segments *SegmentRegistry,
	session *SessionData,
	providerData map[string]any,
) string {
	var parts []string
	for i := range tree {
		rendered := renderNode(&tree[i], segments, session, providerData)
		if rendered != nil {
			parts = append(parts, *rendered)
		}
	}
	return strings.Join(parts, "")
}
