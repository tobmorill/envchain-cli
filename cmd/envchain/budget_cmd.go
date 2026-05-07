package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/user/envchain-cli/internal/budget"
)

const defaultBudgetDir = ".envchain/budgets"

func runBudget(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: envchain budget <set|get|delete|check> [args...]")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	m := budget.New(filepath.Join(home, defaultBudgetDir))

	switch args[0] {
	case "set":
		return runBudgetSet(m, args[1:])
	case "get":
		return runBudgetGet(m, args[1:])
	case "delete":
		return runBudgetDelete(m, args[1:])
	case "check":
		return runBudgetCheck(m, args[1:])
	default:
		return fmt.Errorf("budget: unknown subcommand %q", args[0])
	}
}

func runBudgetSet(m *budget.Manager, args []string) error {
	if len(args) != 2 {
		return errors.New("usage: envchain budget set <project> <limit-bytes>")
	}
	limit, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("budget: invalid limit %q: %w", args[1], err)
	}
	if err := m.Set(args[0], limit); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "budget set: %s → %d bytes\n", args[0], limit)
	return nil
}

func runBudgetGet(m *budget.Manager, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: envchain budget get <project>")
	}
	r, err := m.Get(args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "project:     %s\nlimit_bytes: %d\n", r.Project, r.LimitBytes)
	return nil
}

func runBudgetDelete(m *budget.Manager, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: envchain budget delete <project>")
	}
	if err := m.Delete(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "budget deleted for %s\n", args[0])
	return nil
}

func runBudgetCheck(m *budget.Manager, args []string) error {
	if len(args) != 2 {
		return errors.New("usage: envchain budget check <project> <used-bytes>")
	}
	used, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("budget: invalid byte count %q: %w", args[1], err)
	}
	if err := m.Check(args[0], used); err != nil {
		if errors.Is(err, budget.ErrExceeded) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return err
	}
	fmt.Fprintln(os.Stdout, "ok: within budget")
	return nil
}
