package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/yourorg/envchain-cli/internal/history"
)

func defaultHistoryDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "history")
}

func runHistory(c *cli.Context) error {
	project := c.Args().First()
	if project == "" {
		return fmt.Errorf("project name required")
	}

	m := history.New(defaultHistoryDir())
	entries, err := m.ReadAll(project)
	if err == history.ErrNotFound {
		fmt.Fprintf(c.App.Writer, "no history found for project %q\n", project)
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading history: %w", err)
	}

	w := tabwriter.NewWriter(c.App.Writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIME\tACTION")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\n",
			e.AccessedAt.Local().Format(time.RFC3339),
			e.Action,
		)
	}
	return w.Flush()
}

func runHistoryClear(c *cli.Context) error {
	project := c.Args().First()
	if project == "" {
		return fmt.Errorf("project name required")
	}

	m := history.New(defaultHistoryDir())
	if err := m.Clear(project); err != nil {
		return fmt.Errorf("clearing history: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "history cleared for project %q\n", project)
	return nil
}
