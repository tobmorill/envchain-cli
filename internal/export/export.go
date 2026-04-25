// Package export provides functionality for exporting and importing
// environment variable chains to and from portable file formats.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/envchain-cli/internal/env"
)

// Format represents the export file format version.
const Format = 1

// Bundle is the serialisable representation of a single chain export.
type Bundle struct {
	Version   int            `json:"version"`
	Project   string         `json:"project"`
	ExportedAt time.Time     `json:"exported_at"`
	Entries   []env.Entry    `json:"entries"`
}

// Exporter writes and reads chain bundles.
type Exporter struct{}

// New returns a new Exporter.
func New() *Exporter {
	return &Exporter{}
}

// Write serialises entries for the given project to w as JSON.
func (e *Exporter) Write(w io.Writer, project string, entries []env.Entry) error {
	if project == "" {
		return fmt.Errorf("export: project name must not be empty")
	}
	b := Bundle{
		Version:    Format,
		Project:    project,
		ExportedAt: time.Now().UTC(),
		Entries:    entries,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("export: encode: %w", err)
	}
	return nil
}

// Read deserialises a Bundle from r, validating the format version.
func (e *Exporter) Read(r io.Reader) (*Bundle, error) {
	var b Bundle
	if err := json.NewDecoder(r).Decode(&b); err != nil {
		return nil, fmt.Errorf("export: decode: %w", err)
	}
	if b.Version != Format {
		return nil, fmt.Errorf("export: unsupported format version %d (want %d)", b.Version, Format)
	}
	if b.Project == "" {
		return nil, fmt.Errorf("export: bundle contains no project name")
	}
	return &b, nil
}
