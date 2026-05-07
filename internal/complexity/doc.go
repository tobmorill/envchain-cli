// Package complexity analyses an env chain and produces a human-readable
// complexity score.
//
// The score is derived from:
//   - the number of entries (base logarithmic penalty)
//   - individual entry findings such as oversized values, shell-interpolation
//     characters, or double-underscore namespace separators
//
// Scores below 3.0 are classified as "low", 3.0–8.0 as "medium", and
// anything above 8.0 as "high".
//
// # Classification thresholds
//
//	Score < 3.0  → Low    (minimal risk, straightforward chain)
//	Score < 8.0  → Medium (moderate complexity, review recommended)
//	Score ≥ 8.0  → High   (significant complexity, refactoring advised)
//
// # Usage
//
// Typical usage involves calling [Analyse] with a slice of env entries and
// inspecting the returned [Report], which exposes the numeric Score and the
// string Classification:
//
//	report := complexity.Analyse(entries)
//	fmt.Printf("complexity: %s (%.2f)\n", report.Classification, report.Score)
package complexity
