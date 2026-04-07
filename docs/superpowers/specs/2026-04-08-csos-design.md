# Claude Sobriety Selector (csos) — Design Spec

A humorous CLI wrapper that launches Claude Code or Codex with different configurations framed as "sobriety levels."

## Overview

`csos` presents an interactive TUI for selecting a "sobriety level" that maps to specific model and effort configurations. Lower numbers = more "drunk" (creative/chaotic), higher numbers = more "sober" (disciplined/precise).

## Levels

| # | Keyword | Name | Description | Command |
|---|---------|------|-------------|---------|
| 0 | `forgot` | Can't remember its name | No memories, no regrets | `CLAUDE_CODE_SIMPLE=1 claude --model opus --effort medium` |
| 1 | `nijikai` | Went for nijikai | Creativity on a budget | `claude --model opus --effort medium` |
| 2 | `nihonshu` | Had some nihonshu | Fresh start energy | `CLAUDE_CODE_AUTO_COMPACT_WINDOW=400000 claude --model opus --effort high` |
| 3 | `beers` | Had a few beers | Your daily driver | `claude --model claude-opus-4-5-20251101 --effort high` |
| 4 | `beer` | Had one beer | Slightly tipsy, slightly verbose | `claude --model sonnet --effort medium` |
| 5 | `sober` | Stone cold sober | Does exactly what you say. Exactly. | `claude --model claude-sonnet-4-5-20250929 --effort medium` |
| 6 | `police` | Just call the police | When Claude needs adult supervision | `codex --model gpt-5.4 -c model_reasoning_effort="xhigh"` |

Level 6 (police) is visually separated as a special "escape hatch" option.

## CLI Interface

```bash
# Interactive mode (default)
csos

# Direct selection by number or keyword
csos -l 3
csos --level nihonshu

# With passthrough arguments (everything after --)
csos -- -p "fix the bug"
csos -l 3 -- -p "fix the bug" --allowedTools Edit,Read

# Help
csos -h
```

**Flags:**
- `-l` / `--level`: Skip interactive UI, select directly by index (0-6) or keyword
- `-h` / `--help`: Show usage help
- `--`: Separator for passthrough arguments

## TUI Design

```
  🍺 Claude Sobriety Selector

  ○ Can't remember its name
      No memories, no regrets
  ○ Went for nijikai
      Creativity on a budget
  ○ Had some nihonshu
      Fresh start energy
  ● Had a few beers
      Your daily driver
  ○ Had one beer
      Slightly tipsy, slightly verbose
  ○ Stone cold sober
      Does exactly what you say. Exactly.
  ─────────────────────────────
  ○ Just call the police
      When Claude needs adult supervision

  ↑/↓ select • enter confirm • q quit
```

**Controls:**
- Arrow keys (↑/↓) or `j/k` to navigate
- Enter to select and launch
- `q` or Esc to quit without launching

**Visual treatment:**
- Small ASCII header with beer emoji
- Selected item highlighted (bold + color)
- Visual separator line before "police" level

## Launch Messages

Brief witty message shown before exec:

| Levels | Message |
|--------|---------|
| 0 | "Forgetting everything you told me..." |
| 1-2 | "Pouring some sake..." |
| 3-4 | "Cracking open a cold one..." |
| 5 | "Putting on the serious face..." |
| 6 | "Calling for backup..." |

## Data Model

```go
type Level struct {
    Index       int
    Name        string
    Description string
    Keyword     string
    Env         []string
    Command     string
    Args        []string
    IsSpecial   bool
}
```

## Execution Flow

1. Parse flags (`-l`, `-h`, `--`)
2. If `-l` provided: resolve level (int or keyword), skip TUI
3. Otherwise: show TUI, wait for selection
4. Print launch message
5. `exec.LookPath()` to resolve command to absolute path (enables PATH search + "not found" errors)
6. Build environment: start from `os.Environ()`, overlay `Level.Env` vars
7. Build argv: absolute path + level args + passthrough args
8. `syscall.Exec()` to replace process (clean TTY passthrough)

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Invalid level number | `"Unknown level: 99. Use 0-6 or keyword."`, exit 1 |
| Invalid keyword | `"Unknown level: vodka. Try: forgot, nijikai, nihonshu, beers, beer, sober, police"`, exit 1 |
| `claude` not in PATH | `"claude not found. Install: https://claude.ai/download"`, exit 1 |
| `codex` not in PATH | `"codex not found. Install: npm i -g @openai/codex"`, exit 1 |
| User quits TUI | Exit 0 silently |

## Technical Choices

- **Language:** Go
- **TUI:** Bubbletea
- **Structure:** Single file (`main.go`, ~200-250 lines)
- **CLI parsing:** `flag` package
- **Process exec:** `syscall.Exec` (replaces process, no wrapper)
- **Platform:** macOS and Linux only (`syscall.Exec` is Unix-specific)

## Design Notes

**Model alias strategy:** Some levels use floating aliases (`opus`, `sonnet`) that track the current best model, while others use pinned versions (`claude-opus-4-5-20251101`) for reproducible behavior. This is intentional — "daily driver" drifts with releases, "stone cold sober" stays fixed.

**Passthrough overrides:** Passthrough arguments are appended after level args, so users can override level settings (e.g., `csos -l 5 -- --model opus`). This is intentional — csos is a launcher, not a guardian. Levels are defaults, not constraints.

## Non-Goals

- Configuration file for custom levels
- Persistent state or history
- Logging
- Subcommands
- Windows support
