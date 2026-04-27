package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/merge"
	"github.com/user/envchain-cli/internal/shell"
	"github.com/urfave/cli/v2"
)

// runMerge loads multiple chains and emits a merged export script.
func runMerge(c *cli.Context) error {
	if c.NArg() < 2 {
		return errors.New("merge requires at least two chain names")
	}

	passphrase, err := resolvePassphrase(c)
	if err != nil {
		return err
	}

	strategyFlag := c.String("strategy")
	var strategy merge.Strategy
	switch strings.ToLower(strategyFlag) {
	case "first", "":
		strategy = merge.StrategyFirst
	case "last":
		strategy = merge.StrategyLast
	case "error":
		strategy = merge.StrategyError
	default:
		return fmt.Errorf("unknown strategy %q: use first, last, or error", strategyFlag)
	}

	storePath := c.String("store")
	cm := chain.New(storePath)

	var named []merge.NamedChain
	for _, name := range c.Args().Slice() {
		entries, err := cm.Load(name, passphrase)
		if err != nil {
			return fmt.Errorf("loading chain %q: %w", name, err)
		}
		nc, err := merge.NewNamedChain(name, entries)
		if err != nil {
			return err
		}
		named = append(named, nc)
	}

	result, err := merge.Merge(named, strategy)
	if err != nil {
		return err
	}

	sh, err := resolveShell(c)
	if err != nil {
		return err
	}

	script, err := shell.ExportScript(sh, result.Entries)
	if err != nil {
		return err
	}

	if c.Bool("show-origins") {
		for _, e := range result.Entries {
			fmt.Fprintf(os.Stderr, "# %s <- %s\n", e.Key, result.Origins[e.Key])
		}
	}

	fmt.Print(script)
	return nil
}
