// Package alias manages short named aliases that map to environment chain names.
//
// Aliases allow users to reference chains by memorable short names instead of
// full project-derived names. For example, "prod" can resolve to
// "myapp-production" so that commands like:
//
//	envchain exec prod -- make deploy
//
// work without typing the full chain name. Aliases are stored in a JSON file
// and are normalised to lowercase before storage and lookup.
package alias
