package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/yourorg/envchain-cli/internal/alias"
	"github.com/urfave/cli/v2"
)

func defaultAliasPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "envchain", "aliases.json")
}

func runAlias(c *cli.Context) error {
	mgr := alias.New(defaultAliasPath())

	switch c.Args().First() {
	case "set":
		args := c.Args().Tail()
		if len(args) < 2 {
			return fmt.Errorf("usage: envchain alias set <alias> <chain>")
		}
		if err := mgr.Set(args[0], args[1]); err != nil {
			return fmt.Errorf("set alias: %w", err)
		}
		fmt.Fprintf(c.App.Writer, "alias %q -> %q saved\n", args[0], args[1])
		return nil

	case "get":
		args := c.Args().Tail()
		if len(args) < 1 {
			return fmt.Errorf("usage: envchain alias get <alias>")
		}
		chain, err := mgr.Get(args[0])
		if err != nil {
			return fmt.Errorf("get alias: %w", err)
		}
		fmt.Fprintln(c.App.Writer, chain)
		return nil

	case "delete", "rm":
		args := c.Args().Tail()
		if len(args) < 1 {
			return fmt.Errorf("usage: envchain alias delete <alias>")
		}
		if err := mgr.Delete(args[0]); err != nil {
			return fmt.Errorf("delete alias: %w", err)
		}
		fmt.Fprintf(c.App.Writer, "alias %q deleted\n", args[0])
		return nil

	case "list", "ls", "":
		all, err := mgr.List()
		if err != nil {
			return fmt.Errorf("list aliases: %w", err)
		}
		if len(all) == 0 {
			fmt.Fprintln(c.App.Writer, "no aliases defined")
			return nil
		}
		keys := make([]string, 0, len(all))
		for k := range all {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		w := tabwriter.NewWriter(c.App.Writer, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ALIAS\tCHAIN")
		for _, k := range keys {
			fmt.Fprintf(w, "%s\t%s\n", k, all[k])
		}
		w.Flush()
		return nil

	default:
		return fmt.Errorf("unknown subcommand %q; use set, get, delete, or list", c.Args().First())
	}
}
