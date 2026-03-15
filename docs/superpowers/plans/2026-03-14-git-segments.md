# Git Status & Worktree Segments Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `git.modified`, `git.staged`, `git.untracked`, and `git.worktree` segments to ccglow by extending the existing git provider.

**Architecture:** Extend `GitData` with four new pointer fields. Add `git status --porcelain` parsing and worktree detection to the existing git provider's `Resolve()`. Add four new segment types that render the new fields. Update F1 preset to include `git.modified` and `git.untracked`.

**Tech Stack:** Go, `os/exec` for git subprocess calls, standard library testing

**Spec:** `docs/superpowers/specs/2026-03-14-git-segments-design.md`

---

## Chunk 1: Provider — git status parsing

### Task 1: Add new fields to GitData

**Files:**
- Modify: `internal/provider/git.go:17-21`

- [ ] **Step 1: Add fields to GitData struct**

```go
type GitData struct {
	Branch     *string
	Insertions *int
	Deletions  *int
	Modified   *int
	Staged     *int
	Untracked  *int
	Worktree   *string
}
```

- [ ] **Step 2: Run existing tests to verify nothing breaks**

Run: `go test ./internal/provider/ -v`
Expected: PASS (no existing tests depend on GitData field count)

- [ ] **Step 3: Commit**

```bash
git add internal/provider/git.go
git commit -m "feat(git): add status and worktree fields to GitData"
```

### Task 2: Implement git status parsing

**Files:**
- Modify: `internal/provider/git.go`
- Create: `internal/provider/git_test.go`

- [ ] **Step 1: Write failing test for status counts**

Create `internal/provider/git_test.go`:

```go
package provider

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func skipWithoutGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
}

// initTempRepo creates a temp dir with git init, an initial commit,
// and returns the path. Caller should defer os.RemoveAll.
func initTempRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v failed: %s", args, out)
		}
	}

	// Create and commit a seed file so HEAD exists
	seed := filepath.Join(dir, "seed.txt")
	if err := os.WriteFile(seed, []byte("seed"), 0644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"git", "add", "seed.txt"},
		{"git", "commit", "-m", "initial"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v failed: %s", args, out)
		}
	}

	return dir
}

func TestGitStatusCounts(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	// Modified: edit an existing tracked file (unstaged)
	os.WriteFile(filepath.Join(dir, "seed.txt"), []byte("changed"), 0644)

	// Staged: create and stage a new file
	os.WriteFile(filepath.Join(dir, "staged.txt"), []byte("new"), 0644)
	cmd := exec.Command("git", "add", "staged.txt")
	cmd.Dir = dir
	cmd.Run()

	// Untracked: create a file without adding
	os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("extra"), 0644)

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*GitData)
	if data.Modified == nil || *data.Modified != 1 {
		t.Errorf("expected 1 modified, got %v", data.Modified)
	}
	if data.Staged == nil || *data.Staged != 1 {
		t.Errorf("expected 1 staged, got %v", data.Staged)
	}
	if data.Untracked == nil || *data.Untracked != 1 {
		t.Errorf("expected 1 untracked, got %v", data.Untracked)
	}
}

func TestGitStatusClean(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*GitData)
	if data.Modified == nil {
		t.Fatal("expected non-nil Modified for clean repo")
	}
	if *data.Modified != 0 {
		t.Errorf("expected 0 modified, got %d", *data.Modified)
	}
	if data.Staged == nil || *data.Staged != 0 {
		t.Errorf("expected 0 staged, got %v", data.Staged)
	}
	if data.Untracked == nil || *data.Untracked != 0 {
		t.Errorf("expected 0 untracked, got %v", data.Untracked)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider/ -run TestGitStatus -v`
Expected: FAIL (Modified/Staged/Untracked are nil)

- [ ] **Step 3: Implement parseGitStatus and wire into Resolve**

Add to `internal/provider/git.go`:

```go
func parseGitStatus(cwd string) (modified, staged, untracked int, err error) {
	out, err := gitExec(cwd, "status", "--porcelain")
	if err != nil {
		return 0, 0, 0, err
	}
	if out == "" {
		return 0, 0, 0, nil
	}
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 2 {
			continue
		}
		if strings.HasPrefix(line, "??") {
			untracked++
			continue
		}
		x, y := line[0], line[1]
		// Column 1: staged changes
		if x == 'M' || x == 'A' || x == 'D' || x == 'R' || x == 'C' {
			staged++
		}
		// Column 2: unstaged changes
		if y == 'M' || y == 'D' || y == 'T' {
			modified++
		}
	}
	return modified, staged, untracked, nil
}

func intPtr(n int) *int { return &n }
```

Update `Resolve()` to call `parseGitStatus` after the existing diff logic:

```go
	if mod, stg, unt, err := parseGitStatus(cwd); err == nil {
		data.Modified = intPtr(mod)
		data.Staged = intPtr(stg)
		data.Untracked = intPtr(unt)
	}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/provider/ -run TestGitStatus -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/provider/git.go internal/provider/git_test.go
git commit -m "feat(git): add status parsing for modified/staged/untracked counts"
```

### Task 3: Implement worktree detection

**Files:**
- Modify: `internal/provider/git.go`
- Modify: `internal/provider/git_test.go`

- [ ] **Step 1: Write failing tests for worktree detection**

Add to `internal/provider/git_test.go`:

```go
func TestGitWorktreeDetection(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	// Create a linked worktree
	wtDir := filepath.Join(t.TempDir(), "my-worktree")
	cmd := exec.Command("git", "worktree", "add", wtDir, "-b", "test-branch")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("worktree add failed: %s", out)
	}

	p := &gitProvider{}
	sess := &types.SessionData{CWD: wtDir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*GitData)
	if data.Worktree == nil {
		t.Fatal("expected non-nil Worktree in linked worktree")
	}
	if *data.Worktree != "my-worktree" {
		t.Errorf("expected worktree name 'my-worktree', got %q", *data.Worktree)
	}
}

func TestGitWorktreeMainCopy(t *testing.T) {
	skipWithoutGit(t)
	dir := initTempRepo(t)

	p := &gitProvider{}
	sess := &types.SessionData{CWD: dir}
	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	data := result.(*GitData)
	if data.Worktree != nil {
		t.Errorf("expected nil Worktree in main copy, got %q", *data.Worktree)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/provider/ -run TestGitWorktree -v`
Expected: FAIL (Worktree is always nil)

- [ ] **Step 3: Implement worktree detection in Resolve**

Add to `internal/provider/git.go`:

```go
import "path/filepath"
```

Add helper function:

```go
func detectWorktree(cwd string) *string {
	gitDir, err := gitExec(cwd, "rev-parse", "--git-dir")
	if err != nil {
		return nil
	}
	commonDir, err := gitExec(cwd, "rev-parse", "--git-common-dir")
	if err != nil {
		return nil
	}
	// Normalize to absolute paths for comparison
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(cwd, gitDir)
	}
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(cwd, commonDir)
	}
	gitDir = filepath.Clean(gitDir)
	commonDir = filepath.Clean(commonDir)

	if gitDir == commonDir {
		return nil // main working copy
	}

	// In a linked worktree — get the worktree root name
	toplevel, err := gitExec(cwd, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil
	}
	name := filepath.Base(toplevel)
	return &name
}
```

Add to `Resolve()` after the status parsing:

```go
	data.Worktree = detectWorktree(cwd)
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/provider/ -run TestGitWorktree -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/provider/git.go internal/provider/git_test.go
git commit -m "feat(git): add linked worktree detection"
```

## Chunk 2: Segments and preset

### Task 4: Add new git segments

**Files:**
- Modify: `internal/segment/segment.go`
- Modify: `internal/segment/segment_test.go`

- [ ] **Step 1: Write failing tests for new segments**

Add to `internal/segment/segment_test.go`:

```go
import (
	"github.com/jheddings/ccglow/internal/provider"
)

func intPtr(n int) *int { return &n }
func strPtr(s string) *string { return &s }

func TestGitModifiedSegment(t *testing.T) {
	seg := &gitModifiedSegment{}
	if seg.Name() != "git.modified" {
		t.Errorf("expected name git.modified, got %s", seg.Name())
	}

	// With data
	ctx := &types.SegmentContext{Provider: &provider.GitData{Modified: intPtr(3)}}
	result := seg.Render(ctx)
	if result == nil || *result != "3" {
		t.Errorf("expected '3', got %v", result)
	}

	// Nil field
	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestGitStagedSegment(t *testing.T) {
	seg := &gitStagedSegment{}
	if seg.Name() != "git.staged" {
		t.Errorf("expected name git.staged, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.GitData{Staged: intPtr(5)}}
	result := seg.Render(ctx)
	if result == nil || *result != "5" {
		t.Errorf("expected '5', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestGitUntrackedSegment(t *testing.T) {
	seg := &gitUntrackedSegment{}
	if seg.Name() != "git.untracked" {
		t.Errorf("expected name git.untracked, got %s", seg.Name())
	}

	ctx := &types.SegmentContext{Provider: &provider.GitData{Untracked: intPtr(2)}}
	result := seg.Render(ctx)
	if result == nil || *result != "2" {
		t.Errorf("expected '2', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func TestGitWorktreeSegment(t *testing.T) {
	seg := &gitWorktreeSegment{}
	if seg.Name() != "git.worktree" {
		t.Errorf("expected name git.worktree, got %s", seg.Name())
	}

	name := "my-worktree"
	ctx := &types.SegmentContext{Provider: &provider.GitData{Worktree: &name}}
	result := seg.Render(ctx)
	if result == nil || *result != "my-worktree" {
		t.Errorf("expected 'my-worktree', got %v", result)
	}

	ctx = &types.SegmentContext{Provider: &provider.GitData{}}
	result = seg.Render(ctx)
	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/segment/ -run "TestGit(Modified|Staged|Untracked|Worktree)" -v`
Expected: FAIL (types not defined)

- [ ] **Step 3: Implement the four segment types**

Add to `internal/segment/segment.go` after the existing git segment definitions:

```go
type gitModifiedSegment struct{}

func (s *gitModifiedSegment) Name() string { return "git.modified" }
func (s *gitModifiedSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil && data.Modified != nil {
		v := fmt.Sprintf("%d", *data.Modified)
		return &v
	}
	return nil
}

type gitStagedSegment struct{}

func (s *gitStagedSegment) Name() string { return "git.staged" }
func (s *gitStagedSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil && data.Staged != nil {
		v := fmt.Sprintf("%d", *data.Staged)
		return &v
	}
	return nil
}

type gitUntrackedSegment struct{}

func (s *gitUntrackedSegment) Name() string { return "git.untracked" }
func (s *gitUntrackedSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil && data.Untracked != nil {
		v := fmt.Sprintf("%d", *data.Untracked)
		return &v
	}
	return nil
}

type gitWorktreeSegment struct{}

func (s *gitWorktreeSegment) Name() string { return "git.worktree" }
func (s *gitWorktreeSegment) Render(ctx *types.SegmentContext) *string {
	if data, ok := ctx.Provider.(*provider.GitData); ok && data != nil {
		return data.Worktree
	}
	return nil
}
```

Register in `RegisterBuiltin()`:

```go
	registry.Register(&gitModifiedSegment{})
	registry.Register(&gitStagedSegment{})
	registry.Register(&gitUntrackedSegment{})
	registry.Register(&gitWorktreeSegment{})
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/segment/ -run "TestGit(Modified|Staged|Untracked|Worktree)" -v`
Expected: PASS

- [ ] **Step 5: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/segment/segment.go internal/segment/segment_test.go
git commit -m "feat(segment): add git.modified, git.staged, git.untracked, git.worktree segments"
```

### Task 5: Update F1 preset

**Files:**
- Modify: `internal/preset/f1.json`

- [ ] **Step 1: Add git.modified and git.untracked to F1 preset**

Insert after the `git.deletions` entry (line 26) and before the closing powerline arrow for the git block (line 28):

```json
    {
      "segment": "git.modified",
      "style": { "color": "#FFD700", "bgcolor": "#3A3A3A", "prefix": " ~" }
    },
    {
      "segment": "git.untracked",
      "style": { "color": "#8BE9FD", "bgcolor": "#3A3A3A", "prefix": " ?" }
    },
```

- [ ] **Step 2: Run preset tests**

Run: `go test ./internal/preset/ -v`
Expected: PASS

- [ ] **Step 3: Run full test suite and build**

Run: `go vet ./... && go test ./... && go build ./...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/preset/f1.json
git commit -m "feat(preset): add git status segments to F1 preset"
```

### Task 6: Final verification

- [ ] **Step 1: Run full build and tests**

Run: `go vet ./... && go test ./... && go build ./...`
Expected: all PASS, binary builds

- [ ] **Step 2: Smoke test with sample input**

Run:
```bash
echo '{"cwd":"/tmp"}' | go run . --preset f1
```
Expected: renders without error, git segments collapse (not in a git repo)

Run:
```bash
echo '{"cwd":"'$(pwd)'"}' | go run . --preset f1
```
Expected: renders with git data visible (assuming ccglow repo has changes)
