// Package shell provides utilities for generating shell integration
// scripts that load envchain environment variables into the current shell.
package shell

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envchain-cli/internal/env"
)

// ShellType represents a supported shell.
type ShellType string

const (
	Bash ShellType = "bash"
	Zsh  ShellType = "zsh"
	Fish ShellType = "fish"
)

// SupportedShells returns all supported shell types.
func SupportedShells() []ShellType {
	return []ShellType{Bash, Zsh, Fish}
}

// IsSupported returns true if the given shell string is supported.
func IsSupported(s string) bool {
	for _, sh := range SupportedShells() {
		if ShellType(strings.ToLower(s)) == sh {
			return true
		}
	}
	return false
}

// ExportScript generates a shell-specific export script for the given entries.
func ExportScript(shell ShellType, entries []env.Entry, w io.Writer) error {
	switch shell {
	case Bash, Zsh:
		return exportPosix(entries, w)
	case Fish:
		return exportFish(entries, w)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
}

func exportPosix(entries []env.Entry, w io.Writer) error {
	for _, e := range entries {
		_, err := fmt.Fprintf(w, "export %s=%q\n", e.Key, e.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func exportFish(entries []env.Entry, w io.Writer) error {
	for _, e := range entries {
		_, err := fmt.Fprintf(w, "set -x %s %q\n", e.Key, e.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
