// Package config provides persistent user-level configuration for the
// envchain CLI tool.
//
// Configuration is stored as a JSON file, by default at
// ~/.envchain/config.json. The file is created with restricted permissions
// (0600) so that only the owning user can read it.
//
// Typical usage:
//
//	mgr, err := config.NewManager("") // uses ~/.envchain
//	if err != nil { ... }
//
//	cfg, err := mgr.Load()
//	if err != nil { ... }
//
//	cfg.DefaultShell = "zsh"
//	if err := mgr.Save(cfg); err != nil { ... }
package config
