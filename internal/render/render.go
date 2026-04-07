package render

import (
	"strings"
	"sync"
	"time"

	"github.com/jheddings/ccglow/internal/command"
	"github.com/jheddings/ccglow/internal/eval"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
	"github.com/rs/zerolog/log"
)

func isEnabled(node *types.SegmentNode, session *types.SessionData) bool {
	if node.EnabledFn != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Warn().Interface("panic", r).Msg("enabledFn panicked")
			}
		}()
		return node.EnabledFn(session)
	}
	if node.Enabled != nil {
		return *node.Enabled
	}
	return true
}

func renderNode(
	node *types.SegmentNode,
	session *types.SessionData,
	env map[string]any,
	defaultFormats map[string]string,
) *string {
	if !isEnabled(node, session) {
		return nil
	}

	// Composite: evaluate when, then render children
	if len(node.Children) > 0 {
		if node.When != "" {
			c := eval.CompileCached(node.When)
			if c == nil {
				return nil // compilation failed
			}
			segEnv := eval.BuildSegmentEnv(env, nil, "")
			if !c.Evaluate(segEnv) {
				return nil
			}
		}

		var parts []string
		for i := range node.Children {
			rendered := renderNode(&node.Children[i], session, env, defaultFormats)
			if rendered != nil {
				parts = append(parts, *rendered)
			}
		}
		if len(parts) == 0 {
			return nil
		}
		joined := strings.Join(parts, "")
		styled := style.Apply(joined, node.Style)
		return &styled
	}

	// Resolve raw value
	var raw any
	var hasValue bool

	if node.Value != nil {
		raw = node.Value
		hasValue = true
	} else if node.Expr != "" {
		result, err := eval.Eval(node.Expr, env)
		if err != nil {
			log.Warn().Err(err).Str("expr", node.Expr).Msg("expr eval failed")
			return nil
		}
		raw = result
		hasValue = true
	} else if node.Command != "" {
		output := command.Run(node.Command, env, session.CWD, command.DefaultTimeout)
		if output == "" {
			return nil
		}
		raw = output
		hasValue = true
	}

	if !hasValue {
		return nil
	}

	// Resolve format: config override > provider default > none
	format := node.Format
	if format == "" && node.Expr != "" {
		format = defaultFormats[node.Expr]
	}

	text := FormatValue(raw, format)
	if text == "" {
		return nil
	}

	// Evaluate when expression
	if node.When != "" {
		c := eval.CompileCached(node.When)
		if c == nil {
			return nil // compilation failed
		}
		segEnv := eval.BuildSegmentEnv(env, raw, text)
		if !c.Evaluate(segEnv) {
			return nil
		}
	}

	styled := style.Apply(text, node.Style)
	return &styled
}

// Options carries optional render-time configuration. Width, when > 0,
// overrides terminal width detection. WidthOffset is subtracted from the
// resolved width to leave room for host chrome.
type Options struct {
	Width       int
	WidthOffset int
}

// chunk is a piece of rendered output. text chunks are resolved strings;
// flex chunks are placeholders resolved at line-finalize time.
type chunk struct {
	text string
	flex bool
	fill string
}

// Tree renders the segment tree with default options. See TreeWith for the
// option-aware variant.
func Tree(
	tree []types.SegmentNode,
	session *types.SessionData,
	env map[string]any,
	defaultFormats map[string]string,
) string {
	return TreeWith(tree, session, env, defaultFormats, Options{})
}

// TreeWith performs a depth-first traversal of the segment tree, resolving
// each node against the environment and default formats. Flex segments at
// the top level are resolved per-line using the width derived from opts.
func TreeWith(
	tree []types.SegmentNode,
	session *types.SessionData,
	env map[string]any,
	defaultFormats map[string]string,
	opts Options,
) string {
	var chunks []chunk
	for i := range tree {
		node := &tree[i]
		if node.Flex {
			if !isFlexEnabled(node, session, env) {
				continue
			}
			fill := node.Fill
			if fill == "" {
				fill = " "
			}
			chunks = append(chunks, chunk{flex: true, fill: fill})
			continue
		}
		rendered := renderNode(node, session, env, defaultFormats)
		if rendered != nil {
			chunks = append(chunks, chunk{text: *rendered})
		}
	}
	return finalize(chunks, opts)
}

// isFlexEnabled checks Enabled and When for a top-level flex node.
func isFlexEnabled(node *types.SegmentNode, session *types.SessionData, env map[string]any) bool {
	if !isEnabled(node, session) {
		return false
	}
	if node.When != "" {
		c := eval.CompileCached(node.When)
		if c == nil {
			return false
		}
		if !c.Evaluate(eval.BuildSegmentEnv(env, nil, "")) {
			return false
		}
	}
	return true
}

// resolveWidth picks the effective render width: explicit opts.Width if
// set, otherwise TerminalWidth(); then subtracts opts.WidthOffset.
func resolveWidth(opts Options) int {
	w := opts.Width
	if w <= 0 {
		w = TerminalWidth()
	}
	if opts.WidthOffset > 0 && w > opts.WidthOffset {
		w -= opts.WidthOffset
	}
	return w
}

// finalize walks the chunk slice line-by-line, resolving any flex chunks by
// distributing the remaining terminal width across them.
func finalize(chunks []chunk, opts Options) string {
	if len(chunks) == 0 {
		return ""
	}
	hasFlex := false
	for _, c := range chunks {
		if c.flex {
			hasFlex = true
			break
		}
	}
	if !hasFlex {
		var b strings.Builder
		for _, c := range chunks {
			b.WriteString(c.text)
		}
		return b.String()
	}

	width := resolveWidth(opts)
	var out strings.Builder
	var line []chunk
	flush := func() {
		out.WriteString(resolveLine(line, width))
		line = line[:0]
	}
	for _, c := range chunks {
		if c.flex {
			line = append(line, c)
			continue
		}
		// split text on newlines so each segment lives on exactly one line
		parts := strings.Split(c.text, "\n")
		for i, p := range parts {
			if p != "" {
				line = append(line, chunk{text: p})
			}
			if i < len(parts)-1 {
				flush()
				out.WriteByte('\n')
			}
		}
	}
	flush()
	return out.String()
}

// resolveLine substitutes flex chunks in a single line with fill strings
// sized to fill the remaining terminal width. Even split across multiple
// flex chunks; collapses to zero on overflow.
func resolveLine(line []chunk, width int) string {
	if len(line) == 0 {
		return ""
	}
	used := 0
	flexCount := 0
	for _, c := range line {
		if c.flex {
			flexCount++
			continue
		}
		used += style.VisibleWidth(c.text)
	}
	remaining := max(width-used, 0)

	var b strings.Builder
	flexIdx := 0
	for _, c := range line {
		if !c.flex {
			b.WriteString(c.text)
			continue
		}
		// Distribute remaining across flex chunks; earlier chunks get the
		// extra cell when the split is uneven.
		share := remaining / flexCount
		if flexIdx < remaining%flexCount {
			share++
		}
		flexIdx++
		if share > 0 {
			b.WriteString(strings.Repeat(c.fill, share))
		}
	}
	return b.String()
}

// BuildEnv resolves all providers concurrently and merges their nested
// results into a single environment map. Returns the merged env and
// flat default format map.
func BuildEnv(
	providers map[string]types.DataProvider,
	session *types.SessionData,
) (env map[string]any, defaultFormats map[string]string) {
	env = make(map[string]any)
	defaultFormats = make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, p := range providers {
		wg.Add(1)
		go func(prov types.DataProvider) {
			defer wg.Done()
			start := time.Now()
			result, err := prov.Resolve(session)
			elapsed := time.Since(start)
			if err != nil {
				log.Warn().Err(err).Str("provider", prov.Name()).Msg("provider resolve failed")
				return
			}
			mu.Lock()
			for k, v := range result.Values {
				// inject __metrics__ into the provider's value subtree
				if m, ok := v.(map[string]any); ok {
					m["__metrics__"] = map[string]any{
						"duration_ms": elapsed.Seconds() * 1000,
					}
				}
				env[k] = v
			}
			for k, v := range result.Formats {
				defaultFormats[k] = v
			}
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return env, defaultFormats
}
