package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/your-org/envchain-cli/internal/blame"
	"github.com/your-org/envchain-cli/internal/store"
)

func defaultBlameDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "blame.db")
}

// runBlame handles the "blame" sub-command.
//
//	envchain blame <project>            – show last modifier
//	envchain blame touch <project>      – record current user as last modifier
//	envchain blame delete <project>     – remove blame record
func runBlame(args []string) error {
	fs := flag.NewFlagSet("blame", flag.ContinueOnError)
	note := fs.String("note", "", "optional note to attach when touching")
	dbPath := fs.String("db", defaultBlameDir(), "path to blame store")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		return fmt.Errorf("blame: expected sub-command or project name")
	}

	st, err := store.New(*dbPath)
	if err != nil {
		return fmt.Errorf("blame: open store: %w", err)
	}
	m := blame.New(st)

	switch remaining[0] {
	case "touch":
		if len(remaining) < 2 {
			return fmt.Errorf("blame touch: project name required")
		}
		if err := m.Touch(remaining[1], *note); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "blame: touched %q\n", remaining[1])
		return nil

	case "delete":
		if len(remaining) < 2 {
			return fmt.Errorf("blame delete: project name required")
		}
		if err := m.Delete(remaining[1]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "blame: deleted record for %q\n", remaining[1])
		return nil

	default:
		// treat the argument as a project name and show its record
		project := remaining[0]
		rec, err := m.Get(project)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "project:  %s\n", rec.Project)
		fmt.Fprintf(os.Stdout, "user:     %s\n", rec.User)
		fmt.Fprintf(os.Stdout, "hostname: %s\n", rec.Hostname)
		fmt.Fprintf(os.Stdout, "changed:  %s\n", rec.ChangedAt.Local().Format("2006-01-02 15:04:05"))
		if rec.Note != "" {
			fmt.Fprintf(os.Stdout, "note:     %s\n", rec.Note)
		}
		return nil
	}
}
