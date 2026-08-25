# Segments

Every segment renders a single piece of data — a branch name, a token count, a
dollar amount. Compose them into any layout you want.

Segments are identified by `provider.field` expressions. The prefix determines
which provider fetches the data; the suffix picks the specific value. If a
segment has nothing to show (empty string, no data available), it silently
collapses out of the output.

## Directory — `pwd`

| Segment     | Description                                                                | Example Output |
| ----------- | -------------------------------------------------------------------------- | -------------- |
| `pwd.name`  | Directory basename                                                         | `ccglow`       |
| `pwd.path`  | Full path prefix (everything before the basename, with trailing slash)     | `~/Projects/`  |
| `pwd.smart` | Smart-truncated path — abbreviates intermediate directories for deep paths | `~/P/…/`       |

`pwd.smart` keeps the first and last path components readable and abbreviates
the middle when nesting gets deep. Pair it with `pwd.name` for a compact but
navigable path display.

## Git — `git`

| Segment          | Description                                       | Example Output |
| ---------------- | ------------------------------------------------- | -------------- |
| `git.branch`     | Current branch name                               | `main`         |
| `git.insertions` | Lines added (staged + unstaged combined)          | `42`           |
| `git.deletions`  | Lines removed (staged + unstaged combined)        | `17`           |
| `git.modified`   | Count of modified (unstaged) files                | `3`            |
| `git.staged`     | Count of staged files                             | `2`            |
| `git.untracked`  | Count of untracked files                          | `5`            |
| `git.owner`      | Repository owner                                  | `jheddings`    |
| `git.repo`       | Repository name                                   | `ccglow`       |
| `git.host`       | Repository host                                   | `github.com`   |
| `git.worktree`   | Linked worktree name (empty in main working copy) | `docs-update`  |

Most git segments require a git repository in the current working directory.
When not in a git repo, they return their zero values (`""` for strings, `0`
for integers).

`git.owner`, `git.repo`, `git.host`, and `git.worktree` come from Claude Code
directly when it supplies them, which avoids two subprocesses per render. They
fall back to parsing the `origin` remote (both SSH and HTTPS formats) and
inspecting the working copy when it doesn't — outside a git repo, with no
`origin` remote, or on older Claude Code versions. `git.host` is only
available from Claude Code, so it stays empty on the fallback path.

## Worktree — `worktree`

| Segment                    | Description                            | Example Output              |
| -------------------------- | -------------------------------------- | --------------------------- |
| `worktree.name`            | Active worktree name                   | `my-feature`                |
| `worktree.path`            | Worktree directory                     | `~/.claude/worktrees/my-feature` |
| `worktree.branch`          | Branch checked out in the worktree     | `worktree-my-feature`       |
| `worktree.original_cwd`    | Directory before entering the worktree | `~/Projects/ccglow`         |
| `worktree.original_branch` | Branch before entering the worktree    | `main`                      |
| `worktree.active`          | Whether a worktree is active (bool)    | `true`                      |

`worktree.name` and `worktree.active` are set for any linked worktree. The
remaining segments need a `--worktree` session, and `branch` /
`original_branch` are additionally empty for hook-based worktrees — so gate
them individually rather than assuming they arrive together.

Use `worktree.active` to make the whole group vanish in ordinary sessions:

```json
{
  "when": "worktree.active",
  "children": [
    { "expr": "worktree.name", "style": { "color": "magenta", "prefix": " ⑂ " } },
    {
      "expr": "worktree.original_branch",
      "when": "worktree.original_branch != ''",
      "style": { "color": "brightblack", "prefix": " ← " }
    }
  ]
}
```

## Workspace — `workspace`

| Segment                 | Description                                | Example Output      |
| ----------------------- | ------------------------------------------ | ------------------- |
| `workspace.project_dir` | Directory where Claude Code was launched   | `~/Projects/ccglow` |
| `workspace.added_dirs`  | Count of `/add-dir` directories (int)      | `2`                 |

`workspace.project_dir` can differ from `pwd.path` when the working directory
changes during a session. `workspace.added_dirs` is a count rather than a list,
since a list of paths doesn't fit a statusline — use it as a `when` condition
or a small badge:

```json
{ "expr": "workspace.added_dirs", "when": "workspace.added_dirs > 0", "format": "+%d dirs" }
```

## Context — `context`

| Segment                     | Description                          | Example Output |
| --------------------------- | ------------------------------------ | -------------- |
| `context.tokens`            | Total token count, human-formatted   | `360K`, `1.2M` |
| `context.size`              | Context window capacity              | `1M`, `200K`   |
| `context.percent.used`      | Usage as integer percentage          | `36%`          |
| `context.percent.remaining` | Remaining capacity as percentage     | `64%`          |
| `context.input`             | Total input tokens, human-formatted  | `162K`         |
| `context.output`            | Total output tokens, human-formatted | `45K`          |
| `context.bar`               | Visual progress bar (10 chars wide)  | `███░░░░░░░`   |

Token formatting scales automatically: raw count below 1K, `nK` for thousands,
`n.nM` for millions.

## Model — `model`

| Segment      | Description        | Example Output          |
| ------------ | ------------------ | ----------------------- |
| `model.name` | Model display name | `Opus 4.6 (1M context)` |
| `model.id`   | Model identifier   | `claude-opus-4-6`       |

## Effort — `effort`

| Segment           | Description                              | Example Output |
| ----------------- | ---------------------------------------- | -------------- |
| `effort.level`    | Reasoning effort level                   | `high`         |
| `effort.thinking` | Whether extended thinking is enabled     | `true`         |
| `effort.fast`     | Whether fast mode is enabled             | `true`         |

Where `model` describes model identity, `effort` describes how the model is
currently behaving. All three values track the live session, so mid-session
`/effort` changes appear on the next render.

`effort.level` is one of `low`, `medium`, `high`, `xhigh`, or `max`. It is
empty when the current model has no effort parameter, so gate it with
`"when": "effort.level != ''"` rather than assuming a default. Ultracode is
not a distinct level and reports as `xhigh`.

`effort.thinking` and `effort.fast` are booleans, which makes them most useful
as `when` conditions on an icon:

```json
{ "value": "🚀", "when": "effort.fast" }
```

## Cost — `cost`

| Segment      | Description                              | Example Output |
| ------------ | ---------------------------------------- | -------------- |
| `cost.usd`   | Session cost formatted USD               | `$12.50`       |
| `cost.total` | Raw numeric cost (float64, format $%.2f) | `$12.50`       |

`cost.total` is the raw numeric value suitable for `when` conditions
(e.g. `"when": "cost.total > 5"`).

## Rate Limits — `limits`

| Segment                   | Description                                | Example Output |
| ------------------------- | ------------------------------------------ | -------------- |
| `limits.session.percent`  | 5-hour window used (float, format `%.0f%%`) | `24%`          |
| `limits.session.reset`    | Time until the 5-hour window resets        | `2h 14m`       |
| `limits.session.reset_at` | Raw reset time, epoch seconds (int)        | `1738425600`   |
| `limits.weekly.percent`   | 7-day window used (float, format `%.0f%%`)  | `41%`          |
| `limits.weekly.reset`     | Time until the 7-day window resets         | `30h 0m`       |
| `limits.weekly.reset_at`  | Raw reset time, epoch seconds (int)        | `1738857600`   |

Claude Code sends these only for Claude.ai subscribers (Pro/Max), and only
after the first API response in the session. Each window is independently
absent, so a session may have `limits.session` populated and `limits.weekly`
empty. All segments return zero values otherwise (`0` percent, `""` reset).

Percentages run 0–100, so `when` conditions compare directly:

```json
{ "expr": "limits.session.percent", "when": "limits.session.percent > 80", "style": { "color": "red" } }
```

A bare percentage next to a `5h` label reads as two unrelated fields. The
shipped presets bind them by putting the label in its own muted node, so the
colour shift groups the pair:

```json
{ "value": "5h:", "when": "limits.session.percent > 0", "style": { "color": "240" } },
{ "expr": "limits.session.percent", "when": "limits.session.percent > 0", "style": { "color": "yellow" } }
```

The label has to be a separate node rather than a `prefix`, since a prefix is
painted in its own node's colour and would come out yellow along with the
value.

The `reset` segments render the time *remaining*, which means they go stale
between renders. Statusline updates are event-driven, so if you display a
countdown, set `refreshInterval` in your `statusLine` settings to re-run
ccglow on a timer. A window whose reset time has already elapsed renders empty
rather than a negative duration.

## Speed — `speed`

| Segment        | Description                        | Example Output       |
| -------------- | ---------------------------------- | -------------------- |
| `speed.input`  | Input token throughput             | `45 t/s`, `1.2K t/s` |
| `speed.output` | Output token throughput            | `82 t/s`             |
| `speed.total`  | Combined input + output throughput | `127 t/s`            |

Speed is calculated from total tokens divided by API duration. Formatting
scales the same way as tokens: raw below 1K, `n.nK t/s` above.

## Session — `session`

| Segment                      | Description                      | Example Output  |
| ---------------------------- | -------------------------------- | --------------- |
| `session.duration.total`     | Wall-clock session time          | `2h 15m`, `45m` |
| `session.duration.api`       | Time spent on API calls          | `8m`, `1h 2m`   |
| `session.duration.total_min` | Wall-clock time in minutes (int) | `135`           |
| `session.duration.api_min`   | API time in minutes (int)        | `8`             |
| `session.id`                 | Session identifier               | `abc-123`       |
| `session.name`               | Session name                     | `auth-refactor` |
| `session.lines-added`        | Total lines added this session   | `1380`          |
| `session.lines-removed`      | Total lines removed this session | `21`            |

The `_min` variants are raw integers suitable for `when` conditions
(e.g. `"when": "session.duration.total_min > 60"`).

`session.name` is the custom name set with `--name` or `/rename` when one
exists, otherwise the AI-generated session title. It is empty for a session
with neither — the default display name (e.g. `my-app-3f`) does not count as
a name.

## Agent — `agent`

| Segment        | Description                            | Example Output      |
| -------------- | -------------------------------------- | ------------------- |
| `agent.name`   | Agent driving the session              | `security-reviewer` |
| `agent.active` | Whether this is an agent session (bool) | `true`              |

Populated when Claude Code runs with `--agent` or with agent settings
configured, and empty otherwise. Use `agent.active` to gate a whole group so
ordinary sessions render unchanged:

```json
{
  "when": "agent.active",
  "children": [{ "expr": "agent.name", "style": { "color": "brightmagenta", "prefix": " 🤖 " } }]
}
```

## Claude — `claude`

| Segment          | Description                     | Example Output |
| ---------------- | ------------------------------- | -------------- |
| `claude.version` | Claude Code application version | `2.1.75`       |
| `claude.style`   | Current output style            | `concise`      |

## System — `system`

| Segment                  | Description                          | Example Output    |
| ------------------------ | ------------------------------------ | ----------------- |
| `system.load.avg1`       | 1-minute load average (format %.2f)  | `1.42`            |
| `system.load.avg5`       | 5-minute load average (format %.2f)  | `2.10`            |
| `system.load.avg15`      | 15-minute load average (format %.2f) | `1.87`            |
| `system.mem.used`        | Used memory, human-formatted         | `12.4G`           |
| `system.mem.total`       | Total memory, human-formatted        | `32G`             |
| `system.mem.percent`     | Memory usage percentage              | `39%`             |
| `system.disk.used`       | Used disk space, human-formatted     | `234G`            |
| `system.disk.total`      | Total disk space, human-formatted    | `1T`              |
| `system.disk.percent`    | Disk usage percentage                | `23%`             |
| `system.battery.percent` | Battery charge percentage            | `85%`             |
| `system.battery.state`   | Battery state                        | `charging`        |
| `system.uptime`          | System uptime, human-formatted       | `3d 14h`, `2h 5m` |

Disk usage is measured at the mount point of the current working directory.
Battery segments return zero values on machines without a battery.

## Node Types

There are four kinds of atomic nodes:

### `expr` — Expression nodes

Evaluate an expression against the provider data environment. This is the
primary way to display provider values.

```json
{ "expr": "git.branch", "style": { "bold": true } }
{ "expr": "context.percent.used", "style": { "prefix": " (" , "suffix": ")" } }
```

### `value` — Static value nodes

Render a fixed string. Use these for separators, icons, and line breaks.

```json
{ "value": "|", "style": { "color": "240" } }
{ "value": "\n" }
{ "value": "\ue0b0", "style": { "color": "#DC0000", "bgcolor": "#3A3A3A" } }
```

### `command` — Shell command nodes

Run a shell command and render its stdout. The command is executed via `sh -c`
with a 2-second timeout. If the command produces empty output, exits non-zero,
or times out, the node collapses.

```json
{ "command": "cat VERSION", "style": { "color": "cyan" } }
{ "command": "date +%H:%M" }
```

#### Variable substitution

Use `${provider.field}` to inject resolved provider values into the command
string. Variables are replaced before execution. Unresolved references become
empty strings.

```json
{
  "command": "gh pr list --repo ${git.owner}/${git.repo} --json number --jq length",
  "when": "text != '' && text != '0'"
}
```

#### Behavior

- **Timeout**: Commands have a 2-second timeout. Long-running commands are
  killed and the node collapses.
- **Working directory**: Commands run in the session's current working directory.
- **Collapse**: Empty stdout, non-zero exit, or timeout all cause the node to
  collapse silently, just like an `expr` node with no data.
- **Precedence**: If both `expr` and `command` are set on the same node, `expr`
  wins. The dispatch order is: Children > Value > Expr > Command.

### `flex` — Elastic spacer

Expands to fill the remaining terminal width on its line. Enables
right-alignment: put some segments to the left, a `flex`, then the segments
you want on the right.

```json
{ "segments": [{ "expr": "pwd.name" }, { "flex": true }, { "expr": "context.percent.used" }] }
```

- Top-level only. A `flex` nested inside a `children` group is ignored.
- Defaults to filling with spaces. Set `fill` to any single character (e.g.
  `"·"` or `"─"`) to use a different fill.
- Each line has its own flex resolution — multi-line layouts using `newline`
  segments can have one or more `flex` per line.
- Multiple flex segments on the same line split the remaining width evenly.
- If non-flex content already exceeds terminal width, flex collapses to zero.
- Terminal width is read from `$COLUMNS`, falling back to a `TIOCGWINSZ`
  ioctl on stdout, then to 80. Claude Code does not currently export
  `COLUMNS` to the statusline subprocess, so exporting it from your shell
  config (e.g. `export COLUMNS`) is the most reliable option.

## Node Properties

### `format`

Expression nodes accept an optional `format` string that controls how the raw
value is displayed. Uses Go's `fmt.Sprintf` syntax.

```json
{ "expr": "git.insertions", "format": "+%d" }
{ "expr": "context.percent.used", "format": "(%d%%)" }
```

If no format is specified, the node uses its default format (declared by
the provider) or falls back to the raw value as a string.

### `when`

Any node can conditionally show or hide based on data. See
**[Conditional Visibility](WHEN.md)** for the full reference.

```json
{ "expr": "git.branch", "when": "git.branch != '' && git.branch != 'main'" }
{ "expr": "context.percent.used", "when": "context.percent.used >= 50" }
{ "expr": "git.modified", "when": "value > 0" }
```

### `enabled`

Set `"enabled": false` on any node to exclude it from rendering. Disabled nodes
are skipped entirely, as if they weren't in the tree. Unlike `when`, this is a
static setting — it doesn't evaluate at runtime.

## Groups and Composites

Any node can have `children`. When it does, it acts as a composite — rendering
all children depth-first and collapsing entirely if every child produces empty
output. Use composites for sections that should appear or disappear together.

```json
{
  "when": "git.branch != ''",
  "style": { "prefix": " | " },
  "children": [
    { "expr": "git.branch", "style": { "bold": true } },
    {
      "expr": "git.insertions",
      "when": "value > 0",
      "style": { "color": "green", "prefix": " +" }
    },
    { "expr": "git.deletions", "when": "value > 0", "style": { "color": "red", "prefix": " -" } }
  ]
}
```

Composites support `when` expressions that can reference any provider's data,
allowing you to gate entire sections on any condition. See
[WHEN.md](WHEN.md) for details.
