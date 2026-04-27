package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/lint"
)

func runLint(cmd *cobra.Command, args []string) error {
	project := args[0]
	passphrase, err := resolvePassphrase(cmd)
	if err != nil {
		return err
	}

	storePath, _ := cmd.Flags().GetString("store")
	mgr, err := chain.New(storePath)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	entries, err := mgr.Load(project, passphrase)
	if err != nil {
		return fmt.Errorf("load chain %q: %w", project, err)
	}

	findings := lint.Check(entries)
	if len(findings) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "✓ no issues found in chain %q\n", project)
		return nil
	}

	hasError := false
	for _, f := range findings {
		fmt.Fprintln(cmd.OutOrStdout(), f.String())
		if f.Severity == lint.Error {
			hasError = true
		}
	}

	if hasError {
		os.Exit(1)
	}
	return nil
}
