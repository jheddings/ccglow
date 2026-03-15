package statusline

import "fmt"

// RegisterBuiltinSegments adds all built-in segment implementations to the registry.
func RegisterBuiltinSegments(registry *SegmentRegistry) {
	registry.Register(&literalSegment{})
	registry.Register(&pwdNameSegment{})
	registry.Register(&pwdPathSegment{})
	registry.Register(&pwdSmartSegment{})
	registry.Register(&gitBranchSegment{})
	registry.Register(&gitInsertionsSegment{})
	registry.Register(&gitDeletionsSegment{})
	registry.Register(&contextTokensSegment{})
	registry.Register(&contextSizeSegment{})
	registry.Register(&contextPercentSegment{})
	registry.Register(&modelNameSegment{})
	registry.Register(&costUSDSegment{})
	registry.Register(&sessionDurationSegment{})
	registry.Register(&sessionLinesAddedSegment{})
	registry.Register(&sessionLinesRemovedSegment{})
}

// --- Literal ---

type literalSegment struct{}

func (s *literalSegment) Name() string { return "literal" }
func (s *literalSegment) Render(ctx *SegmentContext) *string {
	if ctx.Props == nil {
		return nil
	}
	if text, ok := ctx.Props["text"].(string); ok {
		return &text
	}
	return nil
}

// --- PWD ---

type pwdNameSegment struct{}

func (s *pwdNameSegment) Name() string { return "pwd.name" }
func (s *pwdNameSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*PwdData); ok && data != nil {
		return &data.Name
	}
	return nil
}

type pwdPathSegment struct{}

func (s *pwdPathSegment) Name() string { return "pwd.path" }
func (s *pwdPathSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*PwdData); ok && data != nil && data.Path != "" {
		return &data.Path
	}
	return nil
}

type pwdSmartSegment struct{}

func (s *pwdSmartSegment) Name() string { return "pwd.smart" }
func (s *pwdSmartSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*PwdData); ok && data != nil && data.Smart != "" {
		return &data.Smart
	}
	return nil
}

// --- Git ---

type gitBranchSegment struct{}

func (s *gitBranchSegment) Name() string { return "git.branch" }
func (s *gitBranchSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*GitData); ok && data != nil {
		return data.Branch
	}
	return nil
}

type gitInsertionsSegment struct{}

func (s *gitInsertionsSegment) Name() string { return "git.insertions" }
func (s *gitInsertionsSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*GitData); ok && data != nil && data.Insertions != nil {
		v := fmt.Sprintf("%d", *data.Insertions)
		return &v
	}
	return nil
}

type gitDeletionsSegment struct{}

func (s *gitDeletionsSegment) Name() string { return "git.deletions" }
func (s *gitDeletionsSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*GitData); ok && data != nil && data.Deletions != nil {
		v := fmt.Sprintf("%d", *data.Deletions)
		return &v
	}
	return nil
}

// --- Context ---

type contextTokensSegment struct{}

func (s *contextTokensSegment) Name() string { return "context.tokens" }
func (s *contextTokensSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*ContextData); ok && data != nil && data.Tokens != "" {
		return &data.Tokens
	}
	return nil
}

type contextSizeSegment struct{}

func (s *contextSizeSegment) Name() string { return "context.size" }
func (s *contextSizeSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*ContextData); ok && data != nil && data.Size != "" {
		return &data.Size
	}
	return nil
}

type contextPercentSegment struct{}

func (s *contextPercentSegment) Name() string { return "context.percent" }
func (s *contextPercentSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*ContextData); ok && data != nil && data.Percent != nil {
		v := fmt.Sprintf("%d%%", *data.Percent)
		return &v
	}
	return nil
}

// --- Model ---

type modelNameSegment struct{}

func (s *modelNameSegment) Name() string { return "model.name" }
func (s *modelNameSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*ModelData); ok && data != nil {
		return data.Name
	}
	return nil
}

// --- Cost ---

type costUSDSegment struct{}

func (s *costUSDSegment) Name() string { return "cost.usd" }
func (s *costUSDSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*CostData); ok && data != nil {
		return data.USD
	}
	return nil
}

// --- Session ---

type sessionDurationSegment struct{}

func (s *sessionDurationSegment) Name() string { return "session.duration" }
func (s *sessionDurationSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*SessionProviderData); ok && data != nil {
		return data.Duration
	}
	return nil
}

type sessionLinesAddedSegment struct{}

func (s *sessionLinesAddedSegment) Name() string { return "session.lines-added" }
func (s *sessionLinesAddedSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*SessionProviderData); ok && data != nil && data.LinesAdded != nil {
		v := fmt.Sprintf("%d", *data.LinesAdded)
		return &v
	}
	return nil
}

type sessionLinesRemovedSegment struct{}

func (s *sessionLinesRemovedSegment) Name() string { return "session.lines-removed" }
func (s *sessionLinesRemovedSegment) Render(ctx *SegmentContext) *string {
	if data, ok := ctx.Provider.(*SessionProviderData); ok && data != nil && data.LinesRemoved != nil {
		v := fmt.Sprintf("%d", *data.LinesRemoved)
		return &v
	}
	return nil
}
