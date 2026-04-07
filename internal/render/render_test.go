package render

import (
	"testing"

	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
)

func TestTree_Empty(t *testing.T) {
	sess := &types.SessionData{CWD: "/tmp"}
	result := Tree(nil, sess, map[string]any{}, map[string]string{})
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestTree_ExprNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/home/user/project"}
	env := map[string]any{
		"pwd": map[string]any{"name": "project"},
	}

	tree := []types.SegmentNode{
		{Expr: "pwd.name"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "project" {
		t.Errorf("expected project, got %q", result)
	}
}

func TestTree_ValueNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Value: "hello"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "hello" {
		t.Errorf("expected hello, got %q", result)
	}
}

func TestTree_ValueNewline(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Value: "\n"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "\n" {
		t.Errorf("expected newline, got %q", result)
	}
}

func TestTree_CompositeCollapse(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"git": map[string]any{
			"branch":     "",
			"insertions": "",
		},
	}

	tree := []types.SegmentNode{
		{
			Style: &types.StyleAttrs{Prefix: " | "},
			Children: []types.SegmentNode{
				{Expr: "git.branch"},
				{Expr: "git.insertions"},
			},
		},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "" {
		t.Errorf("expected empty (collapsed composite), got %q", result)
	}
}

func TestTree_DisabledNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"pwd": map[string]any{"name": "tmp"},
	}

	disabled := false
	tree := []types.SegmentNode{
		{Expr: "pwd.name", Enabled: &disabled},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "" {
		t.Errorf("expected empty for disabled node, got %q", result)
	}
}

func TestTree_MissingExpr(t *testing.T) {
	sess := &types.SessionData{CWD: "/tmp"}

	tree := []types.SegmentNode{
		{Expr: "nonexistent.segment"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "" {
		t.Errorf("expected empty for missing segment, got %q", result)
	}
}

func TestTree_ExprNode_Value(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"name": "hello"},
	}

	tree := []types.SegmentNode{{Expr: "test.name"}}
	result := Tree(tree, sess, env, map[string]string{})
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_ExprEmptyCollapses(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"name": ""},
	}

	tree := []types.SegmentNode{{Expr: "test.name"}}
	result := Tree(tree, sess, env, map[string]string{})
	if result != "" {
		t.Errorf("expected empty (collapsed), got %q", result)
	}
}

func TestTree_ExprWithFormat(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"count": 42},
	}

	tree := []types.SegmentNode{{Expr: "test.count", Format: "+%d"}}
	result := Tree(tree, sess, env, map[string]string{})
	if result != "+42" {
		t.Errorf("expected '+42', got %q", result)
	}
}

func TestTree_ExprDefaultFormat(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"fmt": map[string]any{"pct": 85},
	}
	defaultFormats := map[string]string{"fmt.pct": "%d%%"}

	tree := []types.SegmentNode{{Expr: "fmt.pct"}}
	result := Tree(tree, sess, env, defaultFormats)
	if result != "85%" {
		t.Errorf("expected '85%%', got %q", result)
	}
}

func TestTree_ExprFormatOverridesDefault(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"fmt": map[string]any{"pct": 85},
	}
	defaultFormats := map[string]string{"fmt.pct": "%d%%"}

	tree := []types.SegmentNode{{Expr: "fmt.pct", Format: "(%d)"}}
	result := Tree(tree, sess, env, defaultFormats)
	if result != "(85)" {
		t.Errorf("expected '(85)', got %q", result)
	}
}

func TestTree_EmptyStringCollapses(t *testing.T) {
	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"plain": ""},
	}

	tree := []types.SegmentNode{{Expr: "test.plain"}}
	result := Tree(tree, sess, env, map[string]string{})
	if result != "" {
		t.Errorf("expected empty (collapsed), got %q", result)
	}
}

func TestTree_ExprWhenPasses(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"count": 75},
	}

	tree := []types.SegmentNode{
		{Expr: "test.count", When: "value >= 50"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "75" {
		t.Errorf("expected '75', got %q", result)
	}
}

func TestTree_ExprWhenFails(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"count": 25},
	}

	tree := []types.SegmentNode{
		{Expr: "test.count", When: "value >= 50"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "" {
		t.Errorf("expected empty (when failed), got %q", result)
	}
}

func TestTree_ExprWhenCrossProvider(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{
			"name":  "feature",
			"count": 5,
		},
	}

	tree := []types.SegmentNode{
		{Expr: "test.name", When: "test.count > 0"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "feature" {
		t.Errorf("expected 'feature', got %q", result)
	}
}

func TestTree_ExprWhenText(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"name": "hello"},
	}

	tree := []types.SegmentNode{
		{Expr: "test.name", When: "text != ''"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_CompositeWhen(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"git": map[string]any{"branch": "main"},
	}

	tree := []types.SegmentNode{
		{
			When: "git.branch != ''",
			Children: []types.SegmentNode{
				{Expr: "git.branch"},
			},
		},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "main" {
		t.Errorf("expected 'main', got %q", result)
	}
}

func TestTree_CompositeWhenFails(t *testing.T) {
	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"git": map[string]any{"branch": ""},
	}

	tree := []types.SegmentNode{
		{
			When: "git.branch != ''",
			Children: []types.SegmentNode{
				{Value: "should not appear"},
			},
		},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "" {
		t.Errorf("expected empty (composite when failed), got %q", result)
	}
}

func TestTree_WhenNoExpression(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"name": "hello"},
	}

	tree := []types.SegmentNode{
		{Expr: "test.name"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_WhenInvalidExpression(t *testing.T) {
	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"name": "hello"},
	}

	tree := []types.SegmentNode{
		{Expr: "test.name", When: ">>>bad<<<"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "" {
		t.Errorf("expected empty (invalid when), got %q", result)
	}
}

func TestBuildEnv(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}
	sess := &types.SessionData{CWD: "/tmp"}

	env, formats := BuildEnv(providers, sess)

	test, ok := env["test"].(map[string]any)
	if !ok {
		t.Fatal("expected test namespace in env")
	}
	if test["name"] != "hello" {
		t.Errorf("expected test.name='hello', got %v", test["name"])
	}
	if formats["test.pct"] != "%d%%" {
		t.Errorf("expected test.pct format, got %q", formats["test.pct"])
	}
}

func TestBuildEnv_Metrics(t *testing.T) {
	providers := map[string]types.DataProvider{
		"test": &testProvider{},
	}
	sess := &types.SessionData{CWD: "/tmp"}

	env, _ := BuildEnv(providers, sess)

	test, ok := env["test"].(map[string]any)
	if !ok {
		t.Fatal("expected test namespace in env")
	}

	metrics, ok := test["__metrics__"].(map[string]any)
	if !ok {
		t.Fatal("expected __metrics__ in test namespace")
	}

	duration, ok := metrics["duration_ms"]
	if !ok {
		t.Fatal("expected duration_ms in __metrics__")
	}

	dur, ok := duration.(float64)
	if !ok {
		t.Fatalf("expected duration_ms to be float64, got %T", duration)
	}
	if dur < 0 {
		t.Errorf("expected non-negative duration, got %f", dur)
	}
}

// --- Command node tests ---

func TestTree_CommandNode(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Command: "echo hello"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_CommandEmptyCollapses(t *testing.T) {
	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Command: "printf ''"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "" {
		t.Errorf("expected empty (collapsed), got %q", result)
	}
}

func TestTree_CommandNonZeroCollapses(t *testing.T) {
	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Command: "exit 1"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "" {
		t.Errorf("expected empty for non-zero exit, got %q", result)
	}
}

func TestTree_CommandWithFormat(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Command: "echo 42", Format: "count: %s"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "count: 42" {
		t.Errorf("expected 'count: 42', got %q", result)
	}
}

func TestTree_CommandWhenPasses(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Command: "echo hello", When: "text != ''"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTree_CommandWhenFails(t *testing.T) {
	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Command: "echo 0", When: "text != '0'"},
	}

	result := Tree(tree, sess, map[string]any{}, map[string]string{})
	if result != "" {
		t.Errorf("expected empty (when failed), got %q", result)
	}
}

func TestTree_CommandWithInterpolation(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"name": "world"},
	}
	tree := []types.SegmentNode{
		{Command: "echo ${test.name}"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "world" {
		t.Errorf("expected 'world', got %q", result)
	}
}

func TestTree_ExprTakesPrecedenceOverCommand(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)

	sess := &types.SessionData{CWD: "/tmp"}
	env := map[string]any{
		"test": map[string]any{"name": "from-expr"},
	}
	tree := []types.SegmentNode{
		{Expr: "test.name", Command: "echo from-command"},
	}

	result := Tree(tree, sess, env, map[string]string{})
	if result != "from-expr" {
		t.Errorf("expected 'from-expr' (expr precedence), got %q", result)
	}
}

func TestTree_FlexFillsRemainingWidth(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)
	t.Setenv("COLUMNS", "20")

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Value: "left"},  // 4 chars
		{Flex: true},     // fills remainder
		{Value: "right"}, // 5 chars
	}

	got := Tree(tree, sess, map[string]any{}, map[string]string{})
	want := "left           right" // 4 + 11 spaces + 5 = 20
	if got != want {
		t.Errorf("flex render = %q (len %d), want %q (len %d)", got, len(got), want, len(want))
	}
}

func TestTree_FlexCustomFill(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)
	t.Setenv("COLUMNS", "10")

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Value: "a"},
		{Flex: true, Fill: "-"},
		{Value: "b"},
	}

	got := Tree(tree, sess, map[string]any{}, map[string]string{})
	want := "a--------b"
	if got != want {
		t.Errorf("flex with fill = %q, want %q", got, want)
	}
}

func TestTree_FlexCollapsesOnOverflow(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)
	t.Setenv("COLUMNS", "5")

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Value: "longer"}, // already 6 > 5
		{Flex: true},
		{Value: "x"},
	}

	got := Tree(tree, sess, map[string]any{}, map[string]string{})
	want := "longerx"
	if got != want {
		t.Errorf("overflow collapse = %q, want %q", got, want)
	}
}

func TestTree_FlexEvenSplit(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)
	t.Setenv("COLUMNS", "12")

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Value: "a"},
		{Flex: true},
		{Value: "b"},
		{Flex: true},
		{Value: "c"},
	}

	got := Tree(tree, sess, map[string]any{}, map[string]string{})
	// 12 - 3 = 9, split across 2 flex => 4 and 5 (or 5 and 4)
	if len(got) != 12 {
		t.Errorf("even split len = %d, want 12 (got %q)", len(got), got)
	}
	if got[0] != 'a' || got[len(got)-1] != 'c' {
		t.Errorf("even split bookends wrong: %q", got)
	}
}

func TestTree_FlexPerLine(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)
	t.Setenv("COLUMNS", "10")

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Value: "a"},
		{Flex: true},
		{Value: "b"},
		{Value: "\n"},
		{Value: "c"},
		{Flex: true},
		{Value: "d"},
	}

	got := Tree(tree, sess, map[string]any{}, map[string]string{})
	want := "a        b\nc        d"
	if got != want {
		t.Errorf("per-line flex = %q, want %q", got, want)
	}
}

func TestTree_FlexNoOpWithoutFlex(t *testing.T) {
	style.SetColorLevel(0)
	defer style.SetColorLevel(1)
	t.Setenv("COLUMNS", "20")

	sess := &types.SessionData{CWD: "/tmp"}
	tree := []types.SegmentNode{
		{Value: "hello"},
	}

	got := Tree(tree, sess, map[string]any{}, map[string]string{})
	if got != "hello" {
		t.Errorf("no-flex line should not be padded, got %q", got)
	}
}

// testProvider implements DataProvider for tests.
type testProvider struct{}

func (p *testProvider) Name() string { return "test" }
func (p *testProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	return &types.ProviderResult{
		Values: map[string]any{
			"test": map[string]any{
				"name":  "hello",
				"count": 42,
				"pct":   85,
			},
		},
		Formats: map[string]string{
			"test.pct": "%d%%",
		},
	}, nil
}
