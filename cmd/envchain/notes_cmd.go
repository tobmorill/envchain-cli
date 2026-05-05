package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain-cli/internal/notes"
	"github.com/envchain/envchain-cli/internal/store"
	"github.com/urfave/cli/v2"
)

func runNotes(c *cli.Context) error {
	subcmd := c.Args().First()
	switch subcmd {
	case "set":
		return runNotesSet(c)
	case "get":
		return runNotesGet(c)
	case "delete":
		return runNotesDelete(c)
	default:
		return fmt.Errorf("notes: unknown subcommand %q; use set, get, or delete", subcmd)
	}
}

func notesManager(c *cli.Context) (*notes.Manager, error) {
	dir := c.String("store")
	if dir == "" {
		dir = defaultStorePath()
	}
	s, err := store.New(dir)
	if err != nil {
		return nil, fmt.Errorf("notes: open store: %w", err)
	}
	return notes.New(s), nil
}

func runNotesSet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: envchain notes set <project> [body...]")
	}
	project := args[0]
	body := strings.Join(args[1:], " ")
	if body == "" {
		body = c.String("message")
	}
	if body == "" {
		return errors.New("notes set: body must not be empty; pass text or use --message")
	}
	pass, err := resolvePassphrase(c)
	if err != nil {
		return err
	}
	m, err := notesManager(c)
	if err != nil {
		return err
	}
	if err := m.Set(project, body, pass); err != nil {
		return fmt.Errorf("notes set: %w", err)
	}
	fmt.Fprintf(os.Stdout, "note saved for project %q\n", project)
	return nil
}

func runNotesGet(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: envchain notes get <project>")
	}
	project := args[0]
	pass, err := resolvePassphrase(c)
	if err != nil {
		return err
	}
	m, err := notesManager(c)
	if err != nil {
		return err
	}
	n, err := m.Get(project, pass)
	if err != nil {
		return fmt.Errorf("notes get: %w", err)
	}
	fmt.Fprintf(os.Stdout, "# %s  (updated %s)\n%s\n",
		n.Project, n.UpdatedAt.Format("2006-01-02 15:04:05 UTC"), n.Body)
	return nil
}

func runNotesDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: envchain notes delete <project>")
	}
	project := args[0]
	m, err := notesManager(c)
	if err != nil {
		return err
	}
	if err := m.Delete(project); err != nil {
		return fmt.Errorf("notes delete: %w", err)
	}
	fmt.Fprintf(os.Stdout, "note deleted for project %q\n", project)
	return nil
}
