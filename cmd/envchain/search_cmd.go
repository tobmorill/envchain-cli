package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/search"
	"github.com/user/envchain-cli/internal/store"
)

// runSearch implements the `envchain search` sub-command.
// Usage:
//
//	envchain search [--keys] [--passphrase <pass>] <query>
//
// Without --keys, searches chain names (projects) for the given query string.
// With --keys, decrypts each chain using the provided passphrase and searches
// within environment variable key names, printing "<project>\t<key>" matches.
func runSearch(args []string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	searchKeys := fs.Bool("keys", false, "search within environment variable keys (requires passphrase)")
	passphraseFlag := fs.String("passphrase", "", "passphrase used to decrypt chains when --keys is set")
	if err := fs.Parse(args); err != nil {
		return err
	}

	query := ""
	if fs.NArg() > 0 {
		query = fs.Arg(0)
	}

	dataDir, err := defaultDataDir()
	if err != nil {
		return fmt.Errorf("data dir: %w", err)
	}

	st, err := store.New(dataDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	cm := chain.New(st)
	sm := search.New(cm)

	names, err := listChainNames(st)
	if err != nil {
		return fmt.Errorf("list chains: %w", err)
	}

	if *searchKeys {
		return runSearchKeys(sm, query, *passphraseFlag, names)
	}

	return runSearchProjects(sm, query, names)
}

// runSearchKeys searches for env var keys matching query across all chains,
// decrypting each chain with the given passphrase.
func runSearchKeys(sm *search.Manager, query, passphraseFlag string, names []string) error {
	passphrase, err := resolvePassphrase(passphraseFlag)
	if err != nil {
		return err
	}
	results, err := sm.FindKeys(query, passphrase, names)
	if err != nil {
		return fmt.Errorf("search keys: %w", err)
	}
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no matches found")
		return nil
	}
	for _, r := range results {
		fmt.Printf("%s\t%s\n", r.Project, r.Key)
	}
	return nil
}

// runSearchProjects searches for chain names (projects) matching query.
func runSearchProjects(sm *search.Manager, query string, names []string) error {
	results := sm.FindProjects(query, names)
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no matches found")
		return nil
	}
	for _, r := range results {
		fmt.Println(r.Project)
	}
	return nil
}
