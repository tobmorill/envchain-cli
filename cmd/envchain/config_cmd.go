package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envchain-cli/internal/config"
	"github.com/yourorg/envchain-cli/internal/shell"
)

// runConfig handles the "config" sub-command, allowing users to get or set
// persistent configuration values.
//
// Usage:
//
//	envchain config get <key>
//	envchain config set <key> <value>
func runConfig(args []string) error {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: envchain config <get|set> <key> [value]")
		fmt.Fprintln(os.Stderr, "keys: default_shell, store_path, passphrase_hint")
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		fs.Usage()
		return fmt.Errorf("insufficient arguments")
	}

	mgr, err := config.NewManager("")
	if err != nil {
		return fmt.Errorf("config manager: %w", err)
	}
	cfg, err := mgr.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	action, key := fs.Arg(0), fs.Arg(1)

	switch action {
	case "get":
		val, err := getField(&cfg, key)
		if err != nil {
			return err
		}
		fmt.Println(val)
	case "set":
		if fs.NArg() < 3 {
			return fmt.Errorf("set requires a value argument")
		}
		value := fs.Arg(2)
		if err := setField(&cfg, key, value); err != nil {
			return err
		}
		if err := mgr.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		fmt.Printf("config: %s set to %q\n", key, value)
	default:
		fs.Usage()
		return fmt.Errorf("unknown action %q", action)
	}
	return nil
}

func getField(cfg *config.Config, key string) (string, error) {
	switch key {
	case "default_shell":
		return cfg.DefaultShell, nil
	case "store_path":
		return cfg.StorePath, nil
	case "passphrase_hint":
		return cfg.PassphraseHint, nil
	}
	return "", fmt.Errorf("unknown config key %q", key)
}

func setField(cfg *config.Config, key, value string) error {
	switch key {
	case "default_shell":
		if !shell.IsSupported(value) {
			return fmt.Errorf("unsupported shell %q; supported: %v", value, shell.SupportedShells)
		}
		cfg.DefaultShell = value
	case "store_path":
		cfg.StorePath = value
	case "passphrase_hint":
		cfg.PassphraseHint = value
	default:
		return fmt.Errorf("unknown config key %q", key)
	}
	return nil
}
