package shell

import (
	"os"
	"path/filepath"
	"strings"
)

// Detect attempts to determine the current user's shell from the SHELL
// environment variable. Returns the detected ShellType and true on success,
// or Bash and false if detection fails or the shell is unsupported.
func Detect() (ShellType, bool) {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return Bash, false
	}

	base := strings.ToLower(filepath.Base(shellPath))

	switch {
	case strings.HasPrefix(base, "bash"):
		return Bash, true
	case strings.HasPrefix(base, "zsh"):
		return Zsh, true
	case strings.HasPrefix(base, "fish"):
		return Fish, true
	default:
		return Bash, false
	}
}

// MustDetect returns the detected shell, falling back to Bash.
func MustDetect() ShellType {
	sh, _ := Detect()
	return sh
}
