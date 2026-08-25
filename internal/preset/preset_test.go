package preset

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jheddings/ccglow/internal/provider"
	"github.com/jheddings/ccglow/internal/render"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
)

// collectStyles walks a segment tree and returns every style attached to it.
func collectStyles(nodes []types.SegmentNode) []*types.StyleAttrs {
	var styles []*types.StyleAttrs
	for i := range nodes {
		node := &nodes[i]
		if node.Style != nil {
			styles = append(styles, node.Style)
		}
		styles = append(styles, collectStyles(node.Children)...)
	}
	return styles
}

// colorRenders reports whether a color string produces any escape code.
// resolveColor silently yields nothing for an unrecognized value, so an
// invalid name is indistinguishable from no color at all -- which is exactly
// what makes a typo like "brightBlue" (the table has "blueBright") invisible.
func colorRenders(color string, background bool) bool {
	var withColor, without *types.StyleAttrs
	if background {
		withColor = &types.StyleAttrs{Background: color}
	} else {
		withColor = &types.StyleAttrs{Color: color}
	}
	without = &types.StyleAttrs{}
	return style.Apply("x", withColor) != style.Apply("x", without)
}

// Every color named in a shipped preset must resolve. An unresolvable color
// renders as unstyled text rather than erroring, so nothing else catches it.
func TestPresetColorsResolve(t *testing.T) {
	style.SetColorLevel(1)

	for _, name := range List() {
		nodes := Get(name)
		if nodes == nil {
			t.Errorf("%s: failed to load", name)
			continue
		}
		for _, st := range collectStyles(nodes) {
			if st.Color != "" && !colorRenders(st.Color, false) {
				t.Errorf("%s: color %q does not resolve", name, st.Color)
			}
			if st.Background != "" && !colorRenders(st.Background, true) {
				t.Errorf("%s: bgcolor %q does not resolve", name, st.Background)
			}
		}
	}
}

// collectExprs walks a segment tree and returns every expr string in it.
func collectExprs(nodes []types.SegmentNode) []string {
	var exprs []string
	for i := range nodes {
		node := &nodes[i]
		if node.Expr != "" {
			exprs = append(exprs, node.Expr)
		}
		exprs = append(exprs, collectExprs(node.Children)...)
	}
	return exprs
}

// lookup walks a dotted key path through the resolved provider env.
func lookup(env map[string]any, key string) bool {
	parts := strings.Split(key, ".")
	current := env
	for i, part := range parts {
		val, ok := current[part]
		if !ok {
			return false
		}
		if i == len(parts)-1 {
			return true
		}
		nested, ok := val.(map[string]any)
		if !ok {
			return false
		}
		current = nested
	}
	return false
}

// Segments fail silent by design: an expr that resolves to nothing collapses
// out of the output rather than erroring. That makes a typo in a preset
// invisible forever, so every expr shipped in a preset is checked against the
// real resolved provider env. Providers populate their full key set with zero
// values even when there is no data, so this passes without live session data.
func TestPresetExprsResolve(t *testing.T) {
	registry := provider.NewRegistry()
	provider.RegisterBuiltin(registry)
	env, _ := render.BuildEnv(registry.All(), &types.SessionData{CWD: t.TempDir()})

	names := List()
	if len(names) == 0 {
		t.Fatal("expected at least one preset")
	}

	for _, name := range names {
		nodes := Get(name)
		if nodes == nil {
			t.Errorf("%s: failed to load", name)
			continue
		}
		exprs := collectExprs(nodes)
		if len(exprs) == 0 {
			t.Errorf("%s: no expr nodes found", name)
			continue
		}
		for _, expr := range exprs {
			if !lookup(env, expr) {
				t.Errorf("%s: expr %q does not resolve against any provider", name, expr)
			}
		}
	}
}

func TestGet_Default(t *testing.T) {
	nodes := Get("default")
	if nodes == nil {
		t.Fatal("expected default preset, got nil")
	}
	if len(nodes) == 0 {
		t.Fatal("expected non-empty segment tree")
	}
}

func TestGet_Minimal(t *testing.T) {
	nodes := Get("minimal")
	if nodes == nil {
		t.Fatal("expected minimal preset, got nil")
	}
}

func TestGet_Full(t *testing.T) {
	nodes := Get("full")
	if nodes == nil {
		t.Fatal("expected full preset, got nil")
	}
}

func TestGet_Glow(t *testing.T) {
	nodes := Get("glow")
	if nodes == nil {
		t.Fatal("expected glow preset, got nil")
	}
	if len(nodes) == 0 {
		t.Fatal("expected non-empty segment tree")
	}
}

func TestGet_Unknown(t *testing.T) {
	nodes := Get("nonexistent")
	if nodes != nil {
		t.Errorf("expected nil for unknown preset, got %v", nodes)
	}
}

func TestGet_MinimalHasSegments(t *testing.T) {
	nodes := Get("minimal")
	if nodes == nil {
		t.Fatal("expected minimal preset")
	}
	if len(nodes) == 0 {
		t.Fatal("expected non-empty segment tree")
	}
}

func TestList(t *testing.T) {
	names := List()
	if len(names) < 3 {
		t.Fatalf("expected at least 3 presets, got %d: %v", len(names), names)
	}

	expected := map[string]bool{"default": false, "minimal": false, "full": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected preset %q in list", name)
		}
	}
}

func TestDump_ReturnsValidJSON(t *testing.T) {
	data, err := Dump("default")
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(data) {
		t.Error("expected valid JSON from Dump")
	}
}

func TestDump_Unknown(t *testing.T) {
	_, err := Dump("nonexistent")
	if err == nil {
		t.Error("expected error for unknown preset")
	}
}
