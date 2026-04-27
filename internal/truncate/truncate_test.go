package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain-cli/internal/truncate"
)

func TestPreviewShortValue(t *testing.T) {
	s := "hello"
	got := truncate.Preview(s)
	if got != s {
		t.Errorf("expected %q unchanged, got %q", s, got)
	}
}

func TestPreviewLongValue(t *testing.T) {
	s := strings.Repeat("x", 64)
	got := truncate.Preview(s)
	if len([]rune(got)) > truncate.DefaultMaxLen+len([]rune(truncate.DefaultMask)) {
		t.Errorf("preview too long: %d runes", len([]rune(got)))
	}
	if !strings.HasSuffix(got, truncate.DefaultMask) {
		t.Errorf("expected mask suffix, got %q", got)
	}
}

func TestPreviewExactBoundary(t *testing.T) {
	s := strings.Repeat("a", truncate.DefaultMaxLen)
	got := truncate.Preview(s)
	if got != s {
		t.Errorf("expected exact-length value unchanged, got %q", got)
	}
}

func TestValueCustomOptions(t *testing.T) {
	opts := &truncate.Options{MaxLen: 5, Mask: "..."}
	got := truncate.Value("abcdefgh", opts)
	if got != "abcde..." {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestValueNilOptsUsesDefaults(t *testing.T) {
	long := strings.Repeat("z", 100)
	got := truncate.Value(long, nil)
	if !strings.HasSuffix(got, truncate.DefaultMask) {
		t.Errorf("expected default mask, got %q", got)
	}
}

func TestRedactNonEmpty(t *testing.T) {
	got := truncate.Redact("super-secret-value")
	if strings.Contains(got, "super") {
		t.Errorf("redacted value should not contain original content, got %q", got)
	}
	for _, ch := range got {
		if ch != '*' {
			t.Errorf("expected only '*', got %q", got)
			break
		}
	}
}

func TestRedactEmpty(t *testing.T) {
	got := truncate.Redact("")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestRedactAllOption(t *testing.T) {
	opts := &truncate.Options{RedactAll: true}
	got := truncate.Value("mysecret", opts)
	for _, ch := range got {
		if ch != '*' {
			t.Errorf("expected all stars, got %q", got)
			break
		}
	}
}
