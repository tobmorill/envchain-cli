package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/envchain/envchain-cli/internal/access"
	"github.com/envchain/envchain-cli/internal/store"
	"github.com/urfave/cli/v2"
)

func defaultAccessDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "access")
}

func runAccess(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "get":
		return runAccessGet(c)
	case "reset":
		return runAccessReset(c)
	default:
		return runAccessGet(c)
	}
}

func runAccessGet(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		var err error
		project, err = resolveProject(c)
		if err != nil {
			return err
		}
	}

	st, err := store.New(defaultAccessDir())
	if err != nil {
		return fmt.Errorf("access store: %w", err)
	}
	m := access.New(st)

	rec, err := m.Get(project)
	if err != nil {
		return fmt.Errorf("no access record for %q: %w", project, err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Project:\t%s\n", rec.Project)
	fmt.Fprintf(w, "Access count:\t%d\n", rec.Count)
	fmt.Fprintf(w, "First used:\t%s\n", rec.FirstUsed.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Last used:\t%s\n", rec.LastUsed.Format("2006-01-02 15:04:05 UTC"))
	return w.Flush()
}

func runAccessReset(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		var err error
		project, err = resolveProject(c)
		if err != nil {
			return err
		}
	}

	st, err := store.New(defaultAccessDir())
	if err != nil {
		return fmt.Errorf("access store: %w", err)
	}
	m := access.New(st)

	if err := m.Reset(project); err != nil {
		return fmt.Errorf("reset access record: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Access record for %q cleared.\n", project)
	return nil
}

func resolveProject(c *cli.Context) (string, error) {
	if p := c.String("project"); p != "" {
		return p, nil
	}
	return "", fmt.Errorf("--project flag is required")
}
