package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/user/envchain-cli/internal/embargo"
	"github.com/user/envchain-cli/internal/store"
)

func defaultEmbargoDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "embargo")
}

func runEmbargo(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain embargo <set|get|delete|check> [project] [start] [end]")
	}
	st, err := store.New(filepath.Join(defaultEmbargoDir(), "embargo.db"))
	if err != nil {
		return fmt.Errorf("embargo: open store: %w", err)
	}
	defer st.Close()
	m := embargo.New(st)

	switch args[0] {
	case "set":
		return runEmbargoSet(m, args[1:])
	case "get":
		return runEmbargoGet(m, args[1:])
	case "delete":
		return runEmbargoDelete(m, args[1:])
	case "check":
		return runEmbargoCheck(m, args[1:])
	default:
		return fmt.Errorf("embargo: unknown sub-command %q", args[0])
	}
}

// parseHHMM parses "HH:MM" into hour and minute integers.
func parseHHMM(s string) (int, int, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time %q: expected HH:MM", s)
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid hour in %q", s)
	}
	mi, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minute in %q", s)
	}
	return h, mi, nil
}

func runEmbargoSet(m *embargo.Manager, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envchain embargo set <project> <start HH:MM> <end HH:MM>")
	}
	sh, sm, err := parseHHMM(args[1])
	if err != nil {
		return err
	}
	eh, em, err := parseHHMM(args[2])
	if err != nil {
		return err
	}
	w := embargo.Window{StartHour: sh, StartMin: sm, EndHour: eh, EndMin: em}
	if err := m.Set(args[0], w); err != nil {
		return err
	}
	fmt.Printf("embargo window set for %q: %s – %s\n", args[0], args[1], args[2])
	return nil
}

func runEmbargoGet(m *embargo.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain embargo get <project>")
	}
	w, err := m.Get(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("project: %s\nwindow:  %02d:%02d – %02d:%02d UTC\n",
		args[0], w.StartHour, w.StartMin, w.EndHour, w.EndMin)
	return nil
}

func runEmbargoDelete(m *embargo.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain embargo delete <project>")
	}
	if err := m.Delete(args[0]); err != nil {
		return err
	}
	fmt.Printf("embargo window removed for %q\n", args[0])
	return nil
}

func runEmbargoCheck(m *embargo.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envchain embargo check <project>")
	}
	if err := m.Check(args[0]); err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Printf("access permitted for %q\n", args[0])
	return nil
}
