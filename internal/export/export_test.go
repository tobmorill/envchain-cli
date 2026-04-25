package export_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envchain-cli/internal/env"
	"github.com/user/envchain-cli/internal/export"
)

func TestWriteAndRead(t *testing.T) {
	ex := export.New()
	entries := []env.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}

	var buf bytes.Buffer
	if err := ex.Write(&buf, "myproject", entries); err != nil {
		t.Fatalf("Write: %v", err)
	}

	b, err := ex.Read(&buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if b.Project != "myproject" {
		t.Errorf("project = %q, want %q", b.Project, "myproject")
	}
	if b.Version != export.Format {
		t.Errorf("version = %d, want %d", b.Version, export.Format)
	}
	if len(b.Entries) != 2 {
		t.Fatalf("entries len = %d, want 2", len(b.Entries))
	}
	if b.Entries[0].Key != "FOO" || b.Entries[0].Value != "bar" {
		t.Errorf("unexpected entry[0]: %+v", b.Entries[0])
	}
}

func TestWriteEmptyProjectReturnsError(t *testing.T) {
	ex := export.New()
	var buf bytes.Buffer
	if err := ex.Write(&buf, "", nil); err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestReadUnsupportedVersion(t *testing.T) {
	payload := `{"version":99,"project":"x","exported_at":"2024-01-01T00:00:00Z","entries":[]}`
	ex := export.New()
	_, err := ex.Read(strings.NewReader(payload))
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
}

func TestReadMissingProject(t *testing.T) {
	payload := `{"version":1,"project":"","exported_at":"2024-01-01T00:00:00Z","entries":[]}`
	ex := export.New()
	_, err := ex.Read(strings.NewReader(payload))
	if err == nil {
		t.Fatal("expected error for missing project name")
	}
}

func TestReadInvalidJSON(t *testing.T) {
	ex := export.New()
	_, err := ex.Read(strings.NewReader("not-json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestWriteEmptyEntries(t *testing.T) {
	ex := export.New()
	var buf bytes.Buffer
	if err := ex.Write(&buf, "empty-project", []env.Entry{}); err != nil {
		t.Fatalf("Write: %v", err)
	}
	b, err := ex.Read(&buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(b.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(b.Entries))
	}
}
