package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/envchain/envchain-cli/internal/store"
	"github.com/envchain/envchain-cli/internal/template"
	"github.com/urfave/cli/v2"
)

func runTemplateCmd(c *cli.Context) error {
	subCmd := c.Args().First()
	switch subCmd {
	case "save":
		return runTemplateSave(c)
	case "load":
		return runTemplateLoad(c)
	case "delete":
		return runTemplateDelete(c)
	default:
		return fmt.Errorf("unknown template subcommand %q; use save|load|delete", subCmd)
	}
}

func templateManager(c *cli.Context) (*template.Manager, error) {
	dbPath := c.String("store")
	st, err := store.New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open store: %w", err)
	}
	return template.New(st), nil
}

func runTemplateSave(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 2 {
		return errors.New("usage: template save <name> <KEY1> [KEY2 ...]")
	}
	name := args[0]
	keys := args[1:]

	m, err := templateManager(c)
	if err != nil {
		return err
	}
	tmpl := template.Template{Name: name, Keys: keys}
	if err := m.Save(tmpl); err != nil {
		if errors.Is(err, template.ErrInvalidName) {
			return fmt.Errorf("invalid template name %q: only letters, digits, hyphens and underscores allowed", name)
		}
		return fmt.Errorf("save template: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "template %q saved with keys: %s\n", name, strings.Join(keys, ", "))
	return nil
}

func runTemplateLoad(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: template load <name>")
	}
	name := args[0]

	m, err := templateManager(c)
	if err != nil {
		return err
	}
	tmpl, err := m.Load(name)
	if err != nil {
		if errors.Is(err, template.ErrNotFound) {
			return fmt.Errorf("template %q not found", name)
		}
		return fmt.Errorf("load template: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "name: %s\nkeys:\n", tmpl.Name)
	for _, k := range tmpl.Keys {
		fmt.Fprintf(c.App.Writer, "  - %s\n", k)
	}
	return nil
}

func runTemplateDelete(c *cli.Context) error {
	args := c.Args().Tail()
	if len(args) < 1 {
		return errors.New("usage: template delete <name>")
	}
	name := args[0]

	m, err := templateManager(c)
	if err != nil {
		return err
	}
	if err := m.Delete(name); err != nil {
		if errors.Is(err, template.ErrNotFound) {
			return fmt.Errorf("template %q not found", name)
		}
		return fmt.Errorf("delete template: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "template %q deleted\n", name)
	return nil
}
