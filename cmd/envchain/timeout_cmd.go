package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/yourorg/envchain-cli/internal/timeout"
	"github.com/urfave/cli/v2"
)

func defaultTimeoutDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "envchain", "timeouts")
}

func runTimeout(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "set":
		return runTimeoutSet(c)
	case "get":
		return runTimeoutGet(c)
	case "delete":
		return runTimeoutDelete(c)
	default:
		return fmt.Errorf("timeout: unknown subcommand %q — use set, get, or delete", subcmd)
	}
}

func runTimeoutSet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 2 {
		return fmt.Errorf("usage: envchain timeout set <project> <minutes>")
	}
	project := args[0]
	minutes, err := strconv.Atoi(args[1])
	if err != nil || minutes < 0 {
		return fmt.Errorf("timeout: minutes must be a non-negative integer")
	}
	m := timeout.New(defaultTimeoutDir())
	rule := timeout.Rule{
		Project:  project,
		Duration: time.Duration(minutes) * time.Minute,
		Enabled:  minutes > 0,
	}
	if err := m.Set(rule); err != nil {
		return err
	}
	if minutes == 0 {
		fmt.Fprintf(c.App.Writer, "timeout disabled for project %q\n", project)
	} else {
		fmt.Fprintf(c.App.Writer, "timeout set to %d minute(s) for project %q\n", minutes, project)
	}
	return nil
}

func runTimeoutGet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain timeout get <project>")
	}
	m := timeout.New(defaultTimeoutDir())
	rule, err := m.Get(args[0])
	if err == timeout.ErrNotFound {
		fmt.Fprintf(c.App.Writer, "no timeout configured for project %q\n", args[0])
		return nil
	}
	if err != nil {
		return err
	}
	status := "disabled"
	if rule.Enabled {
		status = fmt.Sprintf("%.0f minute(s)", rule.Duration.Minutes())
	}
	fmt.Fprintf(c.App.Writer, "project: %s\ntimeout: %s\n", rule.Project, status)
	return nil
}

func runTimeoutDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain timeout delete <project>")
	}
	m := timeout.New(defaultTimeoutDir())
	if err := m.Delete(args[0]); err == timeout.ErrNotFound {
		fmt.Fprintf(c.App.Writer, "no timeout rule found for project %q\n", args[0])
		return nil
	} else if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "timeout rule deleted for project %q\n", args[0])
	return nil
}
