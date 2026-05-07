package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain-cli/internal/immutable"
	"github.com/envchain/envchain-cli/internal/store"
	"github.com/urfave/cli/v2"
)

const defaultImmutableDir = ".envchain"

func runImmutable(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "set":
		return runImmutableSet(c)
	case "get":
		return runImmutableGet(c)
	case "delete":
		return runImmutableDelete(c)
	default:
		return fmt.Errorf("immutable: unknown subcommand %q (use set|get|delete)", subcmd)
	}
}

func immutableManager(c *cli.Context) (*immutable.Manager, error) {
	dir := c.String("store")
	if dir == "" {
		dir = defaultImmutableDir
	}
	st, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("immutable: open store: %w", err)
	}
	return immutable.New(st), nil
}

func runImmutableSet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 2 {
		return fmt.Errorf("usage: immutable set <project> <KEY> [KEY...]")
	}
	project := args[0]
	keys := args[1:]
	pass, err := resolvePassphrase(c, false)
	if err != nil {
		return err
	}
	mgr, err := immutableManager(c)
	if err != nil {
		return err
	}
	if err := mgr.Set(project, keys, pass); err != nil {
		return fmt.Errorf("immutable set: %w", err)
	}
	fmt.Fprintf(os.Stdout, "immutable keys set for project %q: %s\n", project, strings.Join(keys, ", "))
	return nil
}

func runImmutableGet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return fmt.Errorf("usage: immutable get <project>")
	}
	project := args[0]
	pass, err := resolvePassphrase(c, false)
	if err != nil {
		return err
	}
	mgr, err := immutableManager(c)
	if err != nil {
		return err
	}
	keys, err := mgr.Get(project, pass)
	if err != nil {
		return fmt.Errorf("immutable get: %w", err)
	}
	if len(keys) == 0 {
		fmt.Fprintln(os.Stdout, "(no immutable keys)")
		return nil
	}
	for _, k := range keys {
		fmt.Fprintln(os.Stdout, k)
	}
	return nil
}

func runImmutableDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return fmt.Errorf("usage: immutable delete <project>")
	}
	project := args[0]
	mgr, err := immutableManager(c)
	if err != nil {
		return err
	}
	if err := mgr.Delete(project); err != nil {
		return fmt.Errorf("immutable delete: %w", err)
	}
	fmt.Fprintf(os.Stdout, "immutable record deleted for project %q\n", project)
	return nil
}
