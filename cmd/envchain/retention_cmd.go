package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/envchain-cli/internal/retention"
	"github.com/envchain-cli/internal/store"
	"github.com/urfave/cli/v2"
)

const defaultRetentionDir = ".envchain/retention"

func retentionManager(c *cli.Context) (*retention.Manager, error) {
	dir := c.String("retention-dir")
	if dir == "" {
		dir = defaultRetentionDir
	}
	st, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("retention store: %w", err)
	}
	return retention.New(st), nil
}

func runRetention(c *cli.Context) error {
	switch c.Args().First() {
	case "set":
		return runRetentionSet(c)
	case "get":
		return runRetentionGet(c)
	case "delete":
		return runRetentionDelete(c)
	default:
		return cli.ShowSubcommandHelp(c)
	}
}

func runRetentionSet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: retention set <project> [--max-age <duration>] [--max-versions <n>]")
	}
	project := args[0]
	m, err := retentionManager(c)
	if err != nil {
		return err
	}
	var maxAge time.Duration
	if s := c.String("max-age"); s != "" {
		maxAge, err = time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("invalid --max-age: %w", err)
		}
	}
	maxVersions := c.Int("max-versions")
	p := retention.Policy{
		Project:     project,
		MaxAge:      maxAge,
		MaxVersions: maxVersions,
	}
	if err := m.Set(p); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "retention policy set for %q\n", project)
	return nil
}

func runRetentionGet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: retention get <project>")
	}
	m, err := retentionManager(c)
	if err != nil {
		return err
	}
	p, err := m.Get(args[0])
	if errors.Is(err, retention.ErrNotFound) {
		fmt.Fprintln(os.Stderr, "no retention policy found")
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "project:      %s\n", p.Project)
	fmt.Fprintf(c.App.Writer, "max_age:      %s\n", p.MaxAge)
	fmt.Fprintf(c.App.Writer, "max_versions: %s\n", strconv.Itoa(p.MaxVersions))
	fmt.Fprintf(c.App.Writer, "updated_at:   %s\n", p.UpdatedAt.Format(time.RFC3339))
	return nil
}

func runRetentionDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: retention delete <project>")
	}
	m, err := retentionManager(c)
	if err != nil {
		return err
	}
	if err := m.Delete(args[0]); errors.Is(err, retention.ErrNotFound) {
		fmt.Fprintln(os.Stderr, "no retention policy found")
		return nil
	} else if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "retention policy deleted for %q\n", args[0])
	return nil
}
