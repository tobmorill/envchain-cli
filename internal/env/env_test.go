package env

import (
	"errors"
	"testing"
)

func TestParseValid(t *testing.T) {
	cases := []struct {
		input    string
		wantKey  string
		wantVal  string
	}{
		{"FOO=bar", "FOO", "bar"},
		{"KEY=value with spaces", "KEY", "value with spaces"},
		{"A=", "A", ""},
		{"DB_URL=postgres://user:pass@host/db", "DB_URL", "postgres://user:pass@host/db"},
		{"MULTI=a=b=c", "MULTI", "a=b=c"},
	}
	for _, tc := range cases {
		e, err := Parse(tc.input)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if e.Key != tc.wantKey || e.Value != tc.wantVal {
			t.Errorf("Parse(%q) = {%q, %q}, want {%q, %q}", tc.input, e.Key, e.Value, tc.wantKey, tc.wantVal)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	cases := []string{
		"",
		"NOEQUALS",
		"=nokey",
	}
	for _, s := range cases {
		_, err := Parse(s)
		if err == nil {
			t.Errorf("Parse(%q) expected error, got nil", s)
		}
		if !errors.Is(err, ErrInvalidEntry) {
			t.Errorf("Parse(%q) error should wrap ErrInvalidEntry, got %v", s, err)
		}
	}
}

func TestEntryString(t *testing.T) {
	e := Entry{Key: "FOO", Value: "bar"}
	if got := e.String(); got != "FOO=bar" {
		t.Errorf("String() = %q, want %q", got, "FOO=bar")
	}
}

func TestParseAll(t *testing.T) {
	lines := []string{"FOO=bar", "BAZ=qux"}
	m, err := ParseAll(lines)
	if err != nil {
		t.Fatalf("ParseAll unexpected error: %v", err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("ParseAll result mismatch: %v", m)
	}
}

func TestParseAllDuplicateKey(t *testing.T) {
	lines := []string{"FOO=bar", "FOO=baz"}
	_, err := ParseAll(lines)
	if err == nil {
		t.Fatal("ParseAll expected error for duplicate key, got nil")
	}
	if !errors.Is(err, ErrInvalidEntry) {
		t.Errorf("expected ErrInvalidEntry, got %v", err)
	}
}

func TestExportScript(t *testing.T) {
	vars := map[string]string{"MY_VAR": "hello world"}
	script := ExportScript(vars)
	expected := `export MY_VAR="hello world"` + "\n"
	if script != expected {
		t.Errorf("ExportScript = %q, want %q", script, expected)
	}
}
