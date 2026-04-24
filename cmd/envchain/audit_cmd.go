package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/yourorg/envchain-cli/internal/audit"
)

// runAudit implements the "envchain audit" sub-command.
// It prints the audit log in a human-readable tabular format.
func runAudit(args []string) error {
	fs := flag.NewFlagSet("audit", flag.ContinueOnError)
	limit := fs.Int("n", 0, "show only the last N entries (0 = all)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	logPath, err := defaultAuditPath()
	if err != nil {
		return err
	}

	logger, err := audit.NewLogger(logPath)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	events, err := logger.ReadAll()
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	if len(events) == 0 {
		fmt.Println("No audit events recorded.")
		return nil
	}

	if *limit > 0 && *limit < len(events) {
		events = events[len(events)-*limit:]
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tKIND\tPROJECT\tMESSAGE")
	for _, e := range events {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Kind,
			e.Project,
			e.Message,
		)
	}
	return w.Flush()
}

// defaultAuditPath returns the platform-appropriate path for the audit log.
func defaultAuditPath() (string, error) {
	dataDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine config dir: %w", err)
	}
	return filepath.Join(dataDir, "envchain", "audit.jsonl"), nil
}
