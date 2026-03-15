# Git Status & Worktree Segments

Adds four new segments to the git provider: `git.modified`, `git.staged`,
`git.untracked`, and `git.worktree`.

Resolves #31 (git status indicator) and #32 (git worktree).

## Approach

Extend the existing `git` provider (Approach A) rather than creating separate
providers. This keeps the `git.*` segment family cohesive and avoids provider
proliferation. The render pipeline only resolves providers referenced by
configured segments, so unused providers have no cost. Note that when any
`git.*` segment is configured, all git subprocess calls run together — this
is acceptable overhead for a statusline that refreshes on each prompt.

## GitData Changes

Add four fields to `GitData` in `internal/provider/git.go`:

```go
type GitData struct {
    Branch     *string
    Insertions *int
    Deletions  *int
    Modified   *int    // modified (unstaged) file count
    Staged     *int    // staged file count
    Untracked  *int    // untracked file count
    Worktree   *string // linked worktree name, nil in main copy
}
```

All new fields use pointer types following the existing pattern. A `nil` value
means no data available, causing the segment to return `nil` and collapse in
composites.

## Provider Changes

Two new git subprocess calls in `Resolve()`:

### git status --porcelain

Parses the two-column status output:
- Column 1 (index/staged): count entries with `M`, `A`, `D`, `R`, or `C` -> `Staged`
- Column 2 (worktree/modified): count entries with `M`, `D`, or `T` -> `Modified`
- Lines starting with `??` -> `Untracked`

Merge conflict entries (`UU`, `AA`, `DD`) are intentionally not counted in any
category — they represent an in-progress merge state, not a clean
modified/staged distinction.

Renamed entries (`R`) may produce lines like `R  old -> new`; parsing uses
only the first two characters of each line, so the extended format is
irrelevant.

When the repo is clean (no output), all three counts are set to `intPtr(0)`
rather than left nil, so segments render `"0"` in a clean repo. This
distinguishes "inside a git repo with no changes" from "not in a git repo"
(where the fields remain nil and segments collapse).

On error, fields remain nil and segments collapse (consistent with existing
error handling in the provider).

Uses the existing `gitExec()` helper with the 5-second timeout.

### git worktree detection

Compares `git rev-parse --git-common-dir` with `git rev-parse --git-dir`:
- If they resolve to different paths, the session is in a linked worktree
- The worktree name is derived from `filepath.Base()` of the output of
  `git rev-parse --show-toplevel` (not `session.CWD`, which may be a
  subdirectory within the worktree)
- If they match (main working copy), `Worktree` remains nil

## New Segments

Four new segment types in `internal/segment/segment.go`:

| Segment | Provider field | Render |
|---------|---------------|--------|
| `git.modified` | `Modified *int` | Decimal string, nil when no data |
| `git.staged` | `Staged *int` | Decimal string, nil when no data |
| `git.untracked` | `Untracked *int` | Decimal string, nil when no data |
| `git.worktree` | `Worktree *string` | String as-is, nil when not in worktree |

All follow the existing segment pattern: type-assert `ctx.Provider` to
`*provider.GitData`, return `nil` when the field is nil.

Registered in `RegisterBuiltin()` alongside existing git segments.

## Preset Update

Update the F1 preset (`internal/preset/f1.json`) to include `git.modified`
and `git.untracked` in the git section of line 1. Add them after
`git.deletions` with styling consistent with the existing git block:

- `git.modified` — yellow text (`#FFD700`) with `~` prefix
- `git.untracked` — cyan text (`#8BE9FD`) with `?` prefix

Both share the `#3A3A3A` background of the git block.

## Testing

### Provider tests (`internal/provider/git_test.go`, new file)

Integration tests that shell out to `git`. Guarded with
`if _, err := exec.LookPath("git"); err != nil { t.Skip("git not available") }`
to avoid failures in CI environments without git.

- **TestGitStatusCounts**: create a temp git repo with `git init`, add and
  commit a file, then modify it (unstaged), stage another change, and create
  an untracked file. Verify `Modified`, `Staged`, and `Untracked` counts.
- **TestGitStatusClean**: clean repo returns zero counts (not nil).
- **TestGitWorktreeDetection**: create a temp repo, add a linked worktree via
  `git worktree add`, resolve from the worktree cwd, verify `Worktree` is
  set to the expected worktree directory name.
- **TestGitWorktreeMainCopy**: resolve from main working copy, verify
  `Worktree` is nil.

### Segment tests (`internal/segment/segment_test.go`, extend existing)

- Test each new segment with a populated `GitData` (verify rendered string).
- Test each new segment with nil fields (verify nil return).

## What Doesn't Change

- Existing segments (`git.branch`, `git.insertions`, `git.deletions`) unchanged
- Provider registry, render pipeline — no changes needed
