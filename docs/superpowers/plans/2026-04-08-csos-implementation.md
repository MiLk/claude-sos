# Claude Sobriety Selector (csos) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a humorous TUI CLI that launches Claude Code or Codex with preconfigured "sobriety levels."

**Architecture:** Single-file Go CLI using bubbletea for the interactive selector. Flag parsing with stdlib `flag`. Process replacement via `syscall.Exec` after resolving the command with `exec.LookPath`.

**Tech Stack:** Go 1.22+, bubbletea (TUI), lipgloss (styling)

---

## File Structure

| File | Responsibility |
|------|----------------|
| `go.mod` | Module definition and dependencies |
| `main.go` | Everything: levels, flags, TUI model, exec logic (~200 lines) |

---

### Task 1: Project Setup

**Files:**
- Create: `go.mod`
- Create: `main.go`

- [ ] **Step 1: Initialize Go module**

```bash
cd /Users/emilien/Dev/claude-sos
go mod init github.com/emilienMusic/claude-sos
```

- [ ] **Step 2: Add dependencies**

```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
```

- [ ] **Step 3: Create minimal main.go**

```go
package main

import "fmt"

func main() {
	fmt.Println("csos")
}
```

- [ ] **Step 4: Verify build**

Run: `go build -o csos && ./csos`
Expected: prints "csos"

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum main.go
git commit -m "feat: initialize csos project with dependencies"
```

---

### Task 2: Level Data Model

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Define Level struct and levels slice**

Add after package/imports:

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

var levels = []Level{
	{0, "Can't remember its name", "No memories, no regrets", "forgot", []string{"CLAUDE_CODE_SIMPLE=1"}, "claude", []string{"--model", "opus", "--effort", "medium"}, false},
	{1, "Went for nijikai", "Creativity on a budget", "nijikai", nil, "claude", []string{"--model", "opus", "--effort", "medium"}, false},
	{2, "Had some nihonshu", "Fresh start energy", "nihonshu", []string{"CLAUDE_CODE_AUTO_COMPACT_WINDOW=400000"}, "claude", []string{"--model", "opus", "--effort", "high"}, false},
	{3, "Had a few beers", "Your daily driver", "beers", nil, "claude", []string{"--model", "claude-opus-4-5-20251101", "--effort", "high"}, false},
	{4, "Had one beer", "Slightly tipsy, slightly verbose", "beer", nil, "claude", []string{"--model", "sonnet", "--effort", "medium"}, false},
	{5, "Stone cold sober", "Does exactly what you say. Exactly.", "sober", nil, "claude", []string{"--model", "claude-sonnet-4-5-20250929", "--effort", "medium"}, false},
	{6, "Just call the police", "When Claude needs adult supervision", "police", nil, "codex", []string{"--model", "gpt-5.4", "-c", `model_reasoning_effort="xhigh"`}, true},
}
```

- [ ] **Step 2: Add helper to find level by keyword or index**

```go
import (
	"fmt"
	"strconv"
)

func findLevel(s string) (*Level, error) {
	if idx, err := strconv.Atoi(s); err == nil {
		if idx >= 0 && idx < len(levels) {
			return &levels[idx], nil
		}
		return nil, fmt.Errorf("unknown level: %d. Use 0-%d or keyword", idx, len(levels)-1)
	}
	for i := range levels {
		if levels[i].Keyword == s {
			return &levels[i], nil
		}
	}
	keywords := ""
	for _, l := range levels {
		keywords += l.Keyword + ", "
	}
	return nil, fmt.Errorf("unknown level: %s. Try: %s", s, keywords[:len(keywords)-2])
}
```

- [ ] **Step 3: Verify build**

Run: `go build -o csos`
Expected: builds without error

- [ ] **Step 4: Commit**

```bash
git add main.go
git commit -m "feat: add level data model with 7 sobriety levels"
```

---

### Task 3: Flag Parsing

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Add flag parsing to main**

Replace main function:

```go
import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	levelFlag = flag.String("l", "", "level by index (0-6) or keyword")
	helpFlag  = flag.Bool("h", false, "show help")
)

func init() {
	flag.StringVar(levelFlag, "level", "", "level by index (0-6) or keyword")
	flag.BoolVar(helpFlag, "help", false, "show help")
}

func main() {
	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	var passthrough []string
	for i, arg := range os.Args[1:] {
		if arg == "--" {
			passthrough = os.Args[i+2:]
			break
		}
	}

	if *levelFlag != "" {
		level, err := findLevel(*levelFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		execute(level, passthrough)
		return
	}

	fmt.Println("TODO: show TUI")
}

func printHelp() {
	fmt.Println(`csos - Claude Sobriety Selector

Usage:
  csos                    Interactive level selection
  csos -l <level>         Direct selection (0-6 or keyword)
  csos -- <args>          Pass arguments to claude/codex

Levels:`)
	for _, l := range levels {
		fmt.Printf("  %d (%s): %s - %s\n", l.Index, l.Keyword, l.Name, l.Description)
	}
}

func execute(level *Level, passthrough []string) {
	fmt.Printf("Would execute: %s %v (passthrough: %v)\n", level.Command, level.Args, passthrough)
}
```

- [ ] **Step 2: Test help flag**

Run: `go build -o csos && ./csos -h`
Expected: shows help with all levels listed

- [ ] **Step 3: Test direct level selection**

Run: `./csos -l 3`
Expected: `Would execute: claude [--model claude-opus-4-5-20251101 --effort high] (passthrough: [])`

Run: `./csos -l beers`
Expected: same output

Run: `./csos -l 99`
Expected: error message about unknown level

- [ ] **Step 4: Test passthrough**

Run: `./csos -l 3 -- -p "test"`
Expected: `Would execute: claude [--model claude-opus-4-5-20251101 --effort high] (passthrough: [-p test])`

- [ ] **Step 5: Commit**

```bash
git add main.go
git commit -m "feat: add flag parsing for -l/--level and passthrough args"
```

---

### Task 4: TUI Model

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Add bubbletea model**

Add imports and model:

```go
import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	cursor   int
	selected *Level
	quitting bool
}

func initialModel() model {
	return model{cursor: 3} // default to "Had a few beers"
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(levels)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = &levels[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}
```

- [ ] **Step 2: Add View method with styling**

```go
var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func (m model) View() string {
	if m.quitting || m.selected != nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("  🍺 Claude Sobriety Selector"))
	b.WriteString("\n\n")

	for i, level := range levels {
		if level.IsSpecial && i > 0 {
			b.WriteString("  " + dimStyle.Render("─────────────────────────────") + "\n")
		}

		cursor := "  ○ "
		name := level.Name
		desc := "      " + level.Description

		if i == m.cursor {
			cursor = "  ● "
			name = selectedStyle.Render(name)
			desc = "      " + selectedStyle.Render(level.Description)
		} else {
			desc = "      " + dimStyle.Render(level.Description)
		}

		b.WriteString(cursor + name + "\n")
		b.WriteString(desc + "\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  ↑/↓ select • enter confirm • q quit"))
	b.WriteString("\n")

	return b.String()
}
```

- [ ] **Step 3: Integrate TUI into main**

Update main() to use TUI:

```go
func main() {
	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	var passthrough []string
	for i, arg := range os.Args[1:] {
		if arg == "--" {
			passthrough = os.Args[i+2:]
			break
		}
	}

	if *levelFlag != "" {
		level, err := findLevel(*levelFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		execute(level, passthrough)
		return
	}

	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	finalModel := m.(model)
	if finalModel.selected != nil {
		execute(finalModel.selected, passthrough)
	}
}
```

- [ ] **Step 4: Verify TUI renders**

Run: `go build -o csos && ./csos`
Expected: Interactive TUI with all levels, arrow navigation works, q quits, enter selects

- [ ] **Step 5: Commit**

```bash
git add main.go
git commit -m "feat: add interactive TUI with bubbletea"
```

---

### Task 5: Launch Messages

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Add launch message function**

```go
func launchMessage(level *Level) string {
	switch level.Index {
	case 0:
		return "Forgetting everything you told me..."
	case 1, 2:
		return "Pouring some sake..."
	case 3, 4:
		return "Cracking open a cold one..."
	case 5:
		return "Putting on the serious face..."
	case 6:
		return "Calling for backup..."
	default:
		return "Launching..."
	}
}
```

- [ ] **Step 2: Update execute to print message**

```go
func execute(level *Level, passthrough []string) {
	fmt.Println(launchMessage(level))
	fmt.Printf("Would execute: %s %v (passthrough: %v)\n", level.Command, level.Args, passthrough)
}
```

- [ ] **Step 3: Test launch messages**

Run: `./csos -l 0`
Expected: "Forgetting everything you told me..."

Run: `./csos -l 6`
Expected: "Calling for backup..."

- [ ] **Step 4: Commit**

```bash
git add main.go
git commit -m "feat: add witty launch messages per level"
```

---

### Task 6: Exec Implementation

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Add exec imports**

```go
import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
```

- [ ] **Step 2: Implement real execute function**

Replace the execute function:

```go
func execute(level *Level, passthrough []string) {
	binary, err := exec.LookPath(level.Command)
	if err != nil {
		if level.Command == "claude" {
			fmt.Fprintln(os.Stderr, "claude not found. Install: https://claude.ai/download")
		} else {
			fmt.Fprintln(os.Stderr, "codex not found. Install: npm i -g @openai/codex")
		}
		os.Exit(1)
	}

	fmt.Println(launchMessage(level))

	// Build argv: command name + level args + passthrough
	argv := append([]string{level.Command}, level.Args...)
	argv = append(argv, passthrough...)

	// Build env: inherit current env, overlay level-specific vars
	env := os.Environ()
	for _, e := range level.Env {
		env = append(env, e)
	}

	// Replace this process
	if err := syscall.Exec(binary, argv, env); err != nil {
		fmt.Fprintln(os.Stderr, "exec failed:", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Test with real claude (if installed)**

Run: `./csos -l 5 -- --version`
Expected: prints launch message, then claude version output

- [ ] **Step 4: Test command not found**

Run: `PATH="" ./csos -l 0`
Expected: "claude not found. Install: https://claude.ai/download"

- [ ] **Step 5: Commit**

```bash
git add main.go
git commit -m "feat: implement real exec with PATH lookup and env inheritance"
```

---

### Task 7: Final Polish

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Ensure all imports are used and code compiles**

Run: `go build -o csos`
Expected: builds without warnings

- [ ] **Step 2: Run gofmt**

```bash
gofmt -w main.go
```

- [ ] **Step 3: Add build constraint for Unix**

Add at top of main.go after package line:

```go
//go:build unix

package main
```

- [ ] **Step 4: Final integration test**

Test full flow:
1. `./csos` - TUI appears, navigate with arrows, press q to quit
2. `./csos` - TUI appears, select level 3, press enter - launches claude
3. `./csos -l nihonshu -- -p "hello"` - direct launch with passthrough
4. `./csos -h` - shows help

- [ ] **Step 5: Commit**

```bash
git add main.go
git commit -m "chore: add unix build constraint, format code"
```

---

### Task 8: README

**Files:**
- Create: `README.md`

- [ ] **Step 1: Create README**

```markdown
# csos - Claude Sobriety Selector

A humorous CLI for launching Claude Code at various "sobriety levels."

## Install

```bash
go install github.com/emilienMusic/claude-sos@latest
```

Or build from source:

```bash
git clone https://github.com/emilienMusic/claude-sos
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
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README with install and usage instructions"
```

---

## Summary

8 tasks, ~35 steps total. Each task produces a working, committable state. The implementation follows TDD-lite (verify behavior before committing) since TUI testing is interactive.
