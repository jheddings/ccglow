package statusline

import "testing"

func setupTestRegistries() (*SegmentRegistry, *ProviderRegistry) {
	seg := NewSegmentRegistry()
	RegisterBuiltinSegments(seg)
	prov := NewProviderRegistry()
	RegisterBuiltinProviders(prov)
	return seg, prov
}

func TestRenderTree_EmptyTree(t *testing.T) {
	seg, _ := setupTestRegistries()
	session := &SessionData{CWD: "/tmp"}
	result := RenderTree(nil, seg, session, map[string]any{})
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestRenderTree_AtomicNode(t *testing.T) {
	SetColorLevel(0)
	defer SetColorLevel(1)

	seg, _ := setupTestRegistries()
	session := &SessionData{CWD: "/home/user/project"}
	providerData := map[string]any{
		"pwd": &PwdData{Name: "project", Path: "/home/user/", Smart: "~/"},
	}

	tree := []SegmentNode{
		{Type: "pwd.name", Provider: "pwd"},
	}

	result := RenderTree(tree, seg, session, providerData)
	if result != "project" {
		t.Errorf("expected project, got %q", result)
	}
}

func TestRenderTree_CompositeCollapse(t *testing.T) {
	SetColorLevel(0)
	defer SetColorLevel(1)

	seg, _ := setupTestRegistries()
	session := &SessionData{CWD: "/tmp"}
	providerData := map[string]any{
		"git": &GitData{}, // no branch, no diffs
	}

	tree := []SegmentNode{
		{
			Type:  "group",
			Style: &StyleAttrs{Prefix: " | "},
			Children: []SegmentNode{
				{Type: "git.branch", Provider: "git"},
				{Type: "git.insertions", Provider: "git"},
			},
		},
	}

	result := RenderTree(tree, seg, session, providerData)
	if result != "" {
		t.Errorf("expected empty (collapsed composite), got %q", result)
	}
}

func TestRenderTree_DisabledNode(t *testing.T) {
	SetColorLevel(0)
	defer SetColorLevel(1)

	seg, _ := setupTestRegistries()
	session := &SessionData{CWD: "/tmp"}
	providerData := map[string]any{
		"pwd": &PwdData{Name: "tmp"},
	}

	disabled := false
	tree := []SegmentNode{
		{Type: "pwd.name", Provider: "pwd", Enabled: &disabled},
	}

	result := RenderTree(tree, seg, session, providerData)
	if result != "" {
		t.Errorf("expected empty for disabled node, got %q", result)
	}
}

func TestRenderTree_Literal(t *testing.T) {
	SetColorLevel(0)
	defer SetColorLevel(1)

	seg, _ := setupTestRegistries()
	session := &SessionData{CWD: "/tmp"}

	tree := []SegmentNode{
		{Type: "literal", Props: map[string]any{"text": "hello"}},
	}

	result := RenderTree(tree, seg, session, map[string]any{})
	if result != "hello" {
		t.Errorf("expected hello, got %q", result)
	}
}

func TestRenderTree_MissingSegment(t *testing.T) {
	seg, _ := setupTestRegistries()
	session := &SessionData{CWD: "/tmp"}

	tree := []SegmentNode{
		{Type: "nonexistent.segment"},
	}

	result := RenderTree(tree, seg, session, map[string]any{})
	if result != "" {
		t.Errorf("expected empty for missing segment, got %q", result)
	}
}
