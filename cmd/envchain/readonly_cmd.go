package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/envchain/envchain-cli/internal/readonly"
	"github.com/envchain/envchain-cli/internal/store"
	"github.com/urfave/cli/v2"
)

func defaultReadonlyDir() string {
	return defaultStorePath("readonly")
}

func runReadonly(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "set":
		return runReadonlySet(c)
	case "unset":
		return runReadonlyUnset(c)
	case "get":
		return runReadonlyGet(c)
	default:
		return fmt.Errorf("unknown subcommand %q; use set, unset, or get", subcmd)
	}
}

func readonlyManager(c *cli.Context) (*readonly.Manager, func(), error) {
	dir := c.String("readonly-dir")
	if dir == "" {
		dir = defaultReadonlyDir()
	}
	st, err := store.New(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("open store: %w", err)
	}
	return readonly.New(st), func() { st.Close() }, nil
}

func runReadonlySet(c *cli.Context) error {
	project := c.Args().Get(1)
	if project == "" {
		return errors.New("usage: readonly set <project>")
	}
	mgr, cleanup, err := readonlyManager(c)
	if err != nil {
		return err
	}
	defer cleanup()
	if err := mgr.Set(project, true); err != nil {
		return fmt.Errorf("set read-only: %w", err)
	}
	fmt.Fprintf(os.Stdout, "project %q marked as read-only\n", project)
	return nil
}

func runReadonlyUnset(c *cli.Context) error {
	project := c.Args().Get(1)
	if project == "" {
		return errors.New("usage: readonly unset <project>")
	}
	mgr, cleanup, err := readonlyManager(c)
	if err != nil {
		return err
	}
	defer cleanup()
	if err := mgr.Set(project, false); err != nil {
		return fmt.Errorf("unset read-only: %w", err)
	}
	fmt.Fprintf(os.Stdout, "project %q is now writable\n", project)
	return nil
}

func runReadonlyGet(c *cli.Context) error {
	project := c.Args().Get(1)
	if project == "" {
		return errors.New("usage: readonly get <project>")
	}
	mgr, cleanup, err := readonlyManager(c)
	if err != nil {
		return err
	}
	defer cleanup()
	rec, err := mgr.Get(project)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	if rec.ReadOnly {
		fmt.Fprintf(os.Stdout, "%s: read-only\n", project)
	} else {
		fmt.Fprintf(os.Stdout, "%s: writable\n", project)
	}
	return nil
}
