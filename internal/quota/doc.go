// Package quota enforces per-project limits on the number of environment
// variable keys and total value bytes stored in a chain.
//
// A Rule describes the upper bounds; Check validates a slice of env.Entry
// values against those bounds and returns a list of Violation descriptors.
//
// Example:
//
//	rule := quota.DefaultRule()
//	rule.MaxKeys = 20
//	violations, err := quota.Check(rule, entries)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, v := range violations {
//		fmt.Println(v)
//	}
package quota
