package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envchain-cli/internal/complexity"
	"github.com/user/envchain-cli/internal/project"
)

func runComplexity(args []string) error {
	fs := flag.NewFlagSet("complexity", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: envchain complexity [--project <name>] [--passphrase-env <var>]")
		fs.PrintDefaults()
	}

	var (
		projectFlag      = fs.String("project", "", "project name (defaults to current directory)")
		passphraseEnvVar = fs.String("passphrase-env", "", "env var holding the passphrase")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	proj := *projectFlag
	if proj == "" {
		var err error
		proj, err = project.ResolveWD()
		if err != nil {
			return fmt.Errorf("resolve project: %w", err)
		}
	}

	passphrase, err := resolvePassphrase(*passphraseEnvVar)
	if err != nil {
		return err
	}

	entries, err := loadChainEntries(proj, passphrase)
	if err != nil {
		return fmt.Errorf("load chain: %w", err)
	}

	result := complexity.Evaluate(entries)

	fmt.Printf("project : %s\n", proj)
	fmt.Printf("score   : %.2f\n", result.Score)
	fmt.Printf("level   : %s\n", result.Level)

	if len(result.Findings) == 0 {
		fmt.Println("findings: none")
		return nil
	}

	fmt.Println("findings:")
	for _, f := range result.Findings {
		fmt.Printf("  %-30s %s (penalty %.1f)\n", f.Key, f.Reason, f.Penalty)
	}
	return nil
}
