package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

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
