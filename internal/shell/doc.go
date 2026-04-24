// Package shell provides shell detection and script generation utilities
// for envchain-cli.
//
// It supports exporting environment variable sets as shell-specific scripts
// that can be evaluated (eval) in the user's current shell session.
//
// Supported shells:
//   - bash
//   - zsh
//   - fish
//
// Example usage:
//
//	sh, _ := shell.Detect()
//	shell.ExportScript(sh, entries, os.Stdout)
//
// The output can be piped into eval:
//
//	eval $(envchain load mychain)
package shell
