# ccnow — Composable Statusline for Claude Code

**Date**: 2026-03-14
**Status**: Draft

## Overview

`ccnow` is a TypeScript CLI package that renders a composable, spaceship-style statusline for Claude Code. It reads session JSON from stdin and outputs a styled, single-line string to stdout.

Design goals:

- **Zero install**: `npx -y ccnow` with sensible defaults
- **CLI-first**: composite flags define layout order (`npx -y ccnow --pwd --sep --git --sep --context`)
- **Configurable**: JSON config file for customization, DSL for presets and power users
- **Extensible**: clean segment interface makes adding new segments trivial
- **Fast**: low-millisecond execution, lazy provider loading, fail-silent

## Architecture

### Render Tree

The fundamental data structure is a tree of segment nodes. Every configuration path — CLI flags, JSON config, DSL presets — produces the same render tree.

```
StatusLine
├── Pwd.smart          (provider: pwd)
├── Sep                (no provider)
├── Git                (composite, enabled: gitAvailable)
│   ├── Git.branch     (provider: git)
│   ├── Literal " ["
│   ├── Git.insertions (provider: git)
│   ├── Literal " "
│   ├── Git.deletions  (provider: git)
│   └── Literal "]"
├── Sep
└── Context            (composite, enabled: true)
    ├── Literal "ctx: "
    ├── Context.tokens  (provider: context)
    ├── Literal " ("
    ├── Context.percent (provider: context)
    └── Literal ")"
```

### Segment Node

```ts
interface SegmentNode {
  type: string;                              // e.g. 'git.branch', 'sep', 'literal'
  provider?: string;                         // e.g. 'git' — omit if session data suffices
  enabled?: boolean | EnabledFn;             // conditional display, defaults to true
  style?: StyleAttrs;                        // color, bold, dim, icon, prefix, suffix
  props?: Record<string, unknown>;           // segment-specific (text for literals, char for sep)
  children?: SegmentNode[];                  // composites only
}

type EnabledFn = (session: SessionData, provider?: unknown) => boolean;
```

- **Atomic segments** return a raw value string from `render()`. The runner applies styling.
- **Composite segments** are groups. They evaluate `enabled` first (gate check), then render children depth-first, and collapse to null if all children are null.
- **Literal segments** output static text. No provider, no data.
- **Separator segments** output a styled character (pipe, arrow, space, etc.).

### Segment Interface

```ts
interface Segment {
  name: string;
  provider?: string;
  render(context: SegmentContext): string | null;
}

interface SegmentContext {
  session: SessionData;
  provider?: unknown;
  config?: SegmentConfig;
}
```

Segments return a raw value string or `null` (nothing to show, error, or `enabled` evaluated to false). The runner wraps the value with style attributes. This separates data logic from presentation.

### Style Attributes

```ts
interface StyleAttrs {
  color?: string;       // 'cyan', 'red', '#ff5500'
  bold?: boolean;
  dim?: boolean;
  italic?: boolean;
  icon?: string;        // prefix glyph, e.g. '\ue0a0'
  prefix?: string;      // text before value, e.g. '+'
  suffix?: string;      // text after value
}
```

Style is declarative and lives in the config/DSL, not in segment code. The runner applies it uniformly via chalk.

### DataProviders

```ts
interface DataProvider {
  name: string;
  resolve(session: SessionData): Promise<unknown>;
}
```

Providers fetch and cache data for segments. Called once per run, result shared across all segments that declare the provider. If `resolve` throws, all dependent segments render as null (fail silent).

**Built-in providers:**

| Provider | Source | Data |
|----------|--------|------|
| `session` | stdin JSON (always available) | cwd, context_window, etc. |
| `git` | git CLI commands against `session.cwd` | branch, insertions, deletions |
| `pwd` | derived from `session.cwd` | name, path, smart-truncated path |
| `context` | derived from session JSON | token count (formatted), percentage |

Session data is the base context passed to every segment automatically. Segments declare an additional provider only if they need data beyond the session JSON. Providers are resolved lazily — only activated if an enabled segment declares them. Async providers run concurrently via `Promise.all`.

## Data Flow

1. **Parse CLI** — resolve preset name, segment toggles, config file path, output format, tee path
2. **Load render tree** — CLI flags, JSON config, or preset DSL → all produce a `SegmentNode` tree
3. **Read stdin** — parse Claude Code session JSON
4. **Tee** (if `--tee` flag) — write raw stdin JSON to file before processing
5. **Resolve providers** — walk tree, collect unique provider names from enabled segments, resolve concurrently
6. **Render** — depth-first tree traversal:
   - Composite: evaluate `enabled` → if false, return null (skip children entirely). If true, render children, collapse if all null.
   - Atomic: call `segment.render(context)` → if null, skip. If string, apply style attrs.
7. **Output** — concatenate non-null results, write to stdout

## Configuration

### Three Tiers

**CLI flags** — what segments and in what order, using composite names for convenience:

```sh
# Default preset
npx -y ccnow

# Select preset
npx -y ccnow --preset=minimal

# Composite toggles (order defines layout)
npx -y ccnow --pwd --sep --git --sep --context

# Config file
npx -y ccnow --config ~/.claude/ccnow.json

# Output format
npx -y ccnow --format=plain

# Debug: save stdin JSON to file
npx -y ccnow --tee /tmp/session.json
```

**JSON config** — serialized render tree for customization without writing code:

```json
{
  "segments": [
    { "segment": "pwd.smart", "color": "cyan", "bold": true },
    { "segment": "sep", "char": "|", "dim": true },
    { "segment": "git", "children": [
      { "segment": "git.branch", "color": "white", "icon": "\ue0a0 " },
      { "segment": "literal", "text": " [" },
      { "segment": "git.insertions", "color": "green", "prefix": "+" },
      { "segment": "literal", "text": " " },
      { "segment": "git.deletions", "color": "red", "prefix": "-" },
      { "segment": "literal", "text": "]" }
    ]},
    { "segment": "sep", "char": "|", "dim": true },
    { "segment": "context", "children": [
      { "segment": "literal", "text": "ctx: " },
      { "segment": "context.tokens", "bold": true },
      { "segment": "literal", "text": " (" },
      { "segment": "context.percent" },
      { "segment": "literal", "text": ")" }
    ]}
  ]
}
```

**DSL** — internal authoring format for presets, future power-user configs:

```ts
import { StatusLine, Pwd, Sep, Git, Branch, Insertions, Deletions,
         Context, Tokens, Percent, Literal } from 'ccnow/dsl'

export default StatusLine(() => [
  Pwd({ style: 'smart', color: 'cyan', bold: true }),
  Sep({ char: '|', dim: true }),
  Git({ enabled: (session) => gitAvailable(session.cwd) })(() => [
    Branch({ color: 'white', icon: '\ue0a0 ' }),
    Literal({ text: ' [' }),
    Insertions({ color: 'green', prefix: '+' }),
    Literal({ text: ' ' }),
    Deletions({ color: 'red', prefix: '-' }),
    Literal({ text: ']' }),
  ]),
  Sep({ char: '|', dim: true }),
  Context()(() => [
    Literal({ text: 'ctx: ' }),
    Tokens({ bold: true }),
    Literal({ text: ' (' }),
    Percent(),
    Literal({ text: ')' }),
  ]),
])
```

The DSL and JSON both hydrate into the same render tree. Presets are named DSL files that ship with the package.

### Config Resolution Priority

CLI flags > config file > preset > built-in default

### Conditional Display

Segments and composites use `enabled` for conditional display:

- `true` / `false` — static enable/disable (JSON and DSL)
- Function `(session, provider?) => boolean` — dynamic condition (DSL only)

Built-in composites ship with sensible defaults (e.g. `Git` checks for git repo availability). Atomic segments return `null` when they have no data, which effectively hides them.

`when` is reserved for future use (potential JSON-expressible conditions).

## Built-in Segments

### Atomic Segments

| Segment | Provider | Description |
|---------|----------|-------------|
| `pwd.name` | pwd | Directory basename |
| `pwd.path` | pwd | Full path |
| `pwd.smart` | pwd | p10k-style truncated path |
| `git.branch` | git | Current branch name |
| `git.insertions` | git | Lines added (staged + unstaged vs HEAD) |
| `git.deletions` | git | Lines removed (staged + unstaged vs HEAD) |
| `context.tokens` | context | Token count, human-formatted (24K, 1.2M) |
| `context.percent` | context | Context window usage percentage |
| `literal` | — | Static text string |
| `sep` | — | Separator character (pipe, arrow, space, etc.) |

### Composite Segments (CLI Shorthand)

| Composite | Expands To | Default `enabled` |
|-----------|------------|-------------------|
| `--pwd` | `pwd.smart` | always |
| `--git` | `git.branch` `[` `git.insertions` `git.deletions` `]` | git repo available |
| `--context` | `ctx:` `context.tokens` `(` `context.percent` `)` | always |
| `--sep` | `sep` with default char | always |

## Presets

| Preset | Description |
|--------|-------------|
| `default` | `pwd.smart \| git \| context` — mirrors existing statusline |
| `minimal` | `pwd.name \| git.branch` |
| `full` | All segments enabled with verbose formatting |

## CLI Reference

```
Usage: ccnow [options]

Options:
  --preset <name>     Use a named preset (default, minimal, full)
  --config <path>     Load JSON config file
  --pwd               Enable pwd composite segment
  --git               Enable git composite segment
  --context           Enable context composite segment
  --sep               Insert separator segment
  --format <type>     Output format: ansi (default), plain, json
  --tee <path>        Write raw stdin JSON to file before processing
  --help              Show help
  --version           Show version
```

## Project Structure

```
ccnow/
  package.json          # Package manifest (app dependencies, bin, npm config)
  justfile              # Project lifecycle tasks (build, test, lint, etc.)
  tsconfig.json
  README.md
  src/
    cli.ts              # Entry point, arg parsing
    runner.ts           # Pipeline orchestrator
    render.ts           # Tree traversal, styling, output
    types.ts            # Shared interfaces
    providers/
      session.ts        # Parses stdin JSON (base context)
      git.ts            # Git CLI commands
      pwd.ts            # Working directory variants
      context.ts        # Token count, percentage formatting
    segments/
      git.branch.ts
      git.insertions.ts
      git.deletions.ts
      context.tokens.ts
      context.percent.ts
      pwd.name.ts
      pwd.path.ts
      pwd.smart.ts
      sep.ts
      literal.ts
    presets/
      default.ts        # DSL-authored default layout
      minimal.ts
      full.ts
    dsl/
      index.ts          # DSL factory functions and tree builder
  docs/
    superpowers/
      specs/
        2026-03-14-ccnow-design.md
```

## Error Handling

- **Fail silent**: any segment that can't produce data returns `null`
- **Provider failure**: if a provider's `resolve` throws, all dependent segments render as `null`
- **Composite collapse**: if all children of a composite are `null`, the composite is `null` (no orphaned separators or literals)
- **Stdin failure**: if stdin is empty or invalid JSON, fall back to a minimal output or exit cleanly
- **Performance**: target sub-50ms execution. Providers resolve concurrently. Git commands use `session.cwd` to avoid filesystem discovery.

## Dependencies

- **chalk** — terminal styling (color, bold, dim, etc.)
- **minimist** or **commander** — CLI argument parsing (lean toward minimist for size)

No other runtime dependencies. Keep the package small for `npx` cold-start performance.

## Future Considerations (Out of Scope)

- Third-party segment authorship (npm packages or local files)
- TS config file support for power users (DSL as user-facing config)
- `when` conditions in JSON (expression-based conditional display)
- Interactive config builder CLI
- Icon theme presets (powerline, nerd-font, ascii)
- `pwd.smart` truncation strategies (p10k-style)
