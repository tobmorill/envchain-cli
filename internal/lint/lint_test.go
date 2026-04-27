package lint_test

import (
	"testing"

	"github.com/user/envchain-cli/internal/env"
	"github.com/user/envchain-cli/internal/lint"
)

func entries(pairs ...string) []env.Entry {
	var out []env.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, env.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestNoFindingsOnCleanEntries(t *testing.T) {
	findings := lint.Check(entries("FOO", "bar", "BAZ", "qux"))
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findings)
	}
}

func TestDuplicateKey(t *testing.T) {
	findings := lint.Check(entries("FOO", "a", "FOO", "b"))
	if !containsSeverity(findings, lint.Error) {
		t.Fatal("expected error finding for duplicate key")
	}
}

func TestKeyWithWhitespace(t *testing.T) {
	findings := lint.Check(entries("FOO BAR", "value"))
	if !containsMessage(findings, "whitespace") {
		t.Fatal("expected finding about whitespace in key")
	}
}

func TestUnexpandedShellVariableDollarBrace(t *testing.T) {
	findings := lint.Check(entries("TOKEN", "${SECRET}"))
	if !containsMessage(findings, "unexpanded") {
		t.Fatal("expected finding about unexpanded variable")
	}
}

func TestUnexpandedShellVariablePlain(t *testing.T) {
	findings := lint.Check(entries("TOKEN", "$SECRET"))
	if !containsMessage(findings, "unexpanded") {
		t.Fatal("expected finding about unexpanded variable")
	}
}

func TestEmptyValue(t *testing.T) {
	findings := lint.Check(entries("EMPTY", ""))
	if !containsMessage(findings, "empty") {
		t.Fatal("expected finding about empty value")
	}
}

func TestSeverityString(t *testing.T) {
	f := lint.Finding{Key: "K", Message: "msg", Severity: lint.Warn}
	s := f.String()
	if s == "" {
		t.Fatal("expected non-empty string representation")
	}
}

// helpers

func containsSeverity(findings []lint.Finding, s lint.Severity) bool {
	for _, f := range findings {
		if f.Severity == s {
			return true
		}
	}
	return false
}

func containsMessage(findings []lint.Finding, substr string) bool {
	for _, f := range findings {
		if len(f.Message) > 0 && contains(f.Message, substr) {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
