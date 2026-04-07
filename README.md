# csos - Claude Sobriety Selector

A humorous CLI for launching Claude Code at various "sobriety levels."

## Install

```bash
go install github.com/MiLk/claude-sos@latest
```

Or build from source:

```bash
git clone https://github.com/MiLk/claude-sos
cd claude-sos
go build -o csos
```

## Usage

```bash
# Interactive selection
csos

# Direct selection by number or keyword
csos -l 3
csos -l beers

# Pass arguments to claude/codex
csos -l 3 -- -p "fix the bug"
```

## Levels

| # | Keyword | Name | Description |
|---|---------|------|-------------|
| 0 | forgot | Can't remember its name | No memories, no regrets |
| 1 | nijikai | Went for nijikai | Creativity on a budget |
| 2 | nihonshu | Had some nihonshu | Fresh start energy |
| 3 | beers | Had a few beers | Your daily driver |
| 4 | beer | Had one beer | Slightly tipsy, slightly verbose |
| 5 | sober | Stone cold sober | Does exactly what you say. Exactly. |
| 6 | police | Just call the police | When Claude needs adult supervision |

## Requirements

- macOS or Linux
- `claude` CLI installed for levels 0-5
- `codex` CLI installed for level 6 (police)
