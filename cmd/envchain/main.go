// Command envchain is a CLI tool for managing per-project environment variable
// sets with encrypted local storage and shell integration.
//
// Usage:
//
//	envchain <command> [flags]
//
// Commands:
//
//	set    <chain> <KEY=VALUE>...  Add or update variables in a chain
//	get    <chain>                 Print export script for a chain
//	list                          List all stored chains
//	delete <chain>                Delete a chain
//	exec   <chain> -- <cmd>...    Run a command with chain variables injected
package main

import (
	"fmt"
	"os"

	"github.com/envchain-cli/envchain/internal/chain"
	"github.com/envchain-cli/envchain/internal/shell"
	"github.com/envchain-cli/envchain/internal/store"

	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Resolve the default store directory (~/.envchain).
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolving home directory: %w", err)
	}
	storeDir := homeDir + "/.envchain"

	s, err := store.New(storeDir)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	mgr := chain.New(s)

	root := &cobra.Command{
		Use:   "envchain",
		Short: "Manage per-project environment variable sets with encrypted local storage",
		SilenceUsage: true,
	}

	var passphrase string
	root.PersistentFlags().StringVarP(&passphrase, "passphrase", "p", "",
		"Passphrase for encrypting/decrypting the chain (prompted if empty)")

	root.AddCommand(
		newSetCmd(mgr, &passphrase),
		newGetCmd(mgr, &passphrase),
		newListCmd(mgr),
		newDeleteCmd(mgr),
		newExecCmd(mgr, &passphrase),
	)

	return root.Execute()
}

// resolvePassphrase returns the passphrase from the flag, the ENVCHAIN_PASS
// environment variable, or prompts the user interactively.
func resolvePassphrase(flagValue *string, chainName string) (string, error) {
	if flagValue != nil && *flagValue != "" {
		return *flagValue, nil
	}
	if v := os.Getenv("ENVCHAIN_PASS"); v != "" {
		return v, nil
	}
	return promptPassphrase(fmt.Sprintf("Passphrase for chain %q: ", chainName))
}

// resolveShell returns the shell name from the --shell flag or auto-detection.
func resolveShell(flagValue string) (string, error) {
	if flagValue != "" {
		if !shell.IsSupported(flagValue) {
			return "", fmt.Errorf("unsupported shell %q; supported: %v", flagValue, shell.SupportedShells)
		}
		return flagValue, nil
	}
	detected, err := shell.Detect()
	if err != nil {
		return "", fmt.Errorf("detecting shell (use --shell to specify): %w", err)
	}
	return detected, nil
}
