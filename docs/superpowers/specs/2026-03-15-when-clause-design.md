# Conditional `when` Expression Syntax

Adds a `when` clause to `SegmentNode` that controls segment visibility
based on provider data or rendered values, using `expr-lang/expr` for
expression evaluation.

Resolves #1.

## Syntax

The `when` field is a string expression on `SegmentNode`, evaluated by
`expr-lang/expr` against the segment's provider data.

```json
{
  "segment": "context.percent",
  "when": ".percent >= 50",
  "style": { "color": "red" }
}
```

### Field references

- **`.field`** — resolves against the segment's provider data struct.
  Registered as variables in the expr environment using dot-prefixed
  lowercase names derived from the provider's exported Go field names.
  Examples: `.percent`, `.insertions`, `.branch`.
- **`value`** — the segment's raw data value from the provider (before
  rendering). The type depends on the provider field. No dot prefix —
  bare word, segment keyword. Evaluated after provider resolution.
- **`text`** — the segment's rendered text output (the unstyled `*string`
  returned by `Segment.Render()`), available as a string. Bare word,
  segment keyword. Evaluated after render.

### Nil pointer coercion

Provider data structs use pointer types for optional fields (`*int`,
`*string`). To keep expressions simple, `BuildEnv` coerces nil
pointers to typed zero values:

- `*int` nil → `0`
- `*string` nil → `""`
- `*float64` nil → `0.0`

Non-nil pointers are dereferenced to their underlying value.

This means `.percent >= 50` works without nil guards — a nil
`*int` becomes `0`, so the comparison is `0 >= 50` → false, and the
segment hides. This matches the intuition that "no data = don't show."

### Expression language

Full `expr-lang/expr` syntax is available:

- **Comparison:** `>=`, `<=`, `>`, `<`, `==`, `!=`
- **Boolean:** `&&`, `||`, `!`
- **Arithmetic:** `+`, `-`, `*`, `/` (e.g., `.insertions + .deletions > 10`)
- **String ops:** `contains`, `startsWith`, `endsWith`, `matches`

Examples:
- `".percent >= 50"` — show when context is >= 50%
- `".insertions > 0 || .deletions > 0"` — show when there are changes
- `".branch != 'main'"` — show when not on main (nil branch → `""`, so `"" != "main"` → true; but the git composite would already collapse when not in a repo)
- `"text != ''"` — show when segment renders something
- `"value > 0"` — show when the raw data value is positive

### Composite nodes

`when` is valid on composite nodes (nodes with `Children`). The
provider data for the composite is resolved the same way as for atomic
segments. Dot-field conditions work normally.

The `value` and `text` keywords are not meaningful on composite nodes
since their output is the concatenation of children. If referenced in
a composite's `when`, `value` evaluates as `nil` and `text` as `""`.

Example — hide the entire git group when not in a repo:
```json
{
  "when": ".branch != ''",
  "children": [
    { "segment": "git.branch", "style": { "prefix": " " } },
    { "segment": "git.insertions", "style": { "prefix": " +" } }
  ]
}
```

### Behavior

- `when` is empty → segment always renders (current behavior)
- Expression evaluates to true → segment renders normally
- Expression evaluates to false → segment returns nil (collapses)
- Expression compilation fails → treat as false, log to stderr
- Provider data is nil → all dot-field variables get zero values;
  `value` still works

## SegmentNode Changes

Add to `SegmentNode` in `internal/types/types.go`:

```go
When string `json:"when,omitempty"`
```

## New Package: `internal/condition/`

Create `internal/condition/condition.go` with:

### `Compile`

```go
func Compile(expression string) (*Condition, error)
```

Compiles the expression string into a reusable `Condition`. Returns an
error if the expression is syntactically invalid. Empty string returns
a nil `*Condition` (always true).

### `Condition.Evaluate`

```go
func (c *Condition) Evaluate(providerData any, renderedValue *string) bool
```

Builds an environment from `providerData` using `BuildEnv`, runs the
compiled expression, returns true if the result is boolean `true`.
Any other result type or runtime error returns false.

### `BuildEnv`

```go
func BuildEnv(providerData any, renderedValue *string) map[string]any
```

Builds the variable environment from provider data using reflection:

1. If `providerData` is a pointer, dereference to the struct.
2. For each exported field, register a dot-prefixed lowercase variable
   using the Go field name (e.g., `Branch` → `.branch`).
3. Dereference pointer fields with nil coercion:
   - `*int` nil → `0`, non-nil → `int` value
   - `*string` nil → `""`, non-nil → `string` value
   - `*float64` nil → `0.0`, non-nil → `float64` value
4. Non-pointer fields registered as-is.
5. Register `value`: for atomic segments, this is the raw provider
   data value that the segment renders (the specific field depends
   on the segment type). For now, `value` is set to `nil` — segments
   do not currently expose which provider field they map to. This
   can be enhanced later by having segments declare their source
   field. Dot-field references (`.percent`, `.branch`) are the
   primary way to access provider data.
6. Register `text` from `renderedValue` (`""` if nil).

If `providerData` is nil or not a struct/pointer-to-struct, return
a map containing only `value` (nil) and `text`.

Note: field names use the Go exported name lowercased, not JSON tags.
This keeps the mapping simple and predictable. Provider field names
are already short and clear (e.g., `Branch`, `Percent`, `Tokens`).

## Render Pipeline Changes

Modify `internal/render/render.go` to evaluate `when` during tree
traversal.

All `when` expressions are evaluated after render. The rendered value
is passed to `Evaluate` so `value` conditions work. If the expression
evaluates to false, the rendered output is discarded (returns nil).

For composite nodes, `when` is evaluated before recursing into
children. If false, the entire subtree is skipped. `value` is `""`
for composites.

### Interaction with `Enabled` / `EnabledFn`

`when` is an additional gate. Evaluation order:
1. `Enabled` / `EnabledFn` → if false, skip (existing behavior)
2. For atomic: `Segment.Render()` → get rendered value
3. `when` → if false, discard rendered value
4. For composite: evaluate `when` before recursing children

### Compilation caching

Compile `when` expressions into a `map[string]*Condition` keyed by
the expression string, built once during render tree setup. This
avoids adding runtime state to `SegmentNode` (which stays a pure
config struct) and deduplicates identical expressions across nodes.

### Error handling

Compilation errors are logged to stderr with `log.Printf` and the
segment is treated as hidden. Runtime evaluation errors also return
false (segment hidden).

## Dependency

Add `github.com/expr-lang/expr` to `go.mod`.

## Testing

### Condition package tests (`internal/condition/condition_test.go`)

**Compile:**
- Valid expression compiles without error
- Invalid expression returns error
- Empty expression returns nil Condition

**BuildEnv:**
- Provider struct with string, int, `*int`, `*string` fields →
  correct dot-prefixed variables
- Nil pointer fields → coerced zero values (0, "")
- Non-nil pointer fields → dereferenced values
- Nil provider → map with only `value`
- Non-struct provider → map with only `value`
- renderedValue non-nil → `value` set
- renderedValue nil → `value` is `""`

**Evaluate:**
- Numeric comparisons: `>=`, `>`, `<`, `<=`, `==`, `!=`
- String comparisons: `==`, `!=`
- Boolean combinators: `&&`, `||`
- Nil-coerced field: `.field >= 50` where field was nil `*int`
- `value` against rendered output
- Empty expression (nil Condition) → true
- Expression returning non-bool → false
- Runtime error → false

### Render integration tests (`internal/render/render_test.go`)

- Segment with `when` that passes → renders normally
- Segment with `when` that fails → returns nil
- Segment with `value` condition
- Composite with `when` that passes → children render
- Composite with `when` that fails → subtree skipped
- Segment with no `when` → unchanged behavior
- Segment with invalid `when` → treated as hidden

## What Doesn't Change

- Existing `Enabled` / `EnabledFn` behavior — `when` is additive
- Provider resolution, segment rendering, styling — all unchanged
- All existing presets and configs — `when` is optional
