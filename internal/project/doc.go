// Package project provides utilities for detecting and resolving the current
// project context from the filesystem.
//
// A project name is determined by one of two strategies, in order of preference:
//
//  1. A .envchain marker file in the target directory. The file should contain
//     the project name as its first line. This allows explicit association of a
//     directory tree with a named envchain project.
//
//  2. The base name of the directory itself, used as a convenient default when
//     no marker file is present.
//
// Valid project names consist of alphanumeric characters, hyphens, and
// underscores only.
package project
