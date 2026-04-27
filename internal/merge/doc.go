// Package merge provides multi-chain environment variable merging.
//
// When a project depends on several envchain chains (e.g. a shared
// "infra" chain and a project-specific chain), Merge combines them
// into a single flat list of env.Entry values.
//
// Three conflict-resolution strategies are available:
//
//   - StrategyFirst  – keep the value from the earliest chain (default).
//   - StrategyLast   – overwrite with the value from the latest chain.
//   - StrategyError  – return an ErrConflict if any key appears twice.
//
// The Result type also carries an Origins map so callers can report
// which chain contributed each key (useful for debugging and auditing).
package merge
