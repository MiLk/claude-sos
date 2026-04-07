//go:build darwin || linux

package main

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

func findLevel(s string) (*Level, error) {
	if idx, err := strconv.Atoi(s); err == nil {
		if idx >= 0 && idx < len(levels) {
			return &levels[idx], nil
		}
		return nil, fmt.Errorf("Unknown level: %d. Use 0-6 or keyword.", idx)
	}
	for i := range levels {
		if levels[i].Keyword == s {
			return &levels[i], nil
		}
	}
	return nil, fmt.Errorf("Unknown level: %s. Try: forgot, nijikai, nihonshu, beers, beer, sober, police", s)
}

var (
	levelFlag = flag.String("l", "", "level by index (0-6) or keyword")
	helpFlag  = flag.Bool("h", false, "show help")
)

func init() {
	flag.StringVar(levelFlag, "level", "", "level by index (0-6) or keyword")
	flag.BoolVar(helpFlag, "help", false, "show help")
}

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

func printHelp() {
	fmt.Println(`csos - Claude Sobriety Selector

Usage:
  csos                        Interactive level selection
  csos -l, --level <level>    Direct selection (0-6 or keyword)
  csos -h, --help             Show this help
  csos -- <args>              Pass args to selected backend (claude or codex)

Note: Passthrough args are forwarded verbatim to the selected level's CLI.
      Level 6 (police) uses codex, which has different flags than claude.

Levels:`)
	for _, l := range levels {
		fmt.Printf("  %d (%s): %s - %s\n", l.Index, l.Keyword, l.Name, l.Description)
	}
}

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

func overlayEnv(base []string, overlay []string) []string {
	env := make(map[string]string)
	for _, e := range base {
		if idx := strings.Index(e, "="); idx != -1 {
			env[e[:idx]] = e[idx+1:]
		}
	}
	for _, e := range overlay {
		if idx := strings.Index(e, "="); idx != -1 {
			env[e[:idx]] = e[idx+1:]
		}
	}
	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, k+"="+v)
	}
	return result
}

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

	// Build env: inherit current env, overlay level-specific vars (replacing duplicates)
	env := overlayEnv(os.Environ(), level.Env)

	// Replace this process
	if err := syscall.Exec(binary, argv, env); err != nil {
		fmt.Fprintln(os.Stderr, "exec failed:", err)
		os.Exit(1)
	}
}
