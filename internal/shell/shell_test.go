package shell_test

import (
	"strings"
	"testing"

	"github.com/user/envchain-cli/internal/env"
	"github.com/user/envchain-cli/internal/shell"
)

var testEntries = []env.Entry{
	{Key: "FOO", Value: "bar"},
	{Key: "DB_URL", Value: "postgres://localhost/mydb"},
}

func TestExportScriptBash(t *testing.T) {
	var sb strings.Builder
	err := shell.ExportScript(shell.Bash, testEntries, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export FOO=") {
		t.Errorf("expected FOO export, got: %s", out)
	}
	if !strings.Contains(out, "export DB_URL=") {
		t.Errorf("expected DB_URL export, got: %s", out)
	}
}

func TestExportScriptZsh(t *testing.T) {
	var sb strings.Builder
	err := shell.ExportScript(shell.Zsh, testEntries, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "export FOO=") {
		t.Error("zsh output should use export syntax")
	}
}

func TestExportScriptFish(t *testing.T) {
	var sb strings.Builder
	err := shell.ExportScript(shell.Fish, testEntries, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "set -x FOO") {
		t.Errorf("expected fish set -x syntax, got: %s", out)
	}
}

func TestExportScriptUnsupported(t *testing.T) {
	var sb strings.Builder
	err := shell.ExportScript(shell.ShellType("powershell"), testEntries, &sb)
	if err == nil {
		t.Error("expected error for unsupported shell")
	}
}

func TestIsSupported(t *testing.T) {
	if !shell.IsSupported("bash") {
		t.Error("bash should be supported")
	}
	if !shell.IsSupported("Zsh") {
		t.Error("Zsh (case-insensitive) should be supported")
	}
	if shell.IsSupported("powershell") {
		t.Error("powershell should not be supported")
	}
}

func TestSupportedShells(t *testing.T) {
	shells := shell.SupportedShells()
	if len(shells) != 3 {
		t.Errorf("expected 3 supported shells, got %d", len(shells))
	}
}
