// Package project provides utilities for detecting and resolving
// the current project context based on the working directory.
package project

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	// MarkerFile is a file that can be placed in a directory to explicitly
	// associate it with an envchain project name.
	MarkerFile = ".envchain"
)

// Resolve returns the project name for the given directory.
// It first checks for a .envchain marker file, then falls back
// to using the base name of the directory.
func Resolve(dir string) (string, error) {
	name, err := readMarker(dir)
	if err != nil {
		return "", err
	}
	if name != "" {
		return name, nil
	}
	return filepath.Base(dir), nil
}

// ResolveWD returns the project name for the current working directory.
func ResolveWD() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return Resolve(wd)
}

// WriteMarker writes a .envchain marker file in dir with the given project name.
func WriteMarker(dir, name string) error {
	path := filepath.Join(dir, MarkerFile)
	return os.WriteFile(path, []byte(name+"\n"), 0o644)
}

// readMarker reads the project name from a .envchain marker file in dir.
// Returns an empty string if no marker file exists.
func readMarker(dir string) (string, error) {
	path := filepath.Join(dir, MarkerFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	name := strings.TrimSpace(string(data))
	return name, nil
}

// IsValidName reports whether name is a valid project identifier.
// Names must be non-empty and contain only alphanumerics, hyphens, and underscores.
func IsValidName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !isAlphanumeric(r) && r != '-' && r != '_' {
			return false
		}
	}
	return true
}

func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
