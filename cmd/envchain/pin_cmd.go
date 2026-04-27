package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/envchain/envchain-cli/internal/pin"
	"github.com/urfave/cli/v2"
)

func defaultPinDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "pins")
}

func runPin(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "set":
		return runPinSet(c)
	case "get":
		return runPinGet(c)
	case "delete", "rm":
		return runPinDelete(c)
	default:
		return fmt.Errorf("pin: unknown subcommand %q — use set, get, or delete", subcmd)
	}
}

func runPinSet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 2 {
		return fmt.Errorf("usage: envchain pin set <project> <KEY> [KEY...]")
	}
	project := args[0]
	keys := args[1:]
	m := pin.New(defaultPinDir())
	if err := m.Set(project, keys); err != nil {
		return fmt.Errorf("pin set: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "pinned %d key(s) for project %q\n", len(keys), project)
	return nil
}

func runPinGet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain pin get <project>")
	}
	project := args[0]
	m := pin.New(defaultPinDir())
	keys, err := m.Get(project)
	if err != nil {
		return fmt.Errorf("pin get: %w", err)
	}
	fmt.Fprintln(c.App.Writer, strings.Join(keys, "\n"))
	return nil
}

func runPinDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain pin delete <project>")
	}
	project := args[0]
	m := pin.New(defaultPinDir())
	if err := m.Delete(project); err != nil {
		return fmt.Errorf("pin delete: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "pins removed for project %q\n", project)
	return nil
}
