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
package complexity
