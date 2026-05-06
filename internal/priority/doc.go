// Package priority provides per-project key priority management for
// envchain-cli. Each environment variable key within a project can be
// assigned a priority level: low, normal (default), or high.
//
// Priority levels are advisory metadata; they do not affect encryption or
// storage but can be consumed by display commands, policy checks, and
// ordering utilities to surface the most important variables first.
//
// Usage:
//
//	st, _ := store.New(dir)
//	m := priority.New(st)
//	_ = m.Set("myproject", "API_KEY", priority.High)
//	lvl, _ := m.Get("myproject", "API_KEY")
//	fmt.Println(lvl) // high
package priority
