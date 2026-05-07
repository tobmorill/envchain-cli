package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/user/envchain-cli/internal/inherit"
	"github.com/user/envchain-cli/internal/store"
)

func runInherit(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain inherit <set|get|delete|chain> [project] [parent]")
	}

	st, err := defaultStore()
	if err != nil {
		return err
	}
	mgr := inherit.New(st)

	switch args[0] {
	case "set":
		return runInheritSet(mgr, args[1:])
	case "get":
		return runInheritGet(mgr, args[1:])
	case "delete":
		return runInheritDelete(mgr, args[1:])
	case "chain":
		return runInheritChain(mgr, args[1:])
	default:
		return fmt.Errorf("inherit: unknown subcommand %q", args[0])
	}
}

func runInheritSet(mgr *inherit.Manager, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envchain inherit set <project> <parent>")
	}
	project, parent := args[0], args[1]
	if err := mgr.Set(project, parent); err != nil {
		if errors.Is(err, inherit.ErrSelfReference) {
			return fmt.Errorf("project %q cannot inherit from itself", project)
		}
		return err
	}
	fmt.Fprintf(os.Stdout, "inherit: %s → %s\n", project, parent)
	return nil
}

func runInheritGet(mgr *inherit.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain inherit get <project>")
	}
	parent, err := mgr.Get(args[0])
	if err != nil {
		return err
	}
	if parent == "" {
		fmt.Fprintf(os.Stdout, "no parent set for %q\n", args[0])
		return nil
	}
	fmt.Fprintln(os.Stdout, parent)
	return nil
}

func runInheritDelete(mgr *inherit.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain inherit delete <project>")
	}
	if err := mgr.Delete(args[0]); err != nil && !errors.Is(err, store.ErrNotFound) {
		return err
	}
	fmt.Fprintf(os.Stdout, "inherit: removed parent for %q\n", args[0])
	return nil
}

func runInheritChain(mgr *inherit.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain inherit chain <project>")
	}
	chain, err := mgr.Chain(args[0])
	if err != nil {
		if errors.Is(err, inherit.ErrCircular) {
			return fmt.Errorf("circular inheritance detected for %q", args[0])
		}
		return err
	}
	if len(chain) == 0 {
		fmt.Fprintf(os.Stdout, "%s has no ancestors\n", args[0])
		return nil
	}
	fmt.Fprintln(os.Stdout, strings.Join(chain, " → "))
	return nil
}
