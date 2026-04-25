package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourorg/envchain-cli/internal/lock"
	"github.com/urfave/cli/v2"
)

const defaultLockTTL = 30 * time.Minute

func defaultLockDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "envchain", "sessions")
}

func runLock(c *cli.Context) error {
	chain := c.Args().First()
	if chain == "" {
		return fmt.Errorf("usage: envchain lock <chain>")
	}
	m := lock.NewManager(defaultLockDir())
	if err := m.Lock(chain); err != nil {
		if err == lock.ErrNotLocked {
			fmt.Fprintf(c.App.Writer, "chain %q is already locked\n", chain)
			return nil
		}
		return fmt.Errorf("lock: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "chain %q locked\n", chain)
	return nil
}

func runUnlock(c *cli.Context) error {
	chain := c.Args().First()
	if chain == "" {
		return fmt.Errorf("usage: envchain unlock [--ttl <duration>] <chain>")
	}
	ttl := c.Duration("ttl")
	if ttl <= 0 {
		ttl = defaultLockTTL
	}
	m := lock.NewManager(defaultLockDir())
	if err := m.Unlock(chain, ttl); err != nil {
		return fmt.Errorf("unlock: %w", err)
	}
	fmt.Fprintf(c.App.Writer, "chain %q unlocked for %s\n", chain, ttl)
	return nil
}

func runLockStatus(c *cli.Context) error {
	chain := c.Args().First()
	if chain == "" {
		return fmt.Errorf("usage: envchain lock-status <chain>")
	}
	m := lock.NewManager(defaultLockDir())
	ok, err := m.IsUnlocked(chain)
	if err != nil {
		return fmt.Errorf("lock-status: %w", err)
	}
	if ok {
		fmt.Fprintf(c.App.Writer, "chain %q is unlocked\n", chain)
	} else {
		fmt.Fprintf(c.App.Writer, "chain %q is locked\n", chain)
	}
	return nil
}
