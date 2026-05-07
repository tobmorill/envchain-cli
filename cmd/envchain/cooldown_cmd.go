package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/user/envchain-cli/internal/cooldown"
	"github.com/user/envchain-cli/internal/store"
)

const defaultCooldownDir = ".envchain/cooldown"

func runCooldown(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: envchain cooldown <set|get|delete|status> [project] [minutes]")
	}

	dir := defaultCooldownDir
	if d := os.Getenv("ENVCHAIN_COOLDOWN_DIR"); d != "" {
		dir = d
	}
	s, err := store.New(dir)
	if err != nil {
		return fmt.Errorf("cooldown: open store: %w", err)
	}
	m := cooldown.New(s)

	switch args[0] {
	case "set":
		return runCooldownSet(m, args[1:])
	case "get":
		return runCooldownGet(m, args[1:])
	case "delete":
		return runCooldownDelete(m, args[1:])
	case "status":
		return runCooldownStatus(m, args[1:])
	default:
		return fmt.Errorf("cooldown: unknown subcommand %q", args[0])
	}
}

func runCooldownSet(m *cooldown.Manager, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envchain cooldown set <project> <minutes>")
	}
	project := args[0]
	mins, err := strconv.Atoi(args[1])
	if err != nil || mins <= 0 {
		return fmt.Errorf("cooldown: minutes must be a positive integer")
	}
	if err := m.Set(project, time.Duration(mins)*time.Minute); err != nil {
		return fmt.Errorf("cooldown: set: %w", err)
	}
	fmt.Printf("cooldown set for %q: %d minute(s)\n", project, mins)
	return nil
}

func runCooldownGet(m *cooldown.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain cooldown get <project>")
	}
	rec, err := m.Get(args[0])
	if errors.Is(err, store.ErrNotFound) {
		fmt.Printf("no cooldown record for %q\n", args[0])
		return nil
	}
	if err != nil {
		return fmt.Errorf("cooldown: get: %w", err)
	}
	status := "expired"
	if rec.IsActive() {
		status = "active"
	}
	fmt.Printf("project:    %s\nstarted:    %s\nduration:   %s\nexpires:    %s\nstatus:     %s\n",
		rec.Project,
		rec.StartedAt.Format(time.RFC3339),
		rec.Duration,
		rec.ExpiresAt().Format(time.RFC3339),
		status,
	)
	return nil
}

func runCooldownDelete(m *cooldown.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain cooldown delete <project>")
	}
	if err := m.Delete(args[0]); err != nil {
		return fmt.Errorf("cooldown: delete: %w", err)
	}
	fmt.Printf("cooldown record for %q removed\n", args[0])
	return nil
}

func runCooldownStatus(m *cooldown.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain cooldown status <project>")
	}
	active, err := m.IsActive(args[0])
	if err != nil {
		return fmt.Errorf("cooldown: status: %w", err)
	}
	if active {
		fmt.Printf("cooldown for %q is ACTIVE\n", args[0])
	} else {
		fmt.Printf("cooldown for %q is not active\n", args[0])
	}
	return nil
}
