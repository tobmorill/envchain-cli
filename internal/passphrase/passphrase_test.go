package passphrase_test

import (
	"os"
	"testing"

	"github.com/nicholasgasior/envchain-cli/internal/passphrase"
)

func TestFromEnvReturnsValue(t *testing.T) {
	t.Setenv("TEST_PASSPHRASE", "s3cr3t")
	got, err := passphrase.FromEnv("TEST_PASSPHRASE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "s3cr3t" {
		t.Errorf("expected %q, got %q", "s3cr3t", got)
	}
}

func TestFromEnvEmptyVar(t *testing.T) {
	os.Unsetenv("TEST_PASSPHRASE_EMPTY")
	_, err := passphrase.FromEnv("TEST_PASSPHRASE_EMPTY")
	if err != passphrase.ErrEmpty {
		t.Errorf("expected ErrEmpty, got %v", err)
	}
}

func TestFromEnvBlankVar(t *testing.T) {
	t.Setenv("TEST_PASSPHRASE_BLANK", "   ")
	_, err := passphrase.FromEnv("TEST_PASSPHRASE_BLANK")
	if err != passphrase.ErrEmpty {
		t.Errorf("expected ErrEmpty for blank value, got %v", err)
	}
}

func TestErrMismatchIsDistinct(t *testing.T) {
	if passphrase.ErrMismatch == passphrase.ErrEmpty {
		t.Error("ErrMismatch and ErrEmpty should be distinct errors")
	}
}

func TestFromEnvTrimsWhitespace(t *testing.T) {
	t.Setenv("TEST_PASSPHRASE_TRIM", "  mypass  ")
	got, err := passphrase.FromEnv("TEST_PASSPHRASE_TRIM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "mypass" {
		t.Errorf("expected trimmed value %q, got %q", "mypass", got)
	}
}
