package render

import (
	"testing"

	"github.com/jheddings/ccglow/internal/provider"
	"github.com/jheddings/ccglow/internal/segment"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
)

func setupTestRegistries() *segment.Registry {
	reg := segment.NewRegistry()
	segment.RegisterBuiltin(reg)
	return reg
}

func TestTree_Empty(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}
	result := Tree(nil, seg, sess, map[string]any{})
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestTree_AtomicNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/home/user/project"}
	providerData := map[string]any{
		"pwd": &provider.PwdData{Name: "project", Path: "/home/user/", Smart: "~/"},
	}

	tree := []types.SegmentNode{
		{Type: "pwd.name", Provider: "pwd"},
	}

	result := Tree(tree, seg, sess, providerData)
	if result != "project" {
		t.Errorf("expected project, got %q", result)
	}
}

func TestTree_CompositeCollapse(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}
	providerData := map[string]any{
		"git": &provider.GitData{},
	}

	tree := []types.SegmentNode{
		{
			Type:  "group",
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Type: "git.branch", Provider: "git"},
				{Type: "git.insertions", Provider: "git"},
			},
		},
	}

	result := Tree(tree, seg, sess, providerData)
	if result != "" {
		t.Errorf("expected empty (collapsed composite), got %q", result)
	}
}

func TestTree_DisabledNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}
	providerData := map[string]any{
		"pwd": &provider.PwdData{Name: "tmp"},
	}

	disabled := false
	tree := []types.SegmentNode{
		{Type: "pwd.name", Provider: "pwd", Enabled: &disabled},
	}

	result := Tree(tree, seg, sess, providerData)
	if result != "" {
		t.Errorf("expected empty for disabled node, got %q", result)
	}
}

func TestTree_Literal(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	tree := []types.SegmentNode{
		{Type: "literal", Props: map[string]any{"text": "hello"}},
	}

	result := Tree(tree, seg, sess, map[string]any{})
	if result != "hello" {
		t.Errorf("expected hello, got %q", result)
	}
}

func TestTree_MissingSegment(t *testing.T) {
	seg := setupTestRegistries()
	sess := &types.SessionData{CWD: "/tmp"}

	tree := []types.SegmentNode{
		{Type: "nonexistent.segment"},
	}

	result := Tree(tree, seg, sess, map[string]any{})
	if result != "" {
		t.Errorf("expected empty for missing segment, got %q", result)
	}
}
