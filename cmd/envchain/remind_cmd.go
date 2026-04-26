package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envchain-cli/internal/remind"
	"github.com/urfave/cli/v2"
)

func defaultRemindDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "reminders")
}

func runRemind(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		return fmt.Errorf("--project is required")
	}
	m := remind.New(defaultRemindDir())
	r, err := m.Get(project)
	if err == remind.ErrNoReminder {
		fmt.Fprintf(c.App.Writer, "no reminder set for project %q\n", project)
		return nil
	}
	if err != nil {
		return err
	}
	if r.IsDue() {
		msg := r.Message
		if msg == "" {
			msg = "review or rotate your environment variables"
		}
		fmt.Fprintf(c.App.Writer, "[REMINDER] %s: %s\n", project, msg)
	} else {
		next := r.LastReset.Add(r.Interval)
		fmt.Fprintf(c.App.Writer, "next reminder for %q due %s\n", project, next.Format(time.RFC1123))
	}
	return nil
}

func runRemindSet(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		return fmt.Errorf("--project is required")
	}
	interval := c.Duration("interval")
	if interval <= 0 {
		return fmt.Errorf("--interval must be a positive duration (e.g. 168h)")
	}
	m := remind.New(defaultRemindDir())
	r := remind.Reminder{
		Project:  project,
		Interval: interval,
		Message:  c.String("message"),
	}
	if err := m.Set(r); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "reminder set for %q every %s\n", project, interval)
	return nil
}

func runRemindReset(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		return fmt.Errorf("--project is required")
	}
	m := remind.New(defaultRemindDir())
	if err := m.Reset(project); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "reminder reset for %q\n", project)
	return nil
}

func runRemindDelete(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		return fmt.Errorf("--project is required")
	}
	m := remind.New(defaultRemindDir())
	if err := m.Delete(project); err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "reminder deleted for %q\n", project)
	return nil
}
