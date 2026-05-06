package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/envchain/envchain-cli/internal/priority"
	"github.com/envchain/envchain-cli/internal/store"
	"github.com/urfave/cli/v2"
)

const defaultPriorityDir = ".envchain"

func priorityManager(c *cli.Context) (*priority.Manager, error) {
	dir := c.String("store")
	if dir == "" {
		dir = defaultPriorityDir
	}
	st, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	return priority.New(st), nil
}

func runPriority(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "set":
		return runPrioritySet(c)
	case "get":
		return runPriorityGet(c)
	case "delete", "del":
		return runPriorityDelete(c)
	case "list", "ls":
		return runPriorityList(c)
	default:
		return fmt.Errorf("priority: unknown subcommand %q (use set|get|delete|list)", subcmd)
	}
}

func runPrioritySet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 3 {
		return errors.New("usage: priority set <project> <key> <low|normal|high>")
	}
	project, key, rawLevel := args[0], args[1], args[2]
	lvl, err := priority.ParseLevel(rawLevel)
	if err != nil {
		return err
	}
	m, err := priorityManager(c)
	if err != nil {
		return err
	}
	if err := m.Set(project, key, lvl); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "priority: %s/%s set to %s\n", project, key, lvl)
	return nil
}

func runPriorityGet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 2 {
		return errors.New("usage: priority get <project> <key>")
	}
	m, err := priorityManager(c)
	if err != nil {
		return err
	}
	lvl, err := m.Get(args[0], args[1])
	if err != nil {
		return err
	}
	fmt.Fprintln(c.App.Writer, lvl)
	return nil
}

func runPriorityDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 2 {
		return errors.New("usage: priority delete <project> <key>")
	}
	m, err := priorityManager(c)
	if err != nil {
		return err
	}
	if err := m.Delete(args[0], args[1]); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "priority: %s/%s reset to normal\n", args[0], args[1])
	return nil
}

func runPriorityList(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: priority list <project>")
	}
	m, err := priorityManager(c)
	if err != nil {
		return err
	}
	all, err := m.GetAll(args[0])
	if err != nil {
		return err
	}
	if len(all) == 0 {
		fmt.Fprintln(c.App.Writer, "(no priority entries)")
		return nil
	}
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tPRIORITY")
	for _, k := range keys {
		fmt.Fprintf(tw, "%s\t%s\n", k, all[k])
	}
	return tw.Flush()
}
