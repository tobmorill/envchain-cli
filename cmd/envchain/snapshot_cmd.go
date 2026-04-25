package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/snapshot"
	"github.com/user/envchain-cli/internal/store"
)

func runSnapshot(args []string) error {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	action := fs.String("action", "save", "Action: save|get|delete")
	label := fs.String("label", "", "Snapshot label (required)")
	project := fs.String("project", "", "Project name (required)")
	dbPath := fs.String("db", defaultStorePath(), "Path to store database")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *project == "" {
		return fmt.Errorf("snapshot: --project is required")
	}
	if *label == "" {
		return fmt.Errorf("snapshot: --label is required")
	}

	st, err := store.New(*dbPath)
	if err != nil {
		return fmt.Errorf("snapshot: open store: %w", err)
	}
	defer st.Close()

	cm := chain.New(st)
	mgr := snapshot.New(st, cm)

	switch *action {
	case "save":
		passphrase, err := resolvePassphrase()
		if err != nil {
			return err
		}
		if err := mgr.Save(*project, *label, passphrase); err != nil {
			return err
		}
		fmt.Printf("Snapshot %q saved for project %q.\n", *label, *project)

	case "get":
		snap, err := mgr.Get(*project, *label)
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Project:\t%s\n", snap.Project)
		fmt.Fprintf(w, "Label:\t%s\n", snap.Label)
		fmt.Fprintf(w, "Created:\t%s\n", snap.CreatedAt.Format(time.RFC3339))
		fmt.Fprintf(w, "Entries:\t%d\n", len(snap.Entries))
		w.Flush()
		for _, e := range snap.Entries {
			fmt.Println(" ", e)
		}

	case "delete":
		if err := mgr.Delete(*project, *label); err != nil {
			return err
		}
		fmt.Printf("Snapshot %q deleted for project %q.\n", *label, *project)

	default:
		return fmt.Errorf("snapshot: unknown action %q (save|get|delete)", *action)
	}
	return nil
}

func defaultStorePath() string {
	home, _ := os.UserHomeDir()
	return home + "/.envchain/store.db"
}
