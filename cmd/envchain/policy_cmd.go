package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-org/envchain-cli/internal/policy"
	"github.com/urfave/cli/v2"
)

func defaultPolicyDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "policies")
}

func runPolicy(c *cli.Context) error {
	subcmd := c.Args().First()
	project := c.Args().Get(1)
	m := policy.New(defaultPolicyDir())

	switch subcmd {
	case "set":
		if project == "" {
			return fmt.Errorf("usage: policy set <project> [--allow k1,k2] [--deny k3] [--pattern regex]")
		}
		r := policy.Rule{
			AllowPattern: c.String("pattern"),
		}
		if raw := c.String("allow"); raw != "" {
			r.AllowKeys = splitCSV(raw)
		}
		if raw := c.String("deny"); raw != "" {
			r.DenyKeys = splitCSV(raw)
		}
		if err := m.Set(project, r); err != nil {
			return err
		}
		fmt.Fprintf(c.App.Writer, "policy set for project %q\n", project)

	case "get":
		if project == "" {
			return fmt.Errorf("usage: policy get <project>")
		}
		r, err := m.Get(project)
		if errors.Is(err, policy.ErrNotFound) {
			return fmt.Errorf("no policy found for project %q", project)
		}
		if err != nil {
			return err
		}
		enc := json.NewEncoder(c.App.Writer)
		enc.SetIndent("", "  ")
		return enc.Encode(r)

	case "delete":
		if project == "" {
			return fmt.Errorf("usage: policy delete <project>")
		}
		err := m.Delete(project)
		if errors.Is(err, policy.ErrNotFound) {
			return fmt.Errorf("no policy found for project %q", project)
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(c.App.Writer, "policy deleted for project %q\n", project)

	case "check":
		// policy check <project> <KEY>
		key := c.Args().Get(2)
		if project == "" || key == "" {
			return fmt.Errorf("usage: policy check <project> <KEY>")
		}
		r, err := m.Get(project)
		if errors.Is(err, policy.ErrNotFound) {
			fmt.Fprintln(c.App.Writer, "allowed (no policy)")
			return nil
		}
		if err != nil {
			return err
		}
		ok, err := policy.Allowed(r, key)
		if err != nil {
			return err
		}
		if ok {
			fmt.Fprintln(c.App.Writer, "allowed")
		} else {
			fmt.Fprintln(c.App.Writer, "denied")
		}

	default:
		return fmt.Errorf("unknown subcommand %q; use set|get|delete|check", subcmd)
	}
	return nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
