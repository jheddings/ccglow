package condition

import (
	"sync"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// Condition is a compiled when expression.
type Condition struct {
	program *vm.Program
}

var (
	evalCache = make(map[string]*vm.Program)
	evalMu    sync.Mutex
)

// Compile compiles an expression string into a reusable Condition.
// Returns nil for empty expressions (always true).
func Compile(expression string) (*Condition, error) {
	if expression == "" {
		return nil, nil
	}

	program, err := expr.Compile(expression, expr.AllowUndefinedVariables())
	if err != nil {
		return nil, err
	}

	return &Condition{program: program}, nil
}

// Evaluate runs the compiled expression against the environment.
// Returns true only if the result is boolean true.
// Nil receiver (empty expression) returns true.
func (c *Condition) Evaluate(env map[string]any) bool {
	if c == nil {
		return true
	}

	result, err := expr.Run(c.program, env)
	if err != nil {
		return false
	}

	b, ok := result.(bool)
	return ok && b
}

// Eval compiles (with caching) and runs an expression against the environment,
// returning the result as any value. Used to resolve expr fields in nodes.
func Eval(expression string, env map[string]any) (any, error) {
	if expression == "" {
		return nil, nil
	}

	evalMu.Lock()
	program, ok := evalCache[expression]
	if !ok {
		var err error
		program, err = expr.Compile(expression, expr.AllowUndefinedVariables())
		if err != nil {
			evalMu.Unlock()
			return nil, err
		}
		evalCache[expression] = program
	}
	evalMu.Unlock()

	return expr.Run(program, env)
}

// BuildSegmentEnv creates the evaluation environment for a single segment's
// when expression. It shallow-copies the nested env and adds value/text keys.
func BuildSegmentEnv(nestedEnv map[string]any, value any, text string) map[string]any {
	env := make(map[string]any, len(nestedEnv)+2)
	for k, v := range nestedEnv {
		env[k] = v
	}
	env["value"] = value
	env["text"] = text
	return env
}
