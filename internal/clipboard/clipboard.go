// Package clipboard provides utilities for copying environment variable
// export scripts to the system clipboard for quick shell integration.
package clipboard

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

// ErrUnsupported is returned when no clipboard command is available on the
// current platform.
var ErrUnsupported = errors.New("clipboard: no supported clipboard command found")

// Write copies text to the system clipboard. It tries platform-specific
// commands in order and returns ErrUnsupported if none are available.
func Write(text string) error {
	cmd, args, err := resolveCommand()
	if err != nil {
		return err
	}

	c := exec.Command(cmd, args...)
	c.Stdin = strings.NewReader(text)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New("clipboard: " + string(out))
	}
	return nil
}

// IsAvailable reports whether a clipboard command is available on the current
// platform without attempting to write anything.
func IsAvailable() bool {
	_, _, err := resolveCommand()
	return err == nil
}

// resolveCommand returns the clipboard command and arguments for the current
// platform. It checks for common utilities in order of preference.
func resolveCommand() (string, []string, error) {
	var candidates []struct {
		cmd  string
		args []string
	}

	switch runtime.GOOS {
	case "darwin":
		candidates = []struct {
			cmd  string
			args []string
		}{
			{"pbcopy", nil},
		}
	case "windows":
		candidates = []struct {
			cmd  string
			args []string
		}{
			{"clip", nil},
		}
	default:
		// Linux / BSD — try xclip then xsel then wl-copy (Wayland)
		candidates = []struct {
			cmd  string
			args []string
		}{
			{"xclip", []string{"-selection", "clipboard"}},
			{"xsel", []string{"--clipboard", "--input"}},
			{"wl-copy", nil},
		}
	}

	for _, c := range candidates {
		if _, err := exec.LookPath(c.cmd); err == nil {
			return c.cmd, c.args, nil
		}
	}
	return "", nil, ErrUnsupported
}
