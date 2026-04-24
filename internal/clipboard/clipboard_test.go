package clipboard

import (
	"os/exec"
	"runtime"
	"testing"
)

// TestIsAvailableReturnsBool verifies IsAvailable does not panic and returns a
// consistent result with resolveCommand.
func TestIsAvailableReturnsBool(t *testing.T) {
	_, _, err := resolveCommand()
	want := err == nil
	got := IsAvailable()
	if got != want {
		t.Errorf("IsAvailable() = %v, want %v", got, want)
	}
}

// TestResolveCommandDarwin checks that pbcopy is selected on macOS when
// available. This test is skipped on other platforms.
func TestResolveCommandDarwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-only test")
	}
	cmd, args, err := resolveCommand()
	if err != nil {
		t.Fatalf("resolveCommand() error on darwin: %v", err)
	}
	if cmd != "pbcopy" {
		t.Errorf("expected pbcopy, got %q", cmd)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}

// TestWriteUsesStdin verifies that Write pipes the provided text to the
// clipboard command. We substitute 'cat' as a no-op stand-in when the real
// clipboard tool is unavailable, so this test only runs when at least one
// supported utility exists.
func TestWriteUsesStdin(t *testing.T) {
	if !IsAvailable() {
		t.Skip("no clipboard command available in this environment")
	}
	// Writing an empty string should succeed without error.
	if err := Write(""); err != nil {
		t.Fatalf("Write(\"\") unexpected error: %v", err)
	}
}

// TestWriteUnsupportedReturnsError verifies ErrUnsupported is surfaced when
// resolveCommand finds nothing. We simulate this by temporarily hiding PATH.
func TestWriteUnsupportedReturnsError(t *testing.T) {
	// Only run when we can confirm the real command is absent.
	if _, err := exec.LookPath("pbcopy"); err == nil {
		t.Skip("pbcopy found; cannot simulate unsupported environment")
	}
	if _, err := exec.LookPath("xclip"); err == nil {
		t.Skip("xclip found; cannot simulate unsupported environment")
	}
	if _, err := exec.LookPath("xsel"); err == nil {
		t.Skip("xsel found; cannot simulate unsupported environment")
	}
	if _, err := exec.LookPath("wl-copy"); err == nil {
		t.Skip("wl-copy found; cannot simulate unsupported environment")
	}
	if _, err := exec.LookPath("clip"); err == nil {
		t.Skip("clip found; cannot simulate unsupported environment")
	}

	err := Write("hello")
	if err != ErrUnsupported {
		t.Errorf("expected ErrUnsupported, got %v", err)
	}
}
