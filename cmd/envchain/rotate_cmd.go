package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/envchain-cli/internal/passphrase"
	"github.com/your-org/envchain-cli/internal/rotate"
	"github.com/your-org/envchain-cli/internal/store"
)

// runRotate handles the "rotate" sub-command.
// Usage: envchain rotate [--all] <chain-name>
func runRotate(args []string) error {
	fs := flag.NewFlagSet("rotate", flag.ContinueOnError)
	all := fs.Bool("all", false, "rotate all chains that share the same old passphrase")
	fs.SetOutput(os.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if !*all && fs.NArg() < 1 {
		return fmt.Errorf("rotate: chain name required (or use --all)")
	}

	oldPass, err := passphrase.Prompt("Current passphrase: ")
	if err != nil {
		return fmt.Errorf("rotate: read current passphrase: %w", err)
	}

	newPass, err := passphrase.PromptConfirm("New passphrase: ", "Confirm new passphrase: ")
	if err != nil {
		return fmt.Errorf("rotate: read new passphrase: %w", err)
	}

	st, err := store.New(defaultStorePath())
	if err != nil {
		return fmt.Errorf("rotate: open store: %w", err)
	}

	rm := rotate.New(st)

	if *all {
		names, err := listChainNames(st)
		if err != nil {
			return fmt.Errorf("rotate: list chains: %w", err)
		}
		if len(names) == 0 {
			fmt.Fprintln(os.Stderr, "rotate: no chains found")
			return nil
		}
		if err := rm.RotateAll(names, oldPass, newPass); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "rotated %d chain(s)\n", len(names))
		return nil
	}

	name := fs.Arg(0)
	if err := rm.Rotate(name, oldPass, newPass); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "rotated chain %q\n", name)
	return nil
}

// listChainNames returns all chain names currently persisted in the store.
// It relies on the store exposing a Keys method that returns raw bucket keys.
func listChainNames(st *store.Store) ([]string, error) {
	keys, err := st.Keys()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
