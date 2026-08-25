package provider

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func worktreeValues(result *types.ProviderResult) map[string]any {
	return result.Values["worktree"].(map[string]any)
}

func workspaceValues(result *types.ProviderResult) map[string]any {
	return result.Values["workspace"].(map[string]any)
}

func TestWorktreeProvider(t *testing.T) {
	p := &worktreeProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		Worktree: &types.WorktreeInfo{
			Name:           "my-feature",
			Path:           "/path/to/.claude/worktrees/my-feature",
			Branch:         "worktree-my-feature",
			OriginalCWD:    "/path/to/project",
			OriginalBranch: "main",
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	wt := worktreeValues(result)
	if wt["name"] != "my-feature" {
		t.Errorf("expected name my-feature, got %v", wt["name"])
	}
	if wt["path"] != "/path/to/.claude/worktrees/my-feature" {
		t.Errorf("unexpected path %v", wt["path"])
	}
	if wt["branch"] != "worktree-my-feature" {
		t.Errorf("expected branch worktree-my-feature, got %v", wt["branch"])
	}
	if wt["original_cwd"] != "/path/to/project" {
		t.Errorf("expected original_cwd /path/to/project, got %v", wt["original_cwd"])
	}
	if wt["original_branch"] != "main" {
		t.Errorf("expected original_branch main, got %v", wt["original_branch"])
	}
	if wt["active"] != true {
		t.Errorf("expected active true, got %v", wt["active"])
	}
}

// worktree is absent outside --worktree sessions.
func TestWorktreeProvider_NoData(t *testing.T) {
	p := &worktreeProvider{}
	result, err := p.Resolve(&types.SessionData{CWD: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}

	wt := worktreeValues(result)
	for _, key := range []string{"name", "path", "branch", "original_cwd", "original_branch"} {
		if wt[key] != "" {
			t.Errorf("expected empty %s, got %v", key, wt[key])
		}
	}
	if wt["active"] != false {
		t.Errorf("expected active false, got %v", wt["active"])
	}
}

// branch and original_branch are absent for hook-based worktrees even when
// the worktree object itself is present, so they resolve independently.
func TestWorktreeProvider_HookBasedWorktree(t *testing.T) {
	p := &worktreeProvider{}
	sess := &types.SessionData{
		CWD:      "/tmp",
		Worktree: &types.WorktreeInfo{Name: "hooked", Path: "/tmp/hooked"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	wt := worktreeValues(result)
	if wt["name"] != "hooked" {
		t.Errorf("expected name hooked, got %v", wt["name"])
	}
	if wt["branch"] != "" {
		t.Errorf("expected empty branch, got %v", wt["branch"])
	}
	if wt["active"] != true {
		t.Errorf("expected active true, got %v", wt["active"])
	}
}

// workspace.git_worktree marks any linked worktree, including ones outside a
// --worktree session, so it counts as active on its own.
func TestWorktreeProvider_ActiveFromWorkspace(t *testing.T) {
	p := &worktreeProvider{}
	sess := &types.SessionData{
		CWD:       "/tmp",
		Workspace: &types.WorkspaceInfo{GitWorktree: "feature-xyz"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	wt := worktreeValues(result)
	if wt["active"] != true {
		t.Errorf("expected active true, got %v", wt["active"])
	}
	if wt["name"] != "feature-xyz" {
		t.Errorf("expected name feature-xyz, got %v", wt["name"])
	}
}

// A --worktree name takes precedence over the linked-worktree name.
func TestWorktreeProvider_NamePrecedence(t *testing.T) {
	p := &worktreeProvider{}
	sess := &types.SessionData{
		CWD:       "/tmp",
		Worktree:  &types.WorktreeInfo{Name: "session-worktree"},
		Workspace: &types.WorkspaceInfo{GitWorktree: "linked-worktree"},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	if got := worktreeValues(result)["name"]; got != "session-worktree" {
		t.Errorf("expected session-worktree, got %v", got)
	}
}

func TestWorkspaceProvider(t *testing.T) {
	p := &workspaceProvider{}
	sess := &types.SessionData{
		CWD: "/tmp",
		Workspace: &types.WorkspaceInfo{
			CurrentDir: "/tmp",
			ProjectDir: "/original/project",
			AddedDirs:  []string{"/extra/one", "/extra/two"},
		},
	}

	result, err := p.Resolve(sess)
	if err != nil {
		t.Fatal(err)
	}

	ws := workspaceValues(result)
	if ws["project_dir"] != "/original/project" {
		t.Errorf("expected /original/project, got %v", ws["project_dir"])
	}
	if ws["added_dirs"] != 2 {
		t.Errorf("expected added_dirs 2, got %v", ws["added_dirs"])
	}
}

func TestWorkspaceProvider_NoData(t *testing.T) {
	p := &workspaceProvider{}
	result, err := p.Resolve(&types.SessionData{CWD: "/tmp"})
	if err != nil {
		t.Fatal(err)
	}

	ws := workspaceValues(result)
	if ws["project_dir"] != "" {
		t.Errorf("expected empty project_dir, got %v", ws["project_dir"])
	}
	if ws["added_dirs"] != 0 {
		t.Errorf("expected added_dirs 0, got %v", ws["added_dirs"])
	}
}

func TestWorktreeProviders_Registered(t *testing.T) {
	registry := NewRegistry()
	RegisterBuiltin(registry)
	for _, name := range []string{"worktree", "workspace"} {
		if _, ok := registry.All()[name]; !ok {
			t.Errorf("%s provider not registered", name)
		}
	}
}
