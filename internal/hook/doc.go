// Package hook provides lifecycle hook management for envchain projects.
//
// Hooks allow users to attach shell commands to a project's pre-load or
// post-load phase. When envchain loads a chain, any registered hooks are
// emitted into the generated shell script so the user's shell executes them
// at the appropriate moment.
//
// Hooks are stored as small JSON files on disk, one file per (project, phase)
// pair, inside the hooks directory (default: ~/.config/envchain/hooks/).
//
// Supported phases:
//
//	pre  – runs before environment variables are exported
//	post – runs after environment variables are exported
package hook
