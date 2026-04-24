// Package clipboard provides cross-platform clipboard write support for
// envchain-cli.
//
// It is used to copy shell export scripts to the system clipboard so that
// users can quickly paste them into a running terminal session without
// sourcing a file.
//
// Supported backends (tried in order):
//
//	- macOS:   pbcopy
//	- Windows: clip
//	- Linux:   xclip, xsel, wl-copy
//
// If no backend is available, Write returns ErrUnsupported and the caller
// should fall back to printing the script to stdout.
package clipboard
