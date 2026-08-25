package provider

import "github.com/jheddings/ccglow/internal/types"

type worktreeProvider struct{}

func (p *worktreeProvider) Name() string { return "worktree" }

func (p *worktreeProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	worktree := map[string]any{
		"name":            "",
		"path":            "",
		"branch":          "",
		"original_cwd":    "",
		"original_branch": "",
		"active":          false,
	}

	result := &types.ProviderResult{
		Values: map[string]any{"worktree": worktree},
	}

	// A linked worktree counts as active even outside a --worktree session,
	// where only workspace.git_worktree is populated.
	if session.Workspace != nil && session.Workspace.GitWorktree != "" {
		worktree["name"] = session.Workspace.GitWorktree
		worktree["active"] = true
	}

	wt := session.Worktree
	if wt == nil {
		return result, nil
	}

	// Each field is independently absent for hook-based worktrees, so they
	// are assigned one at a time rather than as a block.
	if wt.Name != "" {
		worktree["name"] = wt.Name
	}
	if wt.Path != "" {
		worktree["path"] = wt.Path
	}
	if wt.Branch != "" {
		worktree["branch"] = wt.Branch
	}
	if wt.OriginalCWD != "" {
		worktree["original_cwd"] = wt.OriginalCWD
	}
	if wt.OriginalBranch != "" {
		worktree["original_branch"] = wt.OriginalBranch
	}
	worktree["active"] = true

	return result, nil
}

type workspaceProvider struct{}

func (p *workspaceProvider) Name() string { return "workspace" }

func (p *workspaceProvider) Resolve(session *types.SessionData) (*types.ProviderResult, error) {
	workspace := map[string]any{
		"project_dir": "",
		"added_dirs":  0,
	}

	result := &types.ProviderResult{
		Values: map[string]any{"workspace": workspace},
	}

	if session.Workspace == nil {
		return result, nil
	}

	if session.Workspace.ProjectDir != "" {
		workspace["project_dir"] = session.Workspace.ProjectDir
	}
	workspace["added_dirs"] = len(session.Workspace.AddedDirs)

	return result, nil
}
