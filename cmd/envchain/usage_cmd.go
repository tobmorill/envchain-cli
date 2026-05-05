package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/envchain/envchain-cli/internal/store"
	"github.com/envchain/envchain-cli/internal/usage"
	"github.com/urfave/cli/v2"
)

func defaultUsageDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "usage")
}

func runUsage(c *cli.Context) error {
	sub := c.Args().First()
	switch sub {
	case "get":
		return runUsageGet(c)
	case "reset":
		return runUsageReset(c)
	default:
		return cli.ShowSubcommandHelp(c)
	}
}

func usageManager(c *cli.Context) (*usage.Manager, error) {
	dir := c.String("usage-dir")
	if dir == "" {
		dir = defaultUsageDir()
	}
	s, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("usage store: %w", err)
	}
	return usage.New(s), nil
}

func runUsageGet(c *cli.Context) error {
	project := c.Args().Get(1)
	if project == "" {
		return fmt.Errorf("usage get <project>")
	}

	m, err := usageManager(c)
	if err != nil {
		return err
	}

	rec, err := m.Get(project)
	if err != nil {
		return err
	}
	if rec == nil {
		fmt.Fprintf(c.App.Writer, "no usage record for project %q\n", project)
		return nil
	}

	w := tabwriter.NewWriter(c.App.Writer, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "project\t%s\n", rec.Project)
	fmt.Fprintf(w, "count\t%d\n", rec.Count)
	fmt.Fprintf(w, "first_used\t%s\n", rec.FirstUsed.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "last_used\t%s\n", rec.LastUsed.Format("2006-01-02 15:04:05 UTC"))
	return w.Flush()
}

func runUsageReset(c *cli.Context) error {
	project := c.Args().Get(1)
	if project == "" {
		return fmt.Errorf("usage reset <project>")
	}

	m, err := usageManager(c)
	if err != nil {
		return err
	}

	if err := m.Reset(project); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "usage record for %q cleared\n", project)
	return nil
}
