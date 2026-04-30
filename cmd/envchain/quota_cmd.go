package main

import (
	"fmt"
	"strconv"

	"github.com/envchain-cli/envchain/internal/chain"
	"github.com/envchain-cli/envchain/internal/quota"
	"github.com/urfave/cli/v2"
)

func runQuota(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		return fmt.Errorf("--project is required")
	}

	maxKeys := c.Int("max-keys")
	maxBytes := c.Int("max-bytes")

	if maxKeys == 0 && maxBytes == 0 {
		// Report current usage only.
		passphrase, err := resolvePassphrase(c)
		if err != nil {
			return err
		}
		cm := chain.New(defaultStorePath(c))
		entries, err := cm.Load(project, passphrase)
		if err != nil {
			return fmt.Errorf("load chain: %w", err)
		}
		totalBytes := 0
		for _, e := range entries {
			totalBytes += len(e.Value)
		}
		fmt.Fprintf(c.App.Writer, "project:     %s\n", project)
		fmt.Fprintf(c.App.Writer, "keys:        %d\n", len(entries))
		fmt.Fprintf(c.App.Writer, "value_bytes: %d\n", totalBytes)
		return nil
	}

	// Enforce supplied limits against the current chain.
	passphrase, err := resolvePassphrase(c)
	if err != nil {
		return err
	}
	cm := chain.New(defaultStorePath(c))
	entries, err := cm.Load(project, passphrase)
	if err != nil {
		return fmt.Errorf("load chain: %w", err)
	}

	r := quota.Rule{
		MaxKeys:       maxKeys,
		MaxValueBytes: maxBytes,
	}
	violations := quota.Check(entries, r)
	if len(violations) == 0 {
		fmt.Fprintln(c.App.Writer, "ok: no quota violations")
		return nil
	}

	for _, v := range violations {
		fmt.Fprintf(c.App.ErrWriter,
			"quota violation [%s]: limit=%s actual=%s — %s\n",
			v.Field,
			strconv.Itoa(v.Limit),
			strconv.Itoa(v.Actual),
			v.Message,
		)
	}
	return fmt.Errorf("%w: %d violation(s) found", quota.ErrQuotaExceeded, len(violations))
}
