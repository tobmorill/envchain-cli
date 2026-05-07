package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/envchain/envchain-cli/internal/store"
	"github.com/envchain/envchain-cli/internal/visibility"
)

const defaultVisibilityDir = ".envchain"

func runVisibility(args []string) error {
	fs := flag.NewFlagSet("visibility", flag.ContinueOnError)
	storeDir := fs.String("store", defaultVisibilityDir, "path to envchain store directory")
	if err := fs.Parse(args); err != nil {
		return err
	}

	subArgs := fs.Args()
	if len(subArgs) == 0 {
		return fmt.Errorf("usage: envchain visibility <set|get|delete|list> [options]")
	}

	st, err := store.New(filepath.Join(*storeDir, "visibility"))
	if err != nil {
		return fmt.Errorf("visibility: open store: %w", err)
	}
	mgr := visibility.New(st)

	switch subArgs[0] {
	case "set":
		return runVisibilitySet(mgr, subArgs[1:])
	case "get":
		return runVisibilityGet(mgr, subArgs[1:])
	case "delete":
		return runVisibilityDelete(mgr, subArgs[1:])
	case "list":
		return runVisibilityList(mgr, subArgs[1:])
	default:
		return fmt.Errorf("visibility: unknown sub-command %q", subArgs[0])
	}
}

func runVisibilitySet(mgr *visibility.Manager, args []string) error {
	fs := flag.NewFlagSet("visibility set", flag.ContinueOnError)
	project := fs.String("project", "", "project name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *project == "" || fs.NArg() < 2 {
		return fmt.Errorf("usage: visibility set -project <name> <key> <visible|hidden>")
	}
	key := fs.Arg(0)
	lvl := visibility.Level(fs.Arg(1))
	if lvl != visibility.LevelVisible && lvl != visibility.LevelHidden {
		return fmt.Errorf("visibility: level must be 'visible' or 'hidden', got %q", lvl)
	}
	if err := mgr.Set(*project, key, lvl); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "visibility: %s/%s set to %s\n", *project, key, lvl)
	return nil
}

func runVisibilityGet(mgr *visibility.Manager, args []string) error {
	fs := flag.NewFlagSet("visibility get", flag.ContinueOnError)
	project := fs.String("project", "", "project name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *project == "" || fs.NArg() < 1 {
		return fmt.Errorf("usage: visibility get -project <name> <key>")
	}
	lvl, err := mgr.Get(*project, fs.Arg(0))
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(lvl))
	return nil
}

func runVisibilityDelete(mgr *visibility.Manager, args []string) error {
	fs := flag.NewFlagSet("visibility delete", flag.ContinueOnError)
	project := fs.String("project", "", "project name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *project == "" || fs.NArg() < 1 {
		return fmt.Errorf("usage: visibility delete -project <name> <key>")
	}
	return mgr.Delete(*project, fs.Arg(0))
}

func runVisibilityList(mgr *visibility.Manager, args []string) error {
	fs := flag.NewFlagSet("visibility list", flag.ContinueOnError)
	project := fs.String("project", "", "project name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *project == "" {
		return fmt.Errorf("usage: visibility list -project <name>")
	}
	all, err := mgr.GetAll(*project)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tLEVEL")
	for k, lvl := range all {
		fmt.Fprintf(w, "%s\t%s\n", k, lvl)
	}
	return w.Flush()
}
