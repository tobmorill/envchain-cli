package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envchain/envchain-cli/internal/expiry"
	"github.com/urfave/cli/v2"
)

func defaultExpiryDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "expiry")
}

func runExpiry(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "set":
		return runExpirySet(c)
	case "get":
		return runExpiryGet(c)
	case "delete":
		return runExpiryDelete(c)
	default:
		return fmt.Errorf("expiry: unknown subcommand %q (use set|get|delete)", subcmd)
	}
}

func runExpirySet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 2 {
		return errors.New("usage: expiry set <project> <duration|RFC3339> [note]")
	}
	project := args[0]
	var expiresAt time.Time
	if d, err := time.ParseDuration(args[1]); err == nil {
		expiresAt = time.Now().Add(d)
	} else if t, err := time.Parse(time.RFC3339, args[1]); err == nil {
		expiresAt = t
	} else {
		return fmt.Errorf("expiry: cannot parse %q as duration or RFC3339 time", args[1])
	}
	note := ""
	if len(args) >= 3 {
		note = args[2]
	}
	m, err := expiry.New(defaultExpiryDir())
	if err != nil {
		return err
	}
	if err := m.Set(project, expiresAt, note); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "expiry set for %q until %s\n", project, expiresAt.Format(time.RFC3339))
	return nil
}

func runExpiryGet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: expiry get <project>")
	}
	m, err := expiry.New(defaultExpiryDir())
	if err != nil {
		return err
	}
	rec, err := m.Get(args[0])
	if errors.Is(err, expiry.ErrNotFound) {
		return fmt.Errorf("no expiry record for %q", args[0])
	}
	if err != nil {
		return err
	}
	status := "active"
	if rec.IsExpired() {
		status = "EXPIRED"
	}
	fmt.Fprintf(c.App.Writer, "project:    %s\nexpires_at: %s\nstatus:     %s\n",
		rec.Project, rec.ExpiresAt.Format(time.RFC3339), status)
	if rec.Note != "" {
		fmt.Fprintf(c.App.Writer, "note:       %s\n", rec.Note)
	}
	return nil
}

func runExpiryDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: expiry delete <project>")
	}
	m, err := expiry.New(defaultExpiryDir())
	if err != nil {
		return err
	}
	if err := m.Delete(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "expiry record for %q deleted\n", args[0])
	return nil
}
