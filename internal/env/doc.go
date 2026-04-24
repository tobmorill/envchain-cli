// Package env provides utilities for parsing, validating, and formatting
// environment variable entries used by envchain-cli.
//
// Entries follow the standard KEY=VALUE convention. The package supports:
//   - Parsing individual KEY=VALUE strings into Entry structs
//   - Batch parsing with duplicate-key detection
//   - Generating shell export snippets for shell integration
//
// Example usage:
//
//	entry, err := env.Parse("DATABASE_URL=postgres://localhost/mydb")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(entry.Key)   // DATABASE_URL
//	fmt.Println(entry.Value) // postgres://localhost/mydb
package env
