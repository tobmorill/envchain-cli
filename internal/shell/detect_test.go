package shell_test

import (
	"os"
	"testing"

	"github.com/user/envchain-cli/internal/shell"
)

func TestDetectBash(t *testing.T) {
	t.Setenv("SHELL", "/bin/bash")
	sh, ok := shell.Detect()
	if !ok {
		t.Error("expected detection to succeed")
	}
	if sh != shell.Bash {
		t.Errorf("expected bash, got %s", sh)
	}
}

func TestDetectZsh(t *testing.T) {
	t.Setenv("SHELL", "/usr/local/bin/zsh")
	sh, ok := shell.Detect()
	if !ok {
		t.Error("expected detection to succeed")
	}
	if sh != shell.Zsh {
		t.Errorf("expected zsh, got %s", sh)
	}
}

func TestDetectFish(t *testing.T) {
	t.Setenv("SHELL", "/usr/bin/fish")
	sh, ok := shell.Detect()
	if !ok {
		t.Error("expected detection to succeed")
	}
	if sh != shell.Fish {
		t.Errorf("expected fish, got %s", sh)
	}
}

func TestDetectUnknown(t *testing.T) {
	t.Setenv("SHELL", "/bin/tcsh")
	sh, ok := shell.Detect()
	if ok {
		t.Error("expected detection to fail for unknown shell")
	}
	if sh != shell.Bash {
		t.Errorf("expected fallback to bash, got %s", sh)
	}
}

func TestDetectMissing(t *testing.T) {
	os.Unsetenv("SHELL")
	sh, ok := shell.Detect()
	if ok {
		t.Error("expected detection to fail when SHELL is unset")
	}
	if sh != shell.Bash {
		t.Errorf("expected fallback to bash, got %s", sh)
	}
}

func TestMustDetect(t *testing.T) {
	t.Setenv("SHELL", "/bin/zsh")
	sh := shell.MustDetect()
	if sh != shell.Zsh {
		t.Errorf("expected zsh, got %s", sh)
	}
}
