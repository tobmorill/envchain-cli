// Package passphrase provides utilities for securely prompting
// and validating passphrases from the user via a terminal.
package passphrase

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ErrMismatch is returned when two passphrase entries do not match.
var ErrMismatch = errors.New("passphrases do not match")

// ErrEmpty is returned when an empty passphrase is provided.
var ErrEmpty = errors.New("passphrase must not be empty")

// Prompt reads a passphrase from the terminal without echoing input.
// The prompt string is written to stderr before reading.
func Prompt(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr) // newline after hidden input
	if err != nil {
		return "", fmt.Errorf("reading passphrase: %w", err)
	}
	p := strings.TrimRight(string(bytes), "\r\n")
	if p == "" {
		return "", ErrEmpty
	}
	return p, nil
}

// PromptConfirm prompts for a passphrase twice and returns it only if
// both entries match. Intended for use when creating a new chain.
func PromptConfirm(prompt, confirmPrompt string) (string, error) {
	p1, err := Prompt(prompt)
	if err != nil {
		return "", err
	}
	p2, err := Prompt(confirmPrompt)
	if err != nil {
		return "", err
	}
	if p1 != p2 {
		return "", ErrMismatch
	}
	return p1, nil
}

// FromEnv reads the passphrase from the given environment variable.
// Returns ErrEmpty if the variable is unset or blank.
func FromEnv(key string) (string, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return "", ErrEmpty
	}
	return v, nil
}
