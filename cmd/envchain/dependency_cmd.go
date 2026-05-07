package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/envchain/envchain-cli/internal/dependency"
)

func defaultDependencyDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "dependency")
}

func runDependency(args []string) error {
	fs := flag.NewFlagSet("dependency", flag.ContinueOnError)
	dir := fs.String("dir", defaultDependencyDir(), "dependency store directory")
	if err := fs.Parse(args); err != nil {
		return err
	}

	sub := fs.Arg(0)
	rest := fs.Args()
	if len(rest) > 0 {
		rest = rest[1:]
	}

	m, err := dependency.New(*dir)
	if err != nil {
		return err
	}

	switch sub {
	case "set":
		if len(rest) < 2 {
			return fmt.Errorf("usage: dependency set <project> <dep>[,<dep>...]")
		}
		project := rest[0]
		deps := splitCSVDeps(rest[1])
		if err := m.Set(project, deps); err != nil {
			if errors.Is(err, dependency.ErrSelfDependency) {
				return fmt.Errorf("project cannot depend on itself")
			}
			return err
		}
		fmt.Fprintf(os.Stdout, "dependencies for %q updated\n", project)
		return nil

	case "get":
		if len(rest) < 1 {
			return fmt.Errorf("usage: dependency get <project>")
		}
		deps, err := m.Get(rest[0])
		if err != nil {
			return err
		}
		if len(deps) == 0 {
			fmt.Fprintf(os.Stdout, "no dependencies recorded for %q\n", rest[0])
			return nil
		}
		for _, d := range deps {
			fmt.Fprintln(os.Stdout, d)
		}
		return nil

	case "delete":
		if len(rest) < 1 {
			return fmt.Errorf("usage: dependency delete <project>")
		}
		if err := m.Delete(rest[0]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "dependencies for %q deleted\n", rest[0])
		return nil

	default:
		return fmt.Errorf("dependency: unknown subcommand %q — use set, get, delete", sub)
	}
}

func splitCSVDeps(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
