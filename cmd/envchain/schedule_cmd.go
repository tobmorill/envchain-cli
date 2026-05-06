package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envchain/envchain-cli/internal/schedule"
	"github.com/urfave/cli/v2"
)

func defaultScheduleDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "envchain", "schedules")
}

func runSchedule(c *cli.Context) error {
	sub := c.Args().First()
	switch sub {
	case "set":
		return runScheduleSet(c)
	case "get":
		return runScheduleGet(c)
	case "delete":
		return runScheduleDelete(c)
	case "due":
		return runScheduleDue(c)
	default:
		return fmt.Errorf("schedule: unknown subcommand %q — use set|get|delete|due", sub)
	}
}

func runScheduleSet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 2 {
		return errors.New("usage: schedule set <project> <interval>")
	}
	project := args[0]
	d, err := time.ParseDuration(args[1])
	if err != nil {
		return fmt.Errorf("schedule: invalid duration %q: %w", args[1], err)
	}
	mgr, err := schedule.New(defaultScheduleDir())
	if err != nil {
		return err
	}
	if err := mgr.Set(project, d); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "schedule set: %s every %s\n", project, d)
	return nil
}

func runScheduleGet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: schedule get <project>")
	}
	mgr, err := schedule.New(defaultScheduleDir())
	if err != nil {
		return err
	}
	rec, err := mgr.Get(args[0])
	if err != nil {
		if errors.Is(err, schedule.ErrNotFound) {
			return fmt.Errorf("no schedule for project %q", args[0])
		}
		return err
	}
	due := "no"
	if rec.IsDue() {
		due = "YES"
	}
	fmt.Fprintf(c.App.Writer, "project:  %s\ninterval: %s\nupdated:  %s\ndue:      %s\n",
		rec.Project, rec.Interval, rec.UpdatedAt.Format(time.RFC3339), due)
	return nil
}

func runScheduleDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: schedule delete <project>")
	}
	mgr, err := schedule.New(defaultScheduleDir())
	if err != nil {
		return err
	}
	if err := mgr.Delete(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "schedule deleted: %s\n", args[0])
	return nil
}

func runScheduleDue(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: schedule due <project>")
	}
	mgr, err := schedule.New(defaultScheduleDir())
	if err != nil {
		return err
	}
	rec, err := mgr.Get(args[0])
	if err != nil {
		if errors.Is(err, schedule.ErrNotFound) {
			fmt.Fprintln(c.App.Writer, "no schedule configured")
			return nil
		}
		return err
	}
	if rec.IsDue() {
		fmt.Fprintln(c.App.Writer, "due")
	} else {
		fmt.Fprintln(c.App.Writer, "not due")
	}
	return nil
}
