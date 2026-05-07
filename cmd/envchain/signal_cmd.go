package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/envchain/envchain-cli/internal/signal"
)

func defaultSignalDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envchain", "signals")
}

func runSignal(args []string) error {
	fs := flag.NewFlagSet("signal", flag.ContinueOnError)
	dir := fs.String("dir", defaultSignalDir(), "signal storage directory")
	if err := fs.Parse(args); err != nil {
		return err
	}

	sub := fs.Arg(0)
	rest := fs.Args()
	if len(rest) > 0 {
		rest = rest[1:]
	}

	m, err := signal.New(*dir)
	if err != nil {
		return err
	}

	switch sub {
	case "raise":
		return runSignalRaise(m, rest)
	case "get":
		return runSignalGet(m, rest)
	case "ack":
		return runSignalAck(m, rest)
	case "delete":
		return runSignalDelete(m, rest)
	default:
		return fmt.Errorf("signal: unknown subcommand %q (raise|get|ack|delete)", sub)
	}
}

func runSignalRaise(m *signal.Manager, args []string) error {
	fs := flag.NewFlagSet("signal raise", flag.ContinueOnError)
	level := fs.String("level", "info", "signal level: info|warn|error")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("usage: signal raise [-level <lvl>] <project> <message>")
	}
	project := fs.Arg(0)
	message := fs.Arg(1)
	return m.Raise(project, message, signal.Level(*level))
}

func runSignalGet(m *signal.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: signal get <project>")
	}
	rec, err := m.Get(args[0])
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Project:\t%s\n", rec.Project)
	fmt.Fprintf(w, "Level:\t%s\n", rec.Level)
	fmt.Fprintf(w, "Message:\t%s\n", rec.Message)
	fmt.Fprintf(w, "Raised:\t%s\n", rec.RaisedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Acknowledged:\t%v\n", rec.Acknowledged)
	return w.Flush()
}

func runSignalAck(m *signal.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: signal ack <project>")
	}
	if err := m.Acknowledge(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "signal acknowledged for project %q\n", args[0])
	return nil
}

func runSignalDelete(m *signal.Manager, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: signal delete <project>")
	}
	if err := m.Delete(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "signal deleted for project %q\n", args[0])
	return nil
}
