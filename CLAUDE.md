# ccnow - AI Development Context

Composable, spaceship-style statusline for Claude Code. TypeScript + Node.js + chalk.

## Architecture

Segment tree model: atomic segments (single data point) and composite segments
(groups with `enabled` gating). DataProviders lazily fetch and cache external data.
A runner orchestrates the pipeline: parse CLI → load render tree → read stdin →
resolve providers → depth-first render → styled output.

**Key concepts**:

- `SegmentNode` — configuration (what to render, how it looks)
- `Segment` — runtime behavior (how to produce a value)
- `DataProvider` — lazy, cached data fetcher (git, pwd, context)
- DSL — internal authoring format for presets (factory functions with trailing closures)

## Project Structure

- `src/types.ts` — all shared interfaces
- `src/cli.ts` — entry point, stdin/stdout
- `src/runner.ts` — pipeline orchestrator
- `src/render.ts` — tree traversal, styling
- `src/segments/` — one file per segment
- `src/providers/` — one file per data provider
- `src/dsl/` — DSL factory functions
- `src/presets/` — named layouts (default, minimal, full)

## Development

**Key commands**:

- `just setup` — install dependencies
- `just build` — compile TypeScript
- `just test` — run all tests
- `just typecheck` — type check without emitting
- `just preflight` — build + typecheck + test (run before PR)
- `just dev` — build and run with sample input
- `just dev-live` — build and run with tee'd session data

**Adding a segment**: Create `src/segments/<name>.ts` implementing `Segment`,
register it in `src/segments/index.ts`.

**Adding a provider**: Create `src/providers/<name>.ts` implementing `DataProvider`,
register it in `src/providers/index.ts`.

## Commit Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <description>
```

Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `style`, `perf`

Scope is optional but encouraged (e.g. `fix(git): ...`, `feat(cli): ...`).

## Branch Naming

Use the same type prefixes as commits, followed by a short description:

```
<type>/<short-description>
```

Examples: `feat/color-themes`, `fix/token-formatting`, `chore/update-deps`

## Guardrails

**Do**:

- Follow TDD — write failing tests first, then implement
- Keep segments focused — one data point per atomic segment
- Return `null` from segments when there's no data (fail silent)
- Run `just preflight` before merging
- Work on feature branches

**Don't**:

- Push directly to main
- Never force-push to main
- Skip tests for new segments or providers
- Put styling logic in segments — segments return raw values, the runner applies style
- Mutate global chalk state — use `setColorLevel()` from `style.ts`
