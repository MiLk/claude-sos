package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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
