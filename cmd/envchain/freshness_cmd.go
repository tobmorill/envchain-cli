package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envchain/envchain-cli/internal/freshness"
	"github.com/envchain/envchain-cli/internal/store"
	"github.com/urfave/cli/v2"
)

func defaultFreshnessDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "freshness")
}

func runFreshness(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "touch":
		return runFreshnessTouch(c)
	case "get":
		return runFreshnessGet(c)
	case "delete":
		return runFreshnessDelete(c)
	default:
		return fmt.Errorf("freshness: unknown subcommand %q (touch|get|delete)", subcmd)
	}
}

func freshnessManager(c *cli.Context) (*freshness.Manager, error) {
	dir := c.String("freshness-dir")
	if dir == "" {
		dir = defaultFreshnessDir()
	}
	st, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("freshness: open store: %w", err)
	}
	return freshness.New(st), nil
}

func runFreshnessTouch(c *cli.Context) error {
	project := c.Args().Get(1)
	if project == "" {
		return errors.New("usage: freshness touch <project>")
	}
	m, err := freshnessManager(c)
	if err != nil {
		return err
	}
	if err := m.Touch(project); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "touched freshness record for %q\n", project)
	return nil
}

func runFreshnessGet(c *cli.Context) error {
	project := c.Args().Get(1)
	if project == "" {
		return errors.New("usage: freshness get <project> [--threshold <duration>]")
	}
	m, err := freshnessManager(c)
	if err != nil {
		return err
	}
	rec, err := m.Get(project)
	if errors.Is(err, freshness.ErrNotFound) {
		fmt.Fprintf(c.App.Writer, "no freshness record for %q\n", project)
		return nil
	}
	if err != nil {
		return err
	}
	thresholdStr := c.String("threshold")
	staleLabel := ""
	if thresholdStr != "" {
		threshold, err := time.ParseDuration(thresholdStr)
		if err != nil {
			return fmt.Errorf("invalid threshold: %w", err)
		}
		if rec.IsStale(threshold) {
			staleLabel = "  [STALE]"
		} else {
			staleLabel = "  [fresh]"
		}
	}
	fmt.Fprintf(c.App.Writer, "project:    %s\ntouched_at: %s%s\n",
		rec.Project, rec.TouchedAt.Format(time.RFC3339), staleLabel)
	return nil
}

func runFreshnessDelete(c *cli.Context) error {
	project := c.Args().Get(1)
	if project == "" {
		return errors.New("usage: freshness delete <project>")
	}
	m, err := freshnessManager(c)
	if err != nil {
		return err
	}
	if err := m.Delete(project); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "deleted freshness record for %q\n", project)
	return nil
}
