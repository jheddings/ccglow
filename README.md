# ccnow

A composable, spaceship-style statusline for [Claude Code](https://claude.ai/code).

## Quick Start

```sh
npx -y ccnow
```

That's it. Sensible defaults, zero config.

## Usage

### CLI Flags

Control which segments appear and in what order:

```sh
npx -y ccnow --pwd --sep --git --sep --context
```

### Presets

```sh
npx -y ccnow --preset=minimal
npx -y ccnow --preset=full
```

### Config File

```sh
npx -y ccnow --config ~/.claude/ccnow.json
```

### Claude Code Setup

Add to your `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "npx -y ccnow",
    "padding": 0
  }
}
```

## Segments

| Segment     | Description                         |
| ----------- | ----------------------------------- |
| `--pwd`     | Working directory (smart-truncated) |
| `--git`     | Branch name + diff stats            |
| `--context` | Context window token usage          |
| `--sep`     | Separator between segments          |

## Config File Format

```json
{
  "segments": [
    { "segment": "pwd.smart", "color": "cyan", "bold": true },
    { "segment": "sep", "char": "|", "dim": true },
    {
      "segment": "git",
      "children": [
        { "segment": "git.branch", "color": "white" },
        { "segment": "git.insertions", "color": "green", "prefix": "+" },
        { "segment": "git.deletions", "color": "red", "prefix": "-" }
      ]
    },
    { "segment": "sep", "char": "|", "dim": true },
    {
      "segment": "context",
      "children": [{ "segment": "context.tokens", "bold": true }, { "segment": "context.percent" }]
    }
  ]
}
```

## License

MIT
