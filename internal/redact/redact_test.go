package redact_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/env"
	"github.com/envchain/envchain-cli/internal/redact"
)

func TestIsSensitiveMatchesDefaultPatterns(t *testing.T) {
	r := redact.New(nil, "")
	sensitive := []string{"DB_PASSWORD", "GITHUB_TOKEN", "AWS_SECRET_ACCESS_KEY", "PRIVATE_KEY", "API_KEY"}
	for _, k := range sensitive {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitiveSafeKeys(t *testing.T) {
	r := redact.New(nil, "")
	safe := []string{"HOME", "PATH", "USER", "PORT", "DEBUG"}
	for _, k := range safe {
		if r.IsSensitive(k) {
			t.Errorf("expected %q to be safe", k)
		}
	}
}

func TestIsSensitiveCaseInsensitive(t *testing.T) {
	r := redact.New(nil, "")
	if !r.IsSensitive("db_password") {
		t.Error("expected lowercase key to be detected as sensitive")
	}
}

func TestApplyMasksSensitiveEntries(t *testing.T) {
	r := redact.New(nil, "[REDACTED]")
	input := []env.Entry{
		{Key: "HOME", Value: "/home/user"},
		{Key: "API_KEY", Value: "supersecret"},
		{Key: "PORT", Value: "8080"},
	}
	out := r.Apply(input)
	if out[0].Value != "/home/user" {
		t.Errorf("HOME should not be masked, got %q", out[0].Value)
	}
	if out[1].Value != "[REDACTED]" {
		t.Errorf("API_KEY should be masked, got %q", out[1].Value)
	}
	if out[2].Value != "8080" {
		t.Errorf("PORT should not be masked, got %q", out[2].Value)
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	r := redact.New(nil, "")
	input := []env.Entry{{Key: "TOKEN", Value: "original"}}
	r.Apply(input)
	if input[0].Value != "original" {
		t.Error("Apply must not mutate the original slice")
	}
}

func TestApplyMapMasksSensitiveKeys(t *testing.T) {
	r := redact.New(nil, "***")
	m := map[string]string{
		"USER":     "alice",
		"PASSWORD": "hunter2",
	}
	out := r.ApplyMap(m)
	if out["USER"] != "alice" {
		t.Errorf("USER should be unchanged, got %q", out["USER"])
	}
	if out["PASSWORD"] != "***" {
		t.Errorf("PASSWORD should be masked, got %q", out["PASSWORD"])
	}
}

func TestCustomPatterns(t *testing.T) {
	r := redact.New([]string{"CUSTOM"}, "X")
	if !r.IsSensitive("MY_CUSTOM_VAR") {
		t.Error("expected custom pattern to match")
	}
	if r.IsSensitive("API_KEY") {
		t.Error("default patterns should not apply when custom patterns provided")
	}
}
