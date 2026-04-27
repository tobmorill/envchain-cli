package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourorg/envchain-cli/internal/hook"
	"github.com/urfave/cli/v2"
)

func defaultHookDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "envchain", "hooks")
}

func runHook(c *cli.Context) error {
	args := c.Args().Slice()
	if len(args) == 0 {
		return fmt.Errorf("usage: envchain hook <set|get|delete|list> ...")
	}

	m := hook.New(defaultHookDir())
	subcmd := args[0]

	switch subcmd {
	case "set":
		if len(args) < 4 {
			return fmt.Errorf("usage: envchain hook set <project> <pre|post> <command>")
		}
		project := args[1]
		phase := hook.Phase(args[2])
		cmd := strings.Join(args[3:], " ")
		h := hook.Hook{Project: project, Phase: phase, Command: cmd}
		if err := m.Set(h); err != nil {
			return fmt.Errorf("hook set: %w", err)
		}
		fmt.Fprintf(c.App.Writer, "hook set for project %q phase %q\n", project, phase)

	case "get":
		if len(args) < 3 {
			return fmt.Errorf("usage: envchain hook get <project> <pre|post>")
		}
		h, ok, err := m.Get(args[1], hook.Phase(args[2]))
		if err != nil {
			return fmt.Errorf("hook get: %w", err)
		}
		if !ok {
			fmt.Fprintln(c.App.Writer, "no hook registered")
			return nil
		}
		fmt.Fprintf(c.App.Writer, "[%s] %s\n", h.Phase, h.Command)

	case "delete":
		if len(args) < 3 {
			return fmt.Errorf("usage: envchain hook delete <project> <pre|post>")
		}
		if err := m.Delete(args[1], hook.Phase(args[2])); err != nil {
			return fmt.Errorf("hook delete: %w", err)
		}
		fmt.Fprintln(c.App.Writer, "hook deleted")

	case "list":
		if len(args) < 2 {
			return fmt.Errorf("usage: envchain hook list <project>")
		}
		hooks, err := m.List(args[1])
		if err != nil {
			return fmt.Errorf("hook list: %w", err)
		}
		if len(hooks) == 0 {
			fmt.Fprintln(c.App.Writer, "no hooks registered")
			return nil
		}
		for _, h := range hooks {
			fmt.Fprintf(c.App.Writer, "%-6s  %s\n", h.Phase, h.Command)
		}

	default:
		return fmt.Errorf("unknown hook subcommand %q", subcmd)
	}
	return nil
}
